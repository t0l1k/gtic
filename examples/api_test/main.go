package main

import (
	_ "embed"
	"fmt"
	"gtic"
	"gtic/ui"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/colornames"
)

var (
	//go:embed  spinner.png
	SpinnerPNG []byte
)

func LoadSprite(a *gtic.API) (ids []gtic.SpriteID) {
	surface, err := gtic.DecodeImage(SpinnerPNG)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 12; i++ {
		rect := image.Rect(i*24, 0, i*24+24, 24)
		id := a.Sprites().Register(
			gtic.NewSprite("Spinner+"+strconv.Itoa(i), surface, rect))
		ids = append(ids, id)
	}
	return ids
}

func CartClipTest(fn func()) *ui.Element {
	var (
		x, y int  = 96, 24
		clip bool = false
	)
	grid := func(a *gtic.API, rect image.Rectangle, sp, st int, fg, bg gtic.RGBA) {
		a.Rect(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(), bg)
		r := rect.Inset(st)
		x0, y0, w, h := gtic.RectXYWH(r)
		x := x0
		y := y0
		for x <= x0+w {
			a.Line(x, y0, x, y0+h, fg)
			x += sp
		}
		x = x0
		for y <= y0+h {
			a.Line(x, y, x0+w, y, fg)
			y += sp
		}
	}

	e := ui.NewElement("cart.clip.test")
	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
		if a.Key(gtic.KeyUP) {
			y--
		}
		if a.Key(gtic.KeyDOWN) {
			y++
		}
		if a.Key(gtic.KeyLEFT) {
			x--
		}
		if a.Key(gtic.KeyRIGHT) {
			x++
		}
		if a.KeyP(gtic.KeyUP) {
			clip = !clip
		}
	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls(a.Pal(9))
		if clip {
			a.Cls()
			a.Clip(60, 20, 100, 100)
			a.Cls(a.Pal(9))
		} else {
			a.Clip()
		}
		grid(a, image.Rect(0, 0, a.Bounds().Width, a.Bounds().Height), 10, 1, a.Pal(14), a.Pal(4))
		a.Print("Press Up To", 76, 84, a.Pal(8))
		a.Print("Toggle Clipping", 72, 94, a.Pal(1))
		a.Rect(x, y, 8, 8, a.Pal(3))
	}
	return e
}

func CartTriTest(fn func()) *ui.Element {
	pir := func(a *gtic.API, x, y, w, h, cx, cy float64) {
		a.Tri(x, y, w/2+cx, h/2+cy, x+w, y, a.Pal(8))
		a.Tri(x+w, y, w/2+cx, h/2+cy, x+w, y+h, a.Pal(9))
		a.Tri(x, y, w/2+cx, h/2+cy, x, y+h, a.Pal(10))
		a.Tri(x, y+h, w/2+cx, h/2+cy, x+w, y+h, a.Pal(11))
	}

	e := ui.NewElement("cart.tri")
	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls()
		for x := 0.0; x < float64(a.Bounds().Width); x += 28 {
			for y := 0.0; y < float64(a.Bounds().Height); y += 28 {
				factor := float64(a.Time()) / 30000 * (x + y + 1)
				cx := 12 * math.Sin(factor)
				cy := 12 * math.Cos(factor)
				pir(a, x, y, 25, 25, x+cx, y+cy)
			}
		}
	}
	return e
}

func CartEllipsTest(fn func()) *ui.Element {

	var (
		cx, cy     int
		Min, Max   int = 8, 48
		vRad, hRad int = 16, 16
	)

	draw := func(t *gtic.API, dt int64) {
		for y := 0; y < int(t.Bounds().Height); y++ {
			for x := 0; x < int(t.Bounds().Width); x++ {
				idx := (x + y + int(dt)) / 8 % len(colornames.Map)
				color := colornames.Map[colornames.Names[idx]]
				t.Pix(x, y, gtic.NewColor(color))
			}
		}
	}

	e := ui.NewElement("cart.ellips.test")
	e.OnInit = func(a *gtic.API) {
		cx, cy = a.Bounds().Width/2, a.Bounds().Height/2
	}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}

		if a.KeyP(gtic.KeyUP) && vRad < Max {
			vRad++
		}
		if a.KeyP(gtic.KeyDOWN) && vRad > Min {
			vRad--
		}
		if a.KeyP(gtic.KeyLEFT) && hRad < Max {
			hRad++
		}
		if a.KeyP(gtic.KeyRIGHT) && hRad > Min {
			hRad--
		}

	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls()
		draw(a, a.Time()/19)
		a.Elli(cx, cy, hRad, vRad, a.Pal(4))
		a.ElliB(cx, cy, vRad, hRad, a.Pal(5))
		a.Print("Arrow to test ellise", 1, 1, colornames.Aqua)
	}
	return e
}

func CartOne(fn func()) *ui.Element {
	e := ui.NewElement("cart.one")

	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		for range 5000 {
			x := rand.IntN(a.Bounds().Width)
			y := rand.IntN(a.Bounds().Height)
			a.Pix(x, y, 0)
		}
		t := float64(a.Time()) / 500
		r := float64(50)
		x := float64(a.Bounds().Width/2) + r*math.Cos(3*t)
		y := float64(a.Bounds().Height/2) + r*math.Sin(4*t)
		a.Circ(int(x), int(y), 7, gtic.NewColor(colornames.Orange))
	}
	return e
}

// Эффект в gtic80 реализуется с BDR
func CartRainbow(fn func()) *ui.Element {
	e := ui.NewElement("cart.four")

	background := func(a *gtic.API, line int) gtic.RGBA {
		t := float64(a.Time() / 200)
		l := float64(line)/float64(a.Bounds().Height)*math.Pi*2 - t
		r := math.Sin(l)*127 + 128
		g := math.Sin(l+math.Pi*2/3)*127 + 128
		b := math.Sin(l+math.Pi*2*2/3)*127 + 128
		return gtic.NewRGBA(uint8(r), uint8(g), uint8(b), 255)
	}

	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		for y := range a.Bounds().Height {
			c := background(a, y)
			for x := range a.Bounds().Width {
				a.Pix(x, y, c)
			}
		}
		a.Rect(10, 10, 100, 100, gtic.NewColor(colornames.Teal))
	}
	return e
}

func CartTextEffects(fn func()) *ui.Element {
	e := ui.NewElement("cart.text.effects")

	printOutline := func(a *gtic.API, msg string, x, y int, c1, c2 gtic.RGBA) {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				a.Print(msg, x+dx, y+dy, c2)
			}
		}
		a.Print(msg, x, y, c1)
	}
	printShadow := func(a *gtic.API, msg string, x, y int, c1, c2 gtic.RGBA) {
		dx, dy := 2, 2
		a.Print(msg, x+dx, y+dy, c2)
		a.Print(msg, x, y, c1)
	}
	printStripes := func(a *gtic.API, msg string, x, y int) {
		for dy := 0; dy <= 8; dy++ {
			a.Clip(0, y+dy, a.Bounds().Width, 1)
			col := gtic.NewRGBA(uint8((dy*30)%255), uint8((dy*60)%255), uint8((dy*90)%255), 255)
			a.Print(msg, x, y, col)
		}
		a.Clip()
	}

	mod := func(a, b int) int {
		res := a % b
		if res < 0 {
			return res + b
		}
		return res
	}
	circp := func(api *gtic.API, x, y, r int, c gtic.RGBA, p ...int) {
		// Функция рисования круга с паттерном
		// Если p не передано, в Go мы задаем дефолтное значение до вызова,
		// либо инициализируем внутри (0xffff эквивалентно 65535)
		pattern := 0xffff
		if len(p) > 0 {
			pattern = p[0]
		}
		for a := x - r; a <= x+r; a++ {
			for b := y - r; b <= y+r; b++ {
				// Вычисляем расстояние
				dx := float64(a - x)
				dy := float64(b - y)
				d := math.Sqrt(dx*dx + dy*dy)

				if d <= float64(r) {
					// В Lua: ((b-y)//1)%4. В Go для целых чисел достаточно использовать mod.
					// Паттерн 4х4: определяем индекс бита от 0 до 15
					bitY := mod(a-x, 4)
					bitX := mod(b-y, 4)
					bit := uint(bitX*4 + bitY)

					// Проверяем, установлен ли бит в маске p (аналог p & (2^bit) != 0)
					if (pattern & (1 << bit)) != 0 {
						api.Pix(a, b, c)
					} else {
						// В оригинальном коде Lua параметра c2 нет в условиях,
						// но если он нужен для заднего фона паттерна, его можно вызвать тут:
						// pix(a, b, c2)
					}
				}
			}
		}
	}

	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls()
		t := float64(a.Time()) / 200
		circp(a, 120, 68, 30, gtic.NewColor(colornames.Yellow), 0xada7)
		circp(a, 100, 68+int(20*math.Sin(t)), 20, gtic.NewColor(colornames.Red), 0xaeab)
		circp(a, 120+int(60*math.Cos(t/2)), 40, 20, gtic.NewColor(colornames.Orange), 0x8520)

		printOutline(a, "Outline!", 50, 50, gtic.NewColor(colornames.Yellow), gtic.NewColor(colornames.Navajowhite))
		printShadow(a, "Shadow!", 50, 70, gtic.NewColor(colornames.Aliceblue), gtic.NewColor(colornames.Dimgray))
		printStripes(a, "Stripes!", 50, 90)
	}
	return e
}

func CartBalls(fn func()) *ui.Element {
	e := ui.NewElement("cart.balls")

	ball := func(x, y, dx, dy, r int, c1, c2 gtic.RGBA) *ui.Element {
		b := ui.NewElement("ball")
		b.OnUpdate = func(t *gtic.API) {
			x += dx
			y += dy
			if x > t.Bounds().Width-r {
				x = t.Bounds().Width - r + 1
				dx = -dx
			} else if x < r {
				x = r
				dx = -dx
			} else if y > t.Bounds().Height-r {
				y = t.Bounds().Height - r + 1
				dy = -dy
			} else if y < r {
				y = r
				dy = -dy
			}
		}
		b.OnDraw = func(t *gtic.API) {
			t.Circ(x, y, r, c1)
			t.Circ(x+r/4, y-r/4, r/4, c2)
		}
		return b
	}
	e.OnInit = func(t *gtic.API) {
		log.Println("Enter i Scene")

		randomColor := func() gtic.RGBA {
			keys := make([]string, 0, len(colornames.Map))
			for name := range colornames.Map {
				keys = append(keys, name)
			}
			randomIndex := rand.IntN(len(keys))
			randomName := keys[randomIndex]
			var randomColor color.Color = colornames.Map[randomName]
			return gtic.NewColor(randomColor)
		}

		var newBallData func() []int
		newBallData = func() []int {
			d := 1
			arr := make([]int, 0)
			arr = append(arr, rand.IntN(220)+10)
			arr = append(arr, rand.IntN(126)+10)
			arr = append(arr, rand.IntN(3)+1*d)
			arr = append(arr, rand.IntN(3)+1*d)
			arr = append(arr, rand.IntN(12)+6)
			if arr[2] == 0 && arr[3] == 0 {
				arr = newBallData()
			}
			d = d * -1
			return arr
		}
		for range 50 {
			b := newBallData()
			x, y, dx, dy, r := b[0], b[1], b[2], b[3], b[4]
			e.Add(ball(x, y, dx, dy, r, randomColor(), randomColor()))
		}
	}
	e.OnUpdate = func(t *gtic.API) {
		if t.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(t *gtic.API) {
		t.Cls(gtic.NewColor(colornames.Navy))
	}
	return e
}

func CartSystemFontDemo(fn func()) *ui.Element {
	const (
		offy  = 4
		offx  = 65
		space = 10
	)
	var (
		aa = 0.0
	)
	circB := func(a *gtic.API) {
		y1 := float64(a.Bounds().Height) / 2.0
		y2 := float64(a.Bounds().Height) / 3.4
		w1 := float64(a.Bounds().Width) / 1.2
		w2 := float64(a.Bounds().Width) / 3
		for i := 0; i < int(w1); i += space {
			x := (w1 - w2) + w2*math.Sin(aa)
			y := y1 + y2*math.Cos(aa)
			r := i + int(a.Time()/int64(y2)%space)
			a.CircB(int(x), int(y), r, a.Pal(8))
			x = (w1 - w2) + w2*math.Sin(aa)
			y = y1 + y2*math.Cos(aa/2)
			r = i + int(a.Time()/int64(y2)%space)
			a.CircB(int(x), int(y), r, a.Pal(8))
		}
		aa += math.Pi / float64(a.Bounds().Width)
	}

	lpad := func(s string, l int, c byte) string {
		if len(s) >= l {
			return s
		}
		return strings.Repeat(string(c), l-len(s)) + s
	}
	fontDemo := func(a *gtic.API) {
		a.Print("SYSTEM\nFONT", 0, 5, a.Pal(12))
		a.Print("0x14604", 0, 21, a.Pal(10), true, 1, true)
		for x := 0; x < 16; x++ {
			off := x * 16
			a.Print(lpad(strconv.Itoa(off), 3, ' '),
				offx-20, x*8+offy+1, a.Pal(14),
				true, 1, true) // <-- было (true, true)
			for y := 0; y < 16; y++ {
				ch := rune(y*16 + x)
				a.Rect(x*11+offx, y*8+offy, 8, 7, a.Pal(15))
				a.Print(string(ch), x*11+offx, y*8+offy, a.Pal(12))
			}
		}
	}
	e := ui.NewElement("cart.system.font.demo")
	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls()
		circB(a)
		fontDemo(a)
	}
	return e
}

func CartSpr(fn func()) *ui.Element {
	type Param struct{ val, min, max int }
	var (
		params       map[string]*Param
		paramOrder   []string
		selIdx       int
		demoSpriteID gtic.SpriteID
	)
	params = map[string]*Param{
		"x":        {val: 100, min: 0, max: 240},
		"y":        {val: 68, min: 0, max: 136},
		"colorkey": {val: 0, min: -1, max: 15},
		"scale":    {val: 1, min: 1, max: 10},
		"flip":     {val: 0, min: 0, max: 3},
		"rotate":   {val: 0, min: 0, max: 3},
	}
	paramOrder = []string{"x", "y", "colorkey", "scale", "flip", "rotate"}

	e := ui.NewElement("cart.spr")
	e.OnInit = func(a *gtic.API) {
		data := "D000000D" + "D0C00C0D" + "DDDDDDDD" + "DDCDDD2D" + "DCCCDDDD" + "DDcdd2dd" + "DDDDDDDD" + "EEEEEEEE"
		sprData := func() (surface gtic.Surface[gtic.RGBA]) {
			pixels := make([]gtic.RGBA, 0)
			for i := 0; i < 64; i++ {
				val, _ := strconv.ParseUint(string(data[i]), 16, 4)
				colIdx := int(val)
				pixels = append(pixels, a.Pal(colIdx))
			}
			return gtic.Surface[gtic.RGBA]{Data: pixels, Width: 8, Height: 8}
		}
		demoSpriteID = a.Sprites().Register(gtic.NewSprite("test.sprite", sprData(), image.Rect(0, 0, 8, 8)))
	}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}

		// 1. Управление и обновление состояния
		if a.Btnp(0, 30, 6) {
			selIdx = (selIdx - 1 + len(paramOrder)) % len(paramOrder)
		}
		if a.Btnp(1, 30, 6) {
			selIdx = (selIdx + 1) % len(paramOrder)
		}

		selKey := paramOrder[selIdx]
		p := params[selKey]

		if a.Btnp(2, 30, 6) {
			p.val--
		}
		if a.Btnp(3, 30, 6) {
			p.val++
		}

		// Clamp
		if p.val > p.max {
			p.val = p.max
		}
		if p.val < p.min {
			p.val = p.min
		}

	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls()
		a.Print("Use up/down to select parameter", 0, 0)
		a.Print("left/right to change its value:", 0, 10)
		for r, key := range paramOrder {
			cur := " "
			if r == selIdx {
				cur = ">"
			}
			text := fmt.Sprintf("%s%s:%d", cur, key, params[key].val)
			a.Print(text, 0, 30+10*r)
		}

		ckIdx := params["colorkey"].val
		var ckColor gtic.RGBA = gtic.Transparent
		if ckIdx >= 0 {
			ckColor = a.Pallete()[ckIdx]
		}
		a.Spr(demoSpriteID,
			params["x"].val,
			params["y"].val,
			ckColor,
			params["scale"].val,
			params["flip"].val,
			params["rotate"].val,
		)
	}
	return e
}

func CartPixA(fn func()) *ui.Element {
	e := ui.NewElement("cart.pix.a")
	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		a.Print("Pix A test")
		for i := 0; i < 6000; i++ {
			second := int(a.Time()) / 1000
			x := rand.IntN(a.Bounds().Width)
			y := rand.IntN(a.Bounds().Height)
			color := (second * x * y) % 60
			a.Pix(int(x), int(y), a.Pal(color))
		}
	}
	return e
}

func CartPixB(fn func()) *ui.Element {
	e := ui.NewElement("cart.pix.b")

	background := func(a *gtic.API) {
		a.Cls()
		for i := 0; i < len(a.Pallete()); i++ {
			a.Rect(9*i, 6*i, 6*i, 3*i, a.Pal(i))
		}
	}

	e.OnInit = func(a *gtic.API) {
		background(a)
	}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		if a.Frame()%12 != 0 {
			return
		}
		a.Print("Pix B test")
		for x := 0; x < a.Bounds().Width; x += 2 {
			for y := 0; y < a.Bounds().Height; y += 2 {
				idx := a.ParsePal(a.Pix(x, y))
				idx = (idx + 1) % len(a.Pallete())
				a.Pix(x, y, a.Pal(idx))
			}
		}
	}
	return e
}

func CartMouseTest(fn func()) *ui.Element {
	e := ui.NewElement("cart.mouse.test")
	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	GR, YE, WH := 6, 4, 12
	barX, barY := 10, 10
	e.OnDraw = func(a *gtic.API) {
		x, y, left, middle, right, sX, sY := a.Mouse()
		barX += int(sX)
		barY += int(sY)
		if barX < 1 {
			barX = 1
		}
		if barY < 1 {
			barY = 1
		}

		a.Cls()
		a.Print("Move Mouse", 10, 10, a.Pal(YE))
		a.Print(fmt.Sprintf("x=%v,y=%v", x, y), 100, 10, a.Pal(WH))
		a.Print("Press Buttons:", 10, 20, a.Pal(YE))
		a.Print(fmt.Sprintf("Left %v", left), 100, 20, a.Pal(WH))
		a.Print(fmt.Sprintf("Middle %v", middle), 100, 40, a.Pal(WH))
		a.Print(fmt.Sprintf("Right %v", right), 100, 30, a.Pal(WH))
		a.Print("Scroll Wheel:", 10, 80, a.Pal(YE))
		a.Print("Scroll X", 100, 80, a.Pal(WH))
		a.Print("Scroll Y", 160, 80, a.Pal(WH))
		a.Rect(100, 136/2-barX, 8, barX, a.Pal(GR))
		a.Rect(160, 136/2-barY, 8, barY, a.Pal(GR))
	}
	return e
}

func CartPrintScaleTest(fn func()) *ui.Element {
	e := ui.NewElement("cart.print.scale.test")
	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
			a.ResetPal()
		}
	}
	var t float64
	e.OnDraw = func(a *gtic.API) {
		txt := "[gtic-80]"
		for i := 1; i < 16; i++ {
			r := uint8(100 + i*10)
			colG := uint8(i * 5)
			b := uint8(32 + 32*math.Sin(t/100.0))
			col := gtic.NewRGBA(r, colG, b, 255)
			a.Pal(i, col)
		}
		a.Cls()
		//  Рендеринг 3D-эха текста сзади вперед
		for z := 15; z > 0; z-- {
			// Рассчитываем волну смещения по вертикали
			yOffset := 12 * float64(z) * math.Sin(float64(z)/5.0+t/50.0)
			// Узнаем ширину текста при текущем масштабе `z`.
			// Чтобы не рисовать на экран, мы можем временно зажать Clip в нулевой размер
			a.Clip(0, 0, 0, 0)
			textWidth := a.Print(txt, 0, -100, a.Pal(z), z)
			a.Clip() // Восстанавливаем экран
			// Центрируем текст по горизонтали относительно ширины экрана (Width)
			posX := (a.Bounds().Width - textWidth) / 2
			posY := int(68.0 + yOffset/1.5)
			// Выводим текст
			a.Print(txt, posX, posY, a.Pal(z), z)
		}
		t += 1 // Инкремент таймера
	}
	return e
}

func CartPrintTest(fn func()) *ui.Element {
	e := ui.NewElement("cart.print.test")
	var (
		msg   = "my perfectly centered text"
		width int
	)

	e.OnInit = func(a *gtic.API) {
		width = a.Print(msg)
	}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls()
		x, y := (a.Bounds().Width-width)/2, (a.Bounds().Height-6)/2
		a.Print(msg, x, y, gtic.NewColor(colornames.Yellow))

		a.Print("FIXED", 0, 0, gtic.NewColor(colornames.White), true)
		a.Print("FIXED", 0, 8, gtic.NewColor(colornames.White), false)
		a.Print("width", 0, 32, gtic.NewColor(colornames.White), true)
		a.Print("width", 0, 40, gtic.NewColor(colornames.White), false)

		for i := 0; i < 30; i += 6 {
			a.Line(i, 0, i, 6, gtic.NewColor(colornames.Blue))
			a.Line(i, 8, i, 16, gtic.NewColor(colornames.Orange))
			a.Line(i, 32, i, 40, gtic.NewColor(colornames.Blue))
			a.Line(i, 40, i, 48, gtic.NewColor(colornames.Orange))
		}
	}
	return e
}

func CartHello(fn func()) *ui.Element {
	var (
		sprIds []gtic.SpriteID
	)
	spinner := func(x, y, ln int, dur time.Duration) *ui.Element {
		var (
			index int
			dt    time.Duration
		)

		s := ui.NewElement("Spinner")
		s.OnInit = func(t *gtic.API) {}
		s.OnUpdate = func(t *gtic.API) {
			dt += t.Delta()
			if dt > dur {
				index = (index + 1) % ln
				dt = 0
			}
		}
		s.OnDraw = func(t *gtic.API) {
			t.Spr(sprIds[index], x, y, t.Pal(16))
		}

		return s
	}

	grid := 8
	e := ui.NewElement("cart.hello")
	e.Add(spinner(8, 9*grid, 12, 100*time.Millisecond))
	e.OnInit = func(a *gtic.API) {
		sprIds = LoadSprite(a)
		log.Println("pal", a.Pallete())
	}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn()
		}
	}
	e.OnDraw = func(a *gtic.API) {
		a.Cls(gtic.NewColor(colornames.Teal))
		for i := 0; i < a.Bounds().Height; i++ {
			a.Line(0, i*grid, a.Bounds().Width, i*grid, a.Pal(1))
			a.Line(i*grid, 0, i*grid, a.Bounds().Height, a.Pal(1))
		}
		a.Print("Hello, World!!!", 0, 0, gtic.NewColor(colornames.Yellow), false, 2)
		tm := a.Time()
		a.Print(fmt.Sprintf("Seconds elapsed:%v", tm/1000), 0, 3*grid, colornames.Orange)
		if tm%500 > 250 {
			a.Print("Warning!", 0, 5*grid, colornames.Red)
		}
		if tm > 2000 {
			a.Print("Fugit inreparabile tempus", 0, 7*grid, colornames.Fuchsia)
		}
	}
	return e
}

func main() {

	var (
		ST             *ui.SceneTree
		hello          *ui.Element
		mouseTest      *ui.Element
		printScaleTest *ui.Element
		printTest      *ui.Element
		pixATest       *ui.Element
		pixBTest       *ui.Element
		sprTest        *ui.Element
		systemFontDemo *ui.Element
		ballsDemo      *ui.Element
		textEffects    *ui.Element
		rainbow        *ui.Element
		cartOne        *ui.Element
		ellipsTest     *ui.Element
		triTest        *ui.Element
		clipTest       *ui.Element
	)

	gtic.BOOT = func(a *gtic.API) {
		ST = ui.NewSceneTree()
		hello = CartHello(func() { ST.Change(a, mouseTest.Property("ID").Get().(ui.ElementID)) })
		mouseTest = CartMouseTest(func() { ST.Change(a, printTest.Property("ID").Get().(ui.ElementID)) })
		printTest = CartPrintTest(func() { ST.Change(a, printScaleTest.Property("ID").Get().(ui.ElementID)) })
		printScaleTest = CartPrintScaleTest(func() { ST.Change(a, pixATest.Property("ID").Get().(ui.ElementID)) })
		pixATest = CartPixA(func() { ST.Change(a, pixBTest.Property("ID").Get().(ui.ElementID)) })
		pixBTest = CartPixB(func() { ST.Change(a, sprTest.Property("ID").Get().(ui.ElementID)) })
		sprTest = CartSpr(func() { ST.Change(a, systemFontDemo.Property("ID").Get().(ui.ElementID)) })
		systemFontDemo = CartSystemFontDemo(func() { ST.Change(a, ballsDemo.Property("ID").Get().(ui.ElementID)) })
		ballsDemo = CartBalls(func() { ST.Change(a, textEffects.Property("ID").Get().(ui.ElementID)) })
		textEffects = CartTextEffects(func() { ST.Change(a, rainbow.Property("ID").Get().(ui.ElementID)) })
		rainbow = CartRainbow(func() { ST.Change(a, cartOne.Property("ID").Get().(ui.ElementID)) })
		cartOne = CartOne(func() { ST.Change(a, ellipsTest.Property("ID").Get().(ui.ElementID)) })
		ellipsTest = CartEllipsTest(func() { ST.Change(a, triTest.Property("ID").Get().(ui.ElementID)) })
		triTest = CartTriTest(func() { ST.Change(a, clipTest.Property("ID").Get().(ui.ElementID)) })
		clipTest = CartClipTest(func() { ST.Change(a, hello.Property("ID").Get().(ui.ElementID)) })
		ST.AddScene(hello)
		ST.AddScene(mouseTest)
		ST.AddScene(printTest)
		ST.AddScene(printScaleTest)
		ST.AddScene(pixATest)
		ST.AddScene(pixBTest)
		ST.AddScene(sprTest)
		ST.AddScene(systemFontDemo)
		ST.AddScene(ballsDemo)
		ST.AddScene(textEffects)
		ST.AddScene(rainbow)
		ST.AddScene(cartOne)
		ST.AddScene(ellipsTest)
		ST.AddScene(triTest)
		ST.AddScene(clipTest)

		ST.Change(a, hello.Property("ID").Get().(ui.ElementID))
	}
	gtic.TIC = func(a *gtic.API) {
		if a.KeyP(gtic.KeyESC) {
			a.Exit()
		}
		ST.TIC(a)
	}

	if err := gtic.Load(gtic.WithTitle("Api Test Examples"), gtic.WithMode(320, 200), gtic.WithScale(2)).Run(); err != nil {
		log.Println(err)
	}
}
