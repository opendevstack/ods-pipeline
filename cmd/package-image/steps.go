package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/pkg/artifact"
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

func setExtraTags() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		extraTagsSpecified, err := shlex.Split(p.opts.extraTags)
		if err != nil {
			return p, fmt.Errorf("parse extra tags (%s): %w", p.opts.extraTags, err)
		}
		p.parsedExtraTags = extraTagsSpecified
		return p, nil
	}
}

func setImageId() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		p.imageId = createImageIdentity(p.ctxt, &p.opts)
		return p, nil
	}
}

// skipIfImageArtifactExists informs to skip next steps if ODS image artifact is already in place.
// In future we might want to check all the expected artifacts, that must exist to do skip properly.
func skipIfImageArtifactExists() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		fmt.Printf("Checking if image artifact for %s exists already ...\n", p.imageName())
		err := imageArtifactExists(p)
		if err == nil {
			return p, &skipRemainingSteps{"image artifact exists already"}
		}
		return p, nil
	}
}

func buildImageAndGenerateTar() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		fmt.Printf("Building image %s ...\n", p.imageName())
		err := p.buildahBuild(os.Stdout, os.Stderr)
		if err != nil {
			return p, fmt.Errorf("buildah bud: %w", err)
		}
		fmt.Printf("Creating local tar folder for image %s ...\n", p.imageName())
		err = p.buildahPushTar(os.Stdout, os.Stderr)
		if err != nil {
			return p, fmt.Errorf("buildah push tar: %w", err)
		}
		d, err := getImageDigestFromFile(p.opts.checkoutDir)
		if err != nil {
			return p, err
		}
		p.imageDigest = d
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
			aquaImage := fmt.Sprintf("%s/%s", p.imageId.nsStreamSha())
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
		fmt.Printf("Pushing image %s ...\n", p.imageName())
		err := p.buildahPush(os.Stdout, os.Stderr)
		if err != nil {
			return p, fmt.Errorf("buildah push: %w", err)
		}
		return p, nil
	}
}

func storeArtifact() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		imageId := p.imageId
		imageShaTag := imageId.shaTag()
		image := artifact.Image{
			Ref:        imageShaTag.imageRef(p.opts.registry),
			Registry:   p.opts.registry,
			Repository: imageId.ImageNamespace,
			Name:       imageId.ImageStream,
			Tag:        imageId.GitCommitSHA,
			Digest:     p.imageDigest,
		}
		err := writeImageDigestToResults(image.Digest)
		if err != nil {
			return p, err
		}

		fmt.Println("Writing image artifact ...")
		imageArtifactFilename := fmt.Sprintf("%s.json", image.Name)
		err = pipelinectxt.WriteJsonArtifact(image, pipelinectxt.ImageDigestsPath, imageArtifactFilename)
		if err != nil {
			return p, err
		}

		fmt.Println("Writing SBOM artifact ...")
		sbomFilename := fmt.Sprintf("%s.%s", image.Name, pipelinectxt.SBOMsFormat)
		sbomFile := filepath.Join(p.opts.checkoutDir, sbomFilename)
		err = pipelinectxt.CopyArtifact(sbomFile, pipelinectxt.SBOMsPath)
		if err != nil {
			return p, fmt.Errorf("copying SBOM report to artifacts failed: %w", err)
		}

		return p, nil
	}
}

func processExtraTags() PackageStep {
	return func(p *packageImage) (*packageImage, error) {
		if len(p.parsedExtraTags) > 0 {
			log.Printf("Processing extra tags missing in registry: %+q", p.parsedExtraTags)
			missingTags, err := p.skopeoMissingTags()
			if err != nil {
				return p, fmt.Errorf("Could not determine missing tags:", err)
			}
			if len(missingTags) == 0 {
				log.Print("No missing extra tags found.")
				return p, nil
			}
			log.Printf("pushing missing extra tags: %+q", missingTags)
			for _, missingTag := range missingTags {
				imageExtraTag := p.imageId.tag(missingTag)
				err = p.skopeoTag(&imageExtraTag, os.Stdout, os.Stderr)
				if err != nil {
					log.Fatal("skopeo push failed: ", err)
				}
			}
			p.logger.Infof("Writing image artifacts for all extra tags ...")
			for _, extraTag := range p.parsedExtraTags {
				imageExtraTag := p.imageId.tag(extraTag)
				ia := artifact.Image{
					Ref:        imageExtraTag.imageRef(p.opts.registry),
					Registry:   p.opts.registry,
					Repository: imageExtraTag.ImageIdentity.ImageNamespace,
					Name:       imageExtraTag.ImageIdentity.ImageStream,
					Tag:        imageExtraTag.Tag,
					Digest:     p.imageDigest,
				}
				imageArtifactFilename := fmt.Sprintf("%s-%s.json", imageExtraTag.ImageIdentity.ImageStream, imageExtraTag.Tag)
				err = pipelinectxt.WriteJsonArtifact(ia, pipelinectxt.ImageDigestsPath, imageArtifactFilename)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		return p, nil
	}
}
