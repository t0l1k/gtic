# Tic-80 простой прототип на golang ebitenengine
Вдохновил Tic-80 на написание этой библиотечки, на golang/ebitenengine.

## Оглавление
- [Tic-80 простой прототип на golang ebitenengine](#tic-80-простой-прототип-на-golang-ebitenengine)
  - [Оглавление](#оглавление)
- [Привет мир пример](#привет-мир-пример)
- [Runtime](#runtime)
- [Cartridge](#cartridge)
- [Drawing](#drawing)
- [Input](#input)
- [Sound](#sound)
- [System](#system)
# Привет мир пример

``` golang

func main() {
	etic.Load(func() *etic.Cart {
		cart := etic.NewCart("Simple",320,200)
		cart.OnBoot = func(t *etic.T) {
			log.Println("Simple cart booted")
		}
		cart.OnTic = func(t *etic.T) {
			if t.BtnP(ebiten.KeyEscape) {
				t.Exit()
			}
			if t.BtnP(ebiten.KeyEnter) {
				t.Reset()
			}
			t.Cls(colornames.Teal)
			t.Print("Hello, World!", 0, 0, colornames.Yellow)
		}
		return cart
	}()).Run()
}
```

# Runtime
``` golang

type Runtime struct {
	api  *Console
	cart Cartridge
}

func Load(cart Cartridge) *Runtime {
	game := &Runtime{
		cart: cart,
	}
	w, h := cart.Size().X, cart.Size().Y
	game.api = loadConsole(float32(w), float32(h))
	return game
}
func (g *Runtime) Update() error {
	if g.api.quit {
		return ebiten.Termination
	}
	if !g.api.booted {
		g.cart.BOOT(g.api)
		g.api.booted = true
	}
	g.api.update()
	g.cart.TIC(g.api)
	return nil
}
func (g *Runtime) Draw(screen *ebiten.Image)      { g.api.draw(screen) }
func (g *Runtime) Layout(inW, inH int) (int, int) { return int(g.api.Width), int(g.api.Height) }
func (g *Runtime) Run() error {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	ebiten.SetWindowSize(int(g.api.Width), int(g.api.Height))
	ebiten.SetWindowTitle(g.cart.Title())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	return ebiten.RunGame(g)
}

```
# Cartridge

```golang

type Cartridge interface {
	Title() string
	Size() Point[float32]
	BOOT(*Console)
	TIC(*Console)
}
type Cart struct {
	title  string
	size   Point[float32]
	OnBoot func(*Console)
	OnTic  func(*Console)
}

func NewCart(title string, w, h float32) *Cart {
	return &Cart{title: title,
		size: Pt(w, h),
	}
}
func (c *Cart) Title() string        { return c.title }
func (c *Cart) Size() Point[float32] { return c.size }
func (c *Cart) BOOT(t *Console) {
	if c.OnBoot != nil {
		c.OnBoot(t)
	}
}
func (c *Cart) TIC(t *Console) {
	if c.OnTic != nil {
		c.OnTic(t)
	}
}

```

# Drawing
- [Circ](#circ) [x]
- [CircB](#circb) [x]
- [Elli](#elli) [x]
- [ElliB](#ellib) [x]
- [Clip](#clip) [x]
- [Cls](#cls) [x]
- [Font](#font) [ ]
- [Line](#line) [x]
- [Map](#map) [ ] 
- [Pix](#pix) [x]
- [Print](#print) [x]
- [Rect](#rect) [x]
- [RectB](#rectb) [x]
- [Spr](#spr) [x]
- [Tri](#tri) [x]
- [TriB](#trib) [x] 

# Input
- [Btn](#btn) [x]
- [BtnP](#btnp) [x]
- [Key](#key) [ ]
- [KeyP](#keyp) [ ]
- [Mouse](#mouse) [x]

# Sound
- [Music](#music) [ ]
- [Sfx](#sfx) [ ]

# System
- [Exit](#exit) [x]
- [Reset](#reset) [x]
- [Time](#time) [x]
- [Tstamp](#tstamp) [ ]
- [Trace](#trace) [x] 

Реализовал несколько примеров использования.

simple  Привет Мир 
demo    Несколько демок нажав Enter следующая сцена по кругу
rotate  Показ использования спрайтов. 
7GUIs   Немного гуи пример счетчик/температуры конвертировать/таймер
Игры 
    Змейка  тут спрайт 8х8.
    Увернись от крипов
    Найди пару




