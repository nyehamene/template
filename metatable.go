package template

func (p *Parser) metatable(start int) (int, error) {
	ast := Def{kind: DefMetatable}
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

	if _, n, err := p.expect(next, TokenBraceLeft); err == nil {
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.attr(next); err == nil {
		next = n
	} else {
		return start, err
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

func (p *Parser) metatablePackage(start int) (int, error) {
	next := start
	if n, err := p.defPackage(next); err == nil {
		next = n
	} else {
		return start, err
	}

	next = p.repeat(next, p.metatable)

	return next, nil
}

func (p *Parser) metatableDef(start int) (int, error) {
	next := p.repeat(start, p.metatable)

	if n, err := p.defTypeOrTempl(next); err == nil {
		next = n
	} else {
		return start, err
	}

	next = p.repeat(next, p.metatable)
	return next, nil
}

func (p *Parser) metatableRecordComp(start int) (int, error) {
	next := p.repeat(start, p.metatable)

	for {
		var err error
		if _, next, err = p.expect(next, TokenIdent); err != nil {
			break
		}
		if _, next, err = p.expect(next, TokenColon); err != nil {
			return start, err
		}
		if _, next, err = p.expect(next, TokenIdent); err != nil {
			return start, err
		}
		if _, next, err = p.expect(next, TokenSemicolon); err != nil {
			break
		}
	}

	next = p.repeat(next, p.metatable)
	return next, nil
}

func (p *Parser) attr(start int) (DefKind, int, error) {
	kind := DefMetatable
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
