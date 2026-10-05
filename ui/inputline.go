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
		input      string
		cursor     int
		x, y, w, _ int
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

	rect.OnChange.Connect(func(a any) {
		r := a.(image.Rectangle)
		x, y, w, _ = gtic.RectXYWH(r)
		charsWidth = min(w/6, charsWidth*6)
		label.Property("rect").Set(r)
	})

	text.OnChange.Connect(func(a any) {
		txt := a.(string)
		fg.Set(colTxt)
		if txt == "" && !focused {
			txt = ph.Get().(string)
			fg.Set(colPh)
		}
		start := min(max(scroll, 0), len(txt))
		end := min(start+charsWidth, len(txt))
		txt = txt[start:end]
		label.Property("text").Set(txt)
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
	e.OnInit = func(a *gtic.API) { cursor = len(input) }
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

		prevInput := text.Get().(string)
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
		if prevInput != input || lastCursor != cursor {
			text.Set(input)
		}
	}
	e.OnDraw = func(a *gtic.API) {
		label.Draw(a)
		if a.Time()%800 < 500 && focused {
			a.Rect(x+(cursor-scroll)*6, y+7, 5, 1, gtic.NewColor(fg.Get().(color.Color)))
		}
	}
	return e
}
