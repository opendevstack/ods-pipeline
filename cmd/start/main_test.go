package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/test/testserver"
)

func TestDirectoryCleaningSparesCache(t *testing.T) {

	tests := map[string]struct {
		fileSystem         fstest.MapFS
		expectedRemoveAlls []string
	}{
		"testCacheSpared": {
			fstest.MapFS{
				".ods-cache/.a":                                {},
				".ods-cache/deps/dep1.txt":                     {},
				".ods-cache/deps/go/gd1.txt":                   {},
				".ods-cache/deps/go/gd1/foo.txt":               {},
				".ods-cache/deps/go/gd2.txt":                   {},
				".ods-cache/deps/npm/hithere_1.0/package.json": {},
				"src/app.js":                                   {},
				"package.json":                                 {},
				".env":                                         {},
			},
			[]string{
				".env",
				"package.json",
				"src",
			},
		},
		"testCacheSparedCaseSensitive": {
			fstest.MapFS{
				".ods-cache/.a":                                {},
				".ods-cache/deps/dep1.txt":                     {},
				".ods-cache/deps/go/gd1.txt":                   {},
				".ods-cache/deps/go/gd1/foo.txt":               {},
				".ods-cache/deps/go/gd2.txt":                   {},
				".ods-Cache/deps/npm/hithere_1.0/package.json": {},
				"src/app.js":                                   {},
				"package.json":                                 {},
				".env":                                         {},
			},
			[]string{
				".env",
				".ods-Cache",
				"package.json",
				"src",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			removed := []string{}
			deleteDirectoryContentsSpareCache(
				FileSystemBase{tc.fileSystem, "."},
				func(path string, isDir bool) error {
					removed = append(removed, path)
					return nil
				})

			if diff := cmp.Diff(tc.expectedRemoveAlls, removed); diff != "" {
				t.Fatalf("expected (-want +got):\n%s", diff)
			}
		})
	}
}
func TestCacheCleaning(t *testing.T) {

	tests := map[string]struct {
		fileSystem         fstest.MapFS
		expectedRemoveAlls []string
	}{
		"testCacheClean": {
			fstest.MapFS{
				".ods-cache/.a":                                {},
				".ods-cache/deps/dep1.txt":                     {},
				".ods-cache/deps/go/gd1.txt":                   {},
				".ods-cache/deps/go/gd1/foo.txt":               {},
				".ods-cache/deps/go/gd2.txt":                   {},
				".ods-cache/deps/npm/hithere_1.0/package.json": {},
			},
			[]string{
				".ods-cache/.a",
				".ods-cache/deps/dep1.txt",
			},
		},
		"testCacheCleanNotRemovingOutside files": {
			fstest.MapFS{
				".ods-cache/.a":                  {},
				".ods-cache/deps/dep1.txt":       {},
				".ods-cache/deps/go/gd1.txt":     {},
				".ods-cache/deps/go/gd1/foo.txt": {},
				".ods-cache/deps/go/gd2.txt":     {},
				"src/app.js":                     {},
				"package.json":                   {},
				".env":                           {},
			},
			[]string{
				".ods-cache/.a",
				".ods-cache/deps/dep1.txt",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			removed := []string{}
			cleanCache(
				FileSystemBase{tc.fileSystem, "."},
				func(path string, isDir bool) error {
					removed = append(removed, path)
					return nil
				})

			if diff := cmp.Diff(tc.expectedRemoveAlls, removed); diff != "" {
				t.Fatalf("expected (-want +got):\n%s", diff)
			}
		})
	}
}

func TestApplyVersionTags(t *testing.T) {
	ctxt := &pipelinectxt.ODSContext{
		Version:      "1.0.0",
		Project:      "PRJ",
		Repository:   "my-repo",
		GitCommitSHA: "8d351a10fb428c0c1239530256e21cf24f136e73",
	}
	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: "s3cr3t", // does not matter
		BaseURL:  srv.Server.URL,
	})

	tests := map[string]struct {
		env             *config.Environment
		prepareServer   func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext)
		checkServer     func(t *testing.T, srv *testserver.TestServer)
		wantError       string
		wantOutContains string
	}{
		"no tagging for DEV stage": {
			env:       &config.Environment{Name: "foo", Stage: config.DevStage},
			wantError: "",
		},
		"no tagging for QA stage when final tag already exists": {
			env: &config.Environment{Name: "foo", Stage: config.QAStage},
			prepareServer: func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext) {
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					200, "start-cmd/tag-list-final-version-exists.json",
				)
			},
			wantError:       "",
			wantOutContains: "Final version tag exists already.",
		},
		"no tagging for PROD stage when final tag already exists": {
			env: &config.Environment{Name: "foo", Stage: config.ProdStage},
			prepareServer: func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext) {
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					200, "start-cmd/tag-list-final-version-exists.json",
				)
			},
			wantError:       "",
			wantOutContains: "Final version tag exists already.",
		},
		"new RC tag for QA stage when no related RC tag exists": {
			env: &config.Environment{Name: "foo", Stage: config.QAStage},
			prepareServer: func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext) {
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					200, "start-cmd/tag-list-unrelated.json",
				)
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					201, "start-cmd/tag-create.json",
				)
			},
			checkServer: func(t *testing.T, srv *testserver.TestServer) {
				tagPayload := lastTagPayload(t, srv)
				wantTag := "v1.0.0-rc.1"
				if tagPayload.Name != wantTag {
					t.Fatalf("want tag: %s, got %s", wantTag, tagPayload.Name)
				}
			},
			wantError:       "",
			wantOutContains: "",
		},
		"next RC tag for QA stage when related RC tags exist": {
			env: &config.Environment{Name: "foo", Stage: config.QAStage},
			prepareServer: func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext) {
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					200, "start-cmd/tag-list-related.json",
				)
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					201, "start-cmd/tag-create.json",
				)
			},
			checkServer: func(t *testing.T, srv *testserver.TestServer) {
				tagPayload := lastTagPayload(t, srv)
				wantTag := "v1.0.0-rc.3"
				if tagPayload.Name != wantTag {
					t.Fatalf("want tag: %s, got %s", wantTag, tagPayload.Name)
				}
			},
			wantError:       "",
			wantOutContains: "",
		},
		"abort for PROD stage when no RC tag exists": {
			env: &config.Environment{Name: "foo", Stage: config.ProdStage},
			prepareServer: func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext) {
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					200, "start-cmd/tag-list-unrelated.json",
				)
			},
			wantError:       "cannot proceed to prod stage: no release candidate tag found for 1.0.0. Deploy to QA before deploying to Prod",
			wantOutContains: "",
		},
		"abort for PROD stage when latest RC tag does not match checked out commit": {
			env: &config.Environment{Name: "foo", Stage: config.ProdStage},
			prepareServer: func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext) {
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					200, "start-cmd/tag-list-related.json",
				)
				ctxt.GitCommitSHA = "8d51122def5632836d1cb1026e879069e10a1e13"
			},
			wantError:       "cannot proceed to prod stage: latest release candidate tag for 1.0.0 does not point to checked out commit, cowardly refusing to deploy",
			wantOutContains: "",
		},
		"final tag for PROD stage when usable RC tag exists": {
			env: &config.Environment{Name: "foo", Stage: config.ProdStage},
			prepareServer: func(t *testing.T, srv *testserver.TestServer, ctxt *pipelinectxt.ODSContext) {
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					200, "start-cmd/tag-list-related.json",
				)
				srv.EnqueueResponse(
					t, "/rest/api/1.0/projects/PRJ/repos/my-repo/tags",
					201, "start-cmd/tag-create.json",
				)
			},
			checkServer: func(t *testing.T, srv *testserver.TestServer) {
				tagPayload := lastTagPayload(t, srv)
				wantTag := "v1.0.0"
				if tagPayload.Name != wantTag {
					t.Fatalf("want tag: %s, got %s", wantTag, tagPayload.Name)
				}
			},
			wantError:       "",
			wantOutContains: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			srv.Reset()
			clonedCtxt := ctxt.Copy()
			if tc.prepareServer != nil {
				tc.prepareServer(t, srv, clonedCtxt)
			}
			var stdout bytes.Buffer
			logger := &logging.LeveledLogger{Level: logging.LevelDebug, StdoutOverride: &stdout}
			err := applyVersionTags(logger, bitbucketClient, clonedCtxt, nil, tc.env)
			if len(tc.wantError) > 0 {
				if err == nil {
					t.Fatalf("want err: %s, got none", tc.wantError)
				}
				if tc.wantError != err.Error() {
					t.Fatalf("want err: %s, got err: %s", tc.wantError, err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
			}
			if len(tc.wantOutContains) > 0 {
				if !strings.Contains(stdout.String(), tc.wantOutContains) {
					t.Fatalf("want out to contain: %s, got out: %s", tc.wantOutContains, stdout.String())
				}
			}
			if tc.checkServer != nil {
				tc.checkServer(t, srv)
			}
		})
	}
}

func lastTagPayload(t *testing.T, srv *testserver.TestServer) bitbucket.TagCreatePayload {
	req, err := srv.LastRequest()
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	var tagPayload bitbucket.TagCreatePayload
	err = json.Unmarshal(body, &tagPayload)
	if err != nil {
		t.Fatal(err)
	}
	return tagPayload
}
