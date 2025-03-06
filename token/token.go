package token

func New(kind Kind, offset int) Token {
	return Token{Kind: kind, Offset: offset}
}

type Token struct {
	Kind   Kind
	Offset int
}

type Kind int

const (
	Invalid Kind = iota
	EOF
	EOL

	SymbolBegin
	// BraceClose close curly brace }
	BraceClose
	// BraceOpen open curly brace {
	BraceOpen
	// BracketClose close square bracket ]
	BracketClose
	// BracketOpen open square bracket [
	BracketOpen
	Colon
	Comma
	Dot
	Eq
	ParenClose
	ParenOpen
	Semicolon
	Space
	SymbolEnd

	KeywordBegin
	Alias
	Import
	Package
	Record
	Templ
	Type
	Using
	Tag
	List
	Html
	KeywordEnd

	LiteralBegin
	Ident
	String
	TextBlock
	LiteralEnd

	Comment
)

var keywords = map[string]Kind{
	"package": Package,
	"type":    Type,
	"record":  Record,
	"templ":   Templ,
	"alias":   Alias,
	"import":  Import,
	"using":   Using,
	"tag":     Tag,
	"list":    List,
	"html":    Html,
}

var whitespaces = map[rune]bool{
	' ':  true,
	'\t': true,
	'\r': true,
	'\v': true,
	'\f': true,
}

func IsSpace(r rune) bool {
	_, ok := whitespaces[r]
	return ok
}

func KeywordKind(ident string) (kind Kind, ok bool) {
	kind, ok = keywords[ident]
	return
}

func Keyword(tok Kind) (string, bool) {
	for n, kw := range keywords {
		if kw == tok {
			return n, true
		}
	}
	return "", false
}

func IsKeyword(tok Kind) bool {
	return tok > KeywordBegin && tok < KeywordEnd
}
