package pipelinectxt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/file"
)

const (
	ArtifactsDir      = "artifacts"
	ArtifactsPath     = BaseDir + "/" + ArtifactsDir
	PipelineRunsDir   = "pipeline-runs"
	PipelineRunsPath  = ArtifactsPath + "/" + PipelineRunsDir
	ImageDigestsDir   = "image-digests"
	ImageDigestsPath  = ArtifactsPath + "/" + ImageDigestsDir
	SonarAnalysisDir  = "sonarqube-analysis"
	SonarAnalysisPath = ArtifactsPath + "/" + SonarAnalysisDir
	AquaScansDir      = "aquasec-scans"
	AquaScansPath     = ArtifactsPath + "/" + AquaScansDir
	CodeCoveragesDir  = "code-coverage"
	CodeCoveragesPath = ArtifactsPath + "/" + CodeCoveragesDir
	XUnitReportsDir   = "xunit-reports"
	XUnitReportsPath  = ArtifactsPath + "/" + XUnitReportsDir
	LintReportsDir    = "lint-reports"
	LintReportsPath   = ArtifactsPath + "/" + LintReportsDir
	DeploymentsDir    = "deployments"
	DeploymentsPath   = ArtifactsPath + "/" + DeploymentsDir
)

// WriteJsonArtifact marshals given "in" struct and writes it into "artifactsPath" as "filename".
func WriteJsonArtifact(in interface{}, artifactsPath, filename string) error {
	out, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("could not marshal artifact: %w", err)
	}
	err = os.MkdirAll(artifactsPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", artifactsPath, err)
	}
	return ioutil.WriteFile(filepath.Join(artifactsPath, filename), out, 0644)
}

// CopyArtifact copies given "sourceFile" into "artifactsPath".
func CopyArtifact(sourceFile, artifactsPath string) error {
	err := os.MkdirAll(artifactsPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", artifactsPath, err)
	}
	return file.Copy(sourceFile, filepath.Join(artifactsPath, filepath.Base(sourceFile)))
}

func ReadArtifactsDir() (map[string][]string, error) {

	artifactsDir := filepath.Join(BaseDir, ArtifactsDir)
	artifactsMap := map[string][]string{}

	items, err := ioutil.ReadDir(artifactsDir)
	if err != nil {
		return artifactsMap, fmt.Errorf("%w", err)
	}

	for _, item := range items {
		if item.IsDir() {
			// artifact subdir here, e.g. "xunit-reports"
			subitems, err := ioutil.ReadDir(filepath.Join(artifactsDir, item.Name()))
			if err != nil {
				log.Fatalf("Failed to read dir %s, %s", item.Name(), err.Error())
			}
			filesInSubDir := []string{}
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					// artifact file here, e.g. "report.xml"
					log.Println(item.Name() + "/" + subitem.Name())
					filesInSubDir = append(filesInSubDir, subitem.Name())
				}
			}

			artifactsMap[item.Name()] = filesInSubDir
		}
	}

	log.Println("artifactsMap: ", artifactsMap)

	// map of directories and files under .ods/artifacts, e.g
	// [
	//	"xunit-reports": ["report.xml"]
	//	"sonarqube-analysis": ["analysis-report.md", "issues-report.csv"],
	// ]

	return artifactsMap, nil
}
