package etic

type SceneTree struct {
	scenes       map[ElementID]IElement
	currentScene ElementID
}

func NewSceneTree() *SceneTree { return &SceneTree{} }
func (s *SceneTree) AddScene(v IElement) {
	if s.scenes == nil {
		s.scenes = make(map[ElementID]IElement)
	}
	id := v.ID().Get()
	s.scenes[id] = v
}
func (s *SceneTree) Change(t *Console, id ElementID) {
	s.currentScene = id
	s.scenes[s.currentScene].Bounds().Set(Rect(0, 0, t.Width, t.Height))
	s.scenes[s.currentScene].Init(t)
}
func (s *SceneTree) TIC(t *Console) {
	s.scenes[s.currentScene].Update(t)
	s.scenes[s.currentScene].Draw(t)
}
