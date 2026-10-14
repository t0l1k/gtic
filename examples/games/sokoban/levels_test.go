package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPacks(t *testing.T) {
	packs, err := loadPacks()
	if err != nil {
		t.Fatal(err)
	}
	if len(packs) != 2 {
		t.Fatalf("got %d packs, want 2", len(packs))
	}
	if packs[0].name != "simple" || len(packs[0].levels) != 3 {
		t.Errorf("simple pack: got name %q and %d levels, want simple and 3", packs[0].name, len(packs[0].levels))
	}
	if packs[1].name != "xsokoban" || len(packs[1].levels) != 90 {
		t.Errorf("xsokoban pack: got name %q and %d levels, want xsokoban and 90", packs[1].name, len(packs[1].levels))
	}
	for _, pack := range packs {
		for i, level := range pack.levels {
			board := NewBoard(level)
			if got, want := len(board.Field), board.Dim.row*board.Dim.column; got != want {
				t.Errorf("%s level %d: field has %d cells, want %d", pack.name, i+1, got, want)
			}
			playerCount := 0
			for _, cell := range board.Field {
				if cell.player {
					playerCount++
				}
			}
			if playerCount != 1 {
				t.Errorf("%s level %d has %d players, want exactly one", pack.name, i+1, playerCount)
			}
		}
	}
}

func TestParseLevelsPadsRowsAndSupportsOccupiedGoals(t *testing.T) {
	levels, err := parseLevels(";level.1\n##\n#*@#\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(levels) != 1 {
		t.Fatalf("got %d levels, want 1", len(levels))
	}
	board := NewBoard(levels[0])
	if board.Dim != (Dim{row: 4, column: 2}) {
		t.Errorf("got dimensions %+v, want 4x2", board.Dim)
	}
	if !board.IsWin() {
		t.Error("box on goal should be treated as completed")
	}
}

func TestParseLevelsRejectsMultiplePlayers(t *testing.T) {
	if _, err := parseLevels("#@+#"); err == nil {
		t.Fatal("expected an error for a level with multiple players")
	}
}

func TestCellStringShowsOccupiedGoals(t *testing.T) {
	playerOnGoal := &Cell{player: true, goal: true}
	if got := playerOnGoal.String(); got != "+" {
		t.Errorf("player on goal rendered as %q, want \"+\"", got)
	}

	boxOnGoal := &Cell{box: true, goal: true}
	if got := boxOnGoal.String(); got != "*" {
		t.Errorf("box on goal rendered as %q, want \"*\"", got)
	}
}

func TestSaveProgress(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config", "sokoban-progress.json")
	want := savedProgress{"simple.xsb": 2, "xsokoban.xsb": 7}
	if err := saveProgress(path, want); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got savedProgress
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	for pack, wantLevel := range want {
		if got[pack] != wantLevel {
			t.Errorf("%s progress: got %d, want %d", pack, got[pack], wantLevel)
		}
	}
}

func TestStartPackResumesFirstUncompletedLevel(t *testing.T) {
	g := &game{
		packs: []levelPack{{
			file:   "simple.xsb",
			levels: []string{"#@#", "#@#", "#@#"},
		}},
		progress: savedProgress{"simple.xsb": 2},
	}

	g.startPack()

	if g.currentLevel != 2 {
		t.Errorf("started at level %d, want first uncompleted level 2", g.currentLevel)
	}
	if !g.playing || g.board == nil {
		t.Fatal("starting a pack should load its board")
	}
}

func TestStartCompletedPackReplaysFromFirstLevel(t *testing.T) {
	g := &game{
		packs: []levelPack{{
			file:   "simple.xsb",
			levels: []string{"#@#", "#@#"},
		}},
		progress: savedProgress{"simple.xsb": 2},
	}

	g.startPack()

	if g.currentLevel != 0 {
		t.Errorf("started at level %d, want level 0 for replay", g.currentLevel)
	}
}

func TestPushBoxDoesNotPushIntoAnotherBox(t *testing.T) {
	board := NewBoard("#######\n#@$$. #\n#######")

	board.Move(Right)

	if got := board.Field.PlayerAt(); got != board.Dim.Idx(1, 1) {
		t.Errorf("player moved to cell %d, want %d", got, board.Dim.Idx(1, 1))
	}
	if !board.Field.IsBoxAt(board.Dim.Idx(2, 1)) || !board.Field.IsBoxAt(board.Dim.Idx(3, 1)) {
		t.Error("blocked push changed the positions of adjacent boxes")
	}
	if board.moves != 0 || board.pushes != 0 || len(board.history) != 0 {
		t.Errorf("blocked push changed counters/history: moves=%d pushes=%d history=%d", board.moves, board.pushes, len(board.history))
	}
}

func TestPushBoxMovesIntoEmptyCell(t *testing.T) {
	board := NewBoard("#######\n#@$ . #\n#######")

	board.Move(Right)

	if got := board.Field.PlayerAt(); got != board.Dim.Idx(2, 1) {
		t.Errorf("player moved to cell %d, want %d", got, board.Dim.Idx(2, 1))
	}
	if !board.Field.IsBoxAt(board.Dim.Idx(3, 1)) {
		t.Error("box was not moved into the empty cell")
	}
}

func TestMoveAndPushCountersCanBeUndone(t *testing.T) {
	board := NewBoard("#######\n#@ $ .#\n#######")

	board.Move(Right)
	board.Move(Right)
	if board.moves != 2 || board.pushes != 1 {
		t.Fatalf("after moving and pushing got moves=%d pushes=%d, want 2 and 1", board.moves, board.pushes)
	}

	if !board.Undo() {
		t.Fatal("undo should restore the push")
	}
	if got := board.Field.PlayerAt(); got != board.Dim.Idx(2, 1) {
		t.Errorf("after undoing push player is at %d, want %d", got, board.Dim.Idx(2, 1))
	}
	if !board.Field.IsBoxAt(board.Dim.Idx(3, 1)) {
		t.Error("undoing push did not restore the box")
	}
	if board.moves != 1 || board.pushes != 0 {
		t.Errorf("after undoing push got moves=%d pushes=%d, want 1 and 0", board.moves, board.pushes)
	}

	if !board.Undo() {
		t.Fatal("undo should restore the previous player move")
	}
	if got := board.Field.PlayerAt(); got != board.Dim.Idx(1, 1) {
		t.Errorf("after undoing move player is at %d, want %d", got, board.Dim.Idx(1, 1))
	}
	if board.moves != 0 || board.pushes != 0 {
		t.Errorf("after undoing all moves got moves=%d pushes=%d, want 0 and 0", board.moves, board.pushes)
	}
	if board.Undo() {
		t.Error("undo should return false when history is empty")
	}
}

func TestUndoRestoresBoardAfterWinningPush(t *testing.T) {
	board := NewBoard("#######\n#@$.  #\n#######")

	board.Move(Right)

	if !board.win || !board.IsWin() {
		t.Fatal("push onto the goal should win the level")
	}
	if !board.Undo() {
		t.Fatal("undo should restore a winning push")
	}
	if board.win || board.IsWin() {
		t.Error("undo should clear the win state")
	}
	if board.moves != 0 || board.pushes != 0 {
		t.Errorf("after undoing winning push got moves=%d pushes=%d, want 0 and 0", board.moves, board.pushes)
	}
}

func TestResetClearsMoveHistoryAndCounters(t *testing.T) {
	board := NewBoard("#######\n#@ $ .#\n#######")
	board.Move(Right)
	board.Move(Right)

	board.Reset()

	if board.moves != 0 || board.pushes != 0 || len(board.history) != 0 {
		t.Errorf("reset left moves=%d pushes=%d history=%d", board.moves, board.pushes, len(board.history))
	}
}
