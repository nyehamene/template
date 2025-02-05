package ast

import "fmt"

func (k DefKind) String() string {
	switch k {
	case Package:
		return "package"
	case Record:
		return "record"
	case Alias:
		return "alias"
	case Template:
		return "templ"
	case Docline:
		return "docline"
	case Docblock:
		return "docblock"
	case Import:
		return "import"
	case Using:
		return "using"
	case Metatable:
		return "metatable"
	default:
		panic(fmt.Sprintf("unreachable: %#v", k))
	}
}

func (a Def) String() string {
	return fmt.Sprintf("(%s) ; %s", a.Kind, a.Name)
}
