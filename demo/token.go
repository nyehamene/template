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
	str := string(source)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	tokens := []template.Token{}
	next := 0
	for {
		var token template.Token
		var offset int
		var err error
		token, offset, err = template.Tokenize(&str, next)
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
		fmt.Printf("%v\n", token)
	}
}
