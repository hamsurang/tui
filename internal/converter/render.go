package converter

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

func RenderHalfBlocks(img image.Image) string {
	bounds := img.Bounds()
	var sb strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			top := img.At(x, y)
			var bottom color.Color
			if y+1 < bounds.Max.Y {
				bottom = img.At(x, y+1)
			} else {
				bottom = color.Black
			}

			tr, tg, tb, _ := top.RGBA()
			br, bg, bb, _ := bottom.RGBA()

			sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀",
				tr>>8, tg>>8, tb>>8,
				br>>8, bg>>8, bb>>8))
		}
		sb.WriteString("\033[0m\n")
	}

	return sb.String()
}
