package template

import (
	"errors"
)

var (
	EOF                   error = errors.New("end of file")
	ErrNoMatch                  = errors.New("no match found")
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

type Token struct {
	kind   TokenKind
	offset int
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

	if t.isEnd(offset) {
		return Token{}, offset, EOF
	}

	if kind, next, err := t.char(offset); err == nil {
		token := Token{kind, offset}
		return token, next, nil
	} else if kind, next, err := t.strOrTextBlock(offset); err == nil {
		token := Token{kind, offset}
		return token, next, nil
	} else if kind, next, err := t.comment(offset); err == nil {
		token := Token{kind, offset}
		return token, next, nil
	} else if kind, next, err := t.space(offset); err == nil {
		token := Token{kind, offset}
		return token, next, nil
	} else if kind, next, err := t.eol(offset); err == nil {
		token := Token{kind, offset}
		return token, next, nil
	} else if kind, next, err := t.ident(offset); err == nil {
		lexeme := src[offset:next]
		// check if ident is a keyword
		if k, ok := keywords[lexeme]; ok {
			kind = k
		}
		token := Token{kind, offset}
		return token, next, nil
	}

	return Token{}, offset, ErrInvalid
}

func (t Tokenizer) char(offset int) (TokenKind, int, error) {
	src := *t.source
	if t.isEnd(offset) {
		return TokenUndefined, offset, EOF
	}

	char := src[offset]
	next := offset + 1

	switch char {
	case ':':
		return TokenColon, next, nil
	case '=':
		return TokenEqual, next, nil
	case '.':
		return TokenPeriod, next, nil
	case ';':
		return TokenSemicolon, next, nil
	case '{':
		return TokenBraceLeft, next, nil
	case '}':
		return TokenBraceRight, next, nil
	case '[':
		return TokenBracketLeft, next, nil
	case ']':
		return TokenBracketRight, next, nil
	case '(':
		return TokenParLeft, next, nil
	case ')':
		return TokenParRight, next, nil
	}
	return TokenUndefined, offset, ErrNoMatch
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
		return TokenIdent, offset, ErrNoMatch
	}

	next := offset + 1
	for !isEndAt(src, next) {
		char = src[next]
		if !isAlphaNumeric(char) {
			break
		}
		next += 1
	}

	return TokenIdent, next, nil
}

func (t Tokenizer) strOrTextBlock(offset int) (TokenKind, int, error) {
	if kind, n, err := t.textBlock(offset); err == nil {
		return kind, n, nil
	} else {
		return t.str(offset)
	}
}

func (t Tokenizer) textBlock(start int) (TokenKind, int, error) {
	kind, next, err := t.textBlockLine(start)
	if err != nil {
		return kind, start, err
	}

	// treat consecutinve text block lines as a single text block
	// skipping leading whitespaces
	offset := next
	src := *t.source

	for !t.isEnd(offset) {
		char := src[offset]
		// skip spaces
		if whitespaces[char] {
			offset += 1
			continue
		}
		if char != '"' {
			// the next line is not a text block line
			break
		}
		// add next text block line
		_, innerNext, err := t.textBlockLine(offset)
		if err != nil {
			break
		}
		next = innerNext
		offset = next
	}

	return TokenTextBlock, next, nil
}

func (t Tokenizer) textBlockLine(offset int) (TokenKind, int, error) {
	if !t.matchAll(offset, '"', '"', '"') {
		return TokenTextBlock, offset, ErrNoMatch
	}

	src := *t.source
	// marker end
	next := offset + 3
	for !t.isEnd(next) {
		char := src[next]
		if char == '\n' {
			break
		}
		next += 1
	}

	// consume eol
	if t.match(next, '\n') {
		next += 1
	}

	return TokenTextBlock, next, nil
}

func (t Tokenizer) str(offset int) (TokenKind, int, error) {
	if t.isEnd(offset) {
		return TokenString, offset, EOF
	}

	if !t.match(offset, '"') {
		return TokenString, offset, ErrNoMatch
	}

	src := *t.source
	// consume opening "
	next := offset + 1

	// consume string text
	for !isEndAt(src, next) {
		char := src[next]
		if char == '"' {
			break
		}
		next += 1
	}

	if src[next] != '"' {
		return TokenString, next, ErrUnterminatedString
	}

	// consume closing "
	next += 1
	return TokenString, next, nil
}

func (t Tokenizer) comment(start int) (TokenKind, int, error) {
	kind, next, err := t.commentline(start)
	if err != nil {
		return kind, start, err
	}

	// treat consecutive line comments as single token
	// skipping leading whitespaces
	offset := next
	src := *t.source
	for !t.isEnd(offset) {
		char := src[offset]
		// skip spaces
		if whitespaces[char] {
			offset += 1
			continue
		}
		if char != '/' {
			// the next line is not a comment line
			break
		}
		// add the next comment line
		_, innerNext, err := t.commentline(offset)
		if err != nil {
			break
		}
		next = innerNext
		offset = next
	}

	return kind, next, nil
}

func (t Tokenizer) commentline(offset int) (TokenKind, int, error) {
	if t.isEnd(offset) {
		return TokenComment, offset, EOF
	}

	if !t.matchAll(offset, '/', '/') {
		return TokenComment, offset, ErrNoMatch
	}

	// consume opening //
	next := offset + 2

	src := *t.source

	// consume comment text
	for !t.isEnd(next) {
		char := src[next]
		if char == '\n' {
			// end comment
			break
		}
		next += 1
	}

	// consume eol
	if t.match(next, '\n') {
		next += 1
	}

	return TokenComment, next, nil
}

func (t Tokenizer) eol(offset int) (TokenKind, int, error) {
	if t.isEnd(offset) {
		return TokenEOL, offset, EOF
	}

	src := *t.source
	if src[offset] != '\n' {
		return TokenEOL, offset, ErrNoMatch
	}

	// skip consecutive eof
	next := offset
	for !isEndAt(src, next) {
		char := src[next]
		if char != '\n' {
			break
		}
		next += 1
	}
	return TokenEOL, next, nil
}

func (t Tokenizer) space(offset int) (TokenKind, int, error) {
	if t.isEnd(offset) {
		return TokenSpace, offset, EOF
	}

	src := *t.source

	char := src[offset]
	if !whitespaces[char] {
		return TokenSpace, offset, ErrNoMatch
	}

	next := offset + 1

	// consume consecutive whitespace
	for !isEndAt(src, next) {
		char = src[next]
		if !whitespaces[char] {
			break
		}
		next += 1
	}

	return TokenSpace, next, nil
}

func (t Tokenizer) isEnd(offset int) bool {
	return isEndAt(*t.source, offset)
}

func (t Tokenizer) match(offset int, b byte) bool {
	if t.isEnd(offset) {
		return false
	}
	src := *t.source
	return src[offset] == b
}

func (t Tokenizer) matchAll(offset int, b byte, more ...byte) bool {
	if !t.match(offset, b) {
		return false
	}
	next := offset + 1
	for i := 0; !t.isEnd(next) && i < len(more); i++ {
		b := more[i]
		if !t.match(next, b) {
			return false
		}
		next += 1
	}
	return (next - offset) == (len(more) + 1)
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
