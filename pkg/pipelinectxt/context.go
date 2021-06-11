package pipelinectxt

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
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
	PullRequestBase string
	PullRequestKey  string
}

// WriteCache writes the ODS context to .ods
func (o *ODSContext) WriteCache(wsDir string) error {
	dotODS := filepath.Join(wsDir, ".ods")
	if _, err := os.Stat(dotODS); os.IsNotExist(err) {
		err = os.Mkdir(dotODS, 0755)
		if err != nil {
			return fmt.Errorf("could not create .ods: %s", err)
		}
	}
	files := map[string]string{
		".ods/project":        o.Project,
		".ods/repository":     o.Repository,
		".ods/component":      o.Component,
		".ods/git-commit-sha": o.GitCommitSHA,
		".ods/git-url":        o.GitURL,
		".ods/git-ref":        o.GitRef,
		".ods/git-full-ref":   o.GitFullRef,
		".ods/namespace":      o.Namespace,
		".ods/pr-base":        o.PullRequestBase,
		".ods/pr-key":         o.PullRequestKey,
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
		".ods/project":        &o.Project,
		".ods/repository":     &o.Repository,
		".ods/component":      &o.Component,
		".ods/git-commit-sha": &o.GitCommitSHA,
		".ods/git-url":        &o.GitURL,
		".ods/git-ref":        &o.GitRef,
		".ods/git-full-ref":   &o.GitFullRef,
		".ods/namespace":      &o.Namespace,
		".ods/pr-base":        &o.PullRequestBase,
		".ods/pr-key":         &o.PullRequestKey,
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

func (o *ODSContext) Assemble(wsDir string) error {
	absoluteWsDir, err := filepath.Abs(wsDir)
	if err != nil {
		return fmt.Errorf("could not get absolute path of %s: %w", wsDir, err)
	}
	wsName := filepath.Base(absoluteWsDir)

	if len(o.Namespace) == 0 {
		ns, err := getTrimmedFileContent(namespaceFile)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", namespaceFile, err)
		}
		o.Namespace = ns
	}
	if len(o.Project) == 0 {
		o.Project = strings.TrimSuffix(o.Namespace, "-cd")
	}
	if len(o.Repository) == 0 {
		o.Repository = wsName
	}
	if len(o.Component) == 0 {
		o.Component = strings.TrimPrefix(o.Repository, o.Project+"-")
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
	return nil
}

func (o *ODSContext) ReadArtifactsDir() (map[string][]string, error) {

	artifactsDir := ".ods/artifacts/"
	artifactsMap := map[string][]string{}

	items, err := ioutil.ReadDir(artifactsDir)
	if err != nil {
		return artifactsMap, fmt.Errorf("%w", err)
	}

	for _, item := range items {
		if item.IsDir() {
			// artifact subdir here, e.g. "xunit-reports"
			subitems, err := ioutil.ReadDir(artifactsDir + item.Name())
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
