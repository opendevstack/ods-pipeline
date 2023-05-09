package config

import (
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/projectpath"
)

func TestReadFromDir(t *testing.T) {
	ods, err := ReadFromDir(filepath.Join(projectpath.Root, "test/testdata/fixtures/config"))
	if err != nil {
		t.Fatal(err)
	}
	want := "pr:opened"
	got := ods.Pipelines[0].Triggers[0].Events[0]
	if got != want {
		t.Fatalf("Got %s, want %s", got, want)
	}
}

func TestReadFromFile(t *testing.T) {
	ods, err := ReadFromFile(filepath.Join(projectpath.Root, "test/testdata/fixtures/config/ods.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	want := "pr:opened"
	got := ods.Pipelines[0].Triggers[0].Events[0]
	if got != want {
		t.Fatalf("Got %s, want %s", got, want)
	}
}

func TestRead(t *testing.T) {
	tests := map[string]struct {
		Fixture   []byte
		WantError string
	}{
		"broken YAML": {
			Fixture: []byte(`pipelines:
			- triggers: []`),
			WantError: "error converting YAML to JSON: yaml: line 2: found character that cannot start any token",
		},
		"extra field": {
			Fixture: []byte(`foo: bar
pipelines:
- triggers: []`),
			WantError: `error unmarshaling JSON: while decoding JSON: json: unknown field "foo"`,
		},
		"duplicate field": {
			Fixture: []byte(`pipelines: []
pipelines:
- triggers: []`),
			WantError: `error converting YAML to JSON: yaml: unmarshal errors:
  line 3: key "pipelines" already set in map`,
		},
		"valid": {
			Fixture: []byte(`pipelines:
- triggers: []`),
			WantError: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Read(tc.Fixture)
			if len(tc.WantError) == 0 && err != nil {
				t.Fatal(err)
			} else if len(tc.WantError) > 0 {
				if err == nil || tc.WantError != err.Error() {
					t.Fatalf("Want error: %s, got: %s", tc.WantError, err)
				}
			}
		})
	}
}
