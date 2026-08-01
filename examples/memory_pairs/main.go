package main

import (
	"etic"
	"etic/examples/memory_pairs/game"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func main() {
	etic.Load(func() *etic.Cart {
		cart := etic.NewCart(game.Title, 720, 480)
		sceneTree := &etic.SceneTree{}
		g := game.NewMainScreen("scene.main.memory.pairs.game")
		sceneTree.AddScene(g)

		cart.OnBoot = func(t *etic.Console) {
			sceneTree.Change(t, g.ID().Get())
		}
		cart.OnTic = func(t *etic.Console) {
			switch {
			case t.BtnP(ebiten.KeyEscape):
				t.Exit()
			case t.BtnP(ebiten.KeyEnter):
				t.Reset()
			}
			t.Cls(colornames.Teal)
			sceneTree.TIC(t)
		}
		return cart
	}()).Run()
}
