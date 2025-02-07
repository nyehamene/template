package template

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tokens
// chars : = ; { } .
// keywords package component tag list html
//          type record templ end
// primary ident string

func TestNext_char(t *testing.T) {
	source := ":=.;{}[]()"
	expected := []Token{
		{source: &source, kind: TokenColon, offset: 0},
		{source: &source, kind: TokenEqual, offset: 1},
		{source: &source, kind: TokenPeriod, offset: 2},
		{source: &source, kind: TokenSemicolon, offset: 3},
		{source: &source, kind: TokenBraceLeft, offset: 4},
		{source: &source, kind: TokenBraceRight, offset: 5},
		{source: &source, kind: TokenBracketLeft, offset: 6},
		{source: &source, kind: TokenBracketRight, offset: 7},
		{source: &source, kind: TokenParLeft, offset: 8},
		{source: &source, kind: TokenParRight, offset: 9},
	}

	got := []Token{}
	end := 0

	for {
		var err error
		var token Token
		var offset int

		token, offset, err = Tokenize(&source, end)

		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		got = append(got, token)
		end = offset
	}

	if end != len(expected) {
		t.Errorf("expected %d got %d", len(expected), end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_space(t *testing.T) {
	source := " \t\r\v\f"
	expected := Token{source: &source, kind: TokenSpace, offset: 0}
	got, n, err := Tokenize(&source, 0)

	if n != len(source) {
		t.Errorf("expected %d got %d", len(source), n)
	}

	if err != nil && err != EOF {
		t.Error(err)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_newline(t *testing.T) {
	source := "\n\n"
	expected := Token{source: &source, kind: TokenEOL, offset: 0}

	got, end, err := Tokenize(&source, 0)

	if err != nil && err != EOF {
		t.Error(err)
	}

	if end != len(source) {
		t.Errorf("expected %d got %d", len(source), end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_newline2(t *testing.T) {
	source := `foo
bar`
	expected := []Token{
		{source: &source, kind: TokenIdent, offset: 0},
		{source: &source, kind: TokenEOL, offset: 3},
		{source: &source, kind: TokenIdent, offset: 4},
	}

	got := []Token{}
	end := 0
	for {
		token, offset, err := Tokenize(&source, end)

		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		got = append(got, token)
		end = offset
	}

	if end != len(source) {
		t.Errorf("expected %d got %d", len(source), end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_ident(t *testing.T) {
	source := "foo bar"
	expected := []Token{
		{source: &source, kind: TokenIdent, offset: 0},
		{source: &source, kind: TokenSpace, offset: 3},
		{source: &source, kind: TokenIdent, offset: 4},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var err error
		var offset int
		token, offset, err = Tokenize(&source, end)

		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		got = append(got, token)
		end = offset
	}

	if end != len(source) {
		t.Errorf("expected %d got %d", len(source), end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_keyword(t *testing.T) {
	source := "package tag list html type templ end"
	//         012345678901234567890123456789012345
	expected := []Token{
		{source: &source, kind: TokenPackage, offset: 0},
		{source: &source, kind: TokenSpace, offset: 7},
		{source: &source, kind: TokenTag, offset: 8},
		{source: &source, kind: TokenSpace, offset: 11},
		{source: &source, kind: TokenList, offset: 12},
		{source: &source, kind: TokenSpace, offset: 16},
		{source: &source, kind: TokenHtml, offset: 17},
		{source: &source, kind: TokenSpace, offset: 21},
		{source: &source, kind: TokenType, offset: 22},
		{source: &source, kind: TokenSpace, offset: 26},
		{source: &source, kind: TokenTempl, offset: 27},
		{source: &source, kind: TokenSpace, offset: 32},
		{source: &source, kind: TokenEnd, offset: 33},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var offset int
		var err error

		token, offset, err = Tokenize(&source, end)

		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		got = append(got, token)
		end = offset
	}

	if end != len(source) {
		t.Errorf("expected %d got %d", len(source), end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_string(t *testing.T) {
	source := `"" ":" "package" "foo" "\n"`
	//         012345678901234567890123 45
	expected := []Token{
		{source: &source, kind: TokenString, offset: 0},
		{source: &source, kind: TokenSpace, offset: 2},
		{source: &source, kind: TokenString, offset: 3},
		{source: &source, kind: TokenSpace, offset: 6},
		{source: &source, kind: TokenString, offset: 7},
		{source: &source, kind: TokenSpace, offset: 16},
		{source: &source, kind: TokenString, offset: 17},
		{source: &source, kind: TokenSpace, offset: 22},
		{source: &source, kind: TokenString, offset: 23},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var offset int
		var err error
		token, offset, err = Tokenize(&source, end)

		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		got = append(got, token)
		end = offset
	}

	if l := len(source); end != l {
		t.Errorf("expected %d got %d", l, end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_comment(t *testing.T) {
	source := `
// line 1
// line 2`
	expected := []Token{
		{source: &source, kind: TokenEOL, offset: 0},
		{source: &source, kind: TokenComment, offset: 1},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var offset int
		var err error

		token, offset, err = Tokenize(&source, end)

		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		got = append(got, token)
		end = offset
	}

	if l := len(source); l != end {
		t.Errorf("expected %d got %d", l, end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}

func TestNext_comment2(t *testing.T) {
	source := `// line 1
	           // line 2`
	expected := []Token{
		{source: &source, kind: TokenComment, offset: 0},
	}
	got := []Token{}
	end := 0

	for {
		var token Token
		var offset int
		var err error

		token, offset, err = Tokenize(&source, end)

		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		got = append(got, token)
		end = offset
	}
	if l := len(source); l != end {
		t.Errorf("expected %d got %d", l, end)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}
