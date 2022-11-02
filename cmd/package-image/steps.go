package main

import (
	"errors"
	"fmt"
	"log"
	"os"
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
				log.Fatal(err)
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
		// 	imageNamespace := opts.imageNamespace
		// 	if len(imageNamespace) == 0 {
		// 		imageNamespace = ctxt.Namespace
		// 	}
		// 	imageStream := opts.imageStream
		// 	if len(imageStream) == 0 {
		// 		imageStream = ctxt.Component
		// 	}
		// 	imageName := fmt.Sprintf("%s:%s", imageStream, ctxt.GitCommitSHA)
		// 	imageRef := fmt.Sprintf(
		// 		"%s/%s/%s",
		// 		opts.registry, imageNamespace, imageName,
		// 	)
		return p, nil
	}
}

func skipIfImageExists() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		// 	fmt.Printf("Checking if image %s exists already ...\n", imageName)
		// 	imageDigest, err := getImageDigestFromRegistry(imageRef, opts)
		// 	if err == nil {
		// 		fmt.Println("Image exists already.")
		// 	} else {
		return p, nil
	}
}

func buildImage() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		// 		fmt.Printf("Building image %s ...\n", imageName)
		// 		err := buildahBuild(opts, imageRef, os.Stdout, os.Stderr)
		// 		if err != nil {
		// 			log.Fatal("buildah bud: ", err)
		// 		}

		// 		d, err := getImageDigestFromFile(workingDir)
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}
		// 		imageDigest = d
		return p, nil
	}
}

func generateSBOM() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		// $ trivy image --format spdx --output result.spdx alpine:3.15
		err := p.generateImageSBOM()
		if err != nil {
			return p, fmt.Errorf("generate SBOM: %w", err)
		}
		return p, nil
	}
}

func scanImage() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		// 		if aquasecInstalled() {
		// 			fmt.Println("Scanning image with Aqua scanner ...")
		// 			aquaImage := fmt.Sprintf("%s/%s", imageNamespace, imageName)
		// 			htmlReportFile := filepath.Join(workingDir, "report.html")
		// 			jsonReportFile := filepath.Join(workingDir, "report.json")
		// 			scanArgs := aquaAssembleScanArgs(opts, aquaImage, htmlReportFile, jsonReportFile)
		// 			scanSuccessful, err := aquaScan(aquasecBin, scanArgs, os.Stdout, os.Stderr)
		// 			if err != nil {
		// 				log.Fatal("aqua scan: ", err)
		// 			}

		// 			if !scanSuccessful && opts.aquasecGate {
		// 				log.Fatalln("Stopping build as successful Aqua scan is required")
		// 			}

		// 			asu, err := aquaScanURL(opts, aquaImage)
		// 			if err != nil {
		// 				log.Fatal("aqua scan URL:", err)
		// 			}
		// 			fmt.Printf("Aqua vulnerability report is at %s ...\n", asu)

		// 			err = copyAquaReportsToArtifacts(htmlReportFile, jsonReportFile)
		// 			if err != nil {
		// 				log.Fatal(err)
		// 			}

		// 			fmt.Println("Creating Bitbucket code insight report ...")
		// 			err = createBitbucketInsightReport(opts, asu, scanSuccessful, ctxt)
		// 			if err != nil {
		// 				log.Fatal(err)
		// 			}
		// 		} else {
		// 			fmt.Println("Aqua is not configured, image will not be scanned for vulnerabilities.")
		// 		}
		return p, nil
	}
}

func pushImage() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		// 		fmt.Printf("Pushing image %s ...\n", imageRef)
		// 		err = buildahPush(opts, workingDir, imageRef, os.Stdout, os.Stderr)
		// 		if err != nil {
		// 			log.Fatal("buildah push: ", err)
		// 		}
		return p, nil
	}
}

func storeArtifact() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		// 	err = writeImageDigestToResults(imageDigest)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	fmt.Println("Writing image artifact ...")
		// 	ia := artifact.Image{
		// 		Image:      imageRef,
		// 		Registry:   opts.registry,
		// 		Repository: imageNamespace,
		// 		Name:       imageStream,
		// 		Tag:        ctxt.GitCommitSHA,
		// 		Digest:     imageDigest,
		// 	}
		// 	imageArtifactFilename := fmt.Sprintf("%s.json", imageStream)
		// 	err = pipelinectxt.WriteJsonArtifact(ia, pipelinectxt.ImageDigestsPath, imageArtifactFilename)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		return p, nil
	}
}
