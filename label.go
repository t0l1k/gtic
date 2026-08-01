package etic

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
)

type Label struct {
	Element
	Text      Property[string]
	Color     Property[color.Color]
	HAlign    Property[text.Align]
	VAlign    Property[text.Align]
	Font      Property[FontKind]
	ShowFrame bool
}

func NewLabel(id ElementID, v string) *Label {
	l := &Label{Element: *NewElement(id),
		Text:   *NewPropertyWithEqual(v, func(a, b string) bool { return false }),
		Color:  *NewProperty[color.Color](colornames.Yellow),
		HAlign: *NewProperty(text.AlignCenter),
		VAlign: *NewProperty(text.AlignCenter),
		Font:   *NewProperty(FontNormal),
	}
	l.OnDraw = func(t *Console) {
		if l.hidden.Get() || l.Bounds().Get().Empty() {
			return
		}
		if l.ShowFrame {
			x, y, w, h := RectF32(l.Bounds().Get())
			t.RectB(x, y, w, h, colornames.Yellow)
		}
		t.DrawString(l.Text.Get(), l.Font.Get(), l.Bounds().Get(), l.HAlign.Get(), l.VAlign.Get(), l.Color.Get())
	}
	return l
}
