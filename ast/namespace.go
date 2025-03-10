package ast

func (n *NamespaceFile) AddImport(d ImportDecl) {
	n.imports = append(n.imports, d)
}

func (n *NamespaceFile) AddUsing(d UsingDecl) {
	n.usings = append(n.usings, d)
}
