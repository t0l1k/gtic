package game

type Field []*Cell

func NewField() Field { return Field{} }
func (f Field) Load(level string) Field {
	for _, rune := range level {
		if rune == '\n' {
			continue
		}
		cell := NewCell()
		switch rune {
		case Floor.Rune():
			cell.SetFloor()
		case Wall.Rune():
			cell.SetWall()
		case Goal.Rune():
			cell.SetGoal()
		case Player.Rune():
			cell.SetPlayer()
		case GoalPlayer.Rune():
			cell.SetGoalPlayer()
		case Box.Rune():
			cell.SetBox()
		case GoalBox.Rune():
			cell.SetGoalBox()
		}
		f = append(f, cell)
	}
	return f
}
func (f Field) Cell(idx int) *Cell { return f[idx] }
func (f Field) PlayerAt() int {
	for i, cell := range f {
		if cell.IsPlayer() || cell.IsGoalPlayer() {
			return i
		}
	}
	return -1
}
func (f Field) IsBoxAt(idx int) bool { return f[idx].IsBox() || f[idx].IsGoalBox() }
func (f Field) GoalCount() int {
	count := 0
	for _, cell := range f {
		if cell.IsGoal() || cell.IsGoalPlayer() || cell.IsGoalBox() {
			count++
		}
	}
	return count
}
func (f Field) MoveBox(a, b int) {
	aValue := f.Cell(a).Get() // floor, goal
	bValue := f.Cell(b).Get() // box, goalBox
	switch {
	case aValue == Wall|Player|GoalPlayer:
		panic("push failed")
	case aValue == Floor && bValue == Box:
		aValue = Floor
		bValue = Box
	case aValue == Floor && bValue == GoalBox:
		aValue = Goal
		bValue = Box
	case aValue == Goal && bValue == Box:
		aValue = Floor
		bValue = GoalBox
	case aValue == Goal && bValue == GoalBox:
		aValue = Goal
		bValue = GoalBox
	}
	f[a].Set(bValue)
	f[b].Set(aValue)
}
func (f Field) MovePlayer(a, b int) {
	aValue := f.Cell(a).Get() // floor, goal
	bValue := f.Cell(b).Get() // player, goalPlayer
	switch {
	case aValue == Wall|Box|GoalBox:
		panic("move failed")
	case aValue == Floor && bValue == Player:
		bValue = Player
	case aValue == Floor && bValue == GoalPlayer:
		aValue = Goal
		bValue = Player
	case aValue == Goal && bValue == Player:
		bValue = GoalPlayer
		aValue = Floor
	case aValue == Goal && bValue == GoalPlayer:
		aValue = Goal
		bValue = GoalPlayer
	}
	f[a].Set(bValue)
	f[b].Set(aValue)
}
