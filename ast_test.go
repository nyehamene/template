package template

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type parserTestCase func(int) (Ast, int, error)

func TestParse_packageDef(t *testing.T) {
	var testcases []Parser
	var wants [][]Ast
	{
		source := `pkg :package :package_tag("home");`
		//         0123456789012345678901234567890
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{{AstPackage, AstIdent, AstTagPackage, 0}})
	}
	{
		source := `pkg :: package_tag("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{{AstPackage, AstIdent, AstTagPackage, 0}})
	}
	{
		source := `pkg :: package_list("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{{AstPackage, AstIdent, AstListPackage, 0}})
	}
	{
		source := `pkg :: package_html("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{{AstPackage, AstIdent, AstHtmlPackage, 0}})
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
				ast, end, err = tc.Package(next)

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
		want := []Ast{{AstTypeDef, AstTypeIdent, AstRecordDef, 0}}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}
	{
		source := "f :: record {};"
		want := []Ast{{AstTypeDef, AstTypeIdent, AstRecordDef, 0}}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}
	{
		source := "f :: record { a: A; };"
		want := []Ast{{AstTypeDef, AstTypeIdent, AstRecordDef, 0}}
		wants = append(wants, want)
		testcases = append(testcases, NewParser(NewTokenizer(source)))
	}
	{
		source := `f :: record { a: A;
		b: B;
		};`
		want := []Ast{{AstTypeDef, AstTypeIdent, AstRecordDef, 0}}
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
				ast, end, err = tc.Def(next)
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
		want := []Ast{{AstTypeDef, AstTypeIdent, AstAliasDef, 0}}
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
				ast, end, err = tc.Def(next)
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
		want := []Ast{{AstTemplateDef, AstIdent, AstTemplateBody, 0}}
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
				ast, end, err = tc.Def(next)
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
		testcases = append(testcases, p.Package)
		testcases = append(testcases, p.Def)
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
		testcases = append(testcases, p.Package)
		testcases = append(testcases, p.Def)
		wants = append(wants, 0)
		wants = append(wants, 0)
	}
	{
		source := `// package comment
		p :: package_tag("home");`

		p := NewParser(NewTokenizer(source))
		testcases = append(testcases, p.Package)
		wants = append(wants, len(source))
	}
	{
		source := `
		// line 1
		// line 2
		B : type : alias A;`

		p := NewParser(NewTokenizer(source))
		testcases = append(testcases, p.Def)
		wants = append(wants, len(source))
	}

	testFunc := func(start, end int, parseAt func(int) (Ast, int, error)) func(*testing.T) {
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
			{AstDocDef, AstIdent, AstDocline, 3},
			{AstTypeDef, AstTypeIdent, AstRecordDef, 24},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.Def)
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
			{AstDocDef, AstIdent, AstDocblock, 3},
			{AstTypeDef, AstTypeIdent, AstRecordDef, 46},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.Def)
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
	var testcases []parserTestCase
	var wants [][]Ast
	var ends []int
	{
		source := `p : import : import("home/pkg");`
		want := []Ast{
			{AstImport, AstIdent, AstImportPackage, 0},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.Import)
		ends = append(ends, len(source))
	}
	{
		source := `p :: import("home/pkg");`
		want := []Ast{
			{AstImport, AstIdent, AstImportPackage, 0},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.Import)
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
	var testcases []parserTestCase
	var wants [][]Ast
	var ends []int
	{
		source := `A : import : using p;`
		want := []Ast{
			{AstUsing, AstTypeIdent, AstIdent, 0},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.Using)
		ends = append(ends, len(source))
	}
	{
		source := `A :: using p;`
		want := []Ast{
			{AstUsing, AstTypeIdent, AstIdent, 0},
		}
		p := NewParser(NewTokenizer(source))
		wants = append(wants, want)
		testcases = append(testcases, p.Using)
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
