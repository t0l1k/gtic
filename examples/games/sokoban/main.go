package main

import (
	"gtic"
	"gtic/examples/games/sokoban/game"
	"log"
)

func main() {
	game, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	gtic.BOOT = func(a *gtic.API) {}
	gtic.TIC = func(a *gtic.API) {
		game.Input(a)
		game.Draw(a)
	}
	if err := gtic.Load(gtic.WithTitle("Sokoban"), gtic.WithMode(320, 200)).Run(); err != nil {
		log.Println(err)
	}
}
