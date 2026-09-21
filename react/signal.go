package react

import "sort"

type SlotID uint64
type SlotFn[T any] func(T)
type Connection func()

type signalSlot[T any] struct {
	fn   SlotFn[T]
	once bool
}

type Signal[T any] struct {
	slots map[SlotID]signalSlot[T]
	next  SlotID
}

func NewSignal[T any]() *Signal[T]                          { return &Signal[T]{slots: make(map[SlotID]signalSlot[T])} }
func (s *Signal[T]) Connect(fn SlotFn[T]) Connection        { return s.connect(fn, false) }
func (s *Signal[T]) ConnectOneShot(fn SlotFn[T]) Connection { return s.connect(fn, true) }
func (s *Signal[T]) connect(fn SlotFn[T], once bool) Connection {
	if fn == nil {
		return func() {}
	}
	if s.slots == nil {
		s.slots = make(map[SlotID]signalSlot[T])
	}
	id := s.next
	s.next++
	s.slots[id] = signalSlot[T]{fn: fn, once: once}
	return func() {
		delete(s.slots, id)
	}
}
func (s *Signal[T]) Emit(value T) {
	if len(s.slots) == 0 {
		return
	}
	ids := make([]int, 0, len(s.slots))
	for id := range s.slots {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	for _, rawID := range ids {
		id := SlotID(rawID)
		slot, ok := s.slots[id]
		if !ok {
			continue
		}
		slot.fn(value)
		if slot.once {
			delete(s.slots, id)
		}
	}
}
func (s *Signal[T]) DisconnectAll() {
	for id := range s.slots {
		delete(s.slots, id)
	}
}
