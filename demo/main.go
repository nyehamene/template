package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
)

func main() {
	const argSize = 2
	const templateSourceArgIndex = 1

	if l := len(os.Args); l != argSize {
		log.Fatalf("Expected %d cmd arg but got %d", argSize, l)
	}

	path := os.Args[templateSourceArgIndex]

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}

	defer f.Close()

	source, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}

	log.Println("\nTokenizing")
	tokenize(source)

	log.Println("\nParsing")
	parse(source)
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

func parse(s []byte) {}
