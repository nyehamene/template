package stack

func New[T any](cap int) Stack[T] {
	items := make([]T, 0, cap)
	return Stack[T]{
		items: items,
		next:  0,
	}
}

type Stack[T any] struct {
	items []T
	next  int
}

func (s Stack[T]) Len() int {
	return s.next
}

func (s Stack[T]) Cap() int {
	return cap(s.items)
}

func (s Stack[T]) Empty() bool {
	return s.Len() == 0
}

func (s *Stack[T]) Push(t T) {
	if size := len(s.items); s.next >= size {
		s.items = append(s.items, t)
		s.next = len(s.items)
		return
	}
	next := s.next
	s.items[next] = t
	s.next = next + 1
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.Empty() {
		var t T
		return t, false
	}
	head := s.next - 1
	item := s.items[head]
	s.next = head
	return item, true
}
