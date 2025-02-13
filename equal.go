package template

func (t Token) Equal(o Token) bool {
	offset := t.offset == o.offset
	kind := t.kind == o.kind
	return offset && kind
}

func (a Ast) Equal(o Ast) bool {
	offset := a.kind == o.kind
	left := a.left == o.left
	right := a.right == o.right
	return offset && left && right
}
