package sonar

import (
	"fmt"
	"io"
	"os"
	"time"

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

	err = copyReportFiles(project, ".ods/artifacts/sonarqube-analysis")
	if err != nil {
		return "", fmt.Errorf("copying report to artifacts failed: %w", err)
	}

	return string(stdout), nil
}

func copyReportFiles(project, destinationDir string) error {
	analysisReportFile := fmt.Sprintf(
		"%s-%s-analysis-report.md",
		currentDate(),
		project,
	)
	err := copyFile(analysisReportFile, destinationDir+"/analysis-report.md")
	if err != nil {
		return fmt.Errorf("copying %s failed: %w", analysisReportFile, err)
	}

	issuesReportFile := fmt.Sprintf(
		"%s-%s-issues-report.csv",
		currentDate(),
		project,
	)
	err = copyFile(issuesReportFile, destinationDir+"/issues-report.csv")
	if err != nil {
		return fmt.Errorf("copying %s failed: %w", issuesReportFile, err)
	}
	return nil
}

// currentDate returns the current date as YYYY-MM-DD
func currentDate() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01-02")
}

// Copy a file
func copyFile(from, to string) error {
	// Open original file
	original, err := os.Open(from)
	if err != nil {
		return err
	}
	defer original.Close()

	// Create new file
	new, err := os.Create(to)
	if err != nil {
		return err
	}
	defer new.Close()

	// Copy contents of original file into new file
	_, err = io.Copy(new, original)
	if err != nil {
		return err
	}
	return nil
}
