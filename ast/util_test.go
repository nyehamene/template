package ast

import "temlang/tem/token"

// the token kind does not matter (dnm)
func tokenDNM(start, end int) token.Token {
	var doesNotMatter token.Kind
	return token.New(doesNotMatter, start, end)
}

type pos struct {
	start int
	end   int
}

func tokenMany(p ...pos) []token.Token {
	toks := []token.Token{}
	for _, v := range p {
		tok := tokenDNM(v.start, v.end)
		toks = append(toks, tok)
	}
	return toks
}

type singletoken interface {
	Set(*NamespaceFile, token.Token)
	Get(NamespaceFile) string
}

func (h *hasType) Set(f *NamespaceFile, tok token.Token) {
	h.SetType(f, tok)
}

func (h hasType) Get(f NamespaceFile) string {
	return h.Type(f)
}

func (h *hasName) Set(f *NamespaceFile, tok token.Token) {
	h.SetName(f, tok)
}

func (h hasName) Get(f NamespaceFile) string {
	return h.Name(f)
}

func (d *PackageDecl) Set(f *NamespaceFile, tok token.Token) {
	d.SetTempl(f, tok)
}

func (d PackageDecl) Get(f NamespaceFile) string {
	return d.Templ(f)
}

func (d *UsingDecl) Set(f *NamespaceFile, tok token.Token) {
	d.SetPkg(f, tok)
}

func (d UsingDecl) Get(f NamespaceFile) string {
	return d.Pkg(f)
}

func (d *AliasDecl) Set(f *NamespaceFile, tok token.Token) {
	d.SetTarget(f, tok)
}

func (d AliasDecl) Get(f NamespaceFile) string {
	return d.Target(f)
}

func (d *AttrDecl) Set(f *NamespaceFile, tok token.Token) {
	d.SetValue(f, tok)
}

func (d AttrDecl) Get(f NamespaceFile) string {
	return d.Value(f)
}

type manytoken interface {
	Set(*NamespaceFile, []token.Token)
	Get(NamespaceFile) []string
}

func (h *hasIdents) Set(f *NamespaceFile, toks []token.Token) {
	h.SetIdents(f, toks)
}

func (h hasIdents) Get(f NamespaceFile) []string {
	return h.Idents(f)
}

type manydecl[T any] interface {
	Set(*NamespaceFile, []Entry[[]token.Token, token.Token])
	Get(NamespaceFile) []T
}

func (d *TemplDecl) Set(f *NamespaceFile, toks []Entry[[]token.Token, token.Token]) {
	d.SetParams(f, toks)
}

func (d TemplDecl) Get(f NamespaceFile) []VarDecl {
	return d.Params(f)
}

func (d *RecordDecl) Set(f *NamespaceFile, toks []Entry[[]token.Token, token.Token]) {
	d.SetFields(f, toks)
}

func (d RecordDecl) Get(f NamespaceFile) []VarDecl {
	return d.Fields(f)
}

func (d *TagDecl) Set(f *NamespaceFile, toks []Entry[[]token.Token, token.Token]) {
	d.SetAttrs(f, toks)
}

func (d TagDecl) Get(f NamespaceFile) []AttrDecl {
	return d.Attrs(f)
}
