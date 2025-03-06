package ast

import "temlang/tem/token"

func TokenIndexTestOnly(tok, txt int) TokenIndex {
	return TokenIndex{tok, txt}
}

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

	// ensure that n.tokens size didn't grow after appending
	// notExist token. i.e. notExist really was added to index
	// zero rather than growing the array then adding at the end
	// TODO: remove if confirmed that notExist was added to index 0
	if len(n.tokens) > int(c) {
		panic("assertion failed")
	}
}

func (n *NamespaceFile) AddToken(tok token.Token, txt string) TokenIndex {
	tokIndex := len(n.tokens)
	txtIndex := len(n.texts)

	n.tokens = append(n.tokens, tok)
	n.texts = append(n.texts, txt)
	return TokenIndex{tokIndex, txtIndex}
}

func (n NamespaceFile) GetToken(t TokenIndex) token.Token {
	tok := n.tokens[t.token]
	return tok
}

func (n NamespaceFile) GetName(t TokenIndex) string {
	txt := n.texts[t.text]
	return txt
}

type TokenIndex struct {
	token int
	text  int
}

// type TokenSlice struct {
// 	offset Index
// 	len    int
// }

type Decl struct {
	Ident TokenIndex
	Type  TokenIndex
	// LeadingDocs  Index
	// TrailingDocs Index
	// Tags         Index
}

type PackageDecl struct {
	Decl
	Name  TokenIndex
	Templ TokenIndex
}

type ImportDecl struct {
	Decl
	Path TokenIndex
	// offset Index
}

// type UsingDecl struct {
// 	Decl
// 	Import Index
//	offset Index
// }

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
