package template

import "log"

func (p *Parser) defType(start int) (next int, err error) {
	var ident Token
	var ok bool

	if ident, next, ok = p.typeDecl(start); !ok {
		return start, ErrNoMatch
	}

	// At this point the parsing is in a type declaration
	// therefore, any error from this point on means the
	// declaration is invalid. And the parse should not try
	// to parse another kind of token at the current position

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrInvalid
	}

	if next, err = p.defRecord(next, ident); err == nil {
		goto matchSemicolon
	}

	if next, err = p.defAlias(next, ident); err != nil {
		return start, ErrInvalid
	}

matchSemicolon:
	if _, next, ok = p.match(next, TokenSemicolon); !ok {
		return start, ErrInvalid
	}

	return next, nil
}

func (p *Parser) defAlias(start int, ident Token) (int, error) {
	var next int
	var ok bool

	if _, next, ok = p.match(start, TokenAlias); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenIdent); !ok {
		return start, ErrInvalid
	}

	ast := Def{kind: DefAlias, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) defRecord(start int, ident Token) (int, error) {
	var next int
	var ok bool
	var err error

	if _, next, ok = p.match(start, TokenRecord); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenBraceLeft); !ok {
		return start, ErrInvalid
	}

	if next, err = p.docRecordComp(next); err != nil {
		// TODO: wrap err in ErrInvalid
		log.Println(err)
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenBraceRight); !ok {
		return start, ErrInvalid
	}

	ast := Def{kind: DefRecord, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) typeDecl(start int) (Token, int, bool) {
	return p.decl(start, TokenType)
}
