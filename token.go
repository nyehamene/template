package template

import (
	"errors"
	"fmt"
)

var (
	EOF        error = errors.New("end of file")
	ErrInvalid       = errors.New("Invalid token")
)

var whitespaces = map[byte]bool{
	' ':  true,
	'\t': true,
	'\r': true,
	'\v': true,
	'\f': true,
}

func Tokenize(s *string, offset int) (Token, int, error) {
	t := Token{
		kind:   TokenUndefined,
		source: s,
		offset: offset,
	}
	return t.next()
}

type TokenKind string

const (
	TokenUndefined    TokenKind = "<>"
	TokenColon                  = ":"
	TokenEqual                  = "="
	TokenPeriod                 = "."
	TokenSemicolon              = ";"
	TokenBraceLeft              = "{"
	TokenBraceRight             = "}"
	TokenBracketLeft            = "{"
	TokenBracketRight           = "}"
	TokenParLeft                = "("
	TokenParRight               = ")"
	TokenSpace                  = "<spc>"
	TokenEOL                    = "<eol>"
	TokenIdent                  = "<ident>"
	TokenPackage                = "<package>"
	TokenTag                    = "<tag>"
	TokenList                   = "<list>"
	TokenHtml                   = "<html>"
	TokenType                   = "<type>"
	TokenTempl                  = "<templ>"
	TokenEnd                    = "<end>"
)

var keywords = map[string]TokenKind{
	"package": TokenPackage,
	"tag":     TokenTag,
	"list":    TokenList,
	"html":    TokenHtml,
	"type":    TokenType,
	"templ":   TokenTempl,
	"end":     TokenEnd,
}

type Token struct {
	source *string
	kind   TokenKind
	offset int
}

func (t Token) String() string {
	return fmt.Sprintf("(%s %d)", t.kind, t.offset)
}

func (t Token) Equal(o Token) bool {
	offset := t.offset == o.offset
	kind := t.kind == o.kind
	return offset && kind
}

// TODO: next should not be a method on Token
func (t Token) next() (Token, int, error) {
	src := *t.source
	kind := TokenUndefined
	start := t.offset
	end := t.offset
	var char byte
	var err error

	if t.isEnd() {
		err = EOF
		goto ret
	}

	char = src[start]

	switch char {
	case ':':
		kind = TokenColon
	case '=':
		kind = TokenEqual
	case '.':
		kind = TokenPeriod
	case ';':
		kind = TokenSemicolon
	case '{':
		kind = TokenBraceLeft
	case '}':
		kind = TokenBraceRight
	case '[':
		kind = TokenBracketLeft
	case ']':
		kind = TokenBracketRight
	case '(':
		kind = TokenParLeft
	case ')':
		kind = TokenParRight
	default:
		if whitespaces[char] {
			kind, end, err = t.space()
			break
		}
		if char == '\n' {
			end = t.eol()
			kind = TokenEOL
			break
		}
		kind, end, err = t.ident()
		{
			// check for keyword
			offset := end + 1
			lexeme := src[start:offset]
			if k := keywords[lexeme]; k != "" {
				kind = k
			}
		}
	}

ret:
	offset := end + 1
	next := Token{
		source: t.source,
		kind:   kind,
		offset: start,
	}
	return next, offset, err
}

func (t Token) ident() (TokenKind, int, error) {
	isAlpha := func(c byte) bool {
		lowercase := c >= 'a' && c <= 'z'
		uppercase := c >= 'A' && c <= 'Z'
		underscore := c == '_'
		return lowercase || uppercase || underscore
	}

	isNumeric := func(c byte) bool {
		return c >= '0' && c <= '9'
	}

	isAlphaNumeric := func(c byte) bool {
		return isAlpha(c) || isNumeric(c)
	}

	src := *t.source
	char := src[t.offset]
	if !isAlpha(char) {
		return TokenUndefined, t.offset, ErrInvalid
	}

	offset := t.offset + 1
	for !t.isEndAt(offset) {
		char = src[offset]
		if !isAlphaNumeric(char) {
			break
		}
		offset += 1
	}

	end := offset - 1
	return TokenIdent, end, nil
}

func (t Token) eol() int {
	src := *t.source
	offset := t.offset + 1
	var char byte
	for !t.isEndAt(offset) {
		char = src[offset]
		if char != '\n' {
			break
		}
		offset += 1
	}
	end := offset - 1
	return end
}

func (t Token) space() (TokenKind, int, error) {
	offset := t.offset + 1
	src := *t.source
	var char byte
	for !t.isEndAt(offset) {
		char = src[offset]
		if !whitespaces[char] {
			break
		}
		offset += 1
	}
	end := offset - 1
	return TokenSpace, end, nil
}

func (t Token) isEnd() bool {
	return t.isEndAt(t.offset)
}

func (t Token) isEndAt(i int) bool {
	return i >= len(*t.source)
}
