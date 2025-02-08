package main

import (
	"log"
	"os"
)

func main() {
	const argSize = 2
	const templateSourceArgIndex = 1

	if l := len(os.Args); l != argSize {
		log.Fatalf("Expected %d cmd arg but got %d", argSize, l)
	}

	path := os.Args[templateSourceArgIndex]
	tokenize(path)
}
