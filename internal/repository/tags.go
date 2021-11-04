package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

// TagListContainsFinalVersion checks if the list of tags contains a tag
// corresponding to the version (without pre-release/build suffix).
func TagListContainsFinalVersion(tags []bitbucket.Tag, version string) bool {
	searchID := fmt.Sprintf("refs/tags/v%s", version)
	for _, t := range tags {
		if t.ID == searchID {
			return true
		}
	}
	return false
}

// LatestReleaseCandidate returns the highest number of all tags of
// format "v<VERSION>-rc.<NUMBER>".
func LatestReleaseCandidate(tags []bitbucket.Tag, version string) (*bitbucket.Tag, int) {
	var highestNumber int
	var latestTag *bitbucket.Tag
	prefix := fmt.Sprintf("refs/tags/v%s-rc.", version)
	for _, t := range tags {
		if strings.HasPrefix(t.ID, prefix) {
			i, err := strconv.Atoi(strings.TrimPrefix(t.ID, prefix))
			if err == nil && i > highestNumber {
				highestNumber = i
				latestTag = &t
			}
		}
	}
	return latestTag, highestNumber
}

// CreateTag creates a Git tag with given name in the reopsitory identified
// in the context.
func CreateTag(bitbucketClient *bitbucket.Client, ctxt *pipelinectxt.ODSContext, name string) (*bitbucket.Tag, error) {
	return bitbucketClient.TagCreate(
		ctxt.Project,
		ctxt.Repository,
		bitbucket.TagCreatePayload{
			Name:       name,
			StartPoint: ctxt.GitCommitSHA,
		},
	)
}
