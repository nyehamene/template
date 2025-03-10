package parser

func (p *Parser) ParseFile() {
	p.pkg()
}

func (p *Parser) pkg() {
	reset := p.Mark()
	_, ok := p.parsePackageDecl()
	if !ok {
		reset()
		p.errorf("namespace must have a package declaration as the first none comment decl")
	}
	p.import_()
}

func (p *Parser) import_() {
	for {
		reset := p.Mark()
		_, ok := p.parseImportDecl()
		if !ok {
			reset()
			break
		}
	}
	p.using()
}

func (p *Parser) using() {
	for {
		reset := p.Mark()
		_, ok := p.parseUsingDecl()
		if !ok {
			reset()
			break
		}
	}
	p.decl()
}

func (p *Parser) decl() {
	for {
		reset := p.Mark()
		ok := p.basicDecl()
		if !ok {
			reset()
			break
		}
	}
}

func (p *Parser) basicDecl() bool {
	return p.doc()
}

func (p *Parser) doc() bool {
	reset := p.Mark()
	_, ok := p.parseDocDecl()
	if ok {
		return p.doc()
	}
	reset()
	return p.tag()
}

func (p *Parser) tag() bool {
	reset := p.Mark()
	_, ok := p.parseTagDecl()
	if ok {
		return p.doc()
	}
	reset()
	return p.mainDecl()
}

func (p *Parser) mainDecl() (ok bool) {
	reset := p.Mark()
	_, ok = p.parseAliasDecl()
	if ok {
		return
	}
	reset()

	reset = p.Mark()
	_, ok = p.parseRecordDecl()
	if ok {
		return
	}
	reset()

	reset = p.Mark()
	_, ok = p.parseTemplDecl()
	if !ok {
		reset()
	}
	return
}
