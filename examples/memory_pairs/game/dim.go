package game

import (
	"fmt"
	"strconv"
	"strings"
)

type Dim struct{ row, column int }

func (p Dim) String() string {
	return fmt.Sprintf("[%v, %v]", float64(p.row), float64(p.column))
}

func ParseDim(value string) Dim {
	v := strings.Split(value, "x")
	r, err := strconv.Atoi(v[0])
	if err != nil {
		return Dim{}
	}
	c, err := strconv.Atoi(v[1])
	if err != nil {
		return Dim{}
	}
	return Dim{row: r, column: c}
}
