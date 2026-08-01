package game

import (
	"etic"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

const playerSpeed float32 = 400

type player struct {
	etic.Element
	hit         etic.Signal[bool]
	speed       float32
	anim        etic.Animation
	walk, fly   []etic.SpriteID
	IsColliding bool
}

func NewPlayer(id etic.ElementID) *player {
	return &player{Element: *etic.NewElement(id), hit: *etic.NewSignal[bool](), anim: *etic.NewAnimation("player.animation")}
}
func (p *player) Setup(a, b []etic.SpriteID) *player {
	p.walk = a
	p.fly = b
	p.anim.Setup(p.walk, 100*time.Millisecond)
	p.Add(&p.anim)
	p.Bounds().OnChange.Connect(func(r etic.Rectangle[float32]) {
		p.anim.Bounds().Set(r)
	})
	return p
}
func (p *player) Init(t *etic.Console) {
	p.speed = playerSpeed
	p.hit.Emit(false)
	p.Show()
	scale := float32(scale)
	bounds := t.Sprite(p.walk[0]).Bounds()
	w, h := float32(bounds.Dx())*scale, float32(bounds.Dy())*scale
	x := (t.Width - w) / 2
	y := (t.Height - h) / 2
	p.anim.SetScale(scale)
	p.Bounds().Set(etic.Rect(x, y, x+w, y+h))
	p.MarkReady()
	log.Println("player:Init", p.ID().Get(), p.Bounds().Get(), p.anim.Bounds().Get())
}
func (p *player) Update(t *etic.Console) {
	if !p.IsReady() {
		p.Init(t)
	}
	if p.IsHidden() {
		return
	}
	p.anim.Update(t)
	var dx, dy float32 = 0.0, 0.0
	switch {
	case t.Btn(ebiten.KeyArrowUp):
		dy--
	case t.Btn(ebiten.KeyArrowDown):
		dy++
	case t.Btn(ebiten.KeyArrowLeft):
		dx--
	case t.Btn(ebiten.KeyArrowRight):
		dx++
	}
	if dx != 0 || dy != 0 {
		p.Move(etic.Pt(dx, dy).Mul(p.speed*float32(t.Tick().Seconds())), t.Width, t.Height)
		var flip float32 = 0.0
		switch {
		case dx != 0:
			p.anim.SetSprites(p.walk)
			if dx < 0 {
				flip = 1
			}
		case dy != 0:
			p.anim.SetSprites(p.fly)
			if dy > 0 {
				flip = 2
			}
		}
		p.anim.SetFlip(flip)
		if p.anim.IsStopped() {
			p.anim.Play()
		}
	} else {
		p.anim.Stop()
	}
}
func (p *player) Move(v etic.Point[float32], w0, h0 float32) {
	r := p.Bounds().Get().Add(v)
	x, y, w, h := etic.RectF32(r)
	pt := etic.Pt(x, y)
	x = etic.Clamp(x, 0, w0-w)
	y = etic.Clamp(y, 0, h0-h)
	if x != pt.X || y != pt.Y { //если касание края пропустить движение
		return
	}
	p.Bounds().Set(r)
}
func (p *player) Draw(t *etic.Console) {
	if p.IsHidden() || p.Bounds().Get().Empty() {
		return
	}
	p.anim.Draw(t)
	col := colornames.Blue
	if p.IsColliding {
		col = colornames.Red
	}
	bounds := p.Bounds().Get()
	x, y, w, h := etic.RectF32(bounds)
	t.RectB(x, y, w, h, col)
}
