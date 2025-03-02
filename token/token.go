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

	BraceClose
	BraceOpen
	BracketClose
	BracketOpen
	Colon
	Comma
	Dot
	Eq
	ParenClose
	ParenOpen
	Semicolon
	Space

	KeywordBegin
	Alias
	Import
	Package
	Record
	Templ
	Type
	Using
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

func IsKeyword(ident string) (kind Kind, ok bool) {
	kind, ok = keywords[ident]
	return
}
