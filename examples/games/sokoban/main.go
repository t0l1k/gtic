package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"gtic"
	"gtic/react"
	"image"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed assets/levels
var levelFiles embed.FS

type levelPack struct {
	name   string
	file   string
	levels []string
}

func loadPacks() ([]levelPack, error) {
	return loadPacksFrom(levelFiles, "assets/levels")
}

func loadPacksFrom(source fs.FS, dir string) ([]levelPack, error) {
	entries, err := fs.ReadDir(source, dir)
	if err != nil {
		return nil, fmt.Errorf("read level collection directory %q: %w", dir, err)
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".xsb" || ext == ".txt" {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	if len(files) == 0 {
		return nil, fmt.Errorf("no .xsb or .txt level collections found in assets/levels")
	}

	packs := make([]levelPack, 0, len(files))
	for _, file := range files {
		levelPath := path.Join(dir, file)
		data, err := fs.ReadFile(source, levelPath)
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
	metadata := false
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
		if strings.HasPrefix(trimmed, ";") {
			if err := flush(); err != nil {
				return nil, err
			}
			metadata = false
			continue
		}
		if trimmed == "" {
			if err := flush(); err != nil {
				return nil, err
			}
			continue
		}
		if isCollectionMetadata(trimmed) {
			if err := flush(); err != nil {
				return nil, err
			}
			metadata = true
			continue
		}
		if metadata && !isLevelRow(line) {
			continue
		}
		if !isLevelRow(line) {
			return nil, fmt.Errorf("invalid level row on line %d", lineNumber)
		}
		metadata = false
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

func isCollectionMetadata(line string) bool {
	lower := strings.ToLower(line)
	for _, field := range []string{"title:", "description:", "author:", "email:", "website:"} {
		if strings.HasPrefix(lower, field) {
			return true
		}
	}
	return false
}

func isLevelRow(line string) bool {
	for _, tile := range line {
		switch tile {
		case ' ', '#', '.', '@', '+', '$', '*':
		default:
			return false
		}
	}
	return len(line) > 0
}

type savedProgress map[string]int

type attemptResult struct {
	Attempt int `json:"attempt"`
	Moves   int `json:"moves"`
	Pushes  int `json:"pushes"`
}

type levelStatistics struct {
	Attempts  int             `json:"attempts"`
	Successes int             `json:"successes"`
	Best      *attemptResult  `json:"best,omitempty"`
	Results   []attemptResult `json:"results,omitempty"`
}

type savedStatistics map[string]map[int]levelStatistics

func loadStatistics(path string) (savedStatistics, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return savedStatistics{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read Sokoban statistics: %w", err)
	}
	var stats savedStatistics
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("decode Sokoban statistics: %w", err)
	}
	if stats == nil {
		stats = savedStatistics{}
	}
	return stats, nil
}

func saveStatistics(path string, stats savedStatistics) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("encode Sokoban statistics: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create statistics directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write Sokoban statistics: %w", err)
	}
	return nil
}

func (s savedStatistics) level(pack string, level int) levelStatistics {
	return s[pack][level]
}

func (s savedStatistics) startAttempt(pack string, level int) int {
	if s[pack] == nil {
		s[pack] = make(map[int]levelStatistics)
	}
	stats := s[pack][level]
	stats.Attempts++
	s[pack][level] = stats
	return stats.Attempts
}

func (s savedStatistics) recordSuccess(pack string, level int, result attemptResult) {
	if s[pack] == nil {
		s[pack] = make(map[int]levelStatistics)
	}
	stats := s[pack][level]
	stats.Successes++
	stats.Results = append(stats.Results, result)
	if stats.Best == nil || result.Moves < stats.Best.Moves ||
		(result.Moves == stats.Best.Moves && result.Pushes < stats.Best.Pushes) {
		best := result
		stats.Best = &best
	}
	s[pack][level] = stats
}

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
	Floor  Tile = iota
	Wall        // #
	Goal        // .
	Player      // @
	Box         // $
)

func (t Tile) String() string { return []string{" ", "#", ".", "@", "$"}[t] }
func (t Tile) Color() int     { return []int{13, 15, 3, 10, 2}[t] }

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
	Dim        Dim
	Field      Field
	win        bool
	level      string
	Moves      int
	Pushes     int
	history    []boardSnapshot
	historyPos int
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
		Field: NewField().Load(level),
		Dim:   NewDim(level),
		level: level,
	}
	b.recordHistory()
	return b
}
func (b *Board) Reset() {
	b.Field = NewField().Load(b.level)
	b.Dim = NewDim(b.level)
	b.win = false
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
func (b *Board) recordHistory() {
	snapshot := boardSnapshot{
		cells:  make([]cellSnapshot, len(b.Field)),
		moves:  b.Moves,
		pushes: b.Pushes,
		win:    b.win,
	}
	for i, cell := range b.Field {
		snapshot.cells[i] = cellSnapshot{player: cell.player, box: cell.box}
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
		cell.player = snapshot.cells[i].player
		cell.box = snapshot.cells[i].box
	}
	b.Moves = snapshot.moves
	b.Pushes = snapshot.pushes
	b.win = snapshot.win
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
	if newPt.X < 0 || newPt.Y < 0 || newPt.X >= b.Dim.row || newPt.Y >= b.Dim.column {
		return false
	}
	newCell := b.Field.Cell(b.Dim.Idx(newPt.X, newPt.Y))
	if newCell.Get() == Wall || newCell.box {
		return false
	}
	boxIdx := b.Dim.Idx(boxPt.X, boxPt.Y)
	b.Field.MoveBox(b.Dim.Idx(newPt.X, newPt.Y), boxIdx)
	b.Field.MovePlayer(boxIdx, b.Field.PlayerAt())
	b.Moves++
	b.Pushes++
	b.win = b.IsWin()
	b.recordHistory()
	return true
}
func (b *Board) Move(direction image.Point) bool {
	playerIdx := b.Field.PlayerAt()
	if playerIdx < 0 {
		return false
	}
	playerPt := image.Pt(b.Dim.Pos(playerIdx))
	newPt := playerPt.Add(direction)
	if newPt.X < 0 || newPt.Y < 0 || newPt.X >= b.Dim.row || newPt.Y >= b.Dim.column {
		return false
	}
	newIdx := b.Dim.Idx(newPt.X, newPt.Y)
	newCell := b.Field.Cell(newIdx)
	switch newCell.Get() {
	case Wall:
		return false
	case Floor:
		if newCell.box {
			return b.PushBox(newIdx, direction)
		}
	}
	b.Field.MovePlayer(newIdx, playerIdx)
	b.Moves++
	b.win = b.IsWin()
	b.recordHistory()
	return true
}
func (b *Board) Draw(a *gtic.API) {
	for i, cell := range b.Field {
		x, y := b.Dim.Pos(i)
		switch {
		case cell.Get() == Wall:
			a.Rect(x*8, y*8, 8, 8, a.Pal(cell.Color()))
		case cell.Get() == Floor:
			a.Rect(x*8, y*8, 8, 8, a.Pal(Tile(Floor).Color()))
			switch {
			case cell.box:
				a.Circ(x*8+4, y*8+4, 4, a.Pal(cell.Color()))
			case cell.player:
				a.Circ(x*8+4, y*8+4, 4, a.Pal(cell.Color()))
			case cell.goal:
				a.Rect(x*8, y*8, 8, 8, a.Pal(Tile(Floor).Color()))
				a.Circ(x*8+4, y*8+4, 4, a.Pal(cell.Color()))
			}
		}
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
	packs          []levelPack
	progress       savedProgress
	progressPath   string
	statistics     savedStatistics
	statisticsPath string
	attempt        int
	resultRecorded bool
	selectedPack   int
	selectedLevel  int
	currentLevel   int
	board          *Board
	playing        bool
	status         string
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
	statisticsPath := filepath.Join(filepath.Dir(progressPath), "sokoban-statistics.json")
	statistics, err := loadStatistics(statisticsPath)
	if err != nil {
		return nil, err
	}
	for _, pack := range packs {
		progress[pack.file] = min(max(progress[pack.file], 0), len(pack.levels))
	}
	return &game{
		packs:          packs,
		progress:       progress,
		progressPath:   progressPath,
		statistics:     statistics,
		statisticsPath: statisticsPath,
		selectedLevel:  min(progress[packs[0].file], len(packs[0].levels)-1),
	}, nil
}

func (g *game) Input(a *gtic.API) {
	if !g.playing {
		switch {
		case a.KeyP(gtic.KeyLEFT):
			g.selectPack(-1)
			g.status = ""
		case a.KeyP(gtic.KeyRIGHT):
			g.selectPack(1)
			g.status = ""
		case a.KeyP(gtic.KeyUP):
			g.selectLevel(-1)
		case a.KeyP(gtic.KeyDOWN):
			g.selectLevel(1)
		case a.KeyP(gtic.KeyRETURN):
			g.startPack()
		}
		if a.KeyP(gtic.KeyESC) {
			a.Exit()
		}
		return
	}

	if a.KeyP(gtic.KeyESC) && g.playing {
		g.playing = false
		g.selectedLevel = g.currentLevel
		return
	}
	if a.KeyP(gtic.KeyR) {
		g.startAttempt()
		return
	}
	if g.board.win && a.KeyP(gtic.KeyRETURN) {
		g.nextLevel()
		return
	}

	wasWin := g.board.win
	g.board.Input(a)
	if !wasWin && g.board.win && !g.resultRecorded {
		g.saveCompletedLevel()
	}
}

func (g *game) selectPack(direction int) {
	g.selectedPack = (g.selectedPack + direction + len(g.packs)) % len(g.packs)
	g.selectedLevel = min(g.progress[g.packs[g.selectedPack].file], len(g.packs[g.selectedPack].levels)-1)
}

func (g *game) selectLevel(direction int) {
	pack := g.packs[g.selectedPack]
	g.selectedLevel = (g.selectedLevel + direction + len(pack.levels)) % len(pack.levels)
	g.status = ""
}

func (g *game) startPack() {
	g.currentLevel = g.selectedLevel
	g.playing = true
	g.status = ""
	g.startAttempt()
}

func (g *game) startAttempt() {
	pack := g.packs[g.selectedPack]
	g.board = NewBoard(pack.levels[g.currentLevel])
	if g.statistics == nil {
		g.statistics = savedStatistics{}
	}
	g.attempt = g.statistics.startAttempt(pack.file, g.currentLevel+1)
	g.resultRecorded = false
	if err := saveStatistics(g.statisticsPath, g.statistics); err != nil {
		g.status = "Statistics save failed; see log."
		log.Printf("save Sokoban statistics: %v", err)
		return
	}
	g.status = ""
}

func (g *game) saveCompletedLevel() {
	pack := g.packs[g.selectedPack]
	result := attemptResult{
		Attempt: g.attempt,
		Moves:   g.board.Moves,
		Pushes:  g.board.Pushes,
	}
	g.statistics.recordSuccess(pack.file, g.currentLevel+1, result)
	g.resultRecorded = true
	if err := saveStatistics(g.statisticsPath, g.statistics); err != nil {
		g.status = "Statistics save failed; see log."
		log.Printf("save Sokoban statistics: %v", err)
	}
	nextLevel := g.currentLevel + 1
	if nextLevel > g.progress[pack.file] {
		g.progress[pack.file] = nextLevel
		if err := saveProgress(g.progressPath, g.progress); err != nil {
			g.status = "Progress save failed; see log."
			log.Printf("save Sokoban progress: %v", err)
			return
		}
	}
}

func (g *game) nextLevel() {
	pack := g.packs[g.selectedPack]
	if g.currentLevel+1 == len(pack.levels) {
		g.board = nil
		g.playing = false
		g.selectedLevel = g.currentLevel
		g.status = "Pack complete."
		return
	}
	g.currentLevel++
	g.selectedLevel = g.currentLevel
	g.startAttempt()
}

func (g *game) Draw(a *gtic.API) {
	a.Cls()
	if !g.playing {
		g.drawPackMenu(a)
		return
	}

	g.board.Draw(a)
	pack := g.packs[g.selectedPack]
	stats := g.statistics.level(pack.file, g.currentLevel+1)
	a.Print(fmt.Sprintf("%s %d/%d", pack.name, g.currentLevel+1, len(pack.levels)), 0, a.Bounds().Height-16, a.Pal(12))
	a.Print(fmt.Sprintf("Moves:%d Pushes:%d", g.board.Moves, g.board.Pushes), 0, a.Bounds().Height-32, a.Pal(12))
	a.Print(fmt.Sprintf("Try:%d Attempts:%d Won:%d", g.attempt, stats.Attempts, stats.Successes), 0, a.Bounds().Height-40, a.Pal(12))
	if stats.Best != nil {
		a.Print(fmt.Sprintf("Best:%d moves/%d pushes", stats.Best.Moves, stats.Best.Pushes), 0, a.Bounds().Height-48, a.Pal(11))
	}
	if g.board.win {
		a.Print("Clear! Enter: next U: undo F: redo", 0, a.Bounds().Height-8, a.Pal(11))
	} else {
		a.Print("Arrows: move U: undo F: redo R: retry", 0, a.Bounds().Height-8, a.Pal(12))
	}
	if g.status != "" {
		a.Print(g.status, 0, a.Bounds().Height-24, a.Pal(6))
	}
}

func (g *game) drawPackMenu(a *gtic.API) {
	a.Print("SOKOBAN", 4, 4, a.Pal(12))
	a.Print("< "+g.packs[g.selectedPack].name+" >", 4, 18, a.Pal(11))
	pack := g.packs[g.selectedPack]
	completed := g.progress[pack.file]
	if completed == len(pack.levels) {
		a.Print(fmt.Sprintf("All %d levels complete", completed), 4, 30, a.Pal(12))
	} else {
		a.Print(fmt.Sprintf("Progress: %d/%d", completed, len(pack.levels)), 4, 30, a.Pal(12))
	}
	a.Print(fmt.Sprintf("Level %d/%d", g.selectedLevel+1, len(pack.levels)), 4, 44, a.Pal(11))
	g.drawLevelPreview(a, pack.levels[g.selectedLevel], 8, 60, 104, 42)

	stats := g.statistics.level(pack.file, g.selectedLevel+1)
	a.Print("Statistics", 120, 60, a.Pal(12))
	a.Print(fmt.Sprintf("Attempts: %d", stats.Attempts), 120, 72, a.Pal(12))
	a.Print(fmt.Sprintf("Won: %d", stats.Successes), 120, 84, a.Pal(12))
	if stats.Best == nil {
		a.Print("Best: --", 120, 96, a.Pal(11))
	} else {
		a.Print(fmt.Sprintf("Best: %d moves", stats.Best.Moves), 120, 96, a.Pal(11))
		a.Print(fmt.Sprintf("%d pushes (#%d)", stats.Best.Pushes, stats.Best.Attempt), 120, 108, a.Pal(11))
	}
	a.Print("Left/Right: collection  Up/Down: level", 4, 116, a.Pal(12))
	a.Print("Enter: play  Esc: quit", 4, 128, a.Pal(12))
	if g.status != "" {
		a.Print(g.status, 120, 44, a.Pal(6))
	}
}

func (g *game) drawLevelPreview(a *gtic.API, level string, x, y, maxWidth, maxHeight int) {
	board := NewBoard(level)
	cellSize := max(1, min(maxWidth/board.Dim.row, maxHeight/board.Dim.column))
	for i, cell := range board.Field {
		cellX, cellY := board.Dim.Pos(i)
		cellX = x + cellX*cellSize
		cellY = y + cellY*cellSize
		background := Tile(Floor).Color()
		if cell.goal {
			background = Tile(Goal).Color()
		}
		if cell.Get() == Wall {
			background = Tile(Wall).Color()
		}
		a.Rect(cellX, cellY, cellSize, cellSize, a.Pal(background))
		switch {
		case cell.player:
			a.Circ(cellX+cellSize/2, cellY+cellSize/2, max(1, cellSize/2-1), a.Pal(Tile(Player).Color()))
		case cell.box:
			a.Circ(cellX+cellSize/2, cellY+cellSize/2, max(1, cellSize/2-1), a.Pal(Tile(Box).Color()))
		}
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
