package etic

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type InputLine struct {
	Label
	State       Property[ElementState]
	OnReturn    SlotFn[*InputLine]
	OnChanged   Signal[string]
	PlaceHolder Property[string]
}

func NewInputLine(id ElementID, placeHolder string, fn SlotFn[*InputLine]) *InputLine {
	var (
		focused, blink bool
		cursorPos      int
		timer          Timer
	)
	i := &InputLine{
		Label:       *NewLabel(id, ""),
		State:       *NewProperty(ElementNormal),
		OnReturn:    fn,
		OnChanged:   *NewSignal[string](),
		PlaceHolder: *NewProperty(placeHolder),
	}
	timer = *NewTimer("textinput.timer")
	timer.Setup(500*time.Millisecond, false, false)
	timer.Timeout.Connect(func(t *Timer) {
		blink = !blink
	})
	i.Add(&timer)
	i.OnInit = func(t *Console) {
		i.ID().Set(id)
		i.Label.HAlign.Set(text.AlignStart)
	}
	i.Text.OnChange.Connect(func(s string) {
		i.OnChanged.Emit(s)
		if cursorPos > len(s) {
			cursorPos = len(s)
		}
	})
	i.OnUpdate = func(t *Console) {
		if i.hidden.Get() || i.State.Get() == ElementDisabled {
			return
		}
		x, y, l, _, _ := t.Mouse()
		inside := Pt(x, y).In(i.Bounds().Get())
		if inside && l && !focused {
			t.Focus(id)
			i.State.Set(ElementPressed)
		}
		focused = t.focusedId == i.ID().Get()
		if !focused {
			blink = false
			i.State.Set(ElementNormal)
			timer.Stop()
			return
		}
		if timer.IsStopped() {
			timer.Start()
		}
		value := i.Text.Get()
		if chars, ok := t.CharsP(); ok {
			s := value
			for _, r := range chars {
				if cursorPos < 0 {
					cursorPos = 0
				}
				if cursorPos > len(value) {
					cursorPos = len(value)
				}
				s = value[:cursorPos] + string(r) + value[cursorPos:]
				cursorPos++
			}
			if value != s {
				value = s
				i.Text.Set(value)
			}
		}
		value = i.Text.Get()
		switch {
		case t.BtnP(ebiten.KeyEnter):
			if i.OnReturn != nil {
				i.OnReturn(i)
			}
		case t.BtnP(ebiten.KeyArrowLeft):
			if cursorPos > 0 {
				cursorPos--
			}
		case t.BtnP(ebiten.KeyArrowRight):
			if cursorPos < len(value) {
				cursorPos++
			}
		case t.BtnP(ebiten.KeyBackspace):
			if cursorPos > 0 && len(value) > 0 {
				result := value[:cursorPos-1] + value[cursorPos:]
				cursorPos--
				i.Text.Set(result)
			}
		case t.BtnP(ebiten.KeyDelete):
			if cursorPos > 0 && len(value) > 0 {
				result := value[:cursorPos] + value[cursorPos:]
				i.Text.Set(result)
			}
		}
	}
	i.OnDraw = func(t *Console) {
		if i.hidden.Get() {
			return
		}
		value := i.Text.Get()
		display := ""
		if len(value) == 0 && !focused {
			display = i.PlaceHolder.Get()
		} else {
			display = value
		}
		x, y, w, h := RectF32(i.Bounds().Get())
		t.Rect(x, y, w, h, i.State.Get().Color())
		t.DrawString(display, i.Label.Font.Get(), i.Bounds().Get(), i.HAlign.Get(), i.VAlign.Get(), i.Color.Get())

		if blink && focused {
			x, y, _, h := RectF32(i.Bounds().Get())
			face := t.Face(i.Font.Get())
			w, _ := text.Measure(display[:cursorPos], face, face.Size*1.2)
			cursorX := x + float32(w)
			cursorH := float32(face.Size)
			cursorY := y + (h-cursorH)/2
			t.Rect(cursorX, cursorY, 3, cursorH, i.Color.Get())
		}
	}
	return i
}
