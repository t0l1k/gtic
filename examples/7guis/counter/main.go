package counter

import (
	"etic"
	"etic/examples/7guis/data"
	"strconv"
)

func Counter(id etic.ElementID, fn etic.SlotFn[*etic.Button]) *etic.Element {
	s := etic.NewElement(id)
	counter := etic.NewProperty(0)
	lblCount := etic.NewLabel("scene.counter.label.count", "0")
	btnInc := etic.NewButton("scene.counter.button.inc", "Inc", func(b *etic.Button) {
		counter.Set(counter.Get() + 1)
	})
	btnDec := etic.NewButton("scene.counter.button.dec", "Dec", func(b *etic.Button) {
		counter.Set(counter.Get() - 1)
	})
	btnReset := etic.NewButton("scene.counter.button.reset", "Reset", func(b *etic.Button) {
		counter.Set(0)
	})
	counter.OnChange.Connect(func(i int) {
		lblCount.Text.Set(strconv.Itoa(i))
	})

	s.Layout().Set(etic.NewVerticalFlexLayout([]float32{0.1, 0.7, 0.2}, 0.05))
	s.Add(etic.TopBar("scene.counter.topbar", data.Counter, data.QuitDemo, fn))
	s.Add(lblCount)

	contBtns := etic.NewElement("scene.counter.buttons.container")
	contBtns.Layout().Set(etic.NewHorizontalBoxLayout(0.05))
	contBtns.Add(btnInc)
	contBtns.Add(btnDec)
	contBtns.Add(btnReset)
	s.Add(contBtns)
	return s
}
