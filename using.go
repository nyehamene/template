package template

func (p Parser) parseUsing(start int) (Ast, int, error) {
	ast := Ast{kind: AstUsing}
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

	if _, n, err := p.expect(next, TokenUsing); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenIdent); err == nil {
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
