package game

import (
	"fmt"
	"gtic"
	"gtic/examples/games/sokoban/level"
	"image"
	"log"
	"path/filepath"

	"golang.org/x/image/colornames"
)

type Game struct {
	Packs          []level.LevelPack
	Progress       level.SavedProgress
	ProgressPath   string
	Statistics     level.SavedStatistics
	StatisticsPath string
	Attempt        int
	ResultRecorded bool
	SelectedPack   int
	SelectedLevel  int
	CurrentLevel   int
	Board          *Board
	Playing        bool
	mouseLeftDown  bool
	mousePath      []image.Point
	mousePathIndex int
	mousePathTick  int
	selectedBox    int
	boxTarget      *image.Point
	blockedTarget  *image.Point
	blockedTicks   int
	showLegend     bool
	status         string
}

const tileSize = 8
const mouseStepInterval = 5
const blockedBlinkDuration = 36

type LegendItem struct {
	Tile  Tile
	Label string
}

var SokobanLegend = []LegendItem{
	{Tile: Floor, Label: "Floor"},
	{Tile: Wall, Label: "Wall"},
	{Tile: Goal, Label: "Goal"},
	{Tile: Player, Label: "Player"},
	{Tile: Box, Label: "Box"},
	{Tile: GoalBox, Label: "Box on goal"},
	// {Tile: GoalPlayer, Label: "Player on goal"},
}

func NewGame() (*Game, error) {
	packs, err := level.LoadPacks()
	if err != nil {
		return nil, err
	}
	progress, progressPath, err := level.LoadProgress()
	if err != nil {
		return nil, err
	}
	statisticsPath := filepath.Join(filepath.Dir(progressPath), "sokoban-statistics.json")
	statistics, err := level.LoadStatistics(statisticsPath)
	if err != nil {
		return nil, err
	}
	for _, pack := range packs {
		progress[pack.File] = min(max(progress[pack.File], 0), len(pack.Levels))
	}
	return &Game{
		Packs:          packs,
		Progress:       progress,
		ProgressPath:   progressPath,
		Statistics:     statistics,
		StatisticsPath: statisticsPath,
		SelectedLevel:  min(progress[packs[0].File], len(packs[0].Levels)-1),
		selectedBox:    -1,
	}, nil
}

func (g *Game) Input(a *gtic.API) {
	mouseX, mouseY, leftDown, _, _, _, _ := a.Mouse()
	mouseClicked := leftDown && !g.mouseLeftDown
	g.mouseLeftDown = leftDown
	if g.blockedTicks > 0 {
		g.blockedTicks--
		if g.blockedTicks == 0 {
			g.blockedTarget = nil
		}
	}

	if !g.Playing {
		if a.KeyP(gtic.KeyH) {
			g.showLegend = !g.showLegend
			return
		}
		if g.showLegend && a.KeyP(gtic.KeyESC) {
			g.showLegend = false
			return
		}
		switch {
		case a.KeyP(gtic.KeyLEFT):
			g.SelectPack(-1)
			g.status = ""
		case a.KeyP(gtic.KeyRIGHT):
			g.SelectPack(1)
			g.status = ""
		case a.KeyP(gtic.KeyUP):
			g.SelectLevel(-1)
		case a.KeyP(gtic.KeyDOWN):
			g.SelectLevel(1)
		case a.KeyP(gtic.KeyRETURN):
			g.StartPack()
		}
		if a.KeyP(gtic.KeyESC) {
			a.Exit()
		}
		return
	}

	if g.showLegend {
		if a.KeyP(gtic.KeyH) || a.KeyP(gtic.KeyESC) {
			g.showLegend = false
		}
		return
	}
	if a.KeyP(gtic.KeyH) {
		g.showLegend = true
		return
	}
	if a.KeyP(gtic.KeyESC) && g.Playing {
		g.cancelMouseAnimation()
		g.clearBoxSelection()
		g.Playing = false
		g.SelectedLevel = g.CurrentLevel
		return
	}
	if a.KeyP(gtic.KeyR) {
		g.cancelMouseAnimation()
		g.StartAttempt()
		return
	}
	if g.Board.Win && a.KeyP(gtic.KeyRETURN) {
		g.NextLevel()
		return
	}

	if len(g.mousePath) > 0 {
		g.updateMouseAnimation()
		return
	}

	if mouseClicked && !g.Board.Win {
		cellX, cellY := mouseX/tileSize, mouseY/tileSize
		if mouseX >= 0 && mouseY >= 0 && cellX < g.Board.Dim.Row && cellY < g.Board.Dim.Column {
			g.handleBoardClick(cellX, cellY)
			return
		}
	}

	wasWin := g.Board.Win
	g.Board.Input(a)
	if !wasWin && g.Board.Win && !g.ResultRecorded {
		g.SaveCompletedLevel()
	}
}

func (g *Game) handleBoardClick(x, y int) {
	g.blockedTarget = nil
	g.blockedTicks = 0
	clicked := g.Board.Dim.Idx(x, y)
	if g.Board.Field.IsBoxAt(clicked) {
		if g.selectedBox == clicked {
			g.clearBoxSelection()
			g.status = "Box deselected."
			return
		}
		g.selectedBox = clicked
		g.boxTarget = nil
		g.status = "Box selected. Click a destination for it."
		return
	}
	if g.selectedBox < 0 || g.selectedBox >= len(g.Board.Field) || !g.Board.Field.IsBoxAt(g.selectedBox) {
		g.clearBoxSelection()
		g.blinkBlockedTarget(x, y, "Select a box first.")
		return
	}

	path, ok := g.Board.PathBoxTo(g.selectedBox, x, y)
	if !ok {
		g.blinkBlockedTarget(x, y, "No path to push the selected box there.")
		return
	}
	if len(path) == 0 {
		g.boxTarget = nil
		g.status = "Box is already at that destination."
		return
	}
	g.status = ""
	g.mousePath = path
	g.mousePathIndex = 0
	g.mousePathTick = 0
	target := image.Pt(x, y)
	g.boxTarget = &target
}

func (g *Game) blinkBlockedTarget(x, y int, status string) {
	target := image.Pt(x, y)
	g.blockedTarget = &target
	g.blockedTicks = blockedBlinkDuration
	g.status = status
}

func (g *Game) updateMouseAnimation() {
	g.mousePathTick++
	if g.mousePathTick < mouseStepInterval {
		return
	}
	g.mousePathTick = 0
	if g.mousePathIndex >= len(g.mousePath) {
		g.cancelMouseAnimation()
		g.boxTarget = nil
		return
	}
	direction := g.mousePath[g.mousePathIndex]
	oldBox := g.selectedBox
	boxPoint := image.Pt(g.Board.Dim.Pos(oldBox))
	playerPoint := image.Pt(g.Board.Dim.Pos(g.Board.Field.PlayerAt()))
	if !g.Board.Move(direction) {
		g.cancelMouseAnimation()
		g.boxTarget = nil
		return
	}
	if playerPoint.Add(direction) == boxPoint {
		boxPoint = boxPoint.Add(direction)
		g.selectedBox = g.Board.Dim.Idx(boxPoint.X, boxPoint.Y)
	}
	g.mousePathIndex++
	if g.Board.Win {
		g.cancelMouseAnimation()
		g.boxTarget = nil
		if !g.ResultRecorded {
			g.SaveCompletedLevel()
		}
		return
	}
	if g.mousePathIndex >= len(g.mousePath) {
		g.cancelMouseAnimation()
		g.boxTarget = nil
	}
}

func (g *Game) cancelMouseAnimation() {
	g.mousePath = nil
	g.mousePathIndex = 0
	g.mousePathTick = 0
}

func (g *Game) clearBoxSelection() {
	g.selectedBox = -1
	g.boxTarget = nil
	g.cancelMouseAnimation()
}

func (g *Game) SelectPack(direction int) {
	g.SelectedPack = (g.SelectedPack + direction + len(g.Packs)) % len(g.Packs)
	g.SelectedLevel = min(g.Progress[g.Packs[g.SelectedPack].File], len(g.Packs[g.SelectedPack].Levels)-1)
}

func (g *Game) SelectLevel(direction int) {
	pack := g.Packs[g.SelectedPack]
	g.SelectedLevel = (g.SelectedLevel + direction + len(pack.Levels)) % len(pack.Levels)
	g.status = ""
}

func (g *Game) StartPack() {
	g.CurrentLevel = g.SelectedLevel
	g.Playing = true
	g.status = ""
	g.StartAttempt()
}

func (g *Game) StartAttempt() {
	pack := g.Packs[g.SelectedPack]
	g.Board = NewBoard(pack.Levels[g.CurrentLevel])
	g.clearBoxSelection()
	g.blockedTarget = nil
	g.blockedTicks = 0
	if g.Statistics == nil {
		g.Statistics = level.SavedStatistics{}
	}
	g.Attempt = g.Statistics.StartAttempt(pack.File, g.CurrentLevel+1)
	g.ResultRecorded = false
	if err := level.SaveStatistics(g.StatisticsPath, g.Statistics); err != nil {
		g.status = "Statistics save failed; see log."
		log.Printf("save Sokoban statistics: %v", err)
		return
	}
	g.status = ""
}

func (g *Game) SaveCompletedLevel() {
	pack := g.Packs[g.SelectedPack]
	result := level.AttemptResult{
		Attempt: g.Attempt,
		Moves:   g.Board.Moves,
		Pushes:  g.Board.Pushes,
	}
	g.Statistics.RecordSuccess(pack.File, g.CurrentLevel+1, result)
	g.ResultRecorded = true
	if err := level.SaveStatistics(g.StatisticsPath, g.Statistics); err != nil {
		g.status = "Statistics save failed; see log."
		log.Printf("save Sokoban statistics: %v", err)
	}
	nextLevel := g.CurrentLevel + 1
	if nextLevel > g.Progress[pack.File] {
		g.Progress[pack.File] = nextLevel
		if err := level.SaveProgress(g.ProgressPath, g.Progress); err != nil {
			g.status = "Progress save failed; see log."
			log.Printf("save Sokoban progress: %v", err)
			return
		}
	}
}

func (g *Game) NextLevel() {
	pack := g.Packs[g.SelectedPack]
	if g.CurrentLevel+1 == len(pack.Levels) {
		g.Board = nil
		g.Playing = false
		g.SelectedLevel = g.CurrentLevel
		g.status = "Pack complete."
		return
	}
	g.CurrentLevel++
	g.SelectedLevel = g.CurrentLevel
	g.StartAttempt()
}

func (g *Game) Draw(a *gtic.API) {
	a.Cls()
	if g.showLegend {
		g.DrawLegend(a)
		return
	}
	if !g.Playing {
		g.drawPackMenu(a)
		return
	}

	g.Board.Draw(a)
	if g.selectedBox >= 0 && g.selectedBox < len(g.Board.Field) && g.Board.Field.IsBoxAt(g.selectedBox) && (g.mousePathTick/2)%2 == 0 {
		x, y := g.Board.Dim.Pos(g.selectedBox)
		a.RectB(x*tileSize, y*tileSize, tileSize, tileSize, gtic.NewColor(colornames.Lime))
	}
	if len(g.mousePath) > 0 && (g.mousePathTick%4) < 2 {
		player := image.Pt(g.Board.Dim.Pos(g.Board.Field.PlayerAt()))
		a.Circ(player.X*tileSize+tileSize/2, player.Y*tileSize+tileSize/2, 1, gtic.NewColor(colornames.Lime))
	}
	if g.boxTarget != nil && (g.mousePathTick/3)%2 == 0 {
		target := *g.boxTarget
		a.RectB(target.X*tileSize, target.Y*tileSize, tileSize, tileSize, gtic.NewColor(colornames.Lime))
	}
	if g.blockedTarget != nil && g.blockedTicks > 0 && (g.blockedTicks/4)%2 == 0 {
		target := *g.blockedTarget
		a.RectB(target.X*tileSize, target.Y*tileSize, tileSize, tileSize, gtic.NewColor(colornames.Red))
	}
	pack := g.Packs[g.SelectedPack]
	stats := g.Statistics.Level(pack.File, g.CurrentLevel+1)
	a.Print(fmt.Sprintf("%s %d/%d", pack.Name, g.CurrentLevel+1, len(pack.Levels)), 0, a.Bounds().Height-16, a.Pal(12))
	a.Print(fmt.Sprintf("Moves:%d Pushes:%d", g.Board.Moves, g.Board.Pushes), 0, a.Bounds().Height-32, a.Pal(12))
	a.Print(fmt.Sprintf("Try:%d Attempts:%d Won:%d", g.Attempt, stats.Attempts, stats.Successes), 0, a.Bounds().Height-40, a.Pal(12))
	if stats.Best != nil {
		a.Print(fmt.Sprintf("Best:%d moves/%d pushes", stats.Best.Moves, stats.Best.Pushes), 0, a.Bounds().Height-48, a.Pal(11))
	}
	if g.Board.Win {
		a.Print("Clear! Enter: next U: undo F: redo H: legend", 0, a.Bounds().Height-8, a.Pal(11))
	} else {
		a.Print("Click box, then destination | Arrows move | U undo | F redo", 0, a.Bounds().Height-8, a.Pal(12))
	}
	if g.status != "" {
		a.Print(g.status, 0, a.Bounds().Height-24, a.Pal(6))
	}
}

func (g *Game) drawPackMenu(a *gtic.API) {
	a.Print("SOKOBAN", 4, 4, a.Pal(12))
	a.Print("< "+g.Packs[g.SelectedPack].Name+" >", 4, 18, a.Pal(11))
	pack := g.Packs[g.SelectedPack]
	completed := g.Progress[pack.File]
	if completed == len(pack.Levels) {
		a.Print(fmt.Sprintf("All %d levels complete", completed), 4, 30, a.Pal(12))
	} else {
		a.Print(fmt.Sprintf("Progress: %d/%d", completed, len(pack.Levels)), 4, 30, a.Pal(12))
	}
	a.Print(fmt.Sprintf("Level %d/%d", g.SelectedLevel+1, len(pack.Levels)), 4, 44, a.Pal(11))
	g.drawLevelPreview(a, pack.Levels[g.SelectedLevel], 8, 60, 104, 42)

	stats := g.Statistics.Level(pack.File, g.SelectedLevel+1)
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
	a.Print("Enter: play  H: elements  Esc: quit", 4, 128, a.Pal(12))
	if g.status != "" {
		a.Print(g.status, 120, 44, a.Pal(6))
	}
}

func (g *Game) DrawLegend(a *gtic.API) {
	a.Print("SOKOBAN - FIELD ELEMENTS", 8, 6, gtic.NewColor(colornames.White))
	for i, item := range SokobanLegend {
		x := 12 + (i/3)*80
		y := 28 + (i%3)*16
		drawLegendTile(a, x, y, item.Tile)
		a.Print(item.Label, x+22, y+2, gtic.NewColor(colornames.White))
	}
	a.Print("XSB: # wall, . goal, @ player, $ box", 8, 80, gtic.NewColor(colornames.Lightgray))
	a.Print("+: player on goal, *: box on goal", 8, 92, gtic.NewColor(colornames.Lightgray))
	a.Print("H or Esc: back", 8, a.Bounds().Height-16, gtic.NewColor(colornames.White))
}

func drawLegendTile(a *gtic.API, x, y int, tile Tile) {
	background := Tile(Floor).Color()
	switch tile {
	case Wall:
		background = Tile(Wall).Color()
	case Goal, GoalBox, GoalPlayer:
		background = Tile(FloorGoal).Color()
	}
	a.Rect(x, y, 8, 8, gtic.NewColor(background))
	switch tile {
	case Goal:
		drawGoal(a, x, y)
	case GoalBox:
		drawBox(a, x, y, true)
	case Box:
		drawBox(a, x, y, false)
	case Player, GoalPlayer:
		drawPlayer(a, x, y, LoadSprites(a))
	}
}

func (g *Game) drawLevelPreview(a *gtic.API, level string, x, y, maxWidth, maxHeight int) {
	board := NewBoard(level)
	cellSize := max(1, min(maxWidth/board.Dim.Row, maxHeight/board.Dim.Column))
	for i, cell := range board.Field {
		cellX, cellY := board.Dim.Pos(i)
		cellX = x + cellX*cellSize
		cellY = y + cellY*cellSize
		background := Tile(Floor).Color()
		if cell.Goal {
			background = Tile(Goal).Color()
		}
		if cell.Get() == Wall {
			background = Tile(Wall).Color()
		}
		a.Rect(cellX, cellY, cellSize, cellSize, gtic.NewColor(background))
		switch {
		case cell.Player:
			a.Circ(cellX+cellSize/2, cellY+cellSize/2, max(1, cellSize/2-1), gtic.NewColor(Tile(Player).Color()))
		case cell.Box:
			a.Circ(cellX+cellSize/2, cellY+cellSize/2, max(1, cellSize/2-1), gtic.NewColor(Tile(Box).Color()))
		}
	}
}
