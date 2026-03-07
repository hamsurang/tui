package converter

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/webp"
)

// Removed binarizeImage and nearestNeighborResize as they're no longer used in colorful rendering.

func fitSize(srcW, srcH, maxW, maxH int) (int, int) {
	ratioW := float64(maxW) / float64(srcW)
	ratioH := float64(maxH) / float64(srcH)
	ratio := ratioW
	if ratioH < ratioW {
		ratio = ratioH
	}
	w := int(float64(srcW) * ratio)
	h := int(float64(srcH) * ratio)
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return w, h
}

func ImageToANSI(path string, width, height int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	targetH := height * 2
	fitW, fitH := fitSize(src.Bounds().Dx(), src.Bounds().Dy(), width, targetH)

	pixelated := PixelateWithHeight(src, fitW, fitH)
	return RenderHalfBlocks(pixelated), nil
}

func ImageToANSIByWidth(path string, targetWidth int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}

	pixelated := Pixelate(src, targetWidth)
	return RenderHalfBlocks(pixelated), nil
}
