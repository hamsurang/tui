package donut

import (
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	thetaSpacing   = 0.07
	phiSpacing     = 0.02
	r1             = 1.0
	r2             = 2.0
	k2             = 5.0
	luminanceChars = ".,-~:;=!*#$@"
	trailLength    = 15
)

type donutInstance struct {
	a, b   float64
	speedA float64
	speedB float64
	posX   float64
	posY   float64
	velX   float64
	velY   float64
	size   float64
	color  lipgloss.Color
	delay  int
	lane   int
	active bool
	trail  []trailPoint
}

type trailPoint struct {
	x, y float64
}

type tickMsg time.Time

type Model struct {
	width  int
	height int
	donuts []donutInstance
	frame  int
}

func NewModel() Model {
	donuts := []donutInstance{
		{a: 0, b: 0, speedA: 0.04, speedB: 0.02, size: 0.30, color: lipgloss.Color("#84B179"), velX: -1.8, velY: 1.2, delay: 0, lane: 0},
		{a: 1.0, b: 0.5, speedA: 0.07, speedB: 0.03, size: 0.21, color: lipgloss.Color("#A2CB8B"), velX: -1.8, velY: 1.2, delay: 25, lane: 1},
		{a: 2.0, b: 1.0, speedA: 0.03, speedB: 0.05, size: 0.24, color: lipgloss.Color("#576A8F"), velX: -1.8, velY: 1.2, delay: 50, lane: 2},
		{a: 0.5, b: 2.0, speedA: 0.06, speedB: 0.01, size: 0.18, color: lipgloss.Color("#C7EABB"), velX: -1.8, velY: 1.2, delay: 75, lane: 3},
		{a: 3.0, b: 1.5, speedA: 0.02, speedB: 0.06, size: 0.225, color: lipgloss.Color("#E8F5BD"), velX: -1.8, velY: 1.2, delay: 100, lane: 4},
		{a: 1.5, b: 0.3, speedA: 0.05, speedB: 0.04, size: 0.20, color: lipgloss.Color("#FF7444"), velX: -1.8, velY: 1.2, delay: 125, lane: 5},
		{a: 0.8, b: 1.8, speedA: 0.03, speedB: 0.06, size: 0.26, color: lipgloss.Color("#84B179"), velX: -1.8, velY: 1.2, delay: 150, lane: 6},
		{a: 2.5, b: 0.7, speedA: 0.06, speedB: 0.02, size: 0.19, color: lipgloss.Color("#A2CB8B"), velX: -1.8, velY: 1.2, delay: 175, lane: 7},
	}
	return Model{
		width:  80,
		height: 24,
		donuts: donuts,
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.spawnDonuts()

	case tickMsg:
		_ = msg
		m.frame++
		for i := range m.donuts {
			d := &m.donuts[i]

			if !d.active && m.frame >= d.delay {
				d.active = true
				m = m.spawnAt(i)
			}

			if !d.active {
				continue
			}

			d.a += d.speedA
			d.b += d.speedB

			d.trail = append(d.trail, trailPoint{x: d.posX, y: d.posY})
			if len(d.trail) > trailLength {
				d.trail = d.trail[len(d.trail)-trailLength:]
			}

			d.posX += d.velX
			d.posY += d.velY

			if d.posX < -20 || d.posY > float64(m.height)+20 {
				m = m.spawnAt(i)
			}
		}
		return m, tick()
	}

	return m, nil
}

func (m Model) spawnDonuts() Model {
	for i := range m.donuts {
		if m.donuts[i].active {
			continue
		}
	}
	return m
}

func (m Model) spawnAt(idx int) Model {
	d := &m.donuts[idx]
	laneSpacing := float64(m.width) / 9.0
	startX := float64(m.width) + 5 + laneSpacing*float64(d.lane)*0.3
	startY := -10.0 - float64(d.lane)*float64(m.height)*0.25
	d.posX = startX
	d.posY = startY
	d.trail = nil
	d.active = true
	return m
}

func (m Model) View() string {
	screen := make([][]rune, m.height)
	zbuffer := make([][]float64, m.height)
	colors := make([][]lipgloss.Color, m.height)

	for y := 0; y < m.height; y++ {
		screen[y] = make([]rune, m.width)
		zbuffer[y] = make([]float64, m.width)
		colors[y] = make([]lipgloss.Color, m.width)
		for x := 0; x < m.width; x++ {
			screen[y][x] = ' '
		}
	}

	for _, d := range m.donuts {
		if !d.active {
			continue
		}
		m.renderTrail(d, screen, colors)
		m.renderDonut(d, screen, zbuffer, colors)
	}

	var result string
	for y := 0; y < m.height; y++ {
		line := ""
		for x := 0; x < m.width; x++ {
			ch := screen[y][x]
			if ch != ' ' {
				style := lipgloss.NewStyle().Foreground(colors[y][x])
				line += style.Render(string(ch))
			} else {
				line += " "
			}
		}
		if y < m.height-1 {
			result += line + "\n"
		} else {
			result += line
		}
	}

	return result
}

func (m Model) renderTrail(d donutInstance, screen [][]rune, colors [][]lipgloss.Color) {
	if len(d.trail) < 2 {
		return
	}

	trailChars := []rune{'*', '*', '#', '=', '=', ';', ';', ':', ':', '~', '~', '-', '-', '.', '.'}

	totalPoints := len(d.trail)
	for i := 0; i < totalPoints-1; i++ {
		p0 := d.trail[i]
		p1 := d.trail[i+1]

		dx := p1.x - p0.x
		dy := p1.y - p0.y
		dist := math.Sqrt(dx*dx + dy*dy)
		steps := int(dist*2) + 1

		for s := 0; s <= steps; s++ {
			t := float64(s) / float64(steps)
			tx := p0.x + t*dx
			ty := p0.y + t*dy

			px := int(tx)
			py := int(ty)

			if px < 0 || px >= m.width || py < 0 || py >= m.height {
				continue
			}

			if screen[py][px] != ' ' {
				continue
			}

			age := totalPoints - 1 - i
			charIdx := age
			if charIdx >= len(trailChars) {
				charIdx = len(trailChars) - 1
			}
			screen[py][px] = trailChars[charIdx]

			fade := float64(age) / float64(totalPoints)
			dimLevel := 255 - int(fade*19)
			if dimLevel < 236 {
				dimLevel = 236
			}
			colors[py][px] = d.color
			if fade > 0.5 {
				colors[py][px] = lipgloss.Color(itoa(dimLevel))
			}
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

func (m Model) renderDonut(d donutInstance, screen [][]rune, zbuffer [][]float64, colors [][]lipgloss.Color) {
	cosA, sinA := math.Cos(d.a), math.Sin(d.a)
	cosB, sinB := math.Cos(d.b), math.Sin(d.b)

	scaledR1 := r1 * d.size
	scaledR2 := r2 * d.size
	k1 := float64(m.width) * k2 * 3 / (8 * (r1 + r2))

	offsetX := d.posX
	offsetY := d.posY

	for theta := 0.0; theta < 2*math.Pi; theta += thetaSpacing {
		cosTheta, sinTheta := math.Cos(theta), math.Sin(theta)

		for phi := 0.0; phi < 2*math.Pi; phi += phiSpacing {
			cosPhi, sinPhi := math.Cos(phi), math.Sin(phi)

			circleX := scaledR2 + scaledR1*cosTheta
			circleY := scaledR1 * sinTheta

			x := circleX*(cosB*cosPhi+sinA*sinB*sinPhi) - circleY*cosA*sinB
			y := circleX*(sinB*cosPhi-sinA*cosB*sinPhi) + circleY*cosA*cosB
			z := k2 + cosA*circleX*sinPhi + circleY*sinA
			ooz := 1 / z

			xp := int(offsetX + k1*ooz*x)
			yp := int(offsetY + k1*ooz*y*0.5)

			if xp < 0 || xp >= m.width || yp < 0 || yp >= m.height {
				continue
			}

			L := cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi - sinA*sinTheta + cosB*(cosA*sinTheta-cosTheta*sinA*sinPhi)

			if L > 0 && ooz > zbuffer[yp][xp] {
				zbuffer[yp][xp] = ooz
				lumIdx := int(L * 8)
				if lumIdx > len(luminanceChars)-1 {
					lumIdx = len(luminanceChars) - 1
				}
				if lumIdx < 0 {
					lumIdx = 0
				}
				screen[yp][xp] = rune(luminanceChars[lumIdx])
				colors[yp][xp] = d.color
			}
		}
	}
}
