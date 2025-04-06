package parser

import (
	"strings"
	"temlang/tem/ast"
	"temlang/tem/token"
)

func writeLocation(w ast.SExprPrinterContext, loc token.Location) {
	w.WriteString("[%-2d, %2d] - [%-2d, %2d]\n",
		loc.Start.Line, loc.Start.Col,
		loc.End.Line, loc.End.Col,
	)
}

func writeLocationToken(w ast.SExprPrinterContext, tok token.Token) {
	writeLocation(w, w.Location(tok))
}

func (d decltree) writeIdents(w ast.SExprPrinterContext, close ...string) {
	if d.idents.Len() < 1 {
		w.WriteString("%s(ERROR) ; ", w.Indentation())
		writeLocation(w, d.loc)
	} else {

		for {
			ident, ok := d.idents.Pop()
			if !ok {
				break
			}
			writeLiteral(w, ident, "")
		}
	}

	w.WriteString("%s(type ; ", w.Indentation())
	writeLocationToken(w, d.dtype)
	close = append(close, ")")
	w.Indent()
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
		w.WriteString("%s(ERROR)", w.Indentation())
	}
	for _, c := range close {
		w.WriteString(c)
	}
	w.WriteString(" ; ")
	writeLocationToken(w, lit)
}

func exprSExpr(w ast.SExprPrinterContext, e Expr, close ...string) {
	w.WriteString("%s(expr", w.Indentation())
	w.WriteString("\n")
	w.Indent()
	defer w.Dedent()

	switch t := e.(type) {
	case pkgexpr:
		w.WriteString("%s(pkg_expr ; ", w.Indentation())
		writeLocationToken(w, t.name)
		w.Indent()
		w.WriteString("%s(name ; ", w.Indentation())
		writeLocationToken(w, t.name)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.name, close...)
		w.Dedent()
		w.Dedent()
	case importexpr:
		w.WriteString("%s(import_expr ; ", w.Indentation())
		writeLocationToken(w, t.path)
		w.Indent()
		w.WriteString("%s(path ; ", w.Indentation())
		writeLocationToken(w, t.path)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.path, close...)
		w.Dedent()
		w.Dedent()
	case usingexpr:
		w.WriteString("%s(using_expr ; ", w.Indentation())
		writeLocationToken(w, t.target)
		w.Indent()
		w.WriteString("%s(target ; ", w.Indentation())
		writeLocationToken(w, t.target)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.target, close...)
		w.Dedent()
		w.Dedent()
	case typeexpr:
		w.WriteString("%s(type_expr ; ", w.Indentation())
		writeLocationToken(w, t.target)
		w.Indent()
		w.WriteString("%s(target ; ", w.Indentation())
		writeLocationToken(w, t.target)
		w.Indent()
		close = append(close, ")))")
		writeLiteral(w, t.target, close...)
		w.Dedent()
		w.Dedent()
	case recordexpr:
		w.WriteString("%s(record_expr ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		w.WriteString("%s(fields", w.Indentation())
		w.WriteString("\n")
		w.Indent()
		i := 0
		for {
			field, ok := t.fields.Pop()
			if !ok {
				break
			}
			if i == t.fields.Len() {
				close = append(close, ")))")
				treeSExpr(w, field, close...)
			} else {
				treeSExpr(w, field)
			}
		}
		w.Dedent()
		w.Dedent()
	case templexpr:
		w.WriteString("%s(templ_expr ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		w.WriteString("%s(params ; ", w.Indentation())
		w.Indent()
		i := 0
		for {
			param, ok := t.params.Pop()
			if !ok {
				break
			}
			if i == t.params.Len() {
				close = append(close, ")))")
				treeSExpr(w, param, close...)
			} else {
				treeSExpr(w, param)
			}
		}
		w.WriteString("%s(elements ; ", w.Indentation())
		i = 0
		for {
			param, ok := t.elements.Pop()
			if !ok {
				break
			}
			if i == t.params.Len() {
				close = append(close, ")))")
				treeSExpr(w, param, close...)
			} else {
				treeSExpr(w, param)
			}
		}
		w.Dedent()
		w.Dedent()
	case litexpr:
		close = append(close, ")")
		writeLiteral(w, token.Token(t), close...)
	case badexpr:
		w.WriteString("%s(ERROR)", w.Indentation())
		w.WriteString(strings.Join(close, ""))
		writeLocation(w, t.loc)
	default:
		w.WriteString(strings.Join(close, ""))
		w.WriteString("%s(ERROR) ; ", w.Indentation())
	}
}

func treeSExpr(w ast.SExprPrinterContext, tree Tree, close ...string) {
	defer w.Dedent()

	switch t := tree.(type) {
	case pkgtree:
		w.WriteString("%s(package_declaration ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		for {
			dir, ok := t.directives.Pop()
			if !ok {
				break
			}
			writeLiteral(w, dir)
		}
		close := append(close, ")")
		exprSExpr(w, t.expr, close...)
	case importtree:
		w.WriteString("%s(import_declaration ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		close := append(close, ")")
		exprSExpr(w, t.expr, close...)
	case usingtree:
		w.WriteString("%s(using_declaration ; ", w.Indentation())
		writeLocation(w, t.loc)
		close := append(close, ")")
		w.Indent()
		exprSExpr(w, t.expr, close...)
	case typetree:
		w.WriteString("%s(type_declaration ; ", w.Indentation())
		writeLocation(w, t.loc)
		close := append(close, ")")
		w.Indent()
		exprSExpr(w, t.expr, close...)
	case recordtree:
		w.WriteString("%s(record_declaration ; ", w.Indentation())
		writeLocation(w, t.loc)
		close := append(close, ")")
		w.Indent()
		exprSExpr(w, t.expr, close...)
	case templtree:
		w.WriteString("%s(templ_declaration ; ", w.Indentation())
		writeLocation(w, t.loc)
		close := append(close, ")")
		w.Indent()
		exprSExpr(w, t.expr, close...)
	case vartree:
		w.WriteString("%s(var_declaration ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		close = append(close, ")")
		decltree(t).writeIdents(w, close...)
	case doctree:
		w.WriteString("%s(doc_declaration) ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		if t.text.Empty() {
			w.WriteString("%s(ERROR) ; ", w.Indentation())
			writeLocation(w, t.loc)
			break
		}
		i := 0
		for {
			text, ok := t.text.Pop()
			i += 1
			if !ok {
				break
			}
			if i == t.text.Len() {
				writeLiteral(w, text, ")")
				break
			}
			writeLiteral(w, text)
		}
	case tagtree:
		w.WriteString("%s(tag_declaratin ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		i := 0
		for {
			attr, ok := t.attrs.Pop()
			if !ok {
				break
			}
			if i == t.attrs.Len() {
				close = append(close, ")")
				treeSExpr(w, attr, close...)
				break
			}
			treeSExpr(w, attr, close...)
		}
	case attrtree:
		w.WriteString("%s(attr_declaratin ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent()
		close = append(close, ")")
		exprSExpr(w, t.value, close...)
	case badtree:
		w.WriteString("%s(ERROR) ; ", w.Indentation())
		writeLocation(w, t.loc)
		w.Indent() // because of the dedent above
	default:
		w.WriteString("%s(ERROR) ; ", w.Indentation())
		w.Indent() // because of the dedent above
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
