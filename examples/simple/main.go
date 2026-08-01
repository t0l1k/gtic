package main

import (
	"etic"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func main() {
	etic.Load(func() *etic.Cart {
		cart := etic.NewCart("Simple", 320, 200)
		counter := etic.NewProperty(0)
		lbl := etic.NewLabel("label", "Press Enter=Reset Console\nEscape=Quit")
		lbl.Font.Set(etic.FontSystem)
		btn := etic.NewButton("button", "Click me!", func(b *etic.Button) {
			counter.Set(counter.Get() + 1)
		})
		counter.OnChange.Connect(func(i int) {
			result := fmt.Sprintf("Clicked: %v", i)
			btn.Text.Set(result)
		})
		lbl.Bounds().Set(etic.Rect(float32(20), 20, 200, 50))
		btn.Bounds().Set(etic.Rect(float32(20), 70, 200, 100))

		cart.OnBoot = func(t *etic.Console) {
			counter.Set(0)
			btn.Text.Set("Click me!")
			log.Println("Simple cart booted")
		}

		cart.OnTic = func(t *etic.Console) {
			if t.BtnP(ebiten.KeyEscape) {
				t.Exit()
			}
			if t.BtnP(ebiten.KeyEnter) {
				t.Reset()
			}
			t.Cls(colornames.Teal)
			t.Print("Hello, World!", 0, 0, colornames.Yellow)
			lbl.Draw(t)
			btn.Update(t)
			btn.Draw(t)
		}
		return cart
	}()).Run()
}
