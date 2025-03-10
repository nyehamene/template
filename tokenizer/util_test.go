package tokenizer_test

import "temlang/tem/token"

func ident(offset, end int) token.Token {
	return newToken(token.Ident, offset, end)
}

func import0(offset, end int) token.Token {
	return newToken(token.Import, offset, end)
}

func package0(offset, end int) token.Token {
	return newToken(token.Package, offset, end)
}

func record(offset, end int) token.Token {
	return newToken(token.Record, offset, end)
}

func templ(offset, end int) token.Token {
	return newToken(token.Templ, offset, end)
}

func type0(offset, end int) token.Token {
	return newToken(token.Type, offset, end)
}

func using(offset, end int) token.Token {
	return newToken(token.Using, offset, end)
}

func str(offset, end int) token.Token {
	return newToken(token.String, offset, end)
}

func textBlock(offset, end int) token.Token {
	return newToken(token.TextBlock, offset, end)
}

func comment(offset, end int) token.Token {
	return newToken(token.Comment, offset, end)
}

// eol automatically inserted semicolon
func eol(offset int) token.Token {
	return newToken(token.Semicolon, offset, offset)
}

func symbol(kind token.Kind, offset int) token.Token {
	return newToken(kind, offset, offset+1)
}

func newToken(kind token.Kind, offset, end int) token.Token {
	return token.New(kind, offset, end)
}
