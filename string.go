package template

import "fmt"

func (t Token) String() string {
	return fmt.Sprintf("%s %d", t.kind, t.offset)
}

func (k TokenKind) String() string {
	switch k {
	case TokenUndefined:
		return "<>"
	case TokenComma:
		return ","
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
		return "spc"
	case TokenEOL:
		return "eol"
	case TokenIdent:
		return "ident"
	case TokenPackage:
		return "package"
	case TokenTag:
		return "package_tag"
	case TokenList:
		return "package_list"
	case TokenHtml:
		return "package_html"
	case TokenType:
		return "type"
	case TokenTempl:
		return "templ"
	case TokenEnd:
		return "end"
	case TokenString:
		return "str"
	case TokenAlias:
		return "alias"
	case TokenComment:
		return "comment"
	case TokenRecord:
		return "record"
	case TokenTextBlock:
		return "text_block"
	case TokenImport:
		return "import"
	case TokenUsing:
		return "using"
	default:
		panic(fmt.Sprintf("unreachable: %#v", k))
	}
}

func (k DefKind) String() string {
	switch k {
	case DefTagPackage:
		return "tag"
	case DefListPackage:
		return "list"
	case DefHtmlPackage:
		return "html"
	case DefRecord:
		return "record"
	case DefAlias:
		return "alias"
	case DefTemplate:
		return "templ"
	case DefDocline:
		return "docline"
	case DefDocblock:
		return "docblock"
	case DefImport:
		return "import"
	case DefUsing:
		return "using"
	case DefMetatable:
		return "metatable"
	default:
		panic(fmt.Sprintf("unreachable: %#v", k))
	}
}

func (a Def) String() string {
	return fmt.Sprintf("%s %s", a.kind, a.left)
}
