package template

import "fmt"

func (p Parser) parsePackage(start int) (Ast, int, error) {
	ast := Ast{}
	next := start

	if token, n, err := p.expect(next, TokenIdent); err == nil {
		ast.left = token
		next = n
	} else {
		// NOTE: for better error reporting ast.left could be set
		// to the erroneous token
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, ok := p.optional(next, TokenPackage); ok {
		next = n
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if kind, n, err := p.packageTempl(next); err == nil {
		ast.kind = kind
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenSemicolon); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	return ast, next, nil
}

func (p Parser) parseImport(start int) (Ast, int, error) {
	ast := Ast{kind: AstImport}
	next := start

	if ident, n, err := p.expect(next, TokenIdent); err == nil {
		ast.left = ident
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, ok := p.optional(next, TokenImport); ok {
		next = n
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenImport); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenParLeft); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenString); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenParRight); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenSemicolon); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	return ast, next, nil
}

func (p Parser) packageTempl(start int) (AstKind, int, error) {
	next := start
	kind := AstTagTemplPackage

	if k, n, ok := p.packageTempl0(next); ok {
		kind = k
		next = n
	} else {
		return kind, start, fmt.Errorf("invalid package def")
	}

	if _, n, err := p.expect(next, TokenParLeft); err == nil {
		next = n
	} else {
		return kind, start, err
	}

	if _, n, err := p.expect(next, TokenString); err == nil {
		next = n
	} else {
		return kind, start, err
	}

	if _, n, err := p.expect(next, TokenParRight); err == nil {
		next = n
	} else {
		return kind, start, err
	}

	return kind, next, nil
}

func (p Parser) packageTempl0(start int) (AstKind, int, bool) {
	handler := p.skipBefore(TokenSpace)(p.tokenizer.next)
	token, n, err := handler(start)
	if err != nil {
		return AstTagTemplPackage, start, false
	}
	switch token.kind {
	case TokenTag:
		return AstTagTemplPackage, n, true
	case TokenList:
		return AstListTemplPackage, n, true
	case TokenHtml:
		return AstHtmlTemplPackage, n, true
	default:
		return AstTagTemplPackage, start, false
	}
}
