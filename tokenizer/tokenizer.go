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
	fmt.Printf("%s\t", msg)
	fmt.Printf("%v\t", ch)
	fmt.Printf("at %d\n", offset)
}

type Tokenizer struct {
	src      []byte
	ch       rune
	chOffset int
	offset   int
	err      ErrorHandler
}

func (t *Tokenizer) Next() token.Token {
	t.advance()
	ch := t.ch

	var kind token.Kind
	offset := t.chOffset

	switch ch {
	case '}':
		kind = token.BraceClose
	case '{':
		kind = token.BraceOpen
	case ']':
		kind = token.BracketClose
	case '[':
		kind = token.BracketOpen
	case ':':
		kind = token.Colon
	case ',':
		kind = token.Comma
	case '.':
		kind = token.Dot
	case '=':
		kind = token.Eq
	case ')':
		kind = token.ParenClose
	case '(':
		kind = token.ParenOpen
	case ';':
		kind = token.Semicolon
	default:
		kind = token.Invalid
		t.err(ch, offset, "Invalid char")
	}

	tok := token.New(kind, offset)
	return tok
}

func (t *Tokenizer) advance() {
	offset := t.offset
	if l := len(t.src); offset >= l {
		t.ch = eof
		t.offset = l
		return
	}

	ch := rune(t.src[offset])
	t.ch = ch
	t.chOffset = offset
	t.offset = offset + 1
}
