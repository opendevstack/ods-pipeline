package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

type PackageStep func(d *packageImage) (*packageImage, error)

func (d *packageImage) runSteps(steps ...PackageStep) error {
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

func setupContext() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		ctxt := &pipelinectxt.ODSContext{}
		err := ctxt.ReadCache(p.opts.checkoutDir)
		if err != nil {
			return p, fmt.Errorf("read cache: %w", err)
		}
		p.ctxt = ctxt

		if p.opts.debug {
			if err := directory.ListFiles(p.opts.certDir, os.Stdout); err != nil {
				p.logger.Errorf(err.Error())
			}
		}

		// TLS verification of the KinD registry is not possible at the moment as
		// requests error out with "server gave HTTP response to HTTPS client".
		if strings.HasPrefix(p.opts.registry, "kind-registry.kind") {
			p.opts.tlsVerify = false
		}

		return p, nil
	}
}

func setImageName() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		imageNamespace := p.opts.imageNamespace
		if len(imageNamespace) == 0 {
			imageNamespace = p.ctxt.Namespace
		}
		imageStream := p.opts.imageStream
		if len(imageStream) == 0 {
			imageStream = p.ctxt.Component
		}
		imageName := fmt.Sprintf("%s:%s", imageStream, p.ctxt.GitCommitSHA)
		p.image.Name = imageStream
		p.image.Repository = imageNamespace
		p.image.Tag = p.ctxt.GitCommitSHA
		p.image.Ref = fmt.Sprintf(
			"%s/%s/%s",
			p.opts.registry, imageNamespace, imageName,
		)
		return p, nil
	}
}

func skipIfImageExists() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		fmt.Printf("Checking if image %s exists already ...\n", p.image.ImageName())
		imageDigest, err := getImageDigestFromRegistry(p.image.Ref, p.opts)
		if err == nil {
			return p, &skipRemainingSteps{"image exists already"}
		}
		p.image.Digest = imageDigest
		return p, nil
	}
}

func buildImageAndGenerateTar() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		fmt.Printf("Building image %s ...\n", p.image.ImageName())
		err := p.buildahBuild(os.Stdout, os.Stderr)
		if err != nil {
			return p, fmt.Errorf("buildah bud: %w", err)
		}
		fmt.Printf("Creating local tar folder for image %s ...\n", p.image.ImageName())
		err = p.buildahPushTar(os.Stdout, os.Stderr)
		if err != nil {
			return p, fmt.Errorf("buildah push tar: %w", err)
		}
		d, err := getImageDigestFromFile(p.opts.checkoutDir)
		if err != nil {
			return p, err
		}
		p.image.Digest = d
		return p, nil
	}
}

func generateSBOM() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		fmt.Println("Generating image SBOM with trivy scanner ...")
		err := p.generateImageSBOM()
		if err != nil {
			return p, fmt.Errorf("generate SBOM: %w", err)
		}
		return p, nil
	}
}

func scanImageWithAqua() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		if aquasecInstalled() {
			fmt.Println("Scanning image with Aqua scanner ...")
			aquaImage := p.image.ImageName()
			htmlReportFile := filepath.Join(p.opts.checkoutDir, "report.html")
			jsonReportFile := filepath.Join(p.opts.checkoutDir, "report.json")
			scanArgs := aquaAssembleScanArgs(p.opts, aquaImage, htmlReportFile, jsonReportFile)
			scanSuccessful, err := aquaScan(aquasecBin, scanArgs, os.Stdout, os.Stderr)
			if err != nil {
				return p, fmt.Errorf("aqua scan: %w", err)
			}

			if !scanSuccessful && p.opts.aquasecGate {
				return p, errors.New("stopping build as successful Aqua scan is required")
			}

			asu, err := aquaScanURL(p.opts, aquaImage)
			if err != nil {
				return p, fmt.Errorf("aqua scan URL: %w", err)
			}
			fmt.Printf("Aqua vulnerability report is at %s ...\n", asu)

			err = copyAquaReportsToArtifacts(htmlReportFile, jsonReportFile)
			if err != nil {
				return p, err
			}

			fmt.Println("Creating Bitbucket code insight report ...")
			err = createBitbucketInsightReport(p.opts, asu, scanSuccessful, p.ctxt)
			if err != nil {
				return p, err
			}
		} else {
			fmt.Println("Aqua is not configured, image will not be scanned for vulnerabilities.")
		}
		return p, nil
	}
}

func pushImage() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		fmt.Printf("Pushing image %s ...\n", p.image.ImageName())
		err := p.buildahPush(os.Stdout, os.Stderr)
		if err != nil {
			return p, fmt.Errorf("buildah push: %w", err)
		}
		return p, nil
	}
}

func storeArtifact() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		err := writeImageDigestToResults(p.image.Digest)
		if err != nil {
			return p, err
		}

		fmt.Println("Writing image artifact ...")
		imageArtifactFilename := fmt.Sprintf("%s.json", p.image.Name)
		err = pipelinectxt.WriteJsonArtifact(p.image, pipelinectxt.ImageDigestsPath, imageArtifactFilename)
		if err != nil {
			return p, err
		}

		fmt.Println("Writing SBOM artifact ...")
		sbomFile := filepath.Join(p.opts.checkoutDir, pipelinectxt.SbomsFilename)
		if _, err := os.Stat(sbomFile); err == nil {
			err := pipelinectxt.CopyArtifact(sbomFile, pipelinectxt.SbomsPath)
			if err != nil {
				return p, fmt.Errorf("copying SBOM report to artifacts failed: %w", err)
			}
		}

		return p, nil
	}
}
