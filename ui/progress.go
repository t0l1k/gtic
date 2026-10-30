package ui

import (
	"fmt"
	"gtic"
	"image"
	"image/color"

	"golang.org/x/image/colornames"
)

func NewProgress(id ElementID, initial, min, max float32, o Orientation) *Element {
	p := NewElement(id)
	value := p.RegisterProperty("value", initial)
	orient := p.RegisterProperty("orientation", o)
	rect := p.Property("rect")
	label := NewLabel("progress.label", fmt.Sprintf("%.1v", value.Get()))
	label.Property("rect").Set(rect.Get().(image.Rectangle))
	label.Property("bg").Set(gtic.Transparent)
	label.Property("fg").Set(gtic.NewColor(colornames.Yellow))
	rect.OnChange.Connect(func(a any) { label.Property("rect").Set(a.(image.Rectangle)) })

	value.OnChange.Connect(func(a any) { label.Property("text").Set(fmt.Sprintf("%.1v", a.(float32))) })

	newXorColor := func(c color.Color) color.Color {
		r, g, b, a := c.RGBA()
		col := gtic.NewRGBA(uint8(r^0xFF), uint8(g^0xFF), uint8(b^0xFF), uint8(a))
		return col
	}

	bg := p.RegisterProperty("bg", gtic.NewRGBA(0, 0, 255, 255))
	bg.OnChange.Connect(func(a any) { label.Property("fg").Set(newXorColor(a.(color.Color))) })

	ratio := func(v, min, max float32) float32 {
		rangeVal := max - min
		if rangeVal == 0 {
			rangeVal = 1
		}
		ratio := (v - min) / rangeVal
		result := gtic.Clamp(ratio, min, max)
		return result
	}

	p.OnDraw = func(t *gtic.API) {
		colBg := bg.Get().(color.Color)
		colFg := newXorColor(colBg)
		x, y, w, h := gtic.RectXYWH(rect.Get().(image.Rectangle))
		t.Rect(x, y, w, h, gtic.NewColor(colBg))

		switch orient.Get() {
		case Vertical:
			fillH := int(float32(h) * ratio(value.Get().(float32), min, max))
			fillY := y + h - fillH
			t.Rect(x, fillY, w, fillH, gtic.NewColor(colFg))

			// На жёлтой заливке текст синий.
			t.Clip(x, fillY, w, fillH)
			label.Property("fg").Set(gtic.NewColor(colBg))
			label.Draw(t)

			// На синем фоне текст жёлтый.
			t.Clip(x, y, w, h-fillH)
			label.Property("fg").Set(gtic.NewColor(colFg))
			label.Draw(t)

		case Horizontal:
			fillW := int(float32(w) * ratio(value.Get().(float32), min, max))
			t.Rect(x, y, fillW, h, gtic.NewColor(colFg))

			// На жёлтой заливке текст синий.
			t.Clip(x, y, fillW, h)
			label.Property("fg").Set(gtic.NewColor(colBg))
			label.Draw(t)

			// На синем фоне текст жёлтый.
			t.Clip(x+fillW, y, w-fillW, h)
			label.Property("fg").Set(gtic.NewColor(colFg))
			label.Draw(t)

		default:
			label.Draw(t)
		}

		t.Clip() // сбросить отсечение: без аргументов Clip включает весь экран
	}
	return p
}
