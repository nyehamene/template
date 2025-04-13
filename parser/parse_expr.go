package parser

import "temlang/tem/token"

func (p *Parser) baseexpr(start, end int) baseexpr {
	return baseexpr{start: start, end: end}
}

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
	case token.Type:
		f = p.parseTypeExpr
	case token.Record:
		f = p.parseRecordExpr
	case token.Templ:
		f = p.parseTemplExpr
	default:
		offset := p.identOffset()
		return p.badexpr(offset)
	}
	return f()
}

func (p *Parser) parsePackageExpr() Expr {
	var name token.Token
	offset := p.offset()

	if !p.expect(token.Package) {
		return p.badexpr(offset)
	}

	var ok bool
	if name, ok = p.expectSurroundParen(token.String); !ok {
		offset := p.offset()
		return p.badexpr(offset)
	}

	b := p.baseexpr(offset, p.prev.End())
	return pkgexpr{baseexpr: b, name: name}
}

func (p *Parser) parseImportExpr() Expr {
	offset := p.offset()
	if !p.expect(token.Import) {
		return p.badexpr(offset)
	}
	var path token.Token
	var ok bool
	if path, ok = p.expectSurroundParen(token.String); !ok {
		offset := p.offset()
		p.errorExpected("string")
		return p.badexpr(offset)
	}
	b := p.baseexpr(offset, p.prev.End())
	return importexpr{baseexpr: b, path: path}
}

func (p *Parser) parseUsingExpr() Expr {
	offset := p.offset()
	if !p.expect(token.Using) {
		return p.badexpr(offset)
	}
	var target token.Token
	var ok bool
	if target, ok = p.expectSurroundParen(token.Ident); !ok {
		p.errorExpected("ident")
		offset := p.offset()
		return p.badexpr(offset)
	}
	b := p.baseexpr(offset, p.prev.End())
	return usingexpr{baseexpr: b, target: target}
}

func (p *Parser) parseTypeExpr() Expr {
	offset := p.offset()
	if !p.expect(token.Type) {
		offset := p.offset()
		return p.badexpr(offset)
	}
	var target token.Token
	var ok bool
	if target, ok = p.expectSurroundParen(token.Ident); !ok {
		offset := p.offset()
		return p.badexpr(offset)
	}
	b := p.baseexpr(offset, p.prev.End())
	return typeexpr{baseexpr: b, target: target}
}

func (p *Parser) parseRecordExpr() Expr {
	offset := p.offset()
	if !p.expect(token.Record) {
		p.errorExpected("record")
		return p.badexpr(offset)
	}
	if !p.expect(token.BraceOpen) {
		p.errorExpected("{")
		return p.badexpr(p.offset())
	}

	var fields TreeQueue

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

	b := p.baseexpr(offset, p.prev.End())
	return recordexpr{baseexpr: b, fields: fields}
}

func (p *Parser) parseTemplExpr() Expr {
	offset := p.offset()
	if !p.expect(token.Templ) {
		p.errorExpected("templ")
		return p.badexpr(offset)
	}

	params := p.parseParamDecl()
	p.expect(token.BraceOpen)
	elements := p.parseElements()
	p.expect(token.BraceClose)

	b := p.baseexpr(offset, p.prev.End())
	return templexpr{baseexpr: b, params: params, elements: elements}
}
