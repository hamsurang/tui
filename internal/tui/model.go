package tui

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hamsurang/tui/internal/config"
	"github.com/hamsurang/tui/internal/converter"
)

type step int

const (
	stepInput step = iota
	stepHeightInput
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
	mode         SetupMode
	imagePath    string
	preview      string
	cursor       int
	heightInput  string
	customHeight int
	step         step
	err          error
	width        int
	height       int
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
				m.step = stepHeightInput
				m.heightInput = ""
				return m, nil
			}

			if m.step == stepHeightInput {
				// Parse the custom input, fallback to 20
				h, err := strconv.Atoi(m.heightInput)
				if err != nil || h <= 0 {
					h = 20
				}

				// Apply reasonable limits
				maxH := (m.height / 2) - 4
				if maxH < 10 {
					maxH = 10
				}
				if h < 5 {
					h = 5
				} else if h > maxH {
					h = maxH
				}

				m.customHeight = h
				m.updatePreview()
				if m.err == nil {
					m.step = stepPreview
				}
				return m, nil
			}

			if m.step == stepPreview {
				// Save config
				cfg, err := config.Load()
				if err != nil {
					cfg = &config.Config{Width: 80, Height: 20, PixelWidth: 60}
				}
				cfg.ImagePath = m.imagePath
				cfg.Height = m.customHeight
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

				m.step = stepSaved
				return m, nil
			}

		case "backspace":
			if m.step == stepInput && len(m.imagePath) > 0 {
				runes := []rune(m.imagePath)
				m.imagePath = string(runes[:len(runes)-1])
			} else if m.step == stepHeightInput && len(m.heightInput) > 0 {
				runes := []rune(m.heightInput)
				m.heightInput = string(runes[:len(runes)-1])
			}

		case "esc":
			if m.step == stepHeightInput {
				m.step = stepInput
			} else if m.step == stepPreview {
				m.step = stepHeightInput
				m.preview = ""
			}

		default:
			if m.step == stepInput && len(msg.Runes) > 0 {
				m.imagePath += string(msg.Runes)
			} else if m.step == stepHeightInput && len(msg.Runes) > 0 {
				// Only accept numbers
				if msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
					m.heightInput += string(msg.Runes)
				}
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
		s += "\n\n" + HelpStyle.Render("enter: next / q: quit")
		return BorderStyle.Render(s)

	case stepHeightInput:
		s := TitleStyle.Render("Set Image Height")
		s += "\n\n"
		s += "Enter desired height (e.g. 20 for Medium, 30 for Large):\n"
		s += fmt.Sprintf("> %s|", m.heightInput)
		s += "\n\n(Leave empty to use default height 20)\n"
		if m.err != nil {
			s += fmt.Sprintf("\n\nError: %v", m.err)
		}
		s += "\n\n" + HelpStyle.Render("enter: preview / esc: back / q: quit")
		return BorderStyle.Render(s)

	case stepPreview:
		s := TitleStyle.Render("Preview")
		s += "\n\n"
		s += m.preview
		s += "\n\n"
		s += fmt.Sprintf("Current Resolution Height: %d", m.customHeight)

		if m.err != nil {
			s += fmt.Sprintf("\n\nError: %v", m.err)
		}

		s += "\n\n" + HelpStyle.Render("enter: confirm & save / esc: back / q: quit")
		return s

	case stepSaved:
		s := TitleStyle.Render("tui-theme setup")
		s += "\n\n"

		if m.mode == ModeInit {
			s += "Configuration saved and ~/.zshrc updated!\n"
		} else if m.mode == ModeSet {
			s += "Image configuration updated successfully!\n"
		} else {
			s += "Configuration saved.\n"
		}

		s += "\n" + HelpStyle.Render("q: quit")
		return BorderStyle.Render(s)

	default:
		return "Done!"
	}
}

func (m *Model) updatePreview() {
	rendered, err := converter.ImageToANSI(m.imagePath, m.width, m.customHeight)
	if err != nil {
		m.err = err
		return
	}
	m.preview = rendered
	m.err = nil
}
