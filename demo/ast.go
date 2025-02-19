package main

import (
	"fmt"
	"lang/template"
	"log"
)

func parse(source []byte) {
	src := string(source)
	t := template.NewTokenizer(src)
	p := template.NewParser(t)

	asts := []template.Ast{}
	next := 0

	pkg, offset, err := p.Parse(next)
	if err != nil {
		panic(err)
	}
	asts = append(asts, pkg)
	next = offset

	for {
		ast, offset, err := p.Parse(next)
		if err == template.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		asts = append(asts, ast)
		next = offset
	}

	for _, ast := range asts {
		fmt.Printf("(%v)\n", ast)
	}
}
