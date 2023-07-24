package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/image"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
)

const (
	aquasecBin                           = "./.ods-cache/bin/aquasec"
	scanComplianceFailureExitCode        = 4
	scanLicenseValidationFailureExitCode = 5
)

// aquaScanURL returns an URL to the given aquaImage.
func aquaScanURL(opts options, aquaImage string) (string, error) {
	aquaURL, err := url.Parse(opts.aquaURL)
	if err != nil {
		return "", fmt.Errorf("parse base URL: %w", err)
	}
	aquaPath := fmt.Sprintf(
		"/#/images/%s/%s/vulns",
		url.QueryEscape(opts.aquaRegistry), url.QueryEscape(aquaImage),
	)
	fullURL, err := aquaURL.Parse(aquaPath)
	if err != nil {
		return "", fmt.Errorf("parse URL path: %w", err)
	}
	return fullURL.String(), nil
}

// aquaScan runs the scan and returns whether there was a policy incompliance or not.
// An error is returned when the scan cannot be started or encounters failures
// unrelated to policy compliance.
func runScan(exe string, args []string, outWriter, errWriter io.Writer) (bool, error) {
	// STDERR contains the scan log output, hence we read it before STDOUT.
	// STDOUT contains the scan summary (incl. ASCII table).
	return command.RunWithSpecialFailureCode(
		exe, args, []string{}, outWriter, errWriter, scanComplianceFailureExitCode,
	)
}

// aquaAssembleScanArgs creates args/flags to pass to the Aqua scanner based on given arguments.
func aquaAssembleScanArgs(opts options, image, htmlReportFile, jsonReportFile string) []string {
	return []string{
		"scan",
		"--dockerless", "--register", "--text",
		fmt.Sprintf("--htmlfile=%s", htmlReportFile),
		fmt.Sprintf("--jsonfile=%s", jsonReportFile),
		"-w", "/tmp",
		fmt.Sprintf("--user=%s", opts.aquaUsername),
		fmt.Sprintf("--password=%s", opts.aquaPassword),
		fmt.Sprintf("--host=%s", opts.aquaURL),
		image,
		fmt.Sprintf("--registry=%s", opts.aquaRegistry),
	}
}

// htmlReportFilename returns the HTML report filename for given image.
func htmlReportFilename(iid image.Identity) string {
	return fmt.Sprintf("%s.html", iid.ImageStream)
}

// htmlReportFilename returns the JSON report filename for given image.
func jsonReportFilename(iid image.Identity) string {
	return fmt.Sprintf("%s.json", iid.ImageStream)
}

// reportFilenames returns the list of scan report filenames.
func reportFilenames(iid image.Identity) []string {
	return []string{htmlReportFilename(iid), jsonReportFilename(iid)}
}

// aquaReportsExist checks whether the reports associated with the image name
// exist in the given artifacts path.
func aquaReportsExist(artifactsPath string, iid image.Identity) bool {
	d := filepath.Join(artifactsPath, pipelinectxt.AquaScansDir)
	for _, f := range reportFilenames(iid) {
		if _, err := os.Stat(filepath.Join(d, f)); err != nil {
			return false
		}
	}
	return true
}

// copyAquaReportsToArtifacts copies the Aqua scan reports to the artifacts directory.
func copyAquaReportsToArtifacts(htmlReportFile, jsonReportFile string) error {
	if _, err := os.Stat(htmlReportFile); err == nil {
		err := pipelinectxt.CopyArtifact(htmlReportFile, pipelinectxt.AquaScansPath)
		if err != nil {
			return fmt.Errorf("copying HTML report to artifacts failed: %w", err)
		}
	}
	if _, err := os.Stat(jsonReportFile); err == nil {
		err := pipelinectxt.CopyArtifact(jsonReportFile, pipelinectxt.AquaScansPath)
		if err != nil {
			return fmt.Errorf("copying JSON report to artifacts failed: %w", err)
		}
	}
	return nil
}
