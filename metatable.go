package template

func (p Parser) metatable(start int) (Ast, int, error) {
	ast := Ast{kind: AstMetatable}
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

	if _, n, err := p.expect(next, TokenBraceLeft); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.attr(next); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, ok := p.optional(next, TokenComma); ok {
		offset := n

		for {
			if _, n, err := p.attr(offset); err == nil {
				offset = n
			} else {
				break
			}

			if _, n, err := p.expect(offset, TokenComma); err == nil {
				offset = n
			} else {
				break
			}
		}

		next = offset
	}

	if _, n, err := p.expect(next, TokenBraceRight); err == nil {
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

func (p Parser) attr(start int) (AstKind, int, error) {
	kind := AstMetatable
	next := start

	if _, n, err := p.expect(next, TokenIdent); err == nil {
		next = n
	} else {
		return kind, start, err
	}

	if _, n, err := p.expect(next, TokenEqual); err == nil {
		next = n
	} else {
		return kind, start, err
	}

	if _, n, err := p.expect(next, TokenString); err == nil {
		next = n
	} else {
		return kind, start, err
	}

	return kind, next, nil
}
