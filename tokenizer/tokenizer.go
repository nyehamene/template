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
		offset:          0,
		rdOffset:        0,
		errCount:        0,
		errFunc:         DefaultErrorHandler,
		semicolonFunc:   DefaultSemicolonHandler,
		insertSemicolon: false,
	}
	for _, opt := range opts {
		opt(&tok)
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

func DefaultErrorHandler(offset int, ch string, msg string) {
	fmt.Printf("token error: ")
	fmt.Printf("%s ", msg)
	fmt.Printf("%s ", ch)
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
	offset          int
	rdOffset        int
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
	prevOffset := t.offset
	prevInsertSemicolon := t.insertSemicolon
	prevRdOffset := t.rdOffset

	return func() {
		t.ch = prev
		t.offset = prevOffset
		t.insertSemicolon = prevInsertSemicolon
		t.rdOffset = prevRdOffset
	}
}

func (t *Tokenizer) error(offset int, lexeme, msg string) {
	t.errFunc(offset, lexeme, msg)
	t.errCount += 1
}

func (t Tokenizer) eof() bool {
	if end := len(t.src); t.rdOffset >= end {
		return true
	}
	return false
}

func (t *Tokenizer) advance() {
	if t.eof() {
		t.ch = eof
		t.rdOffset = len(t.src)
		t.offset = t.rdOffset
		return
	}

	offset := t.rdOffset
	t.ch = rune(t.src[offset])
	t.offset = offset
	t.rdOffset += 1
}

func (t *Tokenizer) skipSpace() {
	for {
		ch := t.ch
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

	ch := t.ch
	offset := t.offset

semiColonInsertion:
	if t.insertSemicolon {
		switch true {
		case ch == '\n' || insertSemiBeforeComment:
			reset()
			fallthrough
		case ch == eof:
			t.insertSemicolon = false
			return token.New(token.Semicolon, offset, t.offset)
		}
	}

	if ch == eof {
		return token.New(token.EOF, offset, t.offset)
	}

	switch ch {
	case '}':
		kind = token.BraceClose
		t.advance()
	case '{':
		kind = token.BraceOpen
		t.advance()
	case ']':
		kind = token.BracketClose
		t.advance()
	case '[':
		kind = token.BracketOpen
		t.advance()
	case ':':
		kind = token.Colon
		t.advance()
	case ',':
		kind = token.Comma
		t.advance()
	case '.':
		kind = token.Dot
		t.advance()
	case '=':
		kind = token.Eq
		t.advance()
	case ')':
		kind = token.ParenClose
		t.advance()
	case '(':
		kind = token.ParenOpen
		t.advance()
	case ';':
		kind = token.Semicolon
		t.advance()
	case '\n':
		kind = token.EOL
		t.advance()
	case '"':
		t.advance()
		if t.ch == '"' {
			t.advance()
			if t.ch == '"' {
				t.advance()
				t.textBlock()
				kind = token.TextBlock
				break
			}
			// empty string literal
			kind = token.String
			break
		}
		t.string()
		kind = token.String
	case '/':
		t.advance()
		if t.ch != '/' {
			kind = token.Invalid
			offset := t.rdOffset - 1
			cmt := string(t.src[offset:t.offset])
			t.error(t.offset, cmt, "Invalid comment marker")
			break
		}

		if t.insertSemicolon {
			insertSemiBeforeComment = true
			goto semiColonInsertion
		}
		t.comment()
		kind = token.Comment
	default:
		t.advance()
		if isLetter(ch) {
			t.ident()
			kind = token.Ident
			lexeme := string(t.src[offset:t.offset])
			if k, ok := token.KeywordKind(lexeme); ok {
				kind = k
			}
			break
		}
		kind = token.Invalid
		t.error(offset, string(ch), "Invalid char")
	}

	t.semicolonFunc(t, kind)

	tok := token.New(kind, offset, t.offset)
	return tok
}

func (t *Tokenizer) ident() {
	for {
		ch := t.ch
		if !isLetter(ch) && !isDigit(ch) {
			break
		}
		t.advance()
	}
}

func (t *Tokenizer) string() {
	markerStart := t.offset - 1

	for {
		ch := t.ch
		if ch == '\n' || ch <= eof {
			str := string(t.src[markerStart:t.offset])
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
	t.consumeUntil('\n')
}

func (t *Tokenizer) comment() {
	t.consumeUntil('\n')
}

func (t *Tokenizer) consumeUntil(r rune) {
	for {
		ch := t.ch
		if ch == r || ch == eof {
			break
		}
		t.advance()
	}
}
