package parser

import "testing"

func TestPackageAfterDeclarationError(t *testing.T) {
	src := `
	u :: using(foo)

	p :: package("a")

	i :: import("b")
	`

	filename := "test.tem"
	_, err := ParseFile(filename, []byte(src))

	if err.Len() == 0 {
		t.Errorf("ParseFile(%v) succeeded unexpectedly", filename)
	}
}

func TestImportAfterDeclarationError(t *testing.T) {
	src := `
	p :: package("a")

	u :: using(foo)

	i :: import("b")
	`

	filename := "test.tem"
	_, err := ParseFile(filename, []byte(src))

	if err.Len() == 0 {
		t.Errorf("ParseFile(%v) succeeded unexpectedly", filename)
	}
}

func TestUsingAfterDeclarationError(t *testing.T) {
	src := `
	p :: package("a")

	i :: import("b")

	t :: type(foo)
	
	u :: using(foo)
	`

	filename := "test.tem"
	_, err := ParseFile(filename, []byte(src))

	if err.Len() == 0 {
		t.Errorf("ParseFile(%v) succeeded unexpectedly", filename)
	}
}

func TestDocumentationWithNoTarget(t *testing.T) {
	src := `
		p :: package("a")

		d : "a documentation"
	`

	filename := "test.tem"
	_, err := ParseFile(filename, []byte(src))

	if err.Len() == 0 {
		t.Errorf("ParseFile(%v) succeeded unexpectedly", filename)
	}
}

func TestTagWithNoTarget(t *testing.T) {
	src := `
		p :: package("a")

		a : { key = "value" }
	`

	filename := "test.tem"
	_, err := ParseFile(filename, []byte(src))

	if err.Len() == 0 {
		t.Errorf("ParseFile(%v) succeeded unexpectedly", filename)
	}
}

func TestDirectivePlacementErrorOne(t *testing.T) {
	t.Skip()
	filename := "directive_error.tem"
	src := `
		p := package("test")

		User: type #type: type(Person)
	`
	_, errs := ParseFile(filename, []byte(src))
	if errs.Len() == 0 {
		t.Error("parser succeeded unexpectedly")
	}
}
