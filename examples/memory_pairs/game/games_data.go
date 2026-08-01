package game

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

type GamesData struct {
	gameID     int
	data       map[Dim][]GameData
	currentDim Dim
}

func NewGamesData() *GamesData {
	d := &GamesData{}
	d.data = make(map[Dim][]GameData)
	for h := 2; h <= 50; h++ {
		for w := h; w <= 50; w++ {
			if (w*h)%2 == 0 {
				board := Dim{row: w, column: h}
				d.data[board] = []GameData{}
			}
		}
	}
	return d
}

func (g *GamesData) SaveGame(game GameData) {
	game.id = g.ID()
	g.gameID++
	g.data[g.CurrentDim()] = append(g.data[g.CurrentDim()], game)
}
func (g *GamesData) Next() Dim {
	for i, v := range g.Levels() {
		if v == g.currentDim && i+1 < len(g.data) {
			return g.Levels()[i+1]
		}
	}
	return g.currentDim
}
func (g *GamesData) NextUnlocked() (int, Dim) {
	arr := g.Levels()
	n := 0
	for i, v := range arr {
		n = i
		if len(g.data[v]) == 0 {
			return n, arr[n]
		}
	}
	return n, arr[n]
}
func (g *GamesData) Levels() []Dim {
	var (
		dims []Dim
	)
	for k := range g.data {
		dims = append(dims, k)
	}
	slices.SortFunc(dims, func(a, b Dim) int {
		return a.row*a.column - b.row*b.column
	})
	return dims
}

func (g GamesData) ID() int              { return g.gameID }
func (g GamesData) CurrentDim() Dim      { return g.currentDim }
func (g *GamesData) SetCurrentDim(v Dim) { g.currentDim = v }
func (g GamesData) String() string {
	var s strings.Builder
	for key, value := range g.data {
		for _, v := range value {
			fmt.Fprintf(&s, "Dim:%v,Games:%v\n", key, v.String())
		}
	}
	return s.String()
}

func (g GamesData) Top10() []string {
	var arr []GameData
	var s []string
	for _, v := range g.data {
		arr = append(arr, v...)
	}
	sort.SliceStable(arr, func(a, b int) bool {
		return arr[a].score > arr[b].score
	})
	if len(arr) > 10 {
		arr = arr[:10]
	}
	for _, v := range arr {
		s = append(s, v.String())
	}
	return s
}
