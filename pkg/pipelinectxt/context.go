package pipelinectxt

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	WIP           = "WIP"
	BaseDir       = ".ods"
	SubreposDir   = "repos"
	SubreposPath  = BaseDir + "/" + SubreposDir
)

type ODSContext struct {
	Project         string
	Repository      string
	Component       string
	Namespace       string
	GitCommitSHA    string
	GitFullRef      string
	GitRef          string
	GitURL          string
	Version         string
	Environment     string
	PullRequestBase string
	PullRequestKey  string
}

// WriteCache writes the ODS context to .ods
func (o *ODSContext) WriteCache(wsDir string) error {
	dotODS := filepath.Join(wsDir, BaseDir)
	if _, err := os.Stat(dotODS); os.IsNotExist(err) {
		err = os.Mkdir(dotODS, 0755)
		if err != nil {
			return fmt.Errorf("could not create %s: %s", BaseDir, err)
		}
	}
	files := map[string]string{
		BaseDir + "/project":        o.Project,
		BaseDir + "/repository":     o.Repository,
		BaseDir + "/component":      o.Component,
		BaseDir + "/environment":    o.Environment,
		BaseDir + "/version":        o.Version,
		BaseDir + "/git-commit-sha": o.GitCommitSHA,
		BaseDir + "/git-url":        o.GitURL,
		BaseDir + "/git-ref":        o.GitRef,
		BaseDir + "/git-full-ref":   o.GitFullRef,
		BaseDir + "/namespace":      o.Namespace,
		BaseDir + "/pr-base":        o.PullRequestBase,
		BaseDir + "/pr-key":         o.PullRequestKey,
	}
	for filename, content := range files {
		err := writeFile(filepath.Join(wsDir, filename), content)
		if err != nil {
			return fmt.Errorf("could not write %s: %w", filename, err)
		}
	}
	return nil
}

// ReadCache reads ODS context from .ods
// TODO: test that this works
func (o *ODSContext) ReadCache(wsDir string) error {
	files := map[string]*string{
		BaseDir + "/project":        &o.Project,
		BaseDir + "/repository":     &o.Repository,
		BaseDir + "/component":      &o.Component,
		BaseDir + "/environment":    &o.Environment,
		BaseDir + "/version":        &o.Version,
		BaseDir + "/git-commit-sha": &o.GitCommitSHA,
		BaseDir + "/git-url":        &o.GitURL,
		BaseDir + "/git-ref":        &o.GitRef,
		BaseDir + "/git-full-ref":   &o.GitFullRef,
		BaseDir + "/namespace":      &o.Namespace,
		BaseDir + "/pr-base":        &o.PullRequestBase,
		BaseDir + "/pr-key":         &o.PullRequestKey,
	}
	for filename, content := range files {
		if len(*content) == 0 {
			cached, err := getTrimmedFileContent(filepath.Join(wsDir, filename))
			if err != nil {
				return fmt.Errorf("could not read %s: %w", filename, err)
			}
			*content = cached
		}

	}
	return nil
}

// Assemble builds an ODS context based on given wsDir directory.
// The information is gathered from the .git directory.
func (o *ODSContext) Assemble(wsDir string) error {
	if len(o.Namespace) == 0 {
		ns, err := getTrimmedFileContent(namespaceFile)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", namespaceFile, err)
		}
		o.Namespace = ns
	}
	if len(o.GitFullRef) == 0 {
		gitHead, err := getTrimmedFileContent(filepath.Join(wsDir, ".git/HEAD"))
		if err != nil {
			return fmt.Errorf("could not read .git/HEAD: %w", err)
		}
		o.GitFullRef = strings.TrimPrefix(gitHead, "ref: ")
	}
	if len(o.GitRef) == 0 {
		gitFullRefParts := strings.SplitN(o.GitFullRef, "/", 3)
		if len(gitFullRefParts) != 3 {
			return fmt.Errorf("cannot extract git ref from .git/HEAD: %s", o.GitFullRef)
		}
		o.GitRef = gitFullRefParts[2]
	}
	if len(o.GitURL) == 0 {
		gitURL, err := readRemoteOriginURL(filepath.Join(wsDir, ".git/config"))
		if err != nil {
			return fmt.Errorf("could not get remote origin URL: %w", err)
		}
		o.GitURL = gitURL
	}
	if len(o.GitCommitSHA) == 0 {
		gitSHA, err := getTrimmedFileContent(filepath.Join(wsDir, ".git", o.GitFullRef))
		if err != nil {
			return fmt.Errorf("could not read .git/%s: %w", o.GitFullRef, err)
		}
		o.GitCommitSHA = gitSHA
	}
	u, err := url.Parse(o.GitURL)
	if err != nil {
		return fmt.Errorf("could not parse remote origin URL: %w", err)
	}
	pathParts := strings.Split(u.Path, "/")
	organisation := pathParts[len(pathParts)-2]
	repository := pathParts[len(pathParts)-1]
	if len(o.Project) == 0 {
		o.Project = strings.ToLower(organisation)
	}
	if len(o.Repository) == 0 {
		o.Repository = filenameWithoutExtension(repository)
	}
	if len(o.Component) == 0 {
		o.Component = strings.TrimPrefix(o.Repository, o.Project+"-")
	}
	if len(o.Version) == 0 {
		o.Version = WIP
	}
	return nil
}

func (o *ODSContext) Copy() *ODSContext {
	n := *o
	return &n
}

func WriteGitIgnore(path string) error {
	odsPipelineIgnore := fmt.Sprintf("/%s\n/.ods-cache\n", BaseDir)
	return writeFile(path, odsPipelineIgnore)
}

func filenameWithoutExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func readRemoteOriginURL(gitConfigFilename string) (string, error) {
	file, err := os.Open(gitConfigFilename)
	if err != nil {
		return "", fmt.Errorf("could not open git config: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(trimmedLine, "url = ") {
			return strings.TrimPrefix(trimmedLine, "url = "), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error scanning file: %w", err)
	}
	return "", errors.New("did not find URL in git config")
}

func writeFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
