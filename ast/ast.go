package ast

import "temlang/tem/token"

type Def struct {
	Kind DefKind
	Name token.Token
}

type DefKind int

const (
	Package DefKind = iota
	Record
	Alias
	Template
	Docline
	Docblock
	Import
	Using
	Metatable
)
