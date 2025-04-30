package tokenizer_test

import (
	"fmt"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
	tu "temlang/tem/tokenizer/internal"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestCase map[string][]token.Token

func HelperRunTestCases(
	t *testing.T,
	testcases map[string][]token.Token,
	opts ...tokenizer.Option,
) {
	for src, expected := range testcases {
		t.Run(fmt.Sprintf("_%s", src), func(t *testing.T) {
			gots := []token.Token{}
			tok := tokenizer.New("", []byte(src), opts...)

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
		"}": {tu.NewSymbol(token.BraceClose, 0)},
		"{": {tu.NewSymbol(token.BraceOpen, 0)},
		"]": {tu.NewSymbol(token.BracketClose, 0)},
		"[": {tu.NewSymbol(token.BracketOpen, 0)},
		":": {tu.NewSymbol(token.Colon, 0)},
		",": {tu.NewSymbol(token.Comma, 0)},
		".": {tu.NewSymbol(token.Dot, 0)},
		"=": {tu.NewSymbol(token.Eq, 0)},
		")": {tu.NewSymbol(token.ParenClose, 0)},
		"(": {tu.NewSymbol(token.ParenOpen, 0)},
		";": {tu.NewSymbol(token.Semicolon, 0)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextIdent(t *testing.T) {
	testcases := TestCase{
		"_":   {tu.NewIdent(0, 1)},
		"i":   {tu.NewIdent(0, 1)},
		"foo": {tu.NewIdent(0, 3)},
		"a12": {tu.NewIdent(0, 3)},
		"_12": {tu.NewIdent(0, 3)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextKeyword(t *testing.T) {
	testcases := TestCase{
		"import ":  {tu.NewImport(0, 6)},
		"package ": {tu.NewPackage(0, 7)},
		"record ":  {tu.NewRecord(0, 6)},
		"templ ":   {tu.NewTempl(0, 5)},
		"type ":    {tu.NewType(0, 4)},
		"using ":   {tu.NewUsing(0, 5)},
		"#tag ":    {tu.NewToken(token.Directive, 0, 4)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextString(t *testing.T) {
	testcases := TestCase{
		`""`:    {tu.NewStr(0, 2)},
		`"i"`:   {tu.NewStr(0, 3)},
		`"foo"`: {tu.NewStr(0, 5)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextTextBlock(t *testing.T) {
	testcases := TestCase{
		`--`:        {tu.NewTextBlock(0, 2)},
		`-- line 1`: {tu.NewTextBlock(0, 9)},
		`-- line 1
		 -- line 2`: {tu.NewTextBlock(0, 9), tu.NewTextBlock(13, 22)},
	}
	HelperRunTestCases(t, testcases, tokenizer.NoSemicolonInsertion())
}

func TestNextComment(t *testing.T) {
	testcases := TestCase{
		`//`:                  {tu.NewComment(0, 2)},
		`// one line comment`: {tu.NewComment(0, 19)},
		`// line 1
		 // line 2`: {tu.NewComment(0, 9), tu.NewComment(13, 22)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextInsertSemicolon(t *testing.T) {
	testcases := TestCase{
		")\n":     {tu.NewSymbol(token.ParenClose, 0), tu.NewEOL(1)},
		"]\n":     {tu.NewSymbol(token.BracketClose, 0), tu.NewEOL(1)},
		"}\n":     {tu.NewSymbol(token.BraceClose, 0), tu.NewEOL(1)},
		"ident\n": {tu.NewIdent(0, 5), tu.NewEOL(5)},
		`""
			`: {tu.NewStr(0, 2), tu.NewEOL(2)},
		`"abc"
			`: {tu.NewStr(0, 5), tu.NewEOL(5)},
		`-- line 1
		 `: {tu.NewTextBlock(0, 9), tu.NewEOL(9)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextInsertSemicolonEOF(t *testing.T) {
	// insert semicolon at eof
	testcases := TestCase{
		")":     {tu.NewSymbol(token.ParenClose, 0), tu.NewEOL(1)},
		"]":     {tu.NewSymbol(token.BracketClose, 0), tu.NewEOL(1)},
		"}":     {tu.NewSymbol(token.BraceClose, 0), tu.NewEOL(1)},
		"ident": {tu.NewIdent(0, 5), tu.NewEOL(5)},
		`""`:    {tu.NewStr(0, 2), tu.NewEOL(2)},
		`"abc"`: {tu.NewStr(0, 5), tu.NewEOL(5)},
		`-- line 1
		 `: {tu.NewTextBlock(0, 9), tu.NewEOL(9)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextSemicolonBeforeTrailingComment(t *testing.T) {
	testcases := TestCase{
		"ident // a comment": {tu.NewIdent(0, 5), tu.NewEOL(6), tu.NewComment(6, 18)},
		") // a comment":     {tu.NewSymbol(token.ParenClose, 0), tu.NewEOL(2), tu.NewComment(2, 14)},
	}
	HelperRunTestCases(t, testcases)
}

func TestNextInsertSemicolonAndNewline(t *testing.T) {
	testcases := TestCase{
		"ident\n": {tu.NewIdent(0, 5), tu.NewEOL(5), tu.NewSymbol(token.EOL, 5)},
	}
	for src, expected := range testcases {
		tok := tokenizer.New("", []byte(src))
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
		"\n\n":  {tu.NewSymbol(token.EOL, 0), tu.NewSymbol(token.EOL, 1)},
		"\n \n": {tu.NewSymbol(token.EOL, 0), tu.NewSymbol(token.EOL, 2)},
	}
	for src, expected := range testcases {
		tok := tokenizer.New("", []byte(src))
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
