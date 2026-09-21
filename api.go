package gtic

import "time"

type API struct {
	config        Config
	vram          *screen
	input         Input
	started, last time.Time
	delta         time.Duration
	booted, quit  bool
	frame         int

	sprites *spriteBank

	pallete         map[int]RGBA
	btnHoldCounters map[int]int

	showTitleFpsTps bool
}

func newAPI(config Config) *API {
	a := &API{
		config:          config,
		vram:            newScreen(config.Mode.Width, config.Mode.Height),
		started:         time.Now(),
		sprites:         newSpriteBank(),
		btnHoldCounters: make(map[int]int),
		showTitleFpsTps: true,
	}
	a.ResetPal()
	return a
}

func (c *API) Config() Config       { return c.config }
func (a *API) Bounds() Mode         { return Mode{a.vram.width, a.vram.height} }
func (a *API) Sprites() *spriteBank { return a.sprites }
