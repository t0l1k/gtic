package gtic

import (
	"image"
	"math"
)

func (a *API) Clip(rect ...int) {
	x, y, w, h := 0, 0, a.Bounds().Width, a.Bounds().Height
	if len(rect) > 0 {
		x, y, w, h = rect[0], rect[1], rect[2], rect[3]
	}
	a.vram.clip = image.Rect(x, y, x+w, y+h)
}
func (a *API) Camera(x, y int) { a.vram.cam.X = x; a.vram.cam.Y = y }

func (a *API) Cls(col ...RGBA) {
	var c RGBA = a.Pal(0)
	if len(col) > 0 {
		c = col[0]
	}
	a.vram.clear(c)
}

func (a *API) Pix(x, y int, col ...RGBA) RGBA {
	if len(col) > 0 {
		a.vram.set(x, y, col[0])
	}
	return a.vram.at(x, y)
}

// Line draws a one-pixel-wide line using Bresenham's algorithm.
func (a *API) Line(x0, y0, x1, y1 int, colour RGBA) {
	dx := abs(x1 - x0)
	sx := sign(x1 - x0)
	dy := -abs(y1 - y0)
	sy := sign(y1 - y0)
	err := dx + dy
	for {
		a.vram.set(x0, y0, colour)
		if x0 == x1 && y0 == y1 {
			return
		}
		error2 := 2 * err
		if error2 >= dy {
			err += dy
			x0 += sx
		}
		if error2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

// Rect draws a filled rectangle. Non-positive dimensions draw nothing.
func (a *API) Rect(x, y, width, height int, colour RGBA) {
	if width <= 0 || height <= 0 {
		return
	}
	for row := y; row < y+height; row++ {
		for column := x; column < x+width; column++ {
			a.vram.set(column, row, colour)
		}
	}
}

// RectB draws a one-pixel rectangle border. Non-positive dimensions draw nothing.
func (a *API) RectB(x, y, width, height int, colour RGBA) {
	if width <= 0 || height <= 0 {
		return
	}
	a.Line(x, y, x+width-1, y, colour)
	a.Line(x, y+height-1, x+width-1, y+height-1, colour)
	a.Line(x, y, x, y+height-1, colour)
	a.Line(x+width-1, y, x+width-1, y+height-1, colour)
}

// Circ draws a filled circle centred at x, y with radius.
func (a *API) Circ(x, y, radius int, colour RGBA) {
	a.Elli(x, y, radius, radius, colour)
}

// CircB draws a one-pixel circle border centred at x, y with radius.
func (a *API) CircB(x, y, radius int, colour RGBA) {
	a.ElliB(x, y, radius, radius, colour)
}

// Elli draws a filled ellipse centred at x, y with horizontal and vertical radii.
func (a *API) Elli(x, y, radiusX, radiusY int, colour RGBA) {
	if radiusX < 0 || radiusY < 0 {
		return
	}
	if radiusX == 0 {
		a.Line(x, y-radiusY, x, y+radiusY, colour)
		return
	}
	if radiusY == 0 {
		a.Line(x-radiusX, y, x+radiusX, y, colour)
		return
	}
	for offsetY := -radiusY; offsetY <= radiusY; offsetY++ {
		vertical := float64(offsetY*offsetY) / float64(radiusY*radiusY)
		halfWidth := int(math.Floor(float64(radiusX) * math.Sqrt(1-vertical)))
		a.Line(x-halfWidth, y+offsetY, x+halfWidth, y+offsetY, colour)
	}
}

// ElliB draws a one-pixel ellipse border centred at x, y with horizontal and
// vertical radii. Border segments are joined so steep ellipses remain connected.
func (a *API) ElliB(x, y, radiusX, radiusY int, colour RGBA) {
	if radiusX < 0 || radiusY < 0 {
		return
	}
	if radiusX == 0 {
		a.Line(x, y-radiusY, x, y+radiusY, colour)
		return
	}
	if radiusY == 0 {
		a.Line(x-radiusX, y, x+radiusX, y, colour)
		return
	}
	previousLeftX, previousRightX := 0, 0
	previousY := 0
	first := true
	for offsetY := -radiusY; offsetY <= radiusY; offsetY++ {
		vertical := float64(offsetY*offsetY) / float64(radiusY*radiusY)
		halfWidth := int(math.Floor(float64(radiusX) * math.Sqrt(1-vertical)))
		leftX, rightX := x-halfWidth, x+halfWidth
		currentY := y + offsetY
		if first {
			a.vram.set(leftX, currentY, colour)
			a.vram.set(rightX, currentY, colour)
			first = false
		} else {
			a.Line(previousLeftX, previousY, leftX, currentY, colour)
			a.Line(previousRightX, previousY, rightX, currentY, colour)
		}
		previousLeftX, previousRightX, previousY = leftX, rightX, currentY
	}
}

// Контурный треугольник — три линии (алгоритм Брезенхэма)
func (a *API) TriB(x1, y1, x2, y2, x3, y3 int, col RGBA) {
	a.Line(x1, y1, x2, y2, col)
	a.Line(x2, y2, x3, y3, col)
	a.Line(x3, y3, x1, y1, col)
}

// Элегантный и быстрый CPU-растеризатор треугольника
func (a *API) Tri(x0, y0, x1, y1, x2, y2 float64, col RGBA) {
	// Сортировка вершин по возрастанию Y (y0 <= y1 <= y2)
	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	if y0 > y2 {
		x0, x2 = x2, x0
		y0, y2 = y2, y0
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	// Защита от деления на ноль для вырожденных треугольников
	if int(y0) == int(y2) {
		return
	}

	// Проверка на простой треугольник с плоским низом
	if int(y1) == int(y2) {
		a.drawTriangleFlatBottom(x0, y0, x1, y1, x2, y2, col)
	} else if int(y0) == int(y1) { // С плоским верхом
		a.drawTriangleFlatTop(x0, y0, x1, y1, x2, y2, col)
	} else {
		// Общий случай: разделение треугольника на два простых
		// Находим точку пересечения горизонтальной линии из (x1, y1) с противоположной гранью (x0,y0)-(x2,y2)
		splitX := x0 + ((y1-y0)/(y2-y0))*(x2-x0)

		a.drawTriangleFlatBottom(x0, y0, x1, y1, splitX, y1, col)
		a.drawTriangleFlatTop(x1, y1, splitX, y1, x2, y2, col)
	}
}

// Вспомогательная функция для отрисовки треугольника с плоским низом
func (a *API) drawTriangleFlatBottom(x0, y0, x1, y1, x2, y2 float64, col RGBA) {
	invslope1 := (x1 - x0) / (y1 - y0)
	invslope2 := (x2 - x0) / (y2 - y0)

	curx1 := x0
	curx2 := x0

	for scanline := int(y0); scanline <= int(y1); scanline++ {
		a.Line(int(curx1), scanline, int(curx2), scanline, col)
		curx1 += invslope1
		curx2 += invslope2
	}
}

// Вспомогательная функция для отрисовки треугольника с плоским верхом
func (a *API) drawTriangleFlatTop(x0, y0, x1, y1, x2, y2 float64, col RGBA) {
	invslope1 := (x2 - x0) / (y2 - y0)
	invslope2 := (x2 - x1) / (y2 - y1)

	curx1 := x2
	curx2 := x2

	for scanline := int(y2); scanline > int(y0); scanline-- {
		a.Line(int(curx1), scanline, int(curx2), scanline, col)
		curx1 -= invslope1
		curx2 -= invslope2
	}
}
