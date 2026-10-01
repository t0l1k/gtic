package ui

import (
	"gtic"
	"time"
)

type TimerState int

const (
	TimerStop TimerState = iota
	TimerRun
	TimerDone
)

// Timer id ElementID, waitTime time.Duration, autoStart, oneShot bool
func NewTimer(id ElementID, waitTime time.Duration, autoStart, oneShot bool) *Element {
	var (
		timeLeft time.Duration
		running  bool
	)
	start := func() {
		running = true
		timeLeft = waitTime
	}
	e := NewElement(id)
	done := e.RegisterProperty("done", false)
	state := e.RegisterProperty("state", TimerStop)
	state.OnChange.Connect(func(a any) {
		s := a.(TimerState)
		switch s {
		case TimerStop:
			running = false
		case TimerRun:
			start()
		case TimerDone:
			done.Set(true)
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
