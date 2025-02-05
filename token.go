package template

import (
	"errors"
	"fmt"
)

var (
	EOF        error = errors.New("end of file")
	ErrInvalid       = errors.New("Invalid token")
)

func Tokenize(s *string, offset int) (Token, int, error) {
	t := Token{
		kind:   TokenUndefined,
		source: s,
		offset: offset,
	}
	return t.Next()
}

type TokenType string

const (
	TokenUndefined TokenType = "<>"
	TokenColon               = "<colon>"
)

type Token struct {
	source *string
	kind   TokenType
	offset int
}

func (t Token) String() string {
	return fmt.Sprintf("<%s %d>", t.kind, t.offset)
}

func (t Token) Equal(o Token) bool {
	offset := t.offset == o.offset
	kind := t.kind == o.kind
	return offset && kind
}

func (t Token) Next() (Token, int, error) {
	src := *t.source
	char := src[t.offset]
	if char == ':' {
		next := Token{
			source: t.source,
			kind:   TokenColon,
			offset: t.offset,
		}
		return next, t.offset + 1, EOF
	}
	return t.undefined(), t.offset, ErrInvalid
}

func (t Token) undefined() Token {
	return Token{
		source: t.source,
		kind:   TokenUndefined,
		offset: t.offset,
	}
}
