package game

import (
	"etic"
	"fmt"
)

type ScreenSelect struct {
	*etic.Element
	levels *etic.Element
}

func NewScreenSelect(id etic.ElementID, gamesData GamesData, fn etic.SlotFn[*etic.Button]) *ScreenSelect {
	f := &ScreenSelect{Element: etic.NewElement(id)}
	f.Layout().Set(etic.NewVerticalFlexLayout([]float32{0.1, 0.9}, 0.01))
	footer := etic.NewElement("footer")
	footer.Layout().Set(etic.NewVerticalBoxLayout(0.01))
	footer.Add(etic.NewLabel("f.s.label.title", Title))
	footer.Add(etic.NewLabel("f.s.label.helper", "Выберите поле. Открываются только пройденные уровни."))
	f.Add(footer)
	f.levels = etic.NewElement("levels")
	f.levels.Layout().Set(etic.NewGridLayout(5, 5, 0.1))
	for _, v := range gamesData.Levels() {
		res := fmt.Sprintf("%vx%v", v.row, v.column)
		btn := etic.NewButton(etic.ElementID("dim."+res), res, fn)
		btn.State.Set(etic.ElementDisabled)
		f.levels.Add(btn)
	}
	f.levels.Children()[0].(*etic.Button).State.Set(etic.ElementNormal)
	f.Add(f.levels)
	return f
}
