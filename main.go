package main

import (
	"os"

	"github.com/hamsurang/tui/cmd"
	"github.com/hamsurang/tui/internal/tui"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "show":
			cmd.Show()
		case "--init":
			cmd.Setup(tui.ModeInit)
		case "--set":
			cmd.Setup(tui.ModeSet)
		default:
			cmd.Setup(tui.ModeNormal)
		}
	} else {
		cmd.Setup(tui.ModeNormal)
	}
}
