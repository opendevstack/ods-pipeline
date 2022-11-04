package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/pkg/config"
)

func TestHelmDiff(t *testing.T) {
	tests := map[string]struct {
		cmdExitCode int
		wantInSync  bool
		wantErr     bool
	}{
		"diff exits with generic exit code": {
			cmdExitCode: diffGenericExitCode,
			wantInSync:  false,
			wantErr:     true,
		},
		"diff exits with drift exit code": {
			cmdExitCode: diffDriftExitCode,
			wantInSync:  false,
			wantErr:     false,
		},
		"diff passes (no drift)": {
			cmdExitCode: 0,
			wantInSync:  true,
			wantErr:     false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			d := &deployHelm{helmBin: "../../test/scripts/exit-with-code.sh"}
			driftDetected, err := d.helmDiff(
				[]string{"", "", strconv.Itoa(tc.cmdExitCode)},
				&stdout, &stderr,
			)
			if tc.wantErr && err == nil {
				t.Fatal("want err, got none")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("want no err, got %s", err)
			}
			if tc.wantInSync != driftDetected {
				t.Fatalf("want success=%v, got success=%v", tc.wantInSync, driftDetected)
			}
		})
	}
}

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
		"diff encounters another error": {
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
			got := cleanHelmDiffOutput(tc.example)
			if diff := cmp.Diff(tc.want, got); diff != "" {
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
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install", debug: true},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--normalize-manifests", "--three-way-merge", "--debug", "--install",
				"b", "c"},
		},
		"with no diff flags": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "", upgradeFlags: "--install"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--normalize-manifests", "--install",
				"b", "c"},
		},
		"with values file": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install"},
			valuesFiles:      []string{"values.dev.yaml"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--normalize-manifests", "--three-way-merge", "--install", "--values=values.dev.yaml",
				"b", "c"},
		},
		"with CLI values": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install"},
			cliValues:        []string{"--set=image.tag=abcdef"},
			want: []string{"--namespace=a", "secrets", "diff", "upgrade",
				"--detailed-exitcode", "--no-color", "--normalize-manifests", "--three-way-merge", "--install", "--set=image.tag=abcdef",
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
				"--detailed-exitcode", "--no-color", "--normalize-manifests",
				"--three-way-merge", "--no-hooks", "--include-tests",
				"--install", "--wait",
				"--values=secrets.yaml", "--values=values.dev.yaml", "--values=secrets.dev.yaml",
				"--set=image.tag=abcdef", "--set=x=y",
				"b", "c"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := &deployHelm{
				releaseNamespace: tc.releaseNamespace,
				releaseName:      tc.releaseName,
				helmArchive:      tc.helmArchive,
				opts:             tc.opts,
				valuesFiles:      tc.valuesFiles,
				cliValues:        tc.cliValues,
				targetConfig:     &config.Environment{},
			}
			got, err := d.assembleHelmDiffArgs()
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
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install --wait", debug: true},
			want: []string{"--namespace=a", "secrets", "upgrade",
				"--kube-apiserver=https://example.com", "--kube-token=s3cr3t",
				"--debug",
				"--install", "--wait",
				"b", "c"},
		},
		"with no upgrade flags": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: ""},
			want: []string{"--namespace=a", "secrets", "upgrade",
				"--kube-apiserver=https://example.com", "--kube-token=s3cr3t",
				"b", "c"},
		},
		"with values file": {
			releaseNamespace: "a",
			releaseName:      "b",
			helmArchive:      "c",
			opts:             options{diffFlags: "--three-way-merge", upgradeFlags: "--install --wait"},
			valuesFiles:      []string{"values.dev.yaml"},
			want: []string{"--namespace=a", "secrets", "upgrade",
				"--kube-apiserver=https://example.com", "--kube-token=s3cr3t",
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
				"--kube-apiserver=https://example.com", "--kube-token=s3cr3t",
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
				"--kube-apiserver=https://example.com", "--kube-token=s3cr3t",
				"--install", "--atomic",
				"--values=secrets.yaml", "--values=values.dev.yaml", "--values=secrets.dev.yaml",
				"--set=image.tag=abcdef", "--set=x=y",
				"b", "c"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := &deployHelm{
				releaseNamespace: tc.releaseNamespace,
				releaseName:      tc.releaseName,
				helmArchive:      tc.helmArchive,
				opts:             tc.opts,
				valuesFiles:      tc.valuesFiles,
				cliValues:        tc.cliValues,
				targetConfig:     &config.Environment{APIServer: "https://example.com", APIToken: "s3cr3t"},
			}
			got, err := d.assembleHelmUpgradeArgs()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("args mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPrintlnSafeHelmCmd(t *testing.T) {
	var stdout bytes.Buffer
	printlnSafeHelmCmd([]string{"diff", "upgrade", "--kube-apiserver=https://example.com", "--kube-token=s3cr3t", "--debug"}, &stdout)
	want := "helm diff upgrade --kube-apiserver=https://example.com --kube-token=*** --debug"
	got := strings.TrimSpace(stdout.String())
	if got != want {
		t.Fatalf("want: '%s', got: '%s'", want, got)
	}
}
