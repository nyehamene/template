package tokenizer_test

import (
	"fmt"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func ident(offset int) token.Token {
	return newToken(token.Ident, offset)
}

func alias(offset int) token.Token {
	return newToken(token.Alias, offset)
}

func import0(offset int) token.Token {
	return newToken(token.Import, offset)
}

func package0(offset int) token.Token {
	return newToken(token.Package, offset)
}

func record(offset int) token.Token {
	return newToken(token.Record, offset)
}

func templ(offset int) token.Token {
	return newToken(token.Templ, offset)
}

func type0(offset int) token.Token {
	return newToken(token.Type, offset)
}

func using(offset int) token.Token {
	return newToken(token.Using, offset)
}

func str(offset int) token.Token {
	return newToken(token.String, offset)
}

func textBlock(offset int) token.Token {
	return newToken(token.TextBlock, offset)
}

func newToken(kind token.Kind, offset int) token.Token {
	return token.New(kind, offset)
}

type TestCase map[string][]token.Token

func HelperRunTestCases(t *testing.T, testcases map[string][]token.Token) {
	i := 0
	for src, expected := range testcases {
		t.Run(fmt.Sprintf("%d\t%s", i, src), func(t *testing.T) {
			gots := []token.Token{}
			tok := tokenizer.New([]byte(src))

			for {
				got := tok.Next()
				if got.Kind == token.EOF {
					break
				}
				gots = append(gots, got)
			}

			if diff := cmp.Diff(expected, gots); diff != "" {
				t.Error(diff)
			}

			if c := tok.ErrorCount(); c > 0 {
				t.Errorf("expected no error but got %d", c)
			}
		})
		i += 1
	}
}

func TestNext(t *testing.T) {
	testcases := TestCase{
		"}": {newToken(token.BraceClose, 0)},
		"{": {newToken(token.BraceOpen, 0)},
		"]": {newToken(token.BracketClose, 0)},
		"[": {newToken(token.BracketOpen, 0)},
		":": {newToken(token.Colon, 0)},
		",": {newToken(token.Comma, 0)},
		".": {newToken(token.Dot, 0)},
		"=": {newToken(token.Eq, 0)},
		")": {newToken(token.ParenClose, 0)},
		"(": {newToken(token.ParenOpen, 0)},
		";": {newToken(token.Semicolon, 0)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextIdent(t *testing.T) {
	testcases := TestCase{
		"_":   {ident(0)},
		"i":   {ident(0)},
		"foo": {ident(0)},
		"a12": {ident(0)},
		"_12": {ident(0)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextKeyword(t *testing.T) {
	testcases := TestCase{
		"alias":   {alias(0)},
		"import":  {import0(0)},
		"package": {package0(0)},
		"record":  {record(0)},
		"templ":   {templ(0)},
		"type":    {type0(0)},
		"using":   {using(0)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextString(t *testing.T) {
	testcases := TestCase{
		`""`:    {str(0)},
		`"i"`:   {str(0)},
		`"foo"`: {str(0)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextTextBlock(t *testing.T) {
	testcases := TestCase{
		`""" line 1
		 """ line 2`: {textBlock(0)},
		`"""
		 """ line 2`: {textBlock(0)},
		`""" line 1
		 """`: {textBlock(0)},
		`"""
		 """`: {textBlock(0)},
	}
	HelperRunTestCases(t, testcases)
}
