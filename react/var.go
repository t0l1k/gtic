package react

// Var Реактивная переменная с проверкой на равенство, при равенстве и повторном обновлении пропускать уведомление подписчиков
type Var[T any] struct {
	v   T
	ver uint64
	eq  Equal[T]
}

func NewVar[T comparable](initial T) *Var[T] {
	return NewVarWithEqual(initial, func(a, b T) bool { return a == b })
}
func NewVarWithEqual[T any](initial T, equals Equal[T]) *Var[T] {
	if equals == nil {
		equals = func(T, T) bool { return false }
	}
	return &Var[T]{v: initial, eq: equals}
}
func (s *Var[T]) Get() T          { return s.v }
func (s *Var[T]) Version() uint64 { return s.ver }

// Set обновляет значение и возвращает true, если оно изменилось.
func (s *Var[T]) Set(v T) bool {
	if s.eq(s.v, v) {
		return false
	}
	s.v = v
	s.ver++
	return true
}

// Токен отслеживает версию, зафиксированную при последнем обращении, для обнаружения изменений.
type Token uint64

func (s *Var[T]) Changed(t *Token) bool {
	if t == nil {
		return false
	}
	if uint64(*t) == s.ver {
		return false
	}
	*t = Token(s.ver)
	return true
}
