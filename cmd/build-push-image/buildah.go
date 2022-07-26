package main

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/command"
)

// buildahBuild builds a local image using the Dockerfile and context directory
// given in opts, tagging the resulting image with given tag.
func buildahBuild(opts options, tag string, outWriter, errWriter io.Writer) error {
	extraArgs, err := shlex.Split(opts.buildahBuildExtraArgs)
	if err != nil {
		return fmt.Errorf("parse extra args (%s): %w", opts.buildahBuildExtraArgs, err)
	}

	args := []string{
		fmt.Sprintf("--storage-driver=%s", opts.storageDriver),
		"bud",
		fmt.Sprintf("--format=%s", opts.format),
		fmt.Sprintf("--tls-verify=%v", opts.tlsVerify),
		fmt.Sprintf("--cert-dir=%s", opts.certDir),
		"--no-cache",
		fmt.Sprintf("--file=%s", opts.dockerfile),
		fmt.Sprintf("--tag=%s", tag),
	}
	args = append(args, extraArgs...)
	nexusArgs, err := nexusBuildArgs(opts)
	if err != nil {
		return fmt.Errorf("add nexus build args: %w", err)
	}
	args = append(args, nexusArgs...)

	if opts.debug {
		args = append(args, "--log-level=debug")
	}
	_, err = command.RunWithStreamingOutput(
		"buildah",
		append(args, opts.contextDir),
		[]string{},
		outWriter, errWriter,
		-1, // no special exit code handling
	)
	return err
}

// buildahPush pushes a local image to the given imageRef.
func buildahPush(opts options, workingDir, imageRef string, outWriter, errWriter io.Writer) error {
	extraArgs, err := shlex.Split(opts.buildahPushExtraArgs)
	if err != nil {
		log.Printf("could not parse extra args (%s): %s", opts.buildahPushExtraArgs, err)
	}
	args := []string{
		fmt.Sprintf("--storage-driver=%s", opts.storageDriver),
		"push",
		fmt.Sprintf("--tls-verify=%v", opts.tlsVerify),
		fmt.Sprintf("--cert-dir=%s", opts.certDir),
		fmt.Sprintf("--digestfile=%s", filepath.Join(workingDir, "image-digest")),
	}
	args = append(args, extraArgs...)
	if opts.debug {
		args = append(args, "--log-level=debug")
	}
	_, err = command.RunWithStreamingOutput(
		"buildah",
		append(args, imageRef, fmt.Sprintf("docker://%s", imageRef)),
		[]string{},
		outWriter, errWriter,
		-1, // no special exit code handling
	)
	return err
}
