package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/command"
)

const (
	buildahBin     = "buildah"
	buildahWorkdir = "/tmp"
)

// buildahBuild builds a local image using the Dockerfile and context directory
// given in opts, tagging the resulting image with given tag.
func (p *packageImage) buildahBuild(outWriter, errWriter io.Writer) error {
	args, err := p.buildahBuildArgs(p.imageRef())
	if err != nil {
		return fmt.Errorf("assemble build args: %w", err)
	}
	return command.RunInDir(buildahBin, args, []string{}, buildahWorkdir, outWriter, errWriter)
}

// buildahPush pushes a local image to a OCI formatted directory for trivy image scans.
func (p *packageImage) buildahPushTar(outWriter, errWriter io.Writer) error {
	args := []string{
		fmt.Sprintf("--storage-driver=%s", p.opts.storageDriver),
		"push",
		fmt.Sprintf("--digestfile=%s", tektonResultsImageDigestFile),
	}
	if p.opts.debug {
		args = append(args, "--log-level=debug")
	}
	args = append(args, p.imageRef(), fmt.Sprintf("oci:%s", filepath.Join(p.opts.checkoutDir, p.imageName())))
	return command.RunInDir(buildahBin, args, []string{}, buildahWorkdir, outWriter, errWriter)
}

// buildahPush pushes a local image to the given imageRef.
func (p *packageImage) buildahPush(outWriter, errWriter io.Writer) error {
	opts := p.opts
	extraArgs, err := shlex.Split(opts.buildahPushExtraArgs)
	if err != nil {
		log.Printf("could not parse extra args (%s): %s", opts.buildahPushExtraArgs, err)
	}
	args := []string{
		fmt.Sprintf("--storage-driver=%s", opts.storageDriver),
		"push",
		fmt.Sprintf("--tls-verify=%v", opts.tlsVerify),
		fmt.Sprintf("--cert-dir=%s", opts.certDir),
	}
	args = append(args, extraArgs...)
	if opts.debug {
		args = append(args, "--log-level=debug")
	}

	source := p.imageId.ImageRefWithSha(opts.registry)
	destination := fmt.Sprintf("docker://%s", source)
	log.Printf("buildah push %s %s", source, destination)
	args = append(args, source, destination)
	return command.RunInDir(buildahBin, args, []string{}, buildahWorkdir, outWriter, errWriter)
}

// buildahBuildArgs assembles the args to be passed to buildah based on
// given options and tag.
func (p *packageImage) buildahBuildArgs(tag string) ([]string, error) {
	if tag == "" {
		return nil, errors.New("tag must not be empty")
	}
	opts := p.opts
	extraArgs, err := shlex.Split(opts.buildahBuildExtraArgs)
	if err != nil {
		return nil, fmt.Errorf("parse extra args (%s): %w", opts.buildahBuildExtraArgs, err)
	}

	absDir, err := filepath.Abs(opts.checkoutDir)
	if err != nil {
		return nil, fmt.Errorf("abs dir: %w", err)
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
	nexusArgs, err := p.nexusBuildArgs()
	if err != nil {
		return nil, fmt.Errorf("add nexus build args: %w", err)
	}
	args = append(args, nexusArgs...)

	if opts.debug {
		args = append(args, "--log-level=debug")
	}
	return append(args, filepath.Join(absDir, opts.contextDir)), nil
}

// nexusBuildArgs computes --build-arg parameters so that the Dockerfile
// can access nexus as determined by the options nexus related
// parameters.
func (p *packageImage) nexusBuildArgs() ([]string, error) {
	args := []string{}
	opts := p.opts
	if strings.TrimSpace(opts.nexusURL) != "" {
		nexusUrl, err := url.Parse(opts.nexusURL)
		if err != nil {
			return nil, fmt.Errorf("could not parse nexus url (%s): %w", opts.nexusURL, err)
		}
		if nexusUrl.Host == "" {
			return nil, fmt.Errorf("could not get host in nexus url (%s)", opts.nexusURL)
		}
		if opts.nexusUsername != "" {
			if opts.nexusPassword == "" {
				nexusUrl.User = url.User(opts.nexusUsername)
			} else {
				nexusUrl.User = url.UserPassword(opts.nexusUsername, opts.nexusPassword)
			}
		}
		nexusAuth := nexusUrl.User.String() // this is encoded as needed.
		a := strings.SplitN(nexusAuth, ":", 2)
		unEscaped := ""
		pwEscaped := ""
		if len(a) > 0 {
			unEscaped = a[0]
		}
		if len(a) > 1 {
			pwEscaped = a[1]
		}
		args = []string{
			fmt.Sprintf("--build-arg=nexusUrl=%s", opts.nexusURL),
			fmt.Sprintf("--build-arg=nexusUsername=%s", unEscaped),
			fmt.Sprintf("--build-arg=nexusPassword=%s", pwEscaped),
			fmt.Sprintf("--build-arg=nexusHost=%s", nexusUrl.Host),
		}
		args = append(args, fmt.Sprintf("--build-arg=nexusAuth=%s", nexusAuth))
		if nexusAuth != "" {
			args = append(args,
				fmt.Sprintf("--build-arg=nexusUrlWithAuth=%s://%s@%s", nexusUrl.Scheme, nexusAuth, nexusUrl.Host))
		} else {
			args = append(args,
				fmt.Sprintf("--build-arg=nexusUrlWithAuth=%s", opts.nexusURL))
		}
	}
	return args, nil
}
