package tokenizer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMatchSucceeded(t *testing.T) {
	src := "abc "
	tok := New([]byte(src))
	tok.skipSpace()

	ok := tok.match(&tok, 'a', 'b', 'c')

	if match := true; ok != match {
		t.Errorf("expected match %v got %v", match, ok)
	}

	if end := len(src); tok.offset != end {
		t.Errorf("expected offset %d got %d", end, tok.offset)
	}

	if chEnd := 3; tok.chOffset != chEnd {
		t.Errorf("expected ch offset %d got %d", chEnd, tok.chOffset)
	}

	if errCount := 0; tok.errCount != errCount {
		t.Errorf("expected errors %d got %d", errCount, tok.errCount)
	}
}

func TestMatchFailed(t *testing.T) {
	src := "abxd"
	tok := New([]byte(src))
	tok.skipSpace()

	prev := tok

	ok := tok.match(&tok, 'a', 'b', 'c')

	if match := false; ok != match {
		t.Errorf("expected match %v got %v", match, ok)
	}

	// func struct field are never equals unless both are nil
	prev.errFunc = nil
	tok.errFunc = nil

	if diff := cmp.Diff(prev, tok, cmp.AllowUnexported(prev, tok)); diff != "" {
		t.Error(diff)
	}
}
