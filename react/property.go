package react

type Equal[T any] func(a, b T) bool
type Property[T any] struct {
	value    T
	eq       Equal[T]
	OnChange Signal[T]
}

func NewProperty[T comparable](initial T) *Property[T] {
	return &Property[T]{
		value:    initial,
		eq:       func(a, b T) bool { return a == b },
		OnChange: *NewSignal[T](),
	}
}
func NewPropertyWithEqual[T any](initial T, eq Equal[T]) *Property[T] {
	return &Property[T]{
		value:    initial,
		eq:       eq,
		OnChange: *NewSignal[T](),
	}
}
func (p Property[T]) Get() T { return p.value }
func (p *Property[T]) Set(v T) bool {
	if p.eq(p.value, v) {
		// log.Println("Property:Set:Equal", v)
		return false
	}
	p.value = v
	p.OnChange.Emit(v)
	// log.Println("Property:Set", v)
	return true
}
