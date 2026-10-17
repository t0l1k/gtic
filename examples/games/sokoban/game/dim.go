package game

import "strings"

type Dim struct{ Row, Column int }

func NewDim(level string) Dim {
	rows := strings.Split(strings.TrimSuffix(level, "\n"), "\n")
	d := Dim{Column: len(rows)}
	for _, row := range rows {
		if len(row) > d.Row {
			d.Row = len(row)
		}
	}
	return d
}
func (d Dim) Idx(x, y int) int       { return y*d.Row + x }
func (d Dim) Pos(idx int) (int, int) { return idx % d.Row, idx / d.Row }
