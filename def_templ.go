package template

func (p *Parser) defTempl(start int) (next int, err error) {
	var ident Token
	var ok bool

	if ident, next, ok = p.templDecl(next); !ok {
		return start, ErrNoMatch
	}

	// At this point the parsing is on a template declaration
	// therefore, if there is an error at this point on the parse
	// should not try to match another type of declaration at the
	// current position

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.templModel(next); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenBraceLeft); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenBraceRight); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenSemicolon); !ok {
		return start, ErrInvalid
	}

	ast := Def{kind: DefTemplate, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) templModel(start int) (Token, int, bool) {
	var token Token
	var next int
	var ok bool

	if token, next, ok = p.match(next, TokenParLeft); !ok {
		return token, start, false
	}

	if _, next, ok = p.match(next, TokenIdent); !ok {
		return token, start, false
	}

	if _, next, ok = p.match(next, TokenParRight); ok {
		return token, start, false
	}

	return token, next, true
}

func (p *Parser) templDecl(start int) (Token, int, bool) {
	return p.decl(start, TokenTempl)
}
