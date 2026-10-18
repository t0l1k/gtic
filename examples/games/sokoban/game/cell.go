package game

import (
	"fmt"
	"gtic/react"
)

type Cell struct{ *react.Var[Tile] }

func NewCell() *Cell                 { return &Cell{Var: react.NewVar(Wall)} }
func (c *Cell) IsFloor() bool        { return c.Var.Get() == Floor }
func (c *Cell) SetFloor() *Cell      { c.Var.Set(Floor); return c }
func (c *Cell) IsGoal() bool         { return c.Var.Get() == Goal }
func (c *Cell) SetGoal() *Cell       { c.Var.Set(Goal); return c }
func (c *Cell) IsGoalPlayer() bool   { return c.Var.Get() == GoalPlayer }
func (c *Cell) SetGoalPlayer() *Cell { c.Var.Set(GoalPlayer); return c }
func (c *Cell) IsGoalBox() bool      { return c.Var.Get() == GoalBox }
func (c *Cell) SetGoalBox() *Cell    { c.Var.Set(GoalBox); return c }
func (c *Cell) IsPlayer() bool       { return c.Var.Get() == Player }
func (c *Cell) SetPlayer() *Cell     { c.Var.Set(Player); return c }
func (c *Cell) IsBox() bool          { return c.Var.Get() == Box }
func (c *Cell) SetBox() *Cell        { c.Var.Set(Box); return c }
func (c *Cell) IsWall() bool         { return c.Var.Get() == Wall }
func (c *Cell) SetWall() *Cell       { c.Var.Set(Wall); return c }
func (c *Cell) String() string {
	switch {
	case c.IsGoalPlayer():
		return Tile(GoalPlayer).String()
	case c.IsGoalBox():
		return Tile(GoalBox).String()
	case c.IsPlayer():
		return Tile(Player).String()
	case c.IsBox():
		return Tile(Box).String()
	case c.IsGoal():
		return Tile(Goal).String()
	case c.IsWall():
		return Tile(Wall).String()
	case c.IsFloor():
		return Tile(Floor).String()
	default:
		res := fmt.Sprintf("unknown cell:%v", c.Get().String())
		panic(res)
	}
}
