package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"sigs.k8s.io/yaml"
)

// helmDiffDetectedMarker is the message Helm prints when helm-diff is
// configured to exit with a non-zero exit code when drift is detected.
const helmDiffDetectedMarker = `Error: identified at least one change, exiting with non-zero exit code (detailed-exitcode parameter enabled)`

// desiredDiffMessage is the message that should be presented to the user.
const desiredDiffMessage = `plugin "diff" identified at least one change`

type helmChart struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// getHelmChart reads given filename into a helmChart struct.
func getHelmChart(filename string) (*helmChart, error) {
	y, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read chart: %w", err)
	}

	var hc *helmChart
	err = yaml.Unmarshal(y, &hc)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
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
func cleanHelmDiffOutput(out []byte) []byte {
	if !bytes.Contains(out, []byte(helmDiffDetectedMarker)) {
		return out
	}
	cleanedOut := bytes.Replace(
		out, []byte(helmDiffDetectedMarker), []byte(desiredDiffMessage), -1,
	)
	r := regexp.MustCompile(`Error: plugin "(diff|secrets)" exited with error[\n]?`)
	cleanedOut = r.ReplaceAll(cleanedOut, []byte{})
	r = regexp.MustCompile(`helm.go:81: \[debug\] plugin "(diff|secrets)" exited with error[\n]?`)
	cleanedOut = r.ReplaceAll(cleanedOut, []byte{})
	return cleanedOut
}

// assembleHelmDiffArgs creates a slice of arguments for "helm diff upgrade".
func assembleHelmDiffArgs(
	releaseNamespace, releaseName, helmArchive string,
	opts options,
	valuesFiles, cliValues []string) ([]string, error) {
	helmDiffArgs := []string{
		"--namespace=" + releaseNamespace,
		"secrets",
		"diff",
		"upgrade",
		"--detailed-exitcode",
		"--no-color",
	}
	helmDiffFlags, err := shlex.Split(opts.diffFlags)
	if err != nil {
		return []string{}, fmt.Errorf("parse diff flags (%s): %s", opts.diffFlags, err)
	}
	helmDiffArgs = append(helmDiffArgs, helmDiffFlags...)
	commonArgs, err := commonHelmUpgradeArgs(releaseName, helmArchive, opts, valuesFiles, cliValues)
	if err != nil {
		return []string{}, fmt.Errorf("upgrade args: %w", err)
	}
	return append(helmDiffArgs, commonArgs...), nil
}

// assembleHelmDiffArgs creates a slice of arguments for "helm upgrade".
func assembleHelmUpgradeArgs(
	releaseNamespace, releaseName, helmArchive string,
	opts options,
	valuesFiles, cliValues []string,
) ([]string, error) {
	helmUpgradeArgs := []string{
		"--namespace=" + releaseNamespace,
		"secrets",
		"upgrade",
	}
	commonArgs, err := commonHelmUpgradeArgs(releaseName, helmArchive, opts, valuesFiles, cliValues)
	if err != nil {
		return []string{}, fmt.Errorf("upgrade args: %w", err)
	}
	return append(helmUpgradeArgs, commonArgs...), nil
}

// commonHelmUpgradeArgs returns arguments common to "helm upgrade" and "helm diff upgrade".
func commonHelmUpgradeArgs(
	releaseName, helmArchive string,
	opts options,
	valuesFiles, cliValues []string,
) ([]string, error) {
	args, err := shlex.Split(opts.upgradeFlags)
	if err != nil {
		return []string{}, fmt.Errorf("parse upgrade flags (%s): %s", opts.upgradeFlags, err)
	}
	for _, vf := range valuesFiles {
		args = append(args, fmt.Sprintf("--values=%s", vf))
	}
	args = append(args, cliValues...)
	args = append(args, releaseName, helmArchive)
	return args, nil
}

// runHelmCmd runs given Helm command.
func runHelmCmd(args []string, targetConfig *config.Environment, debug bool) (outBytes, errBytes []byte, err error) {
	if debug {
		args = append([]string{"--debug"}, args...)
	}
	printableArgs := args
	if targetConfig.APIServer != "" {
		printableArgs = append(
			[]string{
				fmt.Sprintf("--kube-apiserver=%s", targetConfig.APIServer),
				"--kube-token=***",
			},
			args...,
		)
		args = append(
			[]string{
				fmt.Sprintf("--kube-apiserver=%s", targetConfig.APIServer),
				fmt.Sprintf("--kube-token=%s", targetConfig.APIToken),
			},
			args...,
		)
	}
	fmt.Println(helmBin, strings.Join(printableArgs, " "))

	var extraEnvs = []string{
		fmt.Sprintf("SOPS_AGE_KEY_FILE=%s", ageKeyFilePath),
		"HELM_DIFF_IGNORE_UNKNOWN_FLAGS=true", // https://github.com/databus23/helm-diff/issues/278
	}
	return command.RunWithExtraEnvs(helmBin, args, extraEnvs)
}

// packageHelmChart creates a Helm package for given chart.
func packageHelmChart(chartDir, ctxtVersion, gitCommitSHA string, debug bool) (string, error) {
	hc, err := getHelmChart(filepath.Join(chartDir, "Chart.yaml"))
	if err != nil {
		return "", fmt.Errorf("could not read chart: %w", err)
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
	stdout, stderr, err := command.Run(helmBin, append(helmPackageArgs, chartDir))
	if err != nil {
		return "", fmt.Errorf(
			"could not package chart %s. stderr: %s, err: %s", chartDir, string(stderr), err,
		)
	}
	fmt.Println(string(stdout))

	helmArchive := fmt.Sprintf("%s-%s.tgz", hc.Name, packageVersion)
	return helmArchive, nil
}
