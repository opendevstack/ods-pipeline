package main

import (
	"fmt"
	"testing"

	"github.com/opendevstack/pipeline/pkg/artifact"
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

func TestGetImageURLs(t *testing.T) {
	srcHost := "image-registry.openshift-image-registry.svc:5000"
	destHost := "default-route-openshift-image-registry.apps.example.com"
	imgArtifact := artifact.Image{
		Ref:        fmt.Sprintf("%s/foo/bar:baz", srcHost),
		Repository: "foo", Name: "bar", Tag: "baz",
	}
	tests := map[string]struct {
		registryHost     string
		releaseNamespace string
		want             string
	}{
		"same cluster, same namespace": {
			registryHost:     "",
			releaseNamespace: "foo",
			want:             fmt.Sprintf("%s/foo/bar:baz", srcHost),
		},
		"same cluster, different namespace": {
			registryHost:     "",
			releaseNamespace: "other",
			want:             fmt.Sprintf("%s/other/bar:baz", srcHost),
		},
		"different cluster, same namespace": {
			registryHost:     destHost,
			releaseNamespace: "foo",
			want:             fmt.Sprintf("%s/foo/bar:baz", destHost),
		},
		"different cluster, different namespace": {
			registryHost:     destHost,
			releaseNamespace: "other",
			want:             fmt.Sprintf("%s/other/bar:baz", destHost),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := getImageDestURL(tc.registryHost, tc.releaseNamespace, imgArtifact)
			if got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
