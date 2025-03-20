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
		"def.tem":           def,
		"def_semicolon.tem": semicolon,
		// "template.tem": tmpl,
		// "template_semicolon": tmplSemicolon,
	}
	for name, src := range srcs {
		run(name, src)
	}
}

func run(name string, source []byte) {
	log.Printf("Parsing %s", name)
	parse(name, source)
	println()
}

func parse(filename string, src []byte) {
	file, errs := parser.ParseFile(filename, src)
	if !errs.Empty() {
		for !errs.Empty() {
			err, ok := errs.Pop()
			if !ok {
				break
			}
			fmt.Printf("%s %s\n", err.Msg, err.Location)
		}
		return
	}
	prt := ast.NewPrinter(file)
	prt.Print()
}
