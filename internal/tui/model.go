package tui

import (
	"fmt"
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
	mode          SetupMode
	imagePath     string
	preview       string
	cursor        int
	resolutionIdx int
	step          step
	err           error
	width         int
	height        int
	donut         donut.Model
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

			case "up", "k":
				if m.step == stepPreview && m.resolutionIdx > 0 {
					m.resolutionIdx--
					m.updatePreview()
				}
			case "down", "j":
				if m.step == stepPreview && m.resolutionIdx < 3 {
					m.resolutionIdx++
					m.updatePreview()
				}

			case "enter":
				if m.step == stepInput && m.imagePath != "" {
					m.resolutionIdx = 1
					m.updatePreview()
					if m.err == nil {
						m.step = stepPreview
					}
				} else if m.step == stepPreview {
					cfg, err := config.Load()
					if err != nil {
						cfg = &config.Config{Width: 80, Height: 20, PixelWidth: 60}
					}
					cfg.ImagePath = m.imagePath
					cfg.Height = m.getTargetHeight()
					if err := config.Save(cfg); err != nil {
						m.err = err
					} else if m.mode == ModeInit {
						if err := config.UpdateZshrc(); err != nil {
							m.err = err
						}
					}

					if m.err == nil {
						m.step = stepSaved
					}
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
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var donutModel tea.Model
	var cmd tea.Cmd
	donutModel, cmd = m.donut.Update(msg)
	m.donut = donutModel.(donut.Model)
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
		panel := BorderStyle.Render(s)
		return m.overlayOnDonut(panel)

	case stepPreview:
		s := TitleStyle.Render("Preview")
		s += "\n\n"
		s += m.preview
		s += "\n\n"
		s += "Current Resolution: "

		options := []string{"Small (15)", "Medium (20)", "Large (30)", "Full"}
		var labeledOptions []string
		for i, res := range options {
			if m.resolutionIdx == i {
				labeledOptions = append(labeledOptions, fmt.Sprintf("[ > %s < ]", res))
			} else {
				labeledOptions = append(labeledOptions, res)
			}
		}

		for i, opt := range labeledOptions {
			s += opt
			if i < len(labeledOptions)-1 {
				s += " | "
			}
		}

		s += "\n(Use up/down arrow keys to change size)"

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

func (m *Model) getTargetHeight() int {
	var targetHeight int
	switch m.resolutionIdx {
	case 0:
		targetHeight = 15
	case 1:
		targetHeight = 20
	case 2:
		targetHeight = 30
	case 3:
		targetHeight = (m.height / 2) - 4
		if targetHeight < 10 {
			targetHeight = 10
		}
	}
	return targetHeight
}

func (m *Model) updatePreview() {
	targetHeight := m.getTargetHeight()
	rendered, err := converter.ImageToANSI(m.imagePath, m.width, targetHeight)
	if err != nil {
		m.err = err
		return
	}
	m.preview = rendered
	m.err = nil
}
