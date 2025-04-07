package main

import (
	_ "embed"
	"fmt"
	"log"
	"temlang/tem/ast"
	"temlang/tem/parser"
)

//go:embed def.tem
var def []byte

//go:embed def_semicolon.tem
var semicolon []byte

// go:embed template.tem
// var tmpl []byte

// go:embed def_semicolon.tem
// var tmplSemicolon []byte

func main() {
	srcs := map[string][]byte{
		"def.tem": def,
		// "def_semicolon.tem": semicolon,
		// "template.tem": tmpl,
		// "template_semicolon": tmplSemicolon,
	}
	for name, src := range srcs {
		str := getString(name, src)
		fmt.Println(str)
	}
}

func getString(name string, source []byte) string {
	log.Printf("Parsing %s", name)
	return parse(name, source)
}

func parse(filename string, src []byte) string {
	file, errs := parser.ParseFile(filename, src)
	for !errs.Empty() {
		err, ok := errs.Pop()
		if !ok {
			break
		}
		fmt.Printf("%s %d\n", err.Message(), err.Offset())
	}
	str := ast.PrintSExpr(file)
	return str
}
