package etic

func TopBar(id ElementID, title, quitLbl string, fn SlotFn[*Button]) *Element {
	t := NewElement(id)
	t.Layout().Set(NewHorizontalFlexLayout([]float32{0.1, 0.9}, 0.05))
	t.Add(NewButton("topbar.button.quit", quitLbl, fn))
	t.Add(NewLabel("topbar.label.title", title))
	return t
}
