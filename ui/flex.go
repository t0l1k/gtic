package ui

import (
	"gtic"
	"image"
)

type Layout interface {
	Apply(image.Rectangle, []IElement)
}

// Макет настройки макета вручную, вызовом Property("rect").OnChange.Connect(func(r image.Rectangle[float32]) { тут настройка макета})
type AbsoluteLayout struct{}

func (l AbsoluteLayout) Apply(r image.Rectangle, children []IElement) {}

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

// FlexDir задаёт направление размещения дочерних элементов.
type FlexDir int

// AlignType задаёт выравнивание поперёк направления Flex.
type AlignType int

// JustifyType задаёт распределение свободного места вдоль направления Flex.
type JustifyType int

const (
	// Row размещает элементы слева направо.
	Row FlexDir = iota
	// Col размещает элементы сверху вниз.
	Col
)

const (
	// AlignStart прижимает элементы к началу поперечной оси.
	AlignStart AlignType = iota
	// AlignCenter центрирует элементы по поперечной оси.
	AlignCenter
	// AlignEnd прижимает элементы к концу поперечной оси.
	AlignEnd
)

const (
	// JustifyStart размещает элементы от начала главной оси.
	JustifyStart JustifyType = iota
	// JustifyCenter центрирует группу элементов на главной оси.
	JustifyCenter
	// JustifyEnd размещает группу у конца главной оси.
	JustifyEnd
	// JustifySpaceBetween равномерно распределяет свободное место между элементами.
	JustifySpaceBetween
	// JustifySpaceAround распределяет свободное место вокруг элементов.
	JustifySpaceAround
)

// Flex делит прямоугольник родителя между дочерними элементами.
// Padding задаёт отступ от краёв в пикселях, а GrowFactors — относительные
// доли главной оси для каждого ребёнка. Align управляет поперечной осью,
// Justify — свободным местом на главной оси. Gap задаёт промежуток между
// детьми; ненулевые RowGap и ColumnGap переопределяют Gap для направления.
// Все промежутки задаются в пикселях.
type Flex struct {
	Dir           FlexDir
	Align         AlignType
	Justify       JustifyType
	Padding       int
	Gap           int
	RowGap        int
	ColumnGap     int
	GrowFactors   []float64
	computedRects [][4]int
}

// NewFlex создаёт макет в направлении dir. Если growFactors пуст, дети
// получают равные доли; иначе для каждого ребёнка нужно указать вес.
// Например, []float64{1, 2} делит главную ось в отношении 1:2.
func NewFlex(dir FlexDir, padding int, align AlignType, justify JustifyType, growFactors []float64) *Flex {
	return &Flex{
		Dir:         dir,
		Padding:     padding,
		Align:       align,
		Justify:     justify,
		GrowFactors: growFactors,
	}
}

func NewVerticalFlexLayout(fraction []float64, padding int) *Flex {
	return NewFlex(Col, padding, AlignCenter, JustifySpaceBetween, fraction)
}
func NewHorizontalFlexLayout(fraction []float64, padding int) *Flex {
	return NewFlex(Row, padding, AlignCenter, JustifySpaceBetween, fraction)
}
func NewVerticalBoxLayout(padding int) *Flex {
	return NewFlex(Col, padding, AlignCenter, JustifySpaceBetween, nil)
}

func NewHorizontalBoxLayout(padding int) *Flex {
	return NewFlex(Row, padding, AlignCenter, JustifySpaceBetween, nil)
}

func (f Flex) mainGap() int {
	if f.Dir == Row {
		if f.ColumnGap != 0 {
			return f.ColumnGap
		}
		return f.Gap
	}
	if f.RowGap != 0 {
		return f.RowGap
	}
	return f.Gap
}

func (f Flex) Apply(rect image.Rectangle, children []IElement) {
	if len(children) == 0 {
		return
	}
	x, y, w, h := gtic.RectXYWH(rect)
	f.computedRects = make([][4]int, len(children))
	itemGap := f.mainGap()

	// общий вес
	totalGrow := 0.0
	for _, g := range f.GrowFactors {
		totalGrow += g
	}
	if totalGrow == 0 {
		for range children {
			f.GrowFactors = append(f.GrowFactors, 1)
		}
		totalGrow = float64(len(children))
	}

	if f.Dir == Row {
		extra := w - f.Padding*2 - itemGap*(len(children)-1)
		// предварительное вычисление общей ширины
		totalW := 0
		for i := range children {
			totalW += int((f.GrowFactors[i] / totalGrow) * float64(extra))
		}

		freeSpace := extra - totalW
		cursor := x + f.Padding
		gap := 0
		switch f.Justify {
		case JustifyCenter:
			cursor += freeSpace / 2
		case JustifyEnd:
			cursor += freeSpace
		case JustifySpaceBetween:
			if len(children) > 1 {
				gap = freeSpace/(len(children)-1) + itemGap
			}
		case JustifySpaceAround:
			if len(children) > 0 {
				gap = freeSpace/len(children) + itemGap
				cursor += (freeSpace / len(children)) / 2
			}
		default:
			gap = itemGap
		}

		for i, c := range children {
			cw := int((f.GrowFactors[i] / totalGrow) * float64(extra))
			ch := h - f.Padding*2
			cx, cy := cursor, y+f.Padding

			switch f.Align {
			case AlignCenter:
				cy = y + (h-ch)/2
			case AlignEnd:
				cy = y + h - ch - f.Padding
			}

			f.computedRects[i] = [4]int{cx, cy, cw, ch}
			r := image.Rect(cx, cy, cx+cw, cy+ch)
			c.Property("rect").Set(r)
			cursor += cw + gap
		}
	} else {
		extra := h - f.Padding*2 - itemGap*(len(children)-1)
		totalH := 0
		for i := range children {
			totalH += int((f.GrowFactors[i] / totalGrow) * float64(extra))
		}

		freeSpace := extra - totalH
		cursor := y + f.Padding
		gap := 0
		switch f.Justify {
		case JustifyCenter:
			cursor += freeSpace / 2
		case JustifyEnd:
			cursor += freeSpace
		case JustifySpaceBetween:
			if len(children) > 1 {
				gap = freeSpace/(len(children)-1) + itemGap
			}
		case JustifySpaceAround:
			if len(children) > 0 {
				gap = freeSpace/len(children) + itemGap
				cursor += (freeSpace / len(children)) / 2
			}
		default:
			gap = itemGap
		}

		for i, c := range children {
			ch := int((f.GrowFactors[i] / totalGrow) * float64(extra))
			cw := w - f.Padding*2
			cx, cy := x+f.Padding, cursor

			switch f.Align {
			case AlignCenter:
				cx = x + (w-cw)/2
			case AlignEnd:
				cx = x + w - cw - f.Padding
			}

			f.computedRects[i] = [4]int{cx, cy, cw, ch}
			r := image.Rect(cx, cy, cx+cw, cy+ch)
			c.Property("rect").Set(r)
			cursor += ch + gap
		}
	}
}
