package tasks

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	kclient "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

var alwaysKeepTmpWorkspacesFlag = flag.Bool("always-keep-tmp-workspaces", false, "Whether to keep temporary workspaces from taskruns even when test is successful")

const (
	taskKindRef = "ClusterTask"
)

func checkODSContext(t *testing.T, repoDir string, want *pipelinectxt.ODSContext) {
	checkODSFileContent(t, repoDir, "component", want.Component)
	checkODSFileContent(t, repoDir, "git-commit-sha", want.GitCommitSHA)
	checkODSFileContent(t, repoDir, "git-full-ref", want.GitFullRef)
	checkODSFileContent(t, repoDir, "git-ref", want.GitRef)
	checkODSFileContent(t, repoDir, "git-url", want.GitURL)
	checkODSFileContent(t, repoDir, "namespace", want.Namespace)
	checkODSFileContent(t, repoDir, "pr-base", want.PullRequestBase)
	checkODSFileContent(t, repoDir, "pr-key", want.PullRequestKey)
	checkODSFileContent(t, repoDir, "project", want.Project)
	checkODSFileContent(t, repoDir, "repository", want.Repository)
}

func checkODSFileContent(t *testing.T, wsDir, filename, want string) {
	checkFileContent(t, filepath.Join(wsDir, pipelinectxt.BaseDir), filename, want)
}

func checkFileContent(t *testing.T, wsDir, filename, want string) {
	got, err := getTrimmedFileContent(filepath.Join(wsDir, filename))
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	if got != want {
		t.Fatalf("got '%s', want '%s' in file %s", got, want, filename)
	}
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func trimmedFileContentOrFatal(t *testing.T, filename string) string {
	c, err := getTrimmedFileContent(filename)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func checkFileContentContains(t *testing.T, wsDir, filename, wantContains string) {
	got, err := getFileContentLean(filepath.Join(wsDir, filename))
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	if !strings.Contains(got, wantContains) {
		t.Fatalf("got '%s', wantContains '%s' in file %s", got, wantContains, filename)
	}
}

func getFileContentLean(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	contentStr := strings.ReplaceAll(string(content), "\t", "")
	contentStr = strings.ReplaceAll(contentStr, "\n", "")
	contentStr = strings.ReplaceAll(contentStr, " ", "")

	return contentStr, nil
}

func runTaskTestCases(t *testing.T, taskName string, testCases map[string]tasktesting.TestCase) {
	c, ns := tasktesting.Setup(t,
		tasktesting.SetupOpts{
			SourceDir:        tasktesting.StorageSourceDir,
			StorageCapacity:  tasktesting.StorageCapacity,
			StorageClassName: tasktesting.StorageClassName,
		},
	)

	tasktesting.CleanupOnInterrupt(func() { tasktesting.TearDown(t, c, ns) }, t.Logf)
	defer tasktesting.TearDown(t, c, ns)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.TaskVariant != "" {
				taskName = fmt.Sprintf("%s-%s", taskName, tc.TaskVariant)
			}
			if tc.Timeout == 0 {
				tc.Timeout = 5 * time.Minute
			}
			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             taskKindRef,
				TaskName:                taskName,
				Clients:                 c,
				Namespace:               ns,
				Timeout:                 tc.Timeout,
				AlwaysKeepTmpWorkspaces: *alwaysKeepTmpWorkspacesFlag,
			})
		})
	}
}

func checkSonarQualityGate(t *testing.T, c *kclient.Clientset, namespace, sonarProject string, qualityGateFlag bool, wantQualityGateStatus string) {

	sonarToken, err := kubernetes.GetSecretKey(c, namespace, "ods-sonar-auth", "password")
	if err != nil {
		t.Fatalf("could not get SonarQube token: %s", err)
	}

	sonarClient := sonar.NewClient(&sonar.ClientConfig{
		APIToken:      sonarToken,
		BaseURL:       "http://localhost:9000", // use localhost instead of ods-test-sonarqube.kind!
		ServerEdition: "community",
	})

	if qualityGateFlag {
		qualityGateResult, err := sonarClient.QualityGateGet(
			sonar.QualityGateGetParams{Project: sonarProject},
		)
		if err != nil {
			t.Fatal(err)
		}
		actualStatus := qualityGateResult.ProjectStatus.Status
		if actualStatus != wantQualityGateStatus {
			t.Fatalf("Got: %s, want: %s", actualStatus, wantQualityGateStatus)
		}

	}

}

func createODSYML(wsDir string, o *config.ODS) error {
	y, err := yaml.Marshal(o)
	if err != nil {
		return err
	}
	filename := filepath.Join(wsDir, "ods.yaml")
	return ioutil.WriteFile(filename, y, 0644)
}

func checkBuildStatus(t *testing.T, c *bitbucket.Client, gitCommit, wantBuildStatus string) {
	buildStatus, err := c.BuildStatusGet(gitCommit)
	if err != nil {
		t.Fatal(err)
	}
	if buildStatus.State != wantBuildStatus {
		t.Fatalf("Got: %s, want: %s", buildStatus.State, wantBuildStatus)
	}

}
