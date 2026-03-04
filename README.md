# SimilarityImages (Go) — Image Similarity Search using Color Histograms

## Overview
This project finds the **top 5 most similar images** to a query image inside a folder-based dataset.
Similarity is computed using **RGB color histograms** (quantized) and a histogram intersection score.

The program is implemented in **Go** and uses **goroutines + channels** to compute histograms in parallel.

---

## Features
- Returns **Top 5 similar images** for a given query image
- Folder-based image database (dataset)
- Color histogram extraction (RGB quantization)
- Parallel processing with goroutines
- Benchmarks multiple `K` values to split workload

---

## How it works (high level)
1. Read the query image and compute its color histogram.
2. Load all images from the dataset folder.
3. Split dataset into `K` slices.
4. Compute histograms in parallel using goroutines.
5. Compare each dataset histogram with the query histogram.
6. Select the **Top 5 highest similarity scores**.

### Histogram details
- RGB values are **quantized to 3 bits per channel** (depth = 3).
- This gives **8 × 8 × 8 = 512 bins**.
- Similarity uses histogram intersection:
  `score = sum(min(H_query[i], H_data[i])) / NORMALIZATION`

---

## Requirements
- Go 1.18+ (recommended)
- JPEG images (the project currently imports `image/jpeg`)

---

## Project structure
SimilarityImages/
├── similaritySearch.go # main program
├── queryImages/ # query images (input)
├── imageDataset2_15_20/ # dataset images (database)
├── execution.txt # notes / execution example
└── README.md

---

## Usage
From the repository root:

```bash
go run similaritySearch.go <query_image_filename> <dataset_folder>

Example Input : **go run similaritySearch.go q01.jpg imageDataset2_15_20**

---

## Output
The program prints the Top 5 similar images with their similarity score and the execution time.

Example Output :
**The 1st similarity is 1044.jpg with a rate of 0.83
The 2nd similarity is 201.jpg  with a rate of 0.81
...
EXECUTION TIME IS:  2.31s**

---

## Notes / Limitations

Notes / Limitations

1. Works best when dataset images have similar size/format.

2. Current implementation focuses on color similarity (not shapes/objects).


