package etic

import (
	"slices"
	"time"
)

type Animation struct {
	Element
	sprites            []SpriteID
	timer              *Timer
	index              int
	flip, angle, scale float32
}

func NewAnimation(id ElementID) *Animation {
	a := &Animation{Element: *NewElement(id)}
	a.timer = NewTimer("animation.timer")
	a.Add(a.timer)
	return a
}
func (a *Animation) Setup(
	sprites []SpriteID,
	waitTime time.Duration,
) *Animation {
	a.sprites = sprites
	a.timer.Setup(waitTime, false, false)
	a.timer.Timeout.Connect(func(t *Timer) {
		a.index = (a.index + 1) % len(a.sprites)
	})
	a.scale = 1
	return a
}
func (a *Animation) SetSprites(sprites []SpriteID) {
	if slices.Equal(a.sprites, sprites) {
		return
	}
	a.sprites = sprites
	a.Reset()
}
func (a *Animation) Move(v Point[float32]) { a.Bounds().Set(a.Bounds().Get().Add(v)) }
func (a *Animation) SetAngle(v float32)    { a.angle = v }
func (a *Animation) SetFlip(v float32)     { a.flip = v }
func (a *Animation) SetScale(v float32)    { a.scale = v }
func (a *Animation) Reset()                { a.index = 0 }
func (t *Animation) IsStopped() bool       { return t.timer.IsStopped() }
func (a *Animation) Play()                 { a.timer.Start() }
func (a *Animation) Stop()                 { a.timer.Stop() }
func (a *Animation) Draw(t *Console) {
	if len(a.sprites) == 0 {
		return
	}
	x, y := a.Bounds().Get().Min.X, a.Bounds().Get().Min.Y
	t.Spr(a.sprites[a.index], x, y,
		SprOpt{
			Scale:  a.scale,
			Flip:   Flip(a.flip),
			Rotate: RotateAngle,
			Angle:  a.angle,
		})
}
