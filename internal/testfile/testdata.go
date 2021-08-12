package testfile

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/projectpath"
)

// ReadFixture returns the contents of the fixture file or fails.
func ReadFixture(t *testing.T, filename string) []byte {
	return readFileOrFatal(t, filepath.Join(projectpath.Root, "test/testdata/fixtures", filename))
}

// ReadGolden returns the contents of the golden file or fails.
func ReadGolden(t *testing.T, filename string) []byte {
	return readFileOrFatal(t, filepath.Join(projectpath.Root, "test/testdata/golden", filename))
}

func readFileOrFatal(t *testing.T, filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
