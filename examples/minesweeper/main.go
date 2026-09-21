package main

import (
	_ "embed"
	"gtic"
	"gtic/ui"
	"image"
	"log"
	"math/rand/v2"
	"strconv"
)

const (
	SprDigit1 = iota // 0..9 — большие цифры для счётчика
	SprDigit2
	SprDigit3
	SprDigit4
	SprDigit5
	SprDigit6
	SprDigit7
	SprDigit8
	SprDigit9
	SprDigit0
	SprFaceIdle   // нормальное лицо
	SprFacePlay   // во время игры
	SprFaceOuch   // "ой"
	SprFaceWin    // победа
	SprFaceDead   // проигрыш
	SprCellClosed // закрытая клетка
	SprCellOpen   // открытая (пустая)
	SprCellFlag   // флаг
	SprCellQ      // вопрос
	SprCellQOpen  // вопрос (нажата)
	SprCellMine   // мина
	SprCellBoom   // взорвавшаяся мина
	SprCellWrong  // флаг ошибочный
	SprNum1       // 1..8 — цифры на поле (SprNum1 + n - 1)
	SprNum2
	SprNum3
	SprNum4
	SprNum5
	SprNum6
	SprNum7
	SprNum8
)

var (
	//go:embed mines.png
	MinesSprites_png []byte
)

func RegisterSprites(a *gtic.API) []gtic.SpriteID {
	// rectsMineSheet — разметка листа (139×84) mines.png
	surface, err := gtic.DecodeImage(MinesSprites_png)
	if err != nil {
		panic(err)
	}
	var ids []gtic.SpriteID
	for i := 0; i < 10; i++ { // ряд 0: большие цифры 13×23, шаг 14
		r := image.Rect(i*14, 0, i*14+13, 23)
		name := "Digit:" + strconv.Itoa(i)
		id := a.Sprites().Register(gtic.NewSprite(name, surface, r))
		ids = append(ids, id)
	}
	for i := 0; i < 5; i++ { // ряд 1: лица 26×26, шаг 27
		r := image.Rect(i*27, 24, i*27+26, 50)
		name := "Face:" + strconv.Itoa(i)
		id := a.Sprites().Register(gtic.NewSprite(name, surface, r))
		ids = append(ids, id)
	}
	for i := 0; i < 8; i++ { // ряд 2: клетки поля 16×16, шаг 17
		r := image.Rect(i*17, 51, i*17+16, 67)
		name := "UpSquare:" + strconv.Itoa(i)
		id := a.Sprites().Register(gtic.NewSprite(name, surface, r))
		ids = append(ids, id)
	}
	for i := 0; i < 8; i++ { // ряд 3: цифры 1–8, 16×16
		r := image.Rect(i*17, 68, i*17+16, 84)
		name := "DownSquare:" + strconv.Itoa(i)
		id := a.Sprites().Register(gtic.NewSprite(name, surface, r))
		ids = append(ids, id)
	}
	return ids
}

type mcell struct {
	mine, open, flag, wrong bool
}

type minesGame struct {
	w, h, mines        int
	cells              []mcell
	cur                image.Point // курсор в клетках
	dead, win, started bool
}

func newMines(w, h, mines int) *minesGame {
	g := &minesGame{w: w, h: h, mines: mines}
	g.reset()
	return g
}

func (g *minesGame) reset() {
	g.cells = make([]mcell, g.w*g.h)
	g.cur = image.Pt(g.w/2, g.h/2)
	g.dead, g.win, g.started = false, false, false
}

func (g *minesGame) count(x, y int) int {
	n := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && ny >= 0 && nx < g.w && ny < g.h && g.cells[ny*g.w+nx].mine {
				n++
			}
		}
	}
	return n
}

// plant расставляет мины, избегая клетки (sx, sy)
func (g *minesGame) plant(sx, sy int) {
	g.started = true
	placed := 0
	for placed < g.mines {
		x := rand.IntN(g.w)
		y := rand.IntN(g.h)
		c := &g.cells[y*g.w+x]
		if c.mine || (absI(x-sx) <= 1 && absI(y-sy) <= 1) {
			continue
		}
		c.mine = true
		placed++
	}
}

func absI(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (g *minesGame) open(x, y int) {
	if g.dead || g.win {
		return
	}
	c := &g.cells[y*g.w+x]
	if c.open || c.flag {
		return
	}
	if !g.started {
		g.plant(x, y)
	}
	if c.mine {
		c.open = true
		g.dead = true
		// показать все мины и ошибочные флаги
		for i := range g.cells {
			if g.cells[i].mine && !g.cells[i].flag {
				g.cells[i].open = true
			}
			if g.cells[i].flag && !g.cells[i].mine {
				g.cells[i].wrong = true
			}
		}
		return
	}
	// заливка пустых
	stack := []image.Point{{x, y}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		cc := &g.cells[p.Y*g.w+p.X]
		if cc.open || cc.flag {
			continue
		}
		cc.open = true
		if g.count(p.X, p.Y) == 0 {
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					nx, ny := p.X+dx, p.Y+dy
					if nx >= 0 && ny >= 0 && nx < g.w && ny < g.h {
						stack = append(stack, image.Pt(nx, ny))
					}
				}
			}
		}
	}
	// проверка победы
	opened := 0
	for i := range g.cells {
		if g.cells[i].open {
			opened++
		}
	}
	if opened == g.w*g.h-g.mines {
		g.win = true
	}
}

func (g *minesGame) flag(x, y int) {
	if g.dead || g.win {
		return
	}
	c := &g.cells[y*g.w+x]
	if !c.open {
		c.flag = !c.flag
	}
}

func NewMinesGame() *ui.Element {
	var (
		sprIds                    []gtic.SpriteID
		ox, oy                    int
		leftPressed, rightPressed bool
	)

	g := newMines(9, 9, 10)

	game := ui.NewElement("minesweeper.game")
	game.OnInit = func(a *gtic.API) {
		sprIds = RegisterSprites(a)
		ox, oy = (a.Bounds().Width-9*16)/2, 34 // поле по центру, под панелью
	}
	const cell = 16

	game.OnUpdate = func(a *gtic.API) {
		mx, my, lp, _, rp, _, _ := a.Mouse()
		gx, gy := (mx-ox)/cell, (my-oy)/cell
		inside := gx >= 0 && gy >= 0 && gx < g.w && gy < g.h
		if inside {
			g.cur = image.Pt(gx, gy)
		}

		// нажатие по кромке (edge), а не удержание
		if lp && !leftPressed && inside {
			g.open(gx, gy)
		}
		if rp && !rightPressed && inside {
			g.flag(gx, gy)
		}
		leftPressed, rightPressed = lp, rp

		switch {
		case a.KeyP(gtic.KeyR):
			g.reset()
		case a.KeyP(gtic.KeyLEFT):
			g.cur.X = max(g.cur.X-1, 0)
		case a.KeyP(gtic.KeyRIGHT):
			g.cur.X = min(g.cur.X+1, g.w-1)
		case a.KeyP(gtic.KeyUP):
			g.cur.Y = max(g.cur.Y-1, 0)
		case a.KeyP(gtic.KeyDOWN):
			g.cur.Y = min(g.cur.Y+1, g.h-1)
		case a.KeyP(gtic.KeyZ):
			g.open(g.cur.X, g.cur.Y)
		case a.KeyP(gtic.KeyX):
			g.flag(g.cur.X, g.cur.Y)
		}
	}
	game.OnDraw = func(a *gtic.API) {
		a.Cls()

		// --- верхняя панель: счётчик мин + лицо ---
		flags := 0
		for i := range g.cells {
			if g.cells[i].flag {
				flags++
			}
		}
		left := max(g.mines-flags, 0)
		aa := left / 100 % 10
		bb := left / 10 % 10
		cc := left % 10

		res := func(v int) int {
			switch v {
			case 0:
				return 9
			default:
				return v - 1
			}
		}
		a.Spr(sprIds[res(aa)], 8, 5)
		a.Spr(sprIds[res(bb)], 22, 5)
		a.Spr(sprIds[res(cc)], 36, 5)

		face := SprFaceIdle
		switch {
		case leftPressed || rightPressed:
			face = SprFacePlay
		case g.dead:
			face = SprFaceDead
		case g.win:
			face = SprFaceWin
		}
		fx := (a.Bounds().Width - 26) / 2
		a.Spr(sprIds[face], fx, 0)

		// --- поле ---
		for y := 0; y < g.h; y++ {
			for x := 0; x < g.w; x++ {
				c := g.cells[y*g.w+x]
				px, py := ox+x*16, oy+y*16
				switch {
				case c.wrong:
					a.Spr(sprIds[SprCellWrong], px, py)
				case c.open && c.mine:
					if g.dead && x == g.cur.X && y == g.cur.Y {
						a.Spr(sprIds[SprCellBoom], px, py)
					} else {
						a.Spr(sprIds[SprCellMine], px, py)
					}
				case c.open:
					a.Spr(sprIds[SprCellOpen], px, py)
					if n := g.count(x, y); n > 0 {
						a.Spr(sprIds[SprNum1+n-1], px, py)
					}
				case c.flag:
					a.Spr(sprIds[SprCellFlag], px, py)
				default:
					a.Spr(sprIds[SprCellClosed], px, py)
				}
			}
		}

		// рамка курсора
		if !g.dead && !g.win {
			cx, cy := ox+g.cur.X*16-1, oy+g.cur.Y*16-1
			for x := 0; x < 18; x++ {
				a.Pix(cx+x, cy, 0) // подставьте ваш set пикселя
				a.Pix(cx+x, cy+17, 0)
				a.Pix(cx, cy+x, 0)
				a.Pix(cx+17, cy+x, 0)
			}
		}
	}

	return game
}

func main() {
	game := NewMinesGame()
	gtic.BOOT = func(a *gtic.API) {
		game.Init(a)
		log.Println("Booted minesweeper game")
	}
	gtic.TIC = func(a *gtic.API) {
		if a.Key(gtic.KeyESC) {
			a.Exit()
		}
		if a.KeyP(gtic.KeyRETURN) {
			a.Reset()
		}
		game.Update(a)
		game.Draw(a)
	}
	if err := gtic.Load(gtic.WithTitle("Minesweeper"), gtic.WithMode(320, 200)).Run(); err != nil {
		log.Println(err)
	}
}
