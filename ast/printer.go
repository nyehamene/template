package ast

import (
	"fmt"
	"strings"
	"temlang/tem/token"
)

type tokenLocation map[token.Token]location

func NewPrinter(f *Namespace) *Printer {
	pm := newPrinted(f)
	ns := []*printed{pm}

	p := &Printer{
		namespaces: ns,
	}

	for _, n := range ns {
		p.init(n)
	}

	return p
}

func Empty(t token.Token) bool {
	return t.Start() == 0 && t.End() == 0
}

type position struct {
	line, col int
}

type location struct {
	start, end position
}

type Printer struct {
	namespaces []*printed
}

func (p *Printer) init(pm *printed) {
	p.addLines(pm)
}

func (p *Printer) addLines(pm *printed) {
	src := pm.src
	if src == "" {
		return
	}

	pm.rdLines = 1
	pm.rdCols = 1

	for _, tok := range pm.toks {
		if Empty(tok) {
			continue
		}
		pm.addPos(tok)
	}
}

func (p *Printer) Print() {
	for _, ns := range p.namespaces {
		txt := ns.print()
		fmt.Println(txt)
	}
}

func newPrinted(f *Namespace) *printed {
	tokens := make([]token.Token, f.tokenLen)
	imps := make([]ImportDecl, len(f.imports))
	usings := make([]UsingDecl, len(f.usings))

	copy(tokens, f.tokens)
	copy(imps, f.imports)
	copy(usings, f.usings)

	sb := &strings.Builder{}
	sb.Grow(f.tokenLen)

	maxTagWidth := 30

	return &printed{
		toks:   tokens,
		imps:   imps,
		usings: usings,
		file:   f.Name,
		pkg:    f.Pkg,
		sb:     sb,
		src:    f.src,
		loc:    tokenLocation{},
		width:  maxTagWidth,
	}
}

var noop = func(int) {}

type writeFunc func() func(int)

type printed struct {
	pkg        PackageDecl
	file       string
	src        string
	rdOffset   int
	rdLines    int
	rdCols     int
	identLevel int
	width      int
	loc        tokenLocation
	toks       []token.Token
	imps       []ImportDecl
	usings     []UsingDecl
	sb         *strings.Builder
}

func (p *printed) addPos(tok token.Token) {
	count := func(end int) position {
		if end == 0 {
			return position{p.rdLines, p.rdCols}
		}
		var (
			offset = p.rdOffset
			lines  = p.rdLines
			cols   = p.rdCols
			src    = p.src[offset:end]
		)
		for _, r := range src {
			if r == '\n' {
				lines += 1
				cols = 0
				continue
			}
			cols += 1
		}
		p.rdOffset = end
		p.rdLines = lines
		p.rdCols = cols
		return position{lines, cols}
	}

	start := count(tok.Start())
	end := count(tok.End())
	loc := location{start: start, end: end}
	p.loc[tok] = loc
}

func (p *printed) print() string {
	fst := p.toks[0]
	lst := p.toks[len(p.toks)-1]

	loc := p.writeContainer("source_file", fst, lst, func() func(int) {
		p.writePackage(p.pkg)(p.width - 1)
		p.writeImport(p.imps)(p.width)
		loc := p.writeUsing(p.usings)
		return loc
	})
	loc(p.width - 1)
	return p.sb.String()
}

func (pm *printed) writePackage(d PackageDecl) func(int) {
	fst := d.idents
	lst := d.name
	drt := d.directive
	ft := pm.toks[fst.token.index]
	lt := pm.toks[lst.token.index]
	return pm.writeContainer("package_declaration", ft, lt, func() func(int) {
		pm.writeIdents(fst)(pm.width)
		pm.writeDirective(drt)(pm.width)
		loc := pm.writeString(lst)
		return loc
	})
}

func (pm *printed) writeImport(ds []ImportDecl) func(int) {
	for i, d := range ds {
		fst := d.idents
		lst := d.name
		ft := pm.toks[fst.token.index]
		lt := pm.toks[lst.token.index]
		loc := pm.write("import_declaration", ft, lt)
		if i == len(ds)-1 {
			return loc
		}
		loc(pm.width)
	}
	return noop
}

func (pm *printed) writeUsing(ds []UsingDecl) func(int) {
	for i, d := range ds {
		fst := d.idents
		lst := d.pkg
		ft := pm.toks[fst.token.index]
		lt := pm.toks[lst.token.index]
		loc := pm.write("using_declaration", ft, lt)
		if i == len(ds)-1 {
			return loc
		}
		loc(pm.width)
	}
	return noop
}

func (pm *printed) writeIdents(t TokenIndex) func(int) {
	return pm.writeToken("identifier", t)
}

func (pm *printed) writeString(t TokenIndex) func(int) {
	return pm.writeToken("string", t)
}

func (pm *printed) writeDirective(t TokenIndex) func(int) {
	return pm.writeToken("directive", t)
}

func (pm *printed) writeContainer(tag string, fst, lst token.Token, fn writeFunc) func(int) {
	fl := pm.loc[fst]
	ll := pm.loc[lst]

	ident := strings.Repeat("  ", pm.identLevel)
	txt := fmt.Sprintf("\n%s(%s", ident, tag)
	padSize := pm.width - len(txt)
	pad := strings.Repeat(" ", padSize)

	pm.sb.WriteString(
		fmt.Sprintf("%s%s; [%-2d, %-2d] - [%-2d, %-2d]",
			txt, pad,
			fl.start.line, fl.start.col,
			ll.end.line, ll.end.col), // pkg: +1 for closing parens )
	)

	pm.identLevel += 1
	pos := fn()
	pm.sb.WriteString(")")
	pm.identLevel -= 1
	return pos
}

func (pm *printed) writeToken(tag string, t TokenIndex) func(int) {
	offset := t.token.index
	end := offset + t.token.len
	toks := pm.toks[offset:end]
	width := pm.width
	for i, tok := range toks {
		loc := pm.write(tag, tok, tok)
		if i == len(toks)-1 {
			return loc
		}
		loc(width)
	}
	return noop
}

func (pm *printed) write(tag string, fst, lst token.Token) func(int) {
	fl := pm.loc[fst]
	ll := pm.loc[lst]
	ident := strings.Repeat("  ", pm.identLevel)
	txt := fmt.Sprintf("\n%s(%s)", ident, tag)
	pm.sb.WriteString(txt)
	return func(width int) {
		padSize := width - len(txt)
		pad := strings.Repeat(" ", padSize)
		str := fmt.Sprintf("%s; [%-2d, %-2d] - [%-2d, %-2d]",
			pad,
			fl.start.line, fl.start.col,
			ll.end.line, ll.end.col)
		pm.sb.WriteString(str)
	}
}
