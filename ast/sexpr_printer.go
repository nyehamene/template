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

	sexpr := make([]string, 0, p.sb.Len())
	chunks := make([][]string, 0, p.sb.Len())
	fmtLines := strings.Builder{}

	fmtLines.Grow(len(sexpr))

	str := p.sb.String()
	lines := strings.Split(str, "\n")

	longestSExpr := 0
	for _, line := range lines {
		chunk := strings.Split(line, ";")
		if l := len(chunk[0]); l > longestSExpr {
			longestSExpr = l
		}
		chunks = append(chunks, chunk)
	}

	for _, chunk := range chunks {
		length := len(chunk)
		if length == 0 {
			continue
		}
		if length > 2 {
			panic("Invalid sexpr")
		}
		if length == 1 {
			fmtLines.WriteString(chunk[0])
			fmtLines.WriteString("\n")
			continue
		}

		e := chunk[0]
		l := chunk[1]
		length = len(e)
		var str string

		if length == longestSExpr {
			str = e
		} else {
			padSize := longestSExpr - len(e)
			pad := strings.Repeat(" ", padSize)
			str = fmt.Sprintf("%s%s", e, pad)
		}
		fmtLines.WriteString(fmt.Sprintf("%s  ; %s\n", str, l))
	}

	fmt.Println(fmtLines.String())
}
