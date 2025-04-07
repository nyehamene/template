package parser

import (
	"fmt"
	"strings"
	"temlang/tem/ast"
	"temlang/tem/dsa/stack"
	"temlang/tem/token"
)

func writePosition(w ast.SExprPrinterContext, loc Position, close ...string) {
	for _, c := range close {
		w.WriteString(c)
	}
	w.WriteString(" ; %s\n", w.Location(loc.Start, loc.End))
}

func writePositionOfToken(w ast.SExprPrinterContext, tok token.Token, close ...string) {
	l := Position{Start: tok.Start(), End: tok.End()}
	writePosition(w, l, close...)
}

func writePositionOfStack[T any](
	w ast.SExprPrinterContext,
	s stack.Stack[T],
	startOffset func(T) int,
	endOffset func(T) int,
	close ...string) {

	var first, last *T
	for {
		tok, ok := s.Pop()
		if !ok {
			break
		}
		if last == nil {
			last = &tok
			continue
		}
		first = &tok
	}
	if last == nil {
		w.WriteString("[]\n")
		return
	}
	if first == nil {
		first = last
	}
	start := startOffset(*first)
	end := endOffset(*last)
	l := Position{start, end}
	writePosition(w, l, close...)
}

func writePositionOfTokenStack(w ast.SExprPrinterContext, ts token.TokenStack, close ...string) {
	startFunc := func(t token.Token) int {
		return t.Start()
	}
	endFunc := func(t token.Token) int {
		return t.End()
	}
	writePositionOfStack(w, ts, startFunc, endFunc, close...)
}

func writeLocationOfTreeStack(w ast.SExprPrinterContext, ts TreeStack, close ...string) {
	startFunc := func(t Tree) int {
		return t.Pos().Start
	}
	endFunc := func(t Tree) int {
		return t.Pos().End
	}
	writePositionOfStack(w, ts, startFunc, endFunc, close...)
}

func writeTokenStack(
	w ast.SExprPrinterContext,
	t token.TokenStack,
	close string,
	closeOthers ...string) {
	i := 0
	for {
		tok, ok := t.Pop()
		i += 1
		if !ok {
			break
		}
		if i == t.Len() {
			closeOthers = append(closeOthers, close)
			writeLiteral(w, tok, closeOthers...)
			writePositionOfTokenStack(w, t)
			break
		}
		writeLiteral(w, tok, closeOthers...)
	}
}

func writeTreeStack(
	w ast.SExprPrinterContext,
	t TreeStack,
	close string,
	closeOthers ...string) {
	i := 0
	for {
		attr, ok := t.Pop()
		if !ok {
			break
		}
		if i == t.Len() {
			closeOthers = append(closeOthers, close)
			treeSExpr(w, attr, closeOthers...)
			break
		}
		treeSExpr(w, attr, closeOthers...)
	}
}

func writeDecl(w ast.SExprPrinterContext, d decltree, close ...string) {
	w.WriteString("%s(identifiers", w.Indentation())
	writePositionOfTokenStack(w, d.idents)

	w.Indent()

	if d.dtype.Kind() == token.Type {
		writeTokenStack(w, d.idents, ")")
		w.Dedent()
		w.WriteString("%s(type)", w.Indentation())
		writePositionOfToken(w, d.dtype, close...)
		return
	}

	writeTokenStack(w, d.idents, ")")
	w.Dedent()

	w.WriteString("%s(type", w.Indentation())
	writePositionOfToken(w, d.dtype)

	w.Indent()
	close = append(close, ")")
	writeLiteral(w, d.dtype, close...)
	w.Dedent()
}

func writeLiteral(w ast.SExprPrinterContext, lit token.Token, close ...string) {
	switch lit.Kind() {
	case token.String:
		w.WriteString("%s(string)", w.Indentation())
	case token.Directive:
		w.WriteString("%s(directive)", w.Indentation())
	case token.Ident:
		w.WriteString("%s(identifier)", w.Indentation())
	case token.Comment:
		w.WriteString("%s(comment)", w.Indentation())
	case token.TextBlock:
		w.WriteString("%s(text_block)", w.Indentation())
	case token.Package:
		w.WriteString("%s(package)", w.Indentation())
	case token.Import:
		w.WriteString("%s(import)", w.Indentation())
	case token.Using:
		w.WriteString("%s(using)", w.Indentation())
	case token.Type:
		w.WriteString("%s(type)", w.Indentation())
	case token.Templ:
		w.WriteString("%s(templ)", w.Indentation())
	default:
		w.WriteString("%s(ERROR", w.Indentation())
		writePositionOfToken(w, lit)
		w.Indent()
		w.WriteString("%s(%s))", w.Indentation(), lit.String())
	}
	for _, c := range close {
		w.WriteString(c)
	}
	w.WriteString(fmt.Sprintf(" '%s' ", lit.Text()))
	writePositionOfToken(w, lit)
}

func exprSExpr(w ast.SExprPrinterContext, e Expr, close ...string) {
	w.WriteString("%s(expr", w.Indentation())
	writePosition(w, e.Pos())
	w.Indent()
	defer w.Dedent()

	switch t := e.(type) {
	case pkgexpr:
		w.WriteString("%s(pkg_expr", w.Indentation())
		writePositionOfToken(w, t.name)
		w.Indent()
		w.WriteString("%s(name", w.Indentation())
		writePositionOfToken(w, t.name)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.name, close...)
		w.Dedent()
		w.Dedent()
	case importexpr:
		w.WriteString("%s(import_expr", w.Indentation())
		writePositionOfToken(w, t.path)
		w.Indent()
		w.WriteString("%s(path", w.Indentation())
		writePositionOfToken(w, t.path)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.path, close...)
		w.Dedent()
		w.Dedent()
	case usingexpr:
		w.WriteString("%s(using_expr", w.Indentation())
		writePositionOfToken(w, t.target)
		w.Indent()
		w.WriteString("%s(target", w.Indentation())
		writePositionOfToken(w, t.target)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.target, close...)
		w.Dedent()
		w.Dedent()
	case typeexpr:
		w.WriteString("%s(type_expr", w.Indentation())
		writePositionOfToken(w, t.target)
		w.Indent()
		w.WriteString("%s(target", w.Indentation())
		writePositionOfToken(w, t.target)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.target, close...)
		w.Dedent()
		w.Dedent()
	case recordexpr:
		w.WriteString("%s(record_expr", w.Indentation())
		writeLocationOfTreeStack(w, t.fields)
		w.Indent()
		w.WriteString("%s(fields", w.Indentation())
		writeLocationOfTreeStack(w, t.fields)
		w.Indent()
		writeTreeStack(w, t.fields, ")))", close...)
		w.Dedent()
		w.Dedent()
	case templexpr:
		w.WriteString("%s(templ_expr", w.Indentation())
		writeLocationOfTreeStack(w, t.params)
		w.Indent()
		w.WriteString("%s(params", w.Indentation())
		writeLocationOfTreeStack(w, t.params)
		w.Indent()
		writeTreeStack(w, t.params, ")))", close...)
		w.WriteString("%s(elements", w.Indentation())
		writeLocationOfTreeStack(w, t.elements)
		writeTreeStack(w, t.elements, ")))", close...)
		w.Dedent()
		w.Dedent()
	case litexpr:
		close = append(close, ")")
		writeLiteral(w, token.Token(t), close...)
	case badexpr:
		w.WriteString("%s(ERROR)", w.Indentation())
		w.WriteString(strings.Join(close, ""))
		writePosition(w, t.Position)
	default:
		w.WriteString(strings.Join(close, ""))
		w.WriteString("%s(ERROR)", w.Indentation())
	}
}

func treeSExpr(w ast.SExprPrinterContext, tree Tree, close ...string) {
	switch t := tree.(type) {
	case pkgtree:
		w.WriteString("%s(package_declaration", w.Indentation())
		writePosition(w, t.Position)

		w.Indent()
		writeDecl(w, t.decltree)

		w.WriteString("%s(directives", w.Indentation())
		writePositionOfTokenStack(w, t.directives)

		w.Indent()
		writeTokenStack(w, t.directives, ")")
		w.Dedent()

		close := append(close, ")")
		exprSExpr(w, t.expr, close...)

		w.Dedent()

	case importtree:
		w.WriteString("%s(import_declaration", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		writeDecl(w, t.decltree)

		close := append(close, ")")
		exprSExpr(w, t.expr, close...)

		w.Dedent()

	case usingtree:
		w.WriteString("%s(using_declaration", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		writeDecl(w, t.decltree)

		close := append(close, ")")
		exprSExpr(w, t.expr, close...)

		w.Dedent()

	case typetree:
		w.WriteString("%s(type_declaration", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		writeDecl(w, t.decltree)

		close := append(close, ")")
		exprSExpr(w, t.expr, close...)

		w.Dedent()

	case recordtree:
		w.WriteString("%s(record_declaration", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		writeDecl(w, t.decltree)

		close := append(close, ")")
		exprSExpr(w, t.expr, close...)

		w.Dedent()

	case templtree:
		w.WriteString("%s(templ_declaration", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		writeDecl(w, t.decltree)

		close := append(close, ")")
		exprSExpr(w, t.expr, close...)

		w.Dedent()

	case vartree:
		w.WriteString("%s(var_declaration", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		close = append(close, ")")
		writeDecl(w, decltree(t), close...)

		w.Dedent()

	case doctree:
		w.WriteString("%s(doc_declaration)", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		w.WriteString("%(identifiers", w.Indentation())
		writePositionOfTokenStack(w, t.idents)

		w.Indent()
		writeTokenStack(w, t.idents, ")")
		w.Dedent()

		w.WriteString("%(documentations", w.Indentation())
		w.Indent()
		writeTokenStack(w, t.text, ")", close...)
		w.Dedent()

		w.Dedent()

	case tagtree:
		w.WriteString("%s(tag_declaratin", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()

		w.WriteString("%s(identifiers", w.Indentation())
		writePositionOfTokenStack(w, t.idents)

		w.Indent()
		writeTokenStack(w, t.idents, ")")
		w.Dedent()

		w.WriteString("%s(attributes", w.Indentation())
		w.Indent()
		writeTreeStack(w, t.attrs, ")", close...)
		w.Dedent()

		w.Dedent()

	case attrtree:
		w.WriteString("%s(attr_declaratin", w.Indentation())
		writePosition(w, t.Position)
		w.Indent()
		close = append(close, ")")
		exprSExpr(w, t.value, close...)

		w.Dedent()

	case badtree:
		w.WriteString("%s(ERROR)", w.Indentation())
		writePosition(w, t.Position)

	default:
		w.WriteString("%s(ERROR)", w.Indentation())
	}
}

func (t pkgtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t importtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t usingtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t typetree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t vartree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t recordtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t templtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t doctree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t tagtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t attrtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}

func (t badtree) WriteSExpr(w ast.SExprPrinterContext) {
	treeSExpr(w, t)
}
