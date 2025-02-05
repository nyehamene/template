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
	source := ":="
	expected := []Token{
		{source: &source, kind: TokenColon, offset: 0},
		{source: &source, kind: TokenEqual, offset: 1},
	}

	got := []Token{}
	offset := 0

	for {
		var err error
		var token Token
		var end int

		token, end, err = Tokenize(&source, offset)

		if err == EOF {
			break
		}

		if err != nil {
			t.Fatal(err)
			break // unreachable
		}

		got = append(got, token)
		offset = end
	}

	if offset != 2 {
		t.Errorf("expected 2 got %d", offset)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}
