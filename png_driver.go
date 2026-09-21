package gtic

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type pngDriver struct{ Path string }

func (d *pngDriver) Run(e *Engine) error {
	e.Tick(Input{})
	s := e.api.vram
	w, h := s.width, s.height
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for i, c := range s.vram {
		r, g, b, a := c.ToBytes()
		x, y := s.position(i)
		img.SetRGBA(x, y, color.RGBA{r, g, b, a})
	}

	f, err := os.Create(d.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
