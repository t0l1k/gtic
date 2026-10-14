package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"gtic"
	"gtic/react"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed simple.xsb xsokoban.xsb
var levelFiles embed.FS

type levelPack struct {
	name   string
	file   string
	levels []string
}

func loadPacks() ([]levelPack, error) {
	files := []string{"simple.xsb", "xsokoban.xsb"}
	packs := make([]levelPack, 0, len(files))
	for _, file := range files {
		data, err := levelFiles.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read level pack %q: %w", file, err)
		}
		levels, err := parseLevels(string(data))
		if err != nil {
			return nil, fmt.Errorf("parse level pack %q: %w", file, err)
		}
		packs = append(packs, levelPack{
			name:   strings.TrimSuffix(file, filepath.Ext(file)),
			file:   file,
			levels: levels,
		})
	}
	return packs, nil
}

func parseLevels(data string) ([]string, error) {
	var (
		levels []string
		rows   []string
	)
	flush := func() error {
		if len(rows) == 0 {
			return nil
		}
		width := 0
		playerCount := 0
		for _, row := range rows {
			width = max(width, len(row))
			for _, tile := range row {
				if tile == '@' || tile == '+' {
					playerCount++
				}
			}
		}
		if playerCount != 1 {
			return fmt.Errorf("level %d has %d players; want exactly one", len(levels)+1, playerCount)
		}
		var level strings.Builder
		for i, row := range rows {
			if i > 0 {
				level.WriteByte('\n')
			}
			level.WriteString(row)
			level.WriteString(strings.Repeat(" ", width-len(row)))
		}
		levels = append(levels, level.String())
		rows = nil
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(data))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\r")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ";") || trimmed == "" {
			if err := flush(); err != nil {
				return nil, err
			}
			continue
		}
		for _, tile := range line {
			switch tile {
			case ' ', '#', '.', '@', '+', '$', '*':
			default:
				return nil, fmt.Errorf("invalid tile %q on line %d", tile, lineNumber)
			}
		}
		rows = append(rows, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := flush(); err != nil {
		return nil, err
	}
	if len(levels) == 0 {
		return nil, fmt.Errorf("no levels found")
	}
	return levels, nil
}

type savedProgress map[string]int

func loadProgress() (savedProgress, string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, "", fmt.Errorf("find user config directory: %w", err)
	}
	path := filepath.Join(configDir, "gtic", "sokoban-progress.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return savedProgress{}, path, nil
	}
	if err != nil {
		return nil, "", fmt.Errorf("read Sokoban progress: %w", err)
	}
	var progress savedProgress
	if err := json.Unmarshal(data, &progress); err != nil {
		return nil, "", fmt.Errorf("decode Sokoban progress: %w", err)
	}
	if progress == nil {
		progress = savedProgress{}
	}
	log.Println("load:", path, progress)
	return progress, path, nil
}

func saveProgress(path string, progress savedProgress) error {
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return fmt.Errorf("encode Sokoban progress: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create progress directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write Sokoban progress: %w", err)
	}
	log.Println("save:", path, progress)
	return nil
}

var (
	Left  = image.Pt(-1, 0)
	Right = image.Pt(1, 0)
	Up    = image.Pt(0, -1)
	Down  = image.Pt(0, 1)
)

type Dim struct{ row, column int }

func NewDim(level string) Dim {
	rows := strings.Split(strings.TrimSuffix(level, "\n"), "\n")
	d := Dim{column: len(rows)}
	for _, row := range rows {
		if len(row) > d.row {
			d.row = len(row)
		}
	}
	return d
}
func (d Dim) Idx(x, y int) int       { return y*d.row + x }
func (d Dim) Pos(idx int) (int, int) { return idx % d.row, idx / d.row }

type Tile int

const (
	Floor  Tile = iota //' ', '#', '.', '@', '+', '$', '*'
	Wall               // #
	Goal               // .
	Player             // @
	Box                // $
)

func (t Tile) String() string { return []string{" ", "#", ".", "@", "$"}[t] }
func (t Tile) Color() int     { return []int{12, 7, 6, 11, 3}[t] }

type Cell struct {
	react.Var[Tile]
	goal, player, box bool
}

func (c *Cell) IsBoxAtGoal() bool { return c.box == true && c.goal == true }
func (c *Cell) Color() int {
	result := -1
	switch {
	case c.player:
		result += Tile(Player).Color()
	case c.box:
		result += Tile(Box).Color()
	case c.goal:
		result += Tile(Goal).Color()
	case c.Get() == Wall:
		result += Tile(Wall).Color()
	default:
		result += Tile(Floor).Color()
	}
	return result
}
func (c *Cell) String() string {
	switch {
	case c.player && c.goal:
		return "+"
	case c.box && c.goal:
		return "*"
	case c.player:
		return Tile(Player).String()
	case c.box:
		return Tile(Box).String()
	case c.goal:
		return Tile(Goal).String()
	case c.Get() == Wall:
		return Tile(Wall).String()
	default:
		return Tile(Floor).String()
	}
}

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
			cell.goal = true
		case '@':
			cell.Set(Floor)
			cell.player = true
		case '+':
			cell.Set(Floor)
			cell.goal = true
			cell.player = true
		case '$':
			cell.Set(Floor)
			cell.box = true
		case '*':
			cell.Set(Floor)
			cell.goal = true
			cell.box = true
		}
		f = append(f, cell)
	}
	return f
}
func (f Field) Cell(idx int) *Cell { return f[idx] }
func (f Field) PlayerAt() int {
	for i, cell := range f {
		if cell.player {
			return i
		}
	}
	return -1
}
func (f Field) IsBoxAt(idx int) bool { return f[idx].box }
func (f Field) GoalCount() int {
	count := 0
	for _, cell := range f {
		if cell.goal {
			count++
		}
	}
	return count
}
func (f Field) MoveBox(a, b int) {
	aValue := f.Cell(a).box
	bValue := f.Cell(b).box
	f[a].box = bValue
	f[b].box = aValue
}
func (f Field) MovePlayer(a, b int) {
	aValue := f.Cell(a).player
	bValue := f.Cell(b).player
	f[a].player = bValue
	f[b].player = aValue
}

type Board struct {
	Dim     Dim
	Field   Field
	win     bool
	level   string
	moves   int
	pushes  int
	history []boardState
}

type boardState struct {
	player int
	boxes  []int
	moves  int
	pushes int
	win    bool
}

func NewBoard(level string) *Board {
	b := &Board{
		Field: NewField().Load(level),
		Dim:   NewDim(level),
		level: level,
	}
	return b
}
func (b *Board) Reset() {
	b.Field = NewField().Load(b.level)
	b.Dim = NewDim(b.level)
	b.win = false
	b.moves = 0
	b.pushes = 0
	b.history = nil
}
func (b *Board) saveState() {
	state := boardState{
		player: b.Field.PlayerAt(),
		moves:  b.moves,
		pushes: b.pushes,
		win:    b.win,
	}
	for i, cell := range b.Field {
		if cell.box {
			state.boxes = append(state.boxes, i)
		}
	}
	b.history = append(b.history, state)
}
func (b *Board) Undo() bool {
	if len(b.history) == 0 {
		return false
	}
	state := b.history[len(b.history)-1]
	b.history = b.history[:len(b.history)-1]
	for _, cell := range b.Field {
		cell.player = false
		cell.box = false
	}
	if state.player >= 0 {
		b.Field[state.player].player = true
	}
	for _, idx := range state.boxes {
		b.Field[idx].box = true
	}
	b.moves = state.moves
	b.pushes = state.pushes
	b.win = state.win
	return true
}
func (b *Board) Input(a *gtic.API) {
	if a.KeyP(gtic.KeyR) {
		b.Reset()
		return
	}
	if a.KeyP(gtic.KeyU) {
		b.Undo()
		return
	}
	if b.win {
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
	if newPt.X < 0 || newPt.Y < 0 || newPt.X >= b.Dim.row || newPt.Y >= b.Dim.column {
		return false
	}
	newCell := b.Field.Cell(b.Dim.Idx(newPt.X, newPt.Y))
	if newCell.Get() == Wall || newCell.box {
		return false
	}
	b.saveState()
	boxIdx := b.Dim.Idx(boxPt.X, boxPt.Y)
	b.Field.MoveBox(b.Dim.Idx(newPt.X, newPt.Y), boxIdx)
	b.Field.MovePlayer(boxIdx, b.Field.PlayerAt())
	b.moves++
	b.pushes++
	b.win = b.IsWin()
	return true
}
func (b *Board) Move(direction image.Point) {
	playerIdx := b.Field.PlayerAt()
	if playerIdx < 0 {
		return
	}
	playerPt := image.Pt(b.Dim.Pos(playerIdx))
	newPt := playerPt.Add(direction)
	if newPt.X < 0 || newPt.Y < 0 || newPt.X >= b.Dim.row || newPt.Y >= b.Dim.column {
		return
	}
	var (
		newCell *Cell
	)
	newCell = b.Field.Cell(b.Dim.Idx(newPt.X, newPt.Y))
	switch newCell.Get() {
	case Wall:
		return
	case Floor:
		if newCell.box {
			b.PushBox(b.Dim.Idx(newPt.X, newPt.Y), direction)
			return
		}
		b.saveState()
		b.Field.MovePlayer(b.Dim.Idx(newPt.X, newPt.Y), playerIdx)
		b.moves++
		b.win = b.IsWin()
	}
}
func (b *Board) Draw(a *gtic.API) {
	grid := 8

	for i, cell := range b.Field {
		x, y := b.Dim.Pos(i)
		switch {
		case cell.Get() == Wall:
			a.Rect(x*grid, y*grid, grid, grid, a.Pal(cell.Color()))
		case cell.Get() == Floor:
			a.Rect(x*grid, y*grid, grid, grid, a.Pal(Tile(Floor).Color()))
			switch {
			case cell.box:
				a.Circ(x*grid+4, y*grid+4, 3, a.Pal(cell.Color()))
			case cell.player:
				a.Circ(x*grid+4, y*grid+4, 3, a.Pal(cell.Color()))
			case cell.goal:
				a.Rect(x*grid, y*grid, grid, grid, a.Pal(Tile(Floor).Color()))
				a.Circ(x*grid+4, y*grid+4, 3, a.Pal(cell.Color()))
			}
		}
	}
	for i := 0; i < a.Bounds().Height; i++ {
		a.Line(0, i*grid, a.Bounds().Width, i*grid, a.Pal(1))
		a.Line(i*grid, 0, i*grid, a.Bounds().Height, a.Pal(1))
	}
}
func (b *Board) String() string {
	result := fmt.Sprintf("Dim[%v,%v]\n", b.Dim.row, b.Dim.column)
	pos := image.Point{}
	for y := 0; y < b.Dim.column; y++ {
		for x := 0; x < b.Dim.row; x++ {
			cell := b.Field.Cell(b.Dim.Idx(x, y))
			switch cell.Get() {
			case Floor:
				switch {
				case cell.player:
					result += cell.String()
					pos = image.Pt(x, y)
				case cell.goal:
					result += cell.String()
				case cell.box:
					result += cell.String()
				default:
					result += cell.String()
				}
			case Wall:
				result += cell.String()
			}
		}
		result += "\n"
	}
	result += fmt.Sprintf("Player At:%v\n", pos)
	return result
}

type game struct {
	packs        []levelPack
	progress     savedProgress
	progressPath string
	selectedPack int
	currentLevel int
	board        *Board
	playing      bool
	status       string
}

func newGame() (*game, error) {
	packs, err := loadPacks()
	if err != nil {
		return nil, err
	}
	progress, progressPath, err := loadProgress()
	if err != nil {
		return nil, err
	}
	for _, pack := range packs {
		progress[pack.file] = min(max(progress[pack.file], 0), len(pack.levels))
	}
	return &game{
		packs:        packs,
		progress:     progress,
		progressPath: progressPath,
	}, nil
}

func (g *game) Input(a *gtic.API) {
	if !g.playing {
		switch {
		case a.KeyP(gtic.KeyLEFT):
			g.selectedPack = (g.selectedPack + len(g.packs) - 1) % len(g.packs)
			g.status = ""
		case a.KeyP(gtic.KeyRIGHT):
			g.selectedPack = (g.selectedPack + 1) % len(g.packs)
			g.status = ""
		case a.KeyP(gtic.KeyRETURN):
			g.startPack()
		}
		if a.KeyP(gtic.KeyESC) {
			a.Exit()
		}
		return
	}

	if a.KeyP(gtic.KeyESC) && g.playing {
		g.board = nil
		g.playing = false
		return
	}
	if a.KeyP(gtic.KeyR) {
		g.board.Reset()
		return
	}
	if g.board.win {
		if a.KeyP(gtic.KeyU) {
			g.board.Undo()
			return
		}
		if a.KeyP(gtic.KeyRETURN) {
			g.nextLevel()
		}
		return
	}

	g.board.Input(a)
	if g.board.IsWin() {
		g.board.win = true
		g.saveCompletedLevel()
	}
}

func (g *game) startPack() {
	pack := g.packs[g.selectedPack]
	g.currentLevel = g.progress[pack.file]
	if g.currentLevel >= len(pack.levels) {
		g.currentLevel = 0
	}
	g.board = NewBoard(pack.levels[g.currentLevel])
	g.playing = true
	g.status = ""
}

func (g *game) saveCompletedLevel() {
	pack := g.packs[g.selectedPack]
	nextLevel := g.currentLevel + 1
	if nextLevel <= g.progress[pack.file] {
		return
	}
	g.progress[pack.file] = nextLevel
	if err := saveProgress(g.progressPath, g.progress); err != nil {
		g.status = "Progress save failed; see log."
		log.Printf("save Sokoban progress: %v", err)
		return
	}
	g.status = ""
}

func (g *game) nextLevel() {
	pack := g.packs[g.selectedPack]
	if g.currentLevel+1 == len(pack.levels) {
		g.board = nil
		g.playing = false
		g.status = "Pack complete. Enter to replay from the start."
		return
	}
	g.currentLevel++
	g.board = NewBoard(pack.levels[g.currentLevel])
	g.status = ""
}

func (g *game) Draw(a *gtic.API) {
	a.Cls(a.Pal(14))
	if !g.playing {
		g.drawPackMenu(a)
		return
	}

	g.board.Draw(a)
	pack := g.packs[g.selectedPack]
	a.Print(fmt.Sprintf("%s %d/%d  Moves: %d  Pushes: %d", pack.name, g.currentLevel+1, len(pack.levels), g.board.moves, g.board.pushes), 0, a.Bounds().Height-16, a.Pal(12))
	if g.board.win {
		a.Print("Level clear! U: undo  Enter: next", 0, a.Bounds().Height-8, a.Pal(11))
	} else {
		a.Print("Arrows: move  U: undo  R: restart", 0, a.Bounds().Height-8, a.Pal(12))
	}
	if g.status != "" {
		a.Print(g.status, 0, a.Bounds().Height-24, a.Pal(6))
	}
}

func (g *game) drawPackMenu(a *gtic.API) {
	a.Print("SOKOBAN", 8, 8, a.Pal(12))
	a.Print("< "+g.packs[g.selectedPack].name+" >", 8, 32, a.Pal(11))
	pack := g.packs[g.selectedPack]
	completed := g.progress[pack.file]
	if completed == len(pack.levels) {
		a.Print(fmt.Sprintf("All %d levels complete", completed), 8, 48, a.Pal(12))
	} else {
		a.Print(fmt.Sprintf("Progress: %d/%d", completed, len(pack.levels)), 8, 48, a.Pal(12))
	}
	a.Print("Left/Right: pack", 8, 72, a.Pal(12))
	if completed == len(pack.levels) {
		a.Print("Enter: replay from start", 8, 88, a.Pal(12))
	} else {
		a.Print("Enter: resume", 8, 88, a.Pal(12))
	}
	a.Print("Esc: quit", 8, 104, a.Pal(12))
	if g.status != "" {
		a.Print(g.status, 8, 120, a.Pal(6))
	}
}

func main() {
	game, err := newGame()
	if err != nil {
		log.Fatal(err)
	}
	gtic.BOOT = func(a *gtic.API) {}
	gtic.TIC = func(a *gtic.API) {
		game.Input(a)
		game.Draw(a)
	}
	if err := gtic.Load(gtic.WithTitle("Sokoban")).Run(); err != nil {
		log.Println(err)
	}
}
