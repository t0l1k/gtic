package etic

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type Flip int

const (
	FlipNone Flip = iota
	FlipH
	FlipV
	FlipHV
)

type Rotate int

const (
	RotateNone Rotate = iota
	Rotate90
	Rotate180
	Rotate270
	RotateAngle
)

type SprOpt struct {
	Scale  float32
	Flip   Flip
	Rotate Rotate // 0/90/180/270/1-360 в градусах
	Angle  float32
}

// Spr id, x, y, [scale, flip, rotate]
// scale 0.5 = в 2 раза меньше, 1=original scale
// flip повернуть по горизонтали/вертикали/оба
// rotate в градусах
// colorkey обработать при загрузке → img уже прозрачный
func (t *Console) Spr(id SpriteID, x, y float32, opt SprOpt) {
	img := t.sprites[id]
	if img == nil {
		return
	}

	sw := float64(img.Bounds().Dx())
	sh := float64(img.Bounds().Dy())

	scale := float64(opt.Scale)
	if scale <= 0 {
		scale = 1
	}

	w := sw * scale
	h := sh * scale

	op := &ebiten.DrawImageOptions{}

	switch {
	case opt.Flip != FlipNone:
		op.GeoM.Scale(scale, scale)
		switch opt.Flip {
		case FlipH:
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(w, 0)
		case FlipV:
			op.GeoM.Scale(1, -1)
			op.GeoM.Translate(0, h)
		case FlipHV:
			op.GeoM.Scale(-1, -1)
			op.GeoM.Translate(w, h)
		}
		op.GeoM.Translate(float64(x), float64(y))

	case opt.Rotate != RotateNone:

		toRad := func(angle float64) float64 { return angle * math.Pi / 180 }

		radian := 00.0
		switch opt.Rotate {
		case Rotate90:
			radian = toRad(90)
		case Rotate180:
			radian = toRad(180)
		case Rotate270:
			radian = toRad(270)
		case RotateAngle:
			radian = toRad(float64(opt.Angle))
		}
		cx, cy := w/2, h/2

		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Rotate(radian)
		op.GeoM.Translate(cx, cy)

		op.GeoM.Translate(float64(x), float64(y))

	default:
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x), float64(y))
	}
	t.screen.DrawImage(img, op)
}
func (t *Console) SpriteBounds() Rectangle[float32] { return t.spriteBounds }
func (t *Console) Sprite(id SpriteID) *ebiten.Image { return t.sprites[id] }

func (t *Console) LoadSprite(id SpriteID, img *ebiten.Image) {
	if t.sprites == nil {
		t.sprites = make(map[SpriteID]*ebiten.Image)
	}
	x, y := 0, 0
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	spr := ebiten.NewImage(w, h)
	rect := image.Rect(x, y, w, h)
	spr.DrawImage(img.SubImage(rect).(*ebiten.Image), &ebiten.DrawImageOptions{})
	t.sprites[id] = spr
}

func (t *Console) LoadSpriteSheet(img *ebiten.Image, w, h int) {
	if t.sprites == nil {
		t.sprites = make(map[SpriteID]*ebiten.Image)
	}
	count := img.Bounds().Dx() / w
	for i := 0; i < count; i++ {
		nSpr := ebiten.NewImage(w, h)
		x, y := i*w, 0
		rect := image.Rect(x, y, x+w, y+h)
		nSpr.DrawImage(img.SubImage(rect).(*ebiten.Image), &ebiten.DrawImageOptions{})
		t.sprites[SpriteID(strconv.Itoa(i))] = nSpr
	}
	for k, v := range t.sprites {
		log.Println("loaded", k, v.Bounds())
	}
}

func LoadImage(value []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(value))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

func ApplyColorKey(img *ebiten.Image, key color.Color) *ebiten.Image {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	out := ebiten.NewImage(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			if c == key {
				out.Set(x, y, color.RGBA{0, 0, 0, 0})
			} else {
				out.Set(x, y, c)
			}
		}
	}
	return out
}
