package ui

import "gtic"

type SceneTree struct {
	scenes       map[ElementID]IElement
	currentScene ElementID
}

func NewSceneTree() *SceneTree { return &SceneTree{} }
func (s *SceneTree) AddScene(v IElement) {
	if s.scenes == nil {
		s.scenes = make(map[ElementID]IElement)
	}
	id := v.Property("ID")
	s.scenes[id.Get().(ElementID)] = v
}
func (s *SceneTree) Change(t *gtic.API, id ElementID) {
	if s.currentScene == id {
		return
	}
	s.currentScene = id
	s.scenes[s.currentScene].Init(t)
}
func (s *SceneTree) TIC(t *gtic.API) {
	s.scenes[s.currentScene].Update(t)
	s.scenes[s.currentScene].Draw(t)
}
