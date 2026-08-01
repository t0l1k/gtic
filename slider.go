package etic

import (
	"golang.org/x/image/colornames"
)

type Slider struct {
	Element
	min, max, step float32
	Value          Property[float32]
	OnChanged      Signal[float32]
	orien          Orientation
}

func NewSlider(id ElementID, initial, min, max, step float32, orient Orientation, fn SlotFn[float32]) *Slider {
	s := &Slider{
		Element:   *NewElement(id),
		OnChanged: *NewSignal[float32](),
		Value:     *NewProperty(initial),
		min:       min,
		max:       max,
		step:      step,
		orien:     orient,
	}
	s.OnChanged.Connect(fn)
	s.Value.OnChange.Connect(func(f float32) {
		value := Clamp(f, s.min, s.max)
		s.OnChanged.Emit(value)
	})

	var (
		track, fill, thumb     Rectangle[float32]
		trackH, thumbW, thumbH float32 = 10, 14, 18
		p                      float32
		thumbState             ElementState = ElementNormal
	)
	s.OnUpdate = func(t *Console) {

		r := s.Bounds().Get()
		if s.orien == Horizontal {
			track = Rect(r.Min.X, r.Min.Y+(thumbH-trackH)/2, r.Max.X, r.Min.Y+(thumbH-trackH)/2+trackH)
		} else {
			track = Rect(r.Min.X+(thumbW-trackH)/2, r.Min.Y, r.Min.X+(thumbW-trackH)/2+trackH, r.Max.Y)
		}

		rangeVal := s.max - s.min
		if rangeVal == 0 {
			rangeVal = 1
		}
		p = (s.Value.Get() - s.min) / rangeVal
		p = Clamp(p, 0, 1)

		fill = track
		if s.orien == Horizontal {
			fill.Max.X = track.Min.X + p*track.Dx()
			thumbX := track.Min.X + (p * (track.Dx() - thumbW))
			thumb = Rect(thumbX, r.Min.Y, thumbX+thumbW, r.Min.Y+thumbH)
		} else {
			fill.Min.Y = track.Max.Y - p*track.Dy()
			thumbY := track.Max.Y - (p * (track.Dy() - thumbH))
			thumb = Rect(r.Min.X, thumbY, r.Min.X+thumbW, thumbY+thumbH)
		}

		x0, y0, mousePressed, _, _ := t.Mouse()
		pt := Pt(x0, y0)
		hover := pt.In(thumb) || pt.In(track)
		active := t.activeId == id

		if hover && mousePressed {
			t.activeId = id
			active = true
		}
		if !mousePressed && active {
			t.activeId = ""
			active = false
		}

		if active && mousePressed {
			if s.orien == Horizontal {
				x := pt.X
				x = Clamp(x, track.Min.X, track.Max.X)
				if track.Dx() <= 0 {
					newVal := s.min
					newVal = roundToStep(newVal, s.step)
					newVal = Clamp(newVal, s.min, s.max)
					if newVal != s.Value.Get() {
						s.Value.Set(newVal)
					}
				} else {
					newVal := s.min + ((x-track.Min.X)/track.Dx())*rangeVal
					newVal = roundToStep(newVal, s.step)
					newVal = Clamp(newVal, s.min, s.max)
					if newVal != s.Value.Get() {
						s.Value.Set(newVal)
					}
				}
			} else {
				y := pt.Y
				y = Clamp(y, track.Min.Y, track.Max.Y)
				if track.Dy() <= 0 {
					newVal := s.min
					newVal = roundToStep(newVal, s.step)
					newVal = Clamp(newVal, s.min, s.max)
					if newVal != s.Value.Get() {
						s.Value.Set(newVal)
					}
				} else {
					newVal := s.max - ((y-track.Min.Y)/track.Dy())*rangeVal
					newVal = roundToStep(newVal, s.step)
					newVal = Clamp(newVal, s.min, s.max)
					if newVal != s.Value.Get() {
						s.Value.Set(newVal)
					}
				}
			}
		}

		switch {
		case hover && !active:
			thumbState = ElementHover
		case active:
			thumbState = ElementPressed
		}
	}
	s.OnDraw = func(t *Console) {
		x, y, w, h := RectF32(track)
		t.Rect(x, y, w, h, colornames.Blue)
		if fill.Dx() > 0 {
			x, y, w, h = RectF32(fill)
			t.Rect(x, y, w, h, colornames.Yellow)
		}
		x, y, w, h = RectF32(thumb)
		t.Rect(x, y, w, h, thumbState.Color())
	}
	return s
}
