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
	accept func(E, ast.Namespace) T,
) {

	for src, expected := range testcase {
		t.Run(src, func(t *testing.T) {
			file := ast.New(src, "test.tem")
			file.Name = "ns"
			file.Path = "testing/ns.tem"

			p := New(file)
			pkg, ok := parse(&p)

			if !ok {
				t.Error("parsing failed")
			}

			got := accept(pkg, *file)

			if diff := cmp.Diff(expected, got); diff != "" {
				t.Error(diff)
			}

		})
	}
}

type pkgName struct {
	Idents []string
	Type   string
	Name   string
	Templ  string
}

func HelperPackage(t *testing.T, testcase TestCase[pkgName]) {
	parse := func(p *Parser) (ast.PackageDecl, bool) {
		return p.parsePackageDecl()
	}

	accept := func(p ast.PackageDecl, f ast.Namespace) pkgName {
		var got pkgName
		got.Idents = p.Idents(f)
		got.Type = p.Type(f)
		got.Name = p.Name(f)
		got.Templ = p.Templ(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type importName struct {
	Idents []string
	Type   string
	Path   string
}

func HelperImport(t *testing.T, testcase TestCase[importName]) {
	parse := func(p *Parser) (ast.ImportDecl, bool) {
		return p.parseImportDecl()
	}

	accept := func(i ast.ImportDecl, f ast.Namespace) importName {
		var got importName
		got.Idents = i.Idents(f)
		got.Type = i.Type(f)
		got.Path = i.Name(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type usingName struct {
	Idents []string
	Type   string
	Pkg    string
}

func HelperUsing(t *testing.T, testcase TestCase[usingName]) {
	parse := func(p *Parser) (ast.UsingDecl, bool) {
		return p.parseUsingDecl()
	}

	accept := func(u ast.UsingDecl, f ast.Namespace) usingName {
		var got usingName
		got.Idents = u.Idents(f)
		got.Type = u.Type(f)
		got.Pkg = u.Pkg(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type typeName struct {
	Idents []string
	Type   string
	Target string
}

func HelperType(t *testing.T, testcase TestCase[typeName]) {
	parse := func(p *Parser) (ast.TypeDecl, bool) {
		return p.parseTypeDecl()
	}

	accept := func(d ast.TypeDecl, f ast.Namespace) typeName {
		var got typeName
		got.Idents = d.Idents(f)
		got.Type = d.Type(f)
		got.Target = d.Target(f)
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

type recordName struct {
	Idents []string
	Type   string
	Fields []varName
}

type varName struct {
	Idents []string
	Type   string
}

func HelperRecord(t *testing.T, testcase TestCase[recordName]) {
	parse := func(p *Parser) (ast.RecordDecl, bool) {
		return p.parseRecordDecl()
	}

	accept := func(d ast.RecordDecl, f ast.Namespace) recordName {
		var got recordName
		got.Idents = d.Idents(f)
		got.Type = d.Type(f)
		for _, v := range d.Fields(f) {
			var (
				i = v.Idents(f)
				t = v.Type(f)
			)
			got.Fields = append(got.Fields, varName{Idents: i, Type: t})
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
		return p.parseDocDecl()
	}

	accept := func(d ast.DocDecl, f ast.Namespace) docName {
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
		return p.parseTagDecl()
	}

	accept := func(d ast.TagDecl, f ast.Namespace) tagName {
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

type templName struct {
	Idents []string
	Params []varName
	Type   string
}

func HelperTempl(t *testing.T, testcase TestCase[templName]) {
	parse := func(p *Parser) (ast.TemplDecl, bool) {
		return p.parseTemplDecl()
	}

	accept := func(d ast.TemplDecl, f ast.Namespace) templName {
		var got templName
		got.Idents = d.Idents(f)
		got.Type = d.Type(f)
		for _, attr := range d.Params(f) {
			var (
				i = attr.Idents(f)
				t = attr.Type(f)
			)
			got.Params = append(got.Params, varName{i, t})
		}
		return got
	}

	HelperParser(t, testcase, parse, accept)
}

func TestPackage(t *testing.T) {
	testcase := TestCase[pkgName]{
		`p : package : #tag package("home")`:  {[]string{"p"}, "package", `"home"`, "#tag"},
		`p : package : #list package("home")`: {[]string{"p"}, "package", `"home"`, "#list"},
		`p : package : #html package("home")`: {[]string{"p"}, "package", `"home"`, "#html"},
		`p :: #tag package("home")`:           {[]string{"p"}, "package", `"home"`, "#tag"},
		`p :: #list package("home")`:          {[]string{"p"}, "package", `"home"`, "#list"},
		`p :: #html package("home")`:          {[]string{"p"}, "package", `"home"`, "#html"},
	}
	HelperPackage(t, testcase)
}

func TestImport(t *testing.T) {
	testcase := TestCase[importName]{
		`i : import : import("lib/one")`: {[]string{"i"}, "import", `"lib/one"`},
		`i :: import("lib/one")`:         {[]string{"i"}, "import", `"lib/one"`},
	}
	HelperImport(t, testcase)
}

func TestUsing(t *testing.T) {
	testcase := TestCase[usingName]{
		"a, bb : using : using(ccc)": {[]string{"a", "bb"}, "using", "ccc"},
		"a, bb :: using(ccc)":        {[]string{"a", "bb"}, "using", "ccc"},
		"a : using : using(ccc)":     {[]string{"a"}, "using", "ccc"},
		"a :: using(ccc)":            {[]string{"a"}, "using", "ccc"},
	}
	HelperUsing(t, testcase)
}

func TestType(t *testing.T) {
	testcase := TestCase[typeName]{
		"a : type : type(t)":    {[]string{"a"}, "type", "t"},
		"a :: type(t)":          {[]string{"a"}, "type", "t"},
		"a, b : type : type(t)": {[]string{"a", "b"}, "type", "t"},
		"a, b :: type(t)":       {[]string{"a", "b"}, "type", "t"},
	}
	HelperType(t, testcase)
}

func TestRecord(t *testing.T) {
	testcase := TestCase[recordName]{
		"A : type : record { name: String; email: String; }": {[]string{"A"}, "record", []varName{{[]string{"name"}, "String"}, {[]string{"email"}, "String"}}},
		"A : type : record { name: String; email: String  }": {[]string{"A"}, "record", []varName{{[]string{"name"}, "String"}, {[]string{"email"}, "String"}}},
		"A :: record { name: String; email: String; }":       {[]string{"A"}, "record", []varName{{[]string{"name"}, "String"}, {[]string{"email"}, "String"}}},
		"A :: record { name: String; email: String  }":       {[]string{"A"}, "record", []varName{{[]string{"name"}, "String"}, {[]string{"email"}, "String"}}},
		"A, B : type : record { name: String; }":             {[]string{"A", "B"}, "record", []varName{{[]string{"name"}, "String"}}},
		"A, B : type : record { name: String  }":             {[]string{"A", "B"}, "record", []varName{{[]string{"name"}, "String"}}},
		"A, B :: record { name: String; }":                   {[]string{"A", "B"}, "record", []varName{{[]string{"name"}, "String"}}},
		"A, B :: record { name: String  }":                   {[]string{"A", "B"}, "record", []varName{{[]string{"name"}, "String"}}},
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

func TestTempl(t *testing.T) {
	testcase := TestCase[templName]{
		`A :: templ(a: A){}`:              {[]string{"A"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A :: (a: A){}`:                   {[]string{"A"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A : templ : templ(a: A){}`:       {[]string{"A"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A : templ : (a: A){}`:            {[]string{"A"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A, B :: templ(a: A){}`:           {[]string{"A", "B"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A, B :: (a: A){}`:                {[]string{"A", "B"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A, B : templ : templ(a: A){}`:    {[]string{"A", "B"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A, B : templ : (a: A){}`:         {[]string{"A", "B"}, []varName{{[]string{"a"}, "A"}}, "templ"},
		`A :: templ(a: type){}`:           {[]string{"A"}, []varName{{[]string{"a"}, "type"}}, "templ"},
		`A :: (a: type){}`:                {[]string{"A"}, []varName{{[]string{"a"}, "type"}}, "templ"},
		`A : templ : templ(a: type){}`:    {[]string{"A"}, []varName{{[]string{"a"}, "type"}}, "templ"},
		`A : templ : (a: type){}`:         {[]string{"A"}, []varName{{[]string{"a"}, "type"}}, "templ"},
		`A, B :: templ(a: type){}`:        {[]string{"A", "B"}, []varName{{[]string{"a"}, "type"}}, "templ"},
		`A, B :: (a: type){}`:             {[]string{"A", "B"}, []varName{{[]string{"a"}, "type"}}, "templ"},
		`A, B : templ : templ(a: type){}`: {[]string{"A", "B"}, []varName{{[]string{"a"}, "type"}}, "templ"},
		`A, B : templ : (a: type){}`:      {[]string{"A", "B"}, []varName{{[]string{"a"}, "type"}}, "templ"},
	}
	HelperTempl(t, testcase)
}
