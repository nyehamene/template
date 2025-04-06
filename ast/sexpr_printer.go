package ast

import (
	"fmt"
	"strings"
	"temlang/tem/token"
)

type SExprPrinterContext interface {
	Location(token.Token) token.Location

	WriteString(string, ...any)

	Indent()

	Dedent()

	Indentation() string
}

type SExpressionPrinter interface {
	WriteSExpr(SExprPrinterContext)
}

type sexprPrinter struct {
	sb          strings.Builder
	indent      int
	indentation string
	// ns     *Namespace
}

func (p *sexprPrinter) Location(token.Token) token.Location {
	l := token.Location{} // TODO get location from namespace
	return l
}

func (p *sexprPrinter) Indent() {
	p.indent += 2
	i := strings.Repeat(" ", p.indent)
	p.indentation = i
}

func (p *sexprPrinter) Dedent() {
	p.indent -= 2
	if p.indent < 0 {
		p.indent = 0
	}
	i := strings.Repeat(" ", p.indent)
	p.indentation = i
}

func (p *sexprPrinter) Indentation() string {
	return p.indentation
}

func (p *sexprPrinter) WriteString(str string, obj ...any) {
	s := fmt.Sprintf(str, obj...)
	p.sb.WriteString(s)
}

func PrintSExpr(n *Namespace) {
	p := sexprPrinter{}

	for _, d := range n.decl {
		d.WriteSExpr(&p)
	}

	fmt.Println(p.sb.String())
}
