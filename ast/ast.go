package ast

import "temlang/tem/token"

type NamespaceFile struct {
	Name     string
	Path     string
	Src      string
	Pkg      PackageDecl
	tokens   []token.Token
	tokenLen int
	texts    []string
	// imports      []ImportDecl
	// usings       []UsingDecl
	// alias        []AliasDecl
	// records      []RecordDecl
	// templs       []TemplDecl
	// vars         []VarDecl
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

	if name, ok := token.Keyword(tok.Kind); ok {
		return name, true
	}

	if tok.Kind > token.SymbolBegin && tok.Kind < token.SymbolEnd {
		return tok.String(), true
	}

	lexemeStart := tok.Start()
	lexemeEnd := tok.End()
	lexeme := string(n.Src[lexemeStart:lexemeEnd])
	return lexeme, true
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
	txt, ok := f.text(tok)
	if !ok {
		// TBD: maybe raise some kind of error
	}

	index := len(f.tokens)
	f.tokens = append(f.tokens, tok)
	d.type_.token.index = index

	index = len(f.texts)
	f.texts = append(f.texts, txt)
	d.type_.text.index = index
}

func (p Decl) Type(f NamespaceFile) string {
	txt := f.texts[p.type_.text.index]
	return txt
}

func (d *Decl) SetIdent(n *NamespaceFile, tok token.Token) {
	txt, ok := n.text(tok)
	if !ok {
		// TBD: maybe raise some kind of error
	}

	index := len(n.tokens)
	n.tokens = append(n.tokens, tok)
	d.ident.token.index = index

	index = len(n.texts)
	n.texts = append(n.texts, txt)
	d.ident.text.index = index
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
	txt, ok := f.text(tok)
	if !ok {
		// TBD: maybe raise some kind of error
	}

	index := len(f.tokens)
	f.tokens = append(f.tokens, tok)
	d.name.token.index = index

	index = len(f.texts)
	f.texts = append(f.texts, txt)
	d.name.text.index = index
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
	txt, ok := f.text(tok)
	if !ok {
		// TBD: maybe raise some kind of error
	}

	index := len(f.tokens)
	f.tokens = append(f.tokens, tok)
	d.templ.token.index = index

	index = len(f.texts)
	f.texts = append(f.texts, txt)
	d.templ.text.index = index
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
	txt, ok := f.text(tok)
	if !ok {
		// TBD: maybe raise some kind of error
	}

	index := len(f.tokens)
	f.tokens = append(f.tokens, tok)
	d.path.token.index = index

	index = len(f.texts)
	f.texts = append(f.texts, txt)
	d.path.text.index = index
}

func (i ImportDecl) Path(f NamespaceFile) string {
	txt := f.texts[i.path.text.index]
	return txt
}

type UsingDecl struct {
	Decl
	idents TokenIndex
	pkg    TokenIndex
}

func (d *UsingDecl) SetPkg(f NamespaceFile, tok token.Token) {
	txt, ok := f.text(tok)
	if !ok {
		// TBD: maybe raise some kind of error
	}

	index := len(f.tokens)
	f.tokens = append(f.tokens, tok)
	d.pkg.token.index = index

	index = len(f.texts)
	f.texts = append(f.texts, txt)
	d.pkg.text.index = index
}

func (u UsingDecl) Pkg(f NamespaceFile) string {
	txt := f.texts[u.pkg.text.index]
	return txt
}

func (u UsingDecl) Idents(f NamespaceFile) []string {
	txts := f.texts[u.idents.text.index:u.idents.text.len]
	return txts
}

// type TypeDecl struct {
// 	Decl
//	offset Index
// }

// type AliasDecl struct {
// 	Vars   TokenSlice
// 	Target Index
//	offset Index
// }

// type RecordDecl struct {
// 	Decl
// 	Fields TokenSlice
//	offset Index
// }

// type VarDecl struct {
// 	Decl
//	offset Index
// }

// type TemplDecl struct {
// 	Decl
// 	Params Index
// 	Start  Index
// 	End    Index
//	offset Index
// }
