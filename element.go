package etic

import (
	"image/color"
	"slices"

	"golang.org/x/image/colornames"
)

type IElement interface {
	Init(*Console)
	Update(*Console)
	Draw(*Console)

	ID() *Property[ElementID]

	Bounds() *Property[Rectangle[float32]]
	Children() []IElement

	IsHidden() bool
	Hide()
	Show()

	IsReady() bool
	MarkReady()
	ClearReady()
}
type ElementID string

type ElementState int

func (b ElementState) Color() color.Color {
	return []color.Color{colornames.Gainsboro, colornames.Gray, colornames.Red, colornames.Blue, colornames.Darkgray}[b]
}
func (s ElementState) String() string {
	return []string{"Normal", "Hover", "Pressed", "Selected", "Disabled"}[s]
}

const (
	ElementNormal ElementState = iota
	ElementHover
	ElementPressed
	ElementSelected
	ElementDisabled
)

type Element struct {
	OnInit   func(*Console)
	OnUpdate func(*Console)
	OnDraw   func(*Console)
	id       Property[ElementID]
	rect     Property[Rectangle[float32]]
	hidden   Property[bool]
	layout   Property[Layout]
	children []IElement
	ready    bool
}

func NewElement(id ElementID) *Element {
	e := &Element{
		id:     *NewProperty(id),
		rect:   *NewPropertyWithEqual(Rectangle[float32]{}, func(a, b Rectangle[float32]) bool { return false }),
		hidden: *NewProperty(false),
		layout: *NewProperty[Layout](AbsoluteLayout{}),
	}
	e.Bounds().OnChange.Connect(func(r Rectangle[float32]) {
		if e.layout.Get() != nil {
			e.Layout().Get().Apply(r, e.Children())
		}
		// log.Println("Element:Init:rect", e.ID.Get(), e.Bounds().Get(), r)
	})
	return e
}

func (e *Element) Children() []IElement {
	children := make([]IElement, len(e.children))
	copy(children, e.children)
	return children
}
func (e *Element) Add(value IElement) {
	e.children = append(e.children, value)
	if e.Bounds().Get().Empty() {
		return
	}
	e.Bounds().Set(e.Bounds().Get())
}
func (e *Element) ResetChildren() {
	for _, v := range e.Children() {
		e.Remove(v)
	}
	e.children = nil
}
func (e *Element) Remove(value IElement) {
	if value == nil {
		return
	}
	e.children = slices.DeleteFunc(e.Children(), func(c IElement) bool { return c == value })
}
func (e *Element) Init(t *Console) {
	if e.OnInit != nil {
		e.OnInit(t)
	}
	e.MarkReady()
	for _, v := range e.Children() {
		v.Init(t)
	}
	// log.Println("Element:Init", e.ID.Get(), e.Bounds().Get())
}
func (e *Element) Update(t *Console) {
	if e.OnUpdate != nil {
		e.OnUpdate(t)
	}
	for _, v := range e.Children() {
		v.Update(t)
	}
}
func (e *Element) Draw(t *Console) {
	if e.OnDraw != nil {
		e.OnDraw(t)
	}
	for _, v := range e.Children() {
		v.Draw(t)
	}
}
func (e *Element) Bounds() *Property[Rectangle[float32]] { return &e.rect }
func (e *Element) Layout() *Property[Layout]             { return &e.layout }

func (e *Element) IsReady() bool { return e.ready }
func (e *Element) MarkReady() {
	e.ready = true
	for _, v := range e.Children() {
		v.MarkReady()
	}
}
func (e *Element) ClearReady() {
	e.ready = false
	for _, v := range e.Children() {
		v.ClearReady()
	}
}

func (e *Element) IsHidden() bool { return e.hidden.Get() }
func (e *Element) Hide() {
	e.hidden.Set(true)
	for _, v := range e.Children() {
		v.Hide()
	}
}
func (e *Element) Show() {
	e.hidden.Set(false)
	for _, v := range e.Children() {
		v.Show()
	}
}

func (e *Element) ID() *Property[ElementID] { return &e.id }
