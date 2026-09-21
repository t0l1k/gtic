package main

import (
	"gtic"
	"gtic/react"
	"gtic/ui"
	"image"
	"log"
	"math/rand/v2"
	"slices"
	"strconv"
)

var (
	w0, h0 = 400, 400
	grid   = 16
	scale  = min(w0, h0) / grid

	Up    = image.Pt(0, -1)
	Down  = image.Pt(0, 1)
	Left  = image.Pt(-1, 0)
	Right = image.Pt(1, 0)
)

func NewGame() *ui.Element {
	var (
		game, snake, food         *ui.Element
		foodSprID, headID, bodyID gtic.SpriteID
	)

	appleSpriteData := "00000000" + "00050500" + "00225220" + "02332332" + "03333333" + "02333332" + "00232320" + "00000000"
	headSprData := "00000000" + "00000000" + "06B00B60" + "65555556" + "55555555" + "55555555" + "65555556" + "06555560"
	bodySprData := "06555560" + "65555556" + "55555555" + "55555555" + "55555555" + "55555555" + "65555556" + "06555560"

	registerSnakeGameSprites := func(a *gtic.API) {

		sprData := func(data string) gtic.Surface[gtic.RGBA] {
			pixels := make([]gtic.RGBA, 0)
			for i := 0; i < 64; i++ {
				val, _ := strconv.ParseUint(string(data[i]), 16, 4)
				colIdx := int(val)
				pixels = append(pixels, a.Pal(colIdx))
			}
			return gtic.Surface[gtic.RGBA]{Data: pixels, Width: 8, Height: 8}
		}
		foodSprID = a.Sprites().Register(gtic.NewSprite("apple.sprite", sprData(appleSpriteData), image.Rect(0, 0, 8, 8)))
		headID = a.Sprites().Register(gtic.NewSprite("snake.head.sprite", sprData(headSprData), image.Rect(0, 0, 8, 8)))
		bodyID = a.Sprites().Register(gtic.NewSprite("snake.body.sprite", sprData(bodySprData), image.Rect(0, 0, 8, 8)))
	}

	setNewFood := func(snakeBody []image.Point) (x, y int) {
		check := func(x, y int) bool {
			for _, v := range snakeBody {
				if v.X == x && v.Y == y {
					return true
				}
			}
			return false
		}
		for {
			x, y = rand.IntN(grid), rand.IntN(grid)
			if !check(x, y) {
				return x, y
			}
		}
	}

	drawField := func(t *gtic.API) {
		for i := 0; i < t.Bounds().Height; i++ {
			t.Line(0, i*scale, t.Bounds().Width, i*scale, t.Pal(1))
			t.Line(i*scale, 0, i*scale, t.Bounds().Height, t.Pal(1))
		}
	}

	food = ui.NewElement("snake.game.food")
	foodPos := food.RegisterProperty("food.pos", image.Pt(8, 8))
	food.OnDraw = func(a *gtic.API) {
		verticalShift := a.Frame() % 10 / 5
		pos := foodPos.Get().(image.Point)
		x, y := pos.X*scale, pos.Y*scale+verticalShift
		a.Spr(foodSprID, x, y, a.Pal(0), scale/8, 0, 0)
	}

	snake = ui.NewElement("snake.game.snakebody")
	snakeDefaultPos := []image.Point{image.Pt(13, 8), image.Pt(14, 8), image.Pt(15, 8)}
	snakeBody := snake.RegisterUncomparableProperty("body", react.NewPropertyWithEqual[any]([]image.Point{}, func(a, b any) bool { return false }))
	snakeDir := snake.RegisterProperty("dir", Left)
	snakeSpeed := snake.RegisterProperty("speed", 10)
	snake.OnInit = func(a *gtic.API) {
		snakeBody.Set(nil)
		snakeBody.Set(snakeDefaultPos)
		snakeDir.Set(Left)
		snakeSpeed.Set(10)
	}
	snake.OnUpdate = func(a *gtic.API) {}
	snake.OnDraw = func(a *gtic.API) {
		for i, v := range snakeBody.Get().([]image.Point) {
			x, y := v.X*scale, v.Y*scale
			rotate := 0
			if i == 0 {
				switch snakeDir.Get().(image.Point) {
				case Up:
					rotate = 0
				case Down:
					rotate = 2
				case Left:
					rotate = 1
				case Right:
					rotate = 3
				}
				a.Spr(headID, x, y, a.Pal(0), scale/8, 0, rotate)
			} else {
				a.Spr(bodyID, x, y, a.Pal(0), scale/8, 0, 0)
			}
		}
	}

	game = ui.NewElement("snake.game")
	game.RegisterProperty("game.score", "0")
	game.RegisterProperty("game.gameover", false)
	game.Add(food)
	game.Add(snake)
	game.OnInit = func(a *gtic.API) {
		registerSnakeGameSprites(a)
		game.Property("game.score").Set("0")
		game.Property("game.gameover").Set(false)
		log.Println("game snake reset")
	}
	game.OnUpdate = func(a *gtic.API) {
		if game.Property("game.gameover").Get().(bool) {
			return
		}
		dir := snakeDir.Get().(image.Point)
		if a.KeyP(gtic.KeyLEFT) && dir.X == 0 {
			snakeDir.Set(Left)
		}
		if a.KeyP(gtic.KeyRIGHT) && dir.X == 0 {
			snakeDir.Set(Right)
		}
		if a.KeyP(gtic.KeyUP) && dir.Y == 0 {
			snakeDir.Set(Up)
		}
		if a.KeyP(gtic.KeyDOWN) && dir.Y == 0 {
			snakeDir.Set(Down)
		}

		if a.Frame()%snakeSpeed.Get().(int) == 0 {
			body := snakeBody.Get().([]image.Point)
			newPos := body[0].Add(snakeDir.Get().(image.Point))
			if !newPos.In(image.Rect(0, 0, grid, grid)) {
				game.Property("game.gameover").Set(true)
			}
			for i := 0; i < len(body); i++ {
				if body[i] == newPos {
					game.Property("game.gameover").Set(true)
				}
			}
			snakeBody.Set(slices.Insert(snakeBody.Get().([]image.Point), 0, newPos))
			pos := foodPos.Get().(image.Point)
			body = snakeBody.Get().([]image.Point)
			if newPos == pos {
				game.Property("game.score").Set(strconv.Itoa(len(snakeBody.Get().([]image.Point)) - 3))
				x, y := setNewFood(body)
				foodPos.Set(image.Pt(x, y))
			} else {
				snakeBody.Set(body[:len(body)-1])
			}
		}
	}
	game.OnDraw = func(a *gtic.API) {
		a.Cls()
		drawField(a)
		str := "Score: " + game.Property("game.score").Get().(string)
		if game.Property("game.gameover").Get().(bool) {
			str += "\nGame Over!!!\nPress Enter to play again."
		}
		a.Print(str, 0, 0, a.Pal(3), false, scale/10)
	}
	return game
}

func main() {
	game := NewGame()
	gtic.BOOT = func(a *gtic.API) { game.Init(a) }
	gtic.TIC = func(a *gtic.API) {
		switch {
		case a.KeyP(gtic.KeyESC):
			a.Exit()
		case a.KeyP(gtic.KeyRETURN):
			a.Reset()
		}
		game.Update(a)
		game.Draw(a)
	}

	if err := gtic.Load(gtic.WithTitle("Snake Game"), gtic.WithMode(w0, h0)).Run(); err != nil {
		log.Println(err)
	}
}
