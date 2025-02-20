package main

import (
	"fmt"
	"lang/template"
)

func parse(source []byte) {
	src := string(source)
	t := template.NewTokenizer(src)
	p := template.NewParser(&t)

	asts, err := p.Parse(0)
	if err != nil {
		panic(err)
	}

	for _, ast := range asts {
		fmt.Printf("(%v)\n", ast)
	}
}
