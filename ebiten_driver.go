package gtic

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
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
	for eticKey, ebitenKey := range keyMap {
		if ebiten.IsKeyPressed(ebitenKey) {
			keys[eticKey] = true
		}
	}

	return Input{
		keys:  keys,
		mouse: mouseInput(),
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
	KeyA:            ebiten.KeyA,
	KeyB:            ebiten.KeyB,
	KeyC:            ebiten.KeyC,
	KeyD:            ebiten.KeyD,
	KeyE:            ebiten.KeyE,
	KeyF:            ebiten.KeyF,
	KeyG:            ebiten.KeyG,
	KeyH:            ebiten.KeyH,
	KeyI:            ebiten.KeyI,
	KeyJ:            ebiten.KeyJ,
	KeyK:            ebiten.KeyK,
	KeyL:            ebiten.KeyL,
	KeyM:            ebiten.KeyM,
	KeyN:            ebiten.KeyN,
	KeyO:            ebiten.KeyO,
	KeyP:            ebiten.KeyP,
	KeyQ:            ebiten.KeyQ,
	KeyR:            ebiten.KeyR,
	KeyS:            ebiten.KeyS,
	KeyT:            ebiten.KeyT,
	KeyU:            ebiten.KeyU,
	KeyV:            ebiten.KeyV,
	KeyW:            ebiten.KeyW,
	KeyX:            ebiten.KeyX,
	KeyY:            ebiten.KeyY,
	KeyZ:            ebiten.KeyZ,
	Key0:            ebiten.Key0,
	Key1:            ebiten.Key1,
	Key2:            ebiten.Key2,
	Key3:            ebiten.Key3,
	Key4:            ebiten.Key4,
	Key5:            ebiten.Key5,
	Key6:            ebiten.Key6,
	Key7:            ebiten.Key7,
	Key8:            ebiten.Key8,
	Key9:            ebiten.Key9,
	KeyMINUS:        ebiten.KeyMinus,
	KeyEQUALS:       ebiten.KeyEqual,
	KeyLEFTBRACKET:  ebiten.KeyBracketLeft,
	KeyRIGHTBRACKET: ebiten.KeyBracketRight,
	KeyBACKSLASH:    ebiten.KeyBackslash,
	KeySEMICOLON:    ebiten.KeySemicolon,
	KeyAPOSTROPHE:   ebiten.KeyQuote,
	KeyGRAVE:        ebiten.KeyBackquote,
	KeyCOMMA:        ebiten.KeyComma,
	KeyPERIOD:       ebiten.KeyPeriod,
	KeySLASH:        ebiten.KeySlash,
	KeyRETURN:       ebiten.KeyEnter,
	KeyBACKSPACE:    ebiten.KeyBackspace,
	KeyDELETE:       ebiten.KeyDelete,
	KeyINSERT:       ebiten.KeyInsert,
	KeyPAGEUP:       ebiten.KeyPageUp,
	KeyPAGEDOWN:     ebiten.KeyPageDown,
	KeyHOME:         ebiten.KeyHome,
	KeyEND:          ebiten.KeyEnd,
	KeyUP:           ebiten.KeyArrowUp,
	KeyDOWN:         ebiten.KeyArrowDown,
	KeyLEFT:         ebiten.KeyArrowLeft,
	KeyRIGHT:        ebiten.KeyArrowRight,
	KeyCAPSLOCK:     ebiten.KeyCapsLock,
	KeyCTRL:         ebiten.KeyControl,
	KeySHIFT:        ebiten.KeyShift,
	KeyALT:          ebiten.KeyAlt,
	KeyESC:          ebiten.KeyEscape,
	KeyF1:           ebiten.KeyF1,
	KeyF2:           ebiten.KeyF2,
	KeyF3:           ebiten.KeyF3,
	KeyF4:           ebiten.KeyF4,
	KeyF5:           ebiten.KeyF5,
	KeyF6:           ebiten.KeyF6,
	KeyF7:           ebiten.KeyF7,
	KeyF8:           ebiten.KeyF8,
	KeyF9:           ebiten.KeyF9,
	KeyF10:          ebiten.KeyF10,
	KeyF11:          ebiten.KeyF11,
	KeyF12:          ebiten.KeyF12,
	KeySPACE:        ebiten.KeySpace,
	KeyTAB:          ebiten.KeyTab,
	KeyMeta:         ebiten.KeyMeta,
	KeyMax:          ebiten.KeyMax,
}
