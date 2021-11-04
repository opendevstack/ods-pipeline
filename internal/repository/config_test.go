package repository

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/pkg/config"
)

type fakeRawClient struct {
	// files contains byte slices for filenames
	files map[string][]byte
}

func (c *fakeRawClient) RawGet(project, repository, filename, gitFullRef string) ([]byte, error) {
	if f, ok := c.files[filename]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("%s not found", filename)
}

func TestGetODSConfig(t *testing.T) {
	bitbucketClient := &fakeRawClient{}

	tests := map[string]struct {
		files   map[string][]byte
		wantErr string
		wantODS *config.ODS
	}{
		"no ODS file": {
			files:   map[string][]byte{},
			wantErr: "ods.yml not found",
		},
		"empty ODS file": {
			files: map[string][]byte{
				"ods.yaml": []byte(""),
			},
			wantErr: "config is empty",
		},
		"ods.yaml file": {
			files: map[string][]byte{
				"ods.yaml": []byte("environments: []"),
			},
			wantErr: "",
			wantODS: &config.ODS{Environments: []config.Environment{}},
		},
		"ods.yml file": {
			files: map[string][]byte{
				"ods.yml": []byte("environments: []"),
			},
			wantErr: "",
			wantODS: &config.ODS{Environments: []config.Environment{}},
		},
		"ods.yaml has precedence over ods.yml file": {
			files: map[string][]byte{
				"ods.yaml": []byte("version: 1.0.0"),
				"ods.yml":  []byte("version: 0.1.0"),
			},
			wantErr: "",
			wantODS: &config.ODS{Version: "1.0.0"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bitbucketClient.files = tc.files
			// As context is read from fake Bitbucket, dir value is unused.
			got, err := GetODSConfig(bitbucketClient, "foo", "bar", "refs/heads/master")
			if tc.wantErr == "" && err != nil {
				t.Fatal(err)
			} else if tc.wantErr != "" && err == nil {
				t.Fatalf("want err: %s, got nothing", tc.wantErr)
			} else if err != nil && tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("want err: %s, got err: %s", tc.wantErr, err)
			}
			if diff := cmp.Diff(tc.wantODS, got); diff != "" {
				t.Fatalf("context mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
