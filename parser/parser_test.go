package parser

import (
	"os"
	"path/filepath"
	"strings"
	"temlang/tem/ast"
	"testing"
)

func TestParseFile(t *testing.T) {
	path := "./testdata"
	list, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range list {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".tem") {
			continue
		}

		filename := filepath.Join(path, d.Name())
		src, err := os.ReadFile(filename)
		if err != nil {
			t.Error(err)
			continue
		}

		_, errs := ParseFile(filename, src)

		if errs.Len() != 0 {
			t.Errorf("ParseFile(%s) parser failed unexpectedly", filename)
		}

		for {
			err, ok := errs.Pop()
			if !ok {
				break
			}
			t.Error(err)
		}
	}
}

func TestDirectivePlacementErrorZero(t *testing.T) {
	t.Skip()
	filename := "directive_error.tem"
	src := `
		p :: package("test")

		User: type: #type type(Person)
	`
	_, errs := ParseFile(filename, []byte(src))
	if errs.Len() == 0 {
		t.Error("parser failed unexpectedly")
	}
}

func TestTopLevelVarDeclaration(t *testing.T) {
	src := `
		p :: package("a")

		name: String
	`

	filename := "test.tem"
	_, err := ParseFile(filename, []byte(src))

	if err.Len() != 0 {
		t.Errorf("ParseFile(%v) failed unexpectedly", filename)
	}
}

func TestAstPrinter(t *testing.T) {
	path := "./testdata"
	list, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range list {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".tem") {
			continue
		}

		filename := filepath.Join(path, d.Name())
		src, err := os.ReadFile(filename)
		if err != nil {
			t.Error(err)
			continue
		}

		file, _ := ParseFile(filename, src)

		printed := ast.PrintSExpr(file)
		if strings.Contains(printed, "ERROR") {
			t.Error("parser failed unexpectedly")
		}
	}
}
