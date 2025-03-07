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
	bytes := []byte(file.Src)
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
	p.errorf("invalid %s declaration at %s, expected %s",
		k,
		res.Get(),
		res.Exp())
}

func errorExpectedSemicolon(p *Parser) {
	p.errorf("invalid %s declaration at %s, expected %s",
		p.currentToken, token.Semicolon)
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

func (p *Parser) consumePackageTempl() (res MatchToken) {
	var templ token.Token

	reset := p.Mark()
	defer func() {
		if !res.Ok() {
			reset()
		}
	}()

	if ok := p.match(token.Templ); !ok {
		return matchresult.NoMatch(templ, token.Templ)
	}
	if ok := p.match(token.ParenOpen); !ok {
		return matchresult.NoMatch(p.currentToken, token.ParenOpen)
	}

	switch p.currentToken.Kind() {
	case token.Tag, token.Html, token.List:
		templ = p.currentToken
		p.advance()
	default:
		return
	}

	if ok := p.match(token.ParenClose); !ok {
		return matchresult.NoMatch(p.currentToken, token.ParenClose)
	}

	return Ok(templ)
}

func (p *Parser) ParsePackage() (ast.PackageDecl, bool) {
	var ident token.Token
	var ty token.Token
	var name token.Token
	var templ token.Token

	var decl ast.PackageDecl
	var res MatchToken

	if res = p.consume(token.Ident); !res.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	ident = res.Get()
	isPackageDecl := false

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res, pkg := p.consumePackageName()
		if res.NoMatch() {
			if isPackageDecl {
				errorInvalidDecl(p, token.Package, res)
			}
			return decl, false
		}
		if res.Invalid() {
			errorInvalidDecl(p, token.Package, res)
			return decl, false
		}

		if ty.Kind() == token.Invalid {
			ty = pkg
		}
		name = res.Get()

		if res = p.consumePackageTempl(); !res.Ok() {
			errorInvalidDecl(p, token.Package, res)
			return decl, false
		}
		templ = res.Get()

	case token.Package:
		ty = p.currentToken
		isPackageDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if res := p.consume(token.Semicolon); !res.Ok() {
		errorExpectedSemicolon(p)
	}

	decl.SetIdent(p.file, ident)
	decl.SetType(p.file, ty)
	decl.SetName(p.file, name)
	decl.SetTempl(p.file, templ)

	return decl, true
}

func (p *Parser) ParseImport() (ast.ImportDecl, bool) {
	var ident token.Token
	var ty token.Token
	var path token.Token

	var decl ast.ImportDecl
	var res MatchToken

	if res = p.consume(token.Ident); !res.Ok() {
		return decl, false
	}

	if !p.match(token.Colon) {
		return decl, false
	}

	ident = res.Get()
	isImportDecl := false

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res = p.consume(token.Import)
		if !res.Ok() {
			if isImportDecl {
				errorInvalidDecl(p, token.Import, res)
			}
			return decl, false
		}

		if ty.Kind() == token.Invalid {
			ty = res.Get()
		}

		if res = p.consume(token.ParenOpen); !res.Ok() {
			errorInvalidDecl(p, token.Import, res)
			return decl, false
		}
		if res = p.consume(token.String); !res.Ok() {
			errorInvalidDecl(p, token.Import, res)
			return decl, false
		}

		path = res.Get()

		if res = p.consume(token.ParenClose); !res.Ok() {
			errorInvalidDecl(p, token.Import, res)
			return decl, false
		}
	case token.Import:
		ty = p.currentToken
		isImportDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if res = p.consume(token.Semicolon); !res.Ok() {
		errorExpectedSemicolon(p)
	}

	decl.SetIdent(p.file, ident)
	decl.SetType(p.file, ty)
	decl.SetPath(p.file, path)

	return decl, true
}

func (p *Parser) ParseUsing() (ast.UsingDecl, bool) {
	var ident token.Token
	var idents []token.Token
	var ty token.Token
	var pkg token.Token

	var decl ast.UsingDecl
	var res MatchToken

	if res = p.consume(token.Ident); !res.Ok() {
		return decl, false
	}
	ident = res.Get()
	idents = append(idents, ident)

	for {
		if res = p.consume(token.Comma); !res.Ok() {
			break
		}
		if res = p.consume(token.Ident); !res.Ok() {
			p.errorf("invalid token: expected %s got %s", res.Exp(), res.Get())
			break
		}
		idents = append(idents, res.Get())
	}

	if !p.match(token.Colon) {
		return decl, false
	}

	isUsingDecl := false

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		res = p.consume(token.Using)
		if !res.Ok() {
			if isUsingDecl {
				errorInvalidDecl(p, token.Using, res)
			}
			return decl, false
		}

		if ty.Kind() == token.Invalid {
			ty = res.Get()
		}
		if res = p.consume(token.ParenOpen); !res.Ok() {
			errorInvalidDecl(p, token.Using, res)
			return decl, false
		}
		if res = p.consume(token.Ident); !res.Ok() {
			errorInvalidDecl(p, token.Using, res)
			return decl, false
		}

		pkg = res.Get()

		if res = p.consume(token.ParenClose); !res.Ok() {
			errorInvalidDecl(p, token.Using, res)
			return decl, false
		}
	case token.Using:
		ty = p.currentToken
		isUsingDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if res = p.consume(token.Semicolon); !res.Ok() {
		errorExpectedSemicolon(p)
	}

	decl.SetIdent(p.file, ident)
	decl.SetType(p.file, ty)
	decl.SetPkg(p.file, pkg)
	decl.SetIdents(p.file, idents)

	return decl, true
}

func (p *Parser) ParseAlias() (ast.AliasDecl, bool) {
	var decl ast.AliasDecl
	var ident, ty, target token.Token
	var idents []token.Token
	var resMany MatchManyToken

	if resMany, ident = p.consumeIdents(); !resMany.Ok() {
		return decl, false
	}

	idents = resMany.Get()

	if !p.match(token.Colon) {
		return decl, false
	}

	var isAliasDecl bool

switchStart:
	switch kind := p.currentToken.Kind(); kind {
	case token.Colon:
		p.advance()
		var res MatchToken
		if res = p.consume(token.Type); !res.Ok() {
			if isAliasDecl {
				errorInvalidDecl(p, token.Alias, res)
			}
			return decl, false
		}
		if ty.Kind() == token.Invalid {
			ty = res.Get()
		}
		if res = p.consume(token.ParenOpen); !res.Ok() {
			errorInvalidDecl(p, token.Alias, res)
			return decl, false
		}
		if res = p.consume(token.Ident); !res.Ok() {
			errorInvalidDecl(p, token.Alias, res)
			return decl, false
		}

		target = res.Get()
		if res = p.consume(token.ParenClose); !res.Ok() {
			errorInvalidDecl(p, token.Alias, res)
			return decl, false
		}
	case token.Type:
		ty = p.currentToken
		isAliasDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if !p.match(token.Semicolon) {
		errorExpectedSemicolon(p)
	}

	decl.SetIdent(p.file, ident)
	decl.SetIdents(p.file, idents)
	decl.SetType(p.file, ty)
	decl.SetTarget(p.file, target)

	return decl, true
}

func (p *Parser) consumeIdents() (MatchManyToken, token.Token) {
	var idents []token.Token
	var ident token.Token

	prev := p.currentToken
	if !p.match(token.Ident) {
		return NoMatchMany(p.currentToken, token.Ident), ident
	}

	ident = prev
	idents = append(idents, prev)

	for {
		var res MatchToken
		if res = p.consume(token.Comma); !res.Ok() {
			break
		}
		if res = p.consume(token.Ident); !res.Ok() {
			return InvalidMany(p.currentToken, token.Ident), ident
		}
		idents = append(idents, res.Get())
	}

	return OkMany(idents), ident
}
