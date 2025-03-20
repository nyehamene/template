package ast

import (
	"strings"
	"temlang/tem/token"
)

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

func (d *hasType) SetType(f *Namespace, tok token.Token) {
	f.addToken(&d.type_, tok)
	f.addText(&d.type_, tok)
}

func (p hasType) Type(f *Namespace) string {
	txt := f.getTextOne(p.type_)
	return txt
}

type hasIdents struct {
	idents TokenIndex
}

func (d *hasIdents) SetIdents(f *Namespace, toks *token.TokenStack) {
	tokOffset := len(f.tokens)
	txtOffset := len(f.texts)
	size := toks.Len()

	for tok, ok := toks.Pop(); ok; {
		f.addToken(&TokenIndex{}, tok)
		f.addText(&TokenIndex{}, tok)
	}

	d.idents.token.index = tokOffset
	d.idents.token.len = size

	d.idents.text.index = txtOffset
	d.idents.text.len = size
}

func (d hasIdents) Idents(f *Namespace) []string {
	txts := f.getText(d.idents)
	return txts
}

type hasName struct {
	name TokenIndex
}

func (d *hasName) SetName(f *Namespace, tok token.Token) {
	f.addToken(&d.name, tok)
	f.addText(&d.name, tok)
}

func (d hasName) Name(f *Namespace) string {
	txt := f.getTextOne(d.name)
	return txt
}

type Decl struct {
	hasIdents
	hasType
	// LeadingDocs  Index
	// TrailingDocs Index
	// Tags         Index
}

type PackageDecl struct {
	Decl
	hasName
	directive TokenIndex
}

func (d *PackageDecl) SetDirective(f *Namespace, tok token.Token) {
	f.addToken(&d.directive, tok)
	f.addText(&d.directive, tok)
}

func (d PackageDecl) Templ(f *Namespace) string {
	txt := f.getTextOne(d.directive)
	return txt
}

type ImportDecl struct {
	Decl
	hasName
}

type UsingDecl struct {
	Decl
	pkg TokenIndex
}

func (d *UsingDecl) SetPkg(f *Namespace, tok token.Token) {
	f.addToken(&d.pkg, tok)
	f.addText(&d.pkg, tok)
}

func (u UsingDecl) Pkg(f *Namespace) string {
	txt := f.getTextOne(u.pkg)
	return txt
}

type TypeDecl struct {
	Decl
	target TokenIndex
}

func (d *TypeDecl) SetTarget(f *Namespace, tok token.Token) {
	f.addToken(&d.target, tok)
	f.addText(&d.target, tok)
}

func (d TypeDecl) Target(f *Namespace) string {
	txt := f.getTextOne(d.target)
	return txt
}

type RecordDecl struct {
	Decl
	fields TokenSlice
}

func (d *RecordDecl) SetFields(f *Namespace, fields TokenSlice) {
	d.fields = fields
}

func (d RecordDecl) Fields(f *Namespace) []VarDecl {
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

func (d *DocDecl) SetContent(f *Namespace, strs ...token.Token) {
	tokIndex := len(f.tokens)
	txtIndex := len(f.texts)
	size := len(strs)
	for _, tok := range strs {
		if kind := tok.Kind(); kind != token.String && kind != token.TextBlock {
			panic("doc content must either be a string or text block")
		}
		f.addToken(&d.content, tok)
		f.addText(&d.content, tok)
	}
	d.content.token.index = tokIndex
	d.content.token.len = size

	d.content.text.index = txtIndex
	d.content.text.len = size
}

func (d DocDecl) Content(f *Namespace) string {
	var (
		index = d.content.text.index
		end   = index + d.content.text.len
		txts  = f.texts[index:end]
	)
	return strings.Join(txts, "\n")
}

type TagDecl struct {
	hasIdents
	attrs TokenSlice
}

func (d *TagDecl) SetAttrs(f *Namespace, attrs TokenSlice) {
	d.attrs = attrs
}

func (d TagDecl) Attrs(f *Namespace) []AttrDecl {
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

func (d *AttrDecl) SetValue(f *Namespace, tok token.Token) {
	if kind := tok.Kind(); kind != token.String && kind != token.TextBlock {
		panic("unreachable")
	}
	f.addToken(&d.value, tok)
	f.addText(&d.value, tok)
}

func (d AttrDecl) Value(f *Namespace) string {
	txt := f.getTextOne(d.value)
	return txt
}

type TemplDecl struct {
	Decl
	params TokenSlice
}

func (d *TemplDecl) SetParams(f *Namespace, params TokenSlice) {
	d.params = params
}

func (d TemplDecl) Params(f *Namespace) []VarDecl {
	var (
		ident  = d.params.index
		end    = ident + d.params.len
		params = f.vars[ident:end]
	)
	return params
}
