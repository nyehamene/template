package ast

import (
	"temlang/tem/token"
)

func New(name string, src []byte) *Namespace {
	n := Namespace{src: src}
	return &n
}

type Namespace struct {
	pkg      PackageDecl
	texts    []string
	vars     []VarDecl
	attrs    []AttrDecl
	imports  []ImportDecl
	usings   []UsingDecl
	types    []TypeDecl
	records  []RecordDecl
	templs   []TemplDecl
	docs     []DocDecl
	tags     []TagDecl
	tokens   []token.Token
	Name     string
	Path     string
	src      []byte // TODO remove src
	tokenLen int
	textLen  int
}

func (n *Namespace) Init() {
	// To avoid too many memory allocations assume that 65%
	// of the size len of n.Src can contain all relevant tokens
	srcSizePercent := 0.65
	c := float64(len(n.src)) * srcSizePercent
	// initialize the namespace file
	n.tokens = make([]token.Token, 0, int(c))
	n.texts = make([]string, 0, int(c))
}

func (n *Namespace) text(tok token.Token) (string, bool) {
	if l := len(n.src); tok.Start() > l || tok.End() > l {
		return "", false
	}

	if name, ok := token.Keyword(tok.Kind()); ok {
		return name, true
	}

	if tok.Kind() > token.SymbolBegin && tok.Kind() < token.SymbolEnd {
		return tok.String(), true
	}

	lexemeStart := tok.Start()
	lexemeEnd := tok.End()
	lexeme := string(n.src[lexemeStart:lexemeEnd])
	return lexeme, true
}

func (f *Namespace) addText(tok *TokenIndex, val token.Token) {
	txt, ok := f.text(val)
	if !ok {
		// TBD: maybe return some kind of error
		panic("unreachable")
	}

	index := len(f.texts)
	f.texts = append(f.texts, txt)
	f.textLen += 1
	tok.text.Index = index
	tok.text.Len = 1
}

func (f *Namespace) addToken(tok *TokenIndex, val token.Token) {
	index := len(f.tokens)
	f.tokens = append(f.tokens, val)
	f.tokenLen += 1
	tok.token.Index = index
	tok.token.Len = 1
}

func (f Namespace) getText(tok TokenIndex) []string {
	index := tok.text.Index
	end := index + tok.text.Len
	txts := f.texts[index:end]
	return txts
}

func (f Namespace) getTextOne(tok TokenIndex) string {
	txts := f.getText(tok)
	if len(txts) > 1 {
		panic("expected only one text")
	}
	return txts[0]
}

func (n *Namespace) Pkg() PackageDecl {
	return n.pkg
}

func (n *Namespace) SetPkg(p PackageDecl) {
	n.pkg = p
}

func (n *Namespace) AddImport(d ImportDecl) {
	n.imports = append(n.imports, d)
}

func (n *Namespace) AddUsing(d UsingDecl) {
	n.usings = append(n.usings, d)
}

func (n *Namespace) AddType(d TypeDecl) {
	n.types = append(n.types, d)
}

func (n *Namespace) AddRecord(d RecordDecl) {
	n.records = append(n.records, d)
}

func (n *Namespace) AddTempl(d TemplDecl) {
	n.templs = append(n.templs, d)
}

func (n *Namespace) AddVar(d VarDecl) {
	n.vars = append(n.vars, d)
}

func (n *Namespace) AddDoc(d DocDecl) {
	n.docs = append(n.docs, d)
}

func (n *Namespace) AddTag(d TagDecl) {
	n.tags = append(n.tags, d)
}

func (n *Namespace) AddAttr(d AttrDecl) {
	n.attrs = append(n.attrs, d)
}

func (n *Namespace) VarLen() int {
	return len(n.vars)
}

func (n *Namespace) AttrLen() int {
	return len(n.attrs)
}
