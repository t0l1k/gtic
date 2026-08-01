package etic

import (
	"fmt"
	"image/color"
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Point[T Number] struct{ X, Y T }

func (p Point[T]) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", float64(p.X), float64(p.Y))
}

func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.X + q.X, p.Y + q.Y}
}
func (p Point[T]) Sub(q Point[T]) Point[T] {
	return Point[T]{p.X - q.X, p.Y - q.Y}
}
func (p Point[T]) Mul(k T) Point[T] {
	return Point[T]{p.X * k, p.Y * k}
}
func (p Point[T]) Div(k T) Point[T] {
	return Point[T]{p.X / k, p.Y / k}
}
func (p Point[T]) In(r Rectangle[T]) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

type Vec2 struct{ Point[float32] }

var (
	Up    = Pt[float32](0, -1)
	Down  = Pt[float32](0, 1)
	Left  = Pt[float32](-1, 0)
	Right = Pt[float32](1, 0)
)

func Vec2FromAngle(angle float32) Point[float32] {
	return Point[float32]{float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))}
}

func (v Point[T]) Len() float32 {
	return float32(math.Hypot(float64(v.X), float64(v.Y)))
}
func (v Point[T]) Normalized() Point[float32] {
	lenght := v.Len()
	if lenght == 0 {
		return Point[float32]{}
	}
	return Point[float32]{float32(v.X) / lenght, float32(v.Y) / lenght}
}

//	func (p Point[Console]) Mod(r Rectangle[Console]) Point[Console] {
//		w, h := r.Dx(), r.Dy()
//		p = p.Sub(r.Min)
//		p.X = p.X % w
//		if p.X < 0 {
//			p.X += w
//		}
//		p.Y = p.Y % h
//		if p.Y < 0 {
//			p.Y += h
//		}
//		return p.Add(r.Min)
//	}
func (p Point[T]) Eq(q Point[T]) bool {
	return p == q
}

func (p Point[T]) Area() T { return p.X * p.Y }

func Pt[T Number](X, Y T) Point[T] {
	return Point[T]{X, Y}
}

type Rectangle[T Number] struct {
	Min, Max Point[T]
}

func (r Rectangle[T]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}
func (r Rectangle[T]) Dx() T {
	return r.Max.X - r.Min.X
}
func (r Rectangle[T]) Dy() T {
	return r.Max.Y - r.Min.Y
}
func (r Rectangle[T]) Size() Point[T] {
	return Point[T]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}
func (r Rectangle[T]) Add(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X + p.X, r.Min.Y + p.Y},
		Point[T]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}
func (r Rectangle[T]) Sub(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X - p.X, r.Min.Y - p.Y},
		Point[T]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}
func (r Rectangle[T]) Inset(n T) Rectangle[T] {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}
func (r Rectangle[T]) Intersect(s Rectangle[T]) Rectangle[T] {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Empty() {
		return Rectangle[T]{}
	}
	return r
}
func (r Rectangle[T]) Union(s Rectangle[T]) Rectangle[T] {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}
func (r Rectangle[T]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}
func (r Rectangle[T]) Eq(s Rectangle[T]) bool {
	return r == s || r.Empty() && s.Empty()
}
func (r Rectangle[T]) Overlaps(s Rectangle[T]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}
func (r Rectangle[T]) In(s Rectangle[T]) bool {
	if r.Empty() {
		return true
	}
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}
func (r Rectangle[T]) Canon() Rectangle[T] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}
func (r Rectangle[T]) At(x, y T) color.Color {
	if (Point[T]{x, y}).In(r) {
		return color.Opaque
	}
	return color.Transparent
}
func (r Rectangle[T]) Bounds() Rectangle[T] {
	return r
}
func Rect[T Number](x0, y0, x1, y1 T) Rectangle[T] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[T]{Point[T]{x0, y0}, Point[T]{x1, y1}}
}
