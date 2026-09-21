package gtic

import "image"

type Number interface {
	~int | ~float64 | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~uint | ~byte | ~uint16 | ~uint32 | ~uint64
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func sign(value int) int {
	if value < 0 {
		return -1
	}
	return 1
}

func RectXYWH(r image.Rectangle) (int, int, int, int) { return r.Min.X, r.Min.Y, r.Dx(), r.Dy() }
