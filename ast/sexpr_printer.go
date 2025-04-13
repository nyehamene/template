package ast

import (
	"fmt"
	"strings"
)

type SExprPrinterContext interface {
	Location(start, end int) string

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
	ns          *Namespace
}

func (p *sexprPrinter) lineOffsetFor(offset int) (line int, lineOffset int) {
	for i, l := range p.ns.lines {
		if l >= offset {
			break
		}
		line = i
		lineOffset = l
	}
	// 1 based line numbers
	line += 1
	return
}

func (p *sexprPrinter) Location(start, end int) string {
	startLine, startOffset := p.lineOffsetFor(start)
	endLine, endOffset := p.lineOffsetFor(end)

	startCol := start - startOffset
	endCol := end - endOffset

	if start == 0 && startCol == 0 {
		startCol = 1
	}

	// 1 based col numbers
	// startCol += 1
	// endCol += 1

	return fmt.Sprintf("%d, %d - %d, %d", startLine, startCol, endLine, endCol)
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

func PrintSExpr(n *Namespace) string {
	p := sexprPrinter{
		ns: n,
	}

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
		if len(chunk) > 2 {
			_ = chunk
		}
		chunks = append(chunks, chunk)
	}

	for _, chunk := range chunks {
		length := len(chunk)
		if length == 0 {
			continue
		}
		if length > 2 {
			fmt.Println(strings.Join(chunk, ""))
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

	return fmtLines.String()
}
