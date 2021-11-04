package gittest

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateFakeGitRepoDir creates a temporary directory, which contains a fake
// Git repository. This could be read with pipelinectxt.Assemble for example.
// The returned function should be deferred to clean up the temp dir.
func CreateFakeGitRepoDir(url, branch, sha string) (string, func(), error) {
	cleanupFunc := func() {}
	dir, err := ioutil.TempDir(".", "test-git-repo-")
	if err != nil {
		return "", cleanupFunc, err
	}
	cleanupFunc = func() {
		os.RemoveAll(dir)
	}
	gitDir := filepath.Join(dir, ".git")
	err = os.MkdirAll(gitDir, 0755)
	if err != nil {
		return "", cleanupFunc, err
	}
	err = ioutil.WriteFile(
		filepath.Join(gitDir, "HEAD"),
		[]byte("ref: refs/heads/"+branch),
		0644,
	)
	if err != nil {
		return "", cleanupFunc, err
	}
	err = ioutil.WriteFile(
		filepath.Join(gitDir, "config"),
		[]byte(`[remote "origin"]
	url = `+url+`
	fetch = +refs/heads/*:refs/remotes/origin/*`),
		0644,
	)
	if err != nil {
		return "", cleanupFunc, err
	}
	gitHeadsDir := filepath.Join(gitDir, "refs/heads")
	err = os.MkdirAll(gitHeadsDir, 0755)
	if err != nil {
		return "", cleanupFunc, err
	}
	err = ioutil.WriteFile(
		filepath.Join(gitHeadsDir, branch),
		[]byte(sha),
		0644,
	)
	if err != nil {
		return "", cleanupFunc, err
	}

	return dir, cleanupFunc, nil
}
