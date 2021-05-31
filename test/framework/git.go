package framework

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
)

func InitAndCommit(wsDir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	os.Chdir(wsDir)
	err = writeFile(".gitignore", ".ods/")
	if err != nil {
		return err
	}
	_, _, err = command.Run("git", []string{"init"})
	if err != nil {
		return err
	}
	_, _, err = command.Run("git", []string{"add", "."})
	if err != nil {
		return err
	}
	_, _, err = command.Run("git", []string{"commit", "-m", "initial commit"})
	if err != nil {
		return err
	}
	return nil
}

func WriteDotOds(wsDir string, projectKey string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	wsName := filepath.Base(wsDir)
	os.Chdir(wsDir)
	err = writeFile(".ods/project", projectKey)
	if err != nil {
		return err
	}
	err = writeFile(".ods/repository", wsName)
	if err != nil {
		return err
	}
	err = writeFile(".ods/component", wsName)
	if err != nil {
		return err
	}
	sha, err := getTrimmedFileContent(".git/refs/heads/master")
	if err != nil {
		return err
	}
	err = writeFile(".ods/git-commit-sha", sha)
	if err != nil {
		return err
	}
	return nil
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
