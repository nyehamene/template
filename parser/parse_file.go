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
	p.pkg(file)
	return file, p.errors
}

func (p *Parser) pkg(f *ast.Namespace) {
	tree := p.parseDoc(p.parsePackageDecl)
	if _, ok := tree.(pkgtree); !ok {
		var loc token.Location
		p.errorf(loc, "no package declaration")
	}

	tree.TreeAst(f)
	p.parseImport(f)
}

func (p *Parser) parseImport(f *ast.Namespace) {
	for p.cur.Kind() != token.EOF {
		tree := p.parseDoc(p.parseImportDecl)
		tree.TreeAst(f)
	}
}

func (p *Parser) parseDoc(f parseDeclSpec) Tree {
declStart:
	idents, ok := p.matchIdents()
	if p.cur.Kind() == token.EOF {
		p.errorf(p.loc(), "unexpected end of file")
		return p.badtree()
	}
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
		// parseTrailingDoc(f, tree, idents)
		return tree
	}
}

func (p *Parser) parseDecl(idents token.TokenStack) Tree {
	return p.parseBasicDecl(idents)
}

func (p *Parser) parseBasicDecl(idents token.TokenStack) Tree {
	var kind token.Kind
	var treeFunc func(decltree, token.TokenStack, Expr) Tree

	exprFunc := p.parseBasicExpr

	badTreeFunc := func(_ token.TokenStack) Tree {
		p.errorExpected(p.loc(), "a declaration")
		return p.badtree()
	}

	switch k := p.cur.Kind(); k {
	case token.Type:
		kind = k
		treeFunc = func(d decltree, _ token.TokenStack, e Expr) Tree {
			if _, ok := e.(recordexpr); ok {
				return recordtree{d, e}
			}
			return typetree{d, e}
		}

	case token.Record:
		kind = k
		exprFunc = p.parseRecordExpr
		treeFunc = func(d decltree, _ token.TokenStack, e Expr) Tree {
			return recordtree{d, e}
		}

	case token.Templ:
		kind = k
		exprFunc = p.parseTemplExpr
		treeFunc = func(d decltree, _ token.TokenStack, e Expr) Tree {
			return templtree{d, e}
		}

	default:
		return badTreeFunc(idents)
	}

	return p.parseGenDecl(idents, kind, exprFunc, treeFunc, badTreeFunc)
}

func (p *Parser) parseGenDecl(
	idents token.TokenStack,
	kind token.Kind,
	exprFunc parseExprSpec,
	treeFunc func(decltree, token.TokenStack, Expr) Tree,
	fallback parseDeclSpec,
) Tree {
	var dtype token.Token
	var expr Expr
	var directives token.TokenStack

	kindCount := 0
declStart:
	switch k := p.cur.Kind(); k {
	case token.Colon:
		p.advance()

		for p.cur.Kind() == token.Directive {
			directives.Push(p.cur)
			p.advance()
		}

		expr = exprFunc()

		// infer the tree based on the expr
		if dtype.Kind() == token.Invalid {
			return p.inferTreeFromExpr(expr, directives, idents)
		}

	case kind:
		if kindCount > 0 {
			p.errorf(p.loc(), "unexpected %s")
			return p.badtree()
		}
		p.advance()
		dtype = p.prev
		goto declStart

	default:
		return fallback(idents)
	}

	if dtype.Kind() == token.Invalid {
		dtype = p.empty(kind)
	}

	p.expectSemicolon()
	d := decltree{idents: idents, dtype: dtype} // TODO add location
	return treeFunc(d, directives, expr)
}

func (p *Parser) inferTreeFromExpr(expr Expr, directives token.TokenStack, idents token.TokenStack) Tree {
	d := decltree{idents: idents, dtype: p.empty(token.Import)} // TODO add location
	var tree Tree
	switch expr.(type) {
	case pkgexpr:
		tree = pkgtree{d, directives, expr}
	case importexpr:
		tree = importtree{d, expr}
	case usingexpr:
		tree = usingtree{d, expr}
	case typeexpr:
		tree = typetree{d, expr}
	case recordexpr:
		tree = recordtree{d, expr}
	case templexpr:
		tree = templtree{d, expr}
	case badexpr:
		tree = badtree{loc: p.loc(), expr: expr}
	default:
		tree = badtree{loc: p.loc(), expr: expr}
	}
	p.expectSemicolon()
	return tree
}
