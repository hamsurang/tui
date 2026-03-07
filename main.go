package main

import (
	"os"

	"github.com/hamsurang/tui/cmd"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "show":
			cmd.Show()
		case "pixeltest":
			cmd.PixelTest()
		default:
			cmd.Setup()
		}
	} else {
		cmd.Setup()
	}
}
