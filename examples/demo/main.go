package main

import (
	_ "embed"
	"etic"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var (
	//go:embed  spinner.png
	SpinnerPNG []byte
)

type Demos struct {
	etic.Cart
	mgr *etic.SceneTree
}

func NewSimple(title string) *Demos {
	return &Demos{Cart: *etic.NewCart(title, 320, 240)}
}

func (c *Demos) BOOT(t *etic.Console) {
	var aScene, bScene, cScene, dScene, eScene, fScene, gScene, hScene, iScene *etic.Element

	t.LoadSpriteSheet(etic.LoadImage(SpinnerPNG), 24, 24)

	{ // scene print time timer
		spinner := func(x, y, ln int, dur time.Duration) *etic.Element {
			var (
				index int
				dt    time.Duration
			)
			s := etic.NewElement("Spinner")
			s.OnInit = func(t *etic.Console) {}
			s.OnUpdate = func(t *etic.Console) {
				dt += t.Tick()
				if dt > dur {
					index = (index + 1) % ln
					dt = 0
				}
			}
			s.OnDraw = func(t *etic.Console) {
				t.Spr(etic.SpriteID(strconv.Itoa(index)), float32(x), float32(y), etic.SprOpt{})
			}

			return s
		}
		aScene = etic.NewElement("a.scene")
		aScene.Add(spinner(1, 100, 12, 100*time.Millisecond))

		aScene.OnInit = func(t *etic.Console) {
			log.Println("Enter a Scene")
		}
		aScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "b.scene")
			}
		}
		aScene.OnDraw = func(t *etic.Console) {
			t.Cls(colornames.Brown)
			t.Print("Hello World!", 1, 1, colornames.Yellow)
			tm := t.Time()
			t.Print(fmt.Sprintf("Seconds elapsed:%v", tm/1000), 1, 25, colornames.Orange)
			if tm%500 > 250 {
				t.Print("Warning!", 1, 50, colornames.Red)
			}
			if tm > 2000 {
				t.Print("Fugit inreparabile tempus", 1, 75, colornames.Fuchsia)
			}
		}
	}
	{ // scene label, button, mouse cursor
		cursor := func() *etic.Element {
			var (
				x, y, r float32
				p       bool
				s       string
			)
			nCursor := etic.NewElement("mouse.cursor")
			nCursor.OnUpdate = func(t *etic.Console) {
				x, y, p, _, _ = t.Mouse()
				if p {
					r += 2
				}
				r--
				r = max(0, min(8, r))
				s = fmt.Sprintf("(%v,%v) %v", x, y, r)
			}
			nCursor.OnDraw = func(t *etic.Console) {
				t.Line(x, 0, x, t.Height, colornames.Aqua)
				t.Line(0, y, t.Width, y, colornames.Aqua)
				t.Circ(x, y, r, colornames.Aqua)
				t.Print(s, 1, 1, colornames.Yellow)
			}
			return nCursor
		}
		bScene = etic.NewElement("b.scene")
		bScene.OnInit = func(t *etic.Console) {
			log.Println("Enter b Scene")
		}
		bScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "c.scene")
			}
		}
		panel := etic.PanelGrid("bScene.panel", 10, colornames.Gray, colornames.Black)
		lbl := etic.NewLabel("bScene.label", "Test")
		btn := etic.NewButton("bScene.button", "Click me", func(b *etic.Button) {
			log.Println("Clicked")
		})
		bScene.OnDraw = func(t *etic.Console) {
			t.Cls(colornames.Navy)
		}
		panel.Add(lbl)
		panel.Add(btn)
		bScene.Add(panel)
		bScene.Add(cursor())
		panel.Bounds().Set(etic.Rect(float32(10), 10, 300, 180))
		lbl.Bounds().Set(etic.Rect(float32(10), 10, 110, 30))
		btn.Bounds().Set(etic.Rect(float32(10), 40, 110, 60))
	}
	{ // colors scene
		normalizeColor := func(c color.Color) color.RGBA {
			r, g, b, a := c.RGBA()
			return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
		}
		nameToIdx := make(map[color.RGBA]int, len(colornames.Names))
		for i, name := range colornames.Names {
			nameToIdx[normalizeColor(colornames.Map[name])] = i
		}
		ready := false
		cScene = etic.NewElement("c.scene")
		cScene.OnInit = func(t *etic.Console) {
			ready = false
			log.Println("Enter e Scene")
		}
		cScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "d.scene")
			}
		}
		cScene.OnDraw = func(t *etic.Console) {
			if !ready {
				t.Cls(colornames.Black)
				for i := 0; i < 15; i++ {
					v := float32(i)
					t.Rect(9*v, 6*v, 6*v, 3*v, colornames.Map[colornames.Names[i]])
				}
				ready = true
			}
			if t.Frame()%12 == 0 {
				for x := 0; x < int(t.Width); x++ {
					for y := 0; y < int(t.Height); y++ {
						col := t.Pix(x, y)
						idx, ok := nameToIdx[normalizeColor(col)]
						if !ok {
							continue
						}
						idx = (idx + 1) % len(colornames.Names)
						t.SetPix(float32(x), float32(y), colornames.Map[colornames.Names[idx]])
					}
				}
			}
		}
	}
	{ // balls scene
		ball := func(x, y, dx, dy, r float32, c int) *etic.Element {
			b := etic.NewElement("ball")
			b.OnUpdate = func(t *etic.Console) {
				x += dx
				y += dy
				if x > t.Width-r {
					x = t.Width - r + 1
					dx = -dx
				} else if x < r {
					x = r
					dx = -dx
				} else if y > t.Height-r {
					y = t.Height - r + 1
					dy = -dy
				} else if y < r {
					y = r
					dy = -dy
				}
			}
			b.OnDraw = func(t *etic.Console) {
				t.Circ(x, y, r, colornames.Map[colornames.Names[c]])
				t.Circ(x+r/4, y-r/4, r/4, colornames.Map[colornames.Names[c+7]])
			}
			return b
		}
		dScene = etic.NewElement("d.scene")
		dScene.OnInit = func(t *etic.Console) {
			log.Println("Enter i Scene")
			var newBallData func() []int
			newBallData = func() []int {
				d := 1
				arr := make([]int, 0)
				arr = append(arr, rand.Intn(220)+10)
				arr = append(arr, rand.Intn(126)+10)
				arr = append(arr, rand.Intn(3)+1*d)
				arr = append(arr, rand.Intn(3)+1*d)
				arr = append(arr, rand.Intn(12)+6)
				arr = append(arr, rand.Intn(len(colornames.Names)-10))
				if arr[2] == 0 && arr[3] == 0 {
					arr = newBallData()
				}
				d = d * -1
				return arr
			}
			for range 50 {
				b := newBallData()
				x, y, dx, dy, r, c := b[0], b[1], b[2], b[3], b[4], b[5]
				dScene.Add(ball(float32(x), float32(y), float32(dx), float32(dy), float32(r), c))
			}
		}
		dScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "e.scene")
			}
		}
		dScene.OnDraw = func(t *etic.Console) {
			t.Cls(colornames.Navy)
		}
	}
	{ // scene test ellipses
		var (
			cx, cy           float32     = t.Width / 2, t.Height / 2
			colRed, colGreen color.Color = colornames.Red, colornames.Green
			Min, Max         float32     = 8, 48
			vRad, hRad       float32     = 16, 16
		)
		eScene = etic.NewElement("e.scene")
		aScene.OnInit = func(t *etic.Console) {
			log.Println("Enter e Scene")
		}
		eScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "f.scene")
			}
			if t.BtnP(ebiten.KeyUp) && vRad < Max {
				vRad++
			}
			if t.BtnP(ebiten.KeyDown) && vRad > Min {
				vRad--
			}
			if t.BtnP(ebiten.KeyLeft) && hRad < Max {
				hRad++
			}
			if t.BtnP(ebiten.KeyRight) && hRad > Min {
				hRad--
			}
		}
		eScene.OnDraw = func(t *etic.Console) {
			t.Cls(colornames.Brown)
			t.Elli(cx, cy, hRad, vRad, colRed)
			t.ElliB(cx, cy, vRad, hRad, colGreen)
			t.Print("Arrow to test ellise", 1, 1, colornames.Aqua)
		}
	}

	{ // scene test triangles
		var (
			x, y, cx, cy float64
		)
		pir := func(x, y, w, h, cx, cy float32) {
			t.Tri(x, y, w/2+cx, h/2+cy, x+w, y, colornames.Navy)
			t.Tri(x+w, y, w/2+cx, h/2+cy, x+w, y+h, colornames.Blue)
			t.Tri(x, y, w/2+cx, h/2+cy, x, y+h, colornames.Blueviolet)
			t.Tri(x, y+h, w/2+cx, h/2+cy, x+w, y+h, colornames.Aqua)
		}
		fScene = etic.NewElement("f.scene")
		fScene.OnInit = func(t *etic.Console) {
			log.Println("Enter f Scene")
		}
		fScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "g.scene")
			}
		}
		fScene.OnDraw = func(t *etic.Console) {
			t.Cls(colornames.Gray)
			for x = 0; x < float64(t.Width); x += 28 {
				for y = 0; y < float64(t.Height); y += 28 {
					cx = 12 * math.Sin(float64(float64(t.Time())/30000*(x+y+1)))
					cy = 12 * math.Cos(float64(float64(t.Time())/30000*(x+y+1)))
					pir(float32(x), float32(y), 25, 25, float32(x+cx), float32(y+cy))
				}
			}
		}
	}
	{ // scene clip
		var (
			x, y float32 = 96, 24
			clip bool    = false
		)
		grid := func(rect etic.Rectangle[float32], sp, st float32, fg, bg color.Color) {
			t.Rect(rect.Min.X, rect.Min.Y, rect.Dx(), rect.Dy(), bg)
			r := rect.Inset(st)
			x0, y0, w, h := etic.RectF32(r)
			x := x0
			y := y0
			for x <= x0+w {
				t.Line(x, y0, x, y0+h, fg)
				x += float32(sp)
			}
			x = x0
			for y <= y0+h {
				t.Line(x, y, x0+w, y, fg)
				y += float32(sp)
			}
		}
		gScene = etic.NewElement("g.scene")
		gScene.OnInit = func(t *etic.Console) {
			log.Println("Enter g Scene")
		}
		gScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "h.scene")
			}
			if t.BtnP(ebiten.KeyUp) {
				y--
				clip = !clip
			}
			if t.BtnP(ebiten.KeyDown) {
				y++
			}
			if t.BtnP(ebiten.KeyLeft) {
				x--
			}
			if t.BtnP(ebiten.KeyRight) {
				x++
			}
		}
		gScene.OnDraw = func(t *etic.Console) {
			t.Cls(colornames.Navy)
			if clip {
				t.Cls(colornames.Black)
				t.Clip(60, 20, 100, 100)
				t.Cls(colornames.Blue)
			} else {
				t.ResetClip()
			}
			grid(etic.Rect(0, 0, t.Width, t.Height), 10, 1, colornames.Gray, colornames.Yellow)
			t.Print("Press Up To", 76, 84, colornames.Blue)
			t.Print("Toggle Clipping", 72, 94, colornames.Fuchsia)
			t.Rect(x, y, 8, 8, colornames.Red)
		}
	}

	{
		hScene = etic.NewElement("h.scene")
		hScene.OnInit = func(t *etic.Console) {
			log.Println("Enter h Scene")
		}
		hScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "i.scene")
			}
		}

		draw := func(dt int64) {
			for y := 0; y < int(t.Height); y++ {
				for x := 0; x < int(t.Width); x++ {
					idx := (x + y + int(dt)) / 8 % len(colornames.Map)
					color := colornames.Map[colornames.Names[idx]]
					t.SetPix(float32(x), float32(y), color)
				}
			}
		}
		hScene.OnDraw = func(t *etic.Console) {
			draw(t.Time() / 19)
		}
	}

	{
		var x, y, i float64
		var (
			HALF_SCR_W = float64(t.Width) / 2.0
			HALF_SCR_H = float64(t.Height) / 2.0
			DEVIATION  = 150.0
			SPEED      = 1 / 500.0
			RECT_COUNT = 70.0
			RECT_STEP  = 4.0
			RECT_COLOR = colornames.Yellow
		)

		iScene = etic.NewElement("i.scene")
		iScene.OnInit = func(t *etic.Console) {
			log.Println("Enter i Scene")
		}
		iScene.OnUpdate = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEnter) {
				c.mgr.Change(t, "a.scene")
			}
		}
		iScene.OnDraw = func(t *etic.Console) {
			t.Cls(colornames.Navy)
			for i = 1.0; i < RECT_COUNT; i++ {
				width := i * RECT_STEP
				height := width / 2
				slowedTime := float64(t.Time()) * SPEED
				x = math.Sin(slowedTime)*DEVIATION/i - width/2
				y = math.Cos(slowedTime)*DEVIATION/i - height/2
				t.RectB(float32(HALF_SCR_W+x), float32(HALF_SCR_H+y), float32(width), float32(height), RECT_COLOR)
			}
		}
	}

	c.mgr = &etic.SceneTree{}
	c.mgr.AddScene(aScene)
	c.mgr.AddScene(bScene)
	c.mgr.AddScene(cScene)
	c.mgr.AddScene(dScene)
	c.mgr.AddScene(eScene)
	c.mgr.AddScene(fScene)
	c.mgr.AddScene(gScene)
	c.mgr.AddScene(hScene)
	c.mgr.AddScene(iScene)
	c.mgr.Change(t, aScene.ID().Get())
}
func (c *Demos) TIC(t *etic.Console) {
	if t.Btn(ebiten.KeyEscape) {
		t.Exit()
	}
	c.mgr.TIC(t)
}

func main() { etic.Load(NewSimple("Demos")).Run() }
