package main

import (
	"testing"
)

func TestArtifactFilename(t *testing.T) {
	tests := map[string]struct {
		filename  string
		chartDir  string
		targetEnv string
		want      string
	}{
		"default chart dir": {
			filename:  "diff",
			chartDir:  "./chart",
			targetEnv: "foo-dev",
			want:      "diff-foo-dev",
		},
		"default chart dir without prefix": {
			filename:  "diff",
			chartDir:  "chart",
			targetEnv: "dev",
			want:      "diff-dev",
		},
		"other chart dir": {
			filename:  "diff",
			chartDir:  "./foo-chart",
			targetEnv: "qa",
			want:      "foo-chart-diff-qa",
		},
		"other chart dir without prefix": {
			filename:  "diff",
			chartDir:  "bar-chart",
			targetEnv: "foo-qa",
			want:      "bar-chart-diff-foo-qa",
		},
		"nested chart dir": {
			filename:  "diff",
			chartDir:  "./some/path/chart",
			targetEnv: "prod",
			want:      "some-path-chart-diff-prod",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := artifactFilename(tc.filename, tc.chartDir, tc.targetEnv)
			if got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
