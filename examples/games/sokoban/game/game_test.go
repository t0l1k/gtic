package game

import (
	"image"
	"testing"
)

func TestPathBoxToPlansWalkingAndPushes(t *testing.T) {
	board := NewBoard("########\n#@ $  .#\n########")
	box := board.Dim.Idx(3, 1)

	path, ok := board.PathBoxTo(box, 5, 1)
	if !ok {
		t.Fatal("expected a push path")
	}
	for _, direction := range path {
		if !board.Move(direction) {
			t.Fatalf("planned move %v was rejected", direction)
		}
	}
	if board.Field.PlayerAt() != board.Dim.Idx(4, 1) ||
		!board.Field.IsBoxAt(board.Dim.Idx(5, 1)) ||
		board.Moves != 3 || board.Pushes != 2 {
		t.Fatalf("after plan: player=%d moves=%d pushes=%d", board.Field.PlayerAt(), board.Moves, board.Pushes)
	}
}

func TestMouseBoxSelectionAnimatesPushPlan(t *testing.T) {
	g := &Game{Board: NewBoard("########\n#@ $  .#\n########"), selectedBox: -1}
	box := g.Board.Dim.Idx(3, 1)
	g.handleBoardClick(3, 1)
	if g.selectedBox != box {
		t.Fatalf("selected box = %d, want %d", g.selectedBox, box)
	}
	g.handleBoardClick(5, 1)

	if len(g.mousePath) != 3 || g.Board.Moves != 0 {
		t.Fatalf("planned path length=%d moves=%d, want 3 and no immediate moves", len(g.mousePath), g.Board.Moves)
	}
	for len(g.mousePath) > 0 {
		for i := 0; i < mouseStepInterval; i++ {
			g.updateMouseAnimation()
		}
	}
	if g.Board.Moves != 3 || g.Board.Pushes != 2 || g.selectedBox != g.Board.Dim.Idx(5, 1) {
		t.Fatalf("animation result: moves=%d pushes=%d selected=%d", g.Board.Moves, g.Board.Pushes, g.selectedBox)
	}
	if g.boxTarget != nil {
		t.Fatal("destination highlight should clear after the animation")
	}
}

func TestBlockedMouseBoxDestinationStartsRedBlink(t *testing.T) {
	g := &Game{Board: NewBoard("#######\n#@#$ .#\n#######"), selectedBox: -1}
	g.handleBoardClick(3, 1)
	g.handleBoardClick(5, 1)
	if g.blockedTarget == nil || *g.blockedTarget != image.Pt(5, 1) {
		t.Fatalf("blocked target = %v, want destination (5,1)", g.blockedTarget)
	}
	if g.blockedTicks != blockedBlinkDuration {
		t.Errorf("blocked blink duration = %d, want %d", g.blockedTicks, blockedBlinkDuration)
	}
	if len(g.mousePath) != 0 || g.Board.Moves != 0 {
		t.Error("blocked destination should not start movement")
	}
}

func TestMouseBoxPathCanBeCancelled(t *testing.T) {
	g := &Game{Board: NewBoard("########\n#@ $  .#\n########"), selectedBox: -1}
	g.handleBoardClick(3, 1)
	g.handleBoardClick(5, 1)
	g.cancelMouseAnimation()

	for i := 0; i < mouseStepInterval*2; i++ {
		g.updateMouseAnimation()
	}
	if g.Board.Moves != 0 {
		t.Errorf("cancelled path moved player %d times", g.Board.Moves)
	}
}
