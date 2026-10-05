package ui

import (
	"image"
)

type Layout interface {
	Apply(parentRect image.Rectangle, children []IElement)
}

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// Макет настройки макета вручную, вызовом Property("rect").OnChange.Connect(func(r image.Rectangle[float32]) { тут настройка макета})
type AbsoluteLayout struct{}

func (l AbsoluteLayout) Apply(r image.Rectangle, children []IElement) {}

type FlexLayout struct {
	orientation Orientation
	Fractions   []float32 // например: []float64{0.6, 0.4}
	Padding     float32   // например: 0.1 = 10%
}

func NewVerticalFlexLayout(fraction []float32, padding float32) *FlexLayout {
	return &FlexLayout{orientation: Vertical, Fractions: fraction, Padding: padding}
}
func NewHorizontalFlexLayout(fraction []float32, padding float32) *FlexLayout {
	return &FlexLayout{orientation: Horizontal, Fractions: fraction, Padding: padding}
}
func (l FlexLayout) Apply(parent image.Rectangle, children []IElement) {
	if len(children) == 0 {
		return
	}
	var (
		totalH, totalV, x0, y0, w0, h0, x, y, w, h, frac float32
		r                                                image.Rectangle
	)
	x0, y0 = float32(parent.Min.X), float32(parent.Min.Y)
	w0, h0 = float32(parent.Dx()), float32(parent.Dy())
	if l.orientation == Vertical {
		totalH = h0
		y = y0
	} else {
		totalV = w0
		x = x0
	}
	idx := 0
	for _, ch := range children {
		if l.orientation == Vertical {
			frac = l.Fractions[idx]
			idx++
			h = totalH * frac
			// область для ребёнка
			r = image.Rect(int(x0), int(y), int(x0+w0), int(y+h))
		} else {
			frac = l.Fractions[idx]
			idx++
			w = totalV * frac
			// область для ребёнка
			r = image.Rect(int(x), int(y0), int(x+w), int(y0+h0))
		}
		// внутренний padding
		padX := int(float32(r.Dx()) * l.Padding)
		padY := int(float32(r.Dy()) * l.Padding)
		ch.Property("rect").Set(image.Rect(r.Min.X+padX, r.Min.Y+padY, r.Max.X-padX, r.Max.Y-padY))
		if l.orientation == Vertical {
			y += h
		} else {
			x += w
		}
	}
}

type BoxLayout struct {
	orientation Orientation
	Padding     float32
}

func NewVerticalBoxLayout(padding float32) *BoxLayout {
	return &BoxLayout{orientation: Vertical, Padding: padding}
}
func NewHorizontalBoxLayout(padding float32) *BoxLayout {
	return &BoxLayout{orientation: Horizontal, Padding: padding}
}

func (l BoxLayout) Apply(r image.Rectangle, children []IElement) {
	n := len(children)
	if n == 0 {
		return
	}
	frac := 1 / float32(n)
	var fraction []float32
	for range children {
		fraction = append(fraction, frac)
	}
	if l.orientation == Vertical {
		vp := NewVerticalFlexLayout(fraction, l.Padding)
		vp.Apply(r, children)
	} else {
		hp := NewHorizontalFlexLayout(fraction, l.Padding)
		hp.Apply(r, children)
	}
}

type StackLayout struct{ Padding float32 }

func NewStackLayout(spacing float32) *StackLayout { return &StackLayout{Padding: spacing} }
func (b *StackLayout) Apply(parent image.Rectangle, children []IElement) {
	if len(children) == 0 {
		return
	}
	x0, y0 := float32(parent.Min.X), float32(parent.Min.Y)
	w0, h0 := float32(parent.Dx()), float32(parent.Dy())
	padX := float32(w0) * b.Padding
	padY := float32(h0) * b.Padding
	for _, c := range children {
		x := x0 + padX
		y := y0 + padY
		w := w0 - padX*2
		h := h0 - padY*2
		r := image.Rect(int(x), int(y), int(x+w), int(y+h))
		c.Property("rect").Set(r)
	}
	found := false
	for _, c := range children {
		if !c.Property("hidden").Get().(bool) && !found {
			c.Property("hidden").Set(true)
			found = true
			continue
		}
		c.Property("hidden").Set(false)
	}
}

type GridLayout struct {
	Rows, Cols int
	Padding    float32
	square     bool
}

func NewGridLayout(rows, cols int, padding float32) *GridLayout {
	return &GridLayout{Rows: rows, Cols: cols, Padding: padding, square: false}
}
func NewSquareGriodLayout(row, col int, spacing float32) *GridLayout {
	return &GridLayout{Padding: spacing, Rows: row, Cols: col, square: true}
}

func (g *GridLayout) Apply(parent image.Rectangle, children []IElement) {
	if g.Rows <= 0 || g.Cols <= 0 || len(children) == 0 {
		return
	}
	x0, y0 := float32(parent.Min.X), float32(parent.Min.Y)
	w0, h0 := float32(parent.Dx()), float32(parent.Dy())

	cellW := w0 / float32(g.Rows)
	cellH := h0 / float32(g.Cols)
	if g.square {
		cellSize := min(cellH, cellW)
		cellW, cellH = cellSize, cellSize

		gridW := cellW * float32(g.Rows)
		gridH := cellH * float32(g.Cols)
		offsetX := int((w0 - gridW) / 2)
		offsetY := int((h0 - gridH) / 2)
		parent = image.Rect(parent.Min.X+offsetX, parent.Min.Y+offsetY, parent.Max.X-offsetX, parent.Max.Y-offsetY)
	}

	for i, ch := range children {
		if i >= g.Rows*g.Cols {
			break // Сетка заполнена
		}
		row := i % g.Rows
		col := i / g.Rows

		x := int(x0 + float32(row)*cellW)
		y := int(y0 + float32(col)*cellH)

		r := image.Rect(x, y, x+int(cellW), y+int(cellH))

		padX := int(float32(r.Dx()) * g.Padding)
		padY := int(float32(r.Dy()) * g.Padding)

		inner := image.Rect(r.Min.X+padX, r.Min.Y+padY, r.Max.X-padX, r.Max.Y-padY)
		ch.Property("rect").Set(inner)
	}
}
