package gtic

import (
	"image"
)

type screen struct {
	width, height int
	vram          []RGBA
	clip          image.Rectangle
	cam           image.Point
	pix           []byte
}

func newScreen(w, h int) *screen {
	return &screen{
		width:  w,
		height: h,
		vram:   make([]RGBA, w*h),
		clip:   image.Rect(0, 0, w, h),
		pix:    make([]byte, w*h*4)}
}
func (s *screen) index(x, y int) int {
	sx := x - s.cam.X
	sy := y - s.cam.Y
	if !image.Pt(x, y).In(s.clip) {
		return -1
	}
	if sx < 0 || sy < 0 || sx >= s.width || sy >= s.height {
		return -1
	}
	return (sy*s.width + sx)
}
func (b *screen) position(idx int) (int, int) { return idx % b.width, idx / b.width }
func (s *screen) at(x, y int) RGBA {
	idx := s.index(x, y)
	if idx < 0 {
		return Transparent
	}
	return s.vram[idx]
}
func (s *screen) set(x, y int, c RGBA) {
	if idx := s.index(x, y); idx >= 0 {
		s.vram[idx] = c
	}
}
func (s *screen) clear(c RGBA) {
	firstLine := s.vram[:s.width]
	for i := 0; i < s.width; i++ {
		firstLine[i] = c
	}
	dst := s.vram[s.width:]
	for i := 1; i < s.height; i++ {
		copy(dst, firstLine)
		dst = dst[s.width:]
	}
}
func (s *screen) blit() []byte {
	dst := s.pix
	for i, c := range s.vram {
		r, g, b, a := c.ToBytes()
		dst[i*4+0] = r
		dst[i*4+1] = g
		dst[i*4+2] = b
		dst[i*4+3] = a
	}
	return dst
}
