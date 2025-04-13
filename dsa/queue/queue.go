package queue

func New[T any](cap int) Queue[T] {
	t := make([]T, cap)
	return Queue[T]{
		items:       t,
		readOffset:  0,
		writeOffset: 0,
	}
}

type Queue[T any] struct {
	items       []T
	readOffset  int
	writeOffset int
}

func (q Queue[T]) Len() int {
	length := q.writeOffset - q.readOffset
	return length
}

func (q Queue[T]) Cap() int {
	return cap(q.items)
}

func (q Queue[T]) Empty() bool {
	return q.readOffset == q.writeOffset
}

func (q *Queue[T]) Push(t T) {
	q.compact()
	if q.writeOffset == len(q.items) {
		q.items = append(q.items, t)
	} else {
		index := q.writeOffset
		q.items[index] = t
	}
	q.writeOffset += 1
}

func (q *Queue[T]) PushAll(ts ...T) {
	for _, t := range ts {
		q.Push(t)
	}
}

func (q *Queue[T]) Pop() (T, bool) {
	if q.Empty() {
		var empty T
		return empty, false
	}

	index := q.readOffset
	t := q.items[index]
	q.readOffset += 1
	return t, true
}

func (q *Queue[T]) Peek() (T, bool) {
	if q.Empty() {
		var empty T
		return empty, false
	}
	index := q.readOffset
	t := q.items[index]
	return t, true
}

func (q *Queue[T]) compact() {
	if q.Empty() {
		return
	}
	if q.writeOffset < 10 {
		return
	}

	itemSize := q.Len()
	size := len(q.items)
	quotient := 2
	halfSize := size / quotient

	if itemSize < halfSize {
		return
	}

	var empty T
	i := 0

	for ; i < itemSize; i++ {
		index := q.readOffset + i
		cur := q.items[index]
		q.items[i] = cur
		q.items[index] = empty
	}
	q.writeOffset = i
	q.readOffset = 0
}
