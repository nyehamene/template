package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	path := "./testdata"
	list, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range list {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".tem") {
			continue
		}

		filename := filepath.Join(path, d.Name())
		src, err := os.ReadFile(filename)
		if err != nil {
			t.Error(err)
			continue
		}

		_, errs := ParseFile(filename, src)

		if errs.Len() != 0 {
			t.Errorf("ParseFile(%s) parser failed unexpectedly", filename)
		}
	}
}
