package queue

func New[T any](capacity int) Queue[T] {
	t := make([]T, capacity)
	return Queue[T]{t, 0, 0}
}

type Queue[T any] struct {
	t []T
	r int
	w int
}

func (q Queue[T]) Len() int {
	length := q.w - q.r
	return length
}

func (q Queue[T]) Cap() int {
	return cap(q.t)
}

func (q Queue[T]) Empty() bool {
	return q.r == q.w
}

func (q *Queue[T]) Push(t T) {
	q.compact()
	if q.w == len(q.t) {
		q.t = append(q.t, t)
	} else {
		index := q.w
		q.t[index] = t
	}
	q.w += 1
}

func (q *Queue[T]) PushAll(ts ...T) {
	for _, t := range ts {
		q.Push(t)
	}
}

func (q *Queue[T]) Pop() (T, bool) {
	if q.Empty() {
		q.r = 0
		q.w = 0
		var empty T
		return empty, false
	}

	index := q.r
	t := q.t[index]
	q.r += 1
	return t, true
}

func (q *Queue[T]) Peek() (T, bool) {
	if q.Empty() {
		var empty T
		return empty, false
	}
	index := q.r
	t := q.t[index]
	return t, true
}

func (q *Queue[T]) compact() {
	size := len(q.t)
	quotient := 2
	halfSize := size / quotient

	if q.r < halfSize {
		return
	}

	var empty T
	i := 0

	for ; i < q.Len(); i++ {
		index := q.r + i
		q.t[i] = q.t[index]
		q.t[index] = empty
	}
	q.w = i
	q.r = 0
}
