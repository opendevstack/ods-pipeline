package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanHelmDiffOutput(t *testing.T) {
	tests := map[string]struct {
		example string
		want    string
	}{
		"diff detected drift": {
			example: `Error: identified at least one change, exiting with non-zero exit code (detailed-exitcode parameter enabled)
Error: plugin "diff" exited with error

[helm-secrets] Removed: ./chart/secrets.dev.yaml.dec
Error: plugin "secrets" exited with error`,
			want: `plugin "diff" identified at least one change

[helm-secrets] Removed: ./chart/secrets.dev.yaml.dec
`,
		},
		"diff detected drift with debug turned on": {
			example: `Error: identified at least one change, exiting with non-zero exit code (detailed-exitcode parameter enabled)
Error: plugin "diff" exited with error
helm.go:81: [debug] plugin "diff" exited with error

[helm-secrets] Removed: ./chart/secrets.dev.yaml.dec
Error: plugin "secrets" exited with error
helm.go:81: [debug] plugin "secrets" exited with error`,
			want: `plugin "diff" identified at least one change

[helm-secrets] Removed: ./chart/secrets.dev.yaml.dec
`,
		},
		"diff encounters an error": {
			example: `Error: This command needs 2 arguments: release name, chart path

Use "diff [command] --help" for more information about a command.

Error: plugin "diff" exited with error`,
			want: `Error: This command needs 2 arguments: release name, chart path

Use "diff [command] --help" for more information about a command.

Error: plugin "diff" exited with error`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := cleanHelmDiffOutput([]byte(tc.example))
			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
