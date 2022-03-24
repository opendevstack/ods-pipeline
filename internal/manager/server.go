package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	kubernetesClient "github.com/opendevstack/pipeline/internal/kubernetes"
	intrepo "github.com/opendevstack/pipeline/internal/repository"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// allowedChangeRefType is the Bitbucket change ref handled by this service.
	allowedChangeRefType = "BRANCH"
	// letterBytes contains letters to use for random strings.
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// Server represents this service, and is a global.
type Server struct {
	KubernetesClient  kubernetesClient.ClientInterface
	TektonClient      tektonClient.ClientInterface
	Namespace         string
	Project           string
	RepoBase          string
	Token             string
	WebhookSecret     string
	TaskKind          tekton.TaskKind
	TaskSuffix        string
	StorageConfig     StorageConfig
	BitbucketClient   bitbucketInterface
	PipelineRunPruner PipelineRunPruner
	Mutex             sync.Mutex
	Logger            logging.LeveledLoggerInterface
}

type StorageConfig struct {
	Provisioner string
	ClassName   string
	Size        string
}

// ServerConfig configures a server.
type ServerConfig struct {
	// Namespace is the Kubernetes namespace in which the server runs.
	Namespace string
	// Project is the Bitbucket project to which this server corresponds.
	Project string
	// RepoBase is the common URL base of all repositories on Bitbucket.
	RepoBase string
	// Token is the Bitbucket personal access token.
	Token string
	// WebhookSecret is the shared Bitbucket secret to validate webhook requests.
	WebhookSecret string
	// TaskKind is the Tekton resource kind for tassks.
	// Either "ClusterTask" or "Task".
	TaskKind string
	// TaskSuffic is the suffix applied to tasks (version information).
	TaskSuffix string
	// StorageConfig describes the config to apply to PVCs.
	StorageConfig StorageConfig
	// KubernetesClient is a Kubernetes client
	KubernetesClient kubernetesClient.ClientInterface
	// TektonClient is a Tekton client
	TektonClient tektonClient.ClientInterface
	// BitbucketClient is a Bitbucket client
	BitbucketClient bitbucketInterface
	// PipelineRunPruner is responsible to prune pipeline runs.
	PipelineRunPruner PipelineRunPruner
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
}

type PipelineData struct {
	Name            string `json:"name"`
	Project         string `json:"project"`
	Component       string `json:"component"`
	Repository      string `json:"repository"`
	Stage           string `json:"stage"`
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
func NewServer(serverConfig ServerConfig) (*Server, error) {
	if serverConfig.Namespace == "" {
		return nil, errors.New("namespace is required")
	}
	if serverConfig.Token == "" {
		return nil, errors.New("token is required")
	}
	if serverConfig.WebhookSecret == "" {
		return nil, errors.New("webhook secret is required")
	}
	if serverConfig.StorageConfig.ClassName == "" {
		return nil, errors.New("storage class name is required")
	}
	if serverConfig.StorageConfig.Size == "" {
		return nil, errors.New("storage size is required")
	}
	if serverConfig.TektonClient == nil {
		return nil, errors.New("tekton client is required")
	}
	if serverConfig.BitbucketClient == nil {
		return nil, errors.New("bitbucket client is required")
	}
	if serverConfig.Logger == nil {
		serverConfig.Logger = &logging.LeveledLogger{Level: logging.LevelError}
	}
	s := &Server{
		KubernetesClient:  serverConfig.KubernetesClient,
		TektonClient:      serverConfig.TektonClient,
		BitbucketClient:   serverConfig.BitbucketClient,
		Namespace:         serverConfig.Namespace,
		Project:           serverConfig.Project,
		RepoBase:          serverConfig.RepoBase,
		Token:             serverConfig.Token,
		WebhookSecret:     serverConfig.WebhookSecret,
		TaskKind:          tekton.TaskKind(serverConfig.TaskKind),
		TaskSuffix:        serverConfig.TaskSuffix,
		StorageConfig:     serverConfig.StorageConfig,
		PipelineRunPruner: serverConfig.PipelineRunPruner,
		Logger:            serverConfig.Logger,
	}
	return s, nil
}

// HandleRoot handles all requests to this service. It performs the following:
// - extract pipeline data from request body
// - extend body with calculated pipeline information
// - create/update pipeline that will be triggerd by event listener
func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {

	requestID := randStringBytes(6)
	s.Logger.Infof("%s ---START---", requestID)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "could not read body"
		s.Logger.Errorf("%s %s: %s", requestID, msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if err := validatePayload(r.Header, body, []byte(s.WebhookSecret)); err != nil {
		msg := "failed to validate incoming request"
		s.Logger.Errorf("%s %s: %s", requestID, msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	req := &requestBitbucket{}
	if err := json.Unmarshal(body, &req); err != nil {
		msg := fmt.Sprintf("cannot parse JSON: %s", err)
		s.Logger.Errorf("%s %s", requestID, msg)
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

		if change.Ref.Type != allowedChangeRefType {
			msg := fmt.Sprintf(
				"Skipping change ref type %s, only %s is supported",
				change.Ref.Type,
				allowedChangeRefType,
			)
			s.Logger.Warnf("%s %s", requestID, msg)
			// According to MDN (https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/418),
			// "some websites use this response for requests they do not wish to handle [...]".
			http.Error(w, msg, http.StatusTeapot)
			return
		}
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
		msg := fmt.Sprintf("Unsupported event key: %s", req.EventKey)
		s.Logger.Warnf("%s %s", requestID, msg)
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

	pData := PipelineData{
		Name:         pipelineName,
		Project:      project,
		Component:    component,
		Repository:   repo,
		GitRef:       gitRef,
		GitFullRef:   gitFullRef,
		RepoBase:     s.RepoBase,
		GitURI:       gitURI,
		Namespace:    s.Namespace,
		PVC:          makePVCName(component),
		TriggerEvent: req.EventKey,
		Comment:      commentText,
	}

	if len(commitSHA) == 0 {
		csha, err := getCommitSHA(s.BitbucketClient, pData.Project, pData.Repository, pData.GitFullRef)
		if err != nil {
			msg := "could not get commit SHA"
			s.Logger.Errorf("%s %s: %s", requestID, msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		commitSHA = csha
	}
	pData.GitSHA = commitSHA

	skip := shouldSkip(s.BitbucketClient, pData.Project, pData.Repository, commitSHA)
	if skip {
		msg := "Commit should be skipped"
		s.Logger.Infof("%s %s: %s", requestID, msg, err)
		// According to MDN (https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/418),
		// "some websites use this response for requests they do not wish to handle [..]".
		http.Error(w, msg, http.StatusTeapot)
		return
	}
	prInfo, err := extractPullRequestInfo(s.BitbucketClient, pData.Project, pData.Repository, commitSHA)
	if err != nil {
		msg := "Could not extract PR info"
		s.Logger.Errorf("%s %s: %s", requestID, msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	pData.PullRequestKey = prInfo.ID
	pData.PullRequestBase = prInfo.Base

	odsConfig, err := intrepo.GetODSConfig(
		s.BitbucketClient,
		pData.Project,
		pData.Repository,
		pData.GitFullRef,
	)

	if err != nil {
		msg := fmt.Sprintf("could not download ODS config for repo %s", pData.Repository)
		s.Logger.Errorf("%s %s: %s", requestID, msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pData.Environment = selectEnvironmentFromMapping(odsConfig.BranchToEnvironmentMapping, pData.GitRef)
	pData.Stage = string(config.DevStage)
	if pData.Environment != "" {
		env, err := odsConfig.Environment(pData.Environment)
		if err != nil {
			msg := fmt.Sprintf("environment misconfiguration: %s", err)
			s.Logger.Errorf("%s %s: %s", requestID, msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		pData.Stage = string(env.Stage)
	}
	pData.Version = odsConfig.Version

	newPipeline := assemblePipeline(odsConfig, pData, s.TaskKind, s.TaskSuffix)

	existingPipeline, err := s.TektonClient.GetPipeline(r.Context(), pData.Name, metav1.GetOptions{})
	if err != nil {
		_, err := s.TektonClient.CreatePipeline(r.Context(), newPipeline, metav1.CreateOptions{})
		if err != nil {
			msg := fmt.Sprintf("cannot create pipeline %s", newPipeline.Name)
			s.Logger.Errorf("%s %s: %s", requestID, msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	} else {
		newPipeline.ResourceVersion = existingPipeline.ResourceVersion
		_, err := s.TektonClient.UpdatePipeline(r.Context(), newPipeline, metav1.UpdateOptions{})
		if err != nil {
			msg := fmt.Sprintf("cannot update pipeline %s", newPipeline.Name)
			s.Logger.Errorf("%s %s: %s", requestID, msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}

	// Create PVC if it does not exist yet
	err = s.createPVCIfRequired(r.Context(), pData)
	if err != nil {
		msg := "cannot create workspace PVC"
		s.Logger.Errorf("%s %s: %s", requestID, msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Aquire lock to avoid calls timewise close to one another to lead to
	// parallel pipeline runs.
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Logger.Infof("%s Starting pruning of pipeline runs related to repository %s ...", requestID, pData.Repository)
	ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	pipelineRuns, err := listPipelineRuns(s.TektonClient, ctxt, pData.Repository)
	if err != nil {
		msg := "could not retrieve existing pipeline runs"
		s.Logger.Errorf("%s %s: %s", requestID, msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	s.Logger.Debugf("%s Found %d pipeline runs related to repository %s.", requestID, len(pipelineRuns.Items), pData.Repository)

	if s.PipelineRunPruner != nil {
		go func() {
			ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			err = s.PipelineRunPruner.Prune(ctxt, pipelineRuns.Items)
			if err != nil {
				s.Logger.Warnf(
					"Pruning pipeline runs of repository %s failed: %s",
					pData.Repository, err,
				)
			}
		}()
	}

	s.Logger.Infof("%s %+v", requestID, pData)

	_, err = createPipelineRun(s.TektonClient, r.Context(), pData)
	if err != nil {
		msg := "cannot create pipeline run"
		s.Logger.Errorf("%s %s: %s", requestID, msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(pData)
	if err != nil {
		s.Logger.Errorf("%s cannot write body: %s", requestID, err)
		return
	}
}

// selectEnvironmentFromMapping selects the environment name matching given branch.
func selectEnvironmentFromMapping(mapping []config.BranchToEnvironmentMapping, branch string) string {
	for _, bem := range mapping {
		if mappingBranchMatch(bem.Branch, branch) {
			return bem.Environment
		}
	}
	return ""
}

func mappingBranchMatch(mappingBranch, testBranch string) bool {
	// exact match
	if mappingBranch == testBranch {
		return true
	}
	// prefix match like "release/*", also catches "*"
	if strings.HasSuffix(mappingBranch, "*") {
		branchPrefix := strings.TrimSuffix(mappingBranch, "*")
		if strings.HasPrefix(testBranch, branchPrefix) {
			return true
		}
	}
	return false
}

func determineProject(serverProject string, projectParam string) string {
	projectParam = strings.ToLower(projectParam)
	if len(projectParam) > 0 {
		return projectParam
	}
	return strings.ToLower(serverProject)
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

// randStringBytes creates a random string of length n.
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
