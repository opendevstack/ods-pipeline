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
	"github.com/opendevstack/pipeline/pkg/config"
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

// createPipelineRun creates a PipelineRun resource
func createPipelineRun(tektonClient tektonClient.ClientPipelineRunInterface, ctxt context.Context, pData PipelineData, needQueueing bool) (*tekton.PipelineRun, error) {
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
func listPipelineRuns(tektonClient tektonClient.ClientPipelineRunInterface, ctxt context.Context, repository string) (*tekton.PipelineRunList, error) {
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
func pipelineLabels(data PipelineData) map[string]string {
	return map[string]string{
		repositoryLabel: data.Repository,
		gitRefLabel:     makeValidLabelValue("", data.GitRef, 63),
		stageLabel:      data.Stage,
	}
}

func assemblePipeline(odsConfig *config.ODS, data PipelineData, taskKind tekton.TaskKind, taskSuffix string) *tekton.Pipeline {

	var tasks []tekton.PipelineTask
	tasks = append(tasks, tekton.PipelineTask{
		Name:    "ods-start",
		TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-start" + taskSuffix},
		Workspaces: []tekton.WorkspacePipelineTaskBinding{
			{Name: "source", Workspace: sharedWorkspaceName},
		},
		Params: []tekton.Param{
			{
				Name: "url",
				Value: tekton.ArrayOrString{
					StringVal: "$(params.git-repo-url)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "git-full-ref",
				Value: tekton.ArrayOrString{
					StringVal: "$(params.git-full-ref)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "project",
				Value: tekton.ArrayOrString{
					StringVal: "$(params.project)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "pr-key",
				Value: tekton.ArrayOrString{
					StringVal: "$(params.pr-key)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "pr-base",
				Value: tekton.ArrayOrString{
					StringVal: "$(params.pr-base)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "pipeline-run-name",
				Value: tekton.ArrayOrString{
					StringVal: "$(context.pipelineRun.name)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "environment",
				Value: tekton.ArrayOrString{
					StringVal: "$(params.environment)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "version",
				Value: tekton.ArrayOrString{
					StringVal: "$(params.version)",
					Type:      tekton.ParamTypeString,
				},
			},
		},
	})
	if len(odsConfig.Pipeline.Tasks) > 0 {
		odsConfig.Pipeline.Tasks[0].RunAfter = append(odsConfig.Pipeline.Tasks[0].RunAfter, "ods-start")
		tasks = append(tasks, odsConfig.Pipeline.Tasks...)
	}

	var finallyTasks []tekton.PipelineTask
	finallyTasks = append(finallyTasks, odsConfig.Pipeline.Finally...)

	finallyTasks = append(finallyTasks, tekton.PipelineTask{
		Name:    "ods-finish",
		TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-finish" + taskSuffix},
		Workspaces: []tekton.WorkspacePipelineTaskBinding{
			{Name: "source", Workspace: sharedWorkspaceName},
		},
		Params: []tekton.Param{
			{
				Name: "pipeline-run-name",
				Value: tekton.ArrayOrString{
					StringVal: "$(context.pipelineRun.name)",
					Type:      tekton.ParamTypeString,
				},
			},
			{
				Name: "aggregate-tasks-status",
				Value: tekton.ArrayOrString{
					StringVal: "$(tasks.status)",
					Type:      tekton.ParamTypeString,
				},
			},
		},
	})

	p := &tekton.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: pipelineLabels(data),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: tektonAPIVersion,
			Kind:       "Pipeline",
		},
		Spec: tekton.PipelineSpec{
			Description: "ODS",
			Params: []tekton.ParamSpec{
				{
					Name: "repository",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.Repository,
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "project",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.Project,
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "component",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.Component,
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "git-repo-url",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.GitURI,
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "git-full-ref",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.GitFullRef,
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "pr-key",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: strconv.Itoa(data.PullRequestKey),
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "pr-base",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.PullRequestBase,
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "environment",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.Environment,
						Type:      tekton.ParamTypeString,
					},
				},
				{
					Name: "version",
					Type: "string",
					Default: &tekton.ArrayOrString{
						StringVal: data.Version,
						Type:      tekton.ParamTypeString,
					},
				},
			},
			Tasks: tasks,
			Workspaces: []tekton.PipelineWorkspaceDeclaration{
				{
					Name: sharedWorkspaceName,
				},
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
