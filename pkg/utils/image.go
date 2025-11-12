package utils

import (
	"image"
	"image/color"
)

func ImageEqual(img1, img2 image.Image) bool {
	if img1 == nil || img2 == nil {
		return img1 == img2
	}

	if !img1.Bounds().Eq(img2.Bounds()) {
		return false
	}

	bounds := img1.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if !colorsEqual(img1.At(x, y), img2.At(x, y)) {
				return false
			}
		}
	}

	return true
}

func colorsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
