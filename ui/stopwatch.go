package ui

import (
	"gtic"
	"time"
)

type StateStopwatch int

const (
	StopwatchStart StateStopwatch = iota
	StopwatchStop
	StopwatchReset
)

func NewStopwatch(id ElementID) *Element {
	var (
		running  bool
		started  time.Time
		duration time.Duration
	)
	e := NewElement(id)
	durationProp := e.RegisterProperty("duration", time.Duration(0))
	state := e.RegisterProperty("state", StopwatchStop)
	state.OnChange.Connect(func(v any) {
		s := v.(StateStopwatch)
		switch s {
		case StopwatchStart:
			started = time.Now()
			running = true
		case StopwatchStop:
			if running {
				duration += time.Since(started)
				running = false
			}
		case StopwatchReset:
			running = false
			duration = 0
		}
	})
	e.OnUpdate = func(a *gtic.API) {
		if !running {
			durationProp.Set(duration)
		} else {
			dur := time.Since(started)
			dur += duration
			durationProp.Set(dur)
		}
	}
	return e
}
