package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/image"
)

func (p *packageImage) skopeoTag(idt *image.IdentityWithTag, outWriter, errWriter io.Writer) error {
	imageRef := idt.ImageRefWithSha(p.opts.registry)
	p.logger.Infof("Tagging image %s with %s", imageRef, idt.Tag)
	tlsVerify := p.opts.tlsVerify
	// TLS verification of the KinD registry is not possible at the moment as
	// requests error out with "server gave HTTP response to HTTPS client".
	if strings.HasPrefix(p.opts.registry, "kind-registry.kind") {
		tlsVerify = false
	}
	args := []string{
		"copy",
		fmt.Sprintf("--src-tls-verify=%v", tlsVerify),
		fmt.Sprintf("--dest-tls-verify=%v", tlsVerify),
	}
	if tlsVerify {
		args = append(args,
			fmt.Sprintf("--src-cert-dir=%v", p.opts.certDir),
			fmt.Sprintf("--dest-cert-dir=%v", p.opts.certDir))
	}
	if p.opts.debug {
		args = append(args, "--debug")
	}
	source := fmt.Sprintf("docker://%s", imageRef)
	destination := fmt.Sprintf("docker://%s", idt.ImageRef(p.opts.registry))

	args = append(args, source, destination)

	p.logger.Infof("skopeo copy %s %s", source, destination)
	err := command.Run("skopeo", args, []string{}, outWriter, errWriter)
	if err != nil {
		return fmt.Errorf("skopeo copy %s to %s: %w", source, destination, err)
	}
	return nil
}
