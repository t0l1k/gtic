package gtic

import (
	"errors"
	"time"
)

type Driver interface{ Run(e *Engine) error }

type Engine struct{ api *API }

func Load(options ...Option) *Engine {
	config := Config{Title: "tic", Mode: Mode{240, 136}, Scale: 3}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(&config); err != nil {
			panic(err)
		}
	}
	return &Engine{api: newAPI(config)}
}

func (r *Engine) Tick(input Input) error {
	r.api.frame++
	now := time.Now()
	if !r.api.last.IsZero() {
		r.api.delta = now.Sub(r.api.last)
	}
	r.api.last = now
	r.api.input = input
	if !r.api.booted {
		r.api.booted = true
		BOOT(r.api)
	}
	TIC(r.api)
	if r.api.quit {
		return errors.New("regular termination")
	}
	return nil
}

func (r *Engine) Run() error {
	var drv Driver
	switch r.api.config.Engine {
	case EbitenEngine:
		drv = newEbitenDriver()
	case PngEngine:
		drv = &pngDriver{Path: "frame.png"}
	case ConsoleEngine:
	}
	return drv.Run(r)
}
