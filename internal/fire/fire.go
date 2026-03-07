package fire

import (
	"math"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	heatLevels = 20
	palette    = " .'`^\",:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
)

type flameZone struct {
	centerX int
	height  float64
	phase   float64
}

type tickMsg time.Time

type Model struct {
	width  int
	height int
	buffer [][]int
	zones  []flameZone
	frame  int
	rng    *rand.Rand
}

func NewModel() Model {
	return Model{
		width:  80,
		height: 24,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) initFire() Model {
	m.buffer = make([][]int, m.height)
	for y := 0; y < m.height; y++ {
		m.buffer[y] = make([]int, m.width)
	}

	zoneCount := 8
	m.zones = nil
	segW := float64(m.width) / float64(zoneCount)
	for i := 0; i < zoneCount; i++ {
		cx := int(segW*float64(i) + segW/2)
		hRatio := 0.6 + m.rng.Float64()*0.6
		m.zones = append(m.zones, flameZone{
			centerX: cx,
			height:  hRatio,
			phase:   m.rng.Float64() * math.Pi * 2,
		})
	}

	return m
}

func (m *Model) heightAt(x int) float64 {
	if len(m.zones) == 0 {
		return 0.8
	}

	var totalW, totalH float64
	for _, z := range m.zones {
		dist := math.Abs(float64(x - z.centerX))
		segW := float64(m.width) / float64(len(m.zones))
		w := math.Exp(-dist * dist / (segW * segW * 0.8))
		totalW += w
		totalH += w * z.height
	}
	if totalW == 0 {
		return 0.8
	}
	return totalH / totalW
}

func (m *Model) windAt(x int) float64 {
	if len(m.zones) == 0 {
		return 0
	}
	frameF := float64(m.frame)
	var totalW, totalWind float64
	segW := float64(m.width) / float64(len(m.zones))
	for _, z := range m.zones {
		dist := math.Abs(float64(x - z.centerX))
		w := math.Exp(-dist * dist / (segW * segW))
		totalW += w
		totalWind += w * math.Sin(frameF*0.05+z.phase) * 0.8
	}
	if totalW == 0 {
		return 0
	}
	return totalWind / totalW
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
		m = m.initFire()

	case tickMsg:
		_ = msg
		m.frame++
		m.updateFire()
		return m, tick()
	}

	return m, nil
}

func (m *Model) updateFire() {
	if m.height == 0 || m.width == 0 || len(m.buffer) == 0 {
		return
	}

	maxHeat := heatLevels - 1
	for x := 0; x < m.width; x++ {
		base := maxHeat - m.rng.Intn(2)
		m.buffer[m.height-1][x] = base
		if m.height >= 2 {
			b2 := base - m.rng.Intn(2)
			if b2 < 0 {
				b2 = 0
			}
			m.buffer[m.height-2][x] = b2
		}
	}

	for y := 0; y < m.height-2; y++ {
		for x := 0; x < m.width; x++ {
			srcY := y + 1 + m.rng.Intn(2)
			if srcY >= m.height {
				srcY = m.height - 1
			}

			wb := m.windAt(x)
			windR := m.rng.Float64()*2.4 - 1.2 + wb
			wind := int(math.Round(windR))
			srcX := x + wind
			if srcX < 0 {
				srcX = 0
			}
			if srcX >= m.width {
				srcX = m.width - 1
			}

			decay := m.rng.Intn(4)
			heat := m.buffer[srcY][srcX] - decay
			if heat < 0 {
				heat = 0
			}

			m.buffer[y][x] = heat
		}
	}
}

func heatToChar(heat int) byte {
	if heat <= 0 {
		return ' '
	}
	idx := heat * (len(palette) - 1) / (heatLevels - 1)
	if idx >= len(palette) {
		idx = len(palette) - 1
	}
	return palette[idx]
}

func (m Model) View() string {
	if len(m.buffer) == 0 {
		return ""
	}

	screen := make([][]byte, m.height)
	for y := 0; y < m.height; y++ {
		screen[y] = make([]byte, m.width)
		for x := range screen[y] {
			screen[y][x] = ' '
		}
	}

	for x := 0; x < m.width; x++ {
		hRatio := m.heightAt(x)
		visibleH := int(float64(m.height) * hRatio)
		if visibleH > m.height {
			visibleH = m.height
		}
		cutoff := m.height - visibleH

		for y := cutoff; y < m.height; y++ {
			heat := m.buffer[y][x]
			if y-cutoff < 3 {
				fade := 3 - (y - cutoff)
				heat -= fade * 3
				if heat < 0 {
					heat = 0
				}
			}
			ch := heatToChar(heat)
			if ch != ' ' {
				screen[y][x] = ch
			}
		}
	}

	buf := make([]byte, 0, m.height*(m.width+1))
	for y := 0; y < m.height; y++ {
		buf = append(buf, screen[y]...)
		if y < m.height-1 {
			buf = append(buf, '\n')
		}
	}

	return string(buf)
}
