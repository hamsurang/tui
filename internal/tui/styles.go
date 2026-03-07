package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF7444")).
			MarginBottom(1)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

var (
	outerColor = lipgloss.Color("#576A8F")
	innerColor = lipgloss.Color("#A2CB8B")
)

var outerDash = []rune{'─', '─', '─', ' ', ' ', ' '}
var innerDash = []rune{'╌', '╌', ' ', '╌', '╌', ' '}

func AnimatedBorder(content string, frame int) string {
	lines := strings.Split(content, "\n")
	contentWidth := 0
	for _, l := range lines {
		w := lipgloss.Width(l)
		if w > contentWidth {
			contentWidth = w
		}
	}

	padH := 2
	padV := 1
	innerW := contentWidth + padH*2
	outerW := innerW + 4

	outerStyle := lipgloss.NewStyle().Foreground(outerColor)
	innerStyle := lipgloss.NewStyle().Foreground(innerColor)

	hLine := func(w int, dash []rune, offset int, style lipgloss.Style) string {
		var b strings.Builder
		for i := 0; i < w; i++ {
			ch := dash[(i+offset)%len(dash)]
			b.WriteString(style.Render(string(ch)))
		}
		return b.String()
	}

	slowFrame := frame / 3
	outerOffset := slowFrame % len(outerDash)
	innerOffset := (slowFrame + 2) % len(innerDash)

	outerTL := outerStyle.Render("╭")
	outerTR := outerStyle.Render("╮")
	outerBL := outerStyle.Render("╰")
	outerBR := outerStyle.Render("╯")
	innerTL := innerStyle.Render("┌")
	innerTR := innerStyle.Render("┐")
	innerBL := innerStyle.Render("└")
	innerBR := innerStyle.Render("┘")

	outerVert := outerStyle.Render("│")
	innerVert := innerStyle.Render("│")

	var result []string

	result = append(result, outerTL+hLine(outerW, outerDash, outerOffset, outerStyle)+outerTR)
	result = append(result, outerVert+" "+innerTL+hLine(innerW, innerDash, innerOffset, innerStyle)+innerTR+" "+outerVert)

	for i := 0; i < padV; i++ {
		result = append(result, outerVert+" "+innerVert+strings.Repeat(" ", innerW)+innerVert+" "+outerVert)
	}

	for _, l := range lines {
		w := lipgloss.Width(l)
		pad := innerW - padH*2 - w
		if pad < 0 {
			pad = 0
		}
		row := outerVert + " " + innerVert + strings.Repeat(" ", padH) + l + strings.Repeat(" ", pad+padH) + innerVert + " " + outerVert
		result = append(result, row)
	}

	for i := 0; i < padV; i++ {
		result = append(result, outerVert+" "+innerVert+strings.Repeat(" ", innerW)+innerVert+" "+outerVert)
	}

	result = append(result, outerVert+" "+innerBL+hLine(innerW, innerDash, innerOffset, innerStyle)+innerBR+" "+outerVert)
	result = append(result, outerBL+hLine(outerW, outerDash, outerOffset, outerStyle)+outerBR)

	return strings.Join(result, "\n")
}
