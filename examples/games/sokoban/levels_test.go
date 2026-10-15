package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoadPacks(t *testing.T) {
	packs, err := loadPacks()
	if err != nil {
		t.Fatal(err)
	}
	entries, err := levelFiles.ReadDir("assets/levels")
	if err != nil {
		t.Fatal(err)
	}
	expected := 0
	for _, entry := range entries {
		if !entry.IsDir() && (strings.EqualFold(filepath.Ext(entry.Name()), ".xsb") || strings.EqualFold(filepath.Ext(entry.Name()), ".txt")) {
			expected++
		}
	}
	if len(packs) != expected {
		t.Fatalf("loaded %d packs, want all %d supported collections", len(packs), expected)
	}
	for _, pack := range packs {
		if len(pack.levels) == 0 {
			t.Errorf("%s has no levels", pack.file)
		}
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

func TestParseTextCollectionMetadata(t *testing.T) {
	data := `Title: Test collection
Description: A test collection
             with wrapped metadata
Author: Test

Title: First level
#####
#@$.#
#####

Title: Second level
#####
#@*.#
#####
`
	levels, err := parseLevels(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(levels) != 2 {
		t.Fatalf("parsed %d levels, want 2", len(levels))
	}
}

func TestLoadPacksDiscoversCollectionsFromDirectory(t *testing.T) {
	files := fstest.MapFS{
		"collections/new-pack.xsb": &fstest.MapFile{Data: []byte("#@#")},
		"collections/another.txt":  &fstest.MapFile{Data: []byte("Title: Another\n#####\n#@$.#\n#####")},
		"collections/readme.md":    &fstest.MapFile{Data: []byte("not a Sokoban collection")},
	}

	packs, err := loadPacksFrom(files, "collections")
	if err != nil {
		t.Fatal(err)
	}
	if len(packs) != 2 {
		t.Fatalf("loaded %d packs, want 2 supported collections", len(packs))
	}
	if packs[0].file != "another.txt" || packs[1].file != "new-pack.xsb" {
		t.Errorf("discovered packs in unexpected order: %q, %q", packs[0].file, packs[1].file)
	}
}

func TestLoadPacksFromDirectoryRequiresSupportedCollection(t *testing.T) {
	files := fstest.MapFS{
		"collections/readme.md": &fstest.MapFile{Data: []byte("not a Sokoban collection")},
	}
	if _, err := loadPacksFrom(files, "collections"); err == nil {
		t.Fatal("expected an error when directory has no supported collections")
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

func TestStartCompletedPackCanReplaySelectedLevel(t *testing.T) {
	g := &game{
		packs: []levelPack{{
			file:   "simple.xsb",
			levels: []string{"#@#", "#@#"},
		}},
		progress:       savedProgress{"simple.xsb": 2},
		statistics:     savedStatistics{},
		statisticsPath: filepath.Join(t.TempDir(), "sokoban-statistics.json"),
		selectedLevel:  1,
	}

	g.startPack()

	if g.currentLevel != 1 {
		t.Errorf("started at level %d, want selected level 1", g.currentLevel)
	}
}

func TestSelectingPackAndLevel(t *testing.T) {
	g := &game{
		packs: []levelPack{
			{file: "simple.xsb", levels: []string{"#@#", "#@#", "#@#"}},
			{file: "xsokoban.xsb", levels: []string{"#@#", "#@#"}},
		},
		progress: savedProgress{"simple.xsb": 2, "xsokoban.xsb": 1},
	}

	g.selectPack(1)
	if g.selectedPack != 1 || g.selectedLevel != 1 {
		t.Errorf("after selecting xsokoban: pack=%d level=%d, want pack 1 level 1", g.selectedPack, g.selectedLevel)
	}
	g.selectLevel(1)
	if g.selectedLevel != 0 {
		t.Errorf("level selection did not wrap: got %d, want 0", g.selectedLevel)
	}
	g.selectPack(-1)
	if g.selectedPack != 0 || g.selectedLevel != 2 {
		t.Errorf("after returning to simple: pack=%d level=%d, want pack 0 level 2", g.selectedPack, g.selectedLevel)
	}
}

func TestStartPackLaunchesSelectedLevel(t *testing.T) {
	g := &game{
		packs: []levelPack{{
			file:   "simple.xsb",
			levels: []string{"#####\n#@  #\n#####", "#####\n# @ #\n#####"},
		}},
		progress:       savedProgress{"simple.xsb": 0},
		statistics:     savedStatistics{},
		statisticsPath: filepath.Join(t.TempDir(), "sokoban-statistics.json"),
		selectedLevel:  1,
	}

	g.startPack()

	if g.currentLevel != 1 || g.attempt != 1 {
		t.Errorf("started level=%d attempt=%d, want selected level 1 attempt 1", g.currentLevel, g.attempt)
	}
	if got := g.board.Field.PlayerAt(); got != g.board.Dim.Idx(2, 1) {
		t.Errorf("player starts at cell %d, want %d", got, g.board.Dim.Idx(2, 1))
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

func TestUndoRedoRestoresMovesAndPushes(t *testing.T) {
	board := NewBoard("#######\n#@$.  #\n#######")

	if !board.Move(Right) {
		t.Fatal("expected push to succeed")
	}
	if board.Moves != 1 || board.Pushes != 1 || !board.win {
		t.Fatalf("after push: moves=%d pushes=%d win=%v", board.Moves, board.Pushes, board.win)
	}
	if !board.Undo() {
		t.Fatal("expected undo to succeed")
	}
	if board.Moves != 0 || board.Pushes != 0 || board.win {
		t.Fatalf("after undo: moves=%d pushes=%d win=%v", board.Moves, board.Pushes, board.win)
	}
	if board.Field.PlayerAt() != board.Dim.Idx(1, 1) || !board.Field.IsBoxAt(board.Dim.Idx(2, 1)) {
		t.Fatal("undo did not restore the original player and box positions")
	}
	if !board.Redo() {
		t.Fatal("expected redo to succeed")
	}
	if board.Moves != 1 || board.Pushes != 1 || !board.win {
		t.Fatalf("after redo: moves=%d pushes=%d win=%v", board.Moves, board.Pushes, board.win)
	}
}

func TestNewMoveAfterUndoClearsRedoHistory(t *testing.T) {
	board := NewBoard("######\n#@  .#\n######")
	board.Move(Right)
	board.Move(Right)
	if !board.Undo() {
		t.Fatal("expected undo to succeed")
	}
	if !board.Move(Left) {
		t.Fatal("expected new move to succeed")
	}
	if board.Redo() {
		t.Fatal("redo should not be available after branching with a new move")
	}
}

func TestSaveStatisticsAndBestSuccessfulResult(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config", "sokoban-statistics.json")
	stats := savedStatistics{}
	if got := stats.startAttempt("simple.xsb", 3); got != 1 {
		t.Fatalf("first attempt number = %d, want 1", got)
	}
	if got := stats.startAttempt("simple.xsb", 3); got != 2 {
		t.Fatalf("second attempt number = %d, want 2", got)
	}
	stats.startAttempt("simple.xsb", 3)
	stats.startAttempt("simple.xsb", 3)
	stats.recordSuccess("simple.xsb", 3, attemptResult{Attempt: 2, Moves: 12, Pushes: 5})
	stats.recordSuccess("simple.xsb", 3, attemptResult{Attempt: 3, Moves: 10, Pushes: 6})
	stats.recordSuccess("simple.xsb", 3, attemptResult{Attempt: 4, Moves: 10, Pushes: 4})
	if err := saveStatistics(path, stats); err != nil {
		t.Fatal(err)
	}
	loaded, err := loadStatistics(path)
	if err != nil {
		t.Fatal(err)
	}
	got := loaded.level("simple.xsb", 3)
	if got.Attempts != 4 || got.Successes != 3 {
		t.Errorf("got attempts=%d successes=%d, want 4 attempts and 3 successes", got.Attempts, got.Successes)
	}
	if got.Best == nil || got.Best.Attempt != 4 || got.Best.Moves != 10 || got.Best.Pushes != 4 {
		t.Errorf("got best result %+v, want attempt 4 with 10 moves and 4 pushes", got.Best)
	}
	if len(got.Results) != 3 {
		t.Errorf("got %d successful results, want 3", len(got.Results))
	}
}

func TestStartingAttemptPersistsAttemptCount(t *testing.T) {
	dir := t.TempDir()
	g := &game{
		packs: []levelPack{{
			file:   "simple.xsb",
			levels: []string{"#####\n#@$.#\n#####"},
		}},
		progress:       savedProgress{},
		statistics:     savedStatistics{},
		statisticsPath: filepath.Join(dir, "sokoban-statistics.json"),
		playing:        true,
	}
	g.startAttempt()
	g.startAttempt()

	loaded, err := loadStatistics(g.statisticsPath)
	if err != nil {
		t.Fatal(err)
	}
	if got := loaded.level("simple.xsb", 1).Attempts; got != 2 {
		t.Errorf("saved attempt count = %d, want 2", got)
	}
	if g.attempt != 2 || g.board.Moves != 0 || g.board.Pushes != 0 {
		t.Errorf("retry state: attempt=%d moves=%d pushes=%d", g.attempt, g.board.Moves, g.board.Pushes)
	}
}

func TestSavingCompletedLevelStoresAttemptResult(t *testing.T) {
	dir := t.TempDir()
	g := &game{
		packs: []levelPack{{
			file:   "simple.xsb",
			levels: []string{"#####\n#@$.#\n#####"},
		}},
		progress:       savedProgress{},
		statistics:     savedStatistics{},
		statisticsPath: filepath.Join(dir, "sokoban-statistics.json"),
		progressPath:   filepath.Join(dir, "sokoban-progress.json"),
		currentLevel:   0,
		attempt:        3,
		board:          NewBoard("#####\n#@$.#\n#####"),
	}
	g.statistics.startAttempt("simple.xsb", 1)
	g.statistics.startAttempt("simple.xsb", 1)
	g.statistics.startAttempt("simple.xsb", 1)
	g.board.Move(Right)
	g.saveCompletedLevel()

	loaded, err := loadStatistics(g.statisticsPath)
	if err != nil {
		t.Fatal(err)
	}
	stats := loaded.level("simple.xsb", 1)
	if stats.Successes != 1 || len(stats.Results) != 1 {
		t.Fatalf("got successes=%d results=%d, want one success", stats.Successes, len(stats.Results))
	}
	result := stats.Results[0]
	if result.Attempt != 3 || result.Moves != 1 || result.Pushes != 1 {
		t.Errorf("got result %+v, want attempt 3 with 1 move and 1 push", result)
	}
	if !g.resultRecorded || g.progress["simple.xsb"] != 1 {
		t.Errorf("completion state: resultRecorded=%v progress=%d", g.resultRecorded, g.progress["simple.xsb"])
	}
}
