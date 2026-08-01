package game

import (
	"etic"
	"fmt"
)

type CellState int

const (
	CellClosed CellState = iota
	CellOpen
	CellMatch
)

type Card struct {
	state etic.Property[CellState]
	value int
}

func NewCard(v int) *Card { return &Card{state: *etic.NewProperty(CellClosed), value: v} }

func (c *Card) String() string {
	switch c.state.Get() {
	case CellClosed:
		return fmt.Sprintf("%v", "*")
	default:
		return fmt.Sprintf("%v", c.value)
	}
}
