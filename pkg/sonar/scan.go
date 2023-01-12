package sonar

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
)

type PullRequest struct {
	Key    string
	Branch string
	Base   string
}

// Scan report
type ReportTask struct {
	ProjectKey    string
	ServerUrl     string
	ServerVersion string
	Branch        string
	DashboardUrl  string
	CeTaskId      string
	CeTaskUrl     string
}

const (
	ScannerworkDir     = ".scannerwork"
	ReportTaskFilename = "report-task.txt"
	ReportTaskFile     = ScannerworkDir + "/" + ReportTaskFilename
)

// Scan scans the source code and uploads the analysis to given SonarQube project.
// If pr is non-nil, information for pull request decoration is sent.
func (c *Client) Scan(sonarProject, branch, commit string, pr *PullRequest, outWriter, errWriter io.Writer) error {
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
	// Both Branch Analysis and Pull Request Analysis are only available
	// starting in Developer Edition, see
	// https://docs.sonarqube.org/latest/branches/overview/ and
	// https://docs.sonarqube.org/latest/analysis/pull-request/.
	if c.clientConfig.ServerEdition != "community" {
		if pr != nil {
			scannerParams = append(
				scannerParams,
				fmt.Sprintf("-Dsonar.pullrequest.key=%s", pr.Key),
				fmt.Sprintf("-Dsonar.pullrequest.branch=%s", pr.Branch),
				fmt.Sprintf("-Dsonar.pullrequest.base=%s", pr.Base),
			)
		} else {
			scannerParams = append(scannerParams, fmt.Sprintf("-Dsonar.branch.name=%s", branch))
		}
	}

	c.logger().Debugf("Scan params: %v", scannerParams)
	// The authentication token of a SonarQube user with "Execute Analysis"
	// permission on the project is passed as "sonar.login" for authentication,
	// see https://docs.sonarqube.org/latest/analysis/analysis-parameters/.
	scannerParams = append(scannerParams, fmt.Sprintf("-Dsonar.login=%s", c.clientConfig.APIToken))

	return command.Run(
		"sonar-scanner", scannerParams,
		[]string{fmt.Sprintf("SONAR_SCANNER_OPTS=%s", strings.Join(c.javaSystemProperties(), " "))},
		outWriter, errWriter,
	)
}

/*
Example of the file located in .scannerwork/report-task.txt:

	projectKey=XXXX-python
	serverUrl=https://sonarqube-ods.XXXX.com
	serverVersion=8.2.0.32929
	branch=dummy
	dashboardUrl=https://sonarqube-ods.XXXX.com/dashboard?id=XXXX-python&branch=dummy
	ceTaskId=AXxaAoUSsjAMlIY9kNmn
	ceTaskUrl=https://sonarqube-ods.XXXX.com/api/ce/task?id=AXxaAoUSsjAMlIY9kNmn
*/
func (c *Client) ExtractComputeEngineTaskID(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	taskIDPrefix := "ceTaskId="
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, taskIDPrefix) {
			return strings.TrimPrefix(line, taskIDPrefix), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("properties file %s does not contain %s", filename, taskIDPrefix)
}
