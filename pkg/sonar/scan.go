package sonar

import (
	"fmt"

	"github.com/opendevstack/pipeline/internal/command"
)

type PullRequest struct {
	Key    string
	Branch string
	Base   string
}

func (c *Client) Scan(project, branch, commit string, pr *PullRequest) (string, error) {
	scannerParams := []string{
		fmt.Sprintf("-Dsonar.host.url=%s", c.clientConfig.BaseURL),
		"-Dsonar.scm.provider=git",
		fmt.Sprintf("-Dsonar.projectKey=%s", project),
		fmt.Sprintf("-Dsonar.projectName=%s", project), // TODO: allow to overwrite to cater for multi-repo?
		fmt.Sprintf("-Dsonar.projectVersion=%s", commit),
	}
	if c.clientConfig.Debug {
		scannerParams = append(scannerParams, "-X")
	}
	if pr != nil {
		scannerParams = append(
			scannerParams,
			fmt.Sprintf("-Dsonar.pullrequest.key=%s", pr.Key),
			fmt.Sprintf("-Dsonar.pullrequest.branch=%s", pr.Branch),
			fmt.Sprintf("-Dsonar.pullrequest.base=%s", pr.Base),
		)
	} else if c.clientConfig.ServerEdition != "community" {
		scannerParams = append(scannerParams, fmt.Sprintf("-Dsonar.branch.name=%s", branch))
	}

	fmt.Printf("scan params: %v", scannerParams)
	scannerParams = append(scannerParams, fmt.Sprintf("-Dsonar.login=%s", c.clientConfig.APIToken))
	stdout, stderr, err := command.Run("sonar-scanner", scannerParams)
	if err != nil {
		fmt.Println(string(stdout))
		fmt.Println(string(stderr))
		return "", fmt.Errorf("scanning failed: %w", err)
	}
	return string(stdout), nil
}
