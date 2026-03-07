package converter

import (
	"image"

	"github.com/disintegration/imaging"
)

func Pixelate(src image.Image, targetWidth int) image.Image {
	if targetWidth < 1 {
		targetWidth = 1
	}
	return imaging.Resize(src, targetWidth, 0, imaging.NearestNeighbor)
}

func PixelateWithHeight(src image.Image, targetWidth, targetHeight int) image.Image {
	if targetWidth < 1 {
		targetWidth = 1
	}
	if targetHeight < 1 {
		targetHeight = 1
	}
	return imaging.Resize(src, targetWidth, targetHeight, imaging.NearestNeighbor)
}
