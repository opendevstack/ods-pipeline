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
	gotStage := ods.Environments[0].Stage
	if gotStage != ProdStage {
		t.Fatalf("Got %s, want prod", gotStage)
	}
}

func TestReadFromFile(t *testing.T) {
	ods, err := ReadFromFile(filepath.Join(projectpath.Root, "test/testdata/fixtures/config/ods.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	gotStage := ods.Environments[0].Stage
	if gotStage != ProdStage {
		t.Fatalf("Got %s, want prod", gotStage)
	}
}

func TestReadFromLegacyFormatFile(t *testing.T) {
	ods, err := ReadFromFile(filepath.Join(projectpath.Root, "test/testdata/fixtures/config/ods-legacy.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	gotStage := ods.Environments[0].Stage
	if gotStage != ProdStage {
		t.Fatalf("Got %s, want prod", gotStage)
	}
	gotNumPipelines := len(ods.Pipeline)
	if gotNumPipelines != 1 {
		t.Fatalf("Got %d pipeline definitions, want 1", gotNumPipelines)
	}
}

func TestRead(t *testing.T) {
	tests := map[string]struct {
		Fixture   []byte
		WantError string
	}{
		"broken YAML": {
			Fixture: []byte(`environments:
			- name: foo
			  stage: dev`),
			WantError: "could not unmarshal config: error converting YAML to JSON: yaml: line 2: found character that cannot start any token",
		},
		"extra field": {
			Fixture: []byte(`foo: bar
environments:
- name: foo
  stage: dev`),
			WantError: `could not unmarshal config: error unmarshaling JSON: while decoding JSON: json: unknown field "foo"`,
		},
		"duplicate field": {
			Fixture: []byte(`environments: []
environments:
- name: foo
  stage: dev`),
			WantError: `could not unmarshal config: error converting YAML to JSON: yaml: unmarshal errors:
  line 3: key "environments" already set in map`,
		},
		"missing stage": {
			Fixture: []byte(`environments:
- name: foo`),
			WantError: "invalid stage value '' for environment foo",
		},
		"blank stage": {
			Fixture: []byte(`environments:
- name: foo
  stage: ""`),
			WantError: "invalid stage value '' for environment foo",
		},
		"wrong stage": {
			Fixture: []byte(`environments:
- name: foo
  stage: staging`),
			WantError: "invalid stage value 'staging' for environment foo",
		},
		"blank name": {
			Fixture: []byte(`environments:
- name: ""
  stage: qa`),
			WantError: "name of environment must not be blank",
		},
		"invalid name - extra characters": {
			Fixture: []byte(`environments:
- name: "Hello World!"
  stage: dev`),
			WantError: "name of environment must match ^[a-z-]*$",
		},
		"invalid name - uppercase": {
			Fixture: []byte(`environments:
- name: "DEVenv"
  stage: dev`),
			WantError: "name of environment must match ^[a-z-]*$",
		},
		"valid": {
			Fixture: []byte(`environments:
- name: foo-qa
  stage: qa`),
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

func TestEnvironment(t *testing.T) {
	o := &ODS{
		Environments: []Environment{
			{
				Name:  "a",
				Stage: "dev",
			},
			{
				Name:  "b",
				Stage: "dev",
			},
		},
	}
	got, err := o.Environment("b")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "b" {
		t.Fatalf("Want env: b, got: %s", got.Name)
	}
}
