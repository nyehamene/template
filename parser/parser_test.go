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

type recordName struct {
	Ident  string
	Idents []string
	Type   string
	Fields []varName
}

type varName struct {
	Ident string
	Type  string
}

func HelperRecord(t *testing.T, testcase TestCase[recordName]) {
	parse := func(p *Parser) (ast.RecordDecl, bool) {
		return p.ParseRecord()
	}

	accept := func(d ast.RecordDecl, f ast.NamespaceFile) recordName {
		var got recordName
		got.Ident = d.Ident(f)
		got.Idents = d.Idents(f)
		got.Type = d.Type(f)
		for _, v := range d.Fields(f) {
			var (
				ident = v.Ident(f)
				ty    = v.Type(f)
			)
			got.Fields = append(got.Fields, varName{Ident: ident, Type: ty})
		}
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type docName struct {
	Idents []string
	Text   string
}

func HelperDoc(t *testing.T, testcase TestCase[docName]) {
	parse := func(p *Parser) (ast.DocDecl, bool) {
		return p.ParseDoc()
	}

	accept := func(d ast.DocDecl, f ast.NamespaceFile) docName {
		var got docName
		got.Idents = d.Idents(f)
		got.Text = d.Content(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type tagName struct {
	Idents []string
	Attrs  []attrName
}

type attrName struct {
	Key   []string
	Value string
}

func HelperTag(t *testing.T, testcase TestCase[tagName]) {
	parse := func(p *Parser) (ast.TagDecl, bool) {
		return p.ParseTag()
	}

	accept := func(d ast.TagDecl, f ast.NamespaceFile) tagName {
		var got tagName
		got.Idents = d.Idents(f)
		for _, attr := range d.Attrs(f) {
			var (
				key   = attr.Idents(f)
				value = attr.Value(f)
			)
			got.Attrs = append(got.Attrs, attrName{key, value})
		}
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

func TestPackage(t *testing.T) {
	testcase := TestCase[pkgName]{
		`p : package : package("home") tag`:  {"p", "package", `"home"`, "tag"},
		`p : package : package("home") list`: {"p", "package", `"home"`, "list"},
		`p : package : package("home") html`: {"p", "package", `"home"`, "html"},
		`p :: package("home") tag`:           {"p", "package", `"home"`, "tag"},
		`p :: package("home") list`:          {"p", "package", `"home"`, "list"},
		`p :: package("home") html`:          {"p", "package", `"home"`, "html"},
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
		"a, bb :: using(ccc)":        {"a", []string{"a", "bb"}, "using", "ccc"},
		"a : using : using(ccc)":     {"a", []string{"a"}, "using", "ccc"},
		"a :: using(ccc)":            {"a", []string{"a"}, "using", "ccc"},
	}
	HelperUsing(t, testcase)
}

func TestTypeAlias(t *testing.T) {
	testcase := TestCase[typeAliasName]{
		"a : type : type(t)":    {"a", []string{"a"}, "alias", "t"},
		"a :: type(t)":          {"a", []string{"a"}, "alias", "t"},
		"a, b : type : type(t)": {"a", []string{"a", "b"}, "alias", "t"},
		"a, b :: type(t)":       {"a", []string{"a", "b"}, "alias", "t"},
	}
	HelperTypeAlias(t, testcase)
}

func TestRecord(t *testing.T) {
	testcase := TestCase[recordName]{
		"A : type : record { name: String; email: String; }": {"A", []string{"A"}, "record", []varName{{"name", "String"}, {"email", "String"}}},
		"A : type : record { name: String; email: String  }": {"A", []string{"A"}, "record", []varName{{"name", "String"}, {"email", "String"}}},

		"A, B : type : record { name: String; }": {"A", []string{"A", "B"}, "record", []varName{{"name", "String"}}},
		"A, B : type : record { name: String  }": {"A", []string{"A", "B"}, "record", []varName{{"name", "String"}}},

		"A :: record { name: String; email: String; }": {"A", []string{"A"}, "record", []varName{{"name", "String"}, {"email", "String"}}},
		"A :: record { name: String; email: String  }": {"A", []string{"A"}, "record", []varName{{"name", "String"}, {"email", "String"}}},

		"A, B :: record { name: String; }": {"A", []string{"A", "B"}, "record", []varName{{"name", "String"}}},
		"A, B :: record { name: String  }": {"A", []string{"A", "B"}, "record", []varName{{"name", "String"}}},
	}
	HelperRecord(t, testcase)
}

func TestDoc(t *testing.T) {
	testcase := TestCase[docName]{
		`A : "line 1"`:   {[]string{"A"}, `"line 1"`},
		`A : """ line 1`: {[]string{"A"}, `""" line 1`},
		`A : """ line 1
		     """ line 2`: {[]string{"A"}, "\"\"\" line 1\n\"\"\" line 2"},
		`A, B : "line 1"`: {[]string{"A", "B"}, `"line 1"`},
	}
	HelperDoc(t, testcase)
}

func TestTag(t *testing.T) {
	testcase := TestCase[tagName]{
		`A : { name = "a"; email = "b" }`: {[]string{"A"}, []attrName{{[]string{"name"}, "\"a\""}, {[]string{"email"}, "\"b\""}}},
		`A, B : { name = "a"; }`:          {[]string{"A", "B"}, []attrName{{[]string{"name"}, "\"a\""}}},
	}
	HelperTag(t, testcase)
}
