package sonar

import (
	"fmt"

	"github.com/opendevstack/pipeline/internal/command"
)

func (c *Client) GenerateReport(project, author, branch string) (string, error) {
	reportParams := []string{
		"-jar", "/usr/local/cnes/cnesreport.jar",
		"-s", c.clientConfig.BaseURL,
		"-t", c.clientConfig.APIToken,
		"-p", project,
		"-a", author,
		branch,
	}
	stdout, stderr, err := command.Run("java", reportParams)
	if err != nil {
		fmt.Println(string(stdout))
		fmt.Println(string(stderr))
		return "", fmt.Errorf("scanning failed: %w", err)
	}
	return string(stdout), nil
}
