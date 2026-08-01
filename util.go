package etic

import "image"

func Clamp[T Number](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func RectF32ToIntRect(r Rectangle[float32]) image.Rectangle {
	var x, y, w, h int
	x, y, w, h = int(r.Min.X), int(r.Min.Y), int(r.Dx()), int(r.Dy())
	return image.Rect(x, y, x+w, y+h)
}

func RectIntToRectF32(r image.Rectangle) Rectangle[float32] {
	var x, y, w, h float32
	x, y, w, h = float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy())
	return Rect(x, y, x+w, y+h)
}

func RectF32(r Rectangle[float32]) (float32, float32, float32, float32) {
	return float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy())
}
func RectInt(r Rectangle[int]) (int, int, int, int) {
	return r.Min.X, r.Min.Y, r.Dx(), r.Dy()
}

func roundToStep(v, step float32) float32 {
	n := int((v / step) + 0.5)
	return float32(n) * step
}
