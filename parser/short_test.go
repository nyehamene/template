package parser_test

import (
	"fmt"
	"temlang/tem/parser"
	"testing"
)

var valids = []string{
	// package
	"p :: package(\"m\")\n",
	`p :: package("m")`,
	`p :: package("m");`,
	`p :: package("m");`,
	`p :: #html package("m");`,
	`p :: #tag package("m");`,
	`p :: #lisp package("m");`,
	`p :: #lisp #dir package("m");`,
	// package with type annotation
	"p: package: package(\"m\")\n",
	`p: package: package("m")`,
	`p: package: package("m");`,
	`p: package: package("m");`,
	`p: package: #html package("m");`,
	`p: package: #tag package("m");`,
	`p: package: #lisp package("m");`,
	`p: package: #lisp #dir package("m");`,
	// import
	"p :: package(\"m\"); i :: import(\"p\")\n",
	`p :: package("m");   i :: import("p")`,
	`p :: package("m");   i :: import("p");`,
	`p :: package("m");   i :: #dir import("p");`,
	`p :: package("m");   i :: #dir1 #dir2 import("p");`,
	// import with type annotation
	"p :: package(\"m\"); i: import: import(\"p\")\n",
	`p :: package("m");   i: import: import("p")`,
	`p :: package("m");   i: import: import("p");`,
	`p :: package("m");   i: import: #dir import("p");`,
	`p :: package("m");   i: import: #dir1 #dir2 import("p");`,
	// using
	"p :: package(\"m\"); i :: import(\"p\")\n     a, b :: using(p)\n",
	`p :: package("m");   i :: import("p");        a, b :: using(p)`,
	`p :: package("m");   i :: import("p");        a, b :: using(p);`,
	`p :: package("m");   i :: #dir import("p");   a, b :: using(p);`,
	`p :: package("m");   i :: #d #d2 import("p"); a, b :: using(p);`,
	// using with type annotation
	"p :: package(\"m\"); i :: import(\"p\")\n     a, b: import: using(p)\n",
	`p :: package("m");   i :: import("p");        a, b: import: using(p)`,
	`p :: package("m");   i :: import("p");        a, b: import: using(p);`,
	`p :: package("m");   i :: #d import("p");     a, b: import: using(p);`,
	`p :: package("m");   i :: #d #d2 import("p"); a, b: import: using(p);`,
	// var
	"p :: package(\"m\"); t : String\n",
	`p :: package("m");   t : String`,
	`p :: package("m");   t : String;`,
	// derived type
	"p :: package(\"m\"); t :: type(String)\n",
	`p :: package("m");   t :: type(String)`,
	`p :: package("m");   t :: type(String);`,
	`p :: package("m");   t :: #d type(String);`,
	`p :: package("m");   t :: #d #d1 type(String);`,
	// derived type with type annotation
	"p :: package(\"m\"); t: type: type(String)\n",
	`p :: package("m");   t: type: type(String)`,
	`p :: package("m");   t: type: type(String);`,
	`p :: package("m");   t: type: #d type(String);`,
	`p :: package("m");   t: type: #d #d2 type(String);`,
	// record type
	"p :: package(\"m\"); t :: record{ a: String\n}\n",
	`p :: package("m");   t :: record{ a: String }`,
	`p :: package("m");   t :: record{ a: String; }`,
	`p :: package("m");   t :: record{ a: String; };`,
	`p :: package("m");   t :: record{ a: String; b: String };`,
	`p :: package("m");   t :: #d record{ a: String; b: String };`,
	`p :: package("m");   t :: #d #d2 record{ a: String; b: String };`,
	// record type with type annotation
	"p :: package(\"m\"); t: type: record{ a: String\n}\n",
	`p :: package("m");   t: type: record{ a: String }`,
	`p :: package("m");   t: type: record{ a: String; }`,
	`p :: package("m");   t: type: record{ a: String; };`,
	`p :: package("m");   t: type: record{ a: String; b: String };`,
	`p :: package("m");   t: type: #d record{ a: String; b: String };`,
	`p :: package("m");   t: type: #d #d2 record{ a: String; b: String };`,
	// template
	"p :: package(\"m\"); c :: templ(m: Model){}\n",
	`p :: package("m");   c :: templ(m: Model){}`,
	`p :: package("m");   c :: templ(m: Model){};`,
	// template with type annotation
	"p :: package(\"m\"); c: templ: templ(m: Model){}\n",
	`p :: package("m");   c: templ: templ(m: Model){}`,
	`p :: package("m");   c: templ: templ(m: Model){};`,
}

func TestValids(t *testing.T) {
	const filename = "valid.tem"
	for i, src := range valids {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, errs := parser.ParseFile(filename, []byte(src))
			if errs.Len() != 0 {
				t.Errorf("ParseFile(%s) parser failed unexpectedly", filename)
			}
		})
	}
}
