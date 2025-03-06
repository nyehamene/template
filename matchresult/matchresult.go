package matchresult

import "temlang/tem/token"

func Ok(tok token.Token) Type {
	return Type{tok, tok.Kind, stateOk}
}

func NoMatch(at token.Token, exp token.Kind) Type {
	return Type{at, exp, stateNoMatch}
}

func Invalid(at token.Token, exp token.Kind) Type {
	return Type{at, exp, stateInvalid}
}

type state int

const (
	stateOk state = iota
	stateNoMatch
	stateInvalid
)

type Type struct {
	tok token.Token
	// token kind expected before an error
	exp   token.Kind
	state state
}

func (m Type) Get() token.Token {
	return m.tok
}

func (m Type) Exp() token.Kind {
	return m.exp
}

func (m Type) Ok() bool {
	return m.state == stateOk
}

func (m Type) NoMatch() bool {
	return m.state == stateNoMatch
}

func (m Type) Invalid() bool {
	return m.state != stateOk && m.state != stateNoMatch
}
