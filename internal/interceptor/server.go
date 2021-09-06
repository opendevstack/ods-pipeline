package interceptor

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	taskKind = "ClusterTask"
)

// Server represents this service, and is a global.
type Server struct {
	OpenShiftClient Client
	Namespace       string
	Project         string
	RepoBase        string
	Token           string
	TaskSuffix      string
	BitbucketClient *bitbucket.Client
}

type PipelineData struct {
	Name            string `json:"name"`
	ResourceVersion int    `json:"resourceVersion"`
	Project         string `json:"project"`
	Component       string `json:"component"`
	Repository      string `json:"repository"`
	Environment     string `json:"environment"`
	Version         string `json:"version"`
	GitRef          string `json:"gitRef"`
	GitFullRef      string `json:"gitFullRef"`
	GitSHA          string `json:"gitSha"`
	RepoBase        string `json:"repoBase"`
	GitURI          string `json:"gitURI"`
	Namespace       string `json:"namespace"`
	PVC             string `json:"pvc"`
	TriggerEvent    string `json:"trigger-event"`
	Comment         string `json:"comment"`
	PullRequestKey  int    `json:"prKey"`
	PullRequestBase string `json:"prBase"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewServer returns a new server.
func NewServer(client Client, namespace, project, repoBase, token, taskSuffix string) *Server {
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: token,
		BaseURL:  strings.TrimSuffix(repoBase, "/scm"),
	})
	return &Server{
		OpenShiftClient: client,
		Namespace:       namespace,
		Project:         project,
		RepoBase:        repoBase,
		Token:           token,
		TaskSuffix:      taskSuffix,
		BitbucketClient: bitbucketClient,
	}
}

type repository struct {
	Project struct {
		Key string `json:"key"`
	} `json:"project"`
	Slug string `json:"slug"`
}
type requestBitbucket struct {
	EventKey   string     `json:"eventKey"`
	Repository repository `json:"repository"`
	Changes    []struct {
		Type string `json:"type"`
		Ref  struct {
			ID        string `json:"id"`
			DisplayID string `json:"displayId"`
			Type      string `json:"type"`
		} `json:"ref"`
		FromHash string `json:"fromHash"`
		ToHash   string `json:"toHash"`
	} `json:"changes"`
	PullRequest *struct {
		FromRef struct {
			Repository   repository `json:"repository"`
			ID           string     `json:"id"`
			DisplayID    string     `json:"displayId"`
			LatestCommit string     `json:"latestCommit"`
		} `json:"fromRef"`
	} `json:"pullRequest"`
	Comment *struct {
		Text string `json:"text"`
	} `json:"comment"`
}

// HandleRoot handles all requests to this service.
func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {

	// read request body into Go object
	// extract pipeline data
	// extend body with data
	// create/update pipeline

	requestID := randStringBytes(6)
	log.Println(requestID, "---START---")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(requestID, err.Error())
		http.Error(w, "could not read body", http.StatusInternalServerError)
		return
	}

	req := &requestBitbucket{}
	if err := json.Unmarshal(body, &req); err != nil {
		msg := fmt.Sprintf("cannot parse JSON: %s", err)
		log.Println(requestID, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var repo string
	var gitRef string
	var gitFullRef string
	var project string
	var projectParam string
	var component string
	var commitSHA string
	commentText := ""

	if req.EventKey == "repo:refs_changed" {
		repo = strings.ToLower(req.Repository.Slug)
		change := req.Changes[0]
		gitRef = strings.ToLower(change.Ref.DisplayID)
		gitFullRef = change.Ref.ID

		projectParam = req.Repository.Project.Key
		commitSHA = change.ToHash
	} else if strings.HasPrefix(req.EventKey, "pr:") {
		repo = strings.ToLower(req.PullRequest.FromRef.Repository.Slug)
		gitRef = strings.ToLower(req.PullRequest.FromRef.DisplayID)
		gitFullRef = req.PullRequest.FromRef.ID

		projectParam = req.PullRequest.FromRef.Repository.Project.Key
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Comment != nil {
			commentText = req.Comment.Text
		}
		commitSHA = req.PullRequest.FromRef.LatestCommit
	} else {
		msg := fmt.Sprintf("Unsupported event key: %s", err)
		log.Println(requestID, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	project = determineProject(s.Project, projectParam)
	component = strings.TrimPrefix(repo, project+"-")

	gitURI := fmt.Sprintf(
		"%s/%s/%s.git",
		s.RepoBase,
		project,
		repo,
	)

	pipelineName := makePipelineName(component, gitRef)
	resourceVersion, err := s.OpenShiftClient.GetPipelineResourceVersion(pipelineName)
	if err != nil {
		msg := "Could not retrieve pipeline resourceVersion"
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, 500)
		return
	}

	pData := PipelineData{
		Name:            pipelineName,
		Project:         project,
		Component:       component,
		Repository:      repo,
		GitRef:          gitRef,
		GitFullRef:      gitFullRef,
		ResourceVersion: resourceVersion,
		RepoBase:        s.RepoBase,
		GitURI:          gitURI,
		Namespace:       s.Namespace,
		PVC:             "ods-pipeline",
		TriggerEvent:    req.EventKey,
		Comment:         commentText,
	}

	if len(commitSHA) == 0 {
		csha, err := getCommitSHA(s.BitbucketClient, pData.Project, pData.Repository, pData.GitFullRef)
		if err != nil {
			msg := "could not get commit SHA"
			log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		commitSHA = csha
	}
	pData.GitSHA = commitSHA

	skip := shouldSkip(s.BitbucketClient, pData.Project, pData.Repository, commitSHA)
	if skip {
		msg := "Commit should be skipped"
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	prInfo, err := extractPullRequestInfo(s.BitbucketClient, pData.Project, pData.Repository, commitSHA)
	if err != nil {
		msg := "Could not extract PR info"
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	pData.PullRequestKey = prInfo.ID
	pData.PullRequestBase = prInfo.Base

	log.Println(requestID, fmt.Sprintf("%+v", pData))

	extendedBody, err := extendBodyWithExtensions(body, pData)
	if err != nil {
		msg := "cannot extend body"
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(extendedBody)
	if err != nil {
		msg := "cannot write body"
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	odsConfig, err := getODSConfig(
		s.BitbucketClient,
		pData.Project,
		pData.Repository,
		pData.GitFullRef,
	)

	if err != nil {
		msg := fmt.Sprintf("could not download ODS config for repo %s", pData.Repository)
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pData.Environment = selectEnvironmentFromMapping(odsConfig.BranchToEnvironmentMapping, pData.GitRef)
	pData.Version = odsConfig.Version

	rendered, err := renderPipeline(odsConfig, pData, s.TaskSuffix)
	if err != nil {
		msg := "Could not render pipeline definition"
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, 500)
		return
	}

	jsonBytes, err := yaml.YAMLToJSON(rendered)
	if err != nil {
		msg := "could not convert YAML to JSON"
		log.Println(requestID, fmt.Sprintf("%s: %s", msg, err))
		http.Error(w, msg, 500)
		return
	}

	createStatusCode, createErr := s.OpenShiftClient.ApplyPipeline(jsonBytes, pData)
	if createErr != nil {
		msg := "Could not create/update pipeline"
		log.Println(requestID, fmt.Sprintf("%s [%d]: %s", msg, createStatusCode, createErr))
		http.Error(w, msg, createStatusCode)
		return
	}
}

func selectEnvironmentFromMapping(mapping []config.BranchToEnvironmentMapping, branch string) string {
	for _, bem := range mapping {
		// exact match
		if bem.Branch == branch {
			return bem.Environment
		}
		// prefix match like "release/*", also catches "*"
		if strings.HasSuffix(bem.Branch, "*") {
			branchPrefix := strings.TrimSuffix(bem.Branch, "*")
			if strings.HasPrefix(branch, branchPrefix) {
				return bem.Environment
			}
		}
	}
	return ""
}

func getCommitSHA(bitbucketClient *bitbucket.Client, project, repository, gitFullRef string) (string, error) {
	commitList, err := bitbucketClient.CommitList(project, repository, bitbucket.CommitListParams{
		Until: gitFullRef,
	})
	if err != nil {
		return "", fmt.Errorf("could not get commit list: %w", err)
	}
	return commitList.Values[0].ID, nil
}

func getODSConfig(bitbucketClient *bitbucket.Client, project, repository, gitFullRef string) (*config.ODS, error) {
	var body []byte
	var getErr error
	for _, c := range config.ODSFileCandidates {
		b, err := bitbucketClient.RawGet(project, repository, c, gitFullRef)
		if err == nil {
			body = b
			getErr = nil
			break
		}
		getErr = err
	}
	if getErr != nil {
		return nil, fmt.Errorf("could not download ODS config for repo %s: %w", repository, getErr)
	}

	if body == nil {
		log.Printf("no ODS config located in repo %s", repository)
		return nil, nil
	}
	return config.Read(body)
}

func determineProject(serverProject string, projectParam string) string {
	projectParam = strings.ToLower(projectParam)
	if len(projectParam) > 0 {
		return projectParam
	}
	return strings.ToLower(serverProject)
}

func extendBodyWithExtensions(body []byte, data PipelineData) ([]byte, error) {
	var payload map[string]interface{}
	err := json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}
	payload["extensions"] = data
	extendedBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return extendedBody, nil
}

// makePipelineName generates the name of the pipeline.
// According to the Kubernetes label rules, a maximum of 63 characters is
// allowed, see https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set.
// Therefore, the name might be truncated. As this could cause potential clashes
// between similar named branches, we put a short part of the branch hash value
// into the name name to make this very unlikely.
// We cut the pipeline name at 55 chars to allow e.g. pipeline runs to add suffixes.
func makePipelineName(component string, branch string) string {
	// Cut all non-alphanumeric characters
	safeCharsRegex := regexp.MustCompile("[^-a-zA-Z0-9]+")
	pipeline := component + "-" + safeCharsRegex.ReplaceAllString(
		strings.Replace(branch, "/", "-", -1),
		"",
	)

	// Enforce maximum length - and if truncation needs to happen,
	// ensure uniqueness of pipeline name as much as possible.
	if len(pipeline) > 55 {
		shortenedPipeline := pipeline[0:47]
		h := sha1.New()
		_, err := h.Write([]byte(pipeline))
		if err != nil {
			return shortenedPipeline
		}
		bs := h.Sum(nil)
		s := fmt.Sprintf("%x", bs)
		pipeline = fmt.Sprintf("%s-%s", shortenedPipeline, s[0:7])
	}
	return pipeline
}

func renderPipeline(odsConfig *config.ODS, data PipelineData, taskSuffix string) ([]byte, error) {

	var tasks []tekton.PipelineTask
	tasks = append(tasks, tekton.PipelineTask{
		Name:    "ods-start",
		TaskRef: &tekton.TaskRef{Kind: taskKind, Name: "ods-start" + taskSuffix},
		Workspaces: []tekton.WorkspacePipelineTaskBinding{
			{Name: "source", Workspace: "shared-workspace"},
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
			{Name: "source", Workspace: "shared-workspace"},
		},
		Params: []tekton.Param{
			{
				Name: "pipeline-run-name",
				Value: tekton.ArrayOrString{
					StringVal: "$(context.pipelineRun.name)",
					Type:      tekton.ParamTypeString,
				},
			},
		},
	})
	p := tekton.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:            data.Name,
			ResourceVersion: strconv.Itoa(data.ResourceVersion),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "tekton.dev/v1beta1",
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
					Name: "shared-workspace",
				},
			},
			Finally: finallyTasks,
		},
	}

	return yaml.Marshal(p)
}

type prInfo struct {
	ID   int
	Base string
}

func extractPullRequestInfo(bitbucketClient *bitbucket.Client, projectKey, repositorySlug, gitCommit string) (prInfo, error) {
	var i prInfo

	prPage, err := bitbucketClient.CommitPullRequestList(projectKey, repositorySlug, gitCommit)
	if err != nil {
		return i, err
	}

	for _, v := range prPage.Values {
		if !v.Open {
			continue
		}
		i.ID = v.ID
		i.Base = v.ToRef.ID
		break
	}

	return i, nil
}

func shouldSkip(bitbucketClient *bitbucket.Client, projectKey, repositorySlug, gitCommit string) bool {
	c, err := bitbucketClient.CommitGet(projectKey, repositorySlug, gitCommit)
	if err != nil {
		return false
	}
	return isCiSkipInCommitMessage(c.Message)
}

/** Looks in commit message for the following strings
 *  '[ci skip]', '[ciskip]', '[ci-skip]', '[ci_skip]',
 *  '[skip ci]', '[skipci]', '[skip-ci]', '[skip_ci]',
 *  '***NO_CI***', '***NO CI***', '***NOCI***', '***NO-CI***'
 */
func isCiSkipInCommitMessage(message string) bool {
	messageLines := strings.Split(message, "\n")
	re := regexp.MustCompile(`[\s\-\_]`)
	subject := strings.ToLower(messageLines[0])
	subject = re.ReplaceAllString(subject, "")

	return strings.Contains(subject, "[ciskip]") ||
		strings.Contains(subject, "[skipci]") ||
		strings.Contains(subject, "***noci***")
}
