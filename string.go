package template

import "fmt"

func (t Token) String() string {
	return fmt.Sprintf("%s %d", t.kind, t.offset)
}

func (t TokenKind) String() string {
	switch t {
	case TokenUndefined:
		return "<>"
	case TokenColon:
		return ":"
	case TokenEqual:
		return "="
	case TokenPeriod:
		return "."
	case TokenSemicolon:
		return ";"
	case TokenBraceLeft:
		return "{"
	case TokenBraceRight:
		return "}"
	case TokenBracketLeft:
		return "{"
	case TokenBracketRight:
		return "}"
	case TokenParLeft:
		return "("
	case TokenParRight:
		return ")"
	case TokenSpace:
		return ":spc"
	case TokenEOL:
		return ":eol"
	case TokenIdent:
		return ":ident"
	case TokenPackage:
		return ":package"
	case TokenTag:
		return ":package_tag"
	case TokenList:
		return ":package_list"
	case TokenHtml:
		return ":package_html"
	case TokenType:
		return ":type"
	case TokenTempl:
		return ":templ"
	case TokenEnd:
		return ":end"
	case TokenString:
		return ":str"
	case TokenAlias:
		return ":alias"
	case TokenComment:
		return ":comment"
	default:
		panic("unreachable")
	}
}

func (k AstKind) String() string {
	switch k {
	case AstPackage:
		return ":package"
	case AstTemplateTag:
		return ":tag"
	case AstTemplateList:
		return ":list"
	case AstTemplateHtml:
		return ":html"
	case AstTypeIdent:
		return ":type_ident"
	case AstTypeDef:
		return ":type"
	case AstRecordDef:
		return ":record"
	case AstAliasDef:
		return ":alias"
	case AstIdent:
		return ":ident"
	default:
		panic("unreachable")
	}
}
