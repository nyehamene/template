package parser

import (
	"fmt"
	"log"
	"temlang/tem/ast"
	"temlang/tem/matchresult"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
)

func New(file *ast.NamespaceFile) Parser {
	file.Init()
	bytes := []byte(file.Src())
	t := tokenizer.New(bytes)
	badTok := token.Token{}
	p := Parser{t, file, badTok, defaultErrorHandler}
	p.advance()
	return p
}

func defaultErrorHandler(msg string, args ...any) {
	log.Printf(msg, args...)
	log.Println()
}

func errorInvalidDecl[T any, E any](p *Parser, k string, res matchresult.Type[T, E]) {
	pos := p.pos()
	p.errorf("invalid %s declaration at %s expected %s %s",
		k, res.Get(), res.Exp(),
		pos.String())
}

func errorExpectedSemicolon(p *Parser, k string) {
	pos := p.pos()
	p.errorf("invalid %s declaration at %s expected %s %s",
		k, p.currentToken,
		token.Semicolon,
		pos.String())
}

type Pos struct {
	line, col int
	file      string
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d:%d in %s]", p.line, p.col, p.file)
}

type Parser struct {
	tokenizer    tokenizer.Tokenizer
	file         *ast.NamespaceFile
	currentToken token.Token
	errorf       func(msg string, args ...any)
}

func (p *Parser) pos() Pos {
	src := p.file.Src()
	if src == "" {
		return Pos{0, 0, p.file.Name}
	}

	var (
		line, col = 1, 0
		sc        = src[0:p.currentToken.End()]
	)

	for _, r := range sc {
		if r == '\n' {
			line += 1
			col = 0
			continue
		}
		col += 1
	}

	return Pos{line, col, p.file.Name}
}

func (p *Parser) Mark() func() {
	prev := p.currentToken
	reset := p.tokenizer.Mark()
	return func() {
		reset()
		p.currentToken = prev
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
	if p.currentToken.Kind() == token.EOF {
		return false
	}
	tok := p.skipNewlineAndComment()
	p.currentToken = tok
	return true
}

func (p *Parser) match(tok token.Kind) bool {
	if p.currentToken.Kind() != tok {
		return false
	}

	p.advance()
	return true
}

func (p *Parser) consume(tok token.Kind) MatchToken {
	matchedToken := p.currentToken
	if !p.match(tok) {
		var empty token.Token
		return matchresult.NoMatch(empty, tok)
	}
	return Ok(matchedToken)
}

func (p *Parser) consumePackageName() MatchToken {
	var name token.Token

	if ok := p.match(token.ParenOpen); !ok {
		return matchresult.Invalid(p.currentToken, token.ParenOpen)
	}

	name = p.currentToken
	if ok := p.match(token.String); !ok {
		return matchresult.Invalid(p.currentToken, token.String)
	}
	if ok := p.match(token.ParenClose); !ok {
		return matchresult.Invalid(p.currentToken, token.ParenClose)
	}

	return Ok(name)
}

func (p *Parser) consumeIdents() (res MatchManyToken) {
	var idents []token.Token

	for {
		prev := p.currentToken
		if res := p.consume(token.Ident); !res.Ok() {
			return NoMatchMany(p.currentToken, token.Ident)
		}

		idents = append(idents, prev)

		if res := p.consume(token.Comma); !res.Ok() {
			break
		}
	}

	return OkMany(idents)
}

func (p *Parser) consumeKwExpr(kwk token.Kind, exprk token.Kind) (MatchToken, token.Token) {
	var empty token.Token

	prevKw := p.currentToken
	if kind := kwk; !p.match(kind) {
		return NoMatch(p.currentToken, kind), empty
	}

	if kind := token.ParenOpen; !p.match(kind) {
		return Invalid(p.currentToken, kind), empty
	}

	prevExpr := p.currentToken
	if kind := exprk; !p.match(kind) {
		return Invalid(p.currentToken, kind), empty
	}

	if kind := token.ParenClose; !p.match(kind) {
		return Invalid(p.currentToken, kind), empty
	}

	return Ok(prevExpr), prevKw
}

func (p *Parser) consumeVarDecl() (idents []token.Token, ty token.Token, ok bool) {
	resMany := p.consumeIdents()
	if !resMany.Ok() {
		return
	}
	if !p.match(token.Colon) {
		return
	}

	idents = resMany.Get()
	prev := p.currentToken
	if !p.match(token.Ident) {
		return
	}

	ty = prev
	ok = true
	return
}

func (p *Parser) consumeTemplParamsDecl() (idents []token.Token, ty token.Token, ok bool) {
	resMany := p.consumeIdents()
	if !resMany.Ok() {
		return
	}
	if !p.match(token.Colon) {
		return
	}

	idents = resMany.Get()
	prev := p.currentToken
	if !p.match(token.Ident) && !p.match(token.Type) {
		return
	}

	ty = prev
	ok = true
	return
}

func (p *Parser) consumeAttrDecl() (keys []token.Token, val token.Token, ok bool) {
	resMany := p.consumeIdents()
	if !resMany.Ok() {
		return
	}

	keys = resMany.Get()
	if !p.match(token.Eq) {
		return
	}

	prev := p.currentToken
	if !p.match(token.String) {
		return
	}

	val = prev
	ok = true
	return
}

func (p *Parser) parsePackageDecl() (ast.PackageDecl, bool) {
	var (
		declKind        = token.Package.String()
		ty, name, templ token.Token
		idents          []token.Token
		decl            ast.PackageDecl
	)

	resMany := p.consumeIdents()
	if !resMany.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	idents = resMany.Get()
	isPackageDecl := false

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()

	rhsSwitch:
		switch kind := p.currentToken.Kind(); kind {
		case token.Directive:
			// TODO: accept more than one directive
			templ = p.currentToken
			p.advance()
			goto rhsSwitch

		case token.Package:
			p.advance()
			res := p.consumePackageName()
			if res.NoMatch() {
				if isPackageDecl {
					errorInvalidDecl(p, declKind, res)
				}
				return decl, false
			}
			if res.Invalid() {
				errorInvalidDecl(p, declKind, res)
				return decl, false
			}
			if ty.Kind() == token.Invalid {
				ty = token.New(token.Package, 0, 0)
			}
			name = res.Get()
		}
	case token.Package:
		ty = p.currentToken
		isPackageDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p, declKind)
	}

	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetName(p.file, name)
	decl.SetTempl(p.file, templ)

	p.file.Pkg = decl

	return decl, true
}

func (p *Parser) parseImportDecl() (ast.ImportDecl, bool) {
	var (
		declKind = token.Import.String()
		decl     ast.ImportDecl
		ty, path token.Token
		idents   []token.Token
	)

	resMany := p.consumeIdents()
	if !resMany.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	idents = resMany.Get()
	isImportDecl := false

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res, kw := p.consumeKwExpr(token.Import, token.String)
		if res.NoMatch() {
			if isImportDecl {
				errorInvalidDecl(p, declKind, res)
			}
			return decl, false
		}
		if res.Invalid() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}
		if ty.Kind() == token.Invalid {
			ty = kw
		}

		path = res.Get()

	case token.Import:
		ty = p.currentToken
		isImportDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p, declKind)
	}

	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetName(p.file, path)

	p.file.AddImport(decl)

	return decl, true
}

func (p *Parser) parseUsingDecl() (ast.UsingDecl, bool) {
	var (
		declKind = token.Using.String()
		ty, pkg  token.Token
		idents   []token.Token
		decl     ast.UsingDecl
	)

	resMany := p.consumeIdents()
	if !resMany.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	idents = resMany.Get()
	isUsingDecl := false

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res, kw := p.consumeKwExpr(token.Using, token.Ident)
		if res.NoMatch() {
			if isUsingDecl {
				errorInvalidDecl(p, declKind, res)
			}
			return decl, false
		}
		if res.Invalid() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}
		if ty.Kind() == token.Invalid {
			ty = kw
		}
		pkg = res.Get()

	case token.Using:
		ty = p.currentToken
		isUsingDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p, declKind)
	}

	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetPkg(p.file, pkg)
	decl.SetIdents(p.file, idents)

	p.file.AddUsing(decl)

	return decl, true
}

func (p *Parser) parseAliasDecl() (ast.AliasDecl, bool) {
	var (
		declKind   = token.Alias.String()
		decl       ast.AliasDecl
		idents     []token.Token
		ty, target token.Token
		resMany    MatchManyToken
	)

	if resMany = p.consumeIdents(); !resMany.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	idents = resMany.Get()

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res, _ := p.consumeKwExpr(token.Type, token.Ident)
		if res.NoMatch() {
			return decl, false
		}
		if res.Invalid() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}
		if ty.Kind() == token.Invalid {
			ty = token.New(token.Alias, 0, 0)
		}
		target = res.Get()

	case token.Type:
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p, declKind)
	}

	decl.SetIdents(p.file, idents)
	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetTarget(p.file, target)

	return decl, true
}

func (p *Parser) parseRecordDecl() (ast.RecordDecl, bool) {
	var (
		declKind = token.Record.String()
		decl     ast.RecordDecl
		idents   []token.Token
		resMany  MatchManyToken
		ty       token.Token
		fields   = []ast.Entry[[]token.Token, token.Token]{}
	)

	if resMany = p.consumeIdents(); !resMany.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	idents = resMany.Get()

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res := p.consume(token.Record)
		if res.NoMatch() {
			return decl, false
		}
		if res.Invalid() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}
		if ty.Kind() == token.Invalid {
			ty = token.New(token.Record, 0, 0)
		}
		if res = p.consume(token.BraceOpen); !res.Ok() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}

		for {
			i, t, ok := p.consumeVarDecl()
			if !ok {
				break
			}

			fields = append(fields, ast.EntryMany(i, t))

			if !p.match(token.Semicolon) {
				errorExpectedSemicolon(p, declKind)
				break
			}
		}

		if res = p.consume(token.BraceClose); !res.Ok() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}

	case token.Type:
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p, declKind)
	}

	decl.SetIdents(p.file, idents)
	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetFields(p.file, fields)

	return decl, true
}

func (p *Parser) parseDocDecl() (ast.DocDecl, bool) {
	var (
		decl          ast.DocDecl
		idents, lines []token.Token
		resMany       MatchManyToken
	)

	if resMany = p.consumeIdents(); !resMany.Ok() {
		return decl, false
	}

	idents = resMany.Get()
	if !p.match(token.Colon) {
		return decl, false
	}

	if str := p.currentToken; p.match(token.String) {
		lines = append(lines, str)
		if !p.match(token.Semicolon) {
			errorExpectedSemicolon(p, "doc")
		}
	} else {
		if p.currentToken.Kind() != token.TextBlock {
			return decl, false
		}
		for {
			res := p.consume(token.TextBlock)
			if !res.Ok() {
				break
			}

			lines = append(lines, res.Get())
			if !p.match(token.Semicolon) {
				errorExpectedSemicolon(p, "doc")
				break
			}
			// NOTE: eol is already skipped by the last call to advance
			// if !p.match(token.EOL) {
			// 	errorExpected(p, "doc", token.EOL)
			// 	break
			// }
		}
	}

	decl.SetIdents(p.file, idents)
	decl.SetContent(p.file, lines...)

	return decl, true
}

func (p *Parser) parseTagDecl() (ast.TagDecl, bool) {
	var (
		decl    ast.TagDecl
		idents  []token.Token
		resMany MatchManyToken
	)

	attrs := []ast.Entry[[]token.Token, token.Token]{}

	if resMany = p.consumeIdents(); !resMany.Ok() {
		return decl, false
	}

	idents = resMany.Get()
	if !p.match(token.Colon) {
		return decl, false
	}
	if !p.match(token.BraceOpen) {
		return decl, false
	}

	for {
		keys, value, ok := p.consumeAttrDecl()
		if !ok {
			break
		}

		attrs = append(attrs, ast.EntryMany(keys, value))
		if !p.match(token.Semicolon) {
			break
		}
	}

	if !p.match(token.BraceClose) {
		return decl, false
	}
	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p, "tag")
		return decl, false
	}

	decl.SetIdents(p.file, idents)
	decl.SetAttrs(p.file, attrs)

	return decl, true
}

func (p *Parser) parseTemplDecl() (ast.TemplDecl, bool) {
	var (
		decl    ast.TemplDecl
		idents  []token.Token
		ty      token.Token
		resMany MatchManyToken
	)

	params := []ast.Entry[[]token.Token, token.Token]{}

	if resMany = p.consumeIdents(); !resMany.Ok() {
		return decl, false
	}

	idents = resMany.Get()
	if !p.match(token.Colon) {
		return decl, false
	}

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
	exprStart:
		switch kind := p.currentToken.Kind(); kind {
		case token.Templ:
			p.advance()
			goto exprStart
		case token.ParenOpen:
			p.advance()

			i, t, ok := p.consumeTemplParamsDecl()
			if !ok {
				return decl, false
			}

			params = append(params, ast.EntryMany(i, t))
			if !p.match(token.ParenClose) {
				return decl, false
			}
		default:
			return decl, false
		}
	case token.Templ:
		ty = p.currentToken
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if ty.Kind() == token.Invalid {
		ty = token.New(token.Templ, 0, 0)
	}
	if !p.match(token.BraceOpen) {
		return decl, false
	}

	if !p.match(token.BraceClose) {
		return decl, false
	}
	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p, "templ")
		return decl, false
	}

	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetParams(p.file, params)

	return decl, true
}
