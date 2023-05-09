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
	Difference HashMethod = iota
)

type ImageHash uint64
type HashGrid [8][8]int

// Adapted from https://towardsdatascience.com/detection-of-duplicate-images-using-image-hash-functions-4d9c53f04a75

// Hash will calculate the hash of an image using the specified method
func Hash(img *image.Image, method HashMethod) ImageHash {
	// based on the specified method type, pick the hash algorithm, run it and return the result
	var hash ImageHash
	switch method {
	case Average:
		hash = averageHash(img)
	case Perceptual:
		hash = perceptualHash(img)
	case Difference:
		hash = differenceHash(img)
	}
	return hash
}

// grayscale - convert an image to grayscale
func grayscale(img *image.Image) *image.Gray {
	grayscaledImage := image.NewGray((*img).Bounds())
	draw.Draw(grayscaledImage, grayscaledImage.Bounds(), (*img), (*img).Bounds().Min, draw.Src)
	return grayscaledImage
}

// resize an image to the speified width and height
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

// averageHash of an image, any pixel above the average pixel value is 1, otherwise 0
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

// differenceHash of an image:
//  reduce image to 8x9,
//  return array where 1 if the pixel to the right is more than the current pixel, otherwise 0
func differenceHash(img *image.Image) ImageHash {
	// https://github.com/JohannesBuchner/imagehash/blob/master/imagehash/__init__.py#L303
	smallGrayImage := resize(grayscale(img), 8, 9)
	hash := ImageHash(0)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if smallGrayImage.GrayAt(i, j).Y > smallGrayImage.GrayAt(i, j+1).Y {
				bitIndex := i*8 + j
				hash += 1 << bitIndex
			}
		}
	}
	return hash
}

func waveletHash(img *image.Image) ImageHash {
	hash := ImageHash(0)

	return hash
}

// imgToFloat convert an image to floating point values between 0 and 1
func imgToFloat(img *image.Gray) [][]float64 {
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
