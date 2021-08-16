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

// Scan scans the source code and uploads the analysis to given SonarQube project.
// If pr is non-nil, information for pull request decoration is sent.
func (c *Client) Scan(sonarProject, branch, commit string, pr *PullRequest) (string, error) {
	scannerParams := []string{
		fmt.Sprintf("-Dsonar.host.url=%s", c.clientConfig.BaseURL),
		"-Dsonar.scm.provider=git",
		fmt.Sprintf("-Dsonar.projectKey=%s", sonarProject),
		fmt.Sprintf("-Dsonar.projectName=%s", sonarProject),
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
		return string(stdout), fmt.Errorf("scanning failed: %w, stderr: %s", err, string(stderr))
	}
	return string(stdout), nil
}
