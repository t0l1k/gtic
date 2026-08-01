package game

import (
	"etic"
	"fmt"
	"strconv"

	"golang.org/x/image/colornames"
)

type ScreenGame struct {
	*etic.Element
	fn    etic.SlotFn[*etic.Button]
	game  *Game
	board *etic.Element
	a, b  *etic.Label
}

func NewScreenGame(id etic.ElementID, fn etic.SlotFn[*etic.Button]) *ScreenGame {
	f := &ScreenGame{Element: etic.NewElement(id)}
	f.Layout().Set(etic.NewVerticalFlexLayout([]float32{0.1, 0.9}, 0.1))
	f.fn = fn
	f.game = NewGame()
	footer := etic.NewElement("footer")
	footer.Layout().Set(etic.NewVerticalBoxLayout(0.01))
	f.a = etic.NewLabel("f.s.label.title", "a")
	f.b = etic.NewLabel("f.s.label.helper", "b")
	footer.Add(f.a)
	footer.Add(f.b)
	f.Add(footer)
	f.board = etic.NewElement("game.board")
	f.Add(f.board)

	f.OnUpdate = func(c *etic.Console) {
		f.game.mismatchTimer.Update(c)
	}
	return f
}

func (s *ScreenGame) New(gd *GameData) {
	s.game.New(gd)
	s.board.ResetChildren()
	r, c := gd.dim.row, gd.dim.column
	s.board.Layout().Set(etic.NewSquareGriodLayout(r, c, 0.1))
	for i, card := range s.game.field {
		btn := etic.NewButton(etic.ElementID("btn."+strconv.Itoa(i)), card.String(), s.fn)
		s.board.Add(btn)

		card.state.OnChange.Connect(func(cs CellState) {
			switch cs {
			case CellClosed:
				btn.Text.Set(card.String())
				btn.Color.Set(colornames.Navy)
			case CellOpen:
				btn.Text.Set(card.String())
				btn.Color.Set(colornames.Red)
			case CellMatch:
				btn.Text.Set(card.String())
				btn.Color.Set(colornames.Green)
			}
		})
	}

	s.game.data.Changed.Connect(func(value string) {
		s.Data()
	})
	s.Data()
}

func (s *ScreenGame) Data() {
	s.a.Text.Set(fmt.Sprintf("Поле %d × %d", s.game.data.dim.row, s.game.data.dim.column))
	s.b.Text.Set(fmt.Sprintf("Нажатий:%v Ходы: %d", s.game.data.clicks, s.game.data.moves))
}
