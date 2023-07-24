package scripts_test

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/opendevstack/ods-pipeline/internal/command"
)

const (
	downloadAquaScannerScript = "../../build/package/scripts/download-aqua-scanner.sh"
	fakeScannerBinary         = `#!/bin/bash
echo 1.7.3`
)

var md5Bin = flag.String("md5bin", "md5", "md5 binary to use")

func TestCachedDownload(t *testing.T) {
	hits := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		fmt.Fprintln(w, fakeScannerBinary)
	}))
	defer ts.Close()
	dir, cleanupDir := tmpDir(t)
	defer cleanupDir()

	t.Log("Fresh run -> download")
	runScriptOrFatal(t, dir, fmt.Sprintf("%s/foo", ts.URL))
	if hits != 1 {
		t.Error("Wanted hit, got none")
	}
	fileExistsInDir(t, dir, "aquasec", ".md5-aquasec")

	t.Log("Second run for same URL -> no download")
	runScriptOrFatal(t, dir, fmt.Sprintf("%s/foo", ts.URL))
	if hits > 1 {
		t.Error("Wanted no further hit, got more")
	}
	fileExistsInDir(t, dir, "aquasec", ".md5-aquasec")

	t.Log("Third run for different URL -> download")
	runScriptOrFatal(t, dir, fmt.Sprintf("%s/bar", ts.URL))
	if hits != 2 {
		t.Error("Wanted further hit, got none")
	}
	fileExistsInDir(t, dir, "aquasec", ".md5-aquasec")
}

func TestSkipDownload(t *testing.T) {
	dir, cleanupDir := tmpDir(t)
	defer cleanupDir()

	t.Log("No URL")
	runScriptOrFatal(t, dir, "")
	fileDoesNotExistInDir(t, dir, "aquasec", ".md5-aquasec")

	t.Log("URL set to 'none'")
	runScriptOrFatal(t, dir, "none")
	fileDoesNotExistInDir(t, dir, "aquasec", ".md5-aquasec")
}

func TestBrokenDownload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	}))
	defer ts.Close()
	dir, cleanupDir := tmpDir(t)
	defer cleanupDir()

	t.Log("Download")
	err := runScript(t, dir, fmt.Sprintf("%s/foo", ts.URL))
	if err == nil {
		t.Fatal("script should error on broken download")
	}
	fileDoesNotExistInDir(t, dir, "aquasec", ".md5-aquasec")
}

// runScriptOrFatal calls runScript, then t.Fatal on error.
func runScriptOrFatal(t *testing.T, dir, url string) {
	if err := runScript(t, dir, url); err != nil {
		t.Fatal(err)
	}
}

// runScript runs the download script against given url
// and places the downloaded file into dir.
func runScript(t *testing.T, dir, url string) error {
	return command.Run(
		downloadAquaScannerScript,
		[]string{
			fmt.Sprintf("--bin-dir=%s", dir),
			fmt.Sprintf("--aqua-scanner-url=%s", url),
		}, []string{fmt.Sprintf("MD5_BIN=%s", *md5Bin)},
		&testingLogWriter{t},
		&testingLogWriter{t},
	)
}

// fileExistsInDir checks if file(s) exist in dir or errors.
func fileExistsInDir(t *testing.T, dir string, files ...string) {
	for _, file := range files {
		f := fmt.Sprintf("%s/%s", dir, file)
		if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
			t.Errorf("Want file %s, got %s", f, err)
		}
	}
}

// fileDoesNotExistInDir checks if file(s) exist in dir or errors.
func fileDoesNotExistInDir(t *testing.T, dir string, files ...string) {
	for _, file := range files {
		f := fmt.Sprintf("%s/%s", dir, file)
		if _, err := os.Stat(f); err == nil {
			t.Errorf("Did not want file %s", f)
		}
	}
}

// tmpDir creates a temp dir or fails.
func tmpDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp(".", "download-aqua-scanner-")
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

// testingLogWriter implements io.Writer so that
// it can proxy to t.Log when an io.Writer is required.
type testingLogWriter struct {
	t *testing.T
}

// Write proxies to t.Logf.
func (f *testingLogWriter) Write(p []byte) (n int, err error) {
	f.t.Helper()
	f.t.Logf("%s", string(p))
	return len(p), nil
}
