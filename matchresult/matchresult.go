package matchresult

func Ok[T any, E any](t T, e E) Type[T, E] {
	return Type[T, E]{t, e, stateOk}
}

func NoMatch[T any, E any](at T, exp E) Type[T, E] {
	return Type[T, E]{at, exp, stateNoMatch}
}

func Invalid[T any, E any](at T, exp E) Type[T, E] {
	return Type[T, E]{at, exp, stateInvalid}
}

type state int

const (
	stateOk state = iota
	stateNoMatch
	stateInvalid
)

type Type[T any, E any] struct {
	value T
	// token kind expected before an error
	exp   E
	state state
}

func (m Type[T, E]) Get() T {
	return m.value
}

func (m Type[T, E]) Exp() E {
	return m.exp
}

func (m Type[T, E]) Ok() bool {
	return m.state == stateOk
}

func (m Type[T, E]) NoMatch() bool {
	return m.state == stateNoMatch
}

func (m Type[T, E]) Invalid() bool {
	return m.state != stateOk && m.state != stateNoMatch
}
