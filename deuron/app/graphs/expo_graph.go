package graphs

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/veandco/go-sdl2/sdl"
)

// A test graph for plotting exponentials to see how they shape

type ExpoGraph struct {
	BaseGraph

	dc     *gg.Context
	pixels *image.RGBA

	originX float64
	originY float64

	lineColor color.RGBA

	dt float64

	hw float64
	hh float64

	a    float64
	tau  float64
	m    float64
	s    float64
	wMax float64
	w    float64
}

func NewExpoGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) IGraph {
	g := new(ExpoGraph)
	g.SetGraphics(renderer, texture)
	g.MarkDirty(true)
	g.rect.W = int32(width)
	g.rect.H = int32(height)

	g.dc = gg.NewContext(width, height)
	g.pixels = g.dc.Image().(*image.RGBA)

	g.originX = 10.0
	g.originY = 10.0

	g.lineColor = color.RGBA{255, 127, 0, 255}

	g.dt = 0.0
	g.hw = float64(g.rect.W / 2)
	g.hh = float64(g.rect.H / 2)
	g.a = 100.0
	g.tau = 100.0
	g.m = 5.0
	g.s = g.m / g.tau
	g.wMax = 20.0

	return g
}

func (g *ExpoGraph) SetSeries(accessor SeriesAccessor) {
}

// Destroy release resources
func (g *ExpoGraph) Destroy() {
}

// DrawAt renders graph to texture
func (g *ExpoGraph) DrawAt(x, y int32) {
	g.rect.X = x
	g.rect.Y = y

	if g.dirty {
		g.dc.SetRGB(0.2, 0.2, 0.2)
		g.dc.Clear()
		// Render the graph onto the image pixels

		g.dc.Identity()
		g.dc.InvertY()

		// g.dc.Translate(g.originX, g.originY)

		g.dc.SetLineWidth(1.0)

		// Axis X
		g.dc.SetRGB(0.85, 0.85, 0.85)
		g.dc.MoveTo(0.0, g.hh)
		g.dc.LineTo(float64(g.rect.W), g.hh)
		g.dc.Stroke()
		// Axis Y
		g.dc.MoveTo(g.hw, 0.0)
		g.dc.LineTo(g.hw, float64(g.rect.H))
		g.dc.Stroke()

		// minor vertical grid lines
		for x := g.hw; x < float64(g.rect.W); x += 10.0 {
			g.dc.SetRGB(0.55, 0.55, 0.55)
			if math.Mod(x, 100) == 12.0 {
				g.dc.SetRGB(0.55, 0.55, 1.0)
			}
			g.dc.MoveTo(x, 0.0)
			g.dc.LineTo(x, float64(g.rect.H))
			g.dc.Stroke()
		}
		for x := g.hw; x > 0.0; x -= 10.0 {
			g.dc.SetRGB(0.55, 0.55, 0.55)
			if math.Mod(x, 100) == 12.0 {
				g.dc.SetRGB(0.55, 0.55, 1.0)
			}
			g.dc.MoveTo(x, 0.0)
			g.dc.LineTo(x, float64(g.rect.H))
			g.dc.Stroke()
		}

		px, py, c, more := g.dataAccessor()
		g.dc.MoveTo(px+g.hw, py+g.hh)
		for more > 0 {
			px, py, c, more = g.dataAccessor()
			g.dc.SetColor(c)
			g.dc.LineTo(px+g.hw, py+g.hh)
		}
		g.dc.Stroke()

		g.texture.Update(&g.rect, g.pixels.Pix, g.pixels.Stride)

		// Now copy the texture onto the target (aka the display)
		g.renderer.Copy(g.texture, &g.rect, &g.rect)

		// g.MarkDirty(false)
	}

}

func (g *ExpoGraph) Check() bool {
	return false
}

func (g *ExpoGraph) SetA(v float64) {
	g.a = v
}

func (g *ExpoGraph) A() float64 {
	return g.a
}

func (g *ExpoGraph) SetTau(v float64) {
	g.tau = v
	g.s = g.m / g.tau
}

func (g *ExpoGraph) Tau() float64 {
	return g.tau
}

func (g *ExpoGraph) SetM(v float64) {
	g.m = v
	g.s = g.m / g.tau
}

func (g *ExpoGraph) M() float64 {
	return g.m
}

func (g *ExpoGraph) SetWMax(v float64) {
	g.wMax = v
}

func (g *ExpoGraph) WMax() float64 {
	return g.wMax
}

func (g *ExpoGraph) dataAccessor() (x, y float64, c color.Color, more int) {

	// -dt causes decay, "+" causes growth
	ev := g.a * math.Exp(-g.s*g.dt)
	// adding (1-w) causes a dip and is similar to soft bounds.
	// ev := g.a * (1.0 - g.w/g.wMax) * math.Exp(-g.s*g.dt)

	dt := g.dt
	g.dt = g.dt + 1.0
	g.w += 1.0
	state := 1

	// if g.w > g.wMax {
	// 	g.w = 0.0
	// 	state = 0
	// 	g.dt = 0.0
	// }

	if g.dt > 500 {
		g.w = 0.0
		state = 0
		g.dt = 0.0
	}

	return dt, ev, g.lineColor, state
}

func (g *ExpoGraph) dataAccessor2() (x, y float64, c color.Color, more int) {

	ev := 10 * math.Exp(g.dt/100.0)
	g.dt = g.dt - 1.0
	state := 1
	if g.dt < 0.0 {
		state = 0
		g.dt = 1000.0
	}

	return g.dt, ev, g.lineColor, state
}
