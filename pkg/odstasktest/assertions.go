package odstasktest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// AssertFilesExist checks that all files named by wantFiles exist in wsDir.
// Any files that do not exist will report a test error.
func AssertFilesExist(t *testing.T, wsDir string, wantFiles ...string) {
	t.Helper()
	for _, wf := range wantFiles {
		filename := filepath.Join(wsDir, wf)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Want %s, but got nothing", filename)
		}
	}
}

// AssertFileContent checks that the file named by filename in the directory
// wsDir has the exact context specified by want.
func AssertFileContent(t *testing.T, wsDir, filename, want string) {
	t.Helper()
	got, err := getTrimmedFileContent(filepath.Join(wsDir, filename))
	if err != nil {
		t.Errorf("get content of %s: %s", filename, err)
		return
	}
	if got != want {
		t.Errorf("got '%s', want '%s' in file %s", got, want, filename)
	}
}

// AssertFileContentContains checks that the file named by filename in the directory
// wsDir contains all of wantContains.
func AssertFileContentContains(t *testing.T, wsDir, filename string, wantContains ...string) {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(wsDir, filename))
	got := string(content)
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	for _, w := range wantContains {
		if !strings.Contains(got, w) {
			t.Fatalf("got '%s', want '%s' contained in file %s", got, w, filename)
		}
	}
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
