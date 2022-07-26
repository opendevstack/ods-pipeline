package main

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"

	"github.com/opendevstack/pipeline/internal/command"
)

const (
	aquasecBin                           = "aquasec"
	scanComplianceFailureExitCode        = 4
	scanLicenseValidationFailureExitCode = 5
)

// aquaScanRunner exists for testing purposes.
type aquaScanRunner interface {
	aquaScanRun(opts options, image, htmlReportFile, jsonReportFile string) ([]byte, []byte, error)
}

// aquaScanRunnerFunc is a type implementing the aquaScanRunner interface (like http.HandlerFunc).
type aquaScanRunnerFunc func(opts options, image, htmlReportFile, jsonReportFile string) ([]byte, []byte, error)

// aquaScanRun calls the function it is implemented on.
func (f aquaScanRunnerFunc) aquaScanRun(opts options, image, htmlReportFile, jsonReportFile string) ([]byte, []byte, error) {
	return f(opts, image, htmlReportFile, jsonReportFile)
}

// aquaScanRun runs an Aqua scan on given image.
func aquaScanRun(opts options, image, htmlReportFile, jsonReportFile string) ([]byte, []byte, error) {
	return command.Run(aquasecBin, []string{
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
	})
}

// aquaScan executes the given aquaScanRunner and handles its result.
// If the scan run did not return an err, the scan is considered successful.
// If the scan did return an err, the scan is not successful, and an error is
// returned unless the scan run indicates that the failure is due to compliance
// reasons.
func aquaScan(as aquaScanRunner, opts options, aquaImage, htmlReportFile, jsonReportFile string) (bool, string, error) {
	scanSuccessful := false
	// STDERR contains the scan log output
	// STDOUT contains the scan summary (incl. ASCII table)
	stdout, stderr, err := as.aquaScanRun(opts, aquaImage, htmlReportFile, jsonReportFile)
	out := fmt.Sprintf("%s\n%s", string(stderr), string(stdout))
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) { // failure during Aqua scan
			if ee.ExitCode() != scanComplianceFailureExitCode {
				return false, out, fmt.Errorf("scan failed: %w", err)
			}
		} else { // error e.g. when binary is not found / executable
			return false, out, fmt.Errorf("scan error: %w", err)
		}
	} else {
		scanSuccessful = true
	}
	return scanSuccessful, out, nil
}

// aquasecInstalled checks whether the Aqua binary is in the $PATH.
func aquasecInstalled() bool {
	_, err := exec.LookPath(aquasecBin)
	return err == nil
}

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
