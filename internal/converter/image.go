package converter

import (
	"image/color"

	"github.com/eliukblau/pixterm/pkg/ansimage"
)

func ImageToANSI(path string, width, height int) (string, error) {
	img, err := ansimage.NewScaledFromFile(
		path, height, width,
		color.Black,
		ansimage.ScaleModeFit,
		ansimage.NoDithering,
	)
	if err != nil {
		return "", err
	}
	return img.Render(), nil
}
