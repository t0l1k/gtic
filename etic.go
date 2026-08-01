package etic

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

type SpriteID string

type Console struct {
	Width, Height         float32
	screen                *ebiten.Image
	sprites               map[SpriteID]*ebiten.Image
	started, lastTick     time.Time
	delta                 time.Duration
	frame                 int
	system, normal, title *text.GoTextFace
	booted, quit          bool
	camX, camY            float32
	clip                  Rectangle[float32]
	spriteBounds          Rectangle[float32]
	focusedId, activeId   ElementID
}

func loadConsole(w, h float32) *Console {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}
	api := &Console{
		Width:   w,
		Height:  h,
		system:  &text.GoTextFace{Source: source, Size: 10},
		normal:  &text.GoTextFace{Source: source, Size: 20},
		title:   &text.GoTextFace{Source: source, Size: 50},
		started: time.Now(),
		screen:  ebiten.NewImage(int(w), int(h)),
		clip:    Rect(0, 0, w, h),
	}
	return api
}

func (t *Console) Screen() *ebiten.Image { return t.screen }
func (t *Console) Pix(x, y int) color.Color {
	if !image.Pt(x, y).In(t.screen.Bounds()) {
		return color.Transparent
	}
	return t.screen.At(x, y)
}

func (t *Console) SetPix(x, y float32, c color.Color) {
	x += t.camX
	y += t.camY
	r := t.screen.Bounds()
	if !Pt(x, y).In(RectIntToRectF32(r)) {
		return
	}
	if !Pt(x, y).In(t.clip) {
		return
	}
	t.screen.Set(int(x), int(y), c)
}

func (t *Console) Line(x0, y0, x1, y1 float32, c color.Color) {
	vector.StrokeLine(t.screen, float32(x0), float32(y0), float32(x1), float32(y1), 1, c, false)
}

func (t *Console) Rect(x, y, w, h float32, c color.Color) {
	vector.FillRect(t.screen, float32(x), float32(y), float32(w), float32(h), c, false)
}

func (t *Console) RectB(x, y, w, h float32, c color.Color) {
	vector.StrokeRect(t.screen, float32(x), float32(y), float32(w), float32(h), 1, c, false)
}

func (t *Console) Circ(x, y, r float32, c color.Color) {
	vector.FillCircle(t.screen, float32(x), float32(y), float32(r), c, false)
}

func (t *Console) CircB(x, y, r float32, c color.Color) {
	vector.StrokeCircle(t.screen, float32(x), float32(y), float32(r), 1, c, false)
}

func (t *Console) Cls(c color.Color) { t.screen.Fill(c) }

// Print(text,x,y,color,[scale=1]) -> width, height
// scale системный шрифт=10*scale
func (t *Console) Print(msg string, x, y float32, c color.Color, args ...float32) (float32, float32) {
	w, h := text.Measure(msg, t.system, t.system.Size*1.2)
	scale := 1.0
	if len(args) > 0 {
		scale = float64(args[0])
	}
	t.DrawString(msg, FontSystem, Rect(x, y, x+float32(w), y+float32(h)), text.AlignStart, text.AlignStart, c, scale)
	return float32(w), float32(h)
}

func (t *Console) Btn(id ebiten.Key) bool { return ebiten.IsKeyPressed(id) }

func (t *Console) BtnP(id ebiten.Key) bool { return inpututil.IsKeyJustPressed(id) }

func (t *Console) CharsP() ([]rune, bool) {
	chars := ebiten.AppendInputChars(nil)
	return chars, len(chars) > 0
}

func (t *Console) Mouse() (float32, float32, bool, bool, bool) {
	x, y := ebiten.CursorPosition()
	left := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	middle := ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	right := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	return float32(x), float32(y), left, middle, right
}

func (t *Console) MouseP(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}

func (t *Console) Time() int64         { return time.Since(t.started).Milliseconds() }
func (t *Console) Tick() time.Duration { return t.delta }
func (t *Console) Frame() int          { return t.frame }

type FontKind int

const (
	FontSystem FontKind = iota
	FontNormal
	FontTitle
)

func (t *Console) Face(kind FontKind) *text.GoTextFace {
	switch kind {
	case FontSystem:
		return t.system
	case FontTitle:
		return t.title
	default:
		return t.normal
	}
}

func (t *Console) DrawString(txt string, face FontKind, rect Rectangle[float32], hAlign, vAlign text.Align, fg color.Color, args ...float64) {
	var (
		x, y, w, h float64
	)
	x, y = float64(rect.Min.X), float64(rect.Min.Y)
	w, h = float64(rect.Dx()), float64(rect.Dy())
	switch hAlign {
	case text.AlignStart:
		x += 0
	case text.AlignCenter:
		x += w / 2
	case text.AlignEnd:
		x += w
	}
	switch vAlign {
	case text.AlignStart:
		y += 0
	case text.AlignCenter:
		y += h / 2
	case text.AlignEnd:
		y += h
	}
	op := &text.DrawOptions{}
	scale := 1.0
	if len(args) > 0 {
		scale = args[0]
	}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(fg)
	op.LayoutOptions.PrimaryAlign = hAlign
	op.LayoutOptions.SecondaryAlign = vAlign
	op.LineSpacing = t.Face(face).Size * 1.2
	text.Draw(t.screen, txt, t.Face(face), op)
	// log.Println("DrawString", txt, rect, hAlign, vAlign, fg)
}

func (t *Console) Reset() { t.booted = false }
func (t *Console) Exit()  { t.quit = true }

func (t *Console) Elli(x, y, a, b float32, clr color.Color) {
	if a <= 0 || b <= 0 {
		return
	}

	base := &vector.Path{}
	base.Arc(0, 0, 1, 0, 2*math.Pi, vector.Clockwise)
	base.Close()

	p := &vector.Path{}
	var geom ebiten.GeoM
	geom.Scale(float64(a), float64(b))
	geom.Translate(float64(x), float64(y))
	p.AddPath(base, &vector.AddPathOptions{GeoM: geom})

	var colorScale ebiten.ColorScale
	colorScale.ScaleWithColor(clr)
	vector.FillPath(t.screen, p, &vector.FillOptions{}, &vector.DrawPathOptions{
		AntiAlias:  false,
		ColorScale: colorScale,
	})
}

func (t *Console) ElliB(x, y, a, b float32, clr color.Color) {
	if a <= 0 || b <= 0 {
		return
	}

	base := &vector.Path{}
	base.Arc(0, 0, 1, 0, 2*math.Pi, vector.Clockwise)
	base.Close()

	p := &vector.Path{}
	var geom ebiten.GeoM
	geom.Scale(float64(a), float64(b))
	geom.Translate(float64(x), float64(y))
	p.AddPath(base, &vector.AddPathOptions{GeoM: geom})

	var colorScale ebiten.ColorScale
	colorScale.ScaleWithColor(clr)
	vector.StrokePath(t.screen, p, &vector.StrokeOptions{Width: 1}, &vector.DrawPathOptions{
		AntiAlias:  false,
		ColorScale: colorScale,
	})
}

func (t *Console) Trib(x1, y1, x2, y2, x3, y3 float32, clr color.Color) {
	p := &vector.Path{}
	p.MoveTo(x1, y1)
	p.LineTo(x2, y2)
	p.LineTo(x3, y3)
	p.Close()
	var colorScale ebiten.ColorScale
	colorScale.ScaleWithColor(clr)
	vector.StrokePath(t.screen, p, &vector.StrokeOptions{Width: 1}, &vector.DrawPathOptions{
		AntiAlias:  true,
		ColorScale: colorScale,
	})
}

func (t *Console) Tri(x1, y1, x2, y2, x3, y3 float32, col color.Color) {
	p := &vector.Path{}
	p.MoveTo(float32(x1), float32(y1))
	p.LineTo(float32(x2), float32(y2))
	p.LineTo(float32(x3), float32(y3))
	p.Close()

	var colorScale ebiten.ColorScale
	colorScale.ScaleWithColor(col)
	vector.FillPath(t.screen, p, &vector.FillOptions{}, &vector.DrawPathOptions{
		AntiAlias:  true,
		ColorScale: colorScale,
	})
}

func (t *Console) Clip(x, y, w, h float32) {
	t.clip = Rect(x, y, x+w, y+h)
}
func (t *Console) ResetClip() {
	t.clip = Rect(0, 0, t.Width, t.Height)
}
func (t *Console) Camera(x, y float32) {
	t.camX = x
	t.camY = y
}
func (t *Console) ResetCamera() {
	t.camX = 0
	t.camY = 0
}

func (t *Console) update() {
	now := time.Now()
	if !t.lastTick.IsZero() {
		t.delta = now.Sub(t.lastTick)
	}
	t.lastTick = now
	t.frame++
}

func (t *Console) draw(surface *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	x, y := t.clip.Min.X, t.clip.Min.Y
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(t.screen.SubImage(RectF32ToIntRect(t.clip)).(*ebiten.Image), op)
}
