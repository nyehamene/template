package template

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPackage(t *testing.T) {
	var parse Parser
	{
		source := `pkg : package: tag("home");`
		//         0123456789012345678901234
		t := NewTokenizer(source)
		parse = NewParser(t)
	}
	expected := []Ast{
		{AstPackage, AstIdent, AstTag, 0},
	}

	got := []Ast{}
	next := 0
	for {
		var ast Ast
		var end int
		var err error
		ast, end, err = parse.Package(next)

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

	if l := len(*parse.tokenizer.source); l != next {
		t.Errorf("expected %d got %d", l, next)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Error(diff)
	}
}
