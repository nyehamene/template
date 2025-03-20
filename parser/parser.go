package parser

import (
	"fmt"
	"log"
	"temlang/tem/dsa/stack"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
)

func New(filename string, src []byte) Parser {
	tok := tokenizer.New(filename, src)
	p := Parser{
		filename:  filename,
		src:       src,
		tokenizer: tok,
		cur:       token.Token{},
		errors:    &token.ErrorQueue{},
	}
	p.errorf = func(pos token.Location, msg string, argv ...any) {
		defaultErrorHandler(p.errors, pos, msg, argv...)
	}
	p.advance()
	return p
}

func defaultErrorHandler(errors *token.ErrorQueue, pos token.Location, msg string, args ...any) {
	err := token.Error{
		Msg:      fmt.Sprintf(msg, args...),
		Location: pos,
	}
	log.Printf("%s", err)
	errors.Push(err)
}

type Parser struct {
	tokenizer tokenizer.Tokenizer
	filename  string
	src       []byte // TODO remove src
	cur       token.Token
	prev      token.Token
	errors    *token.ErrorQueue
	errorf    func(p token.Location, msg string, args ...any)
}

func (p *Parser) errorExpected(pos token.Location, msg string) {
	p.errorf(pos, fmt.Sprintf("expected '%s' got %s", msg, p.cur.Kind()))
}

func (p *Parser) expectSemicolon() {
	if k := p.cur.Kind(); k != token.BraceClose {
		p.expect(token.Semicolon)
	}
}

func (p *Parser) Mark() func() {
	prev := p.cur
	reset := p.tokenizer.Mark()
	return func() {
		reset()
		p.cur = prev
	}
}

func (p *Parser) skipNewlineAndComment() token.Token {
	var next token.Token
	for {
		next = p.tokenizer.Next()
		kind := next.Kind()
		if kind != token.EOL && kind != token.Comment {
			break
		}
	}
	return next
}

func (p *Parser) advance() bool {
	if p.cur.Kind() == token.EOF {
		return false
	}
	tok := p.skipNewlineAndComment()
	p.prev = p.cur
	p.cur = tok
	return true
}

func (p *Parser) match(tok token.Kind) bool {
	if p.cur.Kind() != tok {
		return false
	}

	p.advance()
	return true
}

func (p *Parser) expect(tok token.Kind) bool {
	if !p.match(tok) {
		p.errorExpected(p.loc(), tok.String())
		return false
	}
	return true
}

func (p *Parser) matchIdents() (token.TokenStack, bool) {
	var idents token.TokenStack

identStart:
	switch k := p.cur.Kind(); k {
	case token.Ident:
		p.advance()
		idents.Push(p.prev)
		switch k = p.cur.Kind(); k {
		case token.Comma:
			p.advance()
			goto identStart
		case token.Colon:
			p.advance()
		}
	default:
		return idents, false
	}
	return idents, true
}

type parseTokenSpec func() bool

func (p *Parser) expectSurround(open token.Kind, close token.Kind, f parseTokenSpec) bool {
	if !p.expect(open) {
		return false
	}
	if !f() {
		return false
	}
	surrounded := p.prev
	if !p.expect(close) {
		return false
	}
	p.prev = surrounded
	return true
}

func (p *Parser) expectSurroundParen(tok token.Kind) bool {
	return p.expectSurround(token.ParenOpen, token.ParenClose, p.matchspec(tok))
}

func (p *Parser) matchspec(tok token.Kind) parseTokenSpec {
	return func() bool {
		return p.match(tok)
	}
}

type parseTreeSpec func() TreeStack

func (p *Parser) expectSurroundTree(open token.Kind, close token.Kind, f parseTreeSpec) TreeStack {
	if !p.expect(open) {
		p.errorExpected(p.loc(), open.String())
		return p.badtreeStack()
	}
	tree := f()
	if !p.expect(close) {
		p.errorExpected(p.loc(), close.String())
		return p.badtreeStack()
	}
	return tree
}

func (p *Parser) expectSurroundTreeBrace(f parseTreeSpec) TreeStack {
	return p.expectSurroundTree(token.BraceOpen, token.BraceClose, f)
}

func (p *Parser) empty(tok token.Kind) token.Token {
	// pos := p.pos()
	// FIX: dtype should have zero length because it not declared in the source code
	// return token.New(tok, pos.Line, pos.Col)
	return token.New(tok, 0, 0)
}

func (p *Parser) badtree() Tree {
	loc := p.loc()
	from := loc.Start
	to := loc.End
	return badtree{from, to}
}

func (p *Parser) badexpr() Expr {
	return badexpr{}
}

func (p *Parser) badtreeStack() TreeStack {
	err := stack.New[Tree](1)
	err.Push(p.badtree())
	return err
}

type parseDeclSpec func(token.TokenStack) Tree

func (p *Parser) parseDocDecl(idents token.TokenStack) Tree {
	var lines token.TokenStack

	if str := p.cur; p.match(token.String) {
		lines.Push(str)
		p.expectSemicolon()
	} else {
		if p.cur.Kind() != token.TextBlock {
			p.errorExpected(p.loc(), "documentation text")
			return p.badtree()
		}
		for p.cur.Kind() == token.TextBlock {
			if !p.expect(token.TextBlock) {
				break
			}
			lines.Push(p.prev)
			p.expectSemicolon()
		}
	}

	// fix: match optional explicit semicolon
	p.match(token.Semicolon)
	return doctree{idents, lines}
}

func (p *Parser) parsePackageDecl(idents token.TokenStack) Tree {
	return p.parseGenDecl(
		idents,
		token.Package,
		p.parsePackageExpr,
		func(d decltree, e Expr) Tree {
			return pkgtree{d, e}
		},
		p.parseImportDecl,
	)
}

func (p *Parser) parseImportDecl(idents token.TokenStack) Tree {
	return p.parseGenDecl(
		idents,
		token.Import,
		p.parseImportExpr,
		func(d decltree, e Expr) Tree {
			return importtree{d, e}
		},
		p.parseUsingDecl,
	)
}

func (p *Parser) parseUsingDecl(idents token.TokenStack) Tree {
	return p.parseGenDecl(
		idents,
		token.Using,
		p.parseUsingExpr,
		func(d decltree, e Expr) Tree {
			return usingtree{d, e}
		},
		p.parseDecl,
	)
}

func (p *Parser) parseVarDecl(idents token.TokenStack) Tree {
	if !p.match(token.Ident) && !p.match(token.Type) {
		p.errorExpected(p.loc(), "var type")
		return p.badtree()
	}
	dtype := p.prev
	return vartree{idents, dtype}
}

func (p *Parser) parseIdents(f parseDeclSpec) Tree {
	idents, ok := p.matchIdents()
	if !ok {
		p.errorExpected(p.loc(), "ident")
		return p.badtree()
	}
	return f(idents)
}

func (p *Parser) parseParamDecl() TreeStack {
	if !p.expect(token.ParenOpen) {
		p.errorExpected(p.loc(), "(")
		return p.badtreeStack()
	}

	params := TreeStack{}
	for {
		param := p.parseIdents(p.parseVarDecl)
		params.Push(param)
		if !p.match(token.Comma) {
			break
		}
	}

	p.expect(token.ParenClose)
	return params
}

// parseAttrDecl should return ast.TokenIndex for the first var added to namespace
func (p *Parser) parseAttrDecl(idents token.TokenStack) Tree {
	p.expect(token.Eq)
	if !p.expect(token.String) {
		p.errorExpected(p.loc(), "attribute value")
		return p.badtree()
	}

	val := litexpr(p.prev)
	return attrtree{idents, val}
}

func (p *Parser) parseTagDecl(idents token.TokenStack) Tree {
	if !p.expect(token.BraceOpen) {
		p.errorExpected(p.loc(), "{")
		return p.badtree()
	}

	var attrs TreeStack
	for {
		attr := p.parseIdents(p.parseAttrDecl)
		attrs.Push(attr)
		p.expectSemicolon()
		if p.cur.Kind() == token.BraceClose {
			break
		}
	}

	if !p.expect(token.BraceClose) {
		p.errorExpected(p.loc(), "{")
		return p.badtree()
	}

	p.expectSemicolon()

	tree := tagtree{
		idents: idents,
		attrs:  attrs,
	}
	return tree
}

func (p *Parser) parseElements() TreeStack {
	// TODO: parse template elements,text, and expression
	var ts TreeStack
	return ts
}

func (p *Parser) loc() token.Location {
	return p.tokenizer.Location(p.cur)
}
