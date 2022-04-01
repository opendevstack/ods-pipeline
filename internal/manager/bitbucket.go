package manager

import (
	"fmt"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
)

type bitbucketInterface interface {
	bitbucket.CommitClientInterface
	bitbucket.RawClientInterface
	bitbucket.RepoClientInterface
}

type repository struct {
	Project struct {
		Key string `json:"key"`
	} `json:"project"`
	Slug string `json:"slug"`
}
type requestBitbucket struct {
	EventKey   string     `json:"eventKey"`
	Repository repository `json:"repository"`
	Changes    []struct {
		Type string `json:"type"`
		Ref  struct {
			ID        string `json:"id"`
			DisplayID string `json:"displayId"`
			Type      string `json:"type"`
		} `json:"ref"`
		FromHash string `json:"fromHash"`
		ToHash   string `json:"toHash"`
	} `json:"changes"`
	PullRequest *struct {
		FromRef struct {
			Repository   repository `json:"repository"`
			ID           string     `json:"id"`
			DisplayID    string     `json:"displayId"`
			LatestCommit string     `json:"latestCommit"`
		} `json:"fromRef"`
	} `json:"pullRequest"`
	Comment *struct {
		Text string `json:"text"`
	} `json:"comment"`
}

func getCommitSHA(bitbucketClient bitbucket.CommitClientInterface, project, repository, gitFullRef string) (string, error) {
	commitList, err := bitbucketClient.CommitList(project, repository, bitbucket.CommitListParams{
		Until: gitFullRef,
	})
	if err != nil {
		return "", fmt.Errorf("could not get commit list: %w", err)
	}
	return commitList.Values[0].ID, nil
}

type prInfo struct {
	ID   int
	Base string
}

func extractPullRequestInfo(bitbucketClient bitbucket.CommitClientInterface, projectKey, repositorySlug, gitCommit string) (prInfo, error) {
	var i prInfo

	prPage, err := bitbucketClient.CommitPullRequestList(projectKey, repositorySlug, gitCommit)
	if err != nil {
		return i, err
	}

	for _, v := range prPage.Values {
		if !v.Open {
			continue
		}
		i.ID = v.ID
		i.Base = v.ToRef.ID
		break
	}

	return i, nil
}

func shouldSkip(bitbucketClient bitbucket.CommitClientInterface, projectKey, repositorySlug, gitCommit string) bool {
	c, err := bitbucketClient.CommitGet(projectKey, repositorySlug, gitCommit)
	if err != nil {
		return false
	}
	return isCiSkipInCommitMessage(c.Message)
}

// getRepoNames retrieves the name of all repositories within the project
// identified by projectKey.
func getRepoNames(bitbucketClient bitbucket.RepoClientInterface, projectKey string) ([]string, error) {
	repos := []string{}
	rl, err := bitbucketClient.RepoList(projectKey)
	if err != nil {
		return repos, err
	}
	for _, n := range rl.Values {
		repos = append(repos, n.Name)
	}
	return repos, nil
}
