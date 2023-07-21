package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opendevstack/ods-pipeline/internal/image"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	"golang.org/x/exp/slog"
)

type AquaScanStep func(d *aquaScan) (*aquaScan, error)

func (s *aquaScan) runSteps(steps ...AquaScanStep) error {
	var skip *skipRemainingSteps
	var err error
	for _, step := range steps {
		s, err = step(s)
		if err != nil {
			if errors.As(err, &skip) {
				slog.Info(err.Error())
				return nil
			}
			return err
		}
	}
	return nil
}

// setupContext creates and ODS context.
func setupContext() AquaScanStep {
	return func(s *aquaScan) (*aquaScan, error) {
		ctxt := &pipelinectxt.ODSContext{}
		err := ctxt.ReadCache(s.opts.checkoutDir)
		if err != nil {
			return s, fmt.Errorf("read cache: %w", err)
		}
		s.ctxt = ctxt

		return s, nil
	}
}

func setImageId() AquaScanStep {
	return func(p *aquaScan) (*aquaScan, error) {
		p.imageId = image.CreateImageIdentity(p.ctxt, p.opts.imageNamespace, p.opts.imageStream)
		return p, nil
	}
}

func skipIfScanArtifactsExist() AquaScanStep {
	return func(s *aquaScan) (*aquaScan, error) {
		if ok := aquaReportsExist(pipelinectxt.AquaScansPath, s.imageId); ok {
			return s, &skipRemainingSteps{fmt.Sprintf("aqua scan artifact exists already for %s", s.imageId.ImageStream)}
		}
		return s, nil
	}
}

// scanImagesWithAqua runs the Aqua scanner over each image artifact.
func scanImagesWithAqua() AquaScanStep {
	return func(s *aquaScan) (*aquaScan, error) {
		slog.Info("Scanning image with Aqua scanner ...")
		aquaImage := s.imageId.NamespaceStreamSha()
		htmlReportFile := filepath.Join(s.opts.checkoutDir, htmlReportFilename(s.imageId))
		jsonReportFile := filepath.Join(s.opts.checkoutDir, jsonReportFilename(s.imageId))
		scanArgs := aquaAssembleScanArgs(s.opts, aquaImage, htmlReportFile, jsonReportFile)
		scanSuccessful, err := runScan(aquasecBin, scanArgs, os.Stdout, os.Stderr)
		if err != nil {
			return s, fmt.Errorf("aqua scan: %w", err)
		}

		if !scanSuccessful && s.opts.aquasecGate {
			return s, errors.New("stopping build as successful Aqua scan is required")
		}

		asu, err := aquaScanURL(s.opts, aquaImage)
		if err != nil {
			return s, fmt.Errorf("aqua scan URL: %w", err)
		}
		slog.Info("Aqua vulnerability report is at " + asu)

		err = copyAquaReportsToArtifacts(htmlReportFile, jsonReportFile)
		if err != nil {
			return s, err
		}

		slog.Info("Creating Bitbucket code insight report ...")
		err = createBitbucketInsightReport(s.opts, asu, scanSuccessful, s.ctxt)
		if err != nil {
			return s, err
		}
		return s, nil
	}
}
