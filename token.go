package template

import (
	"errors"
	"fmt"
)

var (
	EOF                   error = errors.New("end of file")
	ErrInvalid                  = errors.New("Invalid token")
	ErrUnterminatedString       = errors.New("Unterminated string")
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
	TokenString                 = "<str>"
	TokenComment                = "<comment>"
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
	case '"':
		kind, end, err = t.str()
	case '/':
		kind, end, err = t.comment()
		{
			// try parsing multi line comment
			offset := end + 1
			start := offset
			for {
				token, innerEnd, err := Tokenize(t.source, start)

				if err != nil {
					break
				}

				if token.kind == TokenSpace {
					start = innerEnd
					continue
				}

				if token.kind != TokenComment {
					break
				}

				offset = innerEnd
				start = offset
			}

			end = offset - 1
		}
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
	for !isEndAt(src, offset) {
		char = src[offset]
		if !isAlphaNumeric(char) {
			break
		}
		offset += 1
	}

	end := offset - 1
	return TokenIdent, end, nil
}

func (t Token) str() (TokenKind, int, error) {
	src := *t.source
	offset := t.offset + 1
	var char byte
	var err error
	for !isEndAt(src, offset) {
		char = src[offset]
		if char == '"' {
			offset += 1
			break
		}
		offset += 1
	}

	end := offset - 1
	if src[end] != '"' {
		err = ErrUnterminatedString
	}

	return TokenString, end, err
}

func (t Token) comment() (TokenKind, int, error) {
	var char byte
	var err error

	src := *t.source
	markerEnd := t.offset + 1
	end := t.offset

	if isEndAt(src, markerEnd) {
		err = ErrInvalid
		goto ret
	}

	if src[markerEnd] != '/' {
		err = ErrInvalid
		goto ret
	}

	end = markerEnd + 1

	for !isEndAt(src, end) {
		char = src[end]
		if char == '\n' {
			break
		}
		end += 1
	}

	// If is isEndAt return true
	if isEndAt(src, end) {
		end -= 1
	}

ret:
	return TokenComment, end, err
}

func (t Token) eol() int {
	src := *t.source
	offset := t.offset + 1
	var char byte
	for !isEndAt(src, offset) {
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
	for !isEndAt(src, offset) {
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
	return isEndAt(*t.source, t.offset)
}

func isEndAt(source string, i int) bool {
	return i >= len(source)
}
