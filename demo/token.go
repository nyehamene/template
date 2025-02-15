package main

import (
	"fmt"
	"lang/template"
	"log"
)

func tokenize(source []byte) {
	str := string(source)
	tokenizer := template.NewTokenizer(str)
	tokens := []template.Token{}
	next := 0
	for {
		token, offset, err := tokenizer.Tokenize(next)
		if err == template.EOF {
			break
		}

		if err != nil {
			line, col := tokenizer.Pos(token)
			log.Fatalf("%v @ [%d, %d]", err, line, col)
		}

		tokens = append(tokens, token)
		next = offset
	}

	for _, token := range tokens {
		line, col := tokenizer.Pos(token)
		fmt.Printf("(%v [%d, %d])\n", token, line, col)
	}
}
