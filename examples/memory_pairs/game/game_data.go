package game

import (
	"etic"
	"fmt"
	"time"
)

type GameData struct {
	id       int
	dim      Dim
	start    time.Time
	duration time.Duration
	clicks   int
	moves    int
	score    int
	unlocked bool
	Changed  *etic.Signal[string]
}

func NewGameData(dim Dim, id int) *GameData {
	return &GameData{dim: dim, id: id, Changed: etic.NewSignal[string](), unlocked: false}
}
func (d *GameData) Start() { d.start = time.Now() }
func (d *GameData) Finish() *GameData {
	baseScore := d.dim.row * d.dim.column * 150
	score := max(100, baseScore-d.moves*25)
	d.duration = time.Since(d.start)
	d.score = score
	d.unlocked = true
	return d
}
func (d *GameData) String() string {
	return fmt.Sprintf("dim:%v, game:id:%v, clicks:%v, moves:%v, score:%v, duration:%v", d.dim, d.id, d.clicks, d.moves, d.score, d.duration.String())
}
