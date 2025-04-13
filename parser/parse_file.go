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
				p.error(p.identOffset(), "expected package declaration")
			}
		case token.Import:
			switch last {
			case token.Package, token.Import:
			default:
				p.error(p.identOffset(), "import must appear before other declarations")
			}
		case token.Using:
			switch last {
			case token.Package, token.Import, token.Using:
			default:
				p.error(p.identOffset(), "using must appear immediately after imports before other declarations")
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
		return p.badtree(p.identOffset())
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
	var directives token.TokenQueue

	// NOTE: assume p.idents is not null
	idents := *p.idents

declStart:
	switch k := p.cur.Kind(); k {
	case token.Colon:
		dtypeOffset := p.offset()
		p.advance()

		for p.cur.Kind() == token.Directive {
			directives.Push(p.cur)
			p.advance()
		}

		expr = p.parseGenExpr()

		// infer the tree based on the expr
		if dtype.Kind() == token.Invalid {
			kind := p.inferTreeFromExpr(expr)
			dtype = p.emptyToken(kind, dtypeOffset)
		}

	case token.Package,
		token.Import,
		token.Using,
		token.Type,
		token.Templ:
		if dtype.Kind() != token.Invalid {
			offset := p.identOffset()
			p.error(offset, fmt.Sprintf("unexpected %s", k))
			return p.badtree(offset)
		}
		p.advance()
		dtype = p.prev
		goto declStart

	case token.Ident:
		if dtype.Kind() != token.Invalid {
			offset := p.identOffset()
			p.error(offset, fmt.Sprintf("unexpected %s", k))
			return p.badtree(offset)
		}
		p.advance()
		dtype = p.prev

	default:
		offset := p.identOffset()
		tree := p.badtree(offset)
		p.advance() // advance to avoid infinite loop
		return tree
	}

	d := p.decltree(idents, dtype)
	p.expectSemicolon()
	p.lastTreeKind = dtype.Kind()

	return p.createTree(dtype, directives, d, expr)
}

func (p *Parser) createTree(kind token.Token, directives token.TokenQueue, d decltree, e Expr) Tree {
	switch kind.Kind() {
	case token.Package:
		return pkgtree{decltree: d, directives: directives, expr: e}
	case token.Import:
		return importtree{decltree: d, expr: e}
	case token.Using:
		return usingtree{decltree: d, expr: e}
	case token.Type:
		return typetree{decltree: d, expr: e}
	case token.Templ:
		return templtree{decltree: d, expr: e}
	case token.Ident:
		if e != nil {
			p.error(p.identOffset(), "expression not allowed in a var declaration")
		}
		return vartree(d)
	default:
		return p.badtree(d.Start)
	}
}

func (p *Parser) emptyToken(kind token.Kind, offset int) token.Token {
	tok := token.New(kind, offset, offset)
	return tok
}

func (p *Parser) inferTreeFromExpr(expr Expr) token.Kind {
	switch expr.(type) {
	case pkgexpr:
		return token.Package
	case importexpr:
		return token.Import
	case usingexpr:
		return token.Using
	case typeexpr:
		return token.Type
	case recordexpr:
		return token.Type
	case templexpr:
		return token.Templ
	default:
		return token.Invalid
	}
}

func (p *Parser) Lines() []int {
	return p.tokenizer.Lines()
}
