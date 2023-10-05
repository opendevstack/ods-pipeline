package manager

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	tektonClient "github.com/opendevstack/ods-pipeline/internal/tekton"
	"github.com/opendevstack/ods-pipeline/pkg/config"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
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
		},
		PVC: "pvc",
	}
	t.Run("non-queued PR", func(t *testing.T) {
		pr, err := createPipelineRun(tc, ctxt, pData, false)
		if err != nil {
			t.Fatal(err)
		}
		if pr.GenerateName != "component-" {
			t.Errorf("Expected generated name to be component-, got: %s", pr.GenerateName)
		}
		if pr.Spec.Status != "" {
			t.Errorf("Expected status to be empty, got: %s", pr.Spec.Status)
		}
		if pr.Labels[repositoryLabel] != pData.Repository {
			t.Errorf("Expected label %s to be %s, got: %s", repositoryLabel, pData.Repository, pr.Labels[repositoryLabel])
		}
		if pr.Labels[gitRefLabel] != pData.GitRef {
			t.Errorf("Expected label %s to be %s, got: %s", gitRefLabel, pData.GitRef, pr.Labels[gitRefLabel])
		}
		workspaceCfg := pr.Spec.Workspaces[0]
		if workspaceCfg.Name != sharedWorkspaceName {
			t.Errorf("Expected generated name to be %s, got: %s", sharedWorkspaceName, workspaceCfg.Name)
		}
		if workspaceCfg.PersistentVolumeClaim.ClaimName != "pvc" {
			t.Errorf("Expected generated name to be pvc, got: %s", workspaceCfg.Name)
		}
		if len(tc.CreatedPipelineRuns) != 1 {
			t.Error("No pipeline run created")
		}
	})

	t.Run("pending PR", func(t *testing.T) {
		pr, err := createPipelineRun(tc, ctxt, pData, true)
		if err != nil {
			t.Fatal(err)
		}
		if pr.Spec.Status != tekton.PipelineRunSpecStatusPending {
			t.Errorf("Expected status to be pending, got: %s", pr.Spec.Status)
		}
		if len(tc.CreatedPipelineRuns) != 2 {
			t.Error("No pipeline run created")
		}
	})

	t.Run("with spec", func(t *testing.T) {
		pData.Params = []tekton.Param{
			tektonStringParam("hello", "world"),
			tektonStringParam("start.clone-depth", "5"),
			tektonStringParam("foo.bar", "baz"),
			tektonStringParam("finish.aggregate-tasks-status", "overriden"),
		}
		pData.PipelineSpec.Tasks = []tekton.PipelineTask{
			{
				Name:    "foo",
				TaskRef: &tekton.TaskRef{Kind: "Task", Name: "foo"},
				Params: []tekton.Param{
					tektonStringParam("some", "value"),
				},
			},
		}
		pr, err := createPipelineRun(tc, ctxt, pData, false)
		if err != nil {
			t.Fatal(err)
		}
		wantParams := tekton.Params{
			{Name: "hello", Value: tekton.ParamValue{Type: "string", StringVal: "world"}},
		}
		if diff := cmp.Diff(wantParams, pr.Spec.Params); diff != "" {
			t.Fatalf("expected params (-want +got):\n%s", diff)
		}
		wantTasks := []tekton.PipelineTask{
			{
				Name:       "start",
				TaskRef:    &tekton.TaskRef{Kind: "Task", Name: "ods-pipeline-start"},
				Params:     append(startTaskParams(), tektonStringParam("clone-depth", "5")),
				Workspaces: tektonDefaultWorkspaceBindings(),
			},
			{
				Name:    "foo",
				TaskRef: &tekton.TaskRef{Kind: "Task", Name: "foo"},
				Params: []tekton.Param{
					tektonStringParam("some", "value"),
					tektonStringParam("bar", "baz"),
				},
				RunAfter: []string{"start"},
			},
		}
		if diff := cmp.Diff(wantTasks, pr.Spec.PipelineSpec.Tasks); diff != "" {
			t.Fatalf("expected tasks (-want +got):\n%s", diff)
		}
		wantFinallyTasks := []tekton.PipelineTask{
			{
				Name:    "finish",
				TaskRef: &tekton.TaskRef{Kind: "Task", Name: "ods-pipeline-finish"},
				Params: []tekton.Param{
					tektonStringParam("pipeline-run-name", "$(context.pipelineRun.name)"),
					tektonStringParam("aggregate-tasks-status", "overriden"),
				},
				Workspaces: tektonDefaultWorkspaceBindings(),
			},
		}
		if diff := cmp.Diff(wantFinallyTasks, pr.Spec.PipelineSpec.Finally); diff != "" {
			t.Fatalf("expected finally (-want +got):\n%s", diff)
		}
	})
}

func TestAssemblePipeline(t *testing.T) {
	taskKind := tekton.NamespacedTaskKind
	cfg := PipelineConfig{
		PipelineInfo: PipelineInfo{
			Project:         "project",
			Component:       "component",
			Repository:      "repo",
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
		PipelineSpec: config.Pipeline{
			Tasks: []tekton.PipelineTask{
				{
					Name:    "build",
					TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-pipeline-go-build"},
					Workspaces: []tekton.WorkspacePipelineTaskBinding{
						{Name: "source", Workspace: sharedWorkspaceName},
					},
				},
			},
			Finally: []tekton.PipelineTask{
				{
					Name:    "final",
					TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "final"},
				},
			},
		},
	}
	got := assemblePipelineSpec(cfg)
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
		},
		Tasks: []tekton.PipelineTask{
			{
				Name:    "start",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-pipeline-start"},
				Params: []tekton.Param{
					tektonStringParam("url", "$(params.git-repo-url)"),
					tektonStringParam("git-full-ref", "$(params.git-full-ref)"),
					tektonStringParam("project", "$(params.project)"),
					tektonStringParam("pr-key", "$(params.pr-key)"),
					tektonStringParam("pr-base", "$(params.pr-base)"),
					tektonStringParam("pipeline-run-name", "$(context.pipelineRun.name)"),
					tektonStringParam("version", "$(params.version)"),
				},
				Workspaces: tektonDefaultWorkspaceBindings(),
			},
			{
				Name:       "build",
				RunAfter:   []string{"start"},
				TaskRef:    &tekton.TaskRef{Kind: taskKind, Name: "ods-pipeline-go-build"},
				Params:     nil,
				Workspaces: tektonDefaultWorkspaceBindings(),
			},
		},
		Finally: []tekton.PipelineTask{
			{
				Name:    "final",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "final"},
				Params:  nil,
			},
			{
				Name:    "finish",
				TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-pipeline-finish"},
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

func TestTasksRunAfterInjection(t *testing.T) {
	tests := map[string]struct {
		cfgTasks []tekton.PipelineTask
		want     []tekton.PipelineTask
	}{
		"one build task": {
			cfgTasks: []tekton.PipelineTask{
				{Name: "build"},
				{Name: "package-image", RunAfter: []string{"build"}},
			},
			want: []tekton.PipelineTask{
				{Name: "start"},
				{Name: "build", RunAfter: []string{"start"}},
				{Name: "package-image", RunAfter: []string{"build"}},
			},
		},
		"parallel build tasks": {
			cfgTasks: []tekton.PipelineTask{
				{Name: "build-one"},
				{Name: "build-two"},
				{Name: "package-image", RunAfter: []string{"build-one", "build-two"}},
			},
			want: []tekton.PipelineTask{
				{Name: "start"},
				{Name: "build-one", RunAfter: []string{"start"}},
				{Name: "build-two", RunAfter: []string{"start"}},
				{Name: "package-image", RunAfter: []string{"build-one", "build-two"}},
			},
		},
		"no configured tasks": {
			cfgTasks: []tekton.PipelineTask{},
			want: []tekton.PipelineTask{
				{Name: "start"},
			},
		},
		"only parallel tasks": {
			cfgTasks: []tekton.PipelineTask{
				{Name: "build-one"},
				{Name: "build-two"},
			},
			want: []tekton.PipelineTask{
				{Name: "start"},
				{Name: "build-one", RunAfter: []string{"start"}},
				{Name: "build-two", RunAfter: []string{"start"}},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cfg := PipelineConfig{PipelineSpec: config.Pipeline{Tasks: tc.cfgTasks}}
			got := assemblePipelineSpec(cfg)
			wantRunAfter := [][]string{}
			for _, task := range tc.want {
				wantRunAfter = append(wantRunAfter, task.RunAfter)
			}
			gotRunAfter := [][]string{}
			for _, task := range got.Tasks {
				gotRunAfter = append(gotRunAfter, task.RunAfter)
			}
			if diff := cmp.Diff(wantRunAfter, gotRunAfter); diff != "" {
				t.Fatalf("expected (-want +got):\n%s", diff)
			}
		})
	}

}

func TestExtractTaskParams(t *testing.T) {
	taskName := "foo"
	params := []tekton.Param{
		tektonStringParam("one", "a"),
		tektonStringParam("foo.two", "b"),
		tektonStringParam("foobar.three", "c"),
		tektonStringParam("foo.four", "d"),
	}
	want := []tekton.Param{
		tektonStringParam("two", "b"),
		tektonStringParam("four", "d"),
	}
	got := extractTaskParams(taskName, params)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("expected (-want +got):\n%s", diff)
	}
}

func TestAppendTriggerBasedParams(t *testing.T) {
	var tasks []tekton.PipelineTask
	params := []tekton.Param{
		tektonStringParam("one", "a"),
		tektonStringParam("foo.two", "b"),
		tektonStringParam("foobar.three", "c"),
		tektonStringParam("foo.four", "d"),
	}
	tasks = append(tasks, tekton.PipelineTask{
		Name: "foo",
		Params: []tekton.Param{
			tektonStringParam("zero", "0"),
			tektonStringParam("four", "should be overriden"),
		},
	})
	got := mergeTriggerBasedParams(tasks, params)
	want := tekton.Params{
		tektonStringParam("zero", "0"),
		tektonStringParam("two", "b"),
		tektonStringParam("four", "d"),
	}
	if diff := cmp.Diff(want, got[0].Params); diff != "" {
		t.Fatalf("expected (-want +got):\n%s", diff)
	}
}
