package parser

import (
	"temlang/tem/matchresult"
	"temlang/tem/token"
)

type MatchToken = matchresult.Type[token.Token, token.Kind]

type MatchManyToken = matchresult.Type[[]token.Token, token.Kind]

func OkMany(toks []token.Token) MatchManyToken {
	var empty token.Kind
	res := matchresult.Ok(toks, empty)
	return MatchManyToken(res)
}

func Ok(tok token.Token) MatchToken {
	return matchresult.Ok(tok, tok.Kind())
}

func NoMatch(tok token.Token, k token.Kind) MatchToken {
	return matchresult.NoMatch(tok, k)
}

func Invalid(tok token.Token, k token.Kind) MatchToken {
	return matchresult.Invalid(tok, k)
}

func NoMatchMany(tok token.Token, k token.Kind) MatchManyToken {
	return matchresult.NoMatch([]token.Token{tok}, k)
}

func InvalidMany(tok token.Token, k token.Kind) MatchManyToken {
	return matchresult.Invalid([]token.Token{tok}, k)
}
