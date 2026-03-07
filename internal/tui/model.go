package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hamsurang/tui/internal/config"
	"github.com/hamsurang/tui/internal/converter"
)

type step int

const (
	stepInput step = iota
	stepPreview
	stepSaved
)

type SetupMode int

const (
	ModeNormal SetupMode = iota
	ModeInit
	ModeSet
)

type Model struct {
	mode      SetupMode
	imagePath string
	preview   string
	cursor    int
	step      step
	err       error
	width     int
	height    int
}

func NewModel(mode SetupMode) Model {
	return Model{
		mode:   mode,
		step:   stepInput,
		width:  80,
		height: 24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Paste && m.step == stepInput {
			m.imagePath += string(msg.Runes)
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.step == stepInput && m.imagePath != "" {
				// Save config
				cfg, err := config.Load()
				if err != nil {
					cfg = &config.Config{Width: 80, Height: 20, PixelWidth: 60}
				}
				cfg.ImagePath = m.imagePath
				if err := config.Save(cfg); err != nil {
					m.err = err
					return m, nil
				}

				if m.mode == ModeInit {
					if err := config.UpdateZshrc(); err != nil {
						m.err = err
						return m, nil
					}
				}

				rendered, err := converter.ImageToANSI(m.imagePath, m.width, 20)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.preview = rendered
				m.step = stepPreview
				m.err = nil
			}

		case "backspace":
			if m.step == stepInput && len(m.imagePath) > 0 {
				runes := []rune(m.imagePath)
				m.imagePath = string(runes[:len(runes)-1])
			}

		case "esc":
			if m.step == stepPreview {
				m.step = stepInput
				m.preview = ""
			}

		default:
			if m.step == stepInput && len(msg.Runes) > 0 {
				m.imagePath += string(msg.Runes)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {
	switch m.step {
	case stepInput:
		s := TitleStyle.Render("tui-theme setup")
		s += "\n\n"
		s += fmt.Sprintf("Image path: %s|", m.imagePath)
		if m.err != nil {
			s += fmt.Sprintf("\n\nError: %v", m.err)
		}
		s += "\n\n" + HelpStyle.Render("enter: preview / q: quit")
		return BorderStyle.Render(s)

	case stepPreview:
		s := TitleStyle.Render("Preview")
		s += "\n\n"
		s += m.preview
		s += "\n\n"

		if m.mode == ModeInit {
			s += "Configuration saved and ~/.zshrc updated!\n"
		} else if m.mode == ModeSet {
			s += "Image configuration updated successfully!\n"
		} else {
			s += "Configuration saved.\n"
		}

		s += "\n" + HelpStyle.Render("esc: back / q: quit")
		return s

	case stepSaved:
		s := TitleStyle.Render("tui-theme setup")
		s += "\n\n"
		s += "Image configuration updated successfully!\n"
		s += "\n" + HelpStyle.Render("q: quit")
		return BorderStyle.Render(s)

	default:
		return "Done!"
	}
}
