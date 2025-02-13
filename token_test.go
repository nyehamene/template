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
		{kind: TokenColon, offset: 0},
		{kind: TokenEqual, offset: 1},
		{kind: TokenPeriod, offset: 2},
		{kind: TokenSemicolon, offset: 3},
		{kind: TokenBraceLeft, offset: 4},
		{kind: TokenBraceRight, offset: 5},
		{kind: TokenBracketLeft, offset: 6},
		{kind: TokenBracketRight, offset: 7},
		{kind: TokenParLeft, offset: 8},
		{kind: TokenParRight, offset: 9},
	}

	tokenizer := NewTokenizer(source)
	got := []Token{}
	end := 0

	for {
		var err error
		var token Token
		var offset int

		token, offset, err = tokenizer.Tokenize(end)

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
	tokenizer := NewTokenizer(source)

	expected := Token{kind: TokenSpace, offset: 0}
	got, n, err := tokenizer.Tokenize(0)

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
	tokenizer := NewTokenizer(source)

	expected := Token{kind: TokenEOL, offset: 0}

	got, end, err := tokenizer.Tokenize(0)

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
	tokenizer := NewTokenizer(source)

	expected := []Token{
		{kind: TokenIdent, offset: 0},
		{kind: TokenEOL, offset: 3},
		{kind: TokenIdent, offset: 4},
	}

	got := []Token{}
	end := 0
	for {
		token, offset, err := tokenizer.Tokenize(end)

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
	tokenizer := NewTokenizer(source)

	expected := []Token{
		{kind: TokenIdent, offset: 0},
		{kind: TokenSpace, offset: 3},
		{kind: TokenIdent, offset: 4},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var err error
		var offset int
		token, offset, err = tokenizer.Tokenize(end)

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
	source := "package package_tag package_list package_html type templ end"
	//         012345678901234567890123456789012345678901234567890123456789
	tokenizer := NewTokenizer(source)

	expected := []Token{
		{kind: TokenPackage, offset: 0},
		{kind: TokenSpace, offset: 7},
		{kind: TokenTag, offset: 8},
		{kind: TokenSpace, offset: 19},
		{kind: TokenList, offset: 20},
		{kind: TokenSpace, offset: 32},
		{kind: TokenHtml, offset: 33},
		{kind: TokenSpace, offset: 45},
		{kind: TokenType, offset: 46},
		{kind: TokenSpace, offset: 50},
		{kind: TokenTempl, offset: 51},
		{kind: TokenSpace, offset: 56},
		{kind: TokenEnd, offset: 57},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var offset int
		var err error

		token, offset, err = tokenizer.Tokenize(end)

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
	tokenizer := NewTokenizer(source)
	expected := []Token{
		{kind: TokenString, offset: 0},
		{kind: TokenSpace, offset: 2},
		{kind: TokenString, offset: 3},
		{kind: TokenSpace, offset: 6},
		{kind: TokenString, offset: 7},
		{kind: TokenSpace, offset: 16},
		{kind: TokenString, offset: 17},
		{kind: TokenSpace, offset: 22},
		{kind: TokenString, offset: 23},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var offset int
		var err error
		token, offset, err = tokenizer.Tokenize(end)

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
	tokenizer := NewTokenizer(source)
	expected := []Token{
		{kind: TokenEOL, offset: 0},
		{kind: TokenComment, offset: 1},
	}

	got := []Token{}
	end := 0
	for {
		var token Token
		var offset int
		var err error

		token, offset, err = tokenizer.Tokenize(end)

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
	tokenizer := NewTokenizer(source)
	expected := []Token{
		{kind: TokenComment, offset: 0},
	}
	got := []Token{}
	end := 0

	for {
		var token Token
		var offset int
		var err error

		token, offset, err = tokenizer.Tokenize(end)

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

func TestPos(t *testing.T) {
	source := "line1\nline2"
	tokenizer := NewTokenizer(source)
	tokens := []Token{}
	nextTokenOffset := 0

	for {
		var token Token
		var end int
		var err error

		token, end, err = tokenizer.Tokenize(nextTokenOffset)
		if err == EOF {
			break
		}

		if err != nil {
			t.Error(err)
			break
		}

		tokens = append(tokens, token)
		nextTokenOffset = end
	}

	{
		expectedLine := 0
		expectedCol := 0
		token := tokens[0]
		gotLine, gotCol := tokenizer.Pos(token)

		if expectedLine != gotLine {
			t.Errorf("expected line number to be %d got %d", expectedLine, gotLine)
		}
		if expectedCol != gotCol {
			t.Errorf("expected line number to be %d got %d", expectedCol, gotCol)
		}
	}

	{
		expectedLine := 0
		expectedCol := 5
		token := tokens[1]
		gotLine, gotCol := tokenizer.Pos(token)

		if expectedLine != gotLine {
			t.Errorf("expected line number to be %d got %d", expectedLine, gotLine)
		}
		if expectedCol != gotCol {
			t.Errorf("expected line number to be %d got %d", expectedCol, gotCol)
		}
	}

	{
		expectedLine := 1
		expectedCol := 0
		token := tokens[2]
		gotLine, gotCol := tokenizer.Pos(token)

		if expectedLine != gotLine {
			t.Errorf("expected line number to be %d got %d", expectedLine, gotLine)
		}
		if expectedCol != gotCol {
			t.Errorf("expected column number to be %d got %d", expectedCol, gotCol)
		}
	}
}
