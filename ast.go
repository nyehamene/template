package template

import (
	"fmt"
)

type Def struct {
	kind DefKind
	left Token
}

type DefKind int

const (
	DefInvalid DefKind = iota
	DefTagPackage
	DefListPackage
	DefHtmlPackage
	DefRecord
	DefAlias
	DefTemplate
	DefDocline
	DefDocblock
	DefImport
	DefUsing
	DefMetatable
)

func NewParser(t *Tokenizer) Parser {
	if t.source == nil {
		panic("source must not be nil")
	}
	return Parser{t, []Def{}}
}

type Parser struct {
	tokenizer *Tokenizer
	ast       []Def
}

func (p *Parser) Parse(start int) ([]Def, error) {
	next := start
	if n, err := p.docPackage(next); err == nil {
		next = n
	} else {
		return p.ast, err
	}

	next = p.repeat(next, p.docImport)
	next = p.repeat(next, p.docUsing)
	next = p.repeat(next, p.docDef)

	if err := p.expectErr(next, EOF); err != nil {
		return p.ast, err
	}

	ast := p.ast
	p.ast = nil

	return ast, nil
}

func (p *Parser) defTypeOrTempl(start int) (next int, err error) {
	var ast Def

	if ast, next, err = p.defType(start); err != nil {

		if err == EOF {
			return start, EOF
		}

		if ast, next, err = p.defTempl(start); err != nil {
			return start, ErrInvalid
		}
	}

	p.ast = append(p.ast, ast)
	return next, nil
}

func (p *Parser) expectErr(start int, err error) error {
	if _, _, e := p.expect(start, TokenIdent); e != err {
		return e
	}
	return nil
}

func (p *Parser) decl(start int, kind TokenKind) (Token, int, error) {
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

type parsehandler func(int) (Token, int, error)

type parsematcher func(parsehandler) parsehandler

func (p *Parser) optional(start int, kind TokenKind, matchers ...parsematcher) (Token, int, bool) {
	t, next, err := p.expect(start, kind, matchers...)
	if err != nil {
		return t, start, false
	}

	if t.kind != kind {
		return t, start, false
	}

	return t, next, true
}

func (p *Parser) expect(start int, kind TokenKind, matchers ...parsematcher) (Token, int, error) {
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

func (p *Parser) skipBefore(kind TokenKind, more ...TokenKind) parsematcher {
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

type parseFunc func(int) (int, error)

func (p *Parser) repeat(start int, fn parseFunc) int {
	next := start
	for {
		if n, err := fn(next); err == nil {
			next = n
		} else {
			break
		}
	}
	return next
}
