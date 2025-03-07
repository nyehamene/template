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
		got.Ident = p.Ident(f)
		got.Type = p.Type(f)
		got.Name = p.Name(f)
		got.Templ = p.Templ(f)
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

	accept := func(i ast.ImportDecl, f ast.NamespaceFile) importName {
		var got importName
		got.Ident = i.Ident(f)
		got.Type = i.Type(f)
		got.Path = i.Path(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type usingName struct {
	Ident  string
	Idents []string
	Type   string
	Pkg    string
}

func HelperUsing(t *testing.T, testcase TestCase[usingName]) {
	parse := func(p *Parser) (ast.UsingDecl, bool) {
		return p.ParseUsing()
	}

	accept := func(u ast.UsingDecl, f ast.NamespaceFile) usingName {
		var got usingName
		got.Ident = u.Ident(f)
		got.Idents = u.Idents(f)
		got.Type = u.Type(f)
		got.Pkg = u.Pkg(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type typeAliasName struct {
	Ident  string
	Idents []string
	Type   string
	Target string
}

func HelperTypeAlias(t *testing.T, testcase TestCase[typeAliasName]) {
	parse := func(p *Parser) (ast.AliasDecl, bool) {
		return p.ParseAlias()
	}

	accept := func(d ast.AliasDecl, f ast.NamespaceFile) typeAliasName {
		var got typeAliasName
		got.Ident = d.Ident(f)
		got.Idents = d.Idents(f)
		got.Type = d.Type(f)
		got.Target = d.Target(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

func TestPackage(t *testing.T) {
	testcase := TestCase[pkgName]{
		`p : package : package("home") templ(tag)`:  {"p", "package", `"home"`, "tag"},
		`p : package : package("home") templ(list)`: {"p", "package", `"home"`, "list"},
		`p : package : package("home") templ(html)`: {"p", "package", `"home"`, "html"},
		`p :: package("home") templ(tag)`:           {"p", "package", `"home"`, "tag"},
		`p :: package("home") templ(list)`:          {"p", "package", `"home"`, "list"},
		`p :: package("home") templ(html)`:          {"p", "package", `"home"`, "html"},
	}
	HelperPackage(t, testcase)
}

func TestImport(t *testing.T) {
	testcase := TestCase[importName]{
		`i : import : import("lib/one")`: {"i", "import", `"lib/one"`},
		`i :: import("lib/one")`:         {"i", "import", `"lib/one"`},
	}
	HelperImport(t, testcase)
}

func TestUsing(t *testing.T) {
	testcase := TestCase[usingName]{
		"a, bb : using : using(ccc)": {"a", []string{"a", "bb"}, "using", "ccc"},
		"a : using : using(ccc)":     {"a", []string{"a"}, "using", "ccc"},
		"a :: using(ccc)":            {"a", []string{"a"}, "using", "ccc"},
		"a, bb :: using(ccc)":        {"a", []string{"a", "bb"}, "using", "ccc"},
	}
	HelperUsing(t, testcase)
}

func TestTypeAlias(t *testing.T) {
	testcase := TestCase[typeAliasName]{
		"a : type : type(t)":    {"a", []string{"a"}, "type", "t"},
		"a :: type(t)":          {"a", []string{"a"}, "type", "t"},
		"a, b : type : type(t)": {"a", []string{"a", "b"}, "type", "t"},
		"a, b :: type(t)":       {"a", []string{"a", "b"}, "type", "t"},
	}
	HelperTypeAlias(t, testcase)
}
