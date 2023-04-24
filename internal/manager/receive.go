package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/opendevstack/pipeline/internal/httpjson"
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
	Project         string `json:"project"`
	Component       string `json:"component"`
	Repository      string `json:"repository"`
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
func (s *BitbucketWebhookReceiver) Handle(w http.ResponseWriter, r *http.Request) (any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, httpjson.NewInternalProblem("could not read body", err)
	}

	if err := validatePayload(r.Header, body, []byte(s.WebhookSecret)); err != nil {
		return nil, httpjson.NewStatusProblem(
			http.StatusUnauthorized, "failed to validate incoming request", err,
		)
	}

	req := &requestBitbucket{}
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, httpjson.NewStatusProblem(
			http.StatusBadRequest, fmt.Sprintf("cannot parse JSON: %s", err), nil,
		)
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
				"skipping change ref type %s, only %s is supported",
				change.Ref.Type, allowedChangeRefType,
			)
			return nil, httpjson.NewStatusProblem(http.StatusUnprocessableEntity, msg, nil)
		}
	} else if strings.HasPrefix(req.EventKey, "pr:") {
		repo = strings.ToLower(req.PullRequest.FromRef.Repository.Slug)
		gitRef = strings.ToLower(req.PullRequest.FromRef.DisplayID)
		gitFullRef = req.PullRequest.FromRef.ID
		projectParam = req.PullRequest.FromRef.Repository.Project.Key
		if req.Comment != nil {
			commentText = req.Comment.Text
		}
		commitSHA = req.PullRequest.FromRef.LatestCommit
	} else {
		return nil, httpjson.NewStatusProblem(
			http.StatusBadRequest, fmt.Sprintf("unsupported event key: %s", req.EventKey), nil,
		)
	}

	project = determineProject(s.Project, projectParam)
	component = strings.TrimPrefix(repo, project+"-")
	pInfo := PipelineInfo{
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
			return nil, httpjson.NewInternalProblem("could not get commit SHA", err)
		}
		commitSHA = csha
	}
	pInfo.GitSHA = commitSHA

	skip := shouldSkip(s.BitbucketClient, pInfo.Project, pInfo.Repository, commitSHA)
	if skip {
		return nil, httpjson.NewStatusProblem(
			http.StatusAccepted, fmt.Sprintf("Commit %s should be skipped", commitSHA), nil,
		)
	}
	prInfo, err := extractPullRequestInfo(s.BitbucketClient, pInfo.Project, pInfo.Repository, commitSHA)
	if err != nil {
		return nil, httpjson.NewInternalProblem("could not extract PR info", err)
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
		return nil, httpjson.NewInternalProblem(
			fmt.Sprintf("could not download ODS config for repo %s", pInfo.Repository), err,
		)
	}
	pInfo.Version = odsConfig.Version

	s.Logger.Infof("%+v", pInfo)

	cfg, err := identifyPipelineConfig(pInfo, odsConfig, component)
	if err != nil {
		return nil, httpjson.NewStatusProblem(
			http.StatusBadRequest, "could not identify pipeline to run", err,
		)
	}
	s.TriggeredPipelines <- *cfg

	return pInfo, nil
}

// identifyPipelineConfig finds the first configuration matching the triggering event
func identifyPipelineConfig(pInfo PipelineInfo, odsConfig *config.ODS, component string) (*PipelineConfig, error) {
	for _, p := range odsConfig.Pipelines {
		if len(p.Triggers) == 0 {
			return &PipelineConfig{
				PipelineInfo: pInfo,
				PVC:          makePVCName(component),
				PipelineSpec: p,
				// no params available
			}, nil
		}
		for _, t := range p.Triggers {
			if triggerMatches(pInfo, t) {
				return &PipelineConfig{
					PipelineInfo: pInfo,
					PVC:          makePVCName(component),
					PipelineSpec: p,
					Params:       t.Params,
				}, nil
			}
		}
	}
	return nil, errors.New("no trigger definition matched webhook event")
}

func triggerMatches(pInfo PipelineInfo, trigger config.Trigger) bool {
	return triggerEventsMatch(pInfo, trigger) && triggerBranchesMatch(pInfo, trigger) &&
		triggerExcludedBranchesDoNotMatch(pInfo, trigger) && triggerPRCommentMatches(pInfo, trigger)
}

func anyPatternMatches(s string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}

	for _, pattern := range patterns {
		if matched, err := path.Match(pattern, s); matched && err == nil {
			return true
		}
	}
	return false
}

func triggerEventsMatch(pInfo PipelineInfo, trigger config.Trigger) bool {
	return anyPatternMatches(pInfo.TriggerEvent, trigger.Events)
}

func triggerBranchesMatch(pInfo PipelineInfo, trigger config.Trigger) bool {
	return anyPatternMatches(pInfo.GitRef, trigger.Branches)
}

func triggerExcludedBranchesDoNotMatch(pInfo PipelineInfo, trigger config.Trigger) bool {
	if len(trigger.ExceptBranches) == 0 {
		return true
	}
	return !anyPatternMatches(pInfo.GitRef, trigger.ExceptBranches)
}

func triggerPRCommentMatches(pInfo PipelineInfo, trigger config.Trigger) bool {
	prefix := trigger.PrComment
	if prefix == nil || *prefix == "" {
		return true
	}

	return strings.HasPrefix(pInfo.Comment, *prefix)
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
