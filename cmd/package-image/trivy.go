package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func (p *packageImage) generateImageSBOM() error {
	// settle for one format and name until we have use cases for multiple formats.
	// trivy support --formats:  table, json, sarif, template, cyclonedx, spdx, spdx-json, github, cosign-vuln (default "table")
	// more args for experimentation via extra args
	extraArgs, err := shlex.Split(p.opts.trivySBOMExtraArgs)
	if err != nil {
		p.logger.Errorf("could not parse extra args (%s): %s", p.opts.trivySBOMExtraArgs, err)
	}
	sbomFile := filepath.Join(p.opts.checkoutDir, pipelinectxt.SbomsFilename)
	args := []string{
		"image",
		fmt.Sprintf("--input=%s.tar", p.image.Name),
		"--format=spdx-json",
		fmt.Sprintf("--output=%s", sbomFile),
	}
	args = append(args, extraArgs...)

	return command.Run(
		p.trivyBin, args, []string{},
		os.Stdout, os.Stderr,
	)
}
