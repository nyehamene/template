package template

import "fmt"

func (p Parser) defAlias(start int) (DefKind, int, bool) {
	next := start

	if _, n, err := p.expect(next, TokenAlias); err == nil {
		next = n
	} else {
		return DefAlias, start, false
	}

	if _, n, err := p.expect(next, TokenIdent); err == nil {
		next = n
	} else {
		return DefAlias, start, false
	}

	return DefAlias, next, true
}

func (p Parser) defType(start int) (Def, int, error) {
	ast := Def{}
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

	if kind, n, ok := p.defRecord(next); ok {
		ast.kind = kind
		next = n
	} else if kind, n, ok := p.defAlias(next); ok {
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
