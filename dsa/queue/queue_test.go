package queue

import "testing"

func TestEmpty(t *testing.T) {
	q := New[int](10)

	if !q.Empty() {
		t.Error("a newly created queue should be empty")
	}
	if _, ok := q.Peek(); ok {
		t.Error("queue should be empty")
	}

	q.Push(1)
	q.Push(2)
	q.Push(3)

	q.Pop()
	q.Pop()
	q.Pop()

	if !q.Empty() {
		t.Error("queue should be empty")
	}
	if _, ok := q.Peek(); ok {
		t.Error("queue should be empty")
	}
}

func TestPushAndPop(t *testing.T) {
	q := New[int](10)

	pushed := 1

	q.Push(pushed)

	v, ok := q.Pop()

	if !ok {
		t.Error("pop should return the first value pushed")
	}

	if v != pushed {
		t.Errorf("expected %d got %d", pushed, v)
	}
}

func TestMorePushThanPop(t *testing.T) {
	q := New[int](10)

	pushedFirst := 1
	pushedSecond := 2
	pushedThird := 3

	q.Push(pushedFirst)
	q.Push(pushedSecond)
	q.Push(pushedThird)

	poppedFirst, okPoppedFirst := q.Pop()
	poppedSecond, okPoppedSecond := q.Pop()

	if q.Empty() {
		t.Error("queue should not be empty")
	}

	if expectedRemain, gotRemain := 1, q.Len(); expectedRemain != gotRemain {
		t.Errorf("expected %d item(s) in the queue got %d",
			expectedRemain, gotRemain)
	}

	if !okPoppedFirst {
		t.Error("pop should return the first value pushed")
	}
	if !okPoppedSecond {
		t.Error("pop should return the first value pushed")
	}

	if pushedFirst != poppedFirst {
		t.Errorf("expected %d got %d", pushedFirst, poppedFirst)
	}
	if pushedSecond != poppedSecond {
		t.Errorf("expected %d got %d", pushedSecond, poppedSecond)
	}
}

func TestMorePopThanPush(t *testing.T) {
	q := New[int](10)
	d, ok := q.Pop()

	if ok {
		t.Error("queue should indicate the the prev pop failed by returning false")
	}
	if d != 0 {
		t.Error("queue should return the default value for the queue type")
	}

	q.Push(1)

	q.Pop()
	d, ok = q.Pop()

	if ok {
		t.Error("queue should indicate the the prev pop failed by returning false")
	}
	if d != 0 {
		t.Error("queue should return the default value for the queue type")
	}
}

func TestMorePushThanCapacity(t *testing.T) {
	capacity := 1
	q := New[int](capacity)

	q.PushAll(1, 2)

	if got := q.Cap(); got <= capacity {
		t.Errorf("expected initial capacity %d to be less than %d", capacity, got)
	}
}

func TestCompactLeft(t *testing.T) {
	capacity := 10
	q := New[int](capacity)

	q.PushAll(1, 2, 3, 4, 5, 6, 7, 8, 9, 0)

	// Pop half the content of the queue
	for range 5 {
		v, ok := q.Pop()
		_, _ = v, ok
	}

	// Pushing a new value should move the content of the
	// queue to the left moving the free spaces to the right
	// therefore this push and subsequent pushes will add items
	// in the free space without increase capacity.
	// This make the queue compact
	q.Push(1)

	if c := q.Cap(); c > capacity {
		t.Errorf("capacity should not increase\nexpected %d got %d", capacity, c)
	}

	if expected, got := 6, q.Len(); expected != got {
		t.Errorf("expected length %d got %d", expected, got)
	}

	v, ok := q.Peek()

	if !ok {
		t.Errorf("pop should succeed")
	}

	if expected, got := 6, v; expected != got {
		t.Errorf("expected value %d got %d", expected, got)
	}
}
