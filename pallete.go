package gtic

func (a *API) Pal(i int, c ...RGBA) RGBA {
	if len(c) > 0 {
		a.pallete[i] = c[0]
	}
	return a.pallete[i]
}

func (a *API) ResetPal() {
	if a.pallete == nil {
		a.pallete = make(map[int]RGBA)
	}
	for k, v := range palleteSweetie {
		a.pallete[k] = v
	}
}

func (a *API) ParsePal(c RGBA) int {
	for i, v := range a.Pallete() {
		if v.Equal(c) {
			return i
		}
	}
	return -1
}

func (a *API) Pallete() (result []RGBA) {
	for i := 0; i < len(a.pallete); i++ {
		result = append(result, a.pallete[i])
	}
	return result
}

var palleteSweetie = map[int]RGBA{
	0:  NewRGBA(26, 28, 44, 255),    // 0: #1A1C2C — Черный (Black)
	1:  NewRGBA(93, 39, 93, 255),    // 1: #5D275D — Пурпурный (Purple)
	2:  NewRGBA(177, 62, 83, 255),   // 2: #B13E53 — Красный (Red)
	3:  NewRGBA(239, 125, 87, 255),  // 3: #EF7D57 — Оранжевый (Orange)
	4:  NewRGBA(255, 205, 117, 255), // 4: #FFCD75 — Желтый (Yellow)
	5:  NewRGBA(167, 240, 112, 255), // 5: #A7F070 — Салатовый (Light green)
	6:  NewRGBA(56, 183, 100, 255),  // 6: #38B764 — Зеленый (Green)
	7:  NewRGBA(37, 113, 121, 255),  // 7: #257179 — Темно-зеленый (Dark green)
	8:  NewRGBA(41, 54, 111, 255),   // 8: #29366F — Темно-синий (Dark blue)
	9:  NewRGBA(59, 93, 201, 255),   // 9: #3B5DC9 — Синий (Blue)
	10: NewRGBA(65, 166, 246, 255),  // 10: #41A6F6 — Голубой (Light blue)
	11: NewRGBA(115, 239, 247, 255), // 11: #73EFF7 — Бирюзовый (Cyan)
	12: NewRGBA(244, 244, 244, 255), // 12: #F4F4F4 — Белый (White)
	13: NewRGBA(148, 176, 194, 255), // 13: #94B0C2 — Светло-серый (Light Grey)
	14: NewRGBA(86, 108, 134, 255),  // 14: #566C86 — Серый (Grey)
	15: NewRGBA(51, 60, 87, 255),    // 15: #333C57 — Темно-серый (Dark Grey)
}
