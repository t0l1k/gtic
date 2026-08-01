package etic

import (
	"image/color"
)

func PanelGrid(id ElementID, spacing float32, bg, fg color.Color) *Element {
	g := NewElement(id)
	g.OnDraw = func(t *Console) {
		if bg != color.Transparent {
			t.Rect(0, 0, t.Width, t.Height, bg)
		}
		var i float32
		for i = 0; i < t.Height; i++ {
			t.Line(0, i*spacing, t.Width, i*spacing, fg)
			t.Line(i*spacing, 0, i*spacing, t.Height, fg)
		}
	}
	return g
}
