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
	tokenLen int
	texts    []string
	vars     []VarDecl
	attrs    []AttrDecl
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

type Decl struct {
	ident TokenIndex
	type_ TokenIndex
	// LeadingDocs  Index
	// TrailingDocs Index
	// Tags         Index
}

func (d *Decl) SetType(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.type_, tok)
	f.addText(&d.type_, tok)
}

func (p Decl) Type(f NamespaceFile) string {
	txt := f.texts[p.type_.text.index]
	return txt
}

func (d *Decl) SetIdent(f *NamespaceFile, tok token.Token) {
	f.addToken(&d.ident, tok)
	f.addText(&d.ident, tok)
}

func (u Decl) Ident(f NamespaceFile) string {
	txt := f.texts[u.ident.text.index]
	return txt
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

type manyIdentDecl struct {
	Decl
	idents TokenIndex
}

func (d *manyIdentDecl) SetIdents(f *NamespaceFile, toks []token.Token) {
	tokOffset := len(f.tokens)
	txtOffset := len(f.texts)

	for _, tok := range toks {
		txt, ok := f.text(tok)
		if !ok {
			panic("unreachable")
		}
		f.tokens = append(f.tokens, tok)
		f.texts = append(f.texts, txt)
	}

	d.idents.token.index = tokOffset
	d.idents.token.len = len(toks)

	d.idents.text.index = txtOffset
	d.idents.text.len = len(toks)
}

func (u manyIdentDecl) Idents(f NamespaceFile) []string {
	var (
		offset = u.idents.text.index
		end    = offset + u.idents.text.len
		txts   = f.texts[offset:end]
	)
	return txts
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
	manyIdentDecl
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
	manyIdentDecl
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
	manyIdentDecl
	fields TokenIndex
}

func (d *RecordDecl) SetFields(f *NamespaceFile, entries []Entry[token.Token, token.Token]) {
	txtIndex := len(f.texts)
	varIndex := len(f.vars)
	for _, e := range entries {
		var (
			ident = e.key
			ty    = e.val
			v     = VarDecl{}
		)
		f.addToken(&v.ident, ident)
		f.addText(&v.ident, ident)

		f.addToken(&v.type_, ty)
		f.addText(&v.type_, ty)
		f.vars = append(f.vars, v)
	}
	d.fields.token.index = varIndex
	d.fields.text.index = txtIndex

	d.fields.token.len = len(entries)
	d.fields.text.len = len(entries)
}

func (d RecordDecl) Fields(f NamespaceFile) []VarDecl {
	if d.fields.token.len != d.fields.text.len {
		panic("unreachable")
	}
	var (
		index = d.fields.token.index
		end   = index + d.fields.token.len
		vars  = f.vars[index:end]
	)
	return vars
}

type VarDecl struct {
	Decl
}

type hasIdents struct {
	idents TokenIndex
}

func (d *hasIdents) SetIdents(f *NamespaceFile, toks []token.Token) {
	tokOffset := len(f.tokens)
	txtOffset := len(f.texts)

	for _, tok := range toks {
		txt, ok := f.text(tok)
		if !ok {
			panic("unreachable")
		}
		f.tokens = append(f.tokens, tok)
		f.texts = append(f.texts, txt)
	}

	d.idents.token.index = tokOffset
	d.idents.token.len = len(toks)

	d.idents.text.index = txtOffset
	d.idents.text.len = len(toks)
}

func (d hasIdents) Idents(f NamespaceFile) []string {
	var (
		offset = d.idents.text.index
		end    = offset + d.idents.text.len
		txts   = f.texts[offset:end]
	)
	return txts
}

type DocDecl struct {
	hasIdents
	content TokenIndex
}

func (d *DocDecl) SetContent(f *NamespaceFile, strs ...token.Token) {
	tokIndex := len(f.tokens)
	txtIndex := len(f.texts)
	for _, tok := range strs {
		if kind := tok.Kind(); kind != token.String && kind != token.TextBlock {
			panic("unreachable")
		}
		f.addToken(&d.content, tok)
		f.addText(&d.content, tok)
	}
	d.content.token.index = tokIndex
	d.content.token.len = len(strs)

	d.content.text.index = txtIndex
	d.content.text.len = len(strs)
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
	attrs TokenIndex
}

func (d *TagDecl) SetAttrs(f *NamespaceFile, attrs []Entry[[]token.Token, token.Token]) {
	txtIndex := len(f.texts)
	attrIndex := len(f.attrs)

	for _, attr := range attrs {
		var (
			ad         = AttrDecl{}
			ty         = attr.val
			adIndexTxt = len(f.texts)
			adIndexTok = len(f.tokens)
		)
		for _, k := range attr.key {
			f.addToken(&ad.idents, k)
			f.addText(&ad.idents, k)
		}

		f.addToken(&ad.value, ty)
		f.addText(&ad.value, ty)

		ad.idents.token.index = adIndexTok
		ad.idents.text.index = adIndexTxt

		ad.idents.token.len = len(attr.key)
		ad.idents.text.len = len(attr.key)

		f.attrs = append(f.attrs, ad)
	}
	d.attrs.token.index = attrIndex
	d.attrs.text.index = txtIndex

	d.attrs.token.len = len(attrs)
	d.attrs.text.len = len(attrs)
}

func (d TagDecl) Attrs(f NamespaceFile) []AttrDecl {
	if d.attrs.token.len != d.attrs.text.len {
		panic("unreachable")
	}
	var (
		index = d.attrs.token.index
		end   = index + d.attrs.token.len
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

// type TemplDecl struct {
// 	Decl
// 	Params Index
// 	Start  Index
// 	End    Index
//	offset Index
// }
