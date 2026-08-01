package etic

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Runtime struct {
	api  *Console
	cart Cartridge
}

func Load(cart Cartridge) *Runtime {
	game := &Runtime{
		cart: cart,
	}
	w, h := cart.Size().X, cart.Size().Y
	game.api = loadConsole(float32(w), float32(h))
	return game
}
func (g *Runtime) Update() error {
	if g.api.quit {
		return ebiten.Termination
	}
	if !g.api.booted {
		g.cart.BOOT(g.api)
		g.api.booted = true
	}
	g.api.update()
	g.cart.TIC(g.api)
	return nil
}
func (g *Runtime) Draw(screen *ebiten.Image)      { g.api.draw(screen) }
func (g *Runtime) Layout(inW, inH int) (int, int) { return int(g.api.Width), int(g.api.Height) }
func (g *Runtime) Run() error {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	ebiten.SetWindowSize(int(g.api.Width), int(g.api.Height))
	ebiten.SetWindowTitle(g.cart.Title())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	return ebiten.RunGame(g)
}
