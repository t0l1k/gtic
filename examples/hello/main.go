package main

import (
	"gtic"
	"log"
)

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
