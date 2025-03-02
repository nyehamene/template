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
	tok.advance()
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

func (t Tokenizer) eof() bool {
	offset := t.offset
	if end := len(t.src); offset >= end {
		return true
	}
	return false
}

func (t *Tokenizer) peek() rune {
	if t.eof() {
		return eof
	}

	offset := t.offset
	ch := rune(t.src[offset])
	return ch
}

func (t *Tokenizer) advance() rune {
	next := t.peek()
	if next == eof {
		t.ch = eof
		t.offset = len(t.src)
		return eof
	}

	t.ch = next
	t.chOffset = t.offset
	t.offset += 1
	return next
}

func (t Tokenizer) match(tok *Tokenizer, ch rune, others ...rune) bool {
	if t.ch != ch {
		return false
	}

	for _, r := range others {
		ch := t.advance()
		if ch != r {
			return false
		}
	}
	for i := 0; i <= len(others); i++ {
		_ = tok.advance()
	}
	return true
}

func (t *Tokenizer) skipSpace() {
	for {
		ch := t.ch
		// TODO: preserve newline when inserting semicolon automatically
		if !token.IsSpace(ch) && ch != '\n' {
			break
		}
		t.advance()
	}
}

func (t *Tokenizer) Next() token.Token {
	t.skipSpace()

	ch := t.ch
	offset := t.chOffset

	if ch == eof {
		return token.New(token.EOF, offset)
	}

	t.advance()

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
		if t.match(t, '"', '"') {
			t.textBlock()
			kind = token.TextBlock
		} else {
			t.string()
			kind = token.String
		}
	case '/':
		t.comment()
		kind = token.Comment
	default:
		if isLetter(ch) {
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
		ch = t.advance()
	}
}

func (t *Tokenizer) string() {
	start := t.offset - 1
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

func (t *Tokenizer) textBlock() {
	textBlockLine := func() {
		for {
			ch := t.ch
			if ch == '\n' || ch <= eof {
				break
			}
			t.advance()
		}
	}

	markerOffset := 3
	minLineCount := 2
	lineCount := 0

	offset := t.offset - markerOffset

	for {
		lineCount += 1
		textBlockLine()
		t.skipSpace()
		if !t.match(t, '"', '"', '"') {
			break
		}
	}

	if lineCount < minLineCount {
		str := string(t.src[offset:t.offset])
		t.error(offset, str, "a text block must have at least 2 lines of text")
	}
}

func (t *Tokenizer) comment() {
	if !t.match(t, '/') {
		offset := t.offset - 1
		cmt := string(t.src[offset:t.offset])
		t.error(t.chOffset, cmt, "Invalid comment")
	}

	for {
		ch := t.ch
		if ch == '\n' {
			t.advance()
			break
		}
		if ch == eof {
			break
		}
		t.advance()
	}
}
