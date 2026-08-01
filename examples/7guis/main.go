package main

import (
	"etic"
	"etic/examples/7guis/counter"
	"etic/examples/7guis/data"
	"etic/examples/7guis/temper"
	"etic/examples/7guis/timer"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func sceneMain(fn etic.SlotFn[*etic.Button]) *etic.Element {
	s := etic.NewElement("scene.main")
	s.Layout().Set(etic.NewVerticalBoxLayout(0.05))
	s.Add(etic.TopBar("main.scene.topbar", data.Title, data.QuitApp, fn))
	for _, txt := range []string{data.Counter, data.TempConv, data.FlightBooker, data.Timer, data.Crud, data.Circle, data.Cells} {
		s.Add(etic.NewButton("scane.main.button", txt, fn))
	}
	return s
}

func main() {
	etic.Load(func() *etic.Cart {
		cart := etic.NewCart("Simple", 400, 400)
		sceneTree := &etic.SceneTree{}

		cart.OnBoot = func(t *etic.Console) {

			sceneCounter := counter.Counter("scene.counter", func(b *etic.Button) {
				switch b.Text.Get() {
				case data.QuitDemo:
					sceneTree.Change(t, "scene.main")
				}
			})

			sceneTemp := temper.TemeratureConv("scene.temperature.converter", func(b *etic.Button) {
				switch b.Text.Get() {
				case data.QuitDemo:
					sceneTree.Change(t, "scene.main")
				}
			})

			sceneTimer := timer.Timer("scene.timer", func(b *etic.Button) {
				switch b.Text.Get() {
				case data.QuitDemo:
					sceneTree.Change(t, "scene.main")
				}
			})

			sm := sceneMain(func(b *etic.Button) {
				switch b.Text.Get() {
				case data.QuitApp:
					t.Exit()
				case data.Counter:
					sceneTree.Change(t, "scene.counter")
				case data.TempConv:
					sceneTree.Change(t, "scene.temperature.converter")
				case data.Timer:
					sceneTree.Change(t, "scene.timer")
				}
			})
			sceneTree.AddScene(sm)

			sceneTree.AddScene(sceneCounter)
			sceneTree.AddScene(sceneTemp)
			sceneTree.AddScene(sceneTimer)
			sceneTree.Change(t, sm.ID().Get())

			log.Println("7GUIS cart booted")
		}
		cart.OnTic = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEscape) {
				t.Exit()
			}
			t.Cls(colornames.Teal)
			sceneTree.TIC(t)
		}
		return cart
	}()).Run()
}
