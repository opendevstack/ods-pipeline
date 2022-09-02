package manager

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/config"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func TestShortenString(t *testing.T) {
	tests := map[string]struct {
		s        string
		max      int
		expected string
	}{
		"short enough": {
			s:        "foobar",
			max:      10,
			expected: "foobar",
		},
		"too long": {
			s:        "some-arbitarily-long-name-that-should-be-way-shorter",
			max:      30,
			expected: "some-arbitarily-long-n-8b85b7c",
		},
		"too long with slight difference in cut off string": {
			s:        "some-arbitarily-long-name-that-should-be-way-shorterx",
			max:      30,
			expected: "some-arbitarily-long-n-50a3b84",
		},
		"exact length": {
			s:        "some-arbitarily-long-name-that",
			max:      30,
			expected: "some-arbitarily-long-name-that",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := fitStringToMaxLength(tc.s, tc.max)
			if tc.expected != got {
				t.Fatalf(
					"Want '%s', got '%s' for (s='%s', max='%d')",
					tc.expected,
					got,
					tc.s,
					tc.max,
				)
			}
		})
	}
}

func TestCreatePipelineRun(t *testing.T) {
	tc := &tektonClient.TestClient{}
	ctxt := context.TODO()
	pData := PipelineConfig{
		PipelineInfo: PipelineInfo{
			Component:  "component",
			Repository: "project-component",
			GitRef:     "branch",
			Stage:      config.DevStage,
		},
		PVC: "pvc",
	}
	pr, err := createPipelineRun(tc, ctxt, pData, tekton.ClusterTaskKind, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if pr.GenerateName != "component-" {
		t.Fatalf("Expected generated name to be component-, got: %s", pr.GenerateName)
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
	pr, err = createPipelineRun(tc, ctxt, pData, tekton.NamespacedTaskKind, "", true)
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

func TestAssemblePipeline(t *testing.T) {
	taskKind := tekton.NamespacedTaskKind
	taskSuffix := "-latest"
	cfg := PipelineConfig{
		PipelineInfo: PipelineInfo{
			Project:         "project",
			Component:       "component",
			Repository:      "repo",
			Stage:           config.DevStage,
			Environment:     "env",
			Version:         "1.0.0",
			GitRef:          "branch",
			GitFullRef:      "refs/heads/branch",
			GitSHA:          "6621c6060715428933a1d20851e0d51614b0a195",
			RepoBase:        "https://bitbucket.example.com/scm",
			GitURI:          "https://bitbucket.example.com/scm/project/repo.git",
			Namespace:       "namespace",
			TriggerEvent:    "repo:refs_changed",
			Comment:         "comment",
			PullRequestKey:  1,
			PullRequestBase: "integration",
		},
		PVC: "pvc",
		Tasks: []tekton.PipelineTask{
			{
				Name:    "build",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-build-go" + taskSuffix},
				Workspaces: []tekton.WorkspacePipelineTaskBinding{
					{Name: "source", Workspace: sharedWorkspaceName},
				},
			},
		},
		Finally: []tekton.PipelineTask{
			{
				Name:    "final",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "final" + taskSuffix},
			},
		},
	}
	got := assemblePipelineSpec(cfg, taskKind, taskSuffix)
	want := &tekton.PipelineSpec{
		Description: "",
		Params: []tekton.ParamSpec{
			tektonStringParamSpec("repository", cfg.Repository),
			tektonStringParamSpec("project", cfg.Project),
			tektonStringParamSpec("component", cfg.Component),
			tektonStringParamSpec("git-repo-url", cfg.GitURI),
			tektonStringParamSpec("git-full-ref", cfg.GitFullRef),
			tektonStringParamSpec("pr-key", strconv.Itoa(cfg.PullRequestKey)),
			tektonStringParamSpec("pr-base", cfg.PullRequestBase),
			tektonStringParamSpec("environment", cfg.Environment),
			tektonStringParamSpec("version", cfg.Version),
		},
		Tasks: []tekton.PipelineTask{
			{
				Name:    "start",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-start-latest"},
				Params: []tekton.Param{
					tektonStringParam("url", "$(params.git-repo-url)"),
					tektonStringParam("git-full-ref", "$(params.git-full-ref)"),
					tektonStringParam("project", "$(params.project)"),
					tektonStringParam("pr-key", "$(params.pr-key)"),
					tektonStringParam("pr-base", "$(params.pr-base)"),
					tektonStringParam("pipeline-run-name", "$(context.pipelineRun.name)"),
					tektonStringParam("environment", "$(params.environment)"),
					tektonStringParam("version", "$(params.version)"),
				},
				Workspaces: tektonDefaultWorkspaceBindings(),
			},
			{
				Name:       "build",
				RunAfter:   []string{"start"},
				TaskRef:    &tekton.TaskRef{Kind: taskKind, Name: "ods-build-go-latest"},
				Params:     nil,
				Workspaces: tektonDefaultWorkspaceBindings(),
			},
		},
		Finally: []tekton.PipelineTask{
			{
				Name:    "final",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "final-latest"},
				Params:  nil,
			},
			{
				Name:    "finish",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-finish-latest"},
				Params: []tekton.Param{
					tektonStringParam("pipeline-run-name", "$(context.pipelineRun.name)"),
					tektonStringParam("aggregate-tasks-status", "$(tasks.status)"),
				},
				Workspaces: tektonDefaultWorkspaceBindings(),
			},
		},
		Workspaces: []tekton.PipelineWorkspaceDeclaration{
			{Name: sharedWorkspaceName},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("expected (-want +got):\n%s", diff)
	}
}
