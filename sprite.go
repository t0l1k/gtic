package gtic

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

type Surface[T any] struct {
	Data          []T
	Width, Height int
}

type SpriteID int
type Sprite struct {
	ID     SpriteID
	name   string
	Width  int
	Height int
	Pixels []RGBA
}

func NewSprite(name string, surface Surface[RGBA], r image.Rectangle) *Sprite {
	w, h, pixels := SpriteFrom(surface.Data, surface.Width, surface.Height, r)
	if len(pixels) != w*h {
		panic("sprite pixel count mismatch")
	}
	return &Sprite{name: name, Width: w, Height: h, Pixels: pixels}
}
func (a *Sprite) Equal(b *Sprite) bool { return a.name == b.name }
func (a *Sprite) String() string {
	return fmt.Sprintf("Spr:%v %v [%v,%v],%v", a.ID, a.name, a.Width, a.Height, len(a.Pixels))
}

type spriteBank struct {
	id      SpriteID
	sprites map[SpriteID]*Sprite
}

func newSpriteBank() *spriteBank { return &spriteBank{sprites: make(map[SpriteID]*Sprite)} }
func (sb *spriteBank) Register(s *Sprite) SpriteID {
	for _, v := range sb.sprites {
		if v.Equal(s) {
			return v.ID
		}
	}
	sb.id++
	s.ID = sb.id
	sb.sprites[sb.id] = s
	return s.ID
}
func (sb *spriteBank) Get(id SpriteID) (*Sprite, bool) {
	s, ok := sb.sprites[id]
	return s, ok
}

// Spr(id SpriteID, x, y int, colorkey RGBA, scale, flip, rotate int)
func (c *API) Spr(id SpriteID, x, y int, args ...interface{}) {
	colorKey := Transparent
	scale := 1
	flip := 0
	rotate := 0

	i := 0
	nextInt := func(def int) int {
		if i < len(args) {
			if v, ok := args[i].(int); ok {
				i++
				return v
			}
		}
		return def
	}
	nextRGBA := func(def RGBA) RGBA {
		if i < len(args) {
			if v, ok := args[i].(RGBA); ok {
				i++
				return v
			}
		}
		return def
	}

	if len(args) > 0 {
		colorKey = nextRGBA(colorKey)
	}
	scale = nextInt(scale)
	flip = nextInt(flip)
	rotate = nextInt(rotate)

	spr, exists := c.Sprites().Get(id)
	if !exists || scale <= 0 {
		return
	}
	w, h := spr.Width, spr.Height
	for sy := 0; sy < h; sy++ {
		for sx := 0; sx < w; sx++ {
			srcX, srcY := sx, sy
			if flip&1 != 0 {
				srcX = w - 1 - srcX
			}
			if flip&2 != 0 {
				srcY = h - 1 - srcY
			}
			for r := 0; r < (rotate % 4); r++ {
				srcX, srcY = w-1-srcY, srcX
			}
			col := spr.Pixels[srcY*spr.Width+srcX]
			if col == colorKey {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					px := x + (sx * scale) + dx
					py := y + (sy * scale) + dy
					c.Pix(px, py, col)
				}
			}
		}
	}
}

// LoadImage загружает PNG
func DecodeImage(pngFile []byte) (Surface[RGBA], error) {
	img, err := png.Decode(bytes.NewReader(pngFile))
	if err != nil {
		return Surface[RGBA]{}, err
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	pix := make([]RGBA, w*h)

	switch m := img.(type) {
	case *image.NRGBA:
		for i, j := 0, 0; i < len(pix); i, j = i+1, j+4 {
			pix[i] = NewRGBA(m.Pix[j], m.Pix[j+1], m.Pix[j+2], m.Pix[j+3])
		}
	case *image.RGBA: // премультиплицирован — распаковываем
		for i, j := 0, 0; i < len(pix); i, j = i+1, j+4 {
			a := m.Pix[j+3]
			pix[i] = NewRGBA(m.Pix[j+0]*0xff/a, m.Pix[j+1]*0xff/a, m.Pix[j+2]*0xff/a, a)
		}
	default: // редкие форматы (палитровые, грейскейл)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, bl, a := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
				pix[y*w+x] = NewRGBA(uint8(r>>8), uint8(g>>8), uint8(bl>>8), uint8(a>>8))
			}
		}
	}
	return Surface[RGBA]{Data: pix, Width: w, Height: h}, nil
}

func SpriteFrom(data []RGBA, w, h int, r image.Rectangle) (int, int, []RGBA) {
	sx, sy := r.Min.X, r.Min.Y
	sw, sh := r.Dx(), r.Dy()
	pix := make([]RGBA, sw*sh)
	for y := 0; y < sh; y++ {
		copy(pix[y*sw:(y+1)*sw], data[(sy+y)*w+sx:(sy+y)*w+sx+sw])
	}
	return sw, sh, pix
}
