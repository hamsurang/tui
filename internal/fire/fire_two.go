package fire

import (
	"math"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	firetwoPalette = " .'`^\",:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	firetwoLevels  = 20
)

type firetwoTickMsg time.Time

type ModelTwo struct {
	width  int
	height int
	heat   [][]float64
	frame  int
	rng    *rand.Rand
}

func NewModelTwo() ModelTwo {
	return ModelTwo{
		width:  80,
		height: 24,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m ModelTwo) Init() tea.Cmd {
	return firetwoTick()
}

func firetwoTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return firetwoTickMsg(t)
	})
}

func (m ModelTwo) initBuffer() ModelTwo {
	m.heat = make([][]float64, m.height)
	for y := 0; y < m.height; y++ {
		m.heat[y] = make([]float64, m.width)
	}
	return m
}

// simplex-inspired 2D noise using permutation table
type noiseGen struct {
	perm [512]int
}

func newNoiseGen(rng *rand.Rand) *noiseGen {
	ng := &noiseGen{}
	p := make([]int, 256)
	for i := range p {
		p[i] = i
	}
	rng.Shuffle(256, func(i, j int) { p[i], p[j] = p[j], p[i] })
	for i := 0; i < 512; i++ {
		ng.perm[i] = p[i&255]
	}
	return ng
}

func (ng *noiseGen) noise2D(x, y float64) float64 {
	xi := int(math.Floor(x)) & 255
	yi := int(math.Floor(y)) & 255
	xf := x - math.Floor(x)
	yf := y - math.Floor(y)

	u := xf * xf * (3 - 2*xf)
	v := yf * yf * (3 - 2*yf)

	aa := ng.perm[ng.perm[xi]+yi]
	ab := ng.perm[ng.perm[xi]+yi+1]
	ba := ng.perm[ng.perm[xi+1]+yi]
	bb := ng.perm[ng.perm[xi+1]+yi+1]

	g := func(hash int, fx, fy float64) float64 {
		switch hash & 3 {
		case 0:
			return fx + fy
		case 1:
			return -fx + fy
		case 2:
			return fx - fy
		default:
			return -fx - fy
		}
	}

	x1 := lerp(g(aa, xf, yf), g(ba, xf-1, yf), u)
	x2 := lerp(g(ab, xf, yf-1), g(bb, xf-1, yf-1), u)
	return lerp(x1, x2, v)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func (ng *noiseGen) fbm(x, y float64, octaves int) float64 {
	var val, amp float64
	freq := 1.0
	amp = 1.0
	for i := 0; i < octaves; i++ {
		val += ng.noise2D(x*freq, y*freq) * amp
		freq *= 2.0
		amp *= 0.5
	}
	return val
}

var firetwoNoise *noiseGen

func (m *ModelTwo) updateFire() {
	if m.height == 0 || m.width == 0 || len(m.heat) == 0 {
		return
	}

	if firetwoNoise == nil {
		firetwoNoise = newNoiseGen(m.rng)
	}

	frameF := float64(m.frame)
	maxHeat := float64(firetwoLevels - 1)

	// noise-based heat source at the bottom 2 rows
	for x := 0; x < m.width; x++ {
		nx := float64(x) * 0.08
		nt := frameF * 0.04
		n := firetwoNoise.fbm(nx, nt, 3)
		// n is roughly in [-1.5, 1.5], normalize to [0, 1]
		n = (n + 1.5) / 3.0
		if n > 1 {
			n = 1
		}
		if n < 0 {
			n = 0
		}

		// create blob-like heat: use threshold to make chunks
		var base float64
		if n > 0.35 {
			base = maxHeat * (0.7 + 0.3*n)
		} else {
			base = maxHeat * n * 0.5
		}
		base += (m.rng.Float64() - 0.5) * 2
		if base > maxHeat {
			base = maxHeat
		}
		if base < 0 {
			base = 0
		}

		m.heat[m.height-1][x] = base
		if m.height >= 2 {
			m.heat[m.height-2][x] = base * (0.85 + m.rng.Float64()*0.15)
		}
	}

	// propagate upward with 3x3 weighted blur for chunky feel
	tmp := make([][]float64, m.height)
	for y := 0; y < m.height; y++ {
		tmp[y] = make([]float64, m.width)
		copy(tmp[y], m.heat[y])
	}

	for y := 0; y < m.height-2; y++ {
		for x := 0; x < m.width; x++ {
			// weighted average from below (3x3 kernel biased downward)
			var total, wSum float64
			for dy := 1; dy <= 3; dy++ {
				sy := y + dy
				if sy >= m.height {
					continue
				}
				for dx := -1; dx <= 1; dx++ {
					sx := x + dx
					if sx < 0 || sx >= m.width {
						continue
					}
					w := 1.0
					if dx == 0 {
						w = 2.0
					}
					if dy == 1 {
						w *= 2.0
					}
					total += tmp[sy][sx] * w
					wSum += w
				}
			}

			avg := total / wSum

			// small noise-based perturbation for organic shape
			nx := float64(x)*0.12 + frameF*0.02
			ny := float64(y)*0.15 - frameF*0.06
			nv := firetwoNoise.noise2D(nx, ny)

			// wind drift
			windNoise := firetwoNoise.noise2D(float64(x)*0.05, frameF*0.03)
			drift := windNoise * 1.2
			driftInt := int(math.Round(drift))
			if driftInt != 0 {
				shiftX := x + driftInt
				if shiftX >= 0 && shiftX < m.width {
					sy := y + 1
					if sy < m.height {
						avg = avg*0.6 + tmp[sy][shiftX]*0.4
					}
				}
			}

			// decay: less decay = chunkier blobs
			decay := 0.6 + m.rng.Float64()*0.8 + nv*0.3
			heat := avg - decay
			if heat < 0 {
				heat = 0
			}

			m.heat[y][x] = heat
		}
	}
}

func (m ModelTwo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m = m.initBuffer()

	case firetwoTickMsg:
		_ = msg
		m.frame++
		m.updateFire()
		return m, firetwoTick()
	}

	return m, nil
}

func firetwoHeatToChar(heat float64) byte {
	if heat < 1.0 {
		return ' '
	}
	idx := int(heat) * (len(firetwoPalette) - 1) / (firetwoLevels - 1)
	if idx >= len(firetwoPalette) {
		idx = len(firetwoPalette) - 1
	}
	if idx < 0 {
		idx = 0
	}
	return firetwoPalette[idx]
}

func (m ModelTwo) View() string {
	if len(m.heat) == 0 {
		return ""
	}

	buf := make([]byte, 0, m.height*(m.width+1))
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			buf = append(buf, firetwoHeatToChar(m.heat[y][x]))
		}
		if y < m.height-1 {
			buf = append(buf, '\n')
		}
	}

	return string(buf)
}
