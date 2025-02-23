package template

func (p *Parser) doc(start int) (int, error) {
	var ident Token
	var kind DefKind
	var next int
	var ok bool

	if ident, next, ok = p.match(start, TokenIdent); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	if kind, next, ok = p.docString(next); !ok {
		return start, ErrNoMatch
	}

	// The parser is in a documentation at this point
	// if there is any further errors that means the documentation
	// is invalid and the should not try to match another kind of
	// declaration in the current position

	if _, next, ok = p.match(next, TokenSemicolon); !ok {
		return start, ErrInvalid
	}

	ast := Def{kind: kind, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) docString(start int) (DefKind, int, bool) {
	if _, n, ok := p.match(start, TokenString); ok {
		return DefDocline, n, true
	}
	if _, n, ok := p.match(start, TokenTextBlock); ok {
		return DefDocblock, n, true
	}
	return DefDocline, start, false
}

func (p *Parser) docPackage(start int) (int, error) {
	var next int
	var err error

	if next, err = p.metatablePackage(start); err != nil {
		return start, err
	}
	next, err = p.repeat(next, p.doc)
	return next, err
}

func (p *Parser) docImport(start int) (int, error) {
	var next int
	var err error

	next, err = p.repeat(start, p.doc)
	if err == ErrInvalid {
		return next, err
	}
	if err == EOF {
		return next, EOF
	}

	if next, err = p.defImport(next); err != nil {
		return start, err
	}

	next, err = p.repeat(next, p.doc)
	return next, err
}

func (p *Parser) docUsing(start int) (int, error) {
	var next int
	var err error

	next, err = p.repeat(start, p.doc)
	if err == ErrInvalid {
		return next, err
	}
	if err == EOF {
		return next, EOF
	}

	if next, err = p.defUsing(next); err != nil {
		return start, err
	}

	next, err = p.repeat(next, p.doc)
	return next, err
}

func (p *Parser) docDef(start int) (int, error) {
	var next int
	var err error

	next, err = p.repeat(start, p.doc)
	if err == ErrInvalid {
		return next, err
	}
	if err == EOF {
		return next, EOF
	}

	if next, err = p.metatableDef(next); err != nil {
		return start, err
	}

	next, err = p.repeat(next, p.doc)
	return next, err
}

func (p *Parser) docRecordComp(start int) (int, error) {
	var next int
	var err error

	next, err = p.repeat(start, p.doc)
	if err == ErrInvalid {
		return next, err
	}
	if err == EOF {
		return next, EOF
	}

	if next, err = p.metatableRecordComp(next); err != nil {
		return start, err
	}

	next, err = p.repeat(next, p.doc)
	return next, err
}
