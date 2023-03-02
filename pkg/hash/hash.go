package hash

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"ssim/pkg/dct"
)

type HashMethod int

const (
	Average    HashMethod = iota
	Perceptual HashMethod = iota
)

type ImageHash uint64
type HashGrid [8][8]int

// Adapted from https://towardsdatascience.com/detection-of-duplicate-images-using-image-hash-functions-4d9c53f04a75

func Hash(img *image.Image, method HashMethod) ImageHash {
	// Grayscale
	// Normalize pixel values
	// Resize the image to 8x8, return that as a string
	//grayscaleImage := image.NewGray(img.Bounds())
	//draw.Draw(grayscaleImage, grayscaleImage.Bounds(), img, img.Bounds().Min, draw.Src)
	//
	//smallImage := image.NewGray(image.Rect(0, 0, 8, 8))
	//draw.BiLinear.Scale(smallImage, smallImage.Rect, grayscaleImage, grayscaleImage.Bounds(), draw.Over, nil)

	// based on the specified method type, pick the hash algorithm, and then return it
	var hash ImageHash
	switch method {
	case Average:
		hash = averageHash(img)
	case Perceptual:
		hash = perceptualHash(img)
	}
	return hash
}
func grayscale(img *image.Image) *image.Gray {
	grayscaledImage := image.NewGray((*img).Bounds())
	draw.Draw(grayscaledImage, grayscaledImage.Bounds(), (*img), (*img).Bounds().Min, draw.Src)
	return grayscaledImage
}
func resize(img *image.Gray, width int, height int) *image.Gray {
	resizedImg := image.NewGray(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return resizedImg
}

func gridToHash(grid *HashGrid) ImageHash {
	out := uint64(0)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			bitIndex := i*8 + j
			if grid[i][j] == 1 {
				out = 1 << bitIndex
			}
		}
	}
	return ImageHash(out)
}
func printGrid(grid *HashGrid) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if grid[i][j] == 1 {
				fmt.Print(" ")
			} else {
				fmt.Print("█")
			}
		}
		fmt.Println()
	}
}
func averageHash(img *image.Image) ImageHash {
	smallGrayImg := resize(grayscale(img), 8, 8)
	hash := ImageHash(0)
	// take the average of the array, anything above is marked as 1, otherwise 0
	total := float64(0)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			total += float64(smallGrayImg.GrayAt(i, j).Y) / 255.
		}
	}
	average := total / 64.
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			bitIndex := i*8 + j
			if float64(smallGrayImg.GrayAt(i, j).Y)/255. > average {
				hash += 1 << bitIndex
			}
		}
	}
	return hash
}
func perceptualHash(img *image.Image) ImageHash {
	// https://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
	// https://github.com/JohannesBuchner/imagehash/blob/master/imagehash/__init__.py#L260-L280
	// Reduce size to 32x32
	// Remove colour
	// dct = dct2(dct2(dct, row-wise), col-wise)
	// dct = extract 8x8 at top left of dct
	// Compute average over the DCT
	// output
	hashSize := 8
	smallGrayImage := resize(grayscale(img), 32, 32)
	pixels := imgToFloat(smallGrayImage)
	freqs := dct.Dct2(dct.Dct2(pixels, 0), 1)
	lowFreqs := make([][]float64, hashSize)
	for i := 0; i < hashSize; i++ {
		lowFreqs[i] = freqs[i][:hashSize]
	}
	total := float64(0)
	for i := 0; i < hashSize; i++ {
		for j := 0; j < hashSize; j++ {
			total += lowFreqs[i][j]
		}
	}
	average := total / 64.
	hash := ImageHash(0)
	for i := 0; i < hashSize; i++ {
		for j := 0; j < hashSize; j++ {
			bitIndex := i*hashSize + j
			if lowFreqs[i][j] > average {
				hash += 1 << bitIndex
			}
		}
	}
	return hash
}

func imgToFloat(img *image.Gray) [][]float64 {
	out := make([][]float64, 0)
	xMax := img.Bounds().Max.X
	yMax := img.Bounds().Max.Y
	for i := 0; i < xMax; i++ {
		out = append(out, make([]float64, yMax))
		for j := 0; j < yMax; j++ {
			out[i][j] = float64(img.GrayAt(i, j).Y) / 255.
		}
	}
	return out
}
