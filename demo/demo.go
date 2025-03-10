package main

import (
	_ "embed"
	"fmt"
	"log"
	"temlang/tem/ast"
	"temlang/tem/parser"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
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
	log.Println("\nTokenizing")
	tokenize(source)

	log.Println("\nParsing")
	parse(name, source)
}

func tokenize(s []byte) {
	tok := tokenizer.New(s)
	for {
		t := tok.Next()
		if t.Kind() == token.EOF {
			break
		}
		fmt.Println(t)
	}
}

func parse(name string, s []byte) {
	file := ast.New(string(s), name)
	par := parser.New(file)
	par.ParseFile()
}
