package main

import (
	"fmt"
	"io"
	"net/url"
	"os/exec"

	"github.com/opendevstack/pipeline/internal/command"
)

const (
	aquasecBin                           = "aquasec"
	scanComplianceFailureExitCode        = 4
	scanLicenseValidationFailureExitCode = 5
)

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

// aquaScan runs the scan and returns whether there was a policy incompliance or not.
// An error is returned when the scan cannot be started or encounters failures
// unrelated to policy compliance.
func aquaScan(exe string, args []string, outWriter, errWriter io.Writer) (bool, error) {
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
