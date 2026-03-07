package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hamsurang/tui/internal/tui"
)

func Setup(mode tui.SetupMode) {
	p := tea.NewProgram(tui.NewModel(mode), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
