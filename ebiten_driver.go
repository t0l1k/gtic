package gtic

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ebitenDriver struct {
	engine *Engine
	frame  *ebiten.Image
}

func newEbitenDriver() *ebitenDriver { return &ebitenDriver{} }

func (r *ebitenDriver) Update() error             { return r.engine.Tick(input(r.engine.api)) }
func (r *ebitenDriver) Draw(screen *ebiten.Image) { screen.WritePixels(r.engine.api.vram.blit()) }
func (r *ebitenDriver) Layout(int, int) (int, int) {
	return r.engine.api.config.Mode.Width, r.engine.api.config.Mode.Height
}
func (r *ebitenDriver) Run(e *Engine) error {
	r.engine = e
	conf := r.engine.api.Config()
	w, h := conf.Mode.Width, conf.Mode.Height
	r.frame = ebiten.NewImage(w, h)
	ebiten.SetWindowTitle(conf.Title)
	ebiten.SetWindowSize(w*conf.Scale, h*conf.Scale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	return ebiten.RunGame(r)
}

func input(a *API) Input {
	if a.showTitleFpsTps {
		title := fmt.Sprintf("%v [TPS %.1f / FPS %.1f]", a.config.Title, ebiten.ActualTPS(), ebiten.ActualFPS())
		ebiten.SetWindowTitle(title)
	}
	keys := make(map[Key]bool, len(keyMap))
	justPressed := make(map[Key]bool)
	for eticKey, ebitenKey := range keyMap {
		if ebiten.IsKeyPressed(ebitenKey) {
			keys[eticKey] = true
		}
		if inpututil.IsKeyJustPressed(ebitenKey) {
			justPressed[eticKey] = true
		}
	}

	return Input{
		keys:        keys,
		justPressed: justPressed,
		mouse:       mouseInput(),
	}
}

func mouseInput() mouse {
	x, y := ebiten.CursorPosition()
	sX, sY := ebiten.Wheel()
	buttons := make(map[MouseButton]bool)
	for button := range ebiten.MouseButtonMax {
		buttons[MouseButton(button)] = ebiten.IsMouseButtonPressed(button)
	}
	return newMouseData(x, y, sX, sY, buttons)
}

var keyMap = map[Key]ebiten.Key{
	KeyA: ebiten.KeyA, KeyB: ebiten.KeyB, KeyC: ebiten.KeyC,
	KeyD: ebiten.KeyD, KeyE: ebiten.KeyE, KeyF: ebiten.KeyF,
	KeyG: ebiten.KeyG, KeyH: ebiten.KeyH, KeyI: ebiten.KeyI,
	KeyJ: ebiten.KeyJ, KeyK: ebiten.KeyK, KeyL: ebiten.KeyL,
	KeyM: ebiten.KeyM, KeyN: ebiten.KeyN, KeyO: ebiten.KeyO,
	KeyP: ebiten.KeyP, KeyQ: ebiten.KeyQ, KeyR: ebiten.KeyR,
	KeyS: ebiten.KeyS, KeyT: ebiten.KeyT, KeyU: ebiten.KeyU,
	KeyV: ebiten.KeyV, KeyW: ebiten.KeyW, KeyX: ebiten.KeyX,
	KeyY: ebiten.KeyY, KeyZ: ebiten.KeyZ,
	Key0: ebiten.Key0, Key1: ebiten.Key1, Key2: ebiten.Key2,
	Key3: ebiten.Key3, Key4: ebiten.Key4, Key5: ebiten.Key5,
	Key6: ebiten.Key6, Key7: ebiten.Key7, Key8: ebiten.Key8,
	Key9:  ebiten.Key9,
	KeyUP: ebiten.KeyArrowUp, KeyDOWN: ebiten.KeyArrowDown,
	KeyLEFT: ebiten.KeyArrowLeft, KeyRIGHT: ebiten.KeyArrowRight,
	KeySPACE: ebiten.KeySpace, KeyRETURN: ebiten.KeyEnter,
	KeyESC: ebiten.KeyEscape,
}
