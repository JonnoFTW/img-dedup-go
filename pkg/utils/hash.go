package utils

import (
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"strings"
)

type HashMethod int

const (
	Average    HashMethod = iota
	Perceptual HashMethod = iota
)

type ImageHash string
type HashGrid [8][8]int

func Hash(img image.Image, method HashMethod) ImageHash {
	// Grayscale
	// Normalize pixel values
	// Resize the image to 8x8, return that as a string
	grayscaleImage := image.NewGray(img.Bounds())
	draw.Draw(grayscaleImage, grayscaleImage.Bounds(), img, img.Bounds().Min, draw.Src)
	smallImage := image.NewGray16(image.Rect(0, 0, 8, 8))
	draw.BiLinear.Scale(smallImage, smallImage.Rect, img, img.Bounds(), draw.Over, nil)
	// based on the specified method type, pick the hash algorithm, and then return it
	var grid *HashGrid
	switch method {
	case Average:
		grid = averageHash(smallImage)
	case Perceptual:
	default:
		grid = perceptualHash(smallImage)
	}
	printGrid(grid)
	return gridToHash(grid)

}
func gridToHash(grid *HashGrid) ImageHash {
	var sb strings.Builder
	sb.Grow(8 * 8)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if grid[i][j] == 1 {
				sb.WriteString("1")
			} else {
				sb.WriteString("0")
			}
		}
	}
	return ImageHash(sb.String())
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
func averageHash(img *image.Gray16) *HashGrid {
	grid := new(HashGrid)
	// take the average of the array, anything above is marked as 1, otherwise 0
	total := uint32(0)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			total += uint32(img.Gray16At(i, j).Y)
		}
	}
	average := total / 64
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if uint32(img.Gray16At(i, j).Y) > average {
				grid[i][j] = 1
			} else {
				grid[i][j] = 0
			}
		}
	}
	return grid
}
func perceptualHash(img *image.Gray16) *HashGrid {
	// TODO
	return nil
}
