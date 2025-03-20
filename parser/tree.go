package parser

import (
	"temlang/tem/ast"
	"temlang/tem/dsa/stack"
	"temlang/tem/token"
)

type TreeStack = stack.Stack[Tree]

type Tree interface {
	TreeAst(*ast.Namespace)
}

type badtree struct {
	from, to token.Pos
}

func (t badtree) TreeAst(_ *ast.Namespace) {}

type decltree struct {
	idents token.TokenStack
	dtype  token.Token
}

type pkgtree struct {
	decltree
	expr Expr
}

func (t pkgtree) TreeAst(_ *ast.Namespace) {}

type importtree struct {
	decltree
	expr Expr
}

func (t importtree) TreeAst(_ *ast.Namespace) {}

type usingtree struct {
	decltree
	expr Expr
}

func (t usingtree) TreeAst(_ *ast.Namespace) {}

type typetree struct {
	decltree
	expr Expr
}

func (t typetree) TreeAst(_ *ast.Namespace) {}

type recordtree struct {
	decltree
	expr Expr
}

func (t recordtree) TreeAst(_ *ast.Namespace) {}

type templtree struct {
	decltree
	expr Expr
}

func (t templtree) TreeAst(_ *ast.Namespace) {}

type vartree decltree

func (t vartree) TreeAst(_ *ast.Namespace) {}

type tagtree struct {
	idents token.TokenStack
	attrs  TreeStack
}

func (t tagtree) TreeAst(_ *ast.Namespace) {}

type attrtree struct {
	idents token.TokenStack
	value  litexpr
}

func (t attrtree) TreeAst(_ *ast.Namespace) {}

type Expr interface {
	ExprAst(*ast.Namespace)
}

type doctree struct {
	idents token.TokenStack
	text   token.TokenStack
}

func (t doctree) TreeAst(_ *ast.Namespace) {}

type badexpr struct{}

func (e badexpr) ExprAst(_ *ast.Namespace) {}

type pkgexpr struct {
	directives token.TokenStack
	name       token.Token
}

func (e pkgexpr) ExprAst(_ *ast.Namespace) {}

type importexpr struct {
	path token.Token
}

func (e importexpr) ExprAst(_ *ast.Namespace) {}

type usingexpr struct {
	target token.Token
}

func (e usingexpr) ExprAst(_ *ast.Namespace) {}

type typeexpr struct {
	target token.Token
}

func (e typeexpr) ExprAst(_ *ast.Namespace) {}

type recordexpr struct {
	fields TreeStack
}

func (e recordexpr) ExprAst(_ *ast.Namespace) {}

type templexpr struct {
	params   TreeStack
	elements TreeStack
}

func (e templexpr) ExprAst(_ *ast.Namespace) {}

type LitExpr interface {
	Expr
	LitValue(*ast.Namespace)
}

type litexpr token.Token

func (e litexpr) ExprAst(_ *ast.Namespace)  {}
func (e litexpr) LitValue(_ *ast.Namespace) {}
