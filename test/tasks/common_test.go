package tasks

import (
	"crypto/sha256"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/directory"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/sonar"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	"golang.org/x/exp/slices"
	kclient "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"
)

var alwaysKeepTmpWorkspacesFlag = flag.Bool("always-keep-tmp-workspaces", false, "Whether to keep temporary workspaces from taskruns even when test is successful")
var outsideKindFlag = flag.Bool("outside-kind", false, "Whether to continue if not in KinD cluster")
var skipSonarQubeFlag = flag.Bool("skip-sonar", false, "Whether to skip SonarQube steps")

const (
	taskKindRef = "Task"
)

// buildTaskParams forces all SonarQube params to be "falsy"
// if the skipSonarQubeFlag is set.
func buildTaskParams(p map[string]string) map[string]string {
	if *skipSonarQubeFlag {
		p["sonar-skip"] = "true"
		p["sonar-quality-gate"] = "false"
	}
	return p
}

// requiredServices takes a variable amount of services and removes
// SonarQube from the resulting slice if the skipSonarQubeFlag is set.
func requiredServices(s ...tasktesting.Service) []tasktesting.Service {
	requiredServices := []tasktesting.Service{tasktesting.Nexus}
	sqIndex := slices.Index(requiredServices, tasktesting.SonarQube)
	if sqIndex != -1 && *skipSonarQubeFlag {
		requiredServices = slices.Delete(requiredServices, sqIndex, sqIndex+1)
	}
	return requiredServices
}

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

func checkFilesExist(t *testing.T, wsDir string, wantFiles ...string) {
	for _, wf := range wantFiles {
		filename := filepath.Join(wsDir, wf)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Fatalf("Want %s, but got nothing", filename)
		}
	}
}

func checkFileHash(t *testing.T, wsDir string, filename string, hash [32]byte) {
	filepath := filepath.Join(wsDir, filename)
	filecontent, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Want %s, but got nothing", filename)
	}
	filehash := sha256.Sum256(filecontent)
	if filehash != hash {
		t.Fatalf("Want %x, but got %x", hash, filehash)
	}
}

func getTrimmedFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
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

func checkFileContentContains(t *testing.T, wsDir, filename string, wantContains ...string) {
	content, err := os.ReadFile(filepath.Join(wsDir, filename))
	got := string(content)
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	for _, w := range wantContains {
		if !strings.Contains(got, w) {
			t.Fatalf("got '%s', want '%s' contained in file %s", got, w, filename)
		}
	}
}

func checkFileContentLeanContains(t *testing.T, wsDir, filename string, wantContains string) {
	got, err := getFileContentLean(filepath.Join(wsDir, filename))
	if err != nil {
		t.Fatalf("could not read %s: %s", filename, err)
	}
	if !strings.Contains(got, wantContains) {
		t.Fatalf("got '%s', want '%s' contained in file %s", got, wantContains, filename)
	}
}

func getFileContentLean(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	contentStr := strings.ReplaceAll(string(content), "\t", "")
	contentStr = strings.ReplaceAll(contentStr, "\n", "")
	contentStr = strings.ReplaceAll(contentStr, " ", "")

	return contentStr, nil
}

func runTaskTestCases(t *testing.T, taskName string, requiredServices []tasktesting.Service, testCases map[string]tasktesting.TestCase) {
	tasktesting.CheckCluster(t, *outsideKindFlag)
	if len(requiredServices) != 0 {
		tasktesting.CheckServices(t, requiredServices)
	}

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
			tn := taskName
			if tc.Timeout == 0 {
				tc.Timeout = 5 * time.Minute
			}
			tasktesting.Run(t, tc, tasktesting.TestOpts{
				TaskKindRef:             taskKindRef,
				TaskName:                tn,
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

	sonarClient, err := sonar.NewClient(&sonar.ClientConfig{
		APIToken:      sonarToken,
		BaseURL:       "http://localhost:9000", // use localhost instead of ods-test-sonarqube.kind!
		ServerEdition: "community",
	})
	if err != nil {
		t.Fatalf("sonar client: %s", err)
	}

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
	return os.WriteFile(filename, y, 0644)
}

func checkBuildStatus(t *testing.T, c *bitbucket.Client, gitCommit, wantBuildStatus string) {
	buildStatusPage, err := c.BuildStatusList(gitCommit)
	buildStatus := buildStatusPage.Values[0]
	if err != nil {
		t.Fatal(err)
	}
	if buildStatus.State != wantBuildStatus {
		t.Fatalf("Got: %s, want: %s", buildStatus.State, wantBuildStatus)
	}
}

func createAppInSubDirectory(t *testing.T, wsDir string, subdir string, sampleApp string) {
	err := os.MkdirAll(filepath.Join(wsDir, subdir), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = directory.Copy(
		filepath.Join(projectpath.Root, "test", tasktesting.TestdataWorkspacesPath, sampleApp),
		filepath.Join(wsDir, subdir),
	)
	if err != nil {
		t.Fatal(err)
	}
}
