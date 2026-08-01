package etic

type Layout interface {
	Apply(parentRect Rectangle[float32], children []IElement)
}

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// Макет настройки макета вручную, вызовом Bounds().OnChange.Connect(func(r Rectangle[float32]) { тут настройка макета})
type AbsoluteLayout struct{}

func (l AbsoluteLayout) Apply(r Rectangle[float32], children []IElement) {}

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
func (l FlexLayout) Apply(parent Rectangle[float32], children []IElement) {
	if len(children) == 0 {
		return
	}
	var (
		totalH, totalV, x, y, w, h, frac float32
		r                                Rectangle[float32]
	)
	if l.orientation == Vertical {
		totalH = parent.Dy()
		y = parent.Min.Y
	} else {
		totalV = parent.Dx()
		x = parent.Min.X
	}
	idx := 0
	for _, ch := range children {
		if l.orientation == Vertical {
			frac = l.Fractions[idx]
			idx++
			h = totalH * frac
			// область для ребёнка
			r = Rect(parent.Min.X, y, parent.Max.X, y+h)
		} else {
			frac = l.Fractions[idx]
			idx++
			w = totalV * frac
			// область для ребёнка
			r = Rect(x, parent.Min.Y, x+w, parent.Max.Y)
		}
		// внутренний padding
		padX := r.Dx() * l.Padding
		padY := r.Dy() * l.Padding
		if child, ok := ch.(interface {
			Bounds() *Property[Rectangle[float32]]
		}); ok {
			child.Bounds().Set(Rect(r.Min.X+padX, r.Min.Y+padY, r.Max.X-padX, r.Max.Y-padY))
		}
		if l.orientation == Vertical {
			y += h
		} else {
			x += w
		}
	}
}

type BoxLayout struct {
	orientation Orientation
	Fractions   []float32
	Padding     float32
}

func NewVerticalBoxLayout(padding float32) *BoxLayout {
	return &BoxLayout{orientation: Vertical, Padding: padding}
}
func NewHorizontalBoxLayout(padding float32) *BoxLayout {
	return &BoxLayout{orientation: Horizontal, Padding: padding}
}

func (l BoxLayout) Apply(r Rectangle[float32], children []IElement) {
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

type StackLayout struct{ spacing float32 }

func NewStackLayout(spacing float32) *StackLayout { return &StackLayout{spacing: spacing} }
func (b *StackLayout) Apply(rect Rectangle[float32], children []IElement) {
	if len(children) == 0 {
		return
	}
	for _, c := range children {
		x := rect.Min.X + b.spacing
		y := rect.Min.Y + b.spacing
		w := rect.Dx() - b.spacing*2
		h := rect.Dy() - b.spacing*2
		r := Rect(x, y, x+w, y+h)
		c.Bounds().Set(r)
	}
	found := false
	for _, c := range children {
		if !c.IsHidden() && !found {
			c.Show()
			found = true
			continue
		}
		c.Hide()
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

func (g *GridLayout) Apply(parent Rectangle[float32], children []IElement) {
	if g.Rows <= 0 || g.Cols <= 0 || len(children) == 0 {
		return
	}

	cellW := parent.Dx() / float32(g.Rows)
	cellH := parent.Dy() / float32(g.Cols)
	if g.square {
		cellSize := min(cellH, cellW)
		cellW, cellH = cellSize, cellSize

		gridW := cellW * float32(g.Rows)
		gridH := cellH * float32(g.Cols)
		offsetX := (parent.Dx() - gridW) / 2
		offsetY := (parent.Dy() - gridH) / 2
		parent = Rect(parent.Min.X+offsetX, parent.Min.Y+offsetY, parent.Max.X-offsetX, parent.Max.Y-offsetY)
	}

	for i, ch := range children {
		if i >= g.Rows*g.Cols {
			break // Сетка заполнена
		}
		row := i % g.Rows
		col := i / g.Rows

		x := parent.Min.X + float32(row)*cellW
		y := parent.Min.Y + float32(col)*cellH

		r := Rect(x, y, x+cellW, y+cellH)

		padX := r.Dx() * g.Padding
		padY := r.Dy() * g.Padding

		inner := Rect(r.Min.X+padX, r.Min.Y+padY, r.Max.X-padX, r.Max.Y-padY)
		ch.Bounds().Set(inner)
	}
}
