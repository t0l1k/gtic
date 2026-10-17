package game

import (
	"gtic/react"
)

type Cell struct {
	react.Var[Tile]
	Goal, Player, Box bool
}

func (c *Cell) IsBoxAtGoal() bool { return c.Box == true && c.Goal == true }
func (c *Cell) String() string {
	switch {
	case c.Player && c.Goal:
		return Tile(GoalPlayer).String()
	case c.Box && c.Goal:
		return Tile(GoalBox).String()
	case c.Player:
		return Tile(Player).String()
	case c.Box:
		return Tile(Box).String()
	case c.Goal:
		return Tile(Goal).String()
	case c.Get() == Wall:
		return Tile(Wall).String()
	default:
		return Tile(Floor).String()
	}
}
