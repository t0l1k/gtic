package ui

import (
	"gtic"
	"gtic/react"
	"slices"
)

type PropName string
type ElementID string
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
	ready      bool
	properties map[PropName]*react.Property[any]
}

func NewElement(id ElementID) *Element {
	e := &Element{}
	e.RegisterProperty("ID", id)
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
	e.ready = true
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
