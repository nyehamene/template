package tokenizer_test

import (
	"fmt"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNext(t *testing.T) {
	testcases := map[string]token.Token{
		"}": token.New(token.BraceClose, 0),
		"{": token.New(token.BraceOpen, 0),
		"]": token.New(token.BracketClose, 0),
		"[": token.New(token.BracketOpen, 0),
		":": token.New(token.Colon, 0),
		",": token.New(token.Comma, 0),
		".": token.New(token.Dot, 0),
		"=": token.New(token.Eq, 0),
		")": token.New(token.ParenClose, 0),
		"(": token.New(token.ParenOpen, 0),
		";": token.New(token.Semicolon, 0),
	}

	i := 0
	for src, expected := range testcases {
		t.Run(fmt.Sprintf("%d\t%s", i, src), func(t *testing.T) {
			tok := tokenizer.New([]byte(src))
			got := tok.Next()

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
		i += 1
	}
}
