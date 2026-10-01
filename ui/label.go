package ui

import (
	"gtic"
	"image"
	"image/color"
	"log"

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
	t := l.RegisterProperty("text", s)
	scale := l.RegisterProperty("scale", 1)
	bg := l.RegisterProperty("bg.color", colornames.Green)
	fg := l.RegisterProperty("fg.color", colornames.Yellow)
	rect := l.Property("rect")
	hAlign := l.RegisterProperty("horizontal.align", AlignCenter)
	vAlign := l.RegisterProperty("vertical.align", AlignCenter)

	vAlg, hAlg := AlignCenter, AlignCenter
	textWidth := 0
	l.OnInit = func(a *gtic.API) {
		t.OnChange.Connect(func(v any) {
			msg := v.(string)
			textWidth = a.Print(msg, 0, -100)
			log.Println("label text len", textWidth)
		})
		hAlign.OnChange.Connect(func(v any) {
			hAlg = v.(Align)
		})
		vAlign.OnChange.Connect(func(v any) {
			vAlg = v.(Align)
		})
		textWidth = a.Print(s, 0, -100)
	}
	l.OnDraw = func(a *gtic.API) {
		if l.Property("hidden").Get().(bool) {
			return
		}
		r := rect.Get().(image.Rectangle)
		x, y, w, h := gtic.RectXYWH(r)
		padX, padY := 0, 0
		switch vAlg {
		case AlignStart:
			padY = 0
		case AlignCenter:
			padY = (h - 8) / 2
		case AlignEnd:
			padY = h - 8
		}
		switch hAlg {
		case AlignStart:
			padX = 0
		case AlignCenter:
			padX = (w - textWidth) / 2
		case AlignEnd:
			padX = w - textWidth
		}

		a.Rect(x, y, w, h, gtic.NewColor(bg.Get().(color.Color)))
		col := gtic.NewColor(fg.Get().(color.Color))
		sc := scale.Get().(int)
		a.Print(t.Get().(string), x+padX, y+padY, col, false, sc, false)
	}
	return l
}
