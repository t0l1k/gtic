package ui

import (
	"gtic"
	"image"
	"image/color"

	"golang.org/x/image/colornames"
)

type Align int

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
)

func NewLabel(id ElementID, s string) *Element {
	l := NewElement(id)
	txt := l.RegisterProperty("text", s)
	scale := l.RegisterProperty("scale", 1)
	bg := l.RegisterProperty("bg", colornames.Green)
	fg := l.RegisterProperty("fg", colornames.Yellow)
	rect := l.Property("rect")
	hAlign := l.RegisterProperty("halign", AlignCenter)
	vAlign := l.RegisterProperty("valign", AlignCenter)
	textWidth := 0
	l.OnInit = func(a *gtic.API) {
		txt.OnChange.Connect(func(v any) {
			msg := v.(string)
			textWidth = a.Print(msg, 0, -100)
		})
		textWidth = a.Print(s, 0, -100)
	}
	l.OnDraw = func(a *gtic.API) {
		if l.Property("hidden").Get().(bool) {
			return
		}
		colBg := bg.Get().(color.Color)
		colFg := fg.Get().(color.Color)
		sc := scale.Get().(int)
		r := rect.Get().(image.Rectangle)
		x, y, w, h := gtic.RectXYWH(r)
		padX, padY := 0, 0
		switch vAlign.Get().(Align) {
		case AlignStart:
			padY = 0
		case AlignCenter:
			padY = (h - 8*sc) / 2
		case AlignEnd:
			padY = h - 8*sc
		}
		switch hAlign.Get().(Align) {
		case AlignStart:
			padX = 0
		case AlignCenter:
			padX = (w - textWidth) / 2
		case AlignEnd:
			padX = w - textWidth
		}

		a.Rect(x, y, w, h, gtic.NewColor(colBg))
		a.Print(txt.Get().(string), x+padX, y+padY, gtic.NewColor(colFg), false, sc, false)
	}
	return l
}
