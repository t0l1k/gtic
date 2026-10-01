package ui

import (
	"gtic"
	"gtic/react"
	"image"
	"image/color"

	"golang.org/x/image/colornames"
)

func NewButton(id ElementID, txt string, fn react.SlotFn[*Element]) *Element {
	b := NewElement(id)
	bText := b.RegisterProperty("text", txt)
	normalCol := b.RegisterProperty("normal.color", colornames.Navy)
	pressedCol := b.RegisterProperty("preseed.color", colornames.Red)
	hoverCol := b.RegisterProperty("hover.color", colornames.Blue)
	textCol := b.RegisterProperty("text.color", colornames.Yellow)
	rect := b.Property("rect")
	hidden := b.Property("hidden")
	state := b.Property("state")
	bg := normalCol
	fg := textCol
	textWidth := 0
	b.OnInit = func(a *gtic.API) {
		state.OnChange.Connect(func(v any) {
			s := v.(ElementState)
			switch s {
			case ElementNormal:
				bg = normalCol
			case ElementHover:
				bg = hoverCol
			case ElementPressed:
				bg = pressedCol
			}
		})
		bText.OnChange.Connect(func(v any) {
			msg := v.(string)
			textWidth = a.Print(msg, 0, -100)
		})
		textWidth = a.Print(txt, 0, -100)
	}
	b.OnUpdate = func(a *gtic.API) {
		if hidden.Get().(bool) || state.Get().(ElementState).IsDisabled() {
			return
		}
		x, y, l, m, r, _, _ := a.Mouse()
		inside := image.Pt(x, y).In(rect.Get().(image.Rectangle))
		switch {
		case !inside && !state.Get().(ElementState).IsNormal():
			state.Set(ElementNormal)
		case inside && (l || m || r) && !state.Get().(ElementState).IsPressed():
			state.Set(ElementPressed)
		case inside && !(l || m || r) && !(state.Get().(ElementState).IsHover() || state.Get().(ElementState).IsPressed()):
			state.Set(ElementHover)
		case inside && !(l || m || r) && state.Get().(ElementState).IsPressed():
			if fn != nil {
				fn(b)
			}
			state.Set(ElementHover)
		}
	}
	b.OnDraw = func(a *gtic.API) {
		if hidden.Get().(bool) {
			return
		}
		r := rect.Get().(image.Rectangle)
		stroke := 1
		off := image.Pt(0, 0)
		if state.Get().(ElementState).IsPressed() {
			off = image.Pt(stroke, stroke)
		}
		r = r.Add(off)
		x0, y0, w0, h0 := gtic.RectXYWH(r)
		a.Rect(x0, y0, w0, h0, gtic.NewColor(bg.Get().(color.Color)))

		padX := (w0 - textWidth) / 2
		padY := (h0 - 8) / 2
		a.Print(bText.Get().(string), x0+padX, y0+padY, gtic.NewColor(fg.Get().(color.Color)))
	}
	return b
}
