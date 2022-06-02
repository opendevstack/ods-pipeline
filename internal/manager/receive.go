package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	intrepo "github.com/opendevstack/pipeline/internal/repository"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
)

const (
	// allowedChangeRefType is the Bitbucket change ref handled by this service.
	allowedChangeRefType = "BRANCH"
)

// BitbucketWebhookReceiver receives webhook requests from Bitbucket.
type BitbucketWebhookReceiver struct {
	// Channel to send new runs to
	TriggeredPipelines chan PipelineConfig
	// Logger is the logger to send logging messages to.
	Logger logging.LeveledLoggerInterface
	// BitbucketClient is a client to interact with Bitbucket.
	BitbucketClient bitbucketInterface
	// WebhookSecret is the shared Bitbucket secret to validate webhook requests.
	WebhookSecret string
	// Namespace is the Kubernetes namespace in which the server runs.
	Namespace string
	// Project is the Bitbucket project to which this server corresponds.
	Project string
	// RepoBase is the common URL base of all repositories on Bitbucket.
	RepoBase string
}

// PipelineInfo holds information about a triggered pipeline.
type PipelineInfo struct {
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
	TriggerEvent    string `json:"trigger-event"`
	Comment         string `json:"comment"`
	PullRequestKey  int    `json:"prKey"`
	PullRequestBase string `json:"prBase"`
}

// Handle handles Bitbucket requests. It extracts pipeline data from the request
// body and sends the gained data to the scheduler.
func (s *BitbucketWebhookReceiver) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "could not read body"
		s.Logger.Errorf("%s: %s", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if err := validatePayload(r.Header, body, []byte(s.WebhookSecret)); err != nil {
		msg := "failed to validate incoming request"
		s.Logger.Errorf("%s: %s", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	req := &requestBitbucket{}
	if err := json.Unmarshal(body, &req); err != nil {
		msg := fmt.Sprintf("cannot parse JSON: %s", err)
		s.Logger.Errorf(msg)
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
			s.Logger.Warnf(msg)
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
		s.Logger.Warnf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	project = determineProject(s.Project, projectParam)
	component = strings.TrimPrefix(repo, project+"-")
	pInfo := PipelineInfo{
		Name:       makePipelineName(component, gitRef),
		Project:    project,
		Component:  component,
		Repository: repo,
		GitRef:     gitRef,
		GitFullRef: gitFullRef,
		RepoBase:   s.RepoBase,
		// Assemble GitURI from scratch instead of using user-supplied URI to
		// protect against attacks from external Bitbucket servers and/or projects.
		GitURI:       fmt.Sprintf("%s/%s/%s.git", s.RepoBase, project, repo),
		Namespace:    s.Namespace,
		TriggerEvent: req.EventKey,
		Comment:      commentText,
	}

	if len(commitSHA) == 0 {
		csha, err := getCommitSHA(s.BitbucketClient, pInfo.Project, pInfo.Repository, pInfo.GitFullRef)
		if err != nil {
			msg := "could not get commit SHA"
			s.Logger.Errorf("%s: %s", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		commitSHA = csha
	}
	pInfo.GitSHA = commitSHA

	skip := shouldSkip(s.BitbucketClient, pInfo.Project, pInfo.Repository, commitSHA)
	if skip {
		msg := fmt.Sprintf("Commit %s should be skipped", commitSHA)
		s.Logger.Infof(msg)
		// According to MDN (https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/418),
		// "some websites use this response for requests they do not wish to handle [..]".
		http.Error(w, msg, http.StatusTeapot)
		return
	}
	prInfo, err := extractPullRequestInfo(s.BitbucketClient, pInfo.Project, pInfo.Repository, commitSHA)
	if err != nil {
		msg := "Could not extract PR info"
		s.Logger.Errorf("%s: %s", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	pInfo.PullRequestKey = prInfo.ID
	pInfo.PullRequestBase = prInfo.Base

	odsConfig, err := intrepo.GetODSConfig(
		s.BitbucketClient,
		pInfo.Project,
		pInfo.Repository,
		pInfo.GitFullRef,
	)

	if err != nil {
		msg := fmt.Sprintf("could not download ODS config for repo %s", pInfo.Repository)
		s.Logger.Errorf("%s: %s", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pInfo.Environment = selectEnvironmentFromMapping(odsConfig.BranchToEnvironmentMapping, pInfo.GitRef)
	pInfo.Stage = string(config.DevStage)
	if pInfo.Environment != "" {
		env, err := odsConfig.Environment(pInfo.Environment)
		if err != nil {
			msg := fmt.Sprintf("environment misconfiguration: %s", err)
			s.Logger.Errorf("%s: %s", msg, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		pInfo.Stage = string(env.Stage)
	}
	pInfo.Version = odsConfig.Version

	s.Logger.Infof("%+v", pInfo)

	cfg := PipelineConfig{
		PipelineInfo: pInfo,
		PVC:          makePVCName(component),
		Tasks:        odsConfig.Pipeline.Tasks,
		Finally:      odsConfig.Pipeline.Finally,
	}
	s.TriggeredPipelines <- cfg

	err = json.NewEncoder(w).Encode(pInfo)
	if err != nil {
		s.Logger.Errorf("cannot write body: %s", err)
		return
	}
}

// determineProject returns the project from given serverProject/projectParam.
func determineProject(serverProject, projectParam string) string {
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
