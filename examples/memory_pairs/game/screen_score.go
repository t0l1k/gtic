package game

import (
	"etic"
	"fmt"
	"strconv"
)

const (
	btnResetGame = "Повторить"
	btnNewGame   = "Новый"
	btnGoMenu    = "Меню"
)

type ScreenScore struct {
	*etic.Element
	title, result *etic.Label
	top           []*etic.Label
}

func NewScreenScore(id etic.ElementID, fn etic.SlotFn[*etic.Button]) *ScreenScore {
	f := &ScreenScore{Element: etic.NewElement(id)}
	f.Layout().Set(etic.NewVerticalFlexLayout([]float32{0.1, 0.8, 0.1}, 0.1))
	topbar := etic.NewElement("topbar")
	topbar.Layout().Set(etic.NewVerticalBoxLayout(0.01))
	f.title = etic.NewLabel("f.s.label.title", "Топ-10")
	topbar.Add(f.title)
	f.result = etic.NewLabel("f.s.label.helper", "f.gameresult")
	topbar.Add(f.result)
	f.Add(topbar)
	contLbls := etic.NewElement("top10.labels")
	contLbls.Layout().Set(etic.NewVerticalBoxLayout(0.01))
	for i := 0; i < 10; i++ {
		lbl := etic.NewLabel(etic.ElementID("top10.label."+strconv.Itoa(i)), "")
		f.top = append(f.top, lbl)
		contLbls.Add(lbl)
	}
	f.Add(contLbls)
	btns := etic.NewElement("buttons")
	btns.Layout().Set(etic.NewHorizontalBoxLayout(0.01))
	labels := []string{btnResetGame, btnNewGame, btnGoMenu}
	for _, v := range labels {
		btns.Add(etic.NewButton(etic.ElementID("btn "+v), v, fn))
	}
	f.Add(btns)
	return f
}

func (s *ScreenScore) SetLastGameScore(d GameData) {
	s.result.Text.Set(fmt.Sprintf("Результат: поле %d, %d очков, %d ходов, %d нажатий %v", d.dim, d.score, d.moves, d.clicks, d.duration))
}

func (s *ScreenScore) SetTop10(d []string) {
	for i, value := range d {
		if i < len(s.top)-1 {
			s.top[i].Text.Set(value)
		}
	}
}
