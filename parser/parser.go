package parser

import (
	"fmt"
	"temlang/tem/dsa/queue"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
)

func New(filename string, src []byte) Parser {
	tok := tokenizer.New(filename, src)
	p := Parser{
		filename:  filename,
		tokenizer: tok,
		cur:       token.Token{},
		errors:    &token.ErrorQueue{},
	}
	p.error = func(offset int, msg string) {
		defaultErrorHandler(p.errors, offset, msg)
	}
	p.advance()
	return p
}

func defaultErrorHandler(errors *token.ErrorQueue, offset int, msg string) {
	err := token.NewError(offset, msg)
	// log.Printf("%s", err)
	errors.Push(err)
}

func queueOffset(q token.TokenQueue) int {
	tok, ok := q.Peek()
	if !ok {
		return -1
	}
	return tok.Start()
}

type Parser struct {
	tokenizer    tokenizer.Tokenizer
	filename     string
	cur          token.Token
	prev         token.Token
	lastTreeKind token.Kind
	idents       *token.TokenQueue
	errors       *token.ErrorQueue
	error        func(offset int, msg string)
}

func (p *Parser) errorExpected(msg string) {
	offset := p.offset()
	expected := p.cur.Kind()
	str := fmt.Sprintf("expected '%s' got %s", msg, expected)
	p.error(offset, str)
}

func (p *Parser) expectSemicolon() {
	if k := p.cur.Kind(); k != token.BraceClose {
		p.expect(token.Semicolon)
	}
}

func (p *Parser) offset() int {
	return p.cur.Start()
}

func (p *Parser) identOffset() int {
	if p.idents == nil {
		return -1
	}
	return queueOffset(*p.idents)
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
		p.errorExpected(tok.String())
		return false
	}
	return true
}

func (p *Parser) matchIdents() bool {
	var idents token.TokenQueue

	p.idents = nil
	// offset := p.offset()

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
		return false
	}
	p.idents = &idents
	// p.identOffset = offset
	return true
}

type parseTokenSpec func() bool

func (p *Parser) expectSurround(open token.Kind, close token.Kind, f parseTokenSpec) (token.Token, bool) {
	var surrounded token.Token
	if !p.expect(open) {
		return surrounded, false
	}
	if !f() {
		return surrounded, false
	}
	surrounded = p.prev
	if !p.expect(close) {
		return surrounded, false
	}
	return surrounded, true
}

func (p *Parser) expectSurroundParen(tok token.Kind) (token.Token, bool) {
	return p.expectSurround(token.ParenOpen, token.ParenClose, p.matchspec(tok))
}

func (p *Parser) matchspec(tok token.Kind) parseTokenSpec {
	return func() bool {
		return p.match(tok)
	}
}

func (p *Parser) locationStartingAt(offset int) Position {
	l := Position{Start: offset, End: p.offset()}
	return l
}

func (p *Parser) decltree(idents token.TokenQueue, dtype token.Token) decltree {
	offset := queueOffset(idents)
	l := p.locationStartingAt(offset)
	return decltree{idents: idents, dtype: dtype, Position: l}
}

func (p *Parser) badtree(offset int) Tree {
	return badtree{Position{Start: offset, End: p.offset()}}
}

func (p *Parser) badexpr(offset int) Expr {
	b := p.baseexpr(offset, p.offset())
	return badexpr{b}
}

func (p *Parser) badtreeStack(offset int) TreeQueue {
	err := queue.New[Tree](1)
	err.Push(p.badtree(offset))
	return err
}

type parseDeclSpec func() Tree

func (p *Parser) parseDocDecl() Tree {
	var lines token.TokenQueue
	offset := p.identOffset()

	if str := p.cur; p.match(token.String) {
		lines.Push(str)
		p.expectSemicolon()
	} else {
		if p.cur.Kind() != token.TextBlock {
			p.errorExpected("documentation text")
			return p.badtree(offset)
		}
		for p.cur.Kind() == token.TextBlock {
			if !p.expect(token.TextBlock) {
				break
			}
			lines.Push(p.prev)
			p.expectSemicolon()
		}
	}

	// FIX: match optional explicit semicolon
	p.match(token.Semicolon)
	// NOTE assume p.idents is not nil at this point
	idents := *p.idents
	pos := Position{Start: offset, End: p.prev.End()}
	return doctree{idents: idents, text: lines, Position: pos}
}

// NOTE this method can be removed
func (p *Parser) parseIdents(f parseDeclSpec) Tree {
	offset := p.offset()
	ok := p.matchIdents()
	if !ok {
		p.errorExpected("ident")
		return p.badtree(offset)
	}
	return f()
}

func (p *Parser) parseVarDecl() Tree {
	if !p.match(token.Ident) && !p.match(token.Type) {
		p.errorExpected("var type")
		return p.badtree(p.identOffset())
	}
	// NOTE: assume p.idents is not nil
	idents := *p.idents
	dtype := p.prev
	d := p.decltree(idents, dtype)
	return vartree(d)
}

func (p *Parser) parseParamDecl() TreeQueue {
	offset := p.offset()
	if !p.expect(token.ParenOpen) {
		p.errorExpected("(")
		return p.badtreeStack(offset)
	}

	params := TreeQueue{}
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
func (p *Parser) parseAttrDecl() Tree {
	ok := p.matchIdents()
	if !ok {
		p.errorExpected("attribute key")
		return p.badtree(p.identOffset())
	}

	if !p.expect(token.Eq) {
		return p.badtree(p.identOffset())
	}

	if !p.expect(token.String) {
		p.errorExpected("attribute value")
		return p.badtree(p.identOffset())
	}

	val := litexpr(p.prev)
	// NOTE: assume p.idents is not nil at this point
	idents := *p.idents
	loc := p.locationStartingAt(p.identOffset())
	return attrtree{idents: idents, value: val, Position: loc}
}

func (p *Parser) parseTagDecl() Tree {
	if !p.expect(token.BraceOpen) {
		p.errorExpected("{")
		return p.badtree(p.identOffset())
	}

	var attrs TreeQueue
	for {
		attr := p.parseAttrDecl()
		attrs.Push(attr)
		p.expectSemicolon()
		if p.cur.Kind() == token.BraceClose {
			break
		}
	}

	if !p.expect(token.BraceClose) {
		p.errorExpected("{")
		return p.badtree(p.identOffset())
	}

	p.expectSemicolon()

	// NOTE assume p.idents is not nil at this point
	idents := *p.idents
	loc := p.locationStartingAt(p.identOffset())
	tree := tagtree{
		idents:   idents,
		attrs:    attrs,
		Position: loc,
	}
	return tree
}

func (p *Parser) parseElements() TreeQueue {
	// TODO: parse template elements,text, and expression
	var ts TreeQueue
	return ts
}
