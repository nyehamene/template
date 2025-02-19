package template

func (p Parser) defRecord(start int) (DefKind, int, bool) {
	var err error
	next := start
	if _, next, err = p.expect(next, TokenRecord); err != nil {
		return DefRecord, start, false
	}
	if _, next, err = p.expect(next, TokenBraceLeft); err != nil {
		return DefRecord, start, false
	}

	for {
		var n int
		var e error
		if _, n, e = p.expect(next, TokenIdent); e != nil {
			break
		}
		if _, next, err = p.expect(n, TokenColon); err != nil {
			return DefRecord, start, false
		}
		if _, next, err = p.expect(next, TokenIdent); err != nil {
			return DefRecord, start, false
		}
		if _, next, err = p.expect(next, TokenSemicolon); err != nil {
			return DefRecord, start, false
		}
	}
	if _, next, err = p.expect(next, TokenBraceRight); err != nil {
		return DefRecord, start, false
	}
	return DefRecord, next, true
}
