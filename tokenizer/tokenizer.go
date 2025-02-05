package tokenizer

import (
	"fmt"
	"temlang/tem/token"
)

const (
	eof = -1
)

func New(s []byte) Tokenizer {
	tok := Tokenizer{
		src:    s,
		ch:     0,
		offset: 0,
		err:    defaultErrorHandler,
	}
	return tok
}

type ErrorHandler func(ch rune, offset int, msg string)

func defaultErrorHandler(ch rune, offset int, msg string) {
	fmt.Println(msg)
	fmt.Printf("\t%v", ch)
	fmt.Printf("\n\t at %d", offset)
}

type Tokenizer struct {
	src      []byte
	ch       rune
	chOffset int
	offset   int
	err      ErrorHandler
}
