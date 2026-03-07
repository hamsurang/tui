package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/hamsurang/tui/internal/config"
	"github.com/hamsurang/tui/internal/converter"
	"github.com/hamsurang/tui/internal/donut"
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
	mode             SetupMode
	imagePath        string
	preview          string
	cursor           int
	pixelWidthInput  string
	customPixelWidth int
	step             step
	err              error
	width            int
	height           int
	donut            donut.Model
	frame            int
}

func NewModel(mode SetupMode) Model {
	return Model{
		mode:   mode,
		step:   stepInput,
		width:  80,
		height: 24,
		donut:  donut.NewModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.donut.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Paste && m.step == stepInput {
			m.imagePath += string(msg.Runes)
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "enter":
				if m.step == stepInput && m.imagePath != "" {
					m.step = stepHeightInput
					m.pixelWidthInput = ""
					return m, nil
				}

				if m.step == stepHeightInput {
					// Parse the custom input, fallback to 60
					w, err := strconv.Atoi(m.pixelWidthInput)
					if err != nil || w <= 0 {
						w = 60
					}

					// Apply reasonable limits
					if w < 8 {
						w = 8
					} else if w > m.width {
						w = m.width
					}

					m.customPixelWidth = w
					m.updatePreview()
					if m.err == nil {
						m.step = stepPreview
					}
					return m, nil
				}

				if m.step == stepPreview {
					cfg, err := config.Load()
					if err != nil {
						cfg = &config.Config{Width: 80, Height: 20, PixelWidth: 60}
					}
					cfg.ImagePath = m.imagePath
					cfg.PixelWidth = m.customPixelWidth
					if err := config.Save(cfg); err != nil {
						m.err = err
						return m, nil
					}
				}

			case "backspace":
				if m.step == stepInput && len(m.imagePath) > 0 {
					runes := []rune(m.imagePath)
					m.imagePath = string(runes[:len(runes)-1])
				} else if m.step == stepHeightInput && len(m.pixelWidthInput) > 0 {
					runes := []rune(m.pixelWidthInput)
					m.pixelWidthInput = string(runes[:len(runes)-1])
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
					if msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
						m.pixelWidthInput += string(msg.Runes)
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var donutModel tea.Model
	var cmd tea.Cmd
	donutModel, cmd = m.donut.Update(msg)
	m.donut = donutModel.(donut.Model)
	if cmd != nil {
		m.frame++
	}
	return m, cmd
}

func (m Model) overlayOnDonut(panel string) string {
	bg := m.donut.View()

	panelHeight := lipgloss.Height(panel)
	panelWidth := lipgloss.Width(panel)

	padTop := (m.height - panelHeight) / 2
	if padTop < 0 {
		padTop = 0
	}
	padLeft := (m.width - panelWidth) / 2
	if padLeft < 0 {
		padLeft = 0
	}

	bgLines := strings.Split(bg, "\n")
	for len(bgLines) < m.height {
		bgLines = append(bgLines, strings.Repeat(" ", m.width))
	}

	panelLines := strings.Split(panel, "\n")
	for i, pLine := range panelLines {
		row := padTop + i
		if row >= len(bgLines) {
			break
		}
		left := ansi.Truncate(bgLines[row], padLeft, "")
		right := ansi.TruncateLeft(bgLines[row], padLeft+panelWidth, "")
		bgLines[row] = left + pLine + right
	}

	if len(bgLines) > m.height {
		bgLines = bgLines[:m.height]
	}
	return strings.Join(bgLines, "\n")
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
		panel := AnimatedBorder(s, m.frame)
		return m.overlayOnDonut(panel)

	case stepHeightInput:
		s := TitleStyle.Render("Set Pixel Width")
		s += "\n\n"
		s += "Enter pixel width (e.g. 40 for Small, 60 for Medium, 80 for Large):\n"
		s += fmt.Sprintf("> %s|", m.pixelWidthInput)
		s += "\n\n(Leave empty to use default: 60)\n"
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
		s += fmt.Sprintf("Pixel Width: %d", m.customPixelWidth)

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
	rendered, err := converter.ImageToANSIByPixelWidth(m.imagePath, m.customPixelWidth)
	if err != nil {
		m.err = err
		return
	}
	m.preview = rendered
	m.err = nil
}
