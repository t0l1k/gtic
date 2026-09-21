## Оглавление
- [Tic-80 простой прототип на golang ebitenengine](#tic-80-простой-прототип-на-golang-ebitenengine)
- [Привет мир пример](#привет-мир-пример)
- [Drawing](#drawing)
- [Input](#input)
- [Sound](#sound)
- [System](#system)


# Tic-80 простой прототип на golang ebitenengine
Вдохновил Tic-80 на написание этой библиотечки, на ebitenengine/png. Это учебный прототип, обертка над ebiten, api tic80. Нет ограничений как в tic80, потому остался только контракт BOOT/TIC. Еще подглядывал на [pi](https://github.com/elgopher/pi), где на ebiten обернули pico8.
Рисуется на холст []RGBA, а после показывается через ebiten или png, есть палитра на 16 цветов, но можно показать любой цвет RGBA. Ввод от мыши/клавиатуры. Звуки/музыка еще в разработке. 

Реализовал несколько примеров.

hello  Привет Мир 
api_test    Несколько демок нажав Enter следующая сцена по кругу(большинство примеров конвертировал из вики tic80)
Змейка
Сапер





# Привет мир пример

``` golang
func main() {
	gtic.BOOT = func(a *gtic.API) { log.Println("Booted hello example") }
	gtic.TIC = func(a *gtic.API) {
		if a.Key(gtic.KeyESC) {
			a.Exit()
		}
		if a.KeyP(gtic.KeyRETURN) {
			a.Reset()
		}
		a.Cls()
		a.Print("Hello World!", 10, 10, a.Pal(4), false, 2)
	}
	if err := gtic.Load(gtic.WithTitle("Hello")).Run(); err != nil {
		log.Println(err)
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
- [Btn](#btn) [ ]
- [BtnP](#btnp) [x]
- [Key](#key) [x]
- [KeyP](#keyp) [x]
- [Mouse](#mouse) [x]

# Sound
- [Music](#music) [ ]
- [Sfx](#sfx) [ ]

# System
- [Exit](#exit) [x]
- [Reset](#reset) [x]
- [Time](#time) [x]
- [Tstamp](#tstamp) [x]
- [Trace](#trace) [ ] 
