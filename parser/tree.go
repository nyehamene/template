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
	loc token.Location
}

type decltree struct {
	idents token.TokenStack
	dtype  token.Token
	loc    token.Location
}

type pkgtree struct {
	decltree
	directives token.TokenStack
	expr       Expr
}

type importtree struct {
	decltree
	expr Expr
}

type usingtree struct {
	decltree
	expr Expr
}

type typetree struct {
	decltree
	expr Expr
}

type recordtree struct {
	decltree
	expr Expr
}

type templtree struct {
	decltree
	expr Expr
}

type vartree decltree

type tagtree struct {
	idents token.TokenStack
	attrs  TreeStack
	loc    token.Location
}

type attrtree struct {
	idents token.TokenStack
	value  litexpr
	loc    token.Location
}

// TODO creae separate stringDoctree and textblockDoctree
type doctree struct {
	idents token.TokenStack
	text   token.TokenStack
	loc    token.Location
}

type Expr interface {
	ExprAst(*ast.Namespace)
}

type badexpr struct {
	loc token.Location
}

type pkgexpr struct {
	name token.Token
}

type importexpr struct {
	path token.Token
}

type usingexpr struct {
	target token.Token
}

type typeexpr struct {
	target token.Token
}

type recordexpr struct {
	fields TreeStack
	loc    token.Location
}

type templexpr struct {
	params   TreeStack
	elements TreeStack
	loc      token.Location
}

type LitExpr interface {
	Expr
	LitValue(*ast.Namespace)
}

type litexpr token.Token

func (t badtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
}

func (t pkgtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
	t.expr.ExprAst(n)
}

func (t importtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
	t.expr.ExprAst(n)
}

func (t usingtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
	t.expr.ExprAst(n)
}

func (t typetree) TreeAst(n *ast.Namespace) {
	n.Add(t)
	t.expr.ExprAst(n)
}

func (t recordtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
	t.expr.ExprAst(n)
}

func (t templtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
	t.expr.ExprAst(n)
}

func (t vartree) TreeAst(n *ast.Namespace) {
	n.Add(t)
}

func (t tagtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
}

func (t doctree) TreeAst(n *ast.Namespace) {
	n.Add(t)
}

func (t attrtree) TreeAst(n *ast.Namespace) {
	n.Add(t)
}

func (e badexpr) ExprAst(*ast.Namespace) {}

func (e pkgexpr) ExprAst(n *ast.Namespace) {
	name := "" // TODO get the package name from pkgexpr
	n.SetPackageName(name)
}

func (e importexpr) ExprAst(*ast.Namespace) {}

func (e usingexpr) ExprAst(*ast.Namespace) {}

func (e typeexpr) ExprAst(*ast.Namespace) {}

func (e recordexpr) ExprAst(*ast.Namespace) {}

func (e templexpr) ExprAst(*ast.Namespace) {}

func (e litexpr) ExprAst(*ast.Namespace) {}

func (e litexpr) LitValue(*ast.Namespace) {}
