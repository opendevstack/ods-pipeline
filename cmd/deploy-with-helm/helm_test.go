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

func TestAssembleHelmDiffArgs(t *testing.T) {
	tests := map[string]struct {
		releaseNamespace string
		releaseName      string
		helmArchive      string
		opts             options
		valuesFiles      []string
		cliValues        []string
		want             []string
	}{
		"default": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--three-way-merge", "--install",
				"b", "c"},
		},
		"with no diff flags": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "", upgradeFlags: "--install"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--install",
				"b", "c"},
		},
		"with values file": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install"},
			valuesFiles:      []string{"values.dev.yaml"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--three-way-merge", "--install", "--values=values.dev.yaml",
				"b", "c"},
		},
		"with CLI values": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install"},
			cliValues:        []string{"--set=image.tag=abcdef"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--three-way-merge", "--install", "--set=image.tag=abcdef",
				"b", "c"},
		},
		"with multiple args": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts: options{
				diffFlags:    "--three-way-merge --no-hooks --include-tests",
				upgradeFlags: "--install --wait",
			},
			valuesFiles: []string{"secrets.yaml", "values.dev.yaml", "secrets.dev.yaml"},
			cliValues:   []string{"--set=image.tag=abcdef", "--set=x=y"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color",
				"--three-way-merge", "--no-hooks", "--include-tests",
				"--install", "--wait",
				"--values=secrets.yaml", "--values=values.dev.yaml", "--values=secrets.dev.yaml",
				"--set=image.tag=abcdef", "--set=x=y",
				"b", "c"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := assembleHelmDiffArgs(
				tc.releaseNamespace, tc.releaseName, tc.helmArchive,
				tc.opts,
				tc.valuesFiles, tc.cliValues)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("args mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAssembleHelmUpgradeArgs(t *testing.T) {
	tests := map[string]struct {
		releaseNamespace string
		releaseName      string
		helmArchive      string
		opts             options
		valuesFiles      []string
		cliValues        []string
		want             []string
	}{
		"default": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install --wait"},
			want: []string{"--namespace=a", "secrets", "upgrade",
				"--install", "--wait",
				"b", "c"},
		},
		"with no upgrade flags": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: ""},
			want: []string{"--namespace=a", "secrets", "upgrade",

				"b", "c"},
		},
		"with values file": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install --wait"},
			valuesFiles:      []string{"values.dev.yaml"},
			want: []string{"--namespace=a", "secrets", "upgrade",
				"--install", "--wait",
				"--values=values.dev.yaml",
				"b", "c"},
		},
		"with CLI values": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install --wait"},
			cliValues:        []string{"--set=image.tag=abcdef"},
			want: []string{"--namespace=a", "secrets", "upgrade",
				"--install", "--wait",
				"--set=image.tag=abcdef",
				"b", "c"},
		},
		"with multiple args": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install --atomic"},
			valuesFiles:      []string{"secrets.yaml", "values.dev.yaml", "secrets.dev.yaml"},
			cliValues:        []string{"--set=image.tag=abcdef", "--set=x=y"},
			want: []string{"--namespace=a", "secrets", "upgrade",
				"--install", "--atomic",
				"--values=secrets.yaml", "--values=values.dev.yaml", "--values=secrets.dev.yaml",
				"--set=image.tag=abcdef", "--set=x=y",
				"b", "c"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := assembleHelmUpgradeArgs(
				tc.releaseNamespace, tc.releaseName, tc.helmArchive,
				tc.opts,
				tc.valuesFiles, tc.cliValues)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("args mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
