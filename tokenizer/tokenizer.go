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

func (t Tokenizer) match(tok *Tokenizer, chs ...rune) bool {
	if len(chs) == 0 {
		return false
	}
	for _, r := range chs {
		ch := t.peek()
		if ch != r {
			return false
		}
		t.advance()
	}
	// replay advances on main tokenizer
	for range chs {
		tok.advance()
	}
	return true
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

func (t *Tokenizer) skipNewlineIfMatch(r ...rune) bool {
	copyT := *t
	if !copyT.match(&copyT, '\n') {
		return false
	}

	// consume \n
	copyT.advance()
	copyT.skipSpace()
	// match text block marker
	if copyT.match(&copyT, r...) {
		// consume last char (a.k.a move to \n)
		t.advance()
		// consume optional space after \n
		t.skipSpace()
		return true
	}
	return false
}

func (t *Tokenizer) Next() token.Token {
	t.skipSpace()
	// consume the next char (not whitespace)
	t.advance()

	ch := t.ch
	offset := t.chOffset

	if t.insertSemicolon && (ch == '\n' || ch == eof) {
		t.insertSemicolon = false
		return token.New(token.Semicolon, offset)
	}

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
	case '\n':
		for {
			t.skipSpace()
			if !t.match(t, '\n') {
				break
			}
		}
		kind = token.EOL
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
	textBlockLine := func(tok *Tokenizer) {
		for {
			ch := tok.peek()
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
		textBlockLine(t)
		if !t.skipNewlineIfMatch('"', '"', '"') {
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

	commentLine := func(t *Tokenizer) {
		for {
			ch := t.peek()
			if ch == '\n' || ch == eof {
				break
			}
			t.advance()
		}
	}

	for {
		commentLine(t)
		if !t.skipNewlineIfMatch('/', '/') {
			break
		}
	}
}
