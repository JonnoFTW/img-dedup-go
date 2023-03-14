package hash

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
	// https://github.com/scikit-image/scikit-image/blob/main/skimage/metrics/_structural_similarity.py
	out := &SSIMResult{}
	if !im1.Bounds().Eq(im2.Bounds()) {

	}
	return nil, out
}
