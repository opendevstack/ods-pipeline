package manager

import (
	"context"
	"testing"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/config"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func TestCreatePipelineRun(t *testing.T) {
	tc := &tektonClient.TestClient{}
	ctxt := context.TODO()
	pData := PipelineData{
		Name:       "foo",
		Repository: "repo",
		GitRef:     "branch",
		Stage:      config.DevStage,
		PVC:        "pvc",
	}
	pr, err := createPipelineRun(tc, ctxt, pData, false)
	if err != nil {
		t.Fatal(err)
	}
	if pr.GenerateName != "foo-" {
		t.Fatalf("Expected generated name to be foo-, got: %s", pr.GenerateName)
	}
	if pr.Spec.PipelineRef.Name != "foo" {
		t.Fatalf("Expected pipeline ref to be foo, got: %s", pr.Spec.PipelineRef.Name)
	}
	if pr.Spec.Status != "" {
		t.Fatalf("Expected status to be empty, got: %s", pr.Spec.Status)
	}
	if pr.Labels[repositoryLabel] != pData.Repository {
		t.Fatalf("Expected label %s to be %s, got: %s", repositoryLabel, pData.Repository, pr.Labels[repositoryLabel])
	}
	if pr.Labels[gitRefLabel] != pData.GitRef {
		t.Fatalf("Expected label %s to be %s, got: %s", gitRefLabel, pData.GitRef, pr.Labels[gitRefLabel])
	}
	if pr.Labels[stageLabel] != pData.Stage {
		t.Fatalf("Expected label %s to be %s, got: %s", stageLabel, pData.Stage, pr.Labels[stageLabel])
	}
	workspaceCfg := pr.Spec.Workspaces[0]
	if workspaceCfg.Name != sharedWorkspaceName {
		t.Fatalf("Expected generated name to be %s, got: %s", sharedWorkspaceName, workspaceCfg.Name)
	}
	if workspaceCfg.PersistentVolumeClaim.ClaimName != "pvc" {
		t.Fatalf("Expected generated name to be pvc, got: %s", workspaceCfg.Name)
	}
	if len(tc.CreatedPipelineRuns) != 1 {
		t.Fatal("No pipeline run created")
	}
	pr, err = createPipelineRun(tc, ctxt, pData, true)
	if err != nil {
		t.Fatal(err)
	}
	if pr.Spec.Status != tekton.PipelineRunSpecStatusPending {
		t.Fatalf("Expected status to be pending, got: %s", pr.Spec.Status)
	}
	if len(tc.CreatedPipelineRuns) != 2 {
		t.Fatal("No pipeline run created")
	}
}
