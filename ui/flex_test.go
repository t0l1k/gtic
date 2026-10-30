package ui

import (
	"image"
	"testing"
)

func TestFlexKeepsChildrenInsideParent(t *testing.T) {
	tests := []struct {
		name    string
		dir     FlexDir
		justify JustifyType
		padding int
		rect    image.Rectangle
	}{
		{"row start", Row, JustifyStart, 5, image.Rect(0, 0, 200, 100)},
		{"row center", Row, JustifyCenter, 5, image.Rect(0, 0, 200, 100)},
		{"row end", Row, JustifyEnd, 5, image.Rect(0, 0, 200, 100)},
		{"row between", Row, JustifySpaceBetween, 5, image.Rect(0, 0, 200, 100)},
		{"row around", Row, JustifySpaceAround, 5, image.Rect(0, 0, 200, 100)},
		{"column start", Col, JustifyStart, 3, image.Rect(0, 0, 100, 200)},
		{"column center", Col, JustifyCenter, 3, image.Rect(0, 0, 100, 200)},
		{"column end", Col, JustifyEnd, 3, image.Rect(0, 0, 100, 200)},
		{"column between", Col, JustifySpaceBetween, 3, image.Rect(0, 0, 100, 200)},
		{"column around", Col, JustifySpaceAround, 3, image.Rect(0, 0, 100, 200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := NewElement("parent")
			for range 4 {
				parent.Add(NewElement("child"))
			}

			parent.Property("layout").Set(NewFlex(tt.dir, tt.padding, AlignCenter, tt.justify, nil))
			parent.Property("rect").Set(tt.rect)

			for i, child := range parent.Children() {
				r := child.Property("rect").Get().(image.Rectangle)
				if r.Min.X < tt.rect.Min.X+tt.padding ||
					r.Min.Y < tt.rect.Min.Y+tt.padding ||
					r.Max.X > tt.rect.Max.X-tt.padding ||
					r.Max.Y > tt.rect.Max.Y-tt.padding {
					t.Errorf("child %d rect %v is outside padded parent %v", i, r, tt.rect)
				}
			}
		})
	}
}

func TestFlexGap(t *testing.T) {
	tests := []struct {
		name       string
		dir        FlexDir
		gap        int
		rowGap     int
		columnGap  int
		parentRect image.Rectangle
		firstRect  image.Rectangle
		secondRect image.Rectangle
	}{
		{
			name:       "gap in row",
			dir:        Row,
			gap:        10,
			parentRect: image.Rect(0, 0, 100, 40),
			firstRect:  image.Rect(0, 0, 45, 40),
			secondRect: image.Rect(55, 0, 100, 40),
		},
		{
			name:       "gap in column",
			dir:        Col,
			gap:        10,
			parentRect: image.Rect(0, 0, 40, 100),
			firstRect:  image.Rect(0, 0, 40, 45),
			secondRect: image.Rect(0, 55, 40, 100),
		},
		{
			name:       "column gap overrides shorthand",
			dir:        Row,
			gap:        10,
			columnGap:  4,
			parentRect: image.Rect(0, 0, 100, 40),
			firstRect:  image.Rect(0, 0, 48, 40),
			secondRect: image.Rect(52, 0, 100, 40),
		},
		{
			name:       "row gap overrides shorthand",
			dir:        Col,
			gap:        10,
			rowGap:     4,
			parentRect: image.Rect(0, 0, 40, 100),
			firstRect:  image.Rect(0, 0, 40, 48),
			secondRect: image.Rect(0, 52, 40, 100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := NewElement("parent")
			first := NewElement("first")
			second := NewElement("second")
			parent.Add(first)
			parent.Add(second)

			layout := NewFlex(tt.dir, 0, AlignStart, JustifyStart, nil)
			layout.Gap = tt.gap
			layout.RowGap = tt.rowGap
			layout.ColumnGap = tt.columnGap
			parent.Property("layout").Set(layout)
			parent.Property("rect").Set(tt.parentRect)

			if got := first.Property("rect").Get().(image.Rectangle); got != tt.firstRect {
				t.Errorf("first child rect = %v, want %v", got, tt.firstRect)
			}
			if got := second.Property("rect").Get().(image.Rectangle); got != tt.secondRect {
				t.Errorf("second child rect = %v, want %v", got, tt.secondRect)
			}
		})
	}
}
