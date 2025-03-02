package tokenizer_test

import (
	"fmt"
	"temlang/tem/tokenizer"
	"testing"
)

func HelperRunTestCasesError(t *testing.T, testcases map[string]int) {
	i := 0
	for src, expected := range testcases {
		t.Run(fmt.Sprintf("%d_%s", i, src), func(t *testing.T) {
			tok := tokenizer.New([]byte(src))
			_ = tok.Next()

			if ec := tok.ErrorCount(); ec != expected {
				t.Errorf("expected %d error(s) got %d", expected, ec)
			}
		})
		i += 1
	}
}

func TestNextTextBlockError(t *testing.T) {
	testcases := map[string]int{
		`"""`:        1,
		`""" line 1`: 1,
	}
	HelperRunTestCasesError(t, testcases)
}

func TestNextStringError(t *testing.T) {
	textcases := map[string]int{
		`"`:    1,
		`"abc`: 1,
	}
	HelperRunTestCasesError(t, textcases)
}
