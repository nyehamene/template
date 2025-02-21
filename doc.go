package template

func (p *Parser) doc(start int) (int, error) {
	def := Def{}
	next := start

	if token, n, err := p.expect(next, TokenIdent); err == nil {
		def.left = token
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return start, err
	}

	if kind, n, err := p.docString(next); err == nil {
		def.kind = kind
		next = n
	} else {
		return start, err
	}

	if _, n, err := p.expect(next, TokenSemicolon); err != nil {
		return start, err
	} else {
		next = n
	}

	p.ast = append(p.ast, def)
	return next, nil
}

func (p *Parser) docString(start int) (DefKind, int, error) {
	if _, n, err := p.expect(start, TokenString); err == nil {
		return DefDocline, n, nil
	} else if _, n, err := p.expect(start, TokenTextBlock); err == nil {
		return DefDocblock, n, nil
	}
	return DefDocline, start, ErrInvalid
}

func (p *Parser) docPackage(start int) (int, error) {
	next := start
	if n, err := p.metatablePackage(next); err == nil {
		next = n
	} else {
		return start, err
	}
	next = p.repeat(next, p.doc)
	return next, nil
}

func (p *Parser) docImport(start int) (int, error) {
	next := p.repeat(start, p.doc)

	if n, err := p.defImport(next); err == nil {
		next = n
	} else {
		return start, err
	}

	next = p.repeat(next, p.doc)
	return next, nil
}

func (p *Parser) docUsing(start int) (int, error) {
	next := p.repeat(start, p.doc)

	if n, err := p.defUsing(next); err == nil {
		next = n
	} else {
		return start, err
	}

	next = p.repeat(next, p.doc)
	return next, nil
}

func (p *Parser) docDef(start int) (int, error) {
	next := p.repeat(start, p.doc)

	if n, err := p.metatableDef(next); err == nil {
		next = n
	} else {
		return start, err
	}

	next = p.repeat(next, p.doc)
	return next, nil
}

func (p *Parser) docRecordComp(start int) (int, error) {
	next := p.repeat(start, p.doc)

	if n, err := p.metatableRecordComp(next); err == nil {
		next = n
	} else {
		return start, err
	}

	next = p.repeat(next, p.doc)
	return next, nil
}
