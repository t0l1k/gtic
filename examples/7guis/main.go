package main

import (
	"gtic"
	"gtic/react"
	"gtic/ui"
	"log"
	"strconv"
)

var (
	Title        = "7Guis Examples"
	Counter      = "counter"
	TempConv     = "temperature converter"
	FlightBooker = "flight booker"
	Timer        = "timer"
	Crud         = "crud"
	Circle       = "circle"
	Cells        = "cells"
	QuitDemo     = "<"
	QuitApp      = "X"
)

func NewCounter(id ui.ElementID, fn react.SlotFn[*ui.Element]) *ui.Element {
	s := ui.NewElement(id)
	counter := react.NewProperty(0)
	lblCount := ui.NewLabel("scene.counter.label.count", "0")
	lblCount.Property("scale").Set(3)
	btnInc := ui.NewButton("scene.counter.button.inc", "Inc", func(b *ui.Element) {
		counter.Set(counter.Get() + 1)
	})
	btnDec := ui.NewButton("scene.counter.button.dec", "Dec", func(b *ui.Element) {
		counter.Set(counter.Get() - 1)
	})
	btnReset := ui.NewButton("scene.counter.button.reset", "Reset", func(b *ui.Element) {
		counter.Set(0)
	})
	counter.OnChange.Connect(func(i int) {
		lblCount.Property("text").Set(strconv.Itoa(i))
	})

	s.Property("layout").Set(ui.NewVerticalFlexLayout([]float32{0.1, 0.7, 0.2}, 0.05))
	s.Add(ui.TopBar("scene.counter.topbar", Counter, QuitDemo, fn))
	s.Add(lblCount)

	contBtns := ui.NewElement("scene.counter.buttons.container")
	contBtns.Property("layout").Set(ui.NewHorizontalBoxLayout(0.05))
	contBtns.Add(btnInc)
	contBtns.Add(btnDec)
	contBtns.Add(btnReset)
	s.Add(contBtns)
	return s
}

func sceneMain(fn react.SlotFn[*ui.Element]) *ui.Element {
	s := ui.NewElement("scene.main")
	s.Property("layout").Set(ui.NewVerticalBoxLayout(0.05))
	s.Add(ui.TopBar("main.scene.topbar", Title, QuitApp, fn))
	for _, txt := range []string{Counter, TempConv, FlightBooker, Timer, Crud, Circle, Cells} {
		s.Add(ui.NewButton("scane.main.button", txt, fn))
	}
	return s
}

func main() {
	SC := ui.NewSceneTree()
	gtic.BOOT = func(a *gtic.API) {
		counter := NewCounter("scene.counter", func(e *ui.Element) {
			txt := e.Property("text")
			switch txt.Get().(string) {
			case QuitDemo:
				SC.Change(a, "scene.main")
			}
		})
		main := sceneMain(func(e *ui.Element) {
			txt := e.Property("text")
			switch txt.Get().(string) {
			case Counter:
				SC.Change(a, counter.Property("ID").Get().(ui.ElementID))
			case QuitApp:
				a.Exit()
			}
		})
		SC.AddScene(main)
		SC.AddScene(counter)
		SC.Change(a, main.Property("ID").Get().(ui.ElementID))
	}
	gtic.TIC = func(a *gtic.API) {
		a.Cls()
		SC.TIC(a)
	}
	if err := gtic.Load(gtic.WithTitle(Title)).Run(); err != nil {
		log.Println(err)
	}
}
