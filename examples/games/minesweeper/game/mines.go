package game

import (
	"gtic/react"
	"image"
	"math/rand/v2"
)

type GameState int

const (
	GameStart GameState = iota
	GamePlay
	GamePause
	GameGameOver
	GameWin
)

func (g GameState) String() string {
	return []string{"Start", "Play", "Pause", "Game Over", "Winned"}[g]
}

type mcell struct{ Mine, Open, Flag, Wrong bool }

type minesGame struct {
	w, h, mines int
	cells       []mcell
	State       *react.Property[GameState]
}

func NewMines(w, h, mines int) *minesGame {
	g := &minesGame{w: w, h: h, mines: mines, State: react.NewProperty(GameStart)}
	g.Reset()
	return g
}

func (g *minesGame) Cells() []mcell       { return g.cells }
func (g *minesGame) Dim() (int, int, int) { return g.w, g.h, g.mines }
func (g *minesGame) Reset() {
	g.cells = make([]mcell, g.w*g.h)
	g.State.Set(GameStart)
}

func (g *minesGame) Count(x, y int) int {
	n := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && ny >= 0 && nx < g.w && ny < g.h && g.cells[ny*g.w+nx].Mine {
				n++
			}
		}
	}
	return n
}

// Plant расставляет мины, избегая клетки (sx, sy)
func (g *minesGame) Plant(sx, sy int) {
	g.State.Set(GamePlay)
	placed := 0
	for placed < g.mines {
		x := rand.IntN(g.w)
		y := rand.IntN(g.h)
		c := &g.cells[y*g.w+x]
		if c.Mine || (absI(x-sx) <= 1 && absI(y-sy) <= 1) {
			continue
		}
		c.Mine = true
		placed++
	}
}

func absI(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (g *minesGame) Open(x, y int) {
	state := g.State.Get()
	if state == GameWin || state == GameGameOver {
		return
	}
	c := &g.cells[y*g.w+x]
	if c.Open || c.Flag {
		return
	}
	if state == GameStart {
		g.Plant(x, y)
	}
	if c.Mine {
		c.Open = true
		g.State.Set(GameGameOver)
		// показать все мины и ошибочные флаги
		for i := range g.cells {
			if g.cells[i].Mine && !g.cells[i].Flag {
				g.cells[i].Open = true
			}
			if g.cells[i].Flag && !g.cells[i].Mine {
				g.cells[i].Wrong = true
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
		if cc.Open || cc.Flag {
			continue
		}
		cc.Open = true
		if g.Count(p.X, p.Y) == 0 {
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
		if g.cells[i].Open {
			opened++
		}
	}
	if opened == g.w*g.h-g.mines {
		g.State.Set(GameWin)
	}
}

func (g *minesGame) Flag(x, y int) {
	state := g.State.Get()
	if state == GameWin || state == GameGameOver {
		return
	}
	c := &g.cells[y*g.w+x]
	if !c.Open {
		c.Flag = !c.Flag
	}
}

func (g *minesGame) Flags() int {
	flags := 0
	for _, cell := range g.Cells() {
		if cell.Flag {
			flags++
		}
	}
	return flags
}
