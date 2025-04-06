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
		t.Errorf("Parser(%v) succeeded unexpectedly", filename)
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
		t.Errorf("Parser(%v) succeeded unexpectedly", filename)
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
		t.Errorf("Parser(%v) succeeded unexpectedly", filename)
	}
}
