package parser

import (
	"temlang/tem/ast"
	"temlang/tem/tokenizer"
)

func NewParser(t *tokenizer.Tokenizer) Parser {
	return Parser{t, []ast.Def{}}
}

type Parser struct {
	tokenizer *tokenizer.Tokenizer
	ast       []ast.Def
}
