package cmd

import (
	"fmt"
	"os"

	"github.com/doyoonlee/tui-theme/internal/config"
	"github.com/doyoonlee/tui-theme/internal/converter"
)

func Show() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(1)
	}

	if cfg.ImagePath == "" {
		fmt.Println("No image configured. Run '.tui' to set up.")
		return
	}

	width := cfg.Width
	if width == 0 {
		width = 80
	}
	height := cfg.Height
	if height == 0 {
		height = 20
	}

	rendered, err := converter.ImageToANSI(cfg.ImagePath, width, height)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Image error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(rendered)
}
