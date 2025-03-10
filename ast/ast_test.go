package ast

import (
	"fmt"
	"temlang/tem/token"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestText(t *testing.T) {
	testcase := []struct {
		src      string
		expected string
		tok      token.Token
	}{
		{"A", "A", tokenDNM(0, 1)},
		{`"A"`, `"A"`, tokenDNM(0, 3)},
		{`package`, `package`, tokenDNM(0, 7)},
	}

	for _, tc := range testcase {
		t.Run(tc.src, func(t *testing.T) {
			file := New(tc.src, "test.tem")
			file.Init()

			txt, ok := file.text(tc.tok)
			if !ok {
				t.Error("should succeed")
			}

			if txt != tc.expected {
				t.Errorf("expected %s got %s", tc.expected, txt)
			}
		})
	}
}

func TestAddText(t *testing.T) {
	// doest not matter
	testcase := []struct {
		src      string
		toks     []token.Token
		tokIndex int
		expected string
	}{
		{"A:Bar", tokenMany(pos{0, 1}, pos{1, 2}, pos{2, 5}), 2, "Bar"},
	}

	for _, tc := range testcase {
		t.Run(tc.src, func(t *testing.T) {
			file := New(tc.src, "test.tem")
			file.Init()

			for _, tok := range tc.toks {
				i := TokenIndex{}
				file.addToken(&i, tok)
				file.addText(&i, tok)
			}

			tok := tc.toks[tc.tokIndex]
			i := TokenIndex{
				token: TokenSlice{
					index: tok.Start(),
					len:   1,
				},
				text: TokenSlice{
					index: tok.Start(),
					len:   1,
				},
			}

			txt := file.getTextOne(i)

			if txt != tc.expected {
				t.Errorf("expected %s got %s", tc.expected, txt)
			}
		})
	}
}

func TestSingle(t *testing.T) {
	var (
		text         = "Foo"
		tokPosStart  = 0
		tokPosEnd    = 3
		tok          = token.New(token.Ident, tokPosStart, tokPosEnd)
		expectedText = "Foo"
		decls        = []singletoken{
			&hasType{},
			&hasName{},
			&PackageDecl{},
			&ImportDecl{},
			&UsingDecl{},
			&TypeDecl{},
		}
	)

	for i, decl := range decls {
		t.Run(fmt.Sprintf("%d_%s", i, text), func(t *testing.T) {

			file := New(text, "test.tem")
			file.Init()
			decl.Set(file, tok)

			if file.tokenLen != 1 {
				t.Errorf("expected file to contain 1 token got %d", file.tokenLen)
			}
			if file.textLen != 1 {
				t.Errorf("expected file to contain 1 text got %d", file.textLen)
			}

			got := file.tokens[0]
			if diff := cmp.Diff(tok, got); diff != "" {
				t.Error(diff)
			}

			txt := decl.Get(*file)

			if txt != expectedText {
				t.Errorf("expected %s got %s", expectedText, txt)
			}
		})
	}
}

func TestMany(t *testing.T) {
	var (
		text           = "A B"
		toks           = []token.Token{token.New(token.Ident, 0, 1), token.New(token.Ident, 2, 3)}
		expectedIndex  = 1
		expectedIdents = []string{"A", "B"}
		decls          = []manytoken{
			&hasIdents{},
		}
	)
	for i, decl := range decls {
		t.Run(fmt.Sprintf("%d_%s", i, text), func(t *testing.T) {

			file := New(text, "test.tem")
			file.Init()
			decl.Set(file, toks)

			if file.tokenLen != 2 {
				t.Errorf("expected file to contain 2 tokens got %d", file.tokenLen)
			}
			if file.textLen != 2 {
				t.Errorf("expected file to contain 2 text got %d", file.textLen)
			}

			expected := toks[expectedIndex]
			got := file.tokens[expectedIndex]
			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}

			idents := decl.Get(*file)
			if diff := cmp.Diff(expectedIdents, idents); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestVar(t *testing.T) {
	var (
		text = "a,b:A"
		toks = []Entry[[]token.Token, token.Token]{
			EntryMany(
				[]token.Token{
					token.New(token.Ident, 0, 1),
					token.New(token.Ident, 2, 3),
				},
				token.New(token.Ident, 4, 5),
			)}
		expectedType   = "A"
		expectedIdents = []string{"a", "b"}
		decls          = []manydecl[VarDecl]{
			&TemplDecl{},
			&RecordDecl{},
		}
	)
	for i, decl := range decls {
		t.Run(fmt.Sprintf("%d_%s", i, text), func(t *testing.T) {

			file := New(text, "test.tem")
			file.Init()
			decl.Set(file, toks)

			if file.tokenLen != 3 {
				t.Errorf("expected file to contain 3 tokens got %d", file.tokenLen)
			}
			if file.textLen != 3 {
				t.Errorf("expected file to contain 3 text got %d", file.textLen)
			}

			vars := decl.Get(*file)
			if l := len(vars); l != 1 {
				t.Fatalf("expected file to contain 1 var got %d", l)
			}

			first := vars[0]
			if ty := first.Type(*file); ty != expectedType {
				t.Errorf("expected type %s got %s", expectedType, ty)
			}

			idents := first.Idents(*file)
			if diff := cmp.Diff(expectedIdents, idents); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAttr(t *testing.T) {
	var (
		text = `a,b="A"`
		toks = []Entry[[]token.Token, token.Token]{
			EntryMany(
				[]token.Token{
					token.New(token.Ident, 0, 1),
					token.New(token.Ident, 2, 3),
				},
				token.New(token.String, 4, 7),
			)}
		expectedValue  = `"A"`
		expectedIdents = []string{"a", "b"}
		decls          = []manydecl[AttrDecl]{
			&TagDecl{},
		}
	)
	for i, decl := range decls {
		t.Run(fmt.Sprintf("%d_%s", i, text), func(t *testing.T) {

			file := New(text, "test.tem")
			file.Init()
			decl.Set(file, toks)

			if file.tokenLen != 3 {
				t.Errorf("expected file to contain 3 tokens got %d", file.tokenLen)
			}
			if file.textLen != 3 {
				t.Errorf("expected file to contain 3 text got %d", file.textLen)
			}

			vars := decl.Get(*file)
			if l := len(vars); l != 1 {
				t.Fatalf("expected file to contain 1 var got %d", l)
			}

			first := vars[0]
			if ty := first.Value(*file); ty != expectedValue {
				t.Errorf("expected value %s got %s", expectedValue, ty)
			}

			idents := first.Idents(*file)
			if diff := cmp.Diff(expectedIdents, idents); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestContent(t *testing.T) {
	var (
		text = `"line 1" """ line 2`
		toks = []token.Token{
			token.New(token.String, 0, 8),
			token.New(token.String, 9, 19),
		}
		expected = "\"line 1\"\n\"\"\" line 2"
		decl     = DocDecl{}
	)

	file := New(text, "test.tem")
	file.Init()
	decl.SetContent(file, toks...)

	if file.tokenLen != 2 {
		t.Errorf("expected file to contain 3 tokens got %d", file.tokenLen)
	}
	if file.textLen != 2 {
		t.Errorf("expected file to contain 3 text got %d", file.textLen)
	}

	content := decl.Content(*file)
	if diff := cmp.Diff(expected, content); diff != "" {
		t.Error(diff)
	}
}
