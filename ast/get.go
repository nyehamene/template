package ast

func (n Namespace) Imports() []ImportDecl {
	return copy0(n.imports)
}

func (n Namespace) Usings() []UsingDecl {
	return copy0(n.usings)
}

func (n Namespace) Types() []TypeDecl {
	return copy0(n.types)
}

func (n Namespace) Records() []RecordDecl {
	return copy0(n.records)
}

func (n Namespace) Docs() []DocDecl {
	return copy0(n.docs)
}

func (n Namespace) Tags() []TagDecl {
	return copy0(n.tags)
}

func (n Namespace) Templs() []TemplDecl {
	return copy0(n.templs)
}

func copy0[T any](src []T) []T {
	cpy := make([]T, 0, len(src))
	copy(cpy, src)
	return cpy
}
