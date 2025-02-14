package template

import (
	"errors"
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

type TokenKind int

const (
	TokenUndefined TokenKind = iota
	TokenColon
	TokenEqual
	TokenPeriod
	TokenSemicolon
	TokenBraceLeft
	TokenBraceRight
	TokenBracketLeft
	TokenBracketRight
	TokenParLeft
	TokenParRight
	TokenSpace
	TokenEOL
	TokenIdent
	TokenPackage
	TokenTag
	TokenList
	TokenHtml
	TokenType
	TokenRecord
	TokenTempl
	TokenEnd
	TokenString
	TokenComment
)

var keywords = map[string]TokenKind{
	"package":      TokenPackage,
	"package_tag":  TokenTag,
	"package_list": TokenList,
	"package_html": TokenHtml,
	"type":         TokenType,
	"record":       TokenRecord,
	"templ":        TokenTempl,
	"end":          TokenEnd,
}

func NewTokenizer(s string) Tokenizer {
	return Tokenizer{source: &s}
}

type Tokenizer struct {
	source *string
}

func (t Tokenizer) Tokenize(offset int) (Token, int, error) {
	if t.source == nil {
		panic("source is nil")
	}
	return t.next(offset)
}

// TODO: next should not be a method on Token
func (t Tokenizer) next(offset int) (Token, int, error) {
	src := *t.source
	kind := TokenUndefined
	start := offset
	end := offset
	var char byte
	var err error

	if t.isEnd(offset) {
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
		kind, end, err = t.str(start)
	case '/':
		kind, end, err = t.comment(start)
		{ // try parsing multi line comment
			offset := end + 1
			start := offset
			for {
				token, innerEnd, err := t.Tokenize(start)

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
			kind, end, err = t.space(start)
			break
		}
		if char == '\n' {
			end = t.eol(start)
			kind = TokenEOL
			break
		}
		kind, end, err = t.ident(start)
		{
			// check for keyword
			offset := end + 1
			lexeme := src[start:offset]
			if k, ok := keywords[lexeme]; ok {
				kind = k
			}
		}
	}

ret:
	next := end + 1
	token := Token{
		kind:   kind,
		offset: start,
	}
	return token, next, err
}

func (t Tokenizer) ident(offset int) (TokenKind, int, error) {
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
	char := src[offset]
	if !isAlpha(char) {
		return TokenUndefined, offset, ErrInvalid
	}

	// NOTE: If offset is passed by value this should'nt be a problem
	next := offset + 1
	for !isEndAt(src, next) {
		char = src[next]
		if !isAlphaNumeric(char) {
			break
		}
		next += 1
	}

	end := next - 1
	return TokenIdent, end, nil
}

func (t Tokenizer) str(offset int) (TokenKind, int, error) {
	src := *t.source
	// NOTE: If offset is passed by value this should'nt be a problem
	next := offset + 1
	var char byte
	var err error
	for !isEndAt(src, next) {
		char = src[next]
		if char == '"' {
			next += 1
			break
		}
		next += 1
	}

	end := next - 1
	if src[end] != '"' {
		err = ErrUnterminatedString
	}

	return TokenString, end, err
}

func (t Tokenizer) comment(offset int) (TokenKind, int, error) {
	var char byte
	var err error

	src := *t.source
	markerEnd := offset + 1
	end := offset

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

func (t Tokenizer) eol(offset int) int {
	src := *t.source
	next := offset + 1
	var char byte
	for !isEndAt(src, next) {
		char = src[next]
		if char != '\n' {
			break
		}
		next += 1
	}
	end := next - 1
	return end
}

func (t Tokenizer) space(offset int) (TokenKind, int, error) {
	next := offset + 1
	src := *t.source
	var char byte
	for !isEndAt(src, next) {
		char = src[next]
		if !whitespaces[char] {
			break
		}
		next += 1
	}
	end := next - 1
	return TokenSpace, end, nil
}

func (t Tokenizer) isEnd(offset int) bool {
	return isEndAt(*t.source, offset)
}

func (t Tokenizer) Pos(token Token) (int, int) {
	src := *t.source
	lineNumber := 0
	colNumber := 0
	end := token.offset
	if t.isEnd(token.offset) {
		end = len(src)
	}

	for _, c := range src[0:end] {
		if c == '\n' {
			lineNumber += 1
			colNumber = 0
			continue
		}
		colNumber += 1
	}

	return lineNumber, colNumber
}

func isEndAt(source string, i int) bool {
	return i >= len(source)
}

type Token struct {
	kind   TokenKind
	offset int
}
