package token

import "fmt"

func (t Token) String() string {
	return t.Kind().String()
}

func (k Kind) String() string {
	switch k {
	case Invalid:
		return "INVALID"
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
	case EOL:
		return "\\n"
	case EOF:
		return "EOF"
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
	case Directive:
		return "directive"
	default:
		panic(fmt.Sprintf("unreachable: %#v", k))
	}
}
