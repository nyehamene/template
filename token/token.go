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
	Package
	Import
	Using
	Alias
	Templ
	Type
	Record
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

var whitespaces = map[byte]bool{
	' ':  true,
	'\t': true,
	'\r': true,
	'\v': true,
	'\f': true,
}

func IsSpace(b byte) bool {
	_, ok := whitespaces[b]
	return ok
}

func IsKeyword(ident string) bool {
	_, ok := keywords[ident]
	return ok
}
