package token

import "fmt"

func (t Token) String() string {
	return fmt.Sprintf("%s %d", t.Kind, t.Offset)
}

func (k Kind) String() string {
	switch k {
	case Invalid:
		return "?invalid"
	case Comma:
		return ","
	case Colon:
		return ":"
	case Eq:
		return "="
	case Dot:
		return "."
	case Semicolon:
		return ";"
	case BracketOpen:
		return "["
	case BracketClose:
		return "]"
	case BraceOpen:
		return "{"
	case BraceClose:
		return "}"
	case ParenOpen:
		return "("
	case ParenClose:
		return ")"
	case Space:
		return "spc"
	case EOL:
		return "eol"
	case Ident:
		return "ident"
	case Package:
		return "package"
	case Type:
		return "type"
	case Templ:
		return "templ"
	case String:
		return "str"
	case Alias:
		return "alias"
	case Comment:
		return "comment"
	case Record:
		return "record"
	case TextBlock:
		return "text_block"
	case Import:
		return "import"
	case Using:
		return "using"
	default:
		panic(fmt.Sprintf("unreachable: %#v", k))
	}
}
