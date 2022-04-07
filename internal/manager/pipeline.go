package manager

import (
	"context"
	"crypto/sha1"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	// Label prefix to use for labels applied by this service.
	labelPrefix = "pipeline.opendevstack.org/"
	// Label specifying the Bitbucket repository related to the pipeline.
	repositoryLabel = labelPrefix + "repository"
	// Label specifying the Git ref (e.g. branch) related to the pipeline.
	gitRefLabel = labelPrefix + "git-ref"
	// Label specifying the target stage of the pipeline.
	stageLabel = labelPrefix + "stage"
	// tektonAPIVersion specifies the Tekton API version in use
	tektonAPIVersion = "tekton.dev/v1beta1"
	// sharedWorkspaceName is the name of the workspace shared by all tasks
	sharedWorkspaceName = "shared-workspace"
)

// PipelineConfig holds configuration for a triggered pipeline.
type PipelineConfig struct {
	PipelineInfo
	PVC     string `json:"pvc"`
	Tasks   []tekton.PipelineTask
	Finally []tekton.PipelineTask
}

// createPipelineRun creates a PipelineRun resource
func createPipelineRun(tektonClient tektonClient.ClientPipelineRunInterface, ctxt context.Context, pData PipelineConfig, needQueueing bool) (*tekton.PipelineRun, error) {
	pr := &tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", pData.Name),
			Labels:       pipelineLabels(pData),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: tektonAPIVersion,
			Kind:       "PipelineRun",
		},
		Spec: tekton.PipelineRunSpec{
			PipelineRef:        &tekton.PipelineRef{Name: pData.Name},
			ServiceAccountName: "pipeline", // TODO
			Workspaces: []tekton.WorkspaceBinding{
				{
					Name: sharedWorkspaceName,
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: pData.PVC,
					},
				},
			},
		},
	}
	if needQueueing {
		pr.Spec.Status = tekton.PipelineRunSpecStatusPending
	}
	return tektonClient.CreatePipelineRun(ctxt, pr, metav1.CreateOptions{})
}

// listPipelineRuns lists pipeline runs associated with repository.
func listPipelineRuns(ctxt context.Context, tektonClient tektonClient.ClientPipelineRunInterface, repository string) (*tekton.PipelineRunList, error) {
	labelMap := map[string]string{repositoryLabel: repository}
	return tektonClient.ListPipelineRuns(
		ctxt, metav1.ListOptions{LabelSelector: labels.Set(labelMap).String()},
	)
}

// makePipelineName generates the name of the pipeline.
// According to the Kubernetes label rules, a maximum of 63 characters is
// allowed, see https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set.
// Therefore, the name might be truncated. As this could cause potential clashes
// between similar named branches, we put a short part of the branch hash value
// into the name name to make this very unlikely.
// We cut the pipeline name at 55 chars to allow e.g. pipeline runs to add suffixes.
func makePipelineName(component string, branch string) string {
	// 55 is derived from K8s label max length minus room for generateName suffix.
	return makeValidLabelValue(component+"-", branch, 55)
}

func makeValidLabelValue(prefix, branch string, maxLength int) string {
	// Cut all non-alphanumeric characters
	safeCharsRegex := regexp.MustCompile("[^-a-zA-Z0-9]+")
	result := prefix + safeCharsRegex.ReplaceAllString(
		strings.Replace(branch, "/", "-", -1),
		"",
	)
	result = fitStringToMaxLength(result, maxLength)
	return strings.ToLower(result)
}

// fitStringToMaxLength ensures s is not longer than max.
// If s is longer than max, it shortenes s and appends a unique, consistent
// suffix so that multiple invocations produce the same result. The length
// of the shortened string will be equal to max.
func fitStringToMaxLength(s string, max int) string {
	if len(s) <= max {
		return s
	}
	suffixLength := 7
	shortened := s[0 : max-suffixLength-1]
	h := sha1.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		return shortened
	}
	bs := h.Sum(nil)
	suffix := fmt.Sprintf("%x", bs)
	return fmt.Sprintf("%s-%s", shortened, suffix[0:suffixLength])
}

// pipelineLabels returns a map of labels to apply to pipelines and related runs.
func pipelineLabels(data PipelineConfig) map[string]string {
	return map[string]string{
		repositoryLabel: data.Repository,
		gitRefLabel:     makeValidLabelValue("", data.GitRef, 63),
		stageLabel:      data.Stage,
	}
}

// assemblePipeline returns a Tekton pipeline based on given PipelineConfig.
func assemblePipeline(cfg PipelineConfig, taskKind tekton.TaskKind, taskSuffix string) *tekton.Pipeline {
	var tasks []tekton.PipelineTask
	tasks = append(tasks, tekton.PipelineTask{
		Name:       "ods-start",
		TaskRef:    &tekton.TaskRef{Kind: taskKind, Name: "ods-start" + taskSuffix},
		Workspaces: tektonDefaultWorkspaceBindings(),
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
	})
	if len(cfg.Tasks) > 0 {
		cfg.Tasks[0].RunAfter = append(cfg.Tasks[0].RunAfter, "ods-start")
		tasks = append(tasks, cfg.Tasks...)
	}

	var finallyTasks []tekton.PipelineTask
	finallyTasks = append(finallyTasks, cfg.Finally...)

	finallyTasks = append(finallyTasks, tekton.PipelineTask{
		Name:       "ods-finish",
		TaskRef:    &tekton.TaskRef{Kind: taskKind, Name: "ods-finish" + taskSuffix},
		Workspaces: tektonDefaultWorkspaceBindings(),
		Params: []tekton.Param{
			tektonStringParam("pipeline-run-name", "$(context.pipelineRun.name)"),
			tektonStringParam("aggregate-tasks-status", "$(tasks.status)"),
		},
	})

	p := &tekton.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cfg.Name,
			Labels: pipelineLabels(cfg),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: tektonAPIVersion,
			Kind:       "Pipeline",
		},
		Spec: tekton.PipelineSpec{
			Description: "ODS",
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
			Tasks: tasks,
			Workspaces: []tekton.PipelineWorkspaceDeclaration{
				{Name: sharedWorkspaceName},
			},
			Finally: finallyTasks,
		},
	}
	return p
}

// sortPipelineRunsDescending sorts pipeline runs by time (descending)
func sortPipelineRunsDescending(pipelineRuns []tekton.PipelineRun) {
	sort.Slice(pipelineRuns, func(i, j int) bool {
		return pipelineRuns[j].CreationTimestamp.Time.Before(pipelineRuns[i].CreationTimestamp.Time)
	})
}

// pipelineRunIsProgressing returns true if the PR is not done, not pending,
// not cancelled, and not timed out.
func pipelineRunIsProgressing(pr tekton.PipelineRun) bool {
	return !(pr.IsDone() || pr.IsPending() || pr.IsCancelled() || pr.IsTimedOut())
}

// tektonStringParam returns a Tekton task parameter.
func tektonStringParam(name, val string) tekton.Param {
	return tekton.Param{Name: name, Value: tekton.ArrayOrString{Type: "string", StringVal: val}}
}

// tektonStringParam returns a Tekton task parameter spec.
func tektonStringParamSpec(name, defaultVal string) tekton.ParamSpec {
	return tekton.ParamSpec{
		Name: name,
		Type: "string",
		Default: &tekton.ArrayOrString{
			Type: tekton.ParamTypeString, StringVal: defaultVal,
		}}
}

// tektonDefaultWorkspaceBindings returns the default workspace bindings for a task.
func tektonDefaultWorkspaceBindings() []tekton.WorkspacePipelineTaskBinding {
	return []tekton.WorkspacePipelineTaskBinding{
		{Name: "source", Workspace: sharedWorkspaceName},
	}
}
