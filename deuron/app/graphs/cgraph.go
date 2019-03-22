package graphs

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
	chart "github.com/wcharczuk/go-chart"
)

// CGraph represents a chart on the view
type CGraph struct {
	BaseGraph

	graph chart.Chart

	// Go-Chart ImageWriter
	imgWriter *RGBAWriter
}

// NewCGraph creates a new graph chart
func NewCGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int32) IGraph {
	// var err error

	g := new(CGraph)
	g.renderer = renderer
	g.texture = texture
	g.rect.W = width
	g.rect.H = height
	g.dirty = false

	g.graph = chart.Chart{
		Width:  int(width),
		Height: int(height),
	}

	g.imgWriter = NewRGBAWriter()

	return g
}

// SetSeries set the chart data
func (g *CGraph) SetSeries(accessor SeriesAccessor) {
	// g.graph.Series = []chart.Series{
	// 	chart.ContinuousSeries{
	// 		XValues: x,
	// 		YValues: y,
	// 	},
	// }

	g.dirty = true
}
func (g *CGraph) GetAccessor() SeriesAccessor {
	return nil
}

// func (g *CGraph) SetSeries(x, y []float64) {
// 	g.graph.Series = []chart.Series{
// 		chart.ContinuousSeries{
// 			XValues: x,
// 			YValues: y,
// 		},
// 	}

// 	g.dirty = true
// }

// Destroy release resources
func (g *CGraph) Destroy() {
	// err := g.texture.Destroy()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// DrawAt renders graph to texture
func (g *CGraph) DrawAt(x, y int32) {
	var err error
	g.rect.X = x
	g.rect.Y = y

	if g.dirty {
		// Render the graph into the image writer
		err := g.graph.Render(chart.PNG, g.imgWriter)
		g.dirty = false
		if err != nil {
			panic(err)
		}
	}

	// Get image from writer so we can update the display's texture.
	g.ige, err = g.imgWriter.Image()
	if err != nil {
		g.Destroy()
		log.Fatal(err)
	}

	// rect = sdl.Rect{X: 0, Y: 0, W: int32(ige.Rect.Bounds().Dx()), H: int32(ige.Rect.Bounds().Dy())}
	g.texture.Update(&g.rect, g.ige.Pix, g.ige.Stride)

	// Now copy the texture onto the target (aka the display)
	g.renderer.Copy(g.texture, &g.rect, &g.rect)
}

func (g *CGraph) Check() bool {
	return false
}
