package ui

import (
	"gtic"
	"image"
)

func NewIcon(id ElementID, spriteId gtic.SpriteID) *Element {
	e := NewElement(id)
	sId := e.RegisterProperty("sprite.id", spriteId)
	rect := e.Property("rect")
	col := e.RegisterProperty("colorkey", gtic.Transparent)
	scale := e.RegisterProperty("scale", 1)
	flip := e.RegisterProperty("flip", 0)
	rotate := e.RegisterProperty("rotate", 0)
	e.OnDraw = func(a *gtic.API) {
		spId := sId.Get().(gtic.SpriteID)
		r := rect.Get().(image.Rectangle)
		x, y, _, _ := gtic.RectXYWH(r)
		c := col.Get().(gtic.RGBA)
		s := scale.Get().(int)
		f := flip.Get().(int)
		rt := rotate.Get().(int)
		a.Spr(spId, x, y, c, s, f, rt)
	}
	return e
}
