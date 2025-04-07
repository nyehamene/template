package parser

import "temlang/tem/token"

type parseExprSpec func(offset int) Expr

func (p *Parser) parseGenExpr(offset int) Expr {
	var f parseExprSpec
	switch k := p.cur.Kind(); k {
	case token.Package:
		f = p.parsePackageExpr
	case token.Import:
		f = p.parseImportExpr
	case token.Using:
		f = p.parseUsingExpr
	default:
		return p.parseBasicExpr(offset)
	}
	return f(offset)
}

func (p *Parser) parsePackageExpr(offset int) Expr {
	var name token.Token

	switch kind := p.cur.Kind(); kind {
	case token.Package:
		p.advance()
		if !p.expectSurroundParen(token.String) {
			p.errorExpected("package name")
			return p.badexpr(offset)
		}
		name = p.prev
	default:
		return p.parseImportExpr(offset)
	}

	return pkgexpr{name}
}

func (p *Parser) parseImportExpr(offset int) Expr {
	if !p.match(token.Import) {
		return p.parseUsingExpr(offset)
	}
	if !p.expectSurroundParen(token.String) {
		p.errorExpected("string")
		return p.badexpr(offset)
	}
	return importexpr{p.prev}
}

func (p *Parser) parseUsingExpr(offset int) Expr {
	if !p.match(token.Using) {
		return p.parseBasicExpr(offset)
	}
	if !p.expectSurroundParen(token.Ident) {
		p.errorExpected("ident")
		return p.badexpr(offset)
	}
	return usingexpr{p.prev}
}

func (p *Parser) parseBasicExpr(offset int) Expr {
	var f parseExprSpec
	switch k := p.cur.Kind(); k {
	case token.Type:
		f = p.parseTypeExpr
	case token.Record:
		f = p.parseRecordExpr
	case token.Templ:
		f = p.parseTemplExpr
	default:
		p.errorExpected("an expr")
		return p.badexpr(offset)
	}
	return f(offset)
}

func (p *Parser) parseTypeExpr(offset int) Expr {
	if !p.expect(token.Type) {
		return p.badexpr(offset)
	}
	if !p.expectSurroundParen(token.Ident) {
		return p.badexpr(offset)
	}
	return typeexpr{p.prev}
}

func (p *Parser) parseRecordExpr(offset int) Expr {
	if !p.expect(token.Record) {
		p.errorExpected("record")
		return p.badexpr(offset)
	}
	if !p.expect(token.BraceOpen) {
		p.errorExpected("{")
		return p.badexpr(offset)
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
		p.errorExpected("}")
		return p.badexpr(offset)
	}

	return recordexpr{fields: fields}
}

func (p *Parser) parseTemplExpr(offset int) Expr {
	if !p.expect(token.Templ) {
		p.errorExpected("templ")
		return p.badexpr(offset)
	}

	params := p.parseParamDecl()
	p.expect(token.BraceOpen)
	elements := p.parseElements()
	p.expect(token.BraceClose)

	return templexpr{params: params, elements: elements}
}
