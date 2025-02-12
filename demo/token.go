package main

import (
	"fmt"
	"io"
	"lang/template"
	"log"
	"os"
)

func tokenize(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}

	defer f.Close()

	source, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	str := string(source)
	tokenizer := template.NewTokenizer(str)
	tokens := []template.Token{}
	next := 0
	for {
		var token template.Token
		var offset int
		var err error
		token, offset, err = tokenizer.Tokenize(next)
		if err == template.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v @ %d", err, next)
		}

		tokens = append(tokens, token)
		next = offset
	}

	for _, token := range tokens {
		line, col := tokenizer.Pos(token)
		fmt.Printf("(%v [%d, %d])\n", token, line, col)
	}
}
