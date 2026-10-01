package ui

import "gtic/react"

func TopBar(id ElementID, title, quitLbl string, fn react.SlotFn[*Element]) *Element {
	t := NewElement(id)
	t.Property("layout").Set(NewHorizontalFlexLayout([]float32{0.1, 0.9}, 0.01))
	t.Add(NewButton("topbar.button.quit", quitLbl, fn))
	t.Add(NewLabel("topbar.label.title", title))
	return t
}
