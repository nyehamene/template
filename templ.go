package template

func (p Parser) defTempl(start int) (Def, int, error) {
	ast := Def{kind: DefTemplate}
	next := start

	if token, n, err := p.templDecl(next); err == nil {
		ast.left = token
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.templModel(next); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenBraceLeft); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	if _, n, err := p.expect(next, TokenBraceRight); err == nil {
		next = n
	} else {
		return ast, start, err
	}

	return ast, next, nil
}

func (p Parser) templModel(start int) (Token, int, error) {
	var leftPar Token
	next := start
	token, n, err := p.expect(next, TokenParLeft)
	if err != nil {
		return token, start, err
	}
	next = n
	leftPar = token

	token, n, err = p.expect(next, TokenIdent)
	if err != nil {
		return token, start, err
	}
	next = n

	token, n, err = p.expect(next, TokenParRight)
	if err != nil {
		return token, start, err
	}
	next = n

	return leftPar, next, nil
}

func (p Parser) templDecl(start int) (Token, int, error) {
	return p.decl(start, TokenTempl)
}
