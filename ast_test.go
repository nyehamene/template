package template

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCase func(int) (Ast, int, error)

func ident(offset int) Token {
	return Token{TokenIdent, offset}
}

func TestParse_packageDef(t *testing.T) {
	var testcases []Parser
	var wants [][]Ast
	{
		source := `pkg :package :package_tag("home");`
		//         0123456789012345678901234567890
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{
			{AstTagTemplPackage, ident(0)},
		})
	}
	{
		source := `pkg :: package_tag("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{
			{AstTagTemplPackage, ident(0)},
		})
	}
	{
		source := `pkg :: package_list("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{
			{AstListTemplPackage, ident(0)},
		})
	}
	{
		source := `pkg :: package_html("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{
			{AstHtmlTemplPackage, ident(0)},
		})
	}
	for i, tc := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			next := 0
			for {
				var ast Ast
				var end int
				var err error
				ast, end, err = tc.parsePackage(next)

				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				got = append(got, ast)
				next = end
			}

			if l := len(*tc.tokenizer.source); l != next {
				t.Errorf("expected %d got %d", l, next)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_typeDef(t *testing.T) {
	var testcases []Parser
	var wants [][]Ast
	{
		source := "f : type : record {};"
		want := []Ast{
			{AstRecord, ident(0)},
		}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}
	{
		source := "f :: record {};"
		want := []Ast{
			{AstRecord, ident(0)},
		}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}
	{
		source := "f :: record { a: A; };"
		want := []Ast{
			{AstRecord, ident(0)},
		}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}
	{
		source := `f :: record { a: A;
		b: B;
		};`
		want := []Ast{
			{AstRecord, ident(0)},
		}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			next := 0
			for {
				var ast Ast
				var end int
				var err error
				ast, end, err = tc.parseDef(next)
				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				got = append(got, ast)
				next = end
			}

			if l := len(*tc.tokenizer.source); l != next {
				t.Errorf("expected %d got %d", l, next)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_typeAliasDef(t *testing.T) {
	var testcases []Parser
	var wants [][]Ast
	{
		source := "A : type : alias Foo;"
		//         012345678901234567890
		want := []Ast{{AstAlias, ident(0)}}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}
	for i, tc := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			next := 0
			for {
				var ast Ast
				var end int
				var err error
				ast, end, err = tc.parseDef(next)
				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				got = append(got, ast)
				next = end
			}

			if l := len(*tc.tokenizer.source); l != next {
				t.Errorf("expected %d got %d", l, next)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_templDef(t *testing.T) {
	var testcases []Parser
	var wants [][]Ast
	{
		source := "render : templ : (User) {};"
		//         012345678901234567890123456
		want := []Ast{{AstTemplate, ident(0)}}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}

	for i, tc := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			next := 0
			for {
				var ast Ast
				var end int
				var err error
				ast, end, err = tc.parseDef(next)
				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				got = append(got, ast)
				next = end
			}

			if l := len(*tc.tokenizer.source); l != next {
				t.Errorf("expected %d got %d", l, next)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_comment(t *testing.T) {
	var testcases []func(int) (Ast, int, error)
	var wants []int
	{
		source := "// single line comment"
		p := NewParser(NewTokenizer(source))
		testcases = append(testcases, p.parsePackage)
		testcases = append(testcases, p.parseDef)
		// Due to the way the parse method back track on error
		// even on EOF when the source contains only space or comment
		// the offset will always be zero
		wants = append(wants, 0)
		wants = append(wants, 0)
	}
	{
		source := `// line 1
		// line 2`
		p := NewParser(NewTokenizer(source))
		testcases = append(testcases, p.parsePackage)
		testcases = append(testcases, p.parseDef)
		wants = append(wants, 0)
		wants = append(wants, 0)
	}
	{
		source := `// package comment
		p :: package_tag("home");`

		p := NewParser(NewTokenizer(source))
		testcases = append(testcases, p.parsePackage)
		wants = append(wants, len(source))
	}
	{
		source := `
		// line 1
		// line 2
		B : type : alias A;`

		p := NewParser(NewTokenizer(source))
		testcases = append(testcases, p.parseDef)
		wants = append(wants, len(source))
	}

	testFunc := func(start, end int, parseAt testCase) func(*testing.T) {
		return func(t *testing.T) {
			next := start
			for {
				_, n, err := parseAt(next)
				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}
				next = n
			}

			if end != next {
				t.Errorf("expected %d got %d", end, next)
			}
		}
	}

	// The parse should skip comments withcout producing any error
	for i, tc := range testcases {
		expected := wants[i]
		testName := fmt.Sprintf("(%d)", i)
		t.Run(testName, testFunc(0, expected, tc))
	}
}

func TestParse_doc(t *testing.T) {
	var testcases []func(int) (Ast, int, error)
	var wants [][]Ast
	{
		source := `
		A : "single line";
		A : type : record {};
		`
		want := []Ast{
			{AstDocline, ident(3)},
			{AstRecord, ident(24)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.parseDef)
	}
	{
		source := `
		A : """
			""" line 1
			""" line 2
			;
		A : type : record {};
		`
		want := []Ast{
			{AstDocblock, ident(3)},
			{AstRecord, ident(46)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.parseDef)
	}

	for i, parseAt := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			next := 0
			for {
				token, n, err := parseAt(next)
				if err == EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				got = append(got, token)
				next = n
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_import(t *testing.T) {
	var testcases []testCase
	var wants [][]Ast
	var ends []int
	{
		source := `p : import : import("home/pkg");`
		want := []Ast{
			{AstImport, ident(0)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.parseImport)
		ends = append(ends, len(source))
	}
	{
		source := `p :: import("home/pkg");`
		want := []Ast{
			{AstImport, ident(0)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.parseImport)
		ends = append(ends, len(source))
	}

	for i, parseAt := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			end := ends[i]
			next := 0
			for {
				token, n, err := parseAt(next)
				if err == EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				got = append(got, token)
				next = n
			}

			if end != next {
				t.Errorf("expected %d got %d", end, next)
			}
			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_using(t *testing.T) {
	var testcases []testCase
	var wants [][]Ast
	var ends []int
	{
		source := `A : import : using p;`
		want := []Ast{
			{AstUsing, ident(0)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.parseUsing)
		ends = append(ends, len(source))
	}
	{
		source := `A :: using p;`
		want := []Ast{
			{AstUsing, ident(0)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.parseUsing)
		ends = append(ends, len(source))
	}

	for i, parseAt := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			end := ends[i]
			next := 0
			for {
				token, n, err := parseAt(next)
				if err == EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				got = append(got, token)
				next = n
			}

			if end != next {
				t.Errorf("expected %d got %d", end, next)
			}
			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_metatable(t *testing.T) {
	var testcases []testCase
	var wants [][]Ast
	var ends []int
	{
		source := `A : { k = "foo" };`
		want := []Ast{
			{AstMetatable, ident(0)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.metatable)
		ends = append(ends, len(source))
	}
	{
		source := `A : { a = "foo", b = "bar" };`
		want := []Ast{
			{AstMetatable, ident(0)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.metatable)
		ends = append(ends, len(source))
	}
	{
		source := `A : { a = "foo", b = "bar", };`
		want := []Ast{
			{AstMetatable, ident(0)},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.metatable)
		ends = append(ends, len(source))
	}

	for i, parseAt := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			got := []Ast{}
			end := ends[i]
			next := 0
			for {
				token, n, err := parseAt(next)
				if err == EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				got = append(got, token)
				next = n
			}

			if end != next {
				t.Errorf("expected %d got %d", end, next)
			}
			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
