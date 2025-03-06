package parser

import (
	"temlang/tem/ast"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestCase[T any] map[string]T

func TestPackage(t *testing.T) {
	type Name struct {
		Ident string
		Type  string
		Name  string
		Templ string
	}
	testcase := TestCase[Name]{
		`p : package : package("home") templ(tag)`:  {"p", "package", `"home"`, "tag"},
		`p : package : package("home") templ(list)`: {"p", "package", `"home"`, "list"},
		`p : package : package("home") templ(html)`: {"p", "package", `"home"`, "html"},
	}

	for src, expected := range testcase {
		t.Run(src, func(t *testing.T) {
			file := ast.NamespaceFile{
				Name: "testns",
				Path: "testns.tem",
				Src:  src,
				Pkg:  ast.PackageDecl{},
			}

			p := New(&file)
			pkg, ok := p.ParsePackage()

			if !ok {
				t.Error("parsing failed")
			}

			var got Name
			{
				got.Ident = p.file.GetName(pkg.Ident)
				got.Type = p.file.GetName(pkg.Type)
				got.Name = p.file.GetName(pkg.Name)
				got.Templ = p.file.GetName(pkg.Templ)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}

		})
	}
}

func TestInferedTypePackage(t *testing.T) {
	type Name struct {
		Ident string
		Name  string
		Type  string
		Templ string
	}
	testcase := TestCase[Name]{
		`p :: package("home") templ(tag)`:  {"p", `"home"`, "", "tag"},
		`p :: package("home") templ(list)`: {"p", `"home"`, "", "list"},
		`p :: package("home") templ(html)`: {"p", `"home"`, "", "html"},
	}

	for src, expected := range testcase {
		t.Run(src, func(t *testing.T) {
			file := ast.NamespaceFile{
				Name: "testns",
				Path: "testns.tem",
				Src:  src,
				Pkg:  ast.PackageDecl{},
			}

			p := New(&file)
			pkg, ok := p.ParsePackage()

			if !ok {
				t.Error("parsing failed")
			}

			var got Name
			{
				got.Ident = p.file.GetName(pkg.Ident)
				got.Name = p.file.GetName(pkg.Name)
				got.Templ = p.file.GetName(pkg.Templ)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}

		})
	}
}

func TestImport(t *testing.T) {
	type Name struct {
		Ident string
		Type  string
		Path  string
	}
	testcase := TestCase[Name]{
		`i : import : import("lib/one")`:   {"i", "import", `"lib/one"`},
		`i : import : import("lib/two")`:   {"i", "import", `"lib/two"`},
		`i : import : import("lib/three")`: {"i", "import", `"lib/three"`},
	}

	for src, expected := range testcase {
		t.Run(src, func(t *testing.T) {
			file := ast.NamespaceFile{
				Name: "testns",
				Path: "testns.tem",
				Src:  src,
			}

			p := New(&file)
			imp, ok := p.ParseImport()

			if !ok {
				t.Error("parsing failed")
			}

			var got Name
			{
				got.Ident = p.file.GetName(imp.Ident)
				got.Type = p.file.GetName(imp.Type)
				got.Path = p.file.GetName(imp.Path)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}

		})
	}
}

func TestInferedTypeImport(t *testing.T) {
	type Name struct {
		Ident string
		Type  string
		Path  string
	}
	testcase := TestCase[Name]{
		`i :: import("lib/one")`:   {"i", "", `"lib/one"`},
		`i :: import("lib/two")`:   {"i", "", `"lib/two"`},
		`i :: import("lib/three")`: {"i", "", `"lib/three"`},
	}

	for src, expected := range testcase {
		t.Run(src, func(t *testing.T) {
			file := ast.NamespaceFile{
				Name: "testns",
				Path: "testns.tem",
				Src:  src,
			}

			p := New(&file)
			imp, ok := p.ParseImport()

			if !ok {
				t.Error("parsing failed")
			}

			var got Name
			{
				got.Ident = p.file.GetName(imp.Ident)
				got.Type = p.file.GetName(imp.Type)
				got.Path = p.file.GetName(imp.Path)
			}

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}

		})
	}
}
