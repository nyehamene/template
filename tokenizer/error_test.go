package tokenizer

import (
	"fmt"
	"testing"
)

type TestCaseError map[string]struct {
	errCount    int
	offsetAfter int
}

func HelperRunTestCasesError(t *testing.T, testcases TestCaseError) {
	i := 0
	for src, expected := range testcases {
		t.Run(fmt.Sprintf("%d_%s", i, src), func(t *testing.T) {
			tok := New([]byte(src))
			_ = tok.Next()

			if errCount := tok.ErrorCount(); errCount != expected.errCount {
				t.Errorf("expected %d error(s) got %d", expected.errCount, errCount)
			}

			if expected.offsetAfter != tok.offset {
				t.Errorf("expected offset %d got %d", expected.offsetAfter, tok.offset)
			}
		})
		i += 1
	}
}

func TestNextStringError(t *testing.T) {
	textcases := TestCaseError{
		`"`:    {errCount: 1, offsetAfter: 1},
		`"abc`: {errCount: 1, offsetAfter: 4},
	}
	HelperRunTestCasesError(t, textcases)
}

func TestNextCommentError(t *testing.T) {
	testcases := TestCaseError{
		`/ invalid comment`: {errCount: 1, offsetAfter: 2},
	}
	HelperRunTestCasesError(t, testcases)
}
