package parser

import "temlang/tem/token"

type Position struct {
	Start int
	End   int
}

func (t badtree) Pos() Position {
	return t.Position
}

func (t decltree) Pos() Position {
	return t.Position
}

func (t tagtree) Pos() Position {
	return t.Position
}

func (t attrtree) Pos() Position {
	return t.Position
}

func (t doctree) Pos() Position {
	return t.Position
}

func (t vartree) Pos() Position {
	return t.Position
}

func (e baseexpr) Pos() Position {
	p := Position{Start: e.start, End: e.end}
	return p
}

func (e litexpr) Pos() Position {
	t := token.Token(e)
	p := Position{Start: t.Start(), End: t.End()}
	return p
}
