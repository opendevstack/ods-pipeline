package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/artifact"
)

func (d *deployHelm) copyImage(imageArtifact artifact.Image, destRegistryToken string, outWriter, errWriter io.Writer) error {
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
	if d.targetConfig.RegistryHost != "" && d.targetConfig.RegistryTLSVerify != nil {
		destRegistryTLSVerify = *d.targetConfig.RegistryTLSVerify
	}
	destImageURL := getImageDestURL(d.targetConfig.RegistryHost, d.releaseNamespace, imageArtifact)
	d.logger.Infof("Source image: %s", srcImageURL)
	d.logger.Infof("Destination image: %s", destImageURL)
	// TODO: for QA and PROD we want to ensure that the SHA recorded in Nexus
	// matches the SHA referenced by the Git commit tag.
	args := []string{
		"copy",
		fmt.Sprintf("--src-tls-verify=%v", srcRegistryTLSVerify),
		fmt.Sprintf("--dest-tls-verify=%v", destRegistryTLSVerify),
	}
	if srcRegistryTLSVerify {
		args = append(args, fmt.Sprintf("--src-cert-dir=%v", d.opts.certDir))
	}
	if destRegistryTLSVerify {
		args = append(args, fmt.Sprintf("--dest-cert-dir=%v", d.opts.certDir))
	}
	if destRegistryToken != "" {
		args = append(args, "--dest-registry-token", destRegistryToken)
	}
	if d.opts.debug {
		args = append(args, "--debug")
	}
	args = append(
		args, fmt.Sprintf("docker://%s", srcImageURL), fmt.Sprintf("docker://%s", destImageURL),
	)
	err := command.Run("skopeo", args, []string{}, outWriter, errWriter)
	if err != nil {
		return fmt.Errorf("skopeo copy %s: %w", srcImageURL, err)
	}
	return nil
}

func getImageDestURL(registryHost, releaseNamespace string, imageArtifact artifact.Image) string {
	if registryHost != "" {
		return fmt.Sprintf("%s/%s/%s:%s", registryHost, releaseNamespace, imageArtifact.Name, imageArtifact.Tag)
	} else {
		return strings.Replace(imageArtifact.Image, "/"+imageArtifact.Repository+"/", "/"+releaseNamespace+"/", -1)
	}
}
