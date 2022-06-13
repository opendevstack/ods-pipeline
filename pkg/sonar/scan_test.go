package sonar

import (
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/projectpath"
)

func TestExtractComputeEngineTaskID(t *testing.T) {

	c := testClient(t, "")
	want := "AVAn5RKqYwETbXvgas-I"
	fixture := filepath.Join(projectpath.Root, "test/testdata/fixtures/sonar", ReportTaskFilename)
	got, err := c.ExtractComputeEngineTaskID(fixture)
	if err != nil {
		t.Fatal(err)
	}

	// check extracted status matches
	if got != want {
		t.Fatalf("want %s, got %s", want, got)
	}
}
