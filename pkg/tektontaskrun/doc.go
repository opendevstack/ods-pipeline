/*
Package tektontaskrun implements ODS Pipeline independent functionality to run
Tekton tasks in a KinD cluster.

Using tektontaskrun it is possible to start a KinD cluster, configure it (e.g.
by setting up a temporary namespace), and running a Tekton task.

tektontaskrun is intended to be used by CLI programs and as a library for
testing Tekton tasks using Go.

Example usage:

	package test

	import (
		"log"
		"os"
		"path/filepath"
		"testing"

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
			ttr.WithWorkspace("source", "my-sample-app"),
			ttr.AfterRun(func(config *ttr.TaskRunConfig, run *tekton.TaskRun) {
				wd := config.WorkspaceConfigs["source"].Dir
				// e.g. check files in workspace ...
			}),
		); err != nil {
			t.Fatal(err)
		}
	}

	// further tests here ...
*/
package tektontaskrun
