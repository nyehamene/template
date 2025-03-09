package parser

import (
	"log"
	"temlang/tem/ast"
	"temlang/tem/matchresult"
	"temlang/tem/token"
	"temlang/tem/tokenizer"
)

// Adance on:
//           SUCCESS | FAILURE
// match   |  Yes    | No
// consume |  Yes    | No
// expect  |  Yes    | Yes

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

func errorInvalidDecl[T any, E any](p *Parser, k token.Kind, res matchresult.Type[T, E]) {
	p.errorf("invalid %s declaration at %s expected %s",
		k,
		res.Get(),
		res.Exp())
}

func errorExpectedSemicolon(p *Parser, k token.Kind) {
	p.errorf("invalid %s declaration at %s expected %s",
		k, p.currentToken, token.Semicolon)
}

type Parser struct {
	tokenizer    tokenizer.Tokenizer
	file         *ast.NamespaceFile
	currentToken token.Token
	errorf       func(msg string, args ...any)
}

func (p *Parser) Mark() func() {
	return p.tokenizer.Mark()
}

func (p *Parser) advance() bool {
	if p.currentToken.Kind() == token.EOF {
		return false
	}
	tok := p.tokenizer.Next()
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

func (p *Parser) consumePackageName() (res MatchToken, pkg token.Token) {
	var name token.Token

	reset := p.Mark()
	defer func() {
		if !res.Ok() {
			reset()
		}
	}()

	prev := p.currentToken
	if ok := p.match(token.Package); !ok {
		return matchresult.NoMatch(name, token.Package), pkg
	}

	pkg = prev

	if ok := p.match(token.ParenOpen); !ok {
		return matchresult.Invalid(p.currentToken, token.ParenOpen), pkg
	}

	prev = p.currentToken
	if ok := p.match(token.String); !ok {
		return matchresult.Invalid(p.currentToken, token.String), pkg
	}

	name = prev

	if ok := p.match(token.ParenClose); !ok {
		return matchresult.Invalid(p.currentToken, token.ParenClose), pkg
	}
	return Ok(name), pkg
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

func (p *Parser) consumePackageTempl() (tok token.Token, ok bool) {
	switch p.currentToken.Kind() {
	case token.Tag, token.Html, token.List:
		tok = p.currentToken
		ok = true
		p.advance()
	default:
		return
	}

	return
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

func (p *Parser) ParsePackage() (ast.PackageDecl, bool) {
	const declKind = token.Package

	var (
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
		res, pkg := p.consumePackageName()
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
			ty = pkg
		}

		name = res.Get()

		tok, ok := p.consumePackageTempl()
		if !ok {
			p.errorf("invalid %s declaration expected on of %s, %s, %s at %s",
				declKind,
				token.Tag,
				token.List,
				token.Html,
				p.currentToken,
			)
			return decl, false
		}
		templ = tok

	case declKind:
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

	return decl, true
}

func (p *Parser) ParseImport() (ast.ImportDecl, bool) {
	const declKind = token.Import

	var (
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
		res, kw := p.consumeKwExpr(declKind, token.String)
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

	case declKind:
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

	return decl, true
}

func (p *Parser) ParseUsing() (ast.UsingDecl, bool) {
	const declKind = token.Using

	var (
		ty, pkg token.Token
		idents  []token.Token
		decl    ast.UsingDecl
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
		res, kw := p.consumeKwExpr(declKind, token.Ident)
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

	case declKind:
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

	return decl, true
}

func (p *Parser) ParseAlias() (ast.AliasDecl, bool) {
	const declKind = token.Alias

	var (
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
	var isAliasDecl bool

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res, _ := p.consumeKwExpr(token.Type, token.Ident)
		if res.NoMatch() {
			if isAliasDecl {
				errorInvalidDecl(p, declKind, res)
			}
			return decl, false
		}
		if res.Invalid() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}
		if ty.Kind() == token.Invalid {
			ty = token.New(declKind, 0, 0)
		}
		target = res.Get()

	case token.Type:
		isAliasDecl = true
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

func (p *Parser) ParseRecord() (ast.RecordDecl, bool) {
	const declKind = token.Record

	var (
		decl    ast.RecordDecl
		idents  []token.Token
		resMany MatchManyToken
		ty      token.Token
		fields  = []ast.Entry[[]token.Token, token.Token]{}
	)

	if resMany = p.consumeIdents(); !resMany.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	idents = resMany.Get()
	var isRecordDecl bool

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res := p.consume(declKind)
		if res.NoMatch() {
			if isRecordDecl {
				errorInvalidDecl(p, declKind, res)
			}
			return decl, false
		}
		if res.Invalid() {
			errorInvalidDecl(p, declKind, res)
			return decl, false
		}
		if ty.Kind() == token.Invalid {
			ty = token.New(declKind, 0, 0)
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
		isRecordDecl = true
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

func (p *Parser) ParseDoc() (ast.DocDecl, bool) {
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
				break
			}
			if !p.match(token.EOL) {
				p.errorf("expected %s at %s", token.EOL, p.currentToken)
				break
			}
		}
	}

	decl.SetIdents(p.file, idents)
	decl.SetContent(p.file, lines...)

	return decl, true
}

func (p *Parser) ParseTag() (ast.TagDecl, bool) {
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
		p.errorf("expected %s at %s", token.Semicolon, p.currentToken)
		return decl, false
	}

	decl.SetIdents(p.file, idents)
	decl.SetAttrs(p.file, attrs)

	return decl, true
}

func (p *Parser) ParseTempl() (ast.TemplDecl, bool) {
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
		p.errorf("expected %s at %s", token.Semicolon, p.currentToken)
		return decl, false
	}

	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetParams(p.file, params)

	return decl, true
}
