package tokenizer_test

import "temlang/tem/token"

func NewIdent(offset, end int) token.Token {
	return NewToken(token.Ident, offset, end)
}

func NewImport(offset, end int) token.Token {
	return NewToken(token.Import, offset, end)
}

func NewPackage(offset, end int) token.Token {
	return NewToken(token.Package, offset, end)
}

func NewRecord(offset, end int) token.Token {
	return NewToken(token.Record, offset, end)
}

func NewTempl(offset, end int) token.Token {
	return NewToken(token.Templ, offset, end)
}

func NewType(offset, end int) token.Token {
	return NewToken(token.Type, offset, end)
}

func NewUsing(offset, end int) token.Token {
	return NewToken(token.Using, offset, end)
}

func NewStr(offset, end int) token.Token {
	return NewToken(token.String, offset, end)
}

func NewTextBlock(offset, end int) token.Token {
	return NewToken(token.TextBlock, offset, end)
}

func NewComment(offset, end int) token.Token {
	return NewToken(token.Comment, offset, end)
}

// eol automatically inserted semicolon
func NewEOL(offset int) token.Token {
	return NewToken(token.Semicolon, offset, offset)
}

func NewSymbol(kind token.Kind, offset int) token.Token {
	return NewToken(kind, offset, offset+1)
}

func NewToken(kind token.Kind, offset, end int) token.Token {
	return token.New(kind, offset, end)
}
