package manager

import (
	"context"
	"crypto/sha1"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	tektonClient "github.com/opendevstack/ods-pipeline/internal/tekton"
	"github.com/opendevstack/ods-pipeline/pkg/config"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/clock"
)

const (
	// Label prefix to use for labels applied by this service.
	labelPrefix = "pipeline.opendevstack.org/"
	// Label specifying the Bitbucket repository related to the pipeline.
	repositoryLabel = labelPrefix + "repository"
	// Label specifying the Git ref (e.g. branch) related to the pipeline.
	gitRefLabel = labelPrefix + "git-ref"
	// tektonAPIVersion specifies the Tekton API version in use
	tektonAPIVersion = "tekton.dev/v1beta1"
	// sharedWorkspaceName is the name of the workspace shared by all tasks
	sharedWorkspaceName = "shared-workspace"
)

// PipelineConfig holds configuration for a triggered pipeline.
type PipelineConfig struct {
	PipelineInfo
	PVC          string
	PipelineSpec config.Pipeline
	Params       []tekton.Param
}

// createPipelineRun creates a PipelineRun resource
func createPipelineRun(
	tektonClient tektonClient.ClientPipelineRunInterface,
	ctxt context.Context,
	cfg PipelineConfig,
	taskKind tekton.TaskKind,
	taskSuffix string,
	needQueueing bool) (*tekton.PipelineRun, error) {
	pr := &tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cfg.Component + "-",
			Labels:       pipelineLabels(cfg),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: tektonAPIVersion,
			Kind:       "PipelineRun",
		},
		Spec: tekton.PipelineRunSpec{
			PipelineSpec:       assemblePipelineSpec(cfg, taskKind, taskSuffix),
			Params:             extractPipelineParams(cfg.Params),
			ServiceAccountName: "pipeline", // TODO
			PodTemplate:        cfg.PipelineSpec.PodTemplate,
			TaskRunSpecs:       cfg.PipelineSpec.TaskRunSpecs,
			Timeouts:           cfg.PipelineSpec.Timeouts,
			Workspaces: []tekton.WorkspaceBinding{
				{
					Name: sharedWorkspaceName,
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: cfg.PVC,
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
	}
}

// assemblePipelineSpec returns a Tekton pipeline based on given PipelineConfig.
func assemblePipelineSpec(cfg PipelineConfig, taskKind tekton.TaskKind, taskSuffix string) *tekton.PipelineSpec {
	var tasks []tekton.PipelineTask
	tasks = append(tasks, tekton.PipelineTask{
		Name:       "start",
		TaskRef:    &tekton.TaskRef{Kind: taskKind, Name: "ods-start" + taskSuffix},
		Params:     startTaskParams(),
		Workspaces: tektonDefaultWorkspaceBindings(),
	})
	if len(cfg.PipelineSpec.Tasks) > 0 {
		// Add "start" to runAfter of the first configured task, and to each further task
		// that does not set runAfter until we hit a task that does.
		for i := range cfg.PipelineSpec.Tasks {
			if i > 0 && len(cfg.PipelineSpec.Tasks[i].RunAfter) > 0 {
				break
			}
			cfg.PipelineSpec.Tasks[i].RunAfter = append(cfg.PipelineSpec.Tasks[i].RunAfter, "start")
		}
		tasks = append(tasks, cfg.PipelineSpec.Tasks...)
	}
	tasks = mergeTriggerBasedParams(tasks, cfg.Params)

	finallyTasks := append([]tekton.PipelineTask{}, cfg.PipelineSpec.Finally...)
	finallyTasks = append(finallyTasks, tekton.PipelineTask{
		Name:       "finish",
		TaskRef:    &tekton.TaskRef{Kind: taskKind, Name: "ods-finish" + taskSuffix},
		Workspaces: tektonDefaultWorkspaceBindings(),
		Params:     finishTaskParams(),
	})
	finallyTasks = mergeTriggerBasedParams(finallyTasks, cfg.Params)

	return &tekton.PipelineSpec{
		Params: []tekton.ParamSpec{
			tektonStringParamSpec("repository", cfg.Repository),
			tektonStringParamSpec("project", cfg.Project),
			tektonStringParamSpec("component", cfg.Component),
			tektonStringParamSpec("git-repo-url", cfg.GitURI),
			tektonStringParamSpec("git-full-ref", cfg.GitFullRef),
			tektonStringParamSpec("pr-key", strconv.Itoa(cfg.PullRequestKey)),
			tektonStringParamSpec("pr-base", cfg.PullRequestBase),
		},
		Tasks: tasks,
		Workspaces: []tekton.PipelineWorkspaceDeclaration{
			{Name: sharedWorkspaceName},
		},
		Finally: finallyTasks,
	}
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
	return !(pr.IsDone() || pr.IsPending() || pr.IsCancelled() || pr.HasTimedOut(context.Background(), clock.RealClock{}))
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

// startTaskParams returns the params for the start task.
func startTaskParams() []tekton.Param {
	return []tekton.Param{
		tektonStringParam("url", "$(params.git-repo-url)"),
		tektonStringParam("git-full-ref", "$(params.git-full-ref)"),
		tektonStringParam("project", "$(params.project)"),
		tektonStringParam("pr-key", "$(params.pr-key)"),
		tektonStringParam("pr-base", "$(params.pr-base)"),
		tektonStringParam("pipeline-run-name", "$(context.pipelineRun.name)"),
		tektonStringParam("version", "$(params.version)"),
	}
}

// startTaskParams returns the params for the finish task.
func finishTaskParams() []tekton.Param {
	return []tekton.Param{
		tektonStringParam("pipeline-run-name", "$(context.pipelineRun.name)"),
		tektonStringParam("aggregate-tasks-status", "$(tasks.status)"),
	}
}

// extractPipelineParams returns only those params which are not prefixed.
func extractPipelineParams(params []tekton.Param) (matching []tekton.Param) {
	for _, p := range params {
		if !strings.Contains(p.Name, ".") {
			matching = append(matching, p)
		}
	}
	return
}

// extractTaskParams returns only those params which are prefixed with taskName.
// The returned param names do not include the prefix.
func extractTaskParams(taskName string, params []tekton.Param) (matching []tekton.Param) {
	for _, p := range params {
		if strings.HasPrefix(p.Name, taskName+".") {
			p.Name = strings.TrimPrefix(p.Name, taskName+".")
			matching = append(matching, p)
		}
	}
	return
}

// mergeTriggerBasedParams appends given params to the tasks' params.
// If the given params contain a param of the same name as an existing param,
// the existing param will be overriden.
func mergeTriggerBasedParams(tasks []tekton.PipelineTask, params []tekton.Param) (extended []tekton.PipelineTask) {
	for _, t := range tasks {
		var mergedParams []tekton.Param
		extractedParams := extractTaskParams(t.Name, params)
		for _, originalParam := range t.Params {
			if !containsParam(extractedParams, originalParam) {
				mergedParams = append(mergedParams, originalParam)
			}
		}
		t.Params = append(mergedParams, extractedParams...)
		extended = append(extended, t)
	}
	return
}

// containsParam checks whether params contains a param with the same name as param's name.
func containsParam(params []tekton.Param, param tekton.Param) bool {
	for _, p := range params {
		if p.Name == param.Name {
			return true
		}
	}
	return false
}
