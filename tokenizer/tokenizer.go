package tokenizer

import (
	"fmt"
	"temlang/tem/token"
)

const (
	eof = rune(-1)
)

func New(s []byte, opts ...Option) Tokenizer {
	// TODO: add namespace file path field to tokenizer struct
	tok := Tokenizer{
		src:             s,
		ch:              0,
		chOffset:        0,
		offset:          0,
		errCount:        0,
		errFunc:         DefaultErrorHandler,
		semicolonFunc:   DefaultSemicolonHandler,
		insertSemicolon: false,
	}
	for _, opt := range opts {
		opt(&tok)
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

func DefaultErrorHandler(offset int, ch string, msg string) {
	fmt.Printf("error: ")
	fmt.Printf("%s ", msg)
	fmt.Printf("%v ", ch)
	fmt.Printf("at %d\n", offset)
}

func DefaultSemicolonHandler(t *Tokenizer, kind token.Kind) {
	switch kind {
	case token.Invalid, token.Comment:
		// preserve insertSemicolon
	case token.Ident,
		token.String,
		token.TextBlock,
		token.BracketClose,
		token.BraceClose,
		token.ParenClose:
		t.insertSemicolon = true
	default:
		t.insertSemicolon = false
	}
}

type ErrorHandler func(offset int, ch string, msg string)

type SemicolonHandler func(*Tokenizer, token.Kind)

type Tokenizer struct {
	src             []byte
	ch              rune
	chOffset        int
	offset          int
	insertSemicolon bool
	errFunc         ErrorHandler
	errCount        int
	semicolonFunc   SemicolonHandler
}

func (t Tokenizer) ErrorCount() int {
	return t.errCount
}

func (t *Tokenizer) Mark() func() {
	prev := t.ch
	prevChOffset := t.chOffset
	prevInsertSemicolon := t.insertSemicolon
	prevOffset := t.offset

	return func() {
		t.ch = prev
		t.chOffset = prevChOffset
		t.insertSemicolon = prevInsertSemicolon
		t.offset = prevOffset
	}
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

func (t *Tokenizer) advance() {
	if t.eof() {
		t.ch = eof
		t.offset = len(t.src)
		t.chOffset = t.offset
		return
	}

	offset := t.offset
	ch := t.src[offset]

	t.ch = rune(ch)
	t.chOffset = offset
	t.offset += 1
}

func (t *Tokenizer) skipSpace() {
	for {
		ch := t.peek()
		if !token.IsSpace(ch) {
			break
		}
		t.advance()
	}
}

func (t *Tokenizer) Next() token.Token {
	var kind token.Kind

	t.skipSpace()

	insertSemiBeforeComment := false
	// used to restore tokenizer state after inserting a semicolon
	// before a trailing comment
	reset := t.Mark()

	// consume the next char (not whitespace)
	t.advance()

	ch := t.ch
	offset := t.chOffset

semiColonInsertion:
	if t.insertSemicolon {
		switch true {
		case ch == '\n' || insertSemiBeforeComment:
			reset()
			fallthrough
		case ch == eof:
			t.insertSemicolon = false
			return token.New(token.Semicolon, offset)
		}
	}

	if ch == eof {
		return token.New(token.EOF, offset)
	}

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
	case '\n':
		kind = token.EOL
	case '"':
		t.advance()
		if t.ch != '"' {
			t.string()
			kind = token.String
			break
		}
		if t.peek() == '"' {
			t.textBlock()
			kind = token.TextBlock
			break
		}
		// empty string literal
		kind = token.String
	case '/':
		// consume second /
		t.advance()
		if t.ch != '/' {
			kind = token.Invalid
			offset := t.offset - 1
			cmt := string(t.src[offset:t.offset])
			t.error(t.chOffset, cmt, "Invalid comment marker")
			break
		}

		if t.insertSemicolon {
			insertSemiBeforeComment = true
			goto semiColonInsertion
		}
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

	t.semicolonFunc(t, kind)

	tok := token.New(kind, offset)
	return tok
}

func (t *Tokenizer) ident() {
	for {
		ch := t.peek()
		if !isLetter(ch) && !isDigit(ch) {
			break
		}
		t.advance()
	}
}

func (t *Tokenizer) string() {
	start := t.offset

	for {
		ch := t.peek()
		if ch == '\n' || ch <= eof {
			str := string(t.src[start:t.offset])
			t.error(t.offset, str, "Unterminated string")
			break
		}
		if ch == '"' {
			t.advance()
			break
		}
		t.advance()
	}
}

func (t *Tokenizer) textBlock() {
	t.consumeUntil('\n')
}

func (t *Tokenizer) comment() {
	t.consumeUntil('\n')
}

func (t *Tokenizer) consumeUntil(r rune) {
	for {
		ch := t.peek()
		if ch == r || ch == eof {
			break
		}
		t.advance()
	}
}
