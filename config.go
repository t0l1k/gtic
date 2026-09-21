package gtic

import (
	"fmt"
	"log"
)

type EticEngine int

const (
	EbitenEngine EticEngine = iota
	PngEngine
	ConsoleEngine
)

type Mode struct{ Width, Height int }

type Config struct {
	Engine EticEngine
	Title  string
	Mode   Mode
	Scale  int
}

type Option func(*Config) error

func WithEngine(e EticEngine) Option {
	return func(config *Config) error {
		eng := EbitenEngine
		switch e {
		case EbitenEngine:
			eng = EbitenEngine
		case PngEngine:
			eng = PngEngine
		case ConsoleEngine:
			eng = ConsoleEngine
		default:
			log.Fatal("unknown driver")
		}
		config.Engine = eng
		return nil
	}
}
func WithMode(w, h int) Option {
	mode := Mode{Width: w, Height: h}
	return func(config *Config) error {
		if mode.Width <= 0 || mode.Height <= 0 {
			return fmt.Errorf("etic: mode dimensions must be positive: %dx%d", mode.Width, mode.Height)
		}
		config.Mode = mode
		return nil
	}
}

func WithTitle(title string) Option {
	return func(config *Config) error {
		if title == "" {
			return fmt.Errorf("etic: cartridge title is required")
		}
		config.Title = title
		return nil
	}
}

func WithScale(scale int) Option {
	return func(config *Config) error {
		if scale <= 0 {
			return fmt.Errorf("etic/ebiten: scale must be positive: %d", scale)
		}
		config.Scale = scale
		return nil
	}
}
