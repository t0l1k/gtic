package game

import (
	"fmt"
	"gtic"
	"image"
	"log"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

var (
	Left  = image.Pt(-1, 0)
	Right = image.Pt(1, 0)
	Up    = image.Pt(0, -1)
	Down  = image.Pt(0, 1)
)

type Board struct {
	Dim         Dim
	Field       Field
	Win         bool
	level       string
	Moves       int
	Pushes      int
	history     []boardSnapshot
	historyPos  int
	playerSprId gtic.SpriteID
}

type cellSnapshot struct {
	player bool
	box    bool
}

type boardSnapshot struct {
	cells  []cellSnapshot
	moves  int
	pushes int
	win    bool
}

func NewBoard(level string) *Board {
	b := &Board{
		Field:       NewField().Load(level),
		Dim:         NewDim(level),
		level:       level,
		playerSprId: -1,
	}
	b.recordHistory()
	return b
}
func (b *Board) Reset() {
	b.Field = NewField().Load(b.level)
	b.Dim = NewDim(b.level)
	b.Win = false
	b.Moves = 0
	b.Pushes = 0
	b.history = nil
	b.historyPos = 0
	b.recordHistory()
}

func (b *Board) Input(a *gtic.API) {
	if a.KeyP(gtic.KeyU) {
		b.Undo()
		return
	}
	if a.KeyP(gtic.KeyF) {
		b.Redo()
		return
	}
	if b.Win {
		return
	}
	if a.KeyP(gtic.KeyLEFT) {
		b.Move(Left)
	}
	if a.KeyP(gtic.KeyRIGHT) {
		b.Move(Right)
	}
	if a.KeyP(gtic.KeyUP) {
		b.Move(Up)
	}
	if a.KeyP(gtic.KeyDOWN) {
		b.Move(Down)
	}
}
func (b *Board) recordHistory() {
	snapshot := boardSnapshot{
		cells:  make([]cellSnapshot, len(b.Field)),
		moves:  b.Moves,
		pushes: b.Pushes,
		win:    b.Win,
	}
	for i, cell := range b.Field {
		snapshot.cells[i] = cellSnapshot{player: cell.Player, box: cell.Box}
	}
	if b.historyPos < len(b.history)-1 {
		b.history = b.history[:b.historyPos+1]
	}
	b.history = append(b.history, snapshot)
	b.historyPos = len(b.history) - 1
}
func (b *Board) restoreHistory(pos int) {
	snapshot := b.history[pos]
	for i, cell := range b.Field {
		cell.Player = snapshot.cells[i].player
		cell.Box = snapshot.cells[i].box
	}
	b.Moves = snapshot.moves
	b.Pushes = snapshot.pushes
	b.Win = snapshot.win
	b.historyPos = pos
}
func (b *Board) Undo() bool {
	if b.historyPos == 0 {
		return false
	}
	b.restoreHistory(b.historyPos - 1)
	return true
}
func (b *Board) Redo() bool {
	if b.historyPos >= len(b.history)-1 {
		return false
	}
	b.restoreHistory(b.historyPos + 1)
	return true
}
func (b *Board) IsWin() bool {
	count := b.Field.GoalCount()
	for _, cell := range b.Field {
		if cell.IsBoxAtGoal() {
			count--
		}
	}
	return count == 0
}
func (b *Board) PushBox(idx int, direction image.Point) bool {
	boxPt := image.Pt(b.Dim.Pos(idx))
	newPt := boxPt.Add(direction)
	if newPt.X < 0 || newPt.Y < 0 || newPt.X >= b.Dim.Row || newPt.Y >= b.Dim.Column {
		return false
	}
	newCell := b.Field.Cell(b.Dim.Idx(newPt.X, newPt.Y))
	if newCell.Get() == Wall || newCell.Box {
		return false
	}
	boxIdx := b.Dim.Idx(boxPt.X, boxPt.Y)
	b.Field.MoveBox(b.Dim.Idx(newPt.X, newPt.Y), boxIdx)
	b.Field.MovePlayer(boxIdx, b.Field.PlayerAt())
	b.Moves++
	b.Pushes++
	b.Win = b.IsWin()
	b.recordHistory()
	log.Println("push", b.String())
	return true
}
func (b *Board) Move(direction image.Point) bool {
	playerIdx := b.Field.PlayerAt()
	if playerIdx < 0 {
		return false
	}
	playerPt := image.Pt(b.Dim.Pos(playerIdx))
	newPt := playerPt.Add(direction)
	if newPt.X < 0 || newPt.Y < 0 || newPt.X >= b.Dim.Row || newPt.Y >= b.Dim.Column {
		return false
	}
	newIdx := b.Dim.Idx(newPt.X, newPt.Y)
	newCell := b.Field.Cell(newIdx)
	switch newCell.Get() {
	case Wall:
		return false
	case Floor:
		if newCell.Box {
			return b.PushBox(newIdx, direction)
		}
	}
	b.Field.MovePlayer(newIdx, playerIdx)
	b.Moves++
	b.Win = b.IsWin()
	b.recordHistory()
	log.Println("move", b.String())
	return true
}

func (b *Board) MoveTo(x, y int) bool {
	path, ok := b.PathTo(x, y)
	if !ok {
		return false
	}
	for _, direction := range path {
		if !b.Move(direction) {
			return false
		}
	}
	return true
}

func (b *Board) PathTo(x, y int) ([]image.Point, bool) {
	if b.Win || x < 0 || y < 0 || x >= b.Dim.Row || y >= b.Dim.Column {
		return nil, false
	}
	start := b.Field.PlayerAt()
	target := b.Dim.Idx(x, y)
	if start < 0 || target >= len(b.Field) {
		return nil, false
	}
	if target == start {
		return nil, true
	}
	targetCell := b.Field.Cell(target)
	if targetCell.Get() == Wall || targetCell.Box {
		return nil, false
	}

	parents := make([]int, len(b.Field))
	for i := range parents {
		parents[i] = -1
	}
	steps := make([]image.Point, len(b.Field))
	parents[start] = start
	queue := []int{start}
	directions := [...]image.Point{Up, Right, Down, Left}

	for len(queue) > 0 && parents[target] == -1 {
		current := queue[0]
		queue = queue[1:]
		currentX, currentY := b.Dim.Pos(current)
		for _, direction := range directions {
			nextX, nextY := currentX+direction.X, currentY+direction.Y
			if nextX < 0 || nextY < 0 || nextX >= b.Dim.Row || nextY >= b.Dim.Column {
				continue
			}
			next := b.Dim.Idx(nextX, nextY)
			cell := b.Field.Cell(next)
			if parents[next] != -1 || cell.Get() == Wall || cell.Box {
				continue
			}
			parents[next] = current
			steps[next] = direction
			queue = append(queue, next)
		}
	}
	if parents[target] == -1 {
		return nil, false
	}

	var path []image.Point
	for current := target; current != start; current = parents[current] {
		path = append(path, steps[current])
	}
	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}
	return path, true
}

type boxPlanState struct {
	player int
	box    int
}

type boxPlanStep struct {
	previous  boxPlanState
	direction image.Point
}

func (b *Board) PathBoxTo(boxIdx, targetX, targetY int) ([]image.Point, bool) {
	if b.Win || boxIdx < 0 || boxIdx >= len(b.Field) || !b.Field[boxIdx].Box ||
		targetX < 0 || targetY < 0 || targetX >= b.Dim.Row || targetY >= b.Dim.Column {
		return nil, false
	}
	target := b.Dim.Idx(targetX, targetY)
	if target >= len(b.Field) {
		return nil, false
	}
	if target == boxIdx {
		return nil, true
	}
	targetCell := b.Field.Cell(target)
	if targetCell.Get() == Wall || targetCell.Box {
		return nil, false
	}
	player := b.Field.PlayerAt()
	if player < 0 {
		return nil, false
	}

	start := boxPlanState{player: player, box: boxIdx}
	queue := []boxPlanState{start}
	parents := map[boxPlanState]boxPlanStep{}
	visited := map[boxPlanState]bool{start: true}
	directions := [...]image.Point{Up, Right, Down, Left}
	var finish boxPlanState
	found := false

	for len(queue) > 0 && !found {
		current := queue[0]
		queue = queue[1:]
		playerX, playerY := b.Dim.Pos(current.player)
		for _, direction := range directions {
			nextX, nextY := playerX+direction.X, playerY+direction.Y
			if nextX < 0 || nextY < 0 || nextX >= b.Dim.Row || nextY >= b.Dim.Column {
				continue
			}
			nextPlayer := b.Dim.Idx(nextX, nextY)
			if nextPlayer >= len(b.Field) {
				continue
			}
			cell := b.Field.Cell(nextPlayer)
			if cell.Get() == Wall {
				continue
			}

			next := boxPlanState{player: nextPlayer, box: current.box}
			if nextPlayer == current.box {
				boxX, boxY := b.Dim.Pos(current.box)
				nextBoxX, nextBoxY := boxX+direction.X, boxY+direction.Y
				if nextBoxX < 0 || nextBoxY < 0 || nextBoxX >= b.Dim.Row || nextBoxY >= b.Dim.Column {
					continue
				}
				nextBox := b.Dim.Idx(nextBoxX, nextBoxY)
				if nextBox >= len(b.Field) {
					continue
				}
				boxCell := b.Field.Cell(nextBox)
				if boxCell.Get() == Wall || (boxCell.Box && nextBox != boxIdx) {
					continue
				}
				next.box = nextBox
			} else if cell.Box && nextPlayer != boxIdx {
				continue
			}
			if visited[next] {
				continue
			}
			visited[next] = true
			parents[next] = boxPlanStep{previous: current, direction: direction}
			if next.box == target {
				finish = next
				found = true
				break
			}
			queue = append(queue, next)
		}
	}
	if !found {
		return nil, false
	}

	var path []image.Point
	for current := finish; current != start; {
		step := parents[current]
		path = append(path, step.direction)
		current = step.previous
	}
	for left, right := 0, len(path)-1; left < right; left, right = left+1, right-1 {
		path[left], path[right] = path[right], path[left]
	}
	return path, true
}

func (b *Board) Draw(a *gtic.API) {
	for i, cell := range b.Field {
		x, y := b.Dim.Pos(i)
		switch {
		case cell.Get() == Wall:
			a.Rect(x*8, y*8, 8, 8, gtic.NewColor(Tile(Wall).Color()))
		case cell.Get() == Floor:
			a.Rect(x*8, y*8, 8, 8, gtic.NewColor(Tile(Floor).Color()))
			switch {
			case cell.Box && cell.Goal:
				drawBox(a, x*8, y*8, true)
			case cell.Player:
				if b.playerSprId == -1 {
					b.playerSprId = LoadSprites(a)
				}
				drawPlayer(a, x*8, y*8, b.playerSprId)
			case cell.Goal:
				drawGoal(a, x*8, y*8)
			case cell.Box:
				drawBox(a, x*8, y*8, false)
			}
		}
	}
}
func (b *Board) String() string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Dim[%v,%v]\n", b.Dim.Row, b.Dim.Column))
	pos := image.Point{}
	for y := 0; y < b.Dim.Column; y++ {
		for x := 0; x < b.Dim.Row; x++ {
			result.WriteString(b.Field.Cell(b.Dim.Idx(x, y)).String())
		}
		result.WriteString("\n")
	}
	result.WriteString(fmt.Sprintf("Player At:%v\n", pos))
	return result.String()
}

func drawPlayer(a *gtic.API, x, y int, id gtic.SpriteID) { a.Spr(id, x, y, a.Pal(0)) }

func drawGoal(a *gtic.API, x, y int) {
	gold := gtic.NewColor(colornames.Gold)
	shadow := gtic.NewColor(colornames.Maroon)

	// Ромб-мишень на полу.
	a.Line(x+4, y+1, x+6, y+4, shadow)
	a.Line(x+6, y+4, x+4, y+6, shadow)
	a.Line(x+4, y+6, x+2, y+4, shadow)
	a.Line(x+2, y+4, x+4, y+1, shadow)

	a.Circ(x+4, y+4, 1, gold)
}

func drawBox(a *gtic.API, x, y int, goal bool) {
	shadow := gtic.NewColor(colornames.Black)
	wood := gtic.NewColor(colornames.Maroon)
	highlight := gtic.NewColor(colornames.Burlywood)
	if goal {
		highlight = gtic.NewColor(colornames.Khaki)
	}

	// Тень, корпус, светлый контур и диагональные доски.
	a.Elli(x+4, y+7, 3, 1, shadow)
	a.Rect(x+1, y+1, 6, 6, wood)
	a.RectB(x+1, y+1, 6, 6, highlight)
	a.Line(x+2, y+2, x+5, y+5, highlight)
	a.Line(x+5, y+2, x+2, y+5, highlight)
}

func LoadSprites(a *gtic.API) (id gtic.SpriteID) {
	PlayerSprite := "00777700" +
		"07777770" +
		"07444470" +
		"0044A400" +
		"00099000" +
		"00999900" +
		"04044040" +
		"07700770"

	registerGameSprites := func(a *gtic.API) {
		sprData := func(data string) gtic.Surface[gtic.RGBA] {
			pixels := make([]gtic.RGBA, 0)
			for i := 0; i < 64; i++ {
				val, _ := strconv.ParseUint(string(data[i]), 16, 4)
				colIdx := int(val)
				pixels = append(pixels, a.Pal(colIdx))
			}
			return gtic.Surface[gtic.RGBA]{Data: pixels, Width: 8, Height: 8}
		}
		id = a.Sprites().Register(gtic.NewSprite("player.sprite", sprData(PlayerSprite), image.Rect(0, 0, 8, 8)))
	}
	registerGameSprites(a)
	return id
}
