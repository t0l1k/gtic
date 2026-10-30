package gtic

import "image"

type Number interface {
	~int | ~float64 | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~uint | ~byte | ~uint16 | ~uint32 | ~uint64
}

func RectXYWH(r image.Rectangle) (int, int, int, int) { return r.Min.X, r.Min.Y, r.Dx(), r.Dy() }

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

func ValueFromPercent(percent, total float32) float32 { return percent * total / 100 }
func PercentOf(value, total float32) float32          { return 100 * value / total }

func sign(value int) int {
	if value < 0 {
		return -1
	}
	return 1
}
