package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/doyoonlee/tui-theme/internal/config"
	"github.com/doyoonlee/tui-theme/internal/converter"
)

func PixelTest() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: tui-theme pixeltest <image_path> [pixel_width] [--save]")
		fmt.Println("")
		fmt.Println("  pixel_width: number of horizontal pixels (default: 60)")
		fmt.Println("  --save:      save pixel art to ~/.tui/cache/")
		fmt.Println("")
		fmt.Println("Examples:")
		fmt.Println("  tui-theme pixeltest photo.png 40")
		fmt.Println("  tui-theme pixeltest photo.png 40 --save")
		os.Exit(1)
	}

	imagePath := os.Args[2]
	pixelWidth := 60
	save := false

	for _, arg := range os.Args[3:] {
		if arg == "--save" {
			save = true
		} else if w, err := strconv.Atoi(arg); err == nil && w > 0 {
			pixelWidth = w
		}
	}

	src, err := imaging.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open image: %v\n", err)
		os.Exit(1)
	}

	bounds := src.Bounds()
	fmt.Fprintf(os.Stderr, "Original: %dx%d → Pixel art: %d wide\n\n",
		bounds.Dx(), bounds.Dy(), pixelWidth)

	pixelated := converter.Pixelate(src, pixelWidth)
	result := converter.RenderHalfBlocks(pixelated)
	fmt.Print(result)

	if save {
		name := strings.TrimSuffix(filepath.Base(imagePath), filepath.Ext(imagePath))
		savedPath, err := converter.SavePixelArt(result, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to save: %v\n", err)
			os.Exit(1)
		}

		absPath, _ := filepath.Abs(imagePath)
		cfg, _ := config.Load()
		cfg.ImagePath = absPath
		cfg.PixelWidth = pixelWidth
		config.Save(cfg)

		fmt.Fprintf(os.Stderr, "\nSaved to %s\n", savedPath)
		fmt.Fprintf(os.Stderr, "Config updated (pixel_width: %d)\n", pixelWidth)
	}
}
