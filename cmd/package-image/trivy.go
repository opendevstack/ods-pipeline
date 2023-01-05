package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

const (
	trivyBin = "trivy"
)

func (p *packageImage) generateImageSBOM() error {
	// settle for one format and name until we have use cases for multiple formats (we use spdx format).
	// trivy support --formats:  table, json, sarif, template, cyclonedx, spdx, spdx-json, github, cosign-vuln (default "table")
	// more args for experimentation via extra args
	extraArgs, err := shlex.Split(p.opts.trivySBOMExtraArgs)
	if err != nil {
		p.logger.Errorf("could not parse extra args (%s): %s", p.opts.trivySBOMExtraArgs, err)
	}
	sbomFilename := fmt.Sprintf("%s.%s", p.image.Name, pipelinectxt.SBOMsFormat)
	sbomFile := filepath.Join(p.opts.checkoutDir, sbomFilename)
	args := []string{
		"image",
		fmt.Sprintf("--format=%s", pipelinectxt.SBOMsFormat),
		fmt.Sprintf("--input=%s", filepath.Join(p.opts.checkoutDir, p.image.Name)),
		fmt.Sprintf("--output=%s", sbomFile),
	}
	if p.opts.debug {
		args = append(args, "--debug=true")
	}
	args = append(args, extraArgs...)
	return command.Run(
		trivyBin, args, []string{},
		os.Stdout, os.Stderr,
	)
}
