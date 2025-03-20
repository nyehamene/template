package token

import (
	"fmt"
)

type Location struct {
	Start, End Pos
	File       string
}

func (l Location) String() string {
	return fmt.Sprintf("%s - %s ; %s]", l.Start, l.End, l.File)
}

type Pos struct {
	Line, Col int
}

func (p Pos) String() string {
	return fmt.Sprintf("[%d-%d]", p.Line, p.Col)
}
