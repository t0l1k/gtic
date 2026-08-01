package etic

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
)

type Progress struct {
	Element
	OnChanged Signal[float32]
	orient    Orientation
	Value     Property[float32]
}

func NewProgress(id ElementID, initial, min, max float32, o Orientation) *Progress {
	p := &Progress{
		Element:   *NewElement(id),
		OnChanged: *NewSignal[float32](),
		orient:    o,
		Value:     *NewProperty(initial),
	}
	ratio := func(v, min, max float32) float32 {
		rangeVal := max - min
		if rangeVal == 0 {
			rangeVal = 1
		}
		ratio := (v - min) / rangeVal
		result := Clamp(ratio, min, max)
		return result
	}

	p.OnDraw = func(t *Console) {
		x, y, w, h := RectF32(p.Bounds().Get())
		t.Rect(x, y, w, h, colornames.Blue)
		switch p.orient {
		case Vertical:
			fillH := h * ratio(p.Value.Get(), min, max)
			t.Rect(x, y+h-fillH, w, fillH, colornames.Yellow)
		case Horizontal:
			fillW := w * ratio(p.Value.Get(), min, max)
			t.Rect(x, y, fillW, h, colornames.Yellow)
		}
		t.DrawString(fmt.Sprintf("%.1v", p.Value.Get()), FontSystem, p.Bounds().Get(), text.AlignCenter, text.AlignCenter, colornames.Black)
	}
	return p
}
