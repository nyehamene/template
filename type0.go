package template

import "fmt"

func (p Parser) aliasDef(start int) (AstKind, int, bool) {
	next := start

	if _, n, err := p.expect(next, TokenAlias); err == nil {
		next = n
	} else {
		return AstAlias, start, false
	}

	if _, n, err := p.expect(next, TokenIdent); err == nil {
		next = n
	} else {
		return AstAlias, start, false
	}

	return AstAlias, next, true
}

func (p Parser) typeDef(start int) (Ast, int, error) {
	ast := Ast{}
	next := start

	if token, n, err := p.typeDecl(next); err == nil {
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

	if kind, n, ok := p.recordDef(next); ok {
		ast.kind = kind
		next = n
	} else if kind, n, ok := p.aliasDef(next); ok {
		ast.kind = kind
		next = n
	} else {
		// TODO: get the offset of the next none space token
		return ast, start, fmt.Errorf("invalid type def")
	}

	return ast, next, nil
}

func (p Parser) typeDecl(start int) (Token, int, error) {
	return p.decl(start, TokenType)
}
