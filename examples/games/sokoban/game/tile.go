package game

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type Tile int

const (
	Floor Tile = iota
	FloorGoal
	Goal       // .
	GoalPlayer // +
	GoalBox    // *
	Player     // @
	Box        // $
	Wall       // #
)

func (t Tile) String() string { return []string{" ", " ", ".", "+", "*", "@", "$", "#"}[t] }
func (t Tile) Color() color.Color {
	return []color.Color{
		colornames.Navy,
		colornames.Midnightblue,
		colornames.Gold,
		colornames.Blue,
		colornames.Goldenrod,
		colornames.Orchid,
		colornames.Maroon,
		colornames.Silver,
	}[t]
}
