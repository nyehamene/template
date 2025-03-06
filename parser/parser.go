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

type Parser struct {
	tokenizer    tokenizer.Tokenizer
	file         *ast.NamespaceFile
	currentToken token.Token
	errorf       func(msg string, args ...any)
}

func (p *Parser) Mark() func() {
	return p.tokenizer.Mark()
}

func (p *Parser) addToken(tok token.Token) ast.TokenIndex {
	var index ast.TokenIndex
	if txt, ok := p.tokenizer.Text(tok); ok {
		index = p.file.AddToken(tok, txt)
	} else {
		p.errorf("invalid token: cannot get text at %s", tok)
	}
	return index
}

func (p *Parser) advance() bool {
	if p.currentToken.Kind == token.EOF {
		return false
	}
	tok := p.tokenizer.Next()
	p.currentToken = tok
	return true
}

func (p *Parser) match(tok token.Kind) bool {
	if p.currentToken.Kind != tok {
		return false
	}

	p.advance()
	return true
}

func (p *Parser) consume(tok token.Kind) matchresult.Type {
	matchedToken := p.currentToken
	if !p.match(tok) {
		var empty token.Token
		return matchresult.NoMatch(empty, tok)
	}
	return matchresult.Ok(matchedToken)
}

func (p *Parser) consumePackageName() (res matchresult.Type) {
	var name token.Token

	reset := p.Mark()
	defer func() {
		if !res.Ok() {
			reset()
		}
	}()

	if ok := p.match(token.Package); !ok {
		return matchresult.NoMatch(name, token.Package)
	}
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
	return matchresult.Ok(name)
}

func (p *Parser) consumePackageTempl() (res matchresult.Type) {
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

	switch p.currentToken.Kind {
	case token.Tag, token.Html, token.List:
		templ = p.currentToken
		p.advance()
	default:
		return
	}

	if ok := p.match(token.ParenClose); !ok {
		return matchresult.NoMatch(p.currentToken, token.ParenClose)
	}

	return matchresult.Ok(templ)
}

func (p *Parser) ParsePackage() (ast.PackageDecl, bool) {
	var ident token.Token
	var type0 token.Token
	var name token.Token
	var templ token.Token

	var decl ast.PackageDecl
	var res matchresult.Type

	if res = p.consume(token.Ident); !res.Ok() {
		return decl, false
	}
	if !p.match(token.Colon) {
		return decl, false
	}

	errFunc := func(r matchresult.Type) {
		p.errorf("invalid package declaraton; expected %s at %s",
			r.Exp(),
			r.Get())
	}

	ident = res.Get()
	isPackageDecl := false

switchStart:
	switch kind := p.currentToken.Kind; kind {
	case token.Colon:
		p.advance()
		res := p.consumePackageName()
		if res.NoMatch() {
			if isPackageDecl {
				errFunc(res)
			}
			return decl, false
		}
		if res.Invalid() {
			errFunc(res)
			return decl, false
		}
		// TBD: assign to type0
		name = res.Get()
		res = p.consumePackageTempl()
		if !res.Ok() {
			errFunc(res)
			return decl, false
		}
		templ = res.Get()
	case token.Package:
		type0 = p.currentToken
		isPackageDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if res := p.consume(token.Semicolon); !res.Ok() {
		p.errorf("invalid package declaraton; expected ; at %s", res.Get())
	}

	identIndex := p.addToken(ident)
	typeIndex := p.addToken(type0)
	nameIndex := p.addToken(name)
	templIndex := p.addToken(templ)

	decl = ast.PackageDecl{
		Decl: ast.Decl{
			Ident: identIndex,
			Type:  typeIndex,
		},
		Name:  nameIndex,
		Templ: templIndex,
	}
	return decl, true
}

func (p *Parser) ParseImport() (ast.ImportDecl, bool) {
	var ident token.Token
	var type0 token.Token
	var path token.Token

	var decl ast.ImportDecl
	var res matchresult.Type

	if res = p.consume(token.Ident); !res.Ok() {
		return decl, false
	}

	if !p.match(token.Colon) {
		return decl, false
	}

	errFunc := func(r matchresult.Type) {
		p.errorf("invalid import declaraton; expected %s at %s",
			r.Exp(),
			r.Get())
	}

	ident = res.Get()
	isImportDecl := false

switchStart:
	switch kind := p.currentToken.Kind; kind {
	case token.Colon:
		p.advance()
		res = p.consume(token.Import)
		if !res.Ok() {
			if isImportDecl {
				errFunc(res)
			}
			return decl, false
		}

		// TBD: assign type0

		if res = p.consume(token.ParenOpen); !res.Ok() {
			errFunc(res)
			return decl, false
		}

		path = p.currentToken

		if res = p.consume(token.String); !res.Ok() {
			errFunc(res)
			return decl, false
		}
		if res = p.consume(token.ParenClose); !res.Ok() {
			errFunc(res)
			return decl, false
		}
	case token.Import:
		type0 = p.currentToken
		isImportDecl = true
		p.advance()
		goto switchStart
	default:
		return decl, false
	}

	if res = p.consume(token.Semicolon); !res.Ok() {
		errFunc(res)
	}

	identIndex := p.addToken(ident)
	typeIndex := p.addToken(type0)
	pathIndex := p.addToken(path)

	decl = ast.ImportDecl{
		Decl: ast.Decl{
			Ident: identIndex,
			Type:  typeIndex,
		},
		Path: pathIndex,
	}

	return decl, true
}
