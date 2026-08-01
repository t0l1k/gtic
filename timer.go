package etic

import (
	"time"
)

type Timer struct {
	Element
	waitTime  time.Duration
	timeLeft  time.Duration
	oneShot   bool
	autostart bool
	running   bool
	Timeout   Signal[*Timer]
}

func NewTimer(id ElementID) *Timer { return &Timer{Element: *NewElement(id)} }
func (t *Timer) Setup(w time.Duration, autostart, oneShot bool) *Timer {
	t.waitTime = w
	t.timeLeft = w
	t.oneShot = oneShot
	t.autostart = autostart
	if t.autostart {
		t.Start()
	}
	return t
}
func (t *Timer) Start() {
	t.timeLeft = t.waitTime
	t.running = true
}
func (t *Timer) Stop()           { t.running = false }
func (t *Timer) IsStopped() bool { return !t.running }
func (t *Timer) Update(tic *Console) {
	if !t.running {
		return
	}
	t.timeLeft -= tic.Tick()
	if t.timeLeft <= 0 {
		t.Timeout.Emit(t)
		if t.oneShot {
			t.running = false
		} else {
			t.timeLeft = t.waitTime
		}
	}
}
