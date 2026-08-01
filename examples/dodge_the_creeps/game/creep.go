package game

import (
	"etic"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"
)

const (
	creepBaseSpeed   float32 = 20
	creepSpeedSpread float32 = 100
)

type creep struct {
	etic.Element
	velocity    etic.Point[float32]
	mobKind     []etic.SpriteID
	anim        etic.Animation
	IsColliding bool
}

func NewCreep(id etic.ElementID) *creep {
	return &creep{Element: *etic.NewElement(id), anim: *etic.NewAnimation("creep.animation")}
}
func (p *creep) Setup(a []etic.SpriteID) *creep {
	p.mobKind = a
	p.anim.Setup(a, 100*time.Millisecond)
	p.anim.SetScale(float32(scale))
	p.Add(&p.anim)
	p.Bounds().OnChange.Connect(func(r etic.Rectangle[float32]) {
		p.anim.Bounds().Set(r)
	})
	return p
}
func (c *creep) Init(tic *etic.Console) {
	side := rand.Intn(4)
	var x, y float32
	switch side {
	case 0: // сверху
		x = rand.Float32() * tic.Width
		y = -50
	case 1: // снизу
		x = rand.Float32() * tic.Width
		y = tic.Height + 50
	case 2: // слева
		x = -50
		y = rand.Float32() * tic.Height
	case 3: // справа
		x = tic.Width + 50
		y = rand.Float32() * tic.Height
	}
	// направление к центру экрана
	pt := etic.Pt(tic.Width/2, tic.Height/2)
	dir := pt.Sub(etic.Pt(x, y)).Normalized()
	speed := float32(creepBaseSpeed + rand.Float32()*creepSpeedSpread)
	c.velocity = dir.Mul(speed)
	// угол движения
	angle := float32(math.Atan2(float64(dir.Y), float64(dir.X)) * 180 / math.Pi)
	c.anim.SetAngle(angle)
	r := tic.Sprite(c.mobKind[0]).Bounds()
	w, h := float32(r.Dx())*float32(scale), float32(r.Dy())*float32(scale)
	c.Bounds().Set(etic.Rect(x, y, x+w, y+h))
	c.MarkReady()
}
func (c *creep) Update(t *etic.Console) {
	if !c.IsReady() {
		c.Init(t)
	}
	if c.IsHidden() {
		return
	}
	r := c.Bounds().Get().Add(c.velocity.Mul(float32(t.Tick().Seconds())))
	c.Bounds().Set(r)
	if c.anim.IsStopped() {
		c.anim.Play()
	}
	c.anim.Update(t)
}
func (p *creep) Draw(t *etic.Console) {
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
