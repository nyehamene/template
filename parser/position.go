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

func (t badexpr) Pos() Position {
	return t.Position
}

func (e pkgexpr) Pos() Position {
	pos := Position{Start: e.name.Start(), End: e.name.End()}
	return pos
}

func (e importexpr) Pos() Position {
	pos := Position{Start: e.path.Start(), End: e.path.End()}
	return pos
}

func (e usingexpr) Pos() Position {
	pos := Position{Start: e.target.Start(), End: e.target.End()}
	return pos
}

func (e typeexpr) Pos() Position {
	pos := Position{Start: e.target.Start(), End: e.target.End()}
	return pos
}

func (e recordexpr) Pos() Position {
	return e.Position
}

func (e templexpr) Pos() Position {
	return e.Position
}

func (e litexpr) Pos() Position {
	t := token.Token(e)
	p := Position{Start: t.Start(), End: t.End()}
	return p
}
