package main

import (
	"os"

	"github.com/doyoonlee/tui-theme/cmd"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "show" {
		cmd.Show()
	} else {
		cmd.Setup()
	}
}
