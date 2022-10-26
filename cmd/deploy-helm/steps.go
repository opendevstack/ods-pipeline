package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/file"
	k "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeployStep func(x *deployHelm) (*deployHelm, error)

func setupContext() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		ctxt := &pipelinectxt.ODSContext{}
		err := ctxt.ReadCache(d.opts.checkoutDir)
		if err != nil {
			return d, err
		}
		d.ctxt = ctxt

		clientset, err := k.NewInClusterClientset()
		if err != nil {
			return d, fmt.Errorf("could not create Kubernetes client: %w", err)
		}
		d.clientset = clientset

		d.logger.Debugf("Certificates:\n%s")
		if d.opts.debug {
			// TODO: make this take a writer
			directory.ListFiles(d.opts.certDir)
		}
		return d, nil
	}
}

func skipOnEmptyEnv() DeployStep {
	return func(l *deployHelm) (*deployHelm, error) {
		if len(l.ctxt.Environment) == 0 {
			return l, &skipFollowingSteps{"No environment to deploy to selected. Skipping deployment ..."}
		}
		return l, nil
	}
}

func determineReleaseTarget() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		rt := &releaseTarget{}
		if len(d.opts.releaseName) > 0 {
			rt.name = d.opts.releaseName
		} else {
			rt.name = d.ctxt.Component
		}
		d.logger.Infof("releaseName=%s", rt.name)

		// read ods.y(a)ml
		odsConfig, err := config.ReadFromDir(d.opts.checkoutDir)
		if err != nil {
			return d, fmt.Errorf("err during ods config reading: %w", err)

		}
		targetConfig, err := odsConfig.Environment(d.ctxt.Environment)
		if err != nil {
			return d, fmt.Errorf("err during namespace extraction: %w", err)
		}

		if targetConfig.APIServer != "" {
			token, err := tokenFromSecret(d.clientset, d.ctxt.Namespace, targetConfig.APICredentialsSecret)
			if err != nil {
				return d, fmt.Errorf("could not get token from secret %s: %w", targetConfig.APICredentialsSecret, err)
			}
			targetConfig.APIToken = token
		}

		rt.config = targetConfig

		rt.namespace = targetConfig.Namespace
		if len(rt.namespace) == 0 {
			rt.namespace = fmt.Sprintf("%s-%s", d.ctxt.Project, targetConfig.Name)
		}
		d.logger.Infof("releaseNamespace=%s", rt.namespace)

		d.releaseTarget = rt
		return d, nil
	}
}

func determineSubrepos() DeployStep {
	return func(x *deployHelm) (*deployHelm, error) {
		x.subrepos = []fs.DirEntry{}
		if _, err := os.Stat(pipelinectxt.SubreposPath); err == nil {
			f, err := os.ReadDir(pipelinectxt.SubreposPath)
			if err != nil {
				return x, fmt.Errorf("cannot read %s: %w", pipelinectxt.SubreposPath, err)
			}
			x.subrepos = f
		}
		return x, nil
	}
}

func determineImageDigests() DeployStep {
	return func(x *deployHelm) (*deployHelm, error) {
		var files []string
		id, err := collectImageDigests(pipelinectxt.ImageDigestsPath)
		if err != nil {
			return x, err
		}
		files = append(files, id...)
		for _, s := range x.subrepos {
			subrepoImageDigestsPath := filepath.Join(pipelinectxt.SubreposPath, s.Name(), pipelinectxt.ImageDigestsPath)
			id, err := collectImageDigests(subrepoImageDigestsPath)
			if err != nil {
				return x, err
			}
			files = append(files, id...)
		}
		x.imageDigests = files
		return x, nil
	}
}

func copyImagesIntoReleaseNamespace() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {

		// Copy images into release namespace if there are any image artifacts.
		if len(d.imageDigests) > 0 {
			// Get destination registry token from secret or file in pod.
			var destRegistryToken string
			if d.targetConfig().APIToken != "" {
				destRegistryToken = d.targetConfig().APIToken
			} else {
				token, err := getTrimmedFileContent(tokenFile)
				if err != nil {
					return d, fmt.Errorf("could not get token from file %s: %w", tokenFile, err)
				}
				destRegistryToken = token
			}

			d.logger.Infof("Copying images into release namespace ...")
			for _, artifactFile := range d.imageDigests {
				var imageArtifact artifact.Image
				artifactContent, err := os.ReadFile(artifactFile)
				if err != nil {
					return d, fmt.Errorf("could not read image artifact file %s: %w", artifactFile, err)
				}
				err = json.Unmarshal(artifactContent, &imageArtifact)
				if err != nil {
					return d, fmt.Errorf(
						"could not unmarshal image artifact file %s: %w.\nFile content:\n%s",
						artifactFile, err, string(artifactContent),
					)
				}
				imageStream := imageArtifact.Name
				d.logger.Infof("Copying image %s ...", imageStream)
				srcImageURL := imageArtifact.Image
				// If the source registry should be TLS verified, the destination
				// should be verified by default as well.
				destRegistryTLSVerify := d.opts.srcRegistryTLSVerify
				srcRegistryTLSVerify := d.opts.srcRegistryTLSVerify
				// TLS verification of the KinD registry is not possible at the moment as
				// requests error out with "server gave HTTP response to HTTPS client".
				if strings.HasPrefix(imageArtifact.Registry, "kind-registry.kind") {
					srcRegistryTLSVerify = false
				}
				if len(d.targetConfig().RegistryHost) > 0 && d.targetConfig().RegistryTLSVerify != nil {
					destRegistryTLSVerify = *d.targetConfig().RegistryTLSVerify
				}
				destImageURL := getImageDestURL(d.targetConfig().RegistryHost, d.releaseNamespace(), imageArtifact)
				fmt.Printf("src=%s\n", srcImageURL)
				fmt.Printf("dest=%s\n", destImageURL)
				// TODO: for QA and PROD we want to ensure that the SHA recorded in Nexus
				// matches the SHA referenced by the Git commit tag.
				skopeoCopyArgs := []string{
					"copy",
					fmt.Sprintf("--src-tls-verify=%v", srcRegistryTLSVerify),
					fmt.Sprintf("--dest-tls-verify=%v", destRegistryTLSVerify),
				}
				if srcRegistryTLSVerify {
					skopeoCopyArgs = append(skopeoCopyArgs, fmt.Sprintf("--src-cert-dir=%v", d.opts.certDir))
				}
				if destRegistryTLSVerify {
					skopeoCopyArgs = append(skopeoCopyArgs, fmt.Sprintf("--dest-cert-dir=%v", d.opts.certDir))
				}
				if len(destRegistryToken) > 0 {
					skopeoCopyArgs = append(skopeoCopyArgs, "--dest-registry-token", destRegistryToken)
				}
				if d.opts.debug {
					skopeoCopyArgs = append(skopeoCopyArgs, "--debug")
				}
				stdout, stderr, err := command.RunBuffered(
					"skopeo", append(
						skopeoCopyArgs,
						fmt.Sprintf("docker://%s", srcImageURL),
						fmt.Sprintf("docker://%s", destImageURL),
					),
				)
				if err != nil {
					return d, fmt.Errorf("copying failed: %w, stderr = %s", err, string(stderr))
				}
				fmt.Println(string(stdout))
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
		stdout, stderr, err := command.RunBuffered(d.helmBin, helmPluginArgs)
		if err != nil {
			fmt.Println(string(stderr))
			log.Fatal(err)
		}
		fmt.Println(string(stdout))
		return d, nil
	}
}

func packageHelmChartWithSubcharts() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		// Collect values to be set via the CLI.
		cliValues := []string{
			fmt.Sprintf("--set=image.tag=%s", d.ctxt.GitCommitSHA),
		}

		d.logger.Infof("Adding dependencies from subrepos into the charts/ directory ...")
		// Find subcharts
		chartsDir := filepath.Join(d.opts.chartDir, "charts")
		if _, err := os.Stat(chartsDir); os.IsNotExist(err) {
			err = os.Mkdir(chartsDir, 0755)
			if err != nil {
				log.Fatalf("could not create %s: %s", chartsDir, err)
			}
		}
		for _, r := range d.subrepos {
			subrepo := filepath.Join(pipelinectxt.SubreposPath, r.Name())
			subchart := filepath.Join(subrepo, d.opts.chartDir)
			if _, err := os.Stat(subchart); os.IsNotExist(err) {
				fmt.Printf("no chart in %s\n", r.Name())
				continue
			}
			gitCommitSHA, err := getTrimmedFileContent(filepath.Join(subrepo, ".ods", "git-commit-sha"))
			if err != nil {
				log.Fatal(err)
			}
			hc, err := getHelmChart(filepath.Join(subchart, "Chart.yaml"))
			if err != nil {
				log.Fatal(err)
			}
			cliValues = append(cliValues, fmt.Sprintf("--set=%s.image.tag=%s", hc.Name, gitCommitSHA))
			if d.releaseName() == d.ctxt.Component {
				cliValues = append(cliValues, fmt.Sprintf("--set=%s.fullnameOverride=%s", hc.Name, hc.Name))
			}
			helmArchive, err := packageHelmChart(subchart, d.ctxt.Version, gitCommitSHA, d.opts.debug)
			if err != nil {
				log.Fatal(err)
			}
			helmArchiveName := filepath.Base(helmArchive)
			fmt.Printf("copying %s into %s\n", helmArchiveName, chartsDir)
			err = file.Copy(helmArchive, filepath.Join(chartsDir, helmArchiveName))
			if err != nil {
				log.Fatal(err)
			}
		}

		subcharts, err := os.ReadDir(chartsDir)
		if err != nil {
			log.Fatal(err)
		}
		if len(subcharts) > 0 {
			fmt.Printf("Contents of %s:\n", chartsDir)
			for _, sc := range subcharts {
				fmt.Println(sc.Name())
			}
		}

		d.logger.Infof("Packaging Helm chart ...")
		helmArchive, err := packageHelmChart(d.opts.chartDir, d.ctxt.Version, d.ctxt.GitCommitSHA, d.opts.debug)
		if err != nil {
			log.Fatal(err)
		}
		d.helmArchive = helmArchive
		d.cliValues = cliValues
		return d, nil
	}
}

func collectValuesFiles() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		d.logger.Infof("Collecting Helm values files ...")
		valuesFiles := []string{}
		valuesFilesCandidates := []string{
			fmt.Sprintf("%s/secrets.yaml", d.opts.chartDir), // equivalent values.yaml is added automatically by Helm
			fmt.Sprintf("%s/values.%s.yaml", d.opts.chartDir, d.targetConfig().Stage),
			fmt.Sprintf("%s/secrets.%s.yaml", d.opts.chartDir, d.targetConfig().Stage),
		}
		if string(d.targetConfig().Stage) != d.targetConfig().Name {
			valuesFilesCandidates = append(
				valuesFilesCandidates,
				fmt.Sprintf("%s/values.%s.yaml", d.opts.chartDir, d.targetConfig().Name),
				fmt.Sprintf("%s/secrets.%s.yaml", d.opts.chartDir, d.targetConfig().Name),
			)
		}
		for _, vfc := range valuesFilesCandidates {
			if _, err := os.Stat(vfc); os.IsNotExist(err) {
				fmt.Printf("%s is not present, skipping.\n", vfc)
			} else {
				fmt.Printf("%s is present, adding.\n", vfc)
				valuesFiles = append(valuesFiles, vfc)
			}
		}
		d.valuesFiles = valuesFiles
		return d, nil
	}
}

func importAgeKey() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		if len(d.opts.ageKeySecret) == 0 {
			d.logger.Infof("Skipping import of age key for helm-secrets as parameter is not set ...")
		} else {
			d.logger.Infof("Storing age key for helm-secrets ...")
			secret, err := d.clientset.CoreV1().Secrets(d.ctxt.Namespace).Get(
				context.TODO(), d.opts.ageKeySecret, metav1.GetOptions{},
			)
			if err != nil {
				d.logger.Infof("No secret %s found, skipping.", d.opts.ageKeySecret)
			} else {
				stderr, err := storeAgeKey(secret, d.opts.ageKeySecretField)
				if err != nil {
					fmt.Println(string(stderr))
					log.Fatal(err)
				}
				d.logger.Infof("Age key secret %s stored.", d.opts.ageKeySecret)
			}
		}
		return d, nil
	}
}

func diffHelmRelease() DeployStep {
	return func(d *deployHelm) (*deployHelm, error) {
		fmt.Printf("Diffing Helm release against %s...\n", d.helmArchive)
		helmDiffArgs, err := d.assembleHelmDiffArgs()
		if err != nil {
			log.Fatal("assemble helm diff args: ", err)
		}
		printlnSafeHelmCmd(helmDiffArgs, os.Stdout)
		// helm-dff stderr contains confusing text about "errors" when drift is
		// detected, therefore we want to collect and polish it before we print it.
		// helm-diff stdout needs to be written into a buffer so that we can both
		// print it and store it later as a deployment artifact.
		var diffStdout, diffStderr bytes.Buffer
		inSync, err := helmDiff(d.helmBin, helmDiffArgs, &diffStdout, &diffStderr)
		fmt.Print(diffStdout.String())
		fmt.Print(cleanHelmDiffOutput(diffStderr.String()))
		if err != nil {
			log.Fatal("helm diff: ", err)
		}
		if inSync {
			fmt.Println("No diff detected, skipping helm upgrade.")
			os.Exit(0)
		}

		err = writeDeploymentArtifact(diffStdout.Bytes(), "diff", d.opts.chartDir, d.targetConfig().Name)
		if err != nil {
			log.Fatal("write diff artifact: ", err)
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
		err = writeDeploymentArtifact(upgradeStdoutBuf.Bytes(), "release", d.opts.chartDir, d.targetConfig().Name)
		if err != nil {
			return d, fmt.Errorf("write release artifact: %w", err)
		}
		return d, nil
	}
}
