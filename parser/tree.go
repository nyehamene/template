package parser

import (
	"temlang/tem/ast"
	"temlang/tem/dsa/queue"
	"temlang/tem/token"
)

type TreeQueue = queue.Queue[Tree]

type Tree interface {
	TreeAst(*ast.Namespace)
	Pos() Position
}

type badtree struct {
	Position
}

type decltree struct {
	idents token.TokenQueue
	dtype  token.Token
	Position
}

type pkgtree struct {
	decltree
	directives token.TokenQueue
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

type templtree struct {
	decltree
	expr Expr
}

type vartree decltree

type tagtree struct {
	idents token.TokenQueue
	attrs  TreeQueue
	Position
}

type attrtree struct {
	idents token.TokenQueue
	value  litexpr
	Position
}

// TODO creae separate stringDoctree and textblockDoctree
type doctree struct {
	idents token.TokenQueue
	text   token.TokenQueue
	Position
}

type Expr interface {
	ExprAst(*ast.Namespace)
	Pos() Position
}

type baseexpr struct {
	start, end int
}

type badexpr struct {
	baseexpr
}

type pkgexpr struct {
	baseexpr
	name token.Token
}

type importexpr struct {
	baseexpr
	path token.Token
}

type usingexpr struct {
	baseexpr
	target token.Token
}

type typeexpr struct {
	baseexpr
	target token.Token
}

type recordexpr struct {
	baseexpr
	fields TreeQueue
}

type templexpr struct {
	baseexpr
	params   TreeQueue
	elements TreeQueue
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
