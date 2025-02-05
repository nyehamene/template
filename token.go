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
	TokenColon               = ":"
	TokenEqual               = "="
	TokenPeriod              = "."
	TokenSemicolon           = ";"
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
	kind := TokenUndefined
	offset := t.offset
	var err error

	src := *t.source

	if t.offset >= len(src) {
		return t.undefined(), t.offset, EOF
	}

	char := src[t.offset]

	switch char {
	case ':':
		kind = TokenColon
	case '=':
		kind = TokenEqual
	case '.':
		kind = TokenPeriod
	case ';':
		kind = TokenSemicolon
	default:
		err = ErrInvalid
	}

	next := Token{
		source: t.source,
		kind:   kind,
		offset: t.offset,
	}
	return next, offset + 1, err
}

func (t Token) undefined() Token {
	return Token{
		source: t.source,
		kind:   TokenUndefined,
		offset: t.offset,
	}
}
