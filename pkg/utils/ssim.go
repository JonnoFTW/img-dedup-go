package utils

import (
	"image"
	_ "image/png"
)

type SSIMResult struct {
	mssim     float64
	grads     image.Image
	ssimImage image.Image
}

func ssim(im1 image.Image, im2 image.Image) (error, *SSIMResult) {
	// TODO
	out := &SSIMResult{}
	if !im1.Bounds().Eq(im2.Bounds()) {

	}
	return nil, out
}
