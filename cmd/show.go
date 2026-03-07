package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/hamsurang/tui/internal/config"
	"github.com/hamsurang/tui/internal/converter"
)

func Show() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	if cfg.ImagePath == "" {
		fmt.Println("No image configured. Run 'tui-theme pixeltest <image> <width> --save' to set up.")
		return
	}

	name := strings.TrimSuffix(filepath.Base(cfg.ImagePath), filepath.Ext(cfg.ImagePath))

	rendered, err := converter.LoadPixelArt(name)
	if err == nil {
		fmt.Print(rendered)
		return
	}

	pixelWidth := cfg.PixelWidth
	if pixelWidth == 0 {
		pixelWidth = 60
	}

	src, err := imaging.Open(cfg.ImagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Image error: %v\n", err)
		os.Exit(1)
	}

	pixelated := converter.Pixelate(src, pixelWidth)
	rendered = converter.RenderHalfBlocks(pixelated)

	converter.SavePixelArt(rendered, name)

	fmt.Print(rendered)
}
