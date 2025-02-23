package template

import "log"

func (p *Parser) defPackage(start int) (next int, err error) {
	var ident Token
	var ok bool

	if ident, next, ok = p.match(start, TokenIdent); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	// optional type
	var okTokenPackage bool
	_, next, okTokenPackage = p.match(next, TokenPackage)

	// If the type is a package, and there is an error that means
	// the package declaration is invalid, therefore, the parse should not
	// try to match other types of declarations in the same position.
	defer func() {
		if err != nil && okTokenPackage {
			err = ErrInvalid
		}
	}()

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	var kind DefKind
	if kind, next, ok = p.packageTempl(next); !ok {
		return start, ErrNoMatch
	}

	// if we successfully parse a package templ type then we are in a package
	// declaration
	okTokenPackage = true

	if _, next, ok = p.match(next, TokenSemicolon); !ok {
		return start, ErrInvalid
	}

	ast := Def{kind: kind, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) defImport(start int) (next int, err error) {
	var ident Token
	var ok bool

	if ident, next, ok = p.match(start, TokenIdent); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	// optional type
	var okTokenImport bool
	_, next, okTokenImport = p.match(next, TokenImport)

	// If there is an error, that means the import declaration
	// is invalid, therefore the parse should not try to parse
	// another type of declaration in the same position.
	defer func() {
		if err != nil && okTokenImport {
			err = ErrInvalid
		}
	}()

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenImport); !ok {
		return start, ErrNoMatch
	}

	// If the import keyword is matched at this point
	// then we are in an import declaration
	okTokenImport = true

	if _, next, ok = p.match(next, TokenParLeft); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenString); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenParRight); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenSemicolon); !ok {
		return start, ErrInvalid
	}

	ast := Def{kind: DefImport, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) defUsing(start int) (next int, err error) {
	var ident Token
	var ok bool

	if ident, next, ok = p.match(start, TokenIdent); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	// optional type
	var okTokenImport bool
	_, next, okTokenImport = p.match(next, TokenImport)

	// If there is an error, that means the import declaration
	// is invalid, therefore the parse should not try to parse
	// another type of declaration in the same position.
	defer func() {
		if err != nil && okTokenImport {
			err = ErrInvalid
		}
	}()

	if _, next, ok = p.match(next, TokenColon); !ok {
		return start, ErrNoMatch
	}

	if _, next, ok = p.match(next, TokenUsing); !ok {
		return start, ErrNoMatch
	}

	// At this point the parse is parsing a using declaration
	// therefore, any errors from here on means the declaration
	// is invalid

	if _, next, ok = p.match(next, TokenIdent); !ok {
		return start, ErrInvalid
	}

	if _, next, ok = p.match(next, TokenSemicolon); !ok {
		return start, ErrInvalid
	}

	ast := Def{kind: DefUsing, left: ident}
	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) packageTempl(start int) (DefKind, int, bool) {
	var next int
	var ok bool
	var kind DefKind

	if kind, next, ok = p.packageTempl0(start); !ok {
		return kind, start, false
	}

	if _, next, ok = p.match(next, TokenParLeft); !ok {
		return kind, start, false
	}

	if _, next, ok = p.match(next, TokenString); !ok {
		return kind, start, false
	}

	if _, next, ok = p.match(next, TokenParRight); !ok {
		return kind, start, false
	}

	return kind, next, true
}

func (p *Parser) packageTempl0(start int) (DefKind, int, bool) {
	handler := func(s int) (Token, int, bool) {
		t, n, err := p.tokenizer.next(s)
		if err != nil {
			log.Println(err)
		}
		return t, n, err == nil
	}

	handler = p.skipBefore(TokenSpace)(handler)
	token, n, ok := handler(start)
	if !ok {
		return DefTagPackage, start, false
	}
	switch token.kind {
	case TokenTag:
		return DefTagPackage, n, true
	case TokenList:
		return DefListPackage, n, true
	case TokenHtml:
		return DefHtmlPackage, n, true
	default:
		return DefTagPackage, start, false
	}
}
