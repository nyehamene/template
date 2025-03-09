package ast

import (
	"strings"
	"temlang/tem/token"
)

type NamespaceFile struct {
	Name     string
	Path     string
	Src      string
	Pkg      PackageDecl
	tokens   []token.Token
	texts    []string
	vars     []VarDecl
	attrs    []AttrDecl
	tokenLen int
	textLen  int
	// imports      []ImportDecl
	// usings       []UsingDecl
	// alias        []AliasDecl
	// records      []RecordDecl
	// templs       []TemplDecl
	// docs         []TokenSlice
	// tags         []TokenSlice
}

func (n *NamespaceFile) Init() {
	// the value at index 0 is used to represent tokens that
	// do not exist
	var noToken token.Token
	var noText string
	// To avoid too many memory allocations assume that 65%
	// of the size len of n.Src can contain all relevant tokens
	srcSizePercent := 0.65
	c := float64(len(n.Src)) * srcSizePercent
	// initialize the namespace file
	n.tokens = make([]token.Token, 0, int(c))
	n.tokens = append(n.tokens, noToken)
	n.texts = make([]string, 0, int(c))
	n.texts = append(n.texts, noText)
}

func (n *NamespaceFile) text(tok token.Token) (string, bool) {
	if l := len(n.Src); tok.Start() > l || tok.End() > l {
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
	lexeme := string(n.Src[lexemeStart:lexemeEnd])
	return lexeme, true
}

func (f *NamespaceFile) addText(tok *TokenIndex, val token.Token) {
	txt, ok := f.text(val)
	if !ok {
		// TBD: maybe raise some kind of error
		panic("unreachable")
	}

	index := len(f.texts)
	f.texts = append(f.texts, txt)
	f.textLen += 1
	tok.text.index = index
	tok.text.len = 1
}

func (f *NamespaceFile) addToken(tok *TokenIndex, val token.Token) {
	index := len(f.tokens)
	f.tokens = append(f.tokens, val)
	f.tokenLen += 1
	tok.token.index = index
	tok.token.len = 1
}

type TokenIndex struct {
	token TokenSlice
	text  TokenSlice
}

type TokenSlice struct {
	index int
	len   int
}

type hasType struct {
	type_ TokenIndex
}

func (d *hasType) SetType(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.type_, tok)
	f.addText(&d.type_, tok)
}

func (p hasType) Type(f NamespaceFile) string {
	txt := f.texts[p.type_.text.index]
	return txt
}

type hasIdents struct {
	idents TokenIndex
}

func (d *hasIdents) SetIdents(f *NamespaceFile, toks []token.Token) {
	tokOffset := len(f.tokens)
	txtOffset := len(f.texts)
	size := len(toks)

	for _, tok := range toks {
		txt, ok := f.text(tok)
		if !ok {
			panic("unreachable")
		}
		f.tokens = append(f.tokens, tok)
		f.texts = append(f.texts, txt)
	}

	d.idents.token.index = tokOffset
	d.idents.token.len = size

	d.idents.text.index = txtOffset
	d.idents.text.len = size
}

func (d hasIdents) Idents(f NamespaceFile) []string {
	var (
		offset = d.idents.text.index
		end    = offset + d.idents.text.len
		txts   = f.texts[offset:end]
	)
	return txts
}

type Decl struct {
	hasIdents
	hasType
	// LeadingDocs  Index
	// TrailingDocs Index
	// Tags         Index
}

type NamedDecl struct {
	Decl
	name TokenIndex
}

func (d *NamedDecl) SetName(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.name, tok)
	f.addText(&d.name, tok)
}

func (p NamedDecl) Name(f NamespaceFile) string {
	txt := f.texts[p.name.text.index]
	return txt
}

type PackageDecl struct {
	NamedDecl
	templ TokenIndex
}

func (d *PackageDecl) SetTempl(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.templ, tok)
	f.addText(&d.templ, tok)
}

func (p PackageDecl) Templ(f NamespaceFile) string {
	txt := f.texts[p.templ.text.index]
	return txt
}

type ImportDecl struct {
	Decl
	path TokenIndex
}

func (d *ImportDecl) SetPath(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.path, tok)
	f.addText(&d.path, tok)
}

func (i ImportDecl) Path(f NamespaceFile) string {
	txt := f.texts[i.path.text.index]
	return txt
}

type UsingDecl struct {
	Decl
	pkg TokenIndex
}

func (d *UsingDecl) SetPkg(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.pkg, tok)
	f.addText(&d.pkg, tok)
}

func (u UsingDecl) Pkg(f NamespaceFile) string {
	txt := f.texts[u.pkg.text.index]
	return txt
}

type AliasDecl struct {
	Decl
	target TokenIndex
}

func (d *AliasDecl) SetTarget(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.target, tok)
	f.addText(&d.target, tok)
}

func (d AliasDecl) Target(f NamespaceFile) string {
	txt := f.texts[d.target.text.index]
	return txt
}

type RecordDecl struct {
	Decl
	fields TokenSlice
}

func (d *RecordDecl) SetFields(f *NamespaceFile, entries []Entry[[]token.Token, token.Token]) {
	varIndex := len(f.vars)
	for _, e := range entries {
		var (
			ident = e.key
			ty    = e.val
			v     = VarDecl{}
		)
		v.SetIdents(f, ident)
		v.SetType(f, ty)
		f.vars = append(f.vars, v)
	}
	d.fields.index = varIndex
	d.fields.len = len(entries)
}

func (d RecordDecl) Fields(f NamespaceFile) []VarDecl {
	var (
		index = d.fields.index
		end   = index + d.fields.len
		vars  = f.vars[index:end]
	)
	return vars
}

type VarDecl struct {
	Decl
}

type DocDecl struct {
	hasIdents
	content TokenIndex
}

func (d *DocDecl) SetContent(f *NamespaceFile, strs ...token.Token) {
	tokIndex := len(f.tokens)
	txtIndex := len(f.texts)
	size := len(strs)
	for _, tok := range strs {
		if kind := tok.Kind(); kind != token.String && kind != token.TextBlock {
			panic("doc content must either a string or text block")
		}
		f.addToken(&d.content, tok)
		f.addText(&d.content, tok)
	}
	d.content.token.index = tokIndex
	d.content.token.len = size

	d.content.text.index = txtIndex
	d.content.text.len = size
}

func (d DocDecl) Content(f NamespaceFile) string {
	var (
		sb    = strings.Builder{}
		index = d.content.text.index
		end   = index + d.content.text.len
		txts  = f.texts[index:end]
		l     = 0
	)
	sb.Grow(d.content.text.len)

	for _, txt := range txts {
		if l > 0 {
			sb.WriteString("\n")
		}
		l, _ = sb.WriteString(txt)
	}
	return sb.String()
}

type TagDecl struct {
	hasIdents
	attrs TokenSlice
}

func (d *TagDecl) SetAttrs(f *NamespaceFile, attrs []Entry[[]token.Token, token.Token]) {
	attrIndex := len(f.attrs)
	for _, attr := range attrs {
		ad := AttrDecl{}
		ad.SetIdents(f, attr.key)
		ad.SetValue(f, attr.val)
		f.attrs = append(f.attrs, ad)
	}
	d.attrs.index = attrIndex
	d.attrs.len = len(attrs)
}

func (d TagDecl) Attrs(f NamespaceFile) []AttrDecl {
	var (
		index = d.attrs.index
		end   = index + d.attrs.len
		attrs = f.attrs[index:end]
	)
	return attrs
}

type AttrDecl struct {
	hasIdents
	value TokenIndex
}

func (d *AttrDecl) SetValue(f *NamespaceFile, tok token.Token) {
	if kind := tok.Kind(); kind != token.String && kind != token.TextBlock {
		panic("unreachable")
	}
	f.addToken(&d.value, tok)
	f.addText(&d.value, tok)
}

func (d AttrDecl) Value(f NamespaceFile) string {
	txt := f.texts[d.value.text.index]
	return txt
}

type TemplDecl struct {
	hasIdents
	hasType
	params TokenSlice
}

func (d *TemplDecl) SetParams(f *NamespaceFile, params []Entry[[]token.Token, token.Token]) {
	if len(params) != 1 {
		panic("templ literal parameter list must contain only one param")
	}

	varIndex := len(f.vars)
	varlen := len(params)

	for _, tok := range params {
		v := VarDecl{}
		v.SetIdents(f, tok.key)
		v.SetType(f, tok.val)
		f.vars = append(f.vars, v)
	}

	d.params.index = varIndex
	d.params.len = varlen
}

func (d TemplDecl) Params(f NamespaceFile) []VarDecl {
	var (
		ident  = d.params.index
		end    = ident + d.params.len
		params = f.vars[ident:end]
	)
	return params
}
