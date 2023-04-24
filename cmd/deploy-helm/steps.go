package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/file"
	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	tokenFile    = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	subchartsDir = "charts"
)

type DeployStep func(d *deployHelm) (*deployHelm, error)

func (d *deployHelm) runSteps(steps ...DeployStep) error {
	var skip *skipRemainingSteps
	var err error
	for _, step := range steps {
		d, err = step(d)
		if err != nil {
			if errors.As(err, &skip) {
				d.logger.Infof(err.Error())
				return nil
			}
			return err
		}
	}
	return nil
}

func setupContext() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		ctxt := &pipelinectxt.ODSContext{}
		err := ctxt.ReadCache(d.opts.checkoutDir)
		if err != nil {
			return d, fmt.Errorf("read cache: %w", err)
		}
		d.ctxt = ctxt

		clientset, err := k.NewInClusterClientset()
		if err != nil {
			return d, fmt.Errorf("create Kubernetes clientset: %w", err)
		}
		d.clientset = clientset

		if d.opts.debug {
			if err := directory.ListFiles(d.opts.certDir, os.Stdout); err != nil {
				log.Fatal(err)
			}
		}
		return d, nil
	}
}

func skipOnEmptyNamespace() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		if d.opts.namespace == "" {
			return d, &skipRemainingSteps{"No namespace given. Skipping deployment ..."}
		}
		return d, nil
	}
}

func setReleaseTarget() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		// Release name
		if d.opts.releaseName != "" {
			d.releaseName = d.opts.releaseName
		} else {
			d.releaseName = d.ctxt.Component
		}
		d.logger.Infof("Release name: %s", d.releaseName)

		// Target environment configuration
		targetConfig := &targetEnvironment{
			APIServer:    d.opts.apiServer,
			Namespace:    d.opts.namespace,
			RegistryHost: d.opts.registryHost,
		}
		if targetConfig.APIServer != "" {
			token, err := tokenFromSecret(d.clientset, d.ctxt.Namespace, d.opts.apiCredentialsSecret)
			if err != nil {
				return d, fmt.Errorf("get API token from secret %s: %w", d.opts.apiCredentialsSecret, err)
			}
			targetConfig.APIToken = token
		}
		d.targetConfig = targetConfig

		// Release namespace
		d.releaseNamespace = targetConfig.Namespace
		pattern := "^[a-z][a-z0-9-]{0,61}[a-z]$"
		matched, err := regexp.MatchString(pattern, d.releaseNamespace)
		if err != nil || !matched {
			return d, fmt.Errorf("release namespace: %s must match %s", d.releaseNamespace, pattern)
		}
		d.logger.Infof("Release namespace: %s", d.releaseNamespace)

		return d, nil
	}
}

func detectSubrepos() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		subrepos, err := pipelinectxt.DetectSubrepos()
		if err != nil {
			return d, fmt.Errorf("detect subrepos: %w", err)
		}
		d.subrepos = subrepos
		return d, nil
	}
}

func detectImageDigests() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		digests, err := pipelinectxt.ReadArtifactFilesIncludingSubrepos(pipelinectxt.ImageDigestsPath, d.subrepos)
		if err != nil {
			return d, fmt.Errorf("collect image digests: %w", err)
		}
		d.imageDigests = digests
		return d, nil
	}
}

func copyImagesIntoReleaseNamespace() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		if len(d.imageDigests) == 0 {
			return d, nil
		}
		// Get destination registry token from secret or file in pod.
		var destRegistryToken string
		if d.targetConfig.APIToken != "" {
			destRegistryToken = d.targetConfig.APIToken
		} else {
			token, err := getTrimmedFileContent(tokenFile)
			if err != nil {
				return d, fmt.Errorf("get token from file %s: %w", tokenFile, err)
			}
			destRegistryToken = token
		}

		d.logger.Infof("Copying images into release namespace ...")
		for _, artifactFile := range d.imageDigests {
			imageArtifact, err := artifact.ReadFromFile(artifactFile)
			if err != nil {
				return d, fmt.Errorf("read image artifact %s: %w", artifactFile, err)
			}
			err = d.copyImage(*imageArtifact, destRegistryToken, os.Stdout, os.Stderr)
			if err != nil {
				return d, fmt.Errorf("copy image %s: %w", imageArtifact.Name, err)
			}
		}

		return d, nil
	}
}

func listHelmPlugins() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		d.logger.Infof("List Helm plugins...")
		helmPluginArgs := []string{"plugin", "list"}
		if d.opts.debug {
			helmPluginArgs = append(helmPluginArgs, "--debug")
		}
		err := command.Run(d.helmBin, helmPluginArgs, []string{}, os.Stdout, os.Stderr)
		if err != nil {
			return d, fmt.Errorf("list Helm plugins: %w", err)
		}
		return d, nil
	}
}

func packageHelmChartWithSubcharts() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		// Collect values to be set via the CLI.
		d.cliValues = []string{
			fmt.Sprintf("--set=image.tag=%s", d.ctxt.GitCommitSHA),
		}

		d.logger.Infof("Adding dependencies from subrepos into the %s/ directory ...", subchartsDir)
		// Find subcharts
		chartsDir := filepath.Join(d.opts.chartDir, subchartsDir)
		if _, err := os.Stat(chartsDir); os.IsNotExist(err) {
			err = os.Mkdir(chartsDir, 0755)
			if err != nil {
				return d, fmt.Errorf("create %s: %s", chartsDir, err)
			}
		}
		for _, r := range d.subrepos {
			subrepo := filepath.Join(pipelinectxt.SubreposPath, r.Name())
			subchart := filepath.Join(subrepo, d.opts.chartDir)
			if _, err := os.Stat(subchart); os.IsNotExist(err) {
				d.logger.Infof("no chart in %s", r.Name())
				continue
			}
			gitCommitSHA, err := getTrimmedFileContent(filepath.Join(subrepo, ".ods", "git-commit-sha"))
			if err != nil {
				return d, fmt.Errorf("get commit SHA of %s: %w", subrepo, err)
			}
			hc, err := getHelmChart(filepath.Join(subchart, "Chart.yaml"))
			if err != nil {
				return d, fmt.Errorf("get Helm chart of %s: %w", subrepo, err)
			}
			d.cliValues = append(d.cliValues, fmt.Sprintf("--set=%s.image.tag=%s", hc.Name, gitCommitSHA))
			if d.releaseName == d.ctxt.Component {
				d.cliValues = append(d.cliValues, fmt.Sprintf("--set=%s.fullnameOverride=%s", hc.Name, hc.Name))
			}
			helmArchive, err := packageHelmChart(subchart, d.ctxt.Version, gitCommitSHA, d.opts.debug)
			if err != nil {
				return d, fmt.Errorf("package Helm chart of %s: %w", subrepo, err)
			}
			helmArchiveName := filepath.Base(helmArchive)
			d.logger.Infof("copying %s into %s", helmArchiveName, chartsDir)
			err = file.Copy(helmArchive, filepath.Join(chartsDir, helmArchiveName))
			if err != nil {
				return d, fmt.Errorf("copy Helm archive of %s: %w", subrepo, err)
			}
		}

		subcharts, err := os.ReadDir(chartsDir)
		if err != nil {
			return d, fmt.Errorf("read %s: %w", chartsDir, err)
		}
		if len(subcharts) > 0 {
			d.logger.Infof("Subcharts in %s:", chartsDir)
			for _, sc := range subcharts {
				d.logger.Infof(sc.Name())
			}
		}

		d.logger.Infof("Packaging Helm chart ...")
		helmArchive, err := packageHelmChart(d.opts.chartDir, d.ctxt.Version, d.ctxt.GitCommitSHA, d.opts.debug)
		if err != nil {
			return d, fmt.Errorf("package Helm chart: %w", err)
		}
		d.helmArchive = helmArchive
		return d, nil
	}
}

func collectValuesFiles() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		d.logger.Infof("Collecting Helm values files ...")
		d.valuesFiles = []string{}
		valuesFilesCandidates := []string{
			fmt.Sprintf("%s/secrets.yaml", d.opts.chartDir), // equivalent values.yaml is added automatically by Helm
			fmt.Sprintf("%s/values.%s.yaml", d.opts.chartDir, d.targetConfig.Namespace),
			fmt.Sprintf("%s/secrets.%s.yaml", d.opts.chartDir, d.targetConfig.Namespace),
		}
		for _, vfc := range valuesFilesCandidates {
			if _, err := os.Stat(vfc); os.IsNotExist(err) {
				d.logger.Infof("%s is not present, skipping.", vfc)
			} else {
				d.logger.Infof("%s is present, adding.", vfc)
				d.valuesFiles = append(d.valuesFiles, vfc)
			}
		}
		return d, nil
	}
}

func importAgeKey() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		if len(d.opts.ageKeySecret) == 0 {
			d.logger.Infof("Skipping import of age key for helm-secrets as parameter is not set ...")
			return d, nil
		}
		d.logger.Infof("Storing age key for helm-secrets ...")
		secret, err := d.clientset.CoreV1().Secrets(d.ctxt.Namespace).Get(
			context.TODO(), d.opts.ageKeySecret, metav1.GetOptions{},
		)
		if err != nil {
			d.logger.Infof("No secret %s found, skipping.", d.opts.ageKeySecret)
			return d, nil
		}
		err = storeAgeKey(secret.Data[d.opts.ageKeySecretField])
		if err != nil {
			return d, fmt.Errorf("store age key: %w", err)
		}
		d.logger.Infof("Age key secret %s stored.", d.opts.ageKeySecret)
		return d, nil
	}
}

func diffHelmRelease() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		d.logger.Infof("Diffing Helm release against %s...", d.helmArchive)
		helmDiffArgs, err := d.assembleHelmDiffArgs()
		if err != nil {
			return d, fmt.Errorf("assemble helm diff args: %w", err)
		}
		printlnSafeHelmCmd(helmDiffArgs, os.Stdout)
		// helm-dff stderr contains confusing text about "errors" when drift is
		// detected, therefore we want to collect and polish it before we print it.
		// helm-diff stdout needs to be written into a buffer so that we can both
		// print it and store it later as a deployment artifact.
		var diffStdoutBuf, diffStderrBuf bytes.Buffer
		diffStdoutWriter := io.MultiWriter(os.Stdout, &diffStdoutBuf)
		inSync, err := d.helmDiff(helmDiffArgs, diffStdoutWriter, &diffStderrBuf)
		fmt.Print(cleanHelmDiffOutput(diffStderrBuf.String()))
		if err != nil {
			return d, fmt.Errorf("helm diff: %w", err)
		}
		if inSync {
			return d, &skipRemainingSteps{"No diff detected, skipping helm upgrade."}
		}

		err = writeDeploymentArtifact(diffStdoutBuf.Bytes(), "diff", d.opts.chartDir, d.targetConfig.Namespace)
		if err != nil {
			return d, fmt.Errorf("write diff artifact: %w", err)
		}
		return d, nil
	}
}

func upgradeHelmRelease() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		d.logger.Infof("Upgrading Helm release to %s...", d.helmArchive)
		helmUpgradeArgs, err := d.assembleHelmUpgradeArgs()
		if err != nil {
			return d, fmt.Errorf("assemble helm upgrade args: %w", err)
		}
		printlnSafeHelmCmd(helmUpgradeArgs, os.Stdout)
		var upgradeStdoutBuf bytes.Buffer
		upgradeStdoutWriter := io.MultiWriter(os.Stdout, &upgradeStdoutBuf)
		err = d.helmUpgrade(helmUpgradeArgs, upgradeStdoutWriter, os.Stderr)
		if err != nil {
			return d, fmt.Errorf("helm upgrade: %w", err)
		}
		err = writeDeploymentArtifact(upgradeStdoutBuf.Bytes(), "release", d.opts.chartDir, d.targetConfig.Namespace)
		if err != nil {
			return d, fmt.Errorf("write release artifact: %w", err)
		}
		return d, nil
	}
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func tokenFromSecret(clientset *kubernetes.Clientset, namespace, name string) (string, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data["token"]), nil
}

func writeDeploymentArtifact(content []byte, filename, chartDir, targetEnv string) error {
	err := os.MkdirAll(pipelinectxt.DeploymentsPath, 0755)
	if err != nil {
		return err
	}
	f := artifactFilename(filename, chartDir, targetEnv) + ".txt"
	return os.WriteFile(filepath.Join(pipelinectxt.DeploymentsPath, f), content, 0644)
}

func artifactFilename(filename, chartDir, targetEnv string) string {
	trimmedChartDir := strings.TrimPrefix(chartDir, "./")
	if trimmedChartDir != "chart" {
		filename = fmt.Sprintf("%s-%s", strings.Replace(trimmedChartDir, "/", "-", -1), filename)
	}
	return fmt.Sprintf("%s-%s", filename, targetEnv)
}
