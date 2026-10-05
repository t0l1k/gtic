package gtic

type Key uint8

const (
	KeyUnknown Key = iota
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9

	KeyMINUS
	KeyEQUALS
	KeyLEFTBRACKET
	KeyRIGHTBRACKET
	KeyBACKSLASH
	KeySEMICOLON
	KeyAPOSTROPHE
	KeyGRAVE
	KeyCOMMA
	KeyPERIOD
	KeySLASH

	KeySPACE
	KeyTAB

	KeyRETURN
	KeyBACKSPACE
	KeyDELETE
	KeyINSERT

	KeyPAGEUP
	KeyPAGEDOWN
	KeyHOME
	KeyEND
	KeyUP
	KeyDOWN
	KeyLEFT
	KeyRIGHT

	KeyCAPSLOCK
	KeyCTRL
	KeySHIFT
	KeyALT
	KeyESC
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyNUMPAD0
	KeyNUMPAD1
	KeyNUMPAD2
	KeyNUMPAD3
	KeyNUMPAD4
	KeyNUMPAD5
	KeyNUMPAD6
	KeyNUMPAD7
	KeyNUMPAD8
	KeyNUMPAD9
	KeyNUMPADPLUS
	KeyNUMPADMINUS
	KeyNUMPADMULTIPLY
	KeyNUMPADDIVIDE
	KeyNUMPADENTER
	KeyNUMPADPERIOD
	KeyMeta
	KeyMax
)

type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonMiddle
	MouseButtonRight
)

type mouse struct {
	X, Y           int
	WheelX, WheelY float64
	Buttons        map[MouseButton]bool
}

func newMouseData(x, y int, wx, wy float64, buttons map[MouseButton]bool) mouse {
	return mouse{X: x, Y: y, WheelX: wx, WheelY: wy, Buttons: buttons}
}

type Input struct {
	mouse mouse
	keys  map[Key]bool
}

func (a *API) Mouse() (x, y int, left, middle, right bool, scrollX, scrollY float64) {
	m := a.input.mouse
	x, y = m.X, m.Y
	left, middle, right = m.Buttons[MouseButtonLeft], m.Buttons[MouseButtonMiddle], m.Buttons[MouseButtonRight]
	scrollX, scrollY = m.WheelX, m.WheelY
	return x, y, left, middle, right, scrollX, scrollY
}

func (a *API) Key(key ...Key) bool {
	if len(key) == 0 {
		for _, pressed := range a.input.keys {
			if pressed {
				return true
			}
		}
		return false
	}
	return a.input.keys[key[0]]
}

func (a *API) KeyP(key Key, param ...int) bool {
	hold, period := 30, 6
	if len(param) > 0 {
		hold, period = param[0], param[1]
	}

	pressed := a.input.keys[key]
	if !pressed {
		a.keyHoldCounters[key] = 0
		return false
	}
	ticks := a.keyHoldCounters[key]
	a.keyHoldCounters[key]++
	if ticks == 0 {
		return true
	}
	if ticks >= hold && (ticks-hold)%period == 0 {
		return true
	}
	return false
}

// Задержка автоповтора кнопок (btnp hold, period)
func (a *API) Btnp(id int, param ...int) bool {
	hold, period := 30, 6
	if len(param) > 0 {
		hold, period = param[0], param[1]
	}
	var pressed bool
	switch id {
	case 0:
		pressed = a.Key(KeyUP)
	case 1:
		pressed = a.Key(KeyDOWN)
	case 2:
		pressed = a.Key(KeyLEFT)
	case 3:
		pressed = a.Key(KeyRIGHT)
	case 4:
		pressed = a.Key(KeyZ) //A
	case 5:
		pressed = a.Key(KeyX) //B
	case 6:
		pressed = a.Key(KeyA) //X
	case 7:
		pressed = a.Key(KeyS) //Y
	}
	if !pressed {
		a.btnHoldCounters[id] = 0
		return false
	}
	ticks := a.btnHoldCounters[id]
	a.btnHoldCounters[id]++
	if ticks == 0 {
		return true
	}
	if ticks >= hold && (ticks-hold)%period == 0 {
		return true
	}
	return false
}
