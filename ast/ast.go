package ast

import (
	"strings"
	"temlang/tem/dsa/stack"
	"temlang/tem/token"
)

type TokenIndex struct {
	token TokenSlice
	text  TokenSlice
}

type TokenSlice struct {
	Index int
	Len   int
}

type hasType struct {
	dtype TokenIndex
}

type hasIdents struct {
	idents TokenIndex
}

type hasName struct {
	name TokenIndex
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

type ImportDecl struct {
	Decl
	hasName
}

type UsingDecl struct {
	Decl
	pkg TokenIndex
}

type TypeDecl struct {
	Decl
	target TokenIndex
}

type RecordDecl struct {
	Decl
	fields TokenSlice
}

type VarDecl struct {
	Decl
}

type DocDecl struct {
	hasIdents
	content TokenIndex
}

type TagDecl struct {
	hasIdents
	attrs TokenSlice
}

type AttrDecl struct {
	hasIdents
	value TokenIndex
}

type TemplDecl struct {
	Decl
	params TokenSlice
}

func (p hasType) Type(f *Namespace) string {
	txt := f.getTextOne(p.dtype)
	return txt
}

func (d hasIdents) Idents(f *Namespace) []string {
	txts := f.getText(d.idents)
	return txts
}

func (d hasName) Name(f *Namespace) string {
	txt := f.getTextOne(d.name)
	return txt
}

func (d PackageDecl) Templ(f *Namespace) string {
	txt := f.getTextOne(d.directive)
	return txt
}

func (u UsingDecl) Pkg(f *Namespace) string {
	txt := f.getTextOne(u.pkg)
	return txt
}

func (d TypeDecl) Target(f *Namespace) string {
	txt := f.getTextOne(d.target)
	return txt
}

func (d RecordDecl) Fields(f *Namespace) []VarDecl {
	var (
		index = d.fields.Index
		end   = index + d.fields.Len
		vars  = f.vars[index:end]
	)
	return vars
}

func (d DocDecl) Content(f *Namespace) string {
	var (
		index = d.content.text.Index
		end   = index + d.content.text.Len
		txts  = f.texts[index:end]
	)
	return strings.Join(txts, "\n")
}

func (d TagDecl) Attrs(f *Namespace) []AttrDecl {
	var (
		index = d.attrs.Index
		end   = index + d.attrs.Len
		attrs = f.attrs[index:end]
	)
	return attrs
}

func (d AttrDecl) Value(f *Namespace) string {
	txt := f.getTextOne(d.value)
	return txt
}

func (d TemplDecl) Params(f *Namespace) []VarDecl {
	var (
		ident  = d.params.Index
		end    = ident + d.params.Len
		params = f.vars[ident:end]
	)
	return params
}

func (d *hasType) SetType(f *Namespace, tok token.Token) {
	f.addToken(&d.dtype, tok)
	f.addText(&d.dtype, tok)
}

func (d *hasIdents) SetIdents(f *Namespace, toks token.TokenStack) {
	tokOffset := len(f.tokens)
	txtOffset := len(f.texts)
	size := toks.Len()

	ignored := &TokenIndex{}
	for !toks.Empty() {
		tok, ok := toks.Pop()
		if !ok {
			panic("unreachable")
		}
		f.addToken(ignored, tok)
		f.addText(ignored, tok)
	}

	d.idents.token.Index = tokOffset
	d.idents.token.Len = size

	d.idents.text.Index = txtOffset
	d.idents.text.Len = size
}

func (d *hasName) SetName(f *Namespace, tok token.Token) {
	f.addToken(&d.name, tok)
	f.addText(&d.name, tok)
}

func (d *PackageDecl) SetDirective(f *Namespace, toks token.TokenStack) {
	for !toks.Empty() {
		tok, ok := toks.Pop()
		if !ok {
			break
		}
		f.addToken(&d.directive, tok)
		f.addText(&d.directive, tok)
	}
}

func (d *UsingDecl) SetPkg(f *Namespace, tok token.Token) {
	f.addToken(&d.pkg, tok)
	f.addText(&d.pkg, tok)
}

func (d *TypeDecl) SetTarget(f *Namespace, tok token.Token) {
	f.addToken(&d.target, tok)
	f.addText(&d.target, tok)
}

func (d *RecordDecl) SetFields(f *Namespace, fields stack.Stack[VarDecl]) {
}

func (d *DocDecl) SetContent(f *Namespace, strs token.TokenStack) {
	tokIndex := len(f.tokens)
	txtIndex := len(f.texts)
	size := strs.Len()
	for !strs.Empty() {
		tok, ok := strs.Pop()
		if !ok {
			panic("unreachable")
		}
		if kind := tok.Kind(); kind != token.String && kind != token.TextBlock {
			panic("doc content must either be a string or text block")
		}
		f.addToken(&d.content, tok)
		f.addText(&d.content, tok)
	}
	d.content.token.Index = tokIndex
	d.content.token.Len = size

	d.content.text.Index = txtIndex
	d.content.text.Len = size
}

func (d *TagDecl) SetAttrs(f *Namespace, attrs stack.Stack[AttrDecl]) {
}

func (d *AttrDecl) SetValue(f *Namespace, tok token.Token) {
	if kind := tok.Kind(); kind != token.String && kind != token.TextBlock {
		panic("unreachable")
	}
	f.addToken(&d.value, tok)
	f.addText(&d.value, tok)
}

func (d *TemplDecl) SetParams(f *Namespace, params stack.Stack[VarDecl]) {
}

func (d *TemplDecl) SetElements(f *Namespace, params stack.Stack[any]) {
}
