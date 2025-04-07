package parser

import (
	"os"
	"temlang/tem/ast"
	"temlang/tem/token"
)

func ParseFile(filename string, src []byte) (*ast.Namespace, *token.ErrorQueue) {
	if src == nil {
		var err error
		src, err = os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
	}

	name := "" // TODO get the namespace name from the filename
	file := ast.New(filename, name)
	p := New(filename, src)
	parse(file, &p)
	return file, p.errors
}

func parse(f *ast.Namespace, p *Parser) {
	var last token.Kind

	for p.cur.Kind() != token.EOF {
		tree := p.parseDoc(p.parseGenDecl)
		tree.TreeAst(f)

		if _, ok := tree.(badtree); ok {
			continue
		}

		switch p.lastTreeKind {
		case token.Package:
			if last != token.Invalid {
				p.errorf(p.loc(), "expected package declaration")
			}
		case token.Import:
			switch last {
			case token.Package, token.Import:
			default:
				p.errorf(p.loc(), "import must appear before other declarations")
			}
		case token.Using:
			switch last {
			case token.Package, token.Import, token.Using:
			default:
				p.errorf(p.loc(), "using must appear immediately after imports before other declarations")
			}
		}

		last = p.lastTreeKind
	}
}

func (p *Parser) parseDoc(f parseDeclSpec) Tree {
declStart:
	idents, ok := p.matchIdents()
	if !ok {
		p.errorExpected(p.loc(), "ident")
		return p.badtree()
	}

	switch k := p.cur.Kind(); k {
	case token.String, token.TextBlock:
		p.parseDocDecl(idents)
		goto declStart
	case token.BraceOpen:
		p.parseTagDecl(idents)
		goto declStart
	default:
		tree := f(idents)
		// FEAT: support trailing documentation and attributes
		// parseTrailingDoc(f, tree, idents)
		return tree
	}
}

func (p *Parser) parseGenDecl(idents token.TokenStack) Tree {
	var dtype token.Token
	var expr Expr
	var directives token.TokenStack

declStart:
	switch k := p.cur.Kind(); k {
	case token.Colon:
		p.advance()

		for p.cur.Kind() == token.Directive {
			directives.Push(p.cur)
			p.advance()
		}

		expr = p.parseGenExpr()

	case token.Package,
		token.Import,
		token.Using,
		token.Type,
		token.Templ:
		if dtype.Kind() != token.Invalid {
			p.errorf(p.loc(), "unexpected %s", k)
			return p.badtree()
		}
		p.advance()
		dtype = p.prev
		goto declStart

	case token.Ident:
		if dtype.Kind() != token.Invalid {
			p.errorf(p.loc(), "unexpected %s", k)
			return p.badtree()
		}
		p.advance()
		dtype = p.prev

	default:
		tree := p.badtree()
		p.advance() // advance to avoid infinit loop
		return tree
	}

	p.expectSemicolon()
	d := decltree{idents: idents, dtype: p.empty(token.Import)} // TODO add location

	// infer the tree based on the expr
	if dtype.Kind() == token.Invalid {
		return p.inferTreeFromExpr(expr, directives, d)
	}

	p.lastTreeKind = dtype.Kind()
	return p.createTree(dtype, directives, d, expr)
}

func (p *Parser) createTree(kind token.Token, directives token.TokenStack, d decltree, e Expr) Tree {
	switch kind.Kind() {
	case token.Package:
		return pkgtree{decltree: d, directives: directives, expr: e}
	case token.Import:
		return importtree{decltree: d, expr: e}
	case token.Using:
		return usingtree{decltree: d, expr: e}
	case token.Type:
		if _, ok := e.(typeexpr); ok {
			return recordtree{decltree: d, expr: e}
		}
		return typetree{decltree: d, expr: e}
	case token.Templ:
		return templtree{decltree: d, expr: e}
	case token.Ident:
		if e != nil {
			p.errorf(p.loc(), "expression not allowed in a var declaration")
		}
		return vartree(d)
	default:
		return badtree{expr: e, loc: p.loc()}
	}
}

func (p *Parser) inferTreeFromExpr(expr Expr, directives token.TokenStack, d decltree) Tree {
	switch expr.(type) {
	case pkgexpr:
		p.lastTreeKind = token.Package
		return pkgtree{d, directives, expr}
	case importexpr:
		p.lastTreeKind = token.Import
		return importtree{d, expr}
	case usingexpr:
		p.lastTreeKind = token.Using
		return usingtree{d, expr}
	case typeexpr:
		p.lastTreeKind = token.Type
		return typetree{d, expr}
	case recordexpr:
		p.lastTreeKind = token.Record
		return recordtree{d, expr}
	case templexpr:
		p.lastTreeKind = token.Templ
		return templtree{d, expr}
	case badexpr:
		return badtree{loc: p.loc(), expr: expr}
	default:
		return badtree{loc: p.loc(), expr: expr}
	}
}
