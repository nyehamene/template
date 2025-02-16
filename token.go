package template

import (
	"errors"
)

var (
	EOF                   error = errors.New("end of file")
	ErrInvalid                  = errors.New("Invalid token")
	ErrUseString                = errors.New("Use a string literal instead")
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
	TokenAlias
	TokenTextBlock
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
	"alias":        TokenAlias,
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
		kind, end, err = t.strOrTextBlock(start)
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

func (t Tokenizer) match(offset int, b byte) bool {
	if t.isEnd(offset) {
		return false
	}
	src := *t.source
	return src[offset] == b
}

func (t Tokenizer) strOrTextBlock(offset int) (TokenKind, int, error) {
	if kind, n, err := t.textBlock(offset); err == nil {
		return kind, n, nil
	} else {
		return t.str(offset)
	}
}

func (t Tokenizer) textBlock(start int) (TokenKind, int, error) {
	kind, n, err := t.textBlockLine(start)
	if err != nil {
		return kind, n, err
	}

	offset := n + 1
	advance := offset

	for {
		token, innerEnd, err := t.Tokenize(advance)
		if err != nil {
			break
		}

		if token.kind == TokenSpace {
			advance = innerEnd
			continue
		}
		if token.kind != TokenTextBlock {
			break
		}

		offset = innerEnd
		advance = offset
	}

	end := offset - 1
	return TokenTextBlock, end, nil
}

func (t Tokenizer) textBlockLine(offset int) (TokenKind, int, error) {
	if t.match(offset, '"') && t.match(offset+1, '"') && t.match(offset+2, '"') {
		goto textline
	}
	return TokenTextBlock, offset, ErrInvalid

textline:
	src := *t.source
	// marker end
	end := offset + 3
	// The textblock marker goes until end of line
	// and must not contain any characters except spaces
	for !t.isEnd(end) {
		char := src[end]
		if char == '\n' {
			break
		}
		end += 1
	}

	// If is isEndAt returns true
	if isEndAt(src, end) {
		end = end - 1
	}

	return TokenTextBlock, end, nil
}

func (t Tokenizer) str(offset int) (TokenKind, int, error) {
	src := *t.source
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

	// If is isEndAt returns true
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
