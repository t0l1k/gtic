package game

import (
	"etic"
	"fmt"
	"math/rand"
	"slices"
	"time"
)

type GameState int

const (
	GameStart GameState = iota
	GamePlay
	GamePause
	GameWin
)

type Game struct {
	state         etic.Property[GameState]
	data          *GameData
	field         []*Card
	first, second int
	mismatchTimer *etic.Timer
}

func NewGame() *Game {
	g := &Game{
		state:         *etic.NewProperty(GameStart),
		mismatchTimer: etic.NewTimer("timer.mismatch.show.delay").Setup(800*time.Millisecond, false, true),
	}
	g.mismatchTimer.Timeout.Connect(func(t *etic.Timer) {
		g.field[g.first].state.Set(CellClosed)
		g.field[g.second].state.Set(CellClosed)
		g.first, g.second = -1, -1
	})
	return g
}
func (g *Game) New(data *GameData) {
	g.state.Set(GameStart)
	g.data = data
	g.field = nil
	makeCards := func(count int) (cards []*Card) {
		for value := 0; value < count/2; value++ {
			cards = append(cards, NewCard(value), NewCard(value))
		}
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
		return cards
	}
	g.field = makeCards(g.data.dim.row * g.data.dim.column)
	g.first = -1
	g.second = -1
}

func (g *Game) Move(i int) {
	card := g.field[i]
	if card.state.Get() == CellOpen || card.state.Get() == CellMatch || g.first == i || !g.mismatchTimer.IsStopped() {
		return
	}
	card.state.Set(CellOpen)
	g.data.clicks++
	g.data.Changed.Emit("gamedata.clicks")
	if g.first == -1 {
		g.first = i
		return
	}

	g.second = i
	g.data.moves++
	g.data.Changed.Emit("gamedata.moves")
	if g.field[g.first].value == g.field[g.second].value {
		g.field[g.first].state.Set(CellMatch)
		g.field[g.second].state.Set(CellMatch)
		g.first, g.second = -1, -1
		if !g.StillCollecting() {
			g.state.Set(GameWin)
		}
		return
	} else {
		g.mismatchTimer.Start()
	}
}

func (g *Game) StillCollecting() bool {
	return slices.ContainsFunc(g.field, func(e *Card) bool {
		return !(e.state.Get() == CellMatch)
	})
}

func (g *Game) Finish() GameData {
	g.data.Finish()
	return *g.data
}

func (g *Game) String() (result string) {
	r, c := g.data.dim.row, g.data.dim.column
	result = fmt.Sprintf("\nРазмер поля [%vx%v]\nНажатий:%v\n", r, c, g.data.clicks)
	for y := 0; y < c; y++ {
		for x := 0; x < r; x++ {
			idx := y*r + x
			result += fmt.Sprintf("[%.2v(%v)]", g.field[idx].state.Get(), g.field[idx].value)
		}
		result += "\n"
	}
	return result
}
