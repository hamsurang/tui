package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// Bypass the cache logic we had before, since users might have changed the target height dynamically.
	// rendered, err := converter.LoadPixelArt(name)
	// if err == nil {
	// 	fmt.Print(rendered)
	// 	return
	// }

	targetWidth := cfg.PixelWidth
	if targetWidth <= 0 {
		targetWidth = 40
	}

	rendered, err := converter.ImageToANSIByWidth(cfg.ImagePath, targetWidth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Image error: %v\n", err)
		os.Exit(1)
	}

	converter.SavePixelArt(rendered, fmt.Sprintf("%s_w%d", name, targetWidth))

	fmt.Print(rendered)
}
