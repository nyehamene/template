package token

func (t Token) Equal(o Token) bool {
	offsetOk := t.Offset == o.Offset
	kindOk := t.Kind == o.Kind
	return offsetOk && kindOk
}
