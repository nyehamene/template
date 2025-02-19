package template

import "fmt"

func NewParser(t Tokenizer) Parser {
	if t.source == nil {
		panic("source must not be nil")
	}
	return Parser{t}
}

type parsehandler func(int) (Token, int, error)

type parsematcher func(parsehandler) parsehandler

type AstKind int

const (
	AstTagTemplPackage AstKind = iota
	AstListTemplPackage
	AstHtmlTemplPackage
	AstRecord
	AstAlias
	AstTemplate
	AstDocline
	AstDocblock
	AstImport
	AstUsing
	AstMetatable
)

type Ast struct {
	kind AstKind
	left Token
}

type Parser struct {
	tokenizer Tokenizer
}

func (p Parser) Parse(start int) (Ast, int, error) {
	panic("unimplemented")
}

func (p Parser) parseDef(start int) (Ast, int, error) {
	var ast Ast
	next := start

	if def, n, err := p.doc(next); err == nil {
		ast = def
		next = n
	} else if def, n, err := p.typeDef(next); err == nil {
		ast = def
		next = n
	} else if def, n, err := p.templDef(next); err == nil {
		ast = def
		next = n
	} else {
		return def, start, err
	}

	if _, n, err := p.expect(next, TokenSemicolon); err != nil {
		return ast, start, err
	} else {
		next = n
	}

	return ast, next, nil
}

func (p Parser) doc(start int) (Ast, int, error) {
	ast := Ast{}
	next := start

	if token, n, err := p.expect(next, TokenIdent); err == nil {
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

	if kind, n, err := p.docString(next); err == nil {
		ast.kind = kind
		next = n
	} else {
		return ast, start, err
	}

	return ast, next, nil
}

func (p Parser) docString(start int) (AstKind, int, error) {
	if _, n, err := p.expect(start, TokenString); err == nil {
		return AstDocline, n, nil
	} else if _, n, err := p.expect(start, TokenTextBlock); err == nil {
		return AstDocblock, n, nil
	}
	return AstDocline, start, ErrInvalid
}

func (p Parser) decl(start int, kind TokenKind) (Token, int, error) {
	var ident Token
	next := start

	if token, n, err := p.expect(next, TokenIdent); err == nil {
		ident = token
		next = n
	} else {
		return token, start, err
	}

	if token, n, err := p.expect(next, TokenColon); err == nil {
		next = n
	} else {
		return token, start, err
	}

	if _, n, ok := p.optional(next, kind); ok {
		next = n
	}
	return ident, next, nil
}

func (p Parser) optional(start int, kind TokenKind, matchers ...parsematcher) (Token, int, bool) {
	t, next, err := p.expect(start, kind, matchers...)
	if err != nil {
		return t, start, false
	}

	if t.kind != kind {
		return t, start, false
	}

	return t, next, true
}

func (p Parser) expect(start int, kind TokenKind, matchers ...parsematcher) (Token, int, error) {
	var main parsehandler = func(start int) (Token, int, error) {
		token, next, err := p.tokenizer.next(start)
		if err != nil {
			return token, start, err
		}
		if token.kind != kind {
			r, c := p.tokenizer.Pos(token)
			return token, start, fmt.Errorf("invalid token %s at [%d, %d]\nHelp: %s", token.kind, r, c, kind)
		}
		return token, next, nil
	}

	var handler parsehandler = main
	for _, matcher := range matchers {
		handler = matcher(handler)
	}

	handler = p.skipBefore(TokenSpace, TokenComment, TokenEOL)(handler)

	return handler(start)
}

func (p Parser) skipBefore(kind TokenKind, more ...TokenKind) parsematcher {
	return func(h parsehandler) parsehandler {
		return func(start int) (Token, int, error) {
			kinds := make([]TokenKind, 0, len(more)+1)
			kinds = append(kinds, kind)
			kinds = append(kinds, more...)
			var beforeToken = start
			for {
				toSkip, next, err := p.tokenizer.next(beforeToken)
				if err != nil {
					return toSkip, start, err
				}
				for _, kind := range kinds {
					if toSkip.kind == kind {
						goto skip
					}
				}
				break
			skip:
				beforeToken = next
			}
			return h(beforeToken)
		}
	}
}
