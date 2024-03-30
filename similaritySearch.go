package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"strings"
	"sync"
	"time"
)

type Histo struct {
	Name string
	H    []int
}
type Similarity struct {
	NameFile string
	rate     float64
}
type Slice struct {
	size, normalsize, rest int
}

func main() {
	fmt.Println()
	args := os.Args
	queryImage := "queryImages//" + args[1]
	images := getImagesFrom(args[2])
	Klist := []int{1, 2, 4, 16, 64, 256, 1048}
	queryHist, _ := computeHistogram(queryImage, 3) //query's Histogram
	/////////////////////////////////////////////
	fmt.Println("WE START THE FULL PROGRAM")
	for _, kValues := range Klist {
		fmt.Println("###############################################################################################################")
		fmt.Println("FOR K=", kValues)
		start := time.Now()
		database := makeImagesForSlice(images, makeSlicesFor(kValues, len(images)))
		fmt.Println("WE START COMPUTING HISTOGRAMS,PLEASE WAIT...................")
		ch := make(chan Histo, len(images))
		var wg sync.WaitGroup
		go func() {
			for index := 0; index < len(database); index++ {
				wg.Add(1)
				go computeHistograms(database[index], 3, ch, &wg)
			}
			wg.Wait() //we wait for all go routine
			close(ch) // we close the channel after inserting all values
		}()
		var similarity []Similarity
		for values := range ch {
			nameS := values.Name
			rarHistQ := compareHistogram(values, queryHist)
			simR := Similarity{nameS, rarHistQ}
			similarity = append(similarity, simR)
		}
		simF := findFiveSimilarity(similarity, 0.1) //here we look for the most similarity image in database
		fmt.Printf("------------------------------------YOUR SIMILARITY IMAGES FOR K=%d----------------------------------------------", kValues)
		fmt.Println()
		for index, simFound := range simF {
			if index == 0 {
				fmt.Printf("The 1st similarity is %s with a rate of %f", simFound.NameFile, simFound.rate)
				fmt.Println()
			} else {
				fmt.Printf("The %dd  similarity is %s with a rate of %f", index+1, simFound.NameFile, simFound.rate)
				fmt.Println()
			}
		}
		fmt.Println("----------------------------------------------------------------------------------------------------------------")
		end := time.Now()
		fmt.Println("EXECUTION TIME IS : ", end.Sub(start))
		fmt.Println()

	}
	fmt.Println("####################################################################################################################")
	fmt.Println("FULL PROGRAM'S DONE")
}
func computeHistogram(imagePath string, depth int) (Histo, error) {
	var Hst []int
	getRGB, errOr := readImage(imagePath, depth)
	if errOr != nil {
		return Histo{"", nil}, errOr
	}
	for index := 0; index < 512; index++ {
		color := find3for(index)
		countColor := 0
		for index2 := 0; index2 < len(getRGB); index2++ {
			if compareBits(getRGB[index2], color) {
				countColor++
			}
		}
		Hst = append(Hst, countColor)
	}
	return Histo{imagePath, Hst}, nil
}

func computeHistograms(imagePath []string, depth int, hChan chan<- Histo, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, value := range imagePath {
		hist, _ := computeHistogram(value, depth)
		hChan <- hist
	}
}
func find3for(number int) [3]int { ///This function help us to find 3 bits for a specific number
	var bits [3]int
	for index := 0; index < 3; index++ {
		if index == 0 {
			listBitsC := combinC()
			for c := 0; c < len(listBitsC); c++ {
				if equal3bits(number, listBitsC[c]) {
					bits = listBitsC[c]
					break
				}
			}
		} else if index == 1 {
			listBitsBC := combinBC()
			for bc := 0; bc < len(listBitsBC); bc++ {
				if equal3bits(number, listBitsBC[bc]) {
					bits = listBitsBC[bc]
					break
				}
			}
		} else {
			listBitsABC := combinABC()
			for abc := 0; abc < len(listBitsABC); abc++ {
				if equal3bits(number, listBitsABC[abc]) {
					bits = listBitsABC[abc]
					break
				}
			}
		}
	}
	return bits

}
func equal3bits(number int, bits [3]int) bool {
	expo2 := 64
	expo1 := 8
	expo0 := 1
	addition := bits[2]*expo0 + bits[1]*expo1 + bits[0]*expo2
	return addition == number
}
func combinC() [][3]int {
	var cComb [][3]int
	a := 0
	b := 0
	for c := 0; c < 8; c++ {
		var newCombinaison = [3]int{a, b, c}
		cComb = append(cComb, newCombinaison)
	}
	return cComb
}
func combinBC() [][3]int {
	var bcComb [][3]int
	a := 0
	var c = []int{0, 1, 2, 3, 4, 5, 6, 7}
	var b = []int{1, 2, 3, 4, 5, 6, 7}
	for index := 0; index < len(b); index++ {
		for index2 := 0; index2 < len(c); index2++ {
			var newCBC = [3]int{a, b[index], c[index2]}
			bcComb = append(bcComb, newCBC)
		}
	}
	return bcComb
}
func combinABC() [][3]int {
	var abcComb [][3]int
	var a = []int{1, 2, 3, 4, 5, 6, 7}
	var b = []int{0, 1, 2, 3, 4, 5, 6, 7}
	var c = []int{0, 1, 2, 3, 4, 5, 6, 7}
	for index := 0; index < len(a); index++ {
		for index1 := 0; index1 < len(b); index1++ {
			for index2 := 0; index2 < len(c); index2++ {
				var abcC = [3]int{a[index], b[index1], c[index2]}
				abcComb = append(abcComb, abcC)
			}
		}
	}
	return abcComb
}
func compareBits(b1 [3]int, b2 [3]int) bool { //this function help us to check if 2 list are the same
	return (b1[0] == b2[0]) && (b1[1] == b2[1]) && (b1[2] == b2[2])
}
func getImagesFrom(nameFile string) []string { //this  function return all images in a dataset
	var files []string
	fileT, err := os.ReadDir(nameFile)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range fileT {
		files = append(files, file.Name())
	}
	return stockImages(files, nameFile)
}
func stockImages(files []string, way string) []string {
	var images []string
	for _, textF := range files {
		if len(textF) < 10 {
			images = append(images, way+"//"+textF)
		}
	}
	return images
}
func readImage(imageName string, dec int) ([][3]int, error) {
	var rgbImage [][3]int
	filesName, errO := os.Open(imageName)
	if errO != nil {
		log.Fatal("Error during open image ", errO)
		return rgbImage, errO
	}
	defer filesName.Close()

	imgT, _, errR := image.Decode(filesName)
	if errR != nil {
		log.Fatal("Error during reading image ", errR)
		return rgbImage, errR
	}
	boundsImage := imgT.Bounds()
	width, height := boundsImage.Max.X, boundsImage.Max.Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := imgT.At(x, y).RGBA()
			r >>= 8
			g >>= 8
			b >>= 8
			red := r >> (8 - dec)
			green := g >> (8 - dec)
			blue := b >> (8 - dec)
			newRGB := [3]int{int(red), int(green), int(blue)}
			rgbImage = append(rgbImage, newRGB)
		}
	}
	return rgbImage, nil
}
func compareHistogram(histoData Histo, histoQuery Histo) float64 { //This func return  a comparaison between two histogram
	var cHist float64
	cHist = 0
	for index := 0; index < len(histoQuery.H); index++ {
		if histoData.H[index] > histoQuery.H[index] {
			cHist = cHist + float64(histoQuery.H[index])
		} else {
			cHist = cHist + float64(histoData.H[index])
		}
	}
	return cHist / 172800
}
func findFiveSimilarity(file []Similarity, rateAppr float64) []Similarity {
	var imagesSimilar []Similarity
	for _, filec := range file {
		if filec.rate > rateAppr {
			imagesSimilar = append(imagesSimilar, filec)
		}
	}
	return classRes(imagesSimilar)
}
func classRes(sim []Similarity) []Similarity {
	var res []Similarity
	for index := 0; index < 5; index++ {
		max := sim[0]
		for _, val := range sim {
			if val.rate > max.rate {
				max = val
			}
		}
		nameF := strings.Split(max.NameFile, "//")
		newC := Similarity{nameF[1], max.rate}
		res = append(res, newC)
		sim = removeElement(sim, max)

	}
	return res
}
func removeElement(simS []Similarity, element Similarity) []Similarity {
	var result []Similarity
	for _, values := range simS {
		if values.NameFile != element.NameFile {
			result = append(result, values)
		}
	}
	return result
}
func makeImagesForSlice(filesContent []string, slice Slice) [][]string { //This func split a file in slice
	var divSlics [][]string
	progress := 0
	for index := 0; index < slice.size; index++ {
		var newOne []string
		for index1 := 0; index1 < slice.normalsize; index1++ {
			newOne = append(newOne, filesContent[progress])
			progress++
		}
		divSlics = append(divSlics, newOne)
	}
	for index2 := 0; index2 < slice.rest; index2++ {
		divSlics[index2] = append(divSlics[index2], filesContent[progress])
		progress++
	}
	return divSlics
}

func makeSlicesFor(k int, lengthData int) Slice { //this func make a slice
	normal := int(lengthData / k)
	addSize := 0
	for index := 0; index < k; index++ {
		addSize = addSize + normal
	}
	restSize := lengthData - addSize
	return Slice{k, normal, restSize}
}
