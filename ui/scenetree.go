package ui

import (
	"gtic"
	"image"
)

type SceneTree struct {
	scenes              map[ElementID]IElement
	currentScene        ElementID
	focusedId, activeId ElementID
}

func NewSceneTree() *SceneTree { return &SceneTree{} }
func (s *SceneTree) AddScene(v IElement) {
	if s.scenes == nil {
		s.scenes = make(map[ElementID]IElement)
	}
	id := v.Property("ID")
	s.scenes[id.Get().(ElementID)] = v
}
func (s *SceneTree) Change(a *gtic.API, id ElementID) {
	if s.currentScene == id {
		return
	}
	s.currentScene = id
	s.scenes[s.currentScene].Property("rect").Set(
		image.Rect(0, 0, a.Bounds().Width, a.Bounds().Height))
	s.scenes[s.currentScene].Init(a)
}
func (s *SceneTree) TIC(a *gtic.API) {
	s.scenes[s.currentScene].Update(a)
	s.scenes[s.currentScene].Draw(a)
}

func (t *SceneTree) Focus(id ElementID) {
	if t.focusedId == id {
		return
	}
	t.focusedId = id
}
func (t *SceneTree) Blur() {
	if t.focusedId == "" {
		return
	}
	t.focusedId = ""
}
