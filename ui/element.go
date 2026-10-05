package ui

import (
	"gtic"
	"gtic/react"
	"image"
	"slices"
)

type PropName string
type ElementID string

type ElementState int

const (
	ElementIdle ElementState = iota
	ElementHover
	ElementPressed
	ElementSelected
	ElementDisabled
)

func (s ElementState) String() string {
	return []string{"Normal", "Hover", "Pressed", "Selected", "Disabled"}[s]
}
func (s ElementState) IsIdle() bool     { return s == ElementIdle }
func (s ElementState) IsHover() bool    { return s == ElementHover }
func (s ElementState) IsPressed() bool  { return s == ElementPressed }
func (s ElementState) IsSelected() bool { return s == ElementSelected }
func (s ElementState) IsDisabled() bool { return s == ElementDisabled }

type IElement interface {
	Init(*gtic.API)
	Update(*gtic.API)
	Draw(*gtic.API)
	Property(PropName) *react.Property[any]
}

type Element struct {
	OnInit     func(*gtic.API)
	OnUpdate   func(*gtic.API)
	OnDraw     func(*gtic.API)
	children   []IElement
	properties map[PropName]*react.Property[any]
}

func NewElement(id ElementID) *Element {
	e := &Element{}
	e.RegisterProperty("ID", id)
	e.RegisterProperty("ready", false)
	e.RegisterProperty("state", ElementIdle)
	e.RegisterProperty("hidden", false)
	layout := e.RegisterProperty("layout", AbsoluteLayout{})
	rect := e.RegisterProperty("rect", image.Rectangle{})
	rect.OnChange.Connect(func(a any) {
		r := a.(image.Rectangle)
		l := layout.Get().(Layout)
		if l != nil {
			l.Apply(r, e.children)
		}
	})
	layout.OnChange.Connect(func(a any) { rect.Set(rect.Get()) })
	return e
}
func (e *Element) Children() []IElement {
	children := make([]IElement, len(e.children))
	copy(children, e.children)
	return children
}
func (e *Element) Add(value IElement) { e.children = append(e.children, value) }
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
func (e *Element) Init(t *gtic.API) {
	if e.OnInit != nil {
		e.OnInit(t)
	}
	e.Property("ready").Set(true)
	for _, v := range e.Children() {
		v.Init(t)
	}
}
func (e *Element) Update(t *gtic.API) {
	if e.OnUpdate != nil {
		e.OnUpdate(t)
	}
	for _, v := range e.Children() {
		v.Update(t)
	}
}
func (e *Element) Draw(t *gtic.API) {
	if e.OnDraw != nil {
		e.OnDraw(t)
	}
	for _, v := range e.Children() {
		v.Draw(t)
	}
}

func (e *Element) RegisterUncomparableProperty(name PropName, prop *react.Property[any]) *react.Property[any] {
	if e.properties == nil {
		e.properties = make(map[PropName]*react.Property[any])
	}
	e.properties[name] = prop
	// log.Println("Element:RegisterUncomparableProperty", name, prop.Get())
	return prop
}
func (e *Element) RegisterProperty(name PropName, value any) *react.Property[any] {
	return e.RegisterUncomparableProperty(name, react.NewProperty(value))
}
func (e *Element) Property(name PropName) *react.Property[any] {
	if _, ok := e.properties[name]; !ok {
		for _, v := range e.Children() {
			if v.Property(name) != nil {
				return v.Property(name)
			}
		}
	} else {
		return e.properties[name]
	}
	return nil
}
func (e *Element) Set(name PropName, value any) {
	e.properties[name].Set(value)
	// log.Println("Property:Set", e.properties, name, value)
}
