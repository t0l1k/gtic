package main

import (
	_ "embed"
	"fmt"
	"gtic"
	"gtic/examples/games/minesweeper/game"
	"gtic/examples/games/minesweeper/res"
	"gtic/react"
	"gtic/ui"
	"image"
	"time"
)

func NewMinesGame() *ui.Element {
	var (
		sprIds                    []gtic.SpriteID
		ox, oy                    int
		leftPressed, rightPressed bool
		btnReset, gameStopwatch   *ui.Element
		faceSprite                *react.Var[int]
	)
	minesGame := game.NewMines(9, 9, 10)
	minesGame.State.OnChange.Connect(func(gs game.GameState) {
		swState := gameStopwatch.Property("state")
		switch gs {
		case game.GameStart:
			swState.Set(ui.StopwatchReset)
			faceSprite.Set(res.SprFaceIdle)
		case game.GamePlay:
			swState.Set(ui.StopwatchStart)
			faceSprite.Set(res.SprFaceIdle)
		case game.GameWin, game.GameGameOver:
			swState.Set(ui.StopwatchStop)
			switch gs {
			case game.GameWin:
				faceSprite.Set(res.SprFaceWin)
			case game.GameGameOver:
				faceSprite.Set(res.SprFaceDead)
			}
		}
	})

	gameStopwatch = ui.NewStopwatch("game.stopwatch")
	minesLeftIcon := func() *ui.Element {
		e := ui.NewElement("topbar.minesleft.icon")
		rect := e.Property("rect")
		e.OnInit = func(a *gtic.API) {}
		e.OnUpdate = func(a *gtic.API) {}
		e.OnDraw = func(a *gtic.API) {
			x, y, w, h := gtic.RectXYWH(rect.Get().(image.Rectangle))
			a.Rect(x, y, w, h, a.Pal(11))
			_, _, mines := minesGame.Dim()
			left := max(mines-minesGame.Flags(), 0)
			aa := left / 100 % 10
			bb := left / 10 % 10
			cc := left % 10

			result := func(v int) int {
				switch v {
				case 0:
					return 9
				default:
					return v - 1
				}
			}
			x += 1
			y += 1
			a.Spr(sprIds[result(aa)], x, y)
			a.Spr(sprIds[result(bb)], x+13, y)
			a.Spr(sprIds[result(cc)], x+13*2, y)
		}
		return e
	}()

	timerIcon := func() *ui.Element {
		e := ui.NewElement("topbar.timer.icon")
		rect := e.Property("rect")
		e.OnInit = func(a *gtic.API) {}
		e.OnUpdate = func(a *gtic.API) {}
		e.OnDraw = func(a *gtic.API) {
			x, y, w, h := gtic.RectXYWH(rect.Get().(image.Rectangle))
			a.Rect(x, y, w, h, a.Pal(11))
			left := int(gameStopwatch.Property("duration").Get().(time.Duration).Seconds())
			aa := left / 100 % 10
			bb := left / 10 % 10
			cc := left % 10

			result := func(v int) int {
				switch v {
				case 0:
					return 9
				default:
					return v - 1
				}
			}
			x += 1
			y += 1
			a.Spr(sprIds[result(aa)], x, y)
			a.Spr(sprIds[result(bb)], x+13, y)
			a.Spr(sprIds[result(cc)], x+13*2, y)
		}
		return e
	}()

	faceIcon := func() *ui.Element {
		e := ui.NewElement("topbar.face.icon")
		faceSprite = react.NewVar(res.SprFaceIdle)
		timerOuch := ui.NewTimer("timer.ouch", 250*time.Millisecond, false, true, func(s string) {
			if s != "done" || minesGame.State.Get() != game.GamePlay {
				return
			}
			faceSprite.Set(res.SprFaceIdle)
		})
		rect := e.Property("rect")
		var pressed bool
		e.OnInit = func(a *gtic.API) {}
		e.OnUpdate = func(a *gtic.API) {
			if minesGame.State.Get() != game.GamePlay {
				return
			}
			if leftPressed && !pressed {
				faceSprite.Set(res.SprFacePlay)
				pressed = true
			}
			if pressed && !leftPressed {
				timerOuch.Property("state").Set(ui.TimerStart)
				faceSprite.Set(res.SprFaceOuch)
				pressed = false
			}
		}
		e.OnDraw = func(a *gtic.API) {
			x, y, w, h := gtic.RectXYWH(rect.Get().(image.Rectangle))
			a.Rect(x, y, w, h, a.Pal(11))
			x += 1
			y += 1
			a.Spr(sprIds[faceSprite.Get()], x, y)
		}
		e.Add(timerOuch)
		return e
	}()

	topBar := func() *ui.Element { // счётчик мин + лицо
		e := ui.NewElement("topbar")
		rect := e.Property("rect")
		e.Add(minesLeftIcon)
		e.Add(faceIcon)
		e.Add(timerIcon)
		e.OnInit = func(a *gtic.API) {
			r := rect.Get().(image.Rectangle)
			x0, y0, w0, _ := gtic.RectXYWH(r)
			x, y, w, h := x0+1, y0+1, 39+2, 23+2
			minesLeftIcon.Property("rect").Set(image.Rect(x, y, x+w, y+h))
			x, y, w, h = x0+1+(w0-26)/2, y0+1, 26+2, 26+2
			faceIcon.Property("rect").Set(image.Rect(x, y, x+w, y+h))
			x, y, w, h = x0+w0-42, y0+1, 39+2, 23+2
			timerIcon.Property("rect").Set(image.Rect(x, y, x+w, y+h))
		}
		e.OnUpdate = func(a *gtic.API) {}
		e.OnDraw = func(a *gtic.API) {
			x, y, w, h := gtic.RectXYWH(rect.Get().(image.Rectangle))
			a.RectB(x, y, w, h, a.Pal(3))
		}
		return e
	}()

	cursor := func() *ui.Element {
		e := ui.NewElement("game.cursor")
		rect := e.Property("rect")
		pos := e.RegisterProperty("pos", image.Point{})
		e.OnInit = func(a *gtic.API) {}
		e.OnUpdate = func(a *gtic.API) {}
		e.OnDraw = func(a *gtic.API) {
			p := pos.Get().(image.Point)
			x, y, w, h := gtic.RectXYWH(rect.Get().(image.Rectangle))
			a.RectB(x+p.X*w, y+p.Y*h, w, h, a.Pal(2))
		}
		return e
	}()
	curPos := cursor.Property("pos")

	gameField := func() *ui.Element {
		e := ui.NewElement("game.field")
		e.OnInit = func(a *gtic.API) {}
		e.OnUpdate = func(a *gtic.API) {}
		e.OnDraw = func(a *gtic.API) {
			row, column, _ := minesGame.Dim()
			cur := curPos.Get().(image.Point)
			// --- поле ---
			for y := 0; y < column; y++ {
				for x := 0; x < row; x++ {
					c := minesGame.Cells()[y*row+x]
					px, py := ox+x*16, oy+y*16
					switch {
					case c.Wrong:
						a.Spr(sprIds[res.SprCellWrong], px, py)
					case c.Open && c.Mine:
						if minesGame.State.Get() == game.GameGameOver && x == cur.X && y == cur.Y {
							a.Spr(sprIds[res.SprCellBoom], px, py)
						} else {
							a.Spr(sprIds[res.SprCellMine], px, py)
						}
					case c.Open:
						a.Spr(sprIds[res.SprCellOpen], px, py)
						if n := minesGame.Count(x, y); n > 0 {
							a.Spr(sprIds[res.SprNum1+n-1], px, py)
						}
					case c.Flag:
						a.Spr(sprIds[res.SprCellFlag], px, py)
					default:
						a.Spr(sprIds[res.SprCellClosed], px, py)
					}
				}
			}
		}
		return e
	}()

	btnReset = ui.NewButton("btn.reset", "New Game", func(e *ui.Element) {
		minesGame.Reset()
		r, c, _ := minesGame.Dim()
		curPos.Set(image.Pt(r/2, c/2))
	})

	mainScene := ui.NewElement("minesweeper.game")
	mainScene.Add(btnReset)
	mainScene.Add(topBar)
	mainScene.Add(gameField)
	mainScene.Add(gameStopwatch)
	mainScene.Add(cursor)

	mainScene.OnInit = func(a *gtic.API) {
		sprIds = res.RegisterSprites(a)
		ox, oy = (a.Bounds().Width-9*16)/2, 34 // поле по центру, под панелью
		btnReset.Property("rect").Set(image.Rect(0, a.Bounds().Height-15, 100, a.Bounds().Height))
		cursor.Property("rect").Set(image.Rect(ox, oy, ox+16, oy+16))
		topBar.Property("rect").Set(image.Rect(0, 0, a.Bounds().Width, 30))
	}
	const cell = 16

	mainScene.OnUpdate = func(a *gtic.API) {
		mx, my, lp, _, rp, _, _ := a.Mouse()
		gx, gy := (mx-ox)/cell, (my-oy)/cell
		r, c, _ := minesGame.Dim()
		inside := gx >= 0 && gy >= 0 && gx < r && gy < c
		if inside {
			curPos.Set(image.Pt(gx, gy))
		}

		// нажатие по кромке (edge), а не удержание
		if lp && !leftPressed && inside {
			minesGame.Open(gx, gy)
		}
		if rp && !rightPressed && inside {
			minesGame.Flag(gx, gy)
		}
		leftPressed, rightPressed = lp, rp

		p := curPos.Get().(image.Point)
		cur := image.Pt(p.X, p.Y)
		switch {
		case a.KeyP(gtic.KeyR):
			minesGame.Reset()
			curPos.Set(image.Pt(r/2, c/2))
		case a.KeyP(gtic.KeyLEFT):
			cur.X = max(cur.X-1, 0)
		case a.KeyP(gtic.KeyRIGHT):
			cur.X = min(cur.X+1, r-1)
		case a.KeyP(gtic.KeyUP):
			cur.Y = max(cur.Y-1, 0)
		case a.KeyP(gtic.KeyDOWN):
			cur.Y = min(cur.Y+1, c-1)
		case a.KeyP(gtic.KeyZ):
			minesGame.Open(cur.X, cur.Y)
		case a.KeyP(gtic.KeyX):
			minesGame.Flag(cur.X, cur.Y)
		}
		if !cur.Eq(p) {
			curPos.Set(cur)
		}
	}
	mainScene.OnDraw = func(a *gtic.API) {}
	return mainScene
}

func main() {
	game := NewMinesGame()
	gtic.BOOT = func(a *gtic.API) { game.Init(a) }
	gtic.TIC = func(a *gtic.API) {
		if a.Key(gtic.KeyESC) {
			a.Exit()
		}
		if a.KeyP(gtic.KeyRETURN) {
			a.Reset()
		}
		game.Update(a)
		a.Cls()
		game.Draw(a)
	}
	if err := gtic.Load(gtic.WithTitle("Minesweeper"), gtic.WithMode(320, 200)).Run(); err != nil {
		fmt.Println(err)
	}
}
