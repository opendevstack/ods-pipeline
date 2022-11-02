package main

import (
	"os"

	"github.com/opendevstack/pipeline/internal/command"
)

func (p *packageImage) generateImageSBOM() error {
	// TODO: support multiple formats?
	// TODO: support more args?
	// TODO: how to name the result?
	// TODO: how to ref the image?
	return command.Run(
		p.trivyBin, []string{"image", "--format=spdx", "--output=result.spdx", "alpine3.15"}, []string{},
		os.Stdout, os.Stderr,
	)
}
