package gtic

import "time"

func (a *API) Exit()                { a.quit = true }
func (a *API) Reset()               { a.booted = false }
func (a *API) Time() int64          { return time.Since(a.started).Milliseconds() }
func (a *API) Tstamp() int64        { return time.Now().Unix() }
func (a *API) Delta() time.Duration { return a.delta }
func (a *API) Frame() int           { return a.frame }
