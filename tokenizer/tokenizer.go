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
		src:      s,
		ch:       0,
		offset:   0,
		errFunc:  defaultErrorHandler,
		errCount: 0,
	}
	return tok
}

func isLetter(c rune) bool {
	// TODO: handle other unicode letter
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		c == '_'
}

func isDigit(c rune) bool {
	// TODO: handle other unicode digit
	return c >= '0' && c <= '9'
}

func defaultErrorHandler(offset int, ch string, msg string) {
	fmt.Printf("%s\t", msg)
	fmt.Printf("%v\t", ch)
	fmt.Printf("at %d\n", offset)
}

type ErrorHandler func(offset int, ch string, msg string)

type Tokenizer struct {
	src      []byte
	ch       rune
	chOffset int
	offset   int
	errFunc  ErrorHandler
	errCount int
}

func (t Tokenizer) ErrorCount() int {
	return t.errCount
}

func (t *Tokenizer) error(offset int, lit, msg string) {
	t.errFunc(offset, lit, msg)
	t.errCount += 1
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

func (t *Tokenizer) Next() token.Token {
	t.advance()

	ch := t.ch
	offset := t.chOffset

	if ch == eof {
		return token.New(token.EOF, offset)
	}

	var kind token.Kind

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
	case '"':
		t.advance()
		t.string()
		kind = token.String
	default:
		if isLetter(ch) {
			t.advance()
			t.ident()
			kind = token.Ident
			lit := string(t.src[offset:t.offset])
			if k, ok := token.IsKeyword(lit); ok {
				kind = k
			}
			break
		}
		kind = token.Invalid
		t.error(offset, string(ch), "Invalid char")
	}

	tok := token.New(kind, offset)
	return tok
}

func (t *Tokenizer) ident() {
	ch := t.ch
	for isLetter(ch) || isDigit(ch) {
		t.advance()
		ch = t.ch
	}
}

func (t *Tokenizer) string() {
	start := t.chOffset - 1
	for {
		ch := t.ch
		if ch == '\n' || ch <= eof {
			str := string(t.src[start:t.offset])
			t.error(t.offset, str, "Unterminated string")
			break
		}
		t.advance()
		if ch == '"' {
			break
		}
	}
}
