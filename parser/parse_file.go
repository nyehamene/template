package parser

import (
	"fmt"
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

	file.SetLines(p.Lines())
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
				p.error(p.offset(), "expected package declaration")
			}
		case token.Import:
			switch last {
			case token.Package, token.Import:
			default:
				p.error(p.offset(), "import must appear before other declarations")
			}
		case token.Using:
			switch last {
			case token.Package, token.Import, token.Using:
			default:
				p.error(p.offset(), "using must appear immediately after imports before other declarations")
			}
		}

		last = p.lastTreeKind
	}
}

func (p *Parser) parseDoc(f parseDeclSpec) Tree {
declStart:
	ok := p.matchIdents()
	if !ok {
		p.errorExpected("ident")
		return p.badtree(p.offset())
	}

	switch k := p.cur.Kind(); k {
	case token.String, token.TextBlock:
		p.parseDocDecl()
		goto declStart
	case token.BraceOpen:
		p.parseTagDecl()
		goto declStart
	default:
		tree := f()
		// FEAT: support trailing documentation and attributes
		// parseTrailingDoc(f, tree, idents)
		return tree
	}
}

func (p *Parser) parseGenDecl() Tree {
	var dtype token.Token
	var expr Expr
	var directives token.TokenStack

	offset := p.identOffset

declStart:
	switch k := p.cur.Kind(); k {
	case token.Colon:
		p.advance()

		for p.cur.Kind() == token.Directive {
			directives.Push(p.cur)
			p.advance()
		}

		expr = p.parseGenExpr(offset)

	case token.Package,
		token.Import,
		token.Using,
		token.Type,
		token.Templ:
		if dtype.Kind() != token.Invalid {
			p.error(p.offset(), fmt.Sprintf("unexpected %s", k))
			return p.badtree(offset)
		}
		p.advance()
		dtype = p.prev
		goto declStart

	case token.Ident:
		if dtype.Kind() != token.Invalid {
			p.error(p.offset(), fmt.Sprintf("unexpected %s", k))
			return p.badtree(offset)
		}
		p.advance()
		dtype = p.prev

	default:
		tree := p.badtree(offset)
		p.advance() // advance to avoid infinite loop
		return tree
	}

	d := p.decltree(dtype)
	p.expectSemicolon()

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
			p.error(p.offset(), "expression not allowed in a var declaration")
		}
		return vartree(d)
	default:
		return p.badtree(d.Start)
	}
}

func (p *Parser) emptyToken(kind token.Kind) token.Token {
	tok := token.New(kind, 0, 0)
	return tok
}

func (p *Parser) inferTreeFromExpr(expr Expr, directives token.TokenStack, d decltree) Tree {
	var tree Tree
	switch expr.(type) {
	case pkgexpr:
		d.dtype = p.emptyToken(token.Package)
		tree = pkgtree{d, directives, expr}
	case importexpr:
		d.dtype = p.emptyToken(token.Import)
		tree = importtree{d, expr}
	case usingexpr:
		d.dtype = p.emptyToken(token.Using)
		tree = usingtree{d, expr}
	case typeexpr:
		d.dtype = p.emptyToken(token.Type)
		tree = typetree{d, expr}
	case recordexpr:
		d.dtype = p.emptyToken(token.Type)
		tree = recordtree{d, expr}
	case templexpr:
		d.dtype = p.emptyToken(token.Templ)
		tree = templtree{d, expr}
	// case badexpr:
	default:
		loc := p.locationStartingAt(d.Start)
		return badtree{loc}
	}
	p.lastTreeKind = d.dtype.Kind()
	return tree
}

func (p *Parser) Lines() []int {
	return p.tokenizer.Lines()
}
