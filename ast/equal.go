package ast

func (a Def) Equal(o Def) bool {
	kind := a.Kind == o.Kind
	left := a.Name == o.Name
	return kind && left
}
