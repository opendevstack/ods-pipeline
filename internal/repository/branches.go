package repository

import (
	"fmt"
	"strings"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

// BestMatchingBranch returns the best matching branch for given subrepository
// and version.
// The best match is, in order of specifity:
// - release branch (if there is a non-WIP version and the branch exists)
// - configured branch (if configured)
// - global default branch
func BestMatchingBranch(bitbucketClient bitbucket.BranchClientInterface, project string, subrepo config.Repository, version string) (string, error) {
	subrepoGitFullRef := config.DefaultBranch
	if len(subrepo.Branch) > 0 {
		subrepoGitFullRef = subrepo.Branch
		if !strings.HasPrefix(subrepoGitFullRef, "refs/") {
			subrepoGitFullRef = fmt.Sprintf("refs/heads/%s", subrepoGitFullRef)
		}
	}
	if version != pipelinectxt.WIP {
		releaseBranch, err := findReleaseBranch(bitbucketClient, project, subrepo.Name, version)
		if err != nil {
			return "", fmt.Errorf("could not detect release branches: %w", err)
		}
		if releaseBranch != "" {
			subrepoGitFullRef = releaseBranch
		}
	}
	return subrepoGitFullRef, nil
}

// findReleaseBranch returns the full Git ref of the release branch corresponding
// to given version. If none is found, it returns an empty string.
func findReleaseBranch(bitbucketClient bitbucket.BranchClientInterface, projectKey, repositorySlug, version string) (string, error) {
	releaseBranch := fmt.Sprintf("release/%s", version)
	branchPage, err := bitbucketClient.BranchList(projectKey, repositorySlug, bitbucket.BranchListParams{
		FilterText:   fmt.Sprintf("release/%s", version),
		BoostMatches: true,
	})
	if err != nil {
		return "", err
	}
	for _, b := range branchPage.Values {
		if b.DisplayID == releaseBranch {
			return b.ID, nil
		}
	}
	return "", nil
}

// LatestCommitForBranch returns the latest commit for given repository/branch. If the
// branch is not found, an error is returned.
func LatestCommitForBranch(bitbucketClient bitbucket.BranchClientInterface, projectKey, repositorySlug, branch string) (string, error) {
	branchPage, err := bitbucketClient.BranchList(projectKey, repositorySlug, bitbucket.BranchListParams{
		FilterText:   branch,
		BoostMatches: true,
	})
	if err != nil {
		return "", err
	}
	for _, b := range branchPage.Values {
		if b.ID == branch {
			return b.LatestCommit, nil
		}
	}
	return "", fmt.Errorf("could not find branch %s", branch)
}
