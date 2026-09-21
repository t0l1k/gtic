package gtic

import "image/color"

const (
	fontW = 6 // ширина глифа и шаг
	fontH = 8 // высота строки
)

// Print(text[, x=0][, y=0][, color=White][, fixed=false][, scale=1][, smallfont=false]) -> width
func (a *API) Print(text string, args ...interface{}) int {
	x, y := 0, 0
	col := NewRGBA(255, 255, 255, 255)
	fixed := false
	scale := 1
	smallfont := false

	i := 0
	nextInt := func(def int) int {
		if i < len(args) {
			if v, ok := args[i].(int); ok {
				i++
				return v
			}
		}
		return def
	}
	nextRGBA := func(def RGBA) RGBA {
		if i < len(args) {
			if v, ok := args[i].(RGBA); ok {
				i++
				return v
			}
		}
		return def
	}
	nextBool := func(def bool) bool {
		if i < len(args) {
			if v, ok := args[i].(bool); ok {
				i++
				return v
			}
		}
		return def
	}

	// x и y первые два int, если есть
	if len(args) > 0 {
		if v, ok := args[0].(int); ok {
			x = v
			i = 1
			y = nextInt(0)
		}
	}
	col = nextRGBA(col)
	fixed = nextBool(fixed)
	scale = nextInt(scale)
	smallfont = nextBool(smallfont)

	if scale < 1 {
		scale = 1
	}

	f := SystemFont
	spaceAdv := 3 // пропорциональный пробел, как в TIC-80
	if fixed {
		spaceAdv = f.Advance
	}

	cx, cy, maxW := x, y, 0
	for _, r := range text {
		switch {
		case r == '\n':
			maxW = max(maxW, cx-x)
			cx, cy = x, cy+f.Height*scale
		case r == ' ':
			cx += spaceAdv * scale
		default:
			a.drawGlyph(f.glyph(r), cx, cy, col, scale)
			cx += f.Advance * scale
		}
	}
	return max(maxW, cx-x)
}

type Font struct {
	glyphs  map[rune]*Sprite
	Advance int
	Height  int
}

func (f *Font) glyph(r rune) *Sprite {
	if g, ok := f.glyphs[r]; ok {
		return g
	}
	return f.glyphs['?'] // fallback — вы это уже видели :)
}

var SystemFont = func() *Font {
	f := &Font{glyphs: make(map[rune]*Sprite, len(systemFontHex)), Advance: fontW, Height: fontH}
	for r, h := range systemFontHex {
		f.glyphs[r] = glyphFromHex(h)
	}
	return f
}()

func (a *API) drawGlyph(s *Sprite, x, y int, color RGBA, scale int) {
	for gy := 0; gy < s.Height; gy++ {
		for gx := 0; gx < s.Width; gx++ {
			_, _, _, aa := s.Pixels[gy*s.Width+gx].ToBytes()
			if aa == 0 {
				continue
			}
			for sy := 0; sy < scale; sy++ {
				for sx := 0; sx < scale; sx++ {
					a.vram.set(x+gx*scale+sx, y+gy*scale+sy, color)
				}
			}
		}
	}
}

// glyphFromHex("A", "1C22223E22222200")
func glyphFromHex(hex string) *Sprite {
	px := make([]RGBA, fontW*fontH)
	for row := 0; row < fontH; row++ {
		var v uint8
		for k := 0; k < 2; k++ {
			c := hex[row*2+k]
			v = v<<4 | unhex(c) // '0'-'9','A'-'F' -> 0..15
		}
		for col := 0; col < fontW; col++ {
			if v&(1<<uint(5-col)) != 0 {
				px[row*fontW+col] = NewColor(color.White)
			}
		}
	}
	return &Sprite{Width: fontW, Height: fontH, Pixels: px}
}
func unhex(c byte) uint8 {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	default:
		return 0
	}
}

var systemFontHex = map[rune]string{
	' ':  "0000000000000000",
	'!':  "0808080808000800",
	'"':  "1212000000000000",
	'#':  "14143F143F141400",
	'$':  "081E281C0A3C0800",
	'%':  "1819020408160600",
	'&':  "0C12140815120D00",
	'\'': "0808000000000000",
	'(':  "0408101010080400",
	')':  "1008040404081000",
	'*':  "00082A1C2A080000",
	'+':  "0008083F08080000",
	',':  "00000000000C0810",
	'-':  "0000003F00000000",
	'.':  "0000000000000C00",
	'/':  "0102040810102000",
	'0':  "1C22262A32221C00",
	'1':  "0818080808081C00",
	'2':  "1C22020408103E00",
	'3':  "1C22020C02221C00",
	'4':  "040C14243E040400",
	'5':  "3E203C0202221C00",
	'6':  "0C10203C22221C00",
	'7':  "3E02040808101000",
	'8':  "1C22221C22221C00",
	'9':  "1C22221E02041800",
	':':  "00000C00000C0000",
	';':  "00000C00000C0810",
	'<':  "0408102010080400",
	'=':  "00003F00003F0000",
	'>':  "1008040204081000",
	'?':  "1C22020C08000800",
	'@':  "1C22021A2A2E1C00",
	'A':  "1C22223E22222200",
	'B':  "3C22223C22223C00",
	'C':  "1C22202020221C00",
	'D':  "3824222222243800",
	'E':  "3E20203C20203E00",
	'F':  "3E20203C20202000",
	'G':  "1C22202E22221E00",
	'H':  "2222223E22222200",
	'I':  "1C08080808081C00",
	'J':  "0E04040404241800",
	'K':  "2224283028242200",
	'L':  "2020202020203E00",
	'M':  "22362A2A22222200",
	'N':  "22322A2622222200",
	'O':  "1C22222222221C00",
	'P':  "3C22223C20202000",
	'Q':  "1C2222222A241A00",
	'R':  "3C22223C28242200",
	'S':  "1C22201C02221C00",
	'T':  "3E08080808080800",
	'U':  "2222222222221C00",
	'V':  "2222222222140800",
	'W':  "2222222A2A2A1400",
	'X':  "2222140814222200",
	'Y':  "2222140808080800",
	'Z':  "3E02040810203E00",
	'[':  "0C08080808080C00",
	'\\': "2010100804020200",
	']':  "0C04040404040C00",
	'^':  "0814220000000000",
	'_':  "000000000000003E",
	'`':  "1008000000000000",
	'a':  "00001C021E221E00",
	'b':  "20202C3222223C00",
	'c':  "00001C2020221C00",
	'd':  "02021A2622221E00",
	'e':  "00001C223E201C00",
	'f':  "0C101C1010101000",
	'g':  "00001E22221E021C",
	'h':  "20202C3222222200",
	'i':  "0800180808081C00",
	'j':  "04000C0404042418",
	'k':  "2020242830282400",
	'l':  "1808080808081C00",
	'm':  "00002C362A2A2A00",
	'n':  "00002C3222222200",
	'o':  "00001C2222221C00",
	'p':  "00003C22223C2020",
	'q':  "00001E22221E0202",
	'r':  "00002C3220202000",
	's':  "00001E201C023C00",
	't':  "10103C1010120C00",
	'u':  "0000222222261A00",
	'v':  "0000222222140800",
	'w':  "000022222A2A1400",
	'x':  "0000221408142200",
	'y':  "00002222221E021C",
	'z':  "00003E0408103E00",
	'{':  "0408081008080400",
	'|':  "0808080808080808",
	'}':  "1008080408081000",
	'~':  "0000192600000000",
}
