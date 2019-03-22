package graphs

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

// SeriesAccessor provides access to series data
type SeriesAccessor func() (x, y float64, c color.Color, state int)

var ran = rand.New(rand.NewSource(1963))

type IGraph interface {
	Destroy()
	DrawAt(x, y int32)
	SetSeries(SeriesAccessor)
	// Accessor() SeriesAccessor
	Check() bool
	MarkDirty(dirty bool)
}

type BaseGraph struct {
	Bounds sdl.Rect
	rect   sdl.Rect

	ige *image.RGBA

	renderer *sdl.Renderer
	texture  *sdl.Texture

	dirty bool
}

func (bg *BaseGraph) SetGraphics(renderer *sdl.Renderer, texture *sdl.Texture) {
	bg.renderer = renderer
	bg.texture = texture
}

func (bg *BaseGraph) MarkDirty(dirty bool) {
	bg.dirty = dirty
}
