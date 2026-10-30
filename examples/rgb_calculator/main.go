package main

import (
	"gtic"
	"gtic/ui"
	"image"
	"log"
	"math/rand/v2"

	"golang.org/x/image/colornames"
)

func NewColorPreview() *ui.Element {
	e := ui.NewElement("color.preview")
	col := e.RegisterProperty("color", gtic.NewColor(colornames.Black))
	rect := e.Property("rect")
	e.OnDraw = func(a *gtic.API) {
		x, y, w, h := gtic.RectXYWH(rect.Get().(image.Rectangle))
		a.Rect(x, y, w, h, col.Get().(gtic.RGBA))
	}
	return e
}

func main() {
	st := ui.NewSceneTree()
	preview := NewColorPreview()
	sR := ui.NewProgress("id.red", 0, 0, 1, ui.Horizontal)
	sG := ui.NewProgress("id.green", 0, 0, 1, ui.Horizontal)
	sB := ui.NewProgress("id.blue", 0, 0, 1, ui.Horizontal)
	sA := ui.NewProgress("id.alpha", 1, 0, 1, ui.Horizontal)
	bRnd := ui.NewButton("button.random.color", "Next Random Color", func(e *ui.Element) {
		nextColor := func() (uint8, uint8, uint8, uint8) {
			return uint8(rand.IntN(255)), uint8(rand.IntN(255)), uint8(rand.IntN(255)), 255
		}
		r, g, b, a := nextColor()
		preview.Property("color").Set(gtic.NewRGBA(r, g, b, a))
		rV := gtic.PercentOf(float32(r), 255) * 0.01
		gV := gtic.PercentOf(float32(g), 255) * 0.01
		bV := gtic.PercentOf(float32(b), 255) * 0.01
		sR.Property("value").Set(rV)
		sG.Property("value").Set(gV)
		sB.Property("value").Set(bV)
		sR.Property("bg").Set(gtic.NewRGBA(r, 0, 0, 255))
		sG.Property("bg").Set(gtic.NewRGBA(0, g, 0, 255))
		sB.Property("bg").Set(gtic.NewRGBA(0, 0, b, 255))
		log.Println("color", r, g, b, a, rV, gV, bV)
	})

	right := ui.NewElement("id.right")
	// Внутренняя колонка: отступ 3 px от краёв и 2 px между строками.
	rightLayout := ui.NewFlex(ui.Col, 3, ui.AlignCenter, ui.JustifySpaceBetween, []float64{1, 1, 1, 1, 0.5})
	rightLayout.Gap = 2
	right.Property("layout").Set(rightLayout)
	right.Add(sR).Add(sG).Add(sB).Add(sA).Add(bRnd)

	root := ui.NewElement("scene.root")
	// Корень: отступ 5 px от краёв, 5 px между колонками; ширины относятся как 1:2.
	rootLayout := ui.NewFlex(ui.Row, 5, ui.AlignCenter, ui.JustifySpaceBetween, []float64{1, 2})
	rootLayout.Gap = 5
	root.Property("layout").Set(rootLayout)
	root.Add(preview).Add(right)

	gtic.BOOT = func(a *gtic.API) {
		st.AddScene(root)
		st.Change(a, root.Property("id").Get().(ui.ElementID))
	}
	gtic.TIC = func(a *gtic.API) {
		if a.Key(gtic.KeyESC) {
			a.Exit()
		}
		if a.KeyP(gtic.KeyRETURN) {
			a.Reset()
		}
		a.Cls()
		st.TIC(a)
	}
	if err := gtic.Load(gtic.WithTitle("RGB Calculator")).Run(); err != nil {
		log.Println(err)
	}
}
