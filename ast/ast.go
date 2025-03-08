package ast

import "temlang/tem/token"

func EntrySame[T any](k T, v T) Entry[T, T] {
	return Entry[T, T]{key: k, val: v}
}

type NamespaceFile struct {
	Name     string
	Path     string
	Src      string
	Pkg      PackageDecl
	tokens   []token.Token
	tokenLen int
	texts    []string
	vars     []VarDecl
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
}

func (f *NamespaceFile) addToken(tok *TokenIndex, val token.Token) {
	index := len(f.tokens)
	f.tokens = append(f.tokens, val)
	f.tokenLen += 1
	tok.token.index = index
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
	offset := u.idents.text.index
	end := u.idents.text.len
	txts := f.texts[offset : offset+end]
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

// type TypeDecl struct {
// 	Decl
//	offset Index
// }

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

type Entry[K any, V any] struct {
	key K
	val V
}

func (e Entry[K, _]) Key() K {
	return e.key
}

func (e Entry[_, V]) Val() V {
	return e.val
}

func (d *RecordDecl) SetFields(f *NamespaceFile, entries []Entry[token.Token, token.Token]) {
	txtIndex := len(f.texts)
	varIndex := len(f.vars)
	for _, e := range entries {
		ident := e.key
		ty := e.val

		v := VarDecl{}
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

	index := d.fields.token.index
	length := d.fields.token.len
	vars := f.vars[index : index+length]
	return vars
}

type VarDecl struct {
	Decl
}

// type TemplDecl struct {
// 	Decl
// 	Params Index
// 	Start  Index
// 	End    Index
//	offset Index
// }
