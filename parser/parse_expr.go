package parser

import "temlang/tem/token"

type parseExprSpec func() Expr

func (p *Parser) parseGenExpr() Expr {
	var f parseExprSpec
	switch k := p.cur.Kind(); k {
	case token.Package:
		f = p.parsePackageExpr
	case token.Import:
		f = p.parseImportExpr
	case token.Using:
		f = p.parseUsingExpr
	default:
		return p.parseBasicExpr()
	}
	return f()
}

func (p *Parser) parsePackageExpr() Expr {
	var name token.Token

	switch kind := p.cur.Kind(); kind {
	case token.Package:
		p.advance()
		if !p.expectSurroundParen(token.String) {
			p.errorExpected(p.loc(), "package name")
			return p.badexpr()
		}
		name = p.prev
	default:
		return p.parseImportExpr()
	}

	return pkgexpr{name}
}

func (p *Parser) parseImportExpr() Expr {
	if !p.match(token.Import) {
		return p.parseUsingExpr()
	}
	if !p.expectSurroundParen(token.String) {
		p.errorExpected(p.loc(), "string")
		return p.badexpr()
	}
	return importexpr{p.prev}
}

func (p *Parser) parseUsingExpr() Expr {
	if !p.match(token.Using) {
		return p.parseBasicExpr()
	}
	if !p.expectSurroundParen(token.Ident) {
		p.errorExpected(p.loc(), "ident")
		return p.badexpr()
	}
	return usingexpr{p.prev}
}

func (p *Parser) parseBasicExpr() Expr {
	var f parseExprSpec
	switch k := p.cur.Kind(); k {
	case token.Type:
		f = p.parseTypeExpr
	case token.Record:
		f = p.parseRecordExpr
	case token.Templ:
		f = p.parseTemplExpr
	default:
		p.errorExpected(p.loc(), "an expr")
		return p.badexpr()
	}
	return f()
}

func (p *Parser) parseTypeExpr() Expr {
	if !p.expect(token.Type) {
		return p.badexpr()
	}
	if !p.expectSurroundParen(token.Ident) {
		return p.badexpr()
	}
	return typeexpr{p.prev}
}

func (p *Parser) parseRecordExpr() Expr {
	if !p.expect(token.Record) {
		p.errorExpected(p.loc(), "record")
		return p.badexpr()
	}
	if !p.expect(token.BraceOpen) {
		p.errorExpected(p.loc(), "{")
		return p.badexpr()
	}

	var fields TreeStack

	for p.cur.Kind() == token.Ident {
		field := p.parseDoc(p.parseVarDecl)
		fields.Push(field)
		p.expectSemicolon()
		if p.cur.Kind() == token.BraceClose {
			break
		}
	}
	if !p.expect(token.BraceClose) {
		p.errorExpected(p.loc(), "}")
		return p.badexpr()
	}

	return recordexpr{fields: fields} // TODO get record location
}

func (p *Parser) parseTemplExpr() Expr {
	if !p.expect(token.Templ) {
		p.errorExpected(p.loc(), "templ")
		return p.badexpr()
	}

	params := p.parseParamDecl()
	p.expect(token.BraceOpen)
	elements := p.parseElements()
	p.expect(token.BraceClose)

	return templexpr{params: params, elements: elements} // TODO get templ location
}
