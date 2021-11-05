package repository

import (
	"testing"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

type fakeBitbucketClient struct {
	branches []bitbucket.Branch
}

func (c *fakeBitbucketClient) BranchList(projectKey string, repositorySlug string, params bitbucket.BranchListParams) (*bitbucket.BranchPage, error) {
	return &bitbucket.BranchPage{
		Values: c.branches,
	}, nil
}

func TestBestMatchingBranch(t *testing.T) {
	bitbucketClient := &fakeBitbucketClient{}

	tests := map[string]struct {
		branches []bitbucket.Branch
		subrepo  config.Repository
		version  string
		want     string
	}{
		"no configured branch and WIP": {
			version: pipelinectxt.WIP,
			want:    config.DefaultBranch,
		},
		"configured branch and WIP": {
			subrepo: config.Repository{Branch: "production"},
			version: pipelinectxt.WIP,
			want:    "refs/heads/production",
		},
		"no configured branch, no/non-matching release branch, version": {
			branches: []bitbucket.Branch{
				{DisplayID: "release/0.1.0", ID: "refs/heads/release/0.1.0"},
			},
			version: "1.0.0",
			want:    config.DefaultBranch,
		},
		"no configured branch, matching release branch, version": {
			branches: []bitbucket.Branch{
				{DisplayID: "release/1.0.0", ID: "refs/heads/release/1.0.0"},
				{DisplayID: "release/0.1.0", ID: "refs/heads/release/0.1.0"},
			},
			version: "1.0.0",
			want:    "refs/heads/release/1.0.0",
		},
		"configured branch, no/non-matching release branch, version": {
			branches: []bitbucket.Branch{
				{DisplayID: "release/0.1.0", ID: "refs/heads/release/0.1.0"},
			},
			subrepo: config.Repository{Branch: "production"},
			version: "1.0.0",
			want:    "refs/heads/production",
		},
		"configured branch, matching release branch, version": {
			branches: []bitbucket.Branch{
				{DisplayID: "release/1.0.0", ID: "refs/heads/release/1.0.0"},
				{DisplayID: "release/0.1.0", ID: "refs/heads/release/0.1.0"},
			},
			subrepo: config.Repository{Branch: "production"},
			version: "1.0.0",
			want:    "refs/heads/release/1.0.0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bitbucketClient.branches = tc.branches
			got, err := BestMatchingBranch(bitbucketClient, "foo", tc.subrepo, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
