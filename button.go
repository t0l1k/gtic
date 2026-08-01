package etic

type Button struct {
	Label
	State   Property[ElementState]
	OnClick SlotFn[*Button]
}

func NewButton(id ElementID, txt string, fn SlotFn[*Button]) *Button {
	b := &Button{Label: *NewLabel(id, txt), OnClick: fn, State: *NewProperty(ElementNormal)}
	b.Text.Set(txt)
	b.OnUpdate = func(t *Console) {
		if b.hidden.Get() || b.State.Get() == ElementDisabled {
			return
		}
		x, y, l, _, _ := t.Mouse()
		inside := Pt(x, y).In(b.Bounds().Get())
		switch {
		case !inside && !(b.State.Get() == ElementNormal):
			b.State.Set(ElementNormal)
		case inside && l && !(b.State.Get() == ElementPressed):
			b.State.Set(ElementPressed)
		case inside && !l && !(b.State.Get() == ElementHover || b.State.Get() == ElementPressed):
			b.State.Set(ElementHover)
		case inside && !l && b.State.Get() == ElementPressed:
			if b.OnClick != nil {
				b.OnClick(b)
			}
			b.State.Set(ElementHover)
		}
	}
	b.OnDraw = func(t *Console) {
		if b.hidden.Get() {
			return
		}
		r := b.Bounds().Get()
		s := b.State.Get()
		var stroke float32 = 3
		col := s.Color()
		off := Pt[float32](0, 0)
		if s == ElementPressed {
			off = Pt(stroke, stroke)
		}
		r = r.Add(off)
		x, y, w, h := RectF32(r)
		t.Rect(x, y, w, h, col)
		t.DrawString(b.Text.Get(), b.Font.Get(), r, b.HAlign.Get(), b.VAlign.Get(), b.Color.Get())
	}
	return b
}
