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

type decltree struct {
	idents token.TokenStack
	dtype  token.Token
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
}

type attrtree struct {
	idents token.TokenStack
	value  litexpr
}

type Expr interface {
	ExprAst()
}

// TODO creae separate stringDoctree and textblockDoctree
type doctree struct {
	idents token.TokenStack
	text   token.TokenStack
}

type badexpr struct{}

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
}

type templexpr struct {
	params   TreeStack
	elements TreeStack
}

type LitExpr interface {
	Expr
	LitValue(*ast.Namespace)
}

type litexpr token.Token

func blowUpIfNotOk(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func toVarDecl(n *ast.Namespace, t TreeStack) stack.Stack[ast.VarDecl] {
	var fields stack.Stack[ast.VarDecl]
	for !t.Empty() {
		field, ok := t.Pop()
		blowUpIfNotOk(ok, "toVarDecl")
		vt, ok := field.(vartree)
		blowUpIfNotOk(ok, "recordtree")
		var vd ast.VarDecl
		vd.SetIdents(n, vt.idents)
		vd.SetType(n, vt.dtype)
		fields.Push(vd)

	}
	return fields
}

func (t badtree) TreeAst(*ast.Namespace) {}

func (t pkgtree) TreeAst(n *ast.Namespace) {
	var d ast.PackageDecl
	d.SetIdents(n, t.idents)
	d.SetType(n, t.dtype)
	d.SetDirective(n, t.directives)
	expr, ok := t.expr.(pkgexpr)
	blowUpIfNotOk(ok, "pkgtree")
	d.SetName(n, expr.name)
	n.SetPkg(d)
}

func (t importtree) TreeAst(n *ast.Namespace) {
	var d ast.ImportDecl
	d.SetIdents(n, t.idents)
	d.SetType(n, t.dtype)
	expr, ok := t.expr.(importexpr)
	blowUpIfNotOk(ok, "importtree")
	d.SetName(n, expr.path)
	n.AddImport(d)
}

func (t usingtree) TreeAst(n *ast.Namespace) {
	var d ast.UsingDecl
	d.SetIdents(n, t.idents)
	d.SetType(n, t.dtype)
	expr, ok := t.expr.(usingexpr)
	blowUpIfNotOk(ok, "importtree")
	d.SetPkg(n, expr.target)
	n.AddUsing(d)
}

func (t typetree) TreeAst(n *ast.Namespace) {
	var d ast.TypeDecl
	d.SetIdents(n, t.idents)
	d.SetType(n, t.dtype)
	expr, ok := t.expr.(typeexpr)
	blowUpIfNotOk(ok, "typetree")
	d.SetTarget(n, expr.target)
	n.AddType(d)
}

func (t recordtree) TreeAst(n *ast.Namespace) {
	var d ast.RecordDecl
	d.SetIdents(n, t.idents)
	d.SetType(n, t.dtype)
	expr, ok := t.expr.(recordexpr)
	blowUpIfNotOk(ok, "recordtree")
	fields := toVarDecl(n, expr.fields)
	d.SetFields(n, fields)
	n.AddRecord(d)
}

func (t templtree) TreeAst(n *ast.Namespace) {
	var d ast.TemplDecl
	d.SetIdents(n, t.idents)
	d.SetType(n, t.dtype)
	expr, ok := t.expr.(templexpr)
	blowUpIfNotOk(ok, "templtre")
	params := toVarDecl(n, expr.params)
	d.SetParams(n, params)
	// elements := toVarDecl(n, expr.elements)
	// d.SetElements(n, elements)
	n.AddTempl(d)
}

func (t vartree) TreeAst(n *ast.Namespace) {
	var d ast.VarDecl
	d.SetIdents(n, t.idents)
	d.SetType(n, t.dtype)
	n.AddVar(d)
}

func (t tagtree) TreeAst(n *ast.Namespace) {
	var d ast.TagDecl
	var attrs stack.Stack[ast.AttrDecl]
	d.SetIdents(n, t.idents)
	for !t.attrs.Empty() {
		attr, ok := t.attrs.Pop()
		blowUpIfNotOk(ok, "tagtree")
		at, ok := attr.(attrtree)
		blowUpIfNotOk(ok, "tagtree")
		var ad ast.AttrDecl
		ad.SetIdents(n, at.idents)
		ad.SetValue(n, token.Token(at.value))
		attrs.Push(ad)
	}
	d.SetAttrs(n, attrs)
}

func (t doctree) TreeAst(n *ast.Namespace) {
	var d ast.DocDecl
	d.SetIdents(n, t.idents)
	d.SetContent(n, t.text)
	n.AddDoc(d)
}

func (t attrtree) TreeAst(*ast.Namespace) {
}

func (e badexpr) ExprAst() {}

func (e pkgexpr) ExprAst() {}

func (e importexpr) ExprAst() {}

func (e usingexpr) ExprAst() {}

func (e typeexpr) ExprAst() {}

func (e recordexpr) ExprAst() {}

func (e templexpr) ExprAst() {}

func (e litexpr) ExprAst() {}

func (e litexpr) LitValue() {}
