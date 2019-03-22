package app

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

type Field struct {
	name   *Text
	value  *Text
	x, y   int32
	valOff int32
}

func NewField(nFont *Font, renderer *sdl.Renderer) *Field {
	f := new(Field)

	f.name = NewText(nFont, renderer)
	f.value = NewText(nFont, renderer)
	return f
}

func (f *Field) Destroy() {
	f.name.Destroy()
	f.value.Destroy()
}

func (f *Field) SetPosition(x, y int32) {
	f.x = x
	f.y = y
}

func (f *Field) SetName(n string) {
	f.name.SetText(n, sdl.Color{R: 200, G: 200, B: 200, A: 255})
	f.valOff = f.name.bounds.W + 5
}

func (f *Field) SetValue(n string) {
	f.value.SetText(n, sdl.Color{R: 200, G: 200, B: 200, A: 255})
}

func (f *Field) Value() string {
	return f.name.text
}

func (f *Field) ValueAsFloat() float64 {
	f64, _ := strconv.ParseFloat(f.name.text, 64)
	return f64
}

func (f *Field) ValueAsInt() int {
	i64, _ := strconv.ParseInt(f.name.text, 10, 64)
	return int(i64)
}

func (f *Field) Draw() {
	f.name.DrawAt(f.x, f.y)
	f.value.DrawAt(f.x+f.valOff, f.y)
}

func (f *Field) DrawAt(x, y int32) {
	f.SetPosition(x, y)
	f.Draw()
}
