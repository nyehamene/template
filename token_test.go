package template

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tokens
// chars : = ; { } .
// keywords package component tag list html
//          type record templ
// primary ident string

func TestNext_char(t *testing.T) {
	source := ":=.;{}[]"
	expected := []Token{
		{source: &source, kind: TokenColon, offset: 0},
		{source: &source, kind: TokenEqual, offset: 1},
		{source: &source, kind: TokenPeriod, offset: 2},
		{source: &source, kind: TokenSemicolon, offset: 3},
		{source: &source, kind: TokenBraceLeft, offset: 4},
		{source: &source, kind: TokenBraceRight, offset: 5},
		{source: &source, kind: TokenBracketLeft, offset: 6},
		{source: &source, kind: TokenBracketRight, offset: 7},
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
			t.Error(err)
			break
		}

		got = append(got, token)
		offset = end
	}

	if offset != len(expected) {
		t.Errorf("expected %d got %d", len(expected), offset)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}
