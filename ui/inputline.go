package ui

import (
	"gtic"
	"gtic/react"
	"image"
	"image/color"

	"golang.org/x/image/colornames"
)

func NewInputLine(sc *SceneTree, id ElementID, placeHolder string, fn react.SlotFn[*Element]) *Element {
	const (
		CharTableLower = "!abcdefghijklmnopqrstuvwxyz0123456789-=[]\\;'`,./ "
		CharTableUpper = "!ABCDEFGHIJKLMNOPQRSTUVWXYZ)!@#$%^&*(_+{}|:\"~<>? "
	)
	var (
		cursor     int
		scroll     int
		charsWidth int = 20
		focused    bool
	)

	e := NewElement("input.line")
	e.RegisterProperty("scene.tree", sc)

	text := e.RegisterUncomparableProperty("text", react.NewPropertyWithEqual("", func(a, b any) bool { return false })) // уведомлять всегда
	ph := e.RegisterProperty("placeholder", placeHolder)
	rect := e.Property("rect")

	label := NewLabel("input.line.visible.text", "")
	label.Property("halign").Set(AlignStart)
	label.Property("valign").Set(AlignStart)
	label.Property("bg").Set(colornames.Navy)
	fg := label.Property("fg")
	colTxt := gtic.NewColor(colornames.Yellow)
	colPh := gtic.NewColor(colornames.Gray)
	fg.Set(colTxt)

	cursorEl := NewElement("input.line.cursor")
	cursorEl.OnDraw = func(a *gtic.API) {
		if !focused || a.Time()%800 >= 500 {
			return
		}
		r := cursorEl.Property("rect").Get().(image.Rectangle)
		a.Rect(r.Min.X, r.Min.Y, r.Dx(), r.Dy(), gtic.NewColor(fg.Get().(color.Color)))
	}
	e.Add(cursorEl)

	rect.OnChange.Connect(func(a any) {
		r := a.(image.Rectangle)
		_, _, w, _ := gtic.RectXYWH(r)
		charsWidth = min(w/6, charsWidth*6)
		label.Property("rect").Set(r)
	})

	text.OnChange.Connect(func(a any) {
		result := a.(string)
		if result == "" {
			cursor = 0
		}
		fg.Set(colTxt)
		if result == "" && !focused {
			result = ph.Get().(string)
			fg.Set(colPh)
		}
		start := min(max(scroll, 0), len(result))
		end := min(start+charsWidth, len(result))
		result = result[start:end]
		label.Property("text").Set(result)
	})

	getCh := func(a *gtic.API) string {
		for i := 0; i < len(CharTableLower); i++ {
			if a.Key(gtic.KeySHIFT) && a.KeyP(gtic.Key(i)) {
				return string(CharTableUpper[i])
			} else if a.KeyP(gtic.Key(i)) {
				return string(CharTableLower[i])
			}
		}
		return ""
	}
	e.OnInit = func(a *gtic.API) {}
	e.OnUpdate = func(a *gtic.API) {
		if a.KeyP(gtic.KeyRETURN) {
			fn(e)
		}
		x, y, l, _, _, _, _ := a.Mouse()
		inside := image.Pt(x, y).In(rect.Get().(image.Rectangle))
		if inside && l && !focused {
			sc.Focus(id)
		}
		focused = sc.focusedId == id
		if !focused {
			return
		}

		input := text.Get().(string)
		lastCursor := cursor
		ch := getCh(a)
		if ch != "" {
			input = input[:cursor] + ch + input[cursor:]
			cursor++
		}
		switch {
		case a.KeyP(gtic.KeyLEFT):
			cursor = max(0, cursor-1)
		case a.KeyP(gtic.KeyRIGHT):
			cursor = min(len(input), cursor+1)
		case a.KeyP(gtic.KeyHOME):
			cursor = 0
		case a.KeyP(gtic.KeyEND):
			cursor = len(input)
		case a.KeyP(gtic.KeyBACKSPACE, 20, 3):
			if cursor > 0 && len(input) > 0 {
				input = input[:cursor-1] + input[cursor:]
				cursor--
			}
		case a.KeyP(gtic.KeyDELETE):
			if cursor >= 0 && cursor < len(input) {
				input = input[:cursor] + input[cursor+1:]
			}
		}
		if cursor < scroll {
			scroll = cursor
		} else if cursor >= scroll+charsWidth {
			scroll = cursor - charsWidth + 1
		} else if len(input) >= cursor && cursor > 0 && scroll > 0 {
			scroll = cursor - charsWidth + 1
		}
		if input != text.Get().(string) || lastCursor != cursor {
			text.Set(input)
		}
		fieldRect := rect.Get().(image.Rectangle)
		x0, y0, _, _ := gtic.RectXYWH(fieldRect)
		cursorX := x0 + (cursor-scroll)*6
		cursorEl.Property("rect").Set(image.Rect(cursorX, y0+7, cursorX+5, y0+8))
	}
	e.OnDraw = func(a *gtic.API) { label.Draw(a) }
	return e
}
