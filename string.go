package template

import "fmt"

func (t Token) String() string {
	return fmt.Sprintf("%s %d", t.kind, t.offset)
}

func (t TokenKind) String() string {
	str := ""
	switch t {
	case TokenUndefined:
		str = "<>"
	case TokenColon:
		str = ":"
	case TokenEqual:
		str = "="
	case TokenPeriod:
		str = "."
	case TokenSemicolon:
		str = ";"
	case TokenBraceLeft:
		str = "{"
	case TokenBraceRight:
		str = "}"
	case TokenBracketLeft:
		str = "{"
	case TokenBracketRight:
		str = "}"
	case TokenParLeft:
		str = "("
	case TokenParRight:
		str = ")"
	case TokenSpace:
		str = ":spc"
	case TokenEOL:
		str = ":eol"
	case TokenIdent:
		str = ":ident"
	case TokenPackage:
		str = ":package"
	case TokenTag:
		str = ":tag"
	case TokenList:
		str = ":list"
	case TokenHtml:
		str = ":html"
	case TokenType:
		str = ":type"
	case TokenTempl:
		str = ":templ"
	case TokenEnd:
		str = ":end"
	case TokenString:
		str = ":str"
	case TokenComment:
		str = ":comment"
	default:
		panic("unreachable")
	}
	return str
}

func (k AstKind) String() string {
	switch k {
	case AstPackage:
		return ":package"
	case AstTag:
		return ":tag"
	case AstIdent:
		return ":ident"
	default:
		panic("unreachable")
	}
}
