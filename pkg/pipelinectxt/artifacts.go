package pipelinectxt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/opendevstack/ods-pipeline/internal/file"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
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
	SBOMsDir                  = "sboms"
	SBOMsPath                 = ArtifactsPath + "/" + SBOMsDir
	SBOMsFormat               = "spdx"
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
	// Repository is the artifact repository from which the manifests were downloaded.
	Repository string `json:"repository"`
	// Artifacts lists all artifacts downloaded.
	Artifacts []ArtifactInfo `json:"artifacts"`
}

// ArtifactInfo represents one artifact.
type ArtifactInfo struct {
	URL       string `json:"url"`
	Directory string `json:"directory"`
	Name      string `json:"name"`
}

// NewArtifactsManifest returns a new ArtifactsManifest instance.
func NewArtifactsManifest(repository string, artifacts ...ArtifactInfo) *ArtifactsManifest {
	return &ArtifactsManifest{Repository: repository, Artifacts: artifacts}
}

// ReadArtifactsManifestFromFile reads an artifact manifest from given filename or errors.
func ReadArtifactsManifestFromFile(filename string) (*ArtifactsManifest, error) {
	body, err := os.ReadFile(filename)
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
	if am.Repository != repository {
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
	return os.WriteFile(filepath.Join(artifactsPath, filename), out, 0644)
}

// CopyArtifact copies given "sourceFile" into "artifactsPath".
func CopyArtifact(sourceFile, artifactsPath string) error {
	err := os.MkdirAll(artifactsPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create %s: %w", artifactsPath, err)
	}
	return file.Copy(sourceFile, filepath.Join(artifactsPath, filepath.Base(sourceFile)))
}

// ReadArtifactsDir reads all artifacts in checkoutDir/ArtifactsPath.
// Only files in subdirectories are considered as artifacts.
// Example: [
//
//	"xunit-reports": ["report.xml"]
//	"sonarqube-analysis": ["analysis-report.md", "issues-report.csv"],
//
// ]
func ReadArtifactsDir(artifactsDir string) (map[string][]string, error) {
	artifactsMap := map[string][]string{}

	items, err := os.ReadDir(artifactsDir)
	if err != nil {
		return artifactsMap, fmt.Errorf("%w", err)
	}

	for _, item := range items {
		if item.IsDir() {
			subpath := filepath.Join(artifactsDir, item.Name())
			filesInSubDir, err := ReadArtifactFilenames(subpath)
			if err != nil {
				return artifactsMap, fmt.Errorf("read artifacts in %s: %w", subpath, err)
			}
			artifactsMap[item.Name()] = filesInSubDir
		}
	}

	return artifactsMap, nil
}

// ReadArtifactFilenames returns all filenames in artifactsPath.
// If artifactsPath does not exist, an empty list is returned.
func ReadArtifactFilenames(artifactsPath string) ([]string, error) {
	if _, err := os.Stat(artifactsPath); errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	var files []string
	f, err := os.ReadDir(artifactsPath)
	if err != nil {
		return files, fmt.Errorf("read dir %s: %w", artifactsPath, err)
	}
	for _, fi := range f {
		if !fi.IsDir() {
			files = append(files, fi.Name())
		}
	}
	return files, nil
}

// ReadArtifactFilesIncludingSubrepos reads artifacts in artifactPath of the current
// repository and all subrepos given. The returned files include path and filename.
func ReadArtifactFilesIncludingSubrepos(artifactPath string, subrepos []fs.DirEntry) ([]string, error) {
	topLevelArtifacts, err := ReadArtifactFilenames(artifactPath)
	if err != nil {
		return []string{}, fmt.Errorf("read artifacts: %w", err)
	}
	artifacts := []string{}
	for _, a := range topLevelArtifacts {
		artifacts = append(artifacts, filepath.Join(artifactPath, a))
	}
	for _, s := range subrepos {
		subrepoArtifactsPath := filepath.Join(SubreposPath, s.Name(), artifactPath)
		subArtifacts, err := ReadArtifactFilenames(subrepoArtifactsPath)
		if err != nil {
			return []string{}, fmt.Errorf("read artifacts of %s: %w", subrepoArtifactsPath, err)
		}
		for _, a := range subArtifacts {
			artifacts = append(artifacts, filepath.Join(subrepoArtifactsPath, a))
		}
	}
	return artifacts, nil
}

// DetectSubrepos returns a slice of subrepo directories
func DetectSubrepos() ([]fs.DirEntry, error) {
	subrepos := []fs.DirEntry{}
	if _, err := os.Stat(SubreposPath); err == nil {
		f, err := os.ReadDir(SubreposPath)
		if err != nil {
			return []fs.DirEntry{}, fmt.Errorf("read %s: %w", SubreposPath, err)
		}
		subrepos = f
	}
	return subrepos, nil
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
// and all found artifacts are downloaded into artifactsDir.
// An artifacts manifest is returned describing the downloaded files.
// When none of the given repositories contains any artifacts under the group,
// no artifacts are downloaded and no error is returned.
// If artifactsDir is an empty string, the searched files are not downloaded but
// the artifacts are still recorded in the returned manifest.
func DownloadGroup(
	nexusClient nexus.ClientInterface,
	repository, group, artifactsDir string,
	logger logging.LeveledLoggerInterface) (*ArtifactsManifest, error) {
	// We want to target all artifacts underneath the group, hence the trailing '*'.
	nexusSearchGroup := fmt.Sprintf("%s/*", group)
	am := NewArtifactsManifest(repository)
	urls, err := searchForAssets(nexusClient, nexusSearchGroup, repository, logger)
	if err != nil {
		return nil, err
	}

	if artifactsDir == "" {
		logger.Debugf("Artifacts will not be downloaded but only added to the manifest ...")
	}

	for _, s := range urls {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}
		_, fileWithSubPath, ok := strings.Cut(u.Path, fmt.Sprintf("%s/", group)) // e.g. "pipeline-runs/foo-zh9gt0.json"
		if !ok {
			return nil, fmt.Errorf("unexpected URL path (must contain group '%s'): %s", group, u.Path)
		}
		if !strings.Contains(fileWithSubPath, "/") {
			return nil, fmt.Errorf("unexpected URL path (must contain a subfolder after the commit SHA): %s", fileWithSubPath)
		}
		aritfactName := path.Base(fileWithSubPath) // e.g. "pipeline-runs"
		artifactType := path.Dir(fileWithSubPath)  // e.g. "foo-zh9gt0.json"
		if artifactsDir != "" {
			artifactsSubPath := filepath.Join(artifactsDir, artifactType)
			if _, err := os.Stat(artifactsSubPath); os.IsNotExist(err) {
				if err := os.MkdirAll(artifactsSubPath, 0755); err != nil {
					return nil, fmt.Errorf("failed to create directory: %s, error: %w", artifactsSubPath, err)
				}
			}
			outfile := filepath.Join(artifactsDir, fileWithSubPath)
			if _, err := nexusClient.Download(s, outfile); err != nil {
				return nil, err
			}
		}
		am.Artifacts = append(am.Artifacts, ArtifactInfo{
			URL:       s,
			Directory: artifactType,
			Name:      aritfactName,
		})
	}
	return am, nil
}

// searchForAssets looks for assets in searchGroup of given repository.
// No error is returned when no assets are found.
func searchForAssets(nexusClient nexus.ClientInterface, searchGroup string, repository string, logger logging.LeveledLoggerInterface) ([]string, error) {
	urls, err := nexusClient.Search(repository, searchGroup)
	if err != nil {
		return nil, err
	}
	if len(urls) > 0 {
		logger.Infof("Found artifacts in repository %q inside group %q ...", repository, searchGroup)
		return urls, nil
	}
	logger.Infof("No artifacts found in repository %q inside group %q.", repository, searchGroup)
	return []string{}, nil
}
