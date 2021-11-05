package pipelinectxt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/file"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"sigs.k8s.io/yaml"
)

const (
	ArtifactsDir              = "artifacts"
	ArtifactsPath             = BaseDir + "/" + ArtifactsDir
	PipelineRunsDir           = "pipeline-runs"
	PipelineRunsPath          = ArtifactsPath + "/" + PipelineRunsDir
	ImageDigestsDir           = "image-digests"
	ImageDigestsPath          = ArtifactsPath + "/" + ImageDigestsDir
	SonarAnalysisDir          = "sonarqube-analysis"
	SonarAnalysisPath         = ArtifactsPath + "/" + SonarAnalysisDir
	AquaScansDir              = "aquasec-scans"
	AquaScansPath             = ArtifactsPath + "/" + AquaScansDir
	CodeCoveragesDir          = "code-coverage"
	CodeCoveragesPath         = ArtifactsPath + "/" + CodeCoveragesDir
	XUnitReportsDir           = "xunit-reports"
	XUnitReportsPath          = ArtifactsPath + "/" + XUnitReportsDir
	LintReportsDir            = "lint-reports"
	LintReportsPath           = ArtifactsPath + "/" + LintReportsDir
	DeploymentsDir            = "deployments"
	DeploymentsPath           = ArtifactsPath + "/" + DeploymentsDir
	ArtifactsManifestFilename = "manifest.json"
)

// ArtifactsManifest represents all downloaded artifacts.
type ArtifactsManifest struct {
	// SourceRepository identifies the repository artifacts where downloaded from
	SourceRepository string         `json:"sourceRepository"`
	Artifacts        []ArtifactInfo `json:"artifacts"`
}

// ArtifactInfo represents one artifact.
type ArtifactInfo struct {
	URL       string `json:"url"`
	Directory string `json:"directory"`
	Name      string `json:"name"`
}

// ReadArtifactsManifestFromFile reads an artifact manifest from given filename or errors.
func ReadArtifactsManifestFromFile(filename string) (*ArtifactsManifest, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}
	var am *ArtifactsManifest
	err = yaml.UnmarshalStrict(body, &am, func(dec *json.Decoder) *json.Decoder {
		dec.DisallowUnknownFields()
		return dec
	})
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal manifest: %w", err)
	}
	return am, nil
}

// Contains checks whether given directory/name is already present in repository.
func (am *ArtifactsManifest) Contains(repository, directory, name string) bool {
	if repository != am.SourceRepository {
		return false
	}
	for _, a := range am.Artifacts {
		if a.Directory == directory && a.Name == name {
			return true
		}
	}
	return false
}

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

// ReadArtifactsDir reads all artifacts in ArtifactsPath.
// Only files in subdirectories are considered as artifacts.
// Example: [
//	"xunit-reports": ["report.xml"]
//	"sonarqube-analysis": ["analysis-report.md", "issues-report.csv"],
// ]
func ReadArtifactsDir() (map[string][]string, error) {

	artifactsMap := map[string][]string{}

	items, err := ioutil.ReadDir(ArtifactsPath)
	if err != nil {
		return artifactsMap, fmt.Errorf("%w", err)
	}

	for _, item := range items {
		if item.IsDir() {
			// artifact subdir here, e.g. "xunit-reports"
			subitems, err := ioutil.ReadDir(filepath.Join(ArtifactsPath, item.Name()))
			if err != nil {
				log.Fatalf("Failed to read dir %s, %s", item.Name(), err.Error())
			}
			filesInSubDir := []string{}
			for _, subitem := range subitems {
				if !subitem.IsDir() {
					// artifact file here, e.g. "report.xml"
					filesInSubDir = append(filesInSubDir, subitem.Name())
				}
			}

			artifactsMap[item.Name()] = filesInSubDir
		}
	}

	return artifactsMap, nil
}

// ArtifactGroupBase returns the group base in which aritfacts are stored for
// the given ODS pipeline context.
func ArtifactGroupBase(ctxt *ODSContext) string {
	return nexus.ArtifactGroupBase(ctxt.Project, ctxt.Repository, ctxt.GitCommitSHA)
}

// ArtifactGroup returns the group in which aritfacts are stored for the given
// ODS pipeline context and the subdir.
func ArtifactGroup(ctxt *ODSContext, subdir string) string {
	return nexus.ArtifactGroup(ctxt.Project, ctxt.Repository, ctxt.GitCommitSHA, subdir)
}

// DownloadGroup searches given repositories in order for assets in given group.
// As soon as one repository has any asset in the group, the search is stopped
// and all fond artifacts are downloaded into artifactsDir.
// An artifacts manifest is returned describing the downloaded files.
// When none of the given repositories contains any artifacts under the group,
// no artifacts are downloaded and no error is returned.
func DownloadGroup(nexusClient nexus.ClientInterface, repositories []string, group, artifactsDir string, logger logging.LeveledLoggerInterface) (*ArtifactsManifest, error) {
	// We want to target all artifacts underneath the group, hence the trailing '*'.
	nexusSearchGroup := fmt.Sprintf("%s/*", group)
	am := &ArtifactsManifest{
		Artifacts: []ArtifactInfo{},
	}
	sourceRepo, urls, err := searchForAssets(nexusClient, nexusSearchGroup, repositories, logger)
	if err != nil {
		return nil, err
	}
	am.SourceRepository = sourceRepo

	for _, s := range urls {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		urlPathParts := strings.Split(u.Path, fmt.Sprintf("%s/", group))
		if len(urlPathParts) != 2 {
			return nil, fmt.Errorf("unexpected URL path (must contain two parts after group): %s", u.Path)
		}
		fileWithSubPath := urlPathParts[1] // e.g. "pipeline-runs/foo-zh9gt0.json"
		if !strings.Contains(fileWithSubPath, "/") {
			return nil, fmt.Errorf("unexpected URL path (must contain a subfolder after the commit SHA): %s", fileWithSubPath)
		}
		aritfactName := path.Base(fileWithSubPath) // e.g. "pipeline-runs"
		artifactType := path.Dir(fileWithSubPath)  // e.g. "foo-zh9gt0.json"
		artifactsSubPath := filepath.Join(artifactsDir, artifactType)
		if _, err := os.Stat(artifactsSubPath); os.IsNotExist(err) {
			if err := os.MkdirAll(artifactsSubPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory: %s, error: %w", artifactsSubPath, err)
			}
		}
		outfile := filepath.Join(artifactsDir, fileWithSubPath)
		_, err = nexusClient.Download(s, outfile)
		if err != nil {
			return nil, err
		}
		am.Artifacts = append(am.Artifacts, ArtifactInfo{
			URL:       s,
			Directory: artifactType,
			Name:      aritfactName,
		})
	}
	return am, nil
}

// searchForAssets looks for assets in searchGroup for each repository in order.
// Once some assets are found, the repository and the found URLs are returned,
// skipping any further repositories that are given.
func searchForAssets(nexusClient nexus.ClientInterface, searchGroup string, repositories []string, logger logging.LeveledLoggerInterface) (string, []string, error) {
	for _, r := range repositories {
		urls, err := nexusClient.Search(r, searchGroup)
		if err != nil {
			return "", nil, err
		}
		if len(urls) > 0 {
			logger.Infof("Found artifacts in repository %s inside group %s ...", r, searchGroup)
			return r, urls, nil
		}
		logger.Infof("No artifacts found in repository %s inside group %s.", r, searchGroup)
	}
	return "", []string{}, nil
}
