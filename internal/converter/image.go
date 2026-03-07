package converter

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/eliukblau/pixterm/pkg/ansimage"
	_ "golang.org/x/image/webp"
)

func binarizeImage(src image.Image, threshold uint8) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			if a>>8 < 128 {
				dst.Set(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}
			lum := uint8((299*(r>>8) + 587*(g>>8) + 114*(b>>8)) / 1000)
			if lum > threshold {
				dst.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				dst.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return dst
}

func nearestNeighborResize(src image.Image, dstW, dstH int) *image.RGBA {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			srcX := srcBounds.Min.X + x*srcW/dstW
			srcY := srcBounds.Min.Y + y*srcH/dstH
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

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

	// ansimage uses ▄ (half-block) so each row = 2 pixel rows
	targetH := height * 2
	fitW, fitH := fitSize(src.Bounds().Dx(), src.Bounds().Dy(), width, targetH)

	resized := nearestNeighborResize(src, fitW, fitH)
	bw := binarizeImage(resized, 128)

	img, err := ansimage.NewFromImage(
		bw,
		color.Black,
		ansimage.NoDithering,
	)
	if err != nil {
		return "", err
	}
	return img.Render(), nil
}
