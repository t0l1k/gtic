package etic

func (t *Console) Focus(id ElementID) {
	if t.focusedId == id {
		return
	}
	t.focusedId = id
}
func (t *Console) Blur() {
	if t.focusedId == "" {
		return
	}
	t.focusedId = ""
}
