/*
Package odstasktest implements ODS Pipeline specific functionality to run
Tekton tasks in a KinD cluster on top of package tektontaskrun.

odstasktest is intended to be used as a library for testing ODS Pipeline
tasks using Go.

Example usage:

	package test

	import (
		"log"
		"os"
		"path/filepath"
		"testing"

		ott "github.com/opendevstack/ods-pipeline/pkg/odstasktest"
		ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
	)

	var (
		namespaceConfig *ttr.NamespaceConfig
		rootPath        = "../.."
	)

	func TestMain(m *testing.M) {
		cc, err := ttr.StartKinDCluster(
			ttr.LoadImage(ttr.ImageBuildConfig{
				Dockerfile: "build/images/Dockerfile.my-task",
				ContextDir: rootPath,
			}),
		)
		if err != nil {
			log.Fatal("Could not start KinD cluster: ", err)
		}
		nc, cleanup, err := ttr.SetupTempNamespace(
			cc,
			ott.StartNexus(),
			ott.InstallODSPipeline(),
			ttr.InstallTaskFromPath(
				filepath.Join(rootPath, "build/tasks/my-task.yaml"),
				nil,
			),
		)
		if err != nil {
			log.Fatal("Could not setup temporary namespace: ", err)
		}
		defer cleanup()
		namespaceConfig = nc
		os.Exit(m.Run())
	}

	func TestMyTask(t *testing.T) {
		if err := ttr.RunTask(
			ttr.InNamespace(namespaceConfig.Name),
			ttr.UsingTask("my-task"),
			ttr.WithStringParams(map[string]string{
				"go-os":       runtime.GOOS,
				"go-arch":     runtime.GOARCH,
			}),
			ott.WithGitSourceWorkspace(t, "../testdata/workspaces/go-sample-app"),
			ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun) {
				ott.AssertFilesExist(
					t, config.WorkspaceConfigs["source"].Dir,
					"docker/Dockerfile",
					"docker/app",
				)
			}),
		); err != nil {
			t.Fatal(err)
		}
	}

	// further tests here ...
*/
package odstasktest
