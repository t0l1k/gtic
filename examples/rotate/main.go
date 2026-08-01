package main

import (
	"etic"
	"etic/examples/dodge_the_creeps/assets"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func RotateSpriteTest(id etic.ElementID) *etic.Element {
	r := etic.NewElement(id)
	var (
		playerSprWalk1 etic.SpriteID = "player.walk.1"
		scale          float32       = 1
	)
	r.OnInit = func(t *etic.Console) {
		t.LoadSprite(playerSprWalk1, etic.LoadImage(assets.PlayerGrey_walk1_png))
	}
	r.OnUpdate = func(t *etic.Console) {}

	drawRotateSprite := func(t *etic.Console) {
		s := t.Sprite(playerSprWalk1).Bounds().Size()
		sw := float32(s.X) * scale
		sh := float32(s.Y) * scale
		x, y := (t.Width-sw)/2, (t.Height-sh)/2
		angle := float64(t.Frame() % 360)
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: etic.FlipNone, Rotate: etic.RotateAngle, Angle: float32(angle)})
	}

	r.OnDraw = func(t *etic.Console) {
		x, y := float32(50), float32(50)
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 0, Rotate: 0}) // no flip
		x += 150
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 1, Rotate: 0}) // flip hor
		x += 150
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 2, Rotate: 0}) // flip ver
		x += 150
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 3, Rotate: 0}) // flip both

		x, y = 50, 300
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 0, Rotate: 0}) // no angle
		x += 150
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 0, Rotate: 1}) // angle 90
		x += 150
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 0, Rotate: 2}) // angle 180
		x += 150
		t.Spr(playerSprWalk1, x, y, etic.SprOpt{Scale: scale, Flip: 0, Rotate: 3}) // angle 270

		drawRotateSprite(t)
	}
	return r
}

func RoteteDemo() *etic.Element {
	r := etic.NewElement("demo.rotate")
	r.Add(etic.PanelGrid("rotate.demo.grid", 10, color.Transparent, colornames.Slategray))
	r.Add(RotateSpriteTest("rotate.demo.sprite"))
	return r
}

func main() {
	etic.Load(func() *etic.Cart {
		cart := etic.NewCart("Rotate sprite", 650, 650)
		demo := RoteteDemo()
		cart.OnBoot = func(t *etic.Console) {
			demo.Init(t)
			log.Println("Cart booted")
		}
		cart.OnTic = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEscape) {
				t.Exit()
			}
			if t.BtnP(ebiten.KeyEnter) {
				t.Reset()
			}
			t.Cls(colornames.Teal)
			t.Print("Press Enter=Reset Console\nEscape=Quit", 0, 0, colornames.Yellow)
			demo.Update(t)
			demo.Draw(t)
		}
		return cart
	}()).Run()
}
