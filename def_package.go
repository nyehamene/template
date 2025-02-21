package template

import "fmt"

func (p *Parser) defPackage(start int) (int, error) {
	ast := Def{}
	next := start

	if ident, n, err := p.expect(next, TokenIdent); err == nil {
		ast.left = ident
		next = n
	} else {
		// NOTE: for better error reporting ast.left could be set
		// to the erroneous token
		return start, err
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, ok := p.optional(next, TokenPackage); ok {
		next = n
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return start, err
	}

	if kind, n, err := p.packageTempl(next); err == nil {
		ast.kind = kind
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenSemicolon); err == nil {
		next = n
	} else {
		return start, err
	}

	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) defImport(start int) (int, error) {
	ast := Def{kind: DefImport}
	next := start

	if ident, n, err := p.expect(next, TokenIdent); err == nil {
		ast.left = ident
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, ok := p.optional(next, TokenImport); ok {
		next = n
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenImport); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenParLeft); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenString); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenParRight); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenSemicolon); err == nil {
		next = n
	} else {
		return start, err
	}

	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) defUsing(start int) (int, error) {
	ast := Def{kind: DefUsing}
	next := start

	if ident, n, err := p.expect(next, TokenIdent); err == nil {
		ast.left = ident
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, ok := p.optional(next, TokenImport); ok {
		next = n
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenUsing); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenIdent); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenSemicolon); err == nil {
		next = n
	} else {
		return start, err
	}

	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) packageTempl(start int) (DefKind, int, error) {
	next := start
	kind := DefTagPackage

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

func (p *Parser) packageTempl0(start int) (DefKind, int, bool) {
	handler := p.skipBefore(TokenSpace)(p.tokenizer.next)
	token, n, err := handler(start)
	if err != nil {
		return DefTagPackage, start, false
	}
	switch token.kind {
	case TokenTag:
		return DefTagPackage, n, true
	case TokenList:
		return DefListPackage, n, true
	case TokenHtml:
		return DefHtmlPackage, n, true
	default:
		return DefTagPackage, start, false
	}
}
