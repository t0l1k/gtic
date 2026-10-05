package ui

import (
	"gtic"
	"gtic/react"
	"time"
)

type TimerState int

const (
	TimerStart TimerState = iota
	TimerStop
	TimerDone
)

func NewTimer(id ElementID, waitTime time.Duration, autoStart, oneShot bool, fn react.SlotFn[string]) *Element {
	var (
		timeLeft time.Duration
		running  bool
	)
	start := func() {
		running = true
		timeLeft = waitTime
	}
	e := NewElement(id)
	state := e.RegisterUncomparableProperty("state", react.NewPropertyWithEqual[any](TimerStop, func(a, b any) bool { return false }))
	state.OnChange.Connect(func(a any) {
		s := a.(TimerState)
		switch s {
		case TimerStart:
			start()
		case TimerStop:
			running = false
		case TimerDone:
			fn("done")
		}
	})
	e.OnInit = func(a *gtic.API) {
		timeLeft = waitTime
		if autoStart {
			start()
		}
	}
	e.OnUpdate = func(a *gtic.API) {
		if !running {
			return
		}
		timeLeft -= a.Delta()
		if timeLeft <= 0 {
			state.Set(TimerDone)
			if oneShot {
				running = false
			} else {
				timeLeft = waitTime
			}
		}
	}
	return e
}
