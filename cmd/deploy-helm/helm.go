package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"sigs.k8s.io/yaml"
)

const (
	// helmDiffDetectedMarker is the message Helm prints when helm-diff is
	// configured to exit with a non-zero exit code when drift is detected.
	helmDiffDetectedMarker = `Error: identified at least one change, exiting with non-zero exit code (detailed-exitcode parameter enabled)`

	// desiredDiffMessage is the message that should be presented to the user.
	desiredDiffMessage = `plugin "diff" identified at least one change`

	// exit code returned from helm-diff when diff is detected.
	diffDriftExitCode = 2

	// exit code returned from helm-diff when there is an error (e.g. invalid resource manifests).
	diffGenericExitCode = 1
)

type helmChart struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// helmDiff runs the diff and returns whether the Helm release is in sync.
// An error is returned when the diff cannot be started or encounters failures
// unrelated to drift (such as invalid resource manifests).
func (d *deployHelm) helmDiff(args []string, outWriter, errWriter io.Writer) (bool, error) {
	return command.RunWithSpecialFailureCode(
		d.helmBin, args, []string{
			fmt.Sprintf("SOPS_AGE_KEY_FILE=%s", ageKeyFilePath),
			"HELM_DIFF_IGNORE_UNKNOWN_FLAGS=true", // https://github.com/databus23/helm-diff/issues/278
		}, outWriter, errWriter, diffDriftExitCode,
	)
}

// helmUpgrade runs given Helm command.
func (d *deployHelm) helmUpgrade(args []string, stdout, stderr io.Writer) error {
	return command.Run(
		d.helmBin, args, []string{fmt.Sprintf("SOPS_AGE_KEY_FILE=%s", ageKeyFilePath)}, stdout, stderr,
	)
}

// assembleHelmDiffArgs creates a slice of arguments for "helm diff upgrade".
func (d *deployHelm) assembleHelmDiffArgs() ([]string, error) {
	helmDiffArgs := []string{
		"--namespace=" + d.releaseNamespace,
		"secrets",
		"diff",
		"upgrade",
		"--detailed-exitcode",
		"--no-color",
		"--normalize-manifests",
	}
	helmDiffFlags, err := shlex.Split(d.opts.diffFlags)
	if err != nil {
		return []string{}, fmt.Errorf("parse diff flags (%s): %s", d.opts.diffFlags, err)
	}
	helmDiffArgs = append(helmDiffArgs, helmDiffFlags...)
	commonArgs, err := d.commonHelmUpgradeArgs()
	if err != nil {
		return []string{}, fmt.Errorf("upgrade args: %w", err)
	}
	return append(helmDiffArgs, commonArgs...), nil
}

// assembleHelmDiffArgs creates a slice of arguments for "helm upgrade".
func (d *deployHelm) assembleHelmUpgradeArgs() ([]string, error) {
	helmUpgradeArgs := []string{
		"--namespace=" + d.releaseNamespace,
		"secrets",
		"upgrade",
	}
	commonArgs, err := d.commonHelmUpgradeArgs()
	if err != nil {
		return []string{}, fmt.Errorf("upgrade args: %w", err)
	}
	return append(helmUpgradeArgs, commonArgs...), nil
}

// commonHelmUpgradeArgs returns arguments common to "helm upgrade" and "helm diff upgrade".
func (d *deployHelm) commonHelmUpgradeArgs() ([]string, error) {
	args, err := shlex.Split(d.opts.upgradeFlags)
	if err != nil {
		return []string{}, fmt.Errorf("parse upgrade flags (%s): %s", d.opts.upgradeFlags, err)
	}
	if d.opts.debug {
		args = append([]string{"--debug"}, args...)
	}
	if d.targetConfig.APIServer != "" {
		args = append(
			[]string{
				fmt.Sprintf("--kube-apiserver=%s", d.targetConfig.APIServer),
				fmt.Sprintf("--kube-token=%s", d.targetConfig.APIToken),
			},
			args...,
		)
	}
	for _, vf := range d.valuesFiles {
		args = append(args, fmt.Sprintf("--values=%s", vf))
	}
	args = append(args, d.cliValues...)
	args = append(args, d.releaseName, d.helmArchive)
	return args, nil
}

// getHelmChart reads given filename into a helmChart struct.
func getHelmChart(filename string) (*helmChart, error) {
	y, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read chart file: %w", err)
	}

	var hc *helmChart
	err = yaml.Unmarshal(y, &hc)
	if err != nil {
		return nil, fmt.Errorf("unmarshal chart: %w", err)
	}
	return hc, nil
}

// getChartVersion extracts the version from given Helm chart.
func getChartVersion(contextVersion string, hc *helmChart) string {
	if len(contextVersion) > 0 && contextVersion != pipelinectxt.WIP {
		return contextVersion
	}
	return hc.Version
}

// cleanHelmDiffOutput removes error messages from the given Helm output.
// Those error messages are confusing, because they do not come from  an actual
// error, but from detecting drift between desired and current Helm state.
func cleanHelmDiffOutput(out string) string {
	if !strings.Contains(out, helmDiffDetectedMarker) {
		return out
	}
	cleanedOut := strings.Replace(
		out, helmDiffDetectedMarker, desiredDiffMessage, -1,
	)
	r := regexp.MustCompile(`Error: plugin "(diff|secrets)" exited with error[\n]?`)
	cleanedOut = r.ReplaceAllString(cleanedOut, "")
	r = regexp.MustCompile(`helm.go:81: \[debug\] plugin "(diff|secrets)" exited with error[\n]?`)
	cleanedOut = r.ReplaceAllString(cleanedOut, "")
	return cleanedOut
}

// printlnSafeHelmCmd prints all args that do not contain sensitive information.
func printlnSafeHelmCmd(args []string, outWriter io.Writer) {
	safeArgs := []string{}
	for _, a := range args {
		if strings.HasPrefix(a, "--kube-token=") {
			safeArgs = append(safeArgs, "--kube-token=***")
		} else {
			safeArgs = append(safeArgs, a)
		}
	}
	fmt.Fprintln(outWriter, helmBin, strings.Join(safeArgs, " "))
}

// packageHelmChart creates a Helm package for given chart.
func packageHelmChart(chartDir, ctxtVersion, gitCommitSHA string, debug bool) (string, error) {
	hc, err := getHelmChart(filepath.Join(chartDir, "Chart.yaml"))
	if err != nil {
		return "", fmt.Errorf("read chart: %w", err)
	}
	chartVersion := getChartVersion(ctxtVersion, hc)
	packageVersion := fmt.Sprintf("%s+%s", chartVersion, gitCommitSHA)
	helmPackageArgs := []string{
		"package",
		fmt.Sprintf("--app-version=%s", gitCommitSHA),
		fmt.Sprintf("--version=%s", packageVersion),
	}
	if debug {
		helmPackageArgs = append(helmPackageArgs, "--debug")
	}
	err = command.Run(helmBin, append(helmPackageArgs, chartDir), []string{}, os.Stdout, os.Stderr)
	if err != nil {
		return "", fmt.Errorf("package chart %s: %w", chartDir, err)
	}

	helmArchive := fmt.Sprintf("%s-%s.tgz", hc.Name, packageVersion)
	return helmArchive, nil
}
