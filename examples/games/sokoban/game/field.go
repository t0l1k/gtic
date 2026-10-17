package game

import "gtic/react"

type Field []*Cell

func NewField() Field { return Field{} }
func (f Field) Load(level string) Field {
	for _, rune := range level {
		if rune == '\n' {
			continue
		}
		cell := &Cell{Var: *react.NewVar(Floor)}
		switch rune {
		case ' ':
			cell.Set(Floor)
		case '#':
			cell.Set(Wall)
		case '.':
			cell.Set(Floor)
			cell.Goal = true
		case '@':
			cell.Set(Floor)
			cell.Player = true
		case '+':
			cell.Set(Floor)
			cell.Goal = true
			cell.Player = true
		case '$':
			cell.Set(Floor)
			cell.Box = true
		case '*':
			cell.Set(Floor)
			cell.Goal = true
			cell.Box = true
		}
		f = append(f, cell)
	}
	return f
}
func (f Field) Cell(idx int) *Cell { return f[idx] }
func (f Field) PlayerAt() int {
	for i, cell := range f {
		if cell.Player {
			return i
		}
	}
	return -1
}
func (f Field) IsBoxAt(idx int) bool { return f[idx].Box }
func (f Field) GoalCount() int {
	count := 0
	for _, cell := range f {
		if cell.Goal {
			count++
		}
	}
	return count
}
func (f Field) MoveBox(a, b int) {
	aValue := f.Cell(a).Box
	bValue := f.Cell(b).Box
	f[a].Box = bValue
	f[b].Box = aValue
}
func (f Field) MovePlayer(a, b int) {
	aValue := f.Cell(a).Player
	bValue := f.Cell(b).Player
	f[a].Player = bValue
	f[b].Player = aValue
}
