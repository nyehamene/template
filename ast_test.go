package template

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPackage(t *testing.T) {
	var testcases []Parser
	var wants [][]Ast
	{
		source := `pkg :package :package_tag("home");`
		//         0123456789012345678901234567890
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{{AstPackage, AstIdent, AstTag, 0}})
	}
	{
		source := `pkg :: package_tag("home");`
		//         012345678901234567890123
		t := NewTokenizer(source)
		testcase := NewParser(t)
		testcases = append(testcases, testcase)
		wants = append(wants, []Ast{{AstPackage, AstIdent, AstTag, 0}})
	}
	for i, tc := range testcases {
		t.Run(fmt.Sprintf("(%d)", i), func(t *testing.T) {
			i := i
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
