package token

func (t Token) Equal(o Token) bool {
	return t.Kind() == o.Kind() &&
		t.start == o.start &&
		t.end == o.end
}
