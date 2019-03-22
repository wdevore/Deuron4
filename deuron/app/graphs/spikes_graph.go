package graphs

import (
	"image"
	"image/color"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"

	"github.com/fogleman/gg"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron4/simulation/samples"
)

// SpikesAccessor provides access to series data
// type SpikesAccessor func() (x, y float64, c color.Color, done bool)

// SpikesGraph renders spikes as pixels.
// Each horizontals line is a timeline of spikes for a single
// synapse.
type SpikesGraph struct {
	BaseGraph

	dc     *gg.Context
	pixels *image.RGBA

	originX float64
	originY float64

	// Spike accessor state variables
	laneOffset  float64
	colr        color.RGBA
	state       int
	poisLaneY   float64
	poisLane    *samples.SamplesLane
	poisIt      sll.Iterator
	poisScanIdx int

	stimLaneY   float64
	stimLane    *samples.SamplesLane
	stimIt      sll.Iterator
	stimScanIdx int

	noiseColor    color.RGBA
	stimulusColor color.RGBA
	unknownColor  color.RGBA

	// Synapse accessor state vars
	exciteColor  color.RGBA
	inhibitColor color.RGBA
}

// NewSpikesGraph renders spikes
func NewSpikesGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) IGraph {
	g := new(SpikesGraph)
	g.SetGraphics(renderer, texture)
	g.MarkDirty(true)
	g.rect.W = int32(width)
	g.rect.H = int32(height)

	g.dc = gg.NewContext(width, height)
	g.pixels = g.dc.Image().(*image.RGBA)

	g.originX = 10.0
	g.originY = 10.0

	g.laneOffset = 2.0
	g.state = 0
	g.exciteColor = color.RGBA{127, 127, 127, 255}
	g.inhibitColor = color.RGBA{127, 127, 255, 255}

	g.unknownColor = color.RGBA{255, 0, 0, 255}
	g.noiseColor = color.RGBA{255, 127, 0, 255}
	g.stimulusColor = color.RGBA{127, 255, 127, 255}

	return g
}

// SetSeries set the chart data
func (g *SpikesGraph) SetSeries(accessor SeriesAccessor) {
	// g.accessor = accessor

	g.MarkDirty(true)
}

// Destroy release resources
func (g *SpikesGraph) Destroy() {
}

// DrawAt renders graph to texture
func (g *SpikesGraph) DrawAt(x, y int32) {
	g.rect.X = x
	g.rect.Y = y

	if g.dirty {
		g.dc.SetRGB(0.2, 0.2, 0.2)
		g.dc.Clear()
		// Render the graph onto the image pixels

		g.dc.Identity()
		g.dc.InvertY()

		g.dc.Translate(g.originX, g.originY)

		g.dc.SetLineWidth(1.0)

		// Left border
		g.dc.SetRGB(0.85, 0.85, 0.85)
		g.dc.MoveTo(-1.0, 0.0)
		g.dc.LineTo(-1.0, float64(g.rect.H)-g.originY*2)
		g.dc.Stroke()
		// Bottom border
		g.dc.MoveTo(0.0, -1.0)
		g.dc.LineTo(float64(g.rect.W)-g.originX*2, -1.0)
		g.dc.Stroke()

		// Draw colored horizontal bars based on the synapse type.

		// Draw noise spikes.
		px, py, c, more := g.poissonAccessor()
		for more > 0 {
			if more == 1 {
				g.dc.SetColor(c)
				// g.dc.DrawPoint(px, py, 1.0)
				tx, ty := g.dc.TransformPoint(px, py)
				g.dc.SetPixel(int(tx), int(ty)+1)
				g.dc.Fill()
			}

			px, py, c, more = g.poissonAccessor()
		}

		// Draw stimulus spikes.
		sx, sy, sc, smore := g.stimAccessor()
		for smore > 0 {
			if smore == 1 {
				g.dc.SetColor(sc)
				tx, ty := g.dc.TransformPoint(sx, sy)
				g.dc.SetPixel(int(tx), int(ty))
				g.dc.Fill()
			}

			sx, sy, sc, smore = g.stimAccessor()
		}

		g.texture.Update(&g.rect, g.pixels.Pix, g.pixels.Stride)

		// Now copy the texture onto the target (aka the display)
		g.renderer.Copy(g.texture, &g.rect, &g.rect)

		g.MarkDirty(false)
	}

}

func (g *SpikesGraph) Check() bool {
	poiLanes := samples.PoiSamples.GetLanes()
	g.poisIt = poiLanes.Iterator()

	if g.poisIt.First() {
		g.poisLane = g.poisIt.Value().(*samples.SamplesLane)

		stimLanes := samples.StimSamples.GetLanes()
		g.stimIt = stimLanes.Iterator()

		if g.stimIt.First() {
			g.stimLane = g.stimIt.Value().(*samples.SamplesLane)
			return true
		}

		return false
	}

	return false
}

func (g *SpikesGraph) poissonAccessor() (x, y float64, c color.Color, more int) {
	if g.poisScanIdx >= samples.PoiSamples.Size() {
		g.poisScanIdx = 0
		if !g.poisIt.Next() {
			g.poisLaneY = 0
			return 0, 0, nil, 0
		}

		g.poisLaneY += 1 + g.laneOffset
		g.poisLane = g.poisIt.Value().(*samples.SamplesLane)
	}

	spike := g.poisLane.Samples[g.poisScanIdx]

	if spike.Value == 1 {
		g.state = 1
	} else {
		g.state = 2
	}

	g.poisScanIdx++
	return spike.Time, g.poisLaneY, g.noiseColor, g.state
}

func (g *SpikesGraph) stimAccessor() (x, y float64, c color.Color, more int) {
	if g.stimScanIdx >= samples.StimSamples.Size() {
		g.stimScanIdx = 0
		if !g.stimIt.Next() {
			g.stimLaneY = 0
			return 0, 0, nil, 0
		}

		g.stimLaneY += 1 + g.laneOffset
		g.stimLane = g.stimIt.Value().(*samples.SamplesLane)
	}

	spike := g.stimLane.Samples[g.stimScanIdx]
	// fmt.Printf("s: %v\n", spike)

	if spike.Value == 1 {
		g.state = 1
	} else {
		g.state = 2
	}

	g.stimScanIdx++
	return spike.Time, g.stimLaneY, g.stimulusColor, g.state
}
