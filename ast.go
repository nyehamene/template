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
	AstTypeDef
	AstRecordDef
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
	var ident Token
	var next int = start
	var err error
	// pkg: package: tag("home")
	if ident, next, err = p.expect(next, TokenIdent, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenColon, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, n, ok := p.optional(next, TokenPackage, p.skip(TokenSpace)); ok {
		next = n
	}
	if _, next, err = p.expect(next, TokenColon, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenTag, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenParLeft, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenString, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenParRight, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenSemicolon, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	return Ast{AstPackage, AstIdent, AstTag, ident.offset}, next, err
}

func (p Parser) TypeDef(start int) (Ast, int, error) {
	var ident Token
	var err error
	var next int = start
	// t : type : record {};
	// t :: record {};
	// t :: record { a: A; };
	if ident, next, err = p.expect(next, TokenIdent, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenColon, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, n, ok := p.optional(next, TokenType, p.skip(TokenSpace)); ok {
		next = n
	}
	if _, next, err = p.expect(next, TokenColon, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenRecord, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	if _, next, err = p.expect(next, TokenBraceLeft, p.skip(TokenSpace)); err != nil {
		return Ast{}, next, err
	}
	for {
		var n int
		var e error
		if _, n, e = p.expect(next, TokenIdent, p.skip(TokenSpace, TokenEOL)); e != nil {
			break
		}
		if _, next, err = p.expect(n, TokenColon, p.skip(TokenSpace)); err != nil {
			return Ast{}, next, err
		}
		if _, next, err = p.expect(next, TokenIdent, p.skip(TokenSpace)); err != nil {
			return Ast{}, next, err
		}
		if _, next, err = p.expect(next, TokenSemicolon, p.skip(TokenSpace)); err != nil {
			return Ast{}, next, err
		}
	}
	_, next, err = p.expect(next, TokenBraceRight, p.skip(TokenSpace, TokenEOL))
	return Ast{AstTypeDef, AstIdent, AstRecordDef, ident.offset}, next, err
}

func (p Parser) optional(start int, kind TokenKind, matchers ...parsematcher) (Token, int, bool) {
	var nothing Token

	t, next, err := p.expect(start, kind, matchers...)
	if err != nil {
		return nothing, start, false
	}

	if t.kind != kind {
		return nothing, start, false
	}

	return t, next, true
}

func (p Parser) expect(start int, kind TokenKind, matchers ...parsematcher) (Token, int, error) {
	var nothing Token
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

func (p Parser) skip(kind TokenKind, others ...TokenKind) parsematcher {
	return func(h parsehandler) parsehandler {
		return func(start int) (Token, int, error) {
			var nothing Token
			var afterToken = start
			for {
				token, next, err := p.tokenizer.next(afterToken)
				if err != nil {
					return nothing, start, err
				}
				if token.kind != kind {
					// and not equal to any
					for _, k := range others {
						if token.kind == k {
							// otherwise
							goto skipToken
						}
					}
					// then stop skipping
					next -= 1
					break
				}
			skipToken:
				afterToken = next
			}
			return h(afterToken)
		}
	}
}
