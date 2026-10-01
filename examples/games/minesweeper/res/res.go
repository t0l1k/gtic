package res

import (
	_ "embed"
	"gtic"
	"image"
	"strconv"
)

const (
	SprDigit1 = iota // 0..9 — большие цифры для счётчика
	SprDigit2
	SprDigit3
	SprDigit4
	SprDigit5
	SprDigit6
	SprDigit7
	SprDigit8
	SprDigit9
	SprDigit0
	SprFaceIdle   // нормальное лицо
	SprFacePlay   // во время игры
	SprFaceOuch   // "ой"
	SprFaceWin    // победа
	SprFaceDead   // проигрыш
	SprCellClosed // закрытая клетка
	SprCellOpen   // открытая (пустая)
	SprCellFlag   // флаг
	SprCellQ      // вопрос
	SprCellQOpen  // вопрос (нажата)
	SprCellMine   // мина
	SprCellBoom   // взорвавшаяся мина
	SprCellWrong  // флаг ошибочный
	SprNum1       // 1..8 — цифры на поле (SprNum1 + n - 1)
	SprNum2
	SprNum3
	SprNum4
	SprNum5
	SprNum6
	SprNum7
	SprNum8
)

var (
	//go:embed mines.png
	MinesSprites_png []byte
)

func RegisterSprites(a *gtic.API) []gtic.SpriteID {
	// rectsMineSheet — разметка листа (139×84) mines.png
	surface, err := gtic.DecodeImage(MinesSprites_png)
	if err != nil {
		panic(err)
	}
	var ids []gtic.SpriteID

	add := func(name string, r image.Rectangle) {
		id := a.Sprites().Register(gtic.NewSprite(name, surface, r))
		ids = append(ids, id)
	}
	for i := 0; i < 10; i++ { // ряд 0: большие цифры 13×23, шаг 14
		add("Digit:"+strconv.Itoa(i), image.Rect(i*14, 0, i*14+13, 23))
	}
	for i := 0; i < 5; i++ { // ряд 1: лица 26×26, шаг 27
		add("Face:"+strconv.Itoa(i), image.Rect(i*27, 24, i*27+26, 50))
	}
	for i := 0; i < 8; i++ { // ряд 2: клетки поля 16×16, шаг 17
		add("UpSquare:"+strconv.Itoa(i), image.Rect(i*17, 51, i*17+16, 67))
	}
	for i := 0; i < 8; i++ { // ряд 3: цифры 1–8, 16×16
		add("DownSquare:"+strconv.Itoa(i), image.Rect(i*17, 68, i*17+16, 84))
	}
	return ids
}
