package template

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testCase func(int) (int, error)

func ident(offset int) Token {
	return Token{TokenIdent, offset}
}

func TestParse_packageDef(t *testing.T) {
	var testcases []Parser
	var wants [][]Def
	{
		source := `pkg :package :package_tag("home");`
		//         0123456789012345678901234567890
		t := NewTokenizer(source)
		testcase := NewParser(&t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Def{
			{DefTagPackage, ident(0)},
		})
	}
	{
		source := `pkg :: package_tag("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(&t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Def{
			{DefTagPackage, ident(0)},
		})
	}
	{
		source := `pkg :: package_list("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(&t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Def{
			{DefListPackage, ident(0)},
		})
	}
	{
		source := `pkg :: package_html("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(&t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Def{
			{DefHtmlPackage, ident(0)},
		})
	}
	for i, p := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			next := 0
			for {
				var end int
				var err error
				end, err = p.defPackage(next)

				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				next = end
			}
			got := p.ast

			if l := len(*p.tokenizer.source); l != next {
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
	var wants [][]Def
	{
		source := "f : type : record {};"
		want := []Def{
			{DefRecord, ident(0)},
		}
		wants = append(wants, want)

		t := NewTokenizer(source)
		testcases = append(testcases, NewParser(&t))
	}
	{
		source := "f :: record {};"
		want := []Def{
			{DefRecord, ident(0)},
		}
		wants = append(wants, want)

		t := NewTokenizer(source)

		testcases = append(testcases, NewParser(&t))
	}
	{
		source := "f :: record { a: A };"
		want := []Def{
			{DefRecord, ident(0)},
		}
		wants = append(wants, want)

		t := NewTokenizer(source)

		testcases = append(testcases, NewParser(&t))
	}
	{
		source := `f :: record { a: A;
		b: B;
		};`
		want := []Def{
			{DefRecord, ident(0)},
		}
		wants = append(wants, want)

		t := NewTokenizer(source)

		testcases = append(testcases, NewParser(&t))
	}

	for i, p := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			next := 0
			for {
				var end int
				var err error
				end, err = p.defTypeOrTempl(next)
				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				next = end
			}
			got := p.ast

			if l := len(*p.tokenizer.source); l != next {
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
	var wants [][]Def
	{
		source := "A : type : alias Foo;"
		//         012345678901234567890
		want := []Def{{DefAlias, ident(0)}}
		wants = append(wants, want)

		t := NewTokenizer(source)

		testcases = append(testcases, NewParser(&t))
	}
	for i, p := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			next := 0
			for {
				var end int
				var err error
				end, err = p.defTypeOrTempl(next)
				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				next = end
			}
			got := p.ast

			if l := len(*p.tokenizer.source); l != next {
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
	var wants [][]Def
	{
		source := "render : templ : (User) {};"
		//         012345678901234567890123456
		want := []Def{{DefTemplate, ident(0)}}
		wants = append(wants, want)

		t := NewTokenizer(source)

		testcases = append(testcases, NewParser(&t))
	}

	for i, p := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			next := 0
			for {
				var end int
				var err error
				end, err = p.defTypeOrTempl(next)
				if err == EOF {
					break
				}

				if err != nil {
					t.Error(err)
					break
				}

				next = end
			}

			got := p.ast

			if l := len(*p.tokenizer.source); l != next {
				t.Errorf("expected %d got %d", l, next)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_comment(t *testing.T) {
	var testcases []func(int) (int, error)
	var wants []int
	{
		source := "// single line comment"

		t := NewTokenizer(source)

		p := NewParser(&t)
		testcases = append(testcases, p.defPackage)
		testcases = append(testcases, p.defTypeOrTempl)
		// Due to the way the parse method back track on error
		// even on EOF when the source contains only space or comment
		// the offset will always be zero
		wants = append(wants, 0)
		wants = append(wants, 0)
	}
	{
		source := `// line 1
		// line 2`

		t := NewTokenizer(source)

		p := NewParser(&t)
		testcases = append(testcases, p.defPackage)
		testcases = append(testcases, p.defTypeOrTempl)
		wants = append(wants, 0)
		wants = append(wants, 0)
	}
	{
		source := `// package comment
		p :: package_tag("home");`

		t := NewTokenizer(source)

		p := NewParser(&t)
		testcases = append(testcases, p.defPackage)
		wants = append(wants, len(source))
	}
	{
		source := `
		// line 1
		// line 2
		B : type : alias A;`

		t := NewTokenizer(source)

		p := NewParser(&t)
		testcases = append(testcases, p.defTypeOrTempl)
		wants = append(wants, len(source))
	}

	testFunc := func(start, end int, parseAt testCase) func(*testing.T) {
		return func(t *testing.T) {
			next := start
			for {
				n, err := parseAt(next)
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
	var testcases []Parser
	var wants [][]Def
	{
		source := `
		A : "single line";
		A : type : record {};`
		want := []Def{
			{DefDocline, ident(3)},
			{DefRecord, ident(24)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
	}
	{
		source := `
		A : """
			""" line 1
			""" line 2
			;
		A : type : record {};`
		want := []Def{
			{DefDocblock, ident(3)},
			{DefRecord, ident(46)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
	}

	for i, parser := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]

			next, err := parser.docDef(0)

			got := parser.ast

			if err != nil {
				t.Error(err)
			}

			if l := len(*parser.tokenizer.source); next != l {
				t.Errorf("expected %d got %d", l, next)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse_import(t *testing.T) {
	var testcases []Parser
	var wants [][]Def
	var ends []int
	{
		source := `p : import : import("home/pkg");`
		want := []Def{
			{DefImport, ident(0)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
		ends = append(ends, len(source))
	}
	{
		source := `p :: import("home/pkg");`
		want := []Def{
			{DefImport, ident(0)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
		ends = append(ends, len(source))
	}

	for i, p := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			end := ends[i]
			next := 0
			for {
				n, err := p.defImport(next)
				if err == EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				next = n
			}
			got := p.ast

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
	var testcases []Parser
	var wants [][]Def
	var ends []int
	{
		source := `A : import : using p;`
		want := []Def{
			{DefUsing, ident(0)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
		ends = append(ends, len(source))
	}
	{
		source := `A :: using p;`
		want := []Def{
			{DefUsing, ident(0)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
		ends = append(ends, len(source))
	}

	for i, p := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			end := ends[i]
			next := 0
			for {
				n, err := p.defUsing(next)
				if err == EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				next = n
			}
			got := p.ast

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
	var testcases []Parser
	var wants [][]Def
	var ends []int
	{
		source := `A : { k = "foo" };`
		want := []Def{
			{DefMetatable, ident(0)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
		ends = append(ends, len(source))
	}
	{
		source := `A : { a = "foo", b = "bar" };`
		want := []Def{
			{DefMetatable, ident(0)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
		ends = append(ends, len(source))
	}
	{
		source := `A : { a = "foo", b = "bar", };`
		want := []Def{
			{DefMetatable, ident(0)},
		}

		t := NewTokenizer(source)

		p := NewParser(&t)
		wants = append(wants, want)
		testcases = append(testcases, p)
		ends = append(ends, len(source))
	}

	for i, p := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			expected := wants[i]
			end := ends[i]
			next := 0
			for {
				n, err := p.metatable(next)
				if err == EOF {
					break
				}
				if err != nil {
					t.Error(err)
					break
				}
				next = n
			}
			got := p.ast

			if end != next {
				t.Errorf("expected %d got %d", end, next)
			}
			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestParse(t *testing.T) {
	f, err := os.Open("demo/record.template")
	if err != nil {
		t.Fatal(err)
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	source := string(buf)

	tz := NewTokenizer(source)
	p := NewParser(&tz)

	asts, err := p.Parse(0)
	if err != nil {
		t.Error(err)
	}

	t.Log(asts)
}
