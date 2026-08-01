package main

import (
	"etic"
	"etic/examples/dodge_the_creeps/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	etic.Load(func() *etic.Cart {
		c := etic.NewCart(game.Title, 480, 720)
		g := game.New()
		c.OnBoot = func(t *etic.Console) {
			game.LoadResources(t)
			g.ClearReady()
		}
		c.OnTic = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEscape) {
				t.Exit()
			}
			g.TIC(t)
		}
		return c
	}()).Run()
}
