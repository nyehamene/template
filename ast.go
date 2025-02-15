package template

import "fmt"

var allowEOLAfter = map[TokenKind]bool{
	TokenSemicolon: true,
	TokenParLeft:   true,
	TokenBraceLeft: true,
}

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
	AstTagPackage
	AstListPackage
	AstHtmlPackage
	AstIdent
	AstTypeIdent
	AstTypeDef
	AstRecordDef
	AstAliasDef
	AstTemplateDef
	AstTemplateBody
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
	var lastToken Token
	var right AstKind
	var err error

	ast := Ast{kind: AstPackage}
	next := start

	if lastToken, next, err = p.expect(next, TokenIdent); err != nil {
		return ast, next, err
	}
	ast.offset = lastToken.offset
	ast.left = AstIdent

	if lastToken, next, err = p.expect(next, TokenColon); err != nil {
		return ast, next, err
	}

	if t, n, ok := p.optional(next, TokenPackage); ok {
		next = n
		lastToken = t
	}

	if lastToken, next, err = p.expect(next, TokenColon); err != nil {
		return ast, next, err
	}

	if right, next, err = p.packageKind(next); err != nil {
		return ast, next, err
	}
	ast.right = right

	if lastToken, next, err = p.expect(next,
		TokenSemicolon,
		p.skipAfter(TokenEOL),
	); err != nil {
		return ast, next, err
	}

	return ast, next, nil
}

func (p Parser) packageKind(start int) (AstKind, int, error) {
	next := start
	kind := AstTagPackage
	if k, n, ok := p.packageKind0(next); ok {
		kind = k
		next = n
	} else {
		// Get the next token and display it in the error message
		nextToken, _, err := p.expect(next, TokenTag)
		return kind, next, fmt.Errorf("expected either tag, list, or html package type but got %v\n\t%v", nextToken, err)
	}

	var err error
	if _, next, err = p.expect(next, TokenParLeft); err != nil {
		return kind, next, err
	}

	if _, next, err = p.expect(next, TokenString); err != nil {
		return kind, next, err
	}

	if _, next, err = p.expect(next, TokenParRight); err != nil {
		return kind, next, err
	}

	return kind, next, nil
}

func (p Parser) packageKind0(start int) (AstKind, int, bool) {
	kindHandler := func(start int) (Token, int, error) {
		return p.tokenizer.next(start)
	}
	handler := p.skipBefore(TokenSpace)(kindHandler)
	token, n, err := handler(start)
	if err != nil {
		return AstPackage, n, false
	}
	switch token.kind {
	case TokenTag:
		return AstTagPackage, n, true
	case TokenList:
		return AstListPackage, n, true
	case TokenHtml:
		return AstHtmlPackage, n, true
	default:
		return AstPackage, n, false
	}
}

func (p Parser) Def(start int) (Ast, int, error) {
	var ast Ast
	next := start

	if def, n, err := p.typeDef(next); err == nil {
		ast = def
		next = n
	} else if def, n, err := p.templDef(next); err == nil {
		ast = def
		next = n
	} else {
		return def, next, err
	}

	if _, n, err := p.expect(next,
		TokenSemicolon,
		p.skipAfter(TokenEOL),
	); err != nil {
		return ast, next, err
	} else {
		next = n
	}

	return ast, next, nil
}

func (p Parser) templDef(start int) (Ast, int, error) {
	var err error

	ast := Ast{kind: AstTemplateDef}
	next := start

	token, next, err := p.templDecl(next)
	if err != nil {
		return ast, next, err
	}
	ast.offset = token.offset
	ast.left = AstIdent

	token, next, err = p.expect(next, TokenColon)
	if err != nil {
		return ast, next, err
	}

	token, next, err = p.templModel(next)
	if err != nil {
		return ast, next, err
	}

	token, next, err = p.expect(next,
		TokenBraceLeft,
		p.skipAfter(TokenEOL),
	)
	if err != nil {
		return ast, next, err
	}

	token, next, err = p.expect(next, TokenBraceRight)
	if err != nil {
		return ast, next, err
	}
	ast.right = AstTemplateBody

	return ast, next, nil
}

func (p Parser) templModel(start int) (Token, int, error) {
	var leftPar Token
	next := start
	token, n, err := p.expect(next,
		TokenParLeft,
		p.skipAfter(TokenEOL),
	)
	if err != nil {
		return token, n, err
	}
	next = n
	leftPar = token

	token, n, err = p.expect(next, TokenIdent)
	if err != nil {
		return token, n, err
	}
	next = n

	token, n, err = p.expect(next, TokenParRight)
	if err != nil {
		return token, n, err
	}
	next = n

	return leftPar, next, nil
}

func (p Parser) templDecl(start int) (Token, int, error) {
	return p.decl(start, TokenTempl)
}

func (p Parser) typeDef(start int) (Ast, int, error) {
	var lastToken Token
	var err error

	ast := Ast{kind: AstTypeDef}
	next := start

	if lastToken, next, err = p.typeDecl(next); err != nil {
		return ast, next, err
	}
	ast.offset = lastToken.offset
	ast.left = AstTypeIdent

	if lastToken, next, err = p.expect(next, TokenColon); err != nil {
		return ast, next, err
	}

	if kind, n, ok := p.recordDef(next); ok {
		next = n
		ast.right = kind
	} else if kind, n, ok := p.aliasDef(next); ok {
		next = n
		ast.right = kind
	} else {
		// TODO: get the offset of the next none space token
		return ast, next, fmt.Errorf("expected a record or alias definition")
	}
	return ast, next, nil
}

func (p Parser) aliasDef(start int) (AstKind, int, bool) {
	var err error
	next := start
	if _, next, err = p.expect(next, TokenAlias); err != nil {
		return AstAliasDef, next, false
	}
	if _, next, err = p.expect(next, TokenIdent); err != nil {
		return AstAliasDef, next, false
	}
	return AstAliasDef, next, true
}

func (p Parser) typeDecl(start int) (Token, int, error) {
	return p.decl(start, TokenType)
}

func (p Parser) recordDef(start int) (AstKind, int, bool) {
	var err error
	next := start
	if _, next, err = p.expect(next, TokenRecord); err != nil {
		return AstRecordDef, next, false
	}
	if _, next, err = p.expect(next,
		TokenBraceLeft,
		p.skipAfter(TokenEOL),
	); err != nil {
		return AstRecordDef, next, false
	}

	for {
		var n int
		var e error
		if _, n, e = p.expect(next, TokenIdent); e != nil {
			break
		}
		if _, next, err = p.expect(n, TokenColon); err != nil {
			return AstRecordDef, next, false
		}
		if _, next, err = p.expect(next, TokenIdent); err != nil {
			return AstRecordDef, next, false
		}
		if _, next, err = p.expect(next,
			TokenSemicolon,
			p.skipAfter(TokenEOL),
		); err != nil {
			return AstRecordDef, next, false
		}
	}
	if _, next, err = p.expect(next, TokenBraceRight); err != nil {
		return AstRecordDef, next, false
	}
	return AstRecordDef, next, true
}

func (p Parser) decl(start int, kind TokenKind) (Token, int, error) {
	var ident Token
	next := start

	token, next, err := p.expect(next, TokenIdent)
	if err != nil {
		return token, next, err
	}
	ident = token

	token, next, err = p.expect(next,
		TokenColon,
		p.skipAfter(TokenEOL),
	)
	if err != nil {
		return token, next, err
	}

	if _, n, ok := p.optional(next, kind); ok {
		next = n
	}
	return ident, next, nil
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
	var main parsehandler = func(start int) (Token, int, error) {
		token, next, err := p.tokenizer.next(start)
		if err != nil {
			return token, start, err
		}
		if token.kind != kind {
			r, c := p.tokenizer.Pos(token)
			return token, start, fmt.Errorf("parse error: expected %s got %s [%d, %d]", kind, token.kind, r, c)
		}
		return token, next, nil
	}

	var handler parsehandler = main
	for _, matcher := range matchers {
		handler = matcher(handler)
	}

	handler = p.skipBefore(TokenSpace)(handler)

	return handler(start)
}

func (p Parser) skipBefore(kind TokenKind) parsematcher {
	return func(h parsehandler) parsehandler {
		return func(start int) (Token, int, error) {
			var beforeToken = start
			for {
				toSkip, next, err := p.tokenizer.next(beforeToken)
				if err != nil {
					return toSkip, start, err
				}
				if toSkip.kind != kind {
					break
				}
				beforeToken = next
			}
			return h(beforeToken)
		}
	}
}

func (p Parser) skipAfter(kind TokenKind) parsematcher {
	return func(h parsehandler) parsehandler {
		return func(start int) (Token, int, error) {
			token, afterToken, err := h(start)
			if err != nil {
				return token, afterToken, err
			}
			for {
				toSkip, next, err := p.tokenizer.next(afterToken)
				if err == EOF {
					break
				}
				if err != nil {
					return toSkip, start, err
				}
				if toSkip.kind != kind {
					break
				}
				afterToken = next
			}
			return token, afterToken, nil
		}
	}
}
