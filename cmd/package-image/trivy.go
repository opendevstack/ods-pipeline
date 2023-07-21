package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/shlex"
	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
)

const (
	trivyBin     = "trivy"
	trivyWorkdir = "/tmp"
)

func (p *packageImage) generateImageSBOM() error {
	// settle for one format and name until we have use cases for multiple formats (we use spdx format).
	// trivy support --formats:  table, json, sarif, template, cyclonedx, spdx, spdx-json, github, cosign-vuln (default "table")
	// more args for experimentation via extra args
	extraArgs, err := shlex.Split(p.opts.trivySBOMExtraArgs)
	if err != nil {
		p.logger.Errorf("could not parse extra args (%s): %s", p.opts.trivySBOMExtraArgs, err)
	}
	sbomFilename := fmt.Sprintf("%s.%s", p.imageNameNoSha(), pipelinectxt.SBOMsFormat)
	p.sbomFile = filepath.Join(trivyWorkdir, sbomFilename)
	args := []string{
		"image",
		fmt.Sprintf("--format=%s", pipelinectxt.SBOMsFormat),
		fmt.Sprintf("--input=%s", filepath.Join(buildahWorkdir, p.imageNameNoSha())),
		fmt.Sprintf("--output=%s", p.sbomFile),
	}
	if p.opts.debug {
		args = append(args, "--debug=true")
	}
	args = append(args, extraArgs...)
	return command.RunInDir(trivyBin, args, []string{}, trivyWorkdir, os.Stdout, os.Stderr)
}
