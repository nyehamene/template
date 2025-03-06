package parser

import (
	"temlang/tem/ast"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestCase[T any] map[string]T

func HelperParser[T any, E any](
	t *testing.T,
	testcase TestCase[T],
	parse func(*Parser) (E, bool),
	accept func(E, ast.NamespaceFile) T,
) {

	for src, expected := range testcase {
		t.Run(src, func(t *testing.T) {
			file := ast.NamespaceFile{
				Name: "ns",
				Path: "testing/ns.tem",
				Src:  src,
			}

			p := New(&file)
			pkg, ok := parse(&p)

			if !ok {
				t.Error("parsing failed")
			}

			got := accept(pkg, file)

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}

		})
	}
}

type pkgName struct {
	Ident string
	Type  string
	Name  string
	Templ string
}

func HelperPackage(t *testing.T, testcase TestCase[pkgName]) {
	parse := func(p *Parser) (ast.PackageDecl, bool) {
		return p.ParsePackage()
	}

	accept := func(p ast.PackageDecl, f ast.NamespaceFile) pkgName {
		var got pkgName
		got.Ident = f.GetName(p.Ident)
		got.Type = f.GetName(p.Type)
		got.Name = f.GetName(p.Name)
		got.Templ = f.GetName(p.Templ)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type importName struct {
	Ident string
	Type  string
	Path  string
}

func HelperImport(t *testing.T, testcase TestCase[importName]) {
	parse := func(p *Parser) (ast.ImportDecl, bool) {
		return p.ParseImport()
	}

	accept := func(p ast.ImportDecl, f ast.NamespaceFile) importName {
		var got importName
		got.Ident = f.GetName(p.Ident)
		got.Type = f.GetName(p.Type)
		got.Path = f.GetName(p.Path)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

func TestPackage(t *testing.T) {
	testcase := TestCase[pkgName]{
		`p : package : package("home") templ(tag)`:  {"p", "package", `"home"`, "tag"},
		`p : package : package("home") templ(list)`: {"p", "package", `"home"`, "list"},
		`p : package : package("home") templ(html)`: {"p", "package", `"home"`, "html"},
	}
	HelperPackage(t, testcase)
}

func TestInferedTypePackage(t *testing.T) {
	testcase := TestCase[pkgName]{
		`p :: package("home") templ(tag)`:  {"p", "", `"home"`, "tag"},
		`p :: package("home") templ(list)`: {"p", "", `"home"`, "list"},
		`p :: package("home") templ(html)`: {"p", "", `"home"`, "html"},
	}
	HelperPackage(t, testcase)
}

func TestImport(t *testing.T) {
	testcase := TestCase[importName]{
		`i : import : import("lib/one")`: {"i", "import", `"lib/one"`},
	}
	HelperImport(t, testcase)
}

func TestInferedTypeImport(t *testing.T) {
	testcase := TestCase[importName]{
		`i :: import("lib/one")`: {"i", "", `"lib/one"`},
	}
	HelperImport(t, testcase)
}
