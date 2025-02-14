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
	AstPackage AstKind = iota
	AstTag
	AstIdent
)

type Ast struct {
	kind   AstKind
	left   AstKind
	right  AstKind
	offset int
}

type Parser struct {
	tokenizer Tokenizer
}

func (p Parser) Package(start int) (Ast, int, error) {
	var err error
	var next int = start
	// pkg: package: tag("home")
	ident, next, err := p.expect(nil, next, TokenIdent, p.optionalSpace)
	_, next, err = p.expect(err, next, TokenColon, p.optionalSpace)
	_, next, _ = p.optional(err, next, TokenPackage, p.optionalSpace)
	_, next, err = p.expect(err, next, TokenColon, p.optionalSpace)
	_, next, err = p.expect(err, next, TokenTag, p.optionalSpace)
	_, next, err = p.expect(err, next, TokenParLeft, p.optionalSpace)
	_, next, err = p.expect(err, next, TokenString, p.optionalSpace)
	_, next, err = p.expect(err, next, TokenParRight, p.optionalSpace)
	_, next, err = p.expect(err, next, TokenSemicolon, p.optionalSpace)
	return Ast{AstPackage, AstIdent, AstTag, ident.offset}, next, err
}

func (p Parser) optional(err error, start int, kind TokenKind, matchers ...parsematcher) (Token, int, bool) {
	var nothing Token
	if err != nil {
		return nothing, start, false
	}

	t, next, err := p.expect(err, start, kind, matchers...)
	if err != nil {
		return nothing, start, false
	}

	if t.kind != kind {
		return nothing, start, false
	}

	return t, next, true
}

func (p Parser) expect(err error, start int, kind TokenKind, matchers ...parsematcher) (Token, int, error) {
	var nothing Token
	if err != nil {
		return nothing, start, err
	}
	var main parsehandler = func(s int) (Token, int, error) {
		var token Token
		var err error
		var next = s
		token, next, err = p.tokenizer.next(s)
		if err != nil {
			return nothing, next, err
		}
		if token.kind != kind {
			r, c := p.tokenizer.Pos(token)
			return nothing, next, fmt.Errorf("parse error: expected %s got %s [%d, %d]", kind, token.kind, r, c)
		}
		return token, next, nil
	}

	var handler parsehandler = main
	for _, matcher := range matchers {
		handler = matcher(handler)
	}

	return handler(start)
}

func (p Parser) optionalSpace(h parsehandler) parsehandler {
	return func(start int) (Token, int, error) {
		var nothing Token
		var afterSpace = start
		for {
			token, next, err := p.tokenizer.next(afterSpace)
			if err != nil {
				return nothing, start, err
			}
			if token.kind != TokenSpace {
				next -= 1
				break
			}
			afterSpace = next
		}
		return h(afterSpace)
	}
}
