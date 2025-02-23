package template

func (p *Parser) metatable(start int) (int, error) {
	var ident Token
	var next int
	var ok bool

	if ident, next, ok = p.match(start, TokenIdent); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenBraceLeft); !ok {
		return start, ErrNoMatch
	}

	// The parser is in a metatable declaration. Any error
	// from this point on means the declaration is invalid.
	// therefore, the parse should not try to parse another
	// kind of declaration at the same position

	if next, ok = p.attr(next); !ok {
		return start, ErrInvalid
	}

	if _, n, ok := p.match(next, TokenComma); ok {
		offset := n

		for {
			if offset, ok = p.attr(offset); !ok {
				break
			}

			if _, offset, ok = p.match(offset, TokenComma); !ok {
				break
			}
		}

		next = offset
	}

	if _, next, ok = p.match(next, TokenBraceRight); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenSemicolon); !ok {
		return start, ErrInvalid
	}

	kind := DefMetatable
	ast := Def{kind: kind, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) metatablePackage(start int) (int, error) {
	var next int
	var err error

	if next, err = p.defPackage(start); err != nil {
		return start, err
	}

	next, err = p.repeat(next, p.metatable)
	return next, err
}

func (p *Parser) metatableDef(start int) (int, error) {
	var next int
	var err error

	next, err = p.repeat(start, p.metatable)
	if err == ErrInvalid {
		return next, err
	}
	if err == EOF {
		return next, EOF
	}

	next, err = p.defTypeOrTempl(next)
	if err != nil {
		return start, err
	}

	next, err = p.repeat(next, p.metatable)
	return next, err
}

func (p *Parser) metatableRecordComp(start int) (int, error) {
	var next int
	var err error

	next, err = p.repeat(start, p.metatable)
	if err == ErrInvalid {
		return next, err
	}
	if err == EOF {
		return next, EOF
	}

	for {
		var ok bool
		if _, next, ok = p.match(next, TokenIdent); !ok {
			break
		}
		if _, next, ok = p.match(next, TokenColon); !ok {
			return start, ErrNoMatch
		}
		if _, next, ok = p.match(next, TokenIdent); !ok {
			return start, ErrInvalid
		}
		if _, next, ok = p.match(next, TokenSemicolon); !ok {
			break
		}
	}

	next, err = p.repeat(next, p.metatable)
	return next, err
}

func (p *Parser) attr(start int) (int, bool) {
	var next int
	var ok bool

	if _, next, ok = p.match(start, TokenIdent); !ok {
		return start, false
	}

	if _, next, ok = p.match(next, TokenEqual); !ok {
		return start, false
	}

	if _, next, ok = p.match(next, TokenString); !ok {
		return start, false
	}

	return next, true
}
