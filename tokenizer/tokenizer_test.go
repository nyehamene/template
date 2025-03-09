package tokenizer_test

import (
	"fmt"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func ident(offset, end int) token.Token {
	return newToken(token.Ident, offset, end)
}

func alias(offset, end int) token.Token {
	return newToken(token.Alias, offset, end)
}

func import0(offset, end int) token.Token {
	return newToken(token.Import, offset, end)
}

func package0(offset, end int) token.Token {
	return newToken(token.Package, offset, end)
}

func record(offset, end int) token.Token {
	return newToken(token.Record, offset, end)
}

func templ(offset, end int) token.Token {
	return newToken(token.Templ, offset, end)
}

func type0(offset, end int) token.Token {
	return newToken(token.Type, offset, end)
}

func using(offset, end int) token.Token {
	return newToken(token.Using, offset, end)
}

func str(offset, end int) token.Token {
	return newToken(token.String, offset, end)
}

func textBlock(offset, end int) token.Token {
	return newToken(token.TextBlock, offset, end)
}

func comment(offset, end int) token.Token {
	return newToken(token.Comment, offset, end)
}

// eol automatically inserted semicolon
func eol(offset int) token.Token {
	return newToken(token.Semicolon, offset, offset)
}

func symbol(kind token.Kind, offset int) token.Token {
	return newToken(kind, offset, offset+1)
}

func newToken(kind token.Kind, offset, end int) token.Token {
	return token.New(kind, offset, end)
}

type TestCase map[string][]token.Token

func HelperRunTestCases(
	t *testing.T,
	testcases map[string][]token.Token,
	opts ...tokenizer.Option,
) {
	for src, expected := range testcases {
		t.Run(fmt.Sprintf("_%s", src), func(t *testing.T) {
			gots := []token.Token{}
			tok := tokenizer.New([]byte(src), opts...)

			for {
				got := tok.Next()
				if got.Kind() == token.EOF {
					break
				}
				if got.Kind() == token.EOL {
					continue
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
	}
}

func TestNext(t *testing.T) {
	testcases := TestCase{
		"}": {symbol(token.BraceClose, 0)},
		"{": {symbol(token.BraceOpen, 0)},
		"]": {symbol(token.BracketClose, 0)},
		"[": {symbol(token.BracketOpen, 0)},
		":": {symbol(token.Colon, 0)},
		",": {symbol(token.Comma, 0)},
		".": {symbol(token.Dot, 0)},
		"=": {symbol(token.Eq, 0)},
		")": {symbol(token.ParenClose, 0)},
		"(": {symbol(token.ParenOpen, 0)},
		";": {symbol(token.Semicolon, 0)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextIdent(t *testing.T) {
	testcases := TestCase{
		"_":   {ident(0, 1)},
		"i":   {ident(0, 1)},
		"foo": {ident(0, 3)},
		"a12": {ident(0, 3)},
		"_12": {ident(0, 3)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextKeyword(t *testing.T) {
	testcases := TestCase{
		"alias ":   {alias(0, 5)},
		"import ":  {import0(0, 6)},
		"package ": {package0(0, 7)},
		"record ":  {record(0, 6)},
		"templ ":   {templ(0, 5)},
		"type ":    {type0(0, 4)},
		"using ":   {using(0, 5)},
		"tag ":     {newToken(token.Tag, 0, 3)},
		"list ":    {newToken(token.List, 0, 4)},
		"html ":    {newToken(token.Html, 0, 4)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextString(t *testing.T) {
	testcases := TestCase{
		`""`:    {str(0, 2)},
		`"i"`:   {str(0, 3)},
		`"foo"`: {str(0, 5)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextTextBlock(t *testing.T) {
	testcases := TestCase{
		`"""`:        {textBlock(0, 3)},
		`""" line 1`: {textBlock(0, 10)},
		`""" line 1
		 """ line 2`: {textBlock(0, 10), textBlock(14, 24)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextComment(t *testing.T) {
	testcases := TestCase{
		`//`:                  {comment(0, 2)},
		`// one line comment`: {comment(0, 19)},
		`// line 1
		 // line 2`: {comment(0, 9), comment(13, 22)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextInsertSemicolon(t *testing.T) {
	testcases := TestCase{
		")\n":     {symbol(token.ParenClose, 0), eol(1)},
		"]\n":     {symbol(token.BracketClose, 0), eol(1)},
		"}\n":     {symbol(token.BraceClose, 0), eol(1)},
		"ident\n": {ident(0, 5), eol(5)},
		`""
			`: {str(0, 2), eol(2)},
		`"abc"
			`: {str(0, 5), eol(5)},
		`""" line 1
		 `: {textBlock(0, 10), eol(10)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextInsertSemicolonEOF(t *testing.T) {
	// insert semicolon at eof
	testcases := TestCase{
		")":     {symbol(token.ParenClose, 0), eol(1)},
		"]":     {symbol(token.BracketClose, 0), eol(1)},
		"}":     {symbol(token.BraceClose, 0), eol(1)},
		"ident": {ident(0, 5), eol(5)},
		`""`:    {str(0, 2), eol(2)},
		`"abc"`: {str(0, 5), eol(5)},
		`""" line 1
		 `: {textBlock(0, 10), eol(10)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextSemicolonBeforeTrailingComment(t *testing.T) {
	testcases := TestCase{
		"ident // a comment": {ident(0, 5), eol(6), comment(6, 18)},
		") // a comment":     {symbol(token.ParenClose, 0), eol(2), comment(2, 14)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextInsertSemicolonAndNewline(t *testing.T) {
	testcases := TestCase{
		"ident\n": {ident(0, 5), eol(5), symbol(token.EOL, 5)},
	}
	for src, expected := range testcases {
		tok := tokenizer.New([]byte(src))
		gots := []token.Token{}
		for {
			got := tok.Next()
			if got.Kind() == token.EOF {
				break
			}
			gots = append(gots, got)
		}
		if diff := cmp.Diff(expected, gots); diff != "" {
			t.Error(diff)
		}
	}
}

func TestNextNewline(t *testing.T) {
	testcases := TestCase{
		"\n\n":  {symbol(token.EOL, 0), symbol(token.EOL, 1)},
		"\n \n": {symbol(token.EOL, 0), symbol(token.EOL, 2)},
	}
	for src, expected := range testcases {
		tok := tokenizer.New([]byte(src))
		gots := []token.Token{}
		for {
			got := tok.Next()
			if got.Kind() == token.EOF {
				break
			}
			gots = append(gots, got)
		}
		if diff := cmp.Diff(expected, gots); diff != "" {
			t.Error(diff)
		}
	}
}
