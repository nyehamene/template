package stack_test

import (
	"temlang/tem/dsa/stack"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStack(t *testing.T) {
	const capacity = 10
	s := stack.New[string](capacity)

	if !s.Empty() {
		t.Errorf("expected empty stack")
	}

	if s.Len() != 0 {
		t.Errorf("expected stack len to be 0")
	}

	if c := s.Cap(); c != capacity {
		t.Errorf("expected stack cap to be %d got %d", capacity, c)
	}

	if _, ok := s.Pop(); ok {
		t.Errorf("expected stack pop to fail")
	}

	const expected = "second"
	s.Push("first")
	s.Push(expected)

	got, ok := s.Pop()

	s.Push("third")

	if !ok {
		t.Errorf("expected stack pop to succeed")
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}

	if s.Empty() {
		t.Errorf("expected stack to not be empty")
	}

	if expected, l := 2, s.Len(); l != expected {
		t.Errorf("expected stack len %d got %d", expected, l)
	}
}
