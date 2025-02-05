package template

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tokens
// chars : = ; { }
// keywords package component tag list html
//          type record templ
// primary ident string

func TestNext_Token(t *testing.T) {
	source := ":"
	expected := Token{
		source: new(string),
		kind:   TokenColon,
		offset: 0,
	}
	got, n, err := Tokenize(&source, 0)

	if err != EOF {
		t.Error(err)
	}

	if n != 1 {
		t.Errorf("expected 1 got %d", n)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}
