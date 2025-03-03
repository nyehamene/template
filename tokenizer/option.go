package tokenizer

type Option func(t *Tokenizer)

func SetSemicolonHandler(h SemicolonHandler) Option {
	return func(t *Tokenizer) {
		t.semicolonFunc = h
	}
}

func SetErrorHandler(h ErrorHandler) Option {
	return func(t *Tokenizer) {
		t.errFunc = h
	}
}
