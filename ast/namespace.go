package ast

func New(file, name string) *Namespace {
	return &Namespace{
		file: file,
		name: name,
	}
}

type Namespace struct {
	pkg   string
	name  string
	file  string
	decl  []SExpressionPrinter
	lines []int
}

func (n *Namespace) SetPackageName(name string) {
	n.pkg = name
}

func (n *Namespace) Add(d SExpressionPrinter) {
	n.decl = append(n.decl, d)
}

func (n *Namespace) SetLines(lines []int) {
	n.lines = lines
}
