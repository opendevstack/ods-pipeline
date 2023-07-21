package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/opendevstack/pipeline/internal/httpjson"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
)

const changeRefTypeTag = "TAG"

type BitbucketWebhookReceiverBase struct {
	// Namespace is the Kubernetes namespace in which the server runs.
	Namespace string
	// Project is the Bitbucket project to which this server corresponds.
	Project string
	// RepoBase is the common URL base of all repositories on Bitbucket.
	RepoBase string
}

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

	BitbucketWebhookReceiverBase
}

// PipelineInfo holds information about a triggered pipeline.
type PipelineInfo struct {
	Project         string `json:"project"`
	Component       string `json:"component"`
	Repository      string `json:"repository"`
	GitRef          string `json:"gitRef"`
	GitFullRef      string `json:"gitFullRef"`
	GitSHA          string `json:"gitSha"`
	RepoBase        string `json:"repoBase"`
	GitURI          string `json:"gitURI"`
	Namespace       string `json:"namespace"`
	TriggerEvent    string `json:"trigger-event"`
	ChangeRefType   string `json:"change-ref-type"`
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

	pInfo, err := readBitbucketRequest(&s.BitbucketWebhookReceiverBase, body)
	if err != nil {
		return nil, err
	}

	if len(pInfo.GitSHA) == 0 {
		csha, err := getCommitSHA(s.BitbucketClient, pInfo.Project, pInfo.Repository, pInfo.GitFullRef)
		if err != nil {
			return nil, httpjson.NewInternalProblem("could not get commit SHA", err)
		}
		pInfo.GitSHA = csha
	}

	skip := shouldSkip(s.BitbucketClient, pInfo.Project, pInfo.Repository, pInfo.GitSHA)
	if skip {
		return nil, httpjson.NewStatusProblem(
			http.StatusAccepted, fmt.Sprintf("Commit %s should be skipped", pInfo.GitSHA), nil,
		)
	}
	prInfo, err := extractPullRequestInfo(s.BitbucketClient, pInfo.Project, pInfo.Repository, pInfo.GitSHA)
	if err != nil {
		return nil, httpjson.NewInternalProblem("could not extract PR info", err)
	}
	pInfo.PullRequestKey = prInfo.ID
	pInfo.PullRequestBase = prInfo.Base

	odsConfig, err := fetchODSConfig(
		s.BitbucketClient,
		pInfo.Project,
		pInfo.Repository,
		pInfo.GitFullRef,
	)
	if err != nil {
		return nil, httpjson.NewInternalProblem(
			fmt.Sprintf("could not fetch ODS config for repo %s", pInfo.Repository), err,
		)
	}

	s.Logger.Infof("%+v", pInfo)

	cfg := identifyPipelineConfig(*pInfo, *odsConfig, pInfo.Component)
	if cfg == nil {
		return nil, httpjson.NewStatusProblem(
			http.StatusAccepted, "Could not identify any pipeline to run as no trigger matched", nil,
		)
	}

	s.TriggeredPipelines <- *cfg
	return pInfo, nil
}

func readBitbucketRequest(bbb *BitbucketWebhookReceiverBase, body []byte) (*PipelineInfo, error) {
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
	var changeRefType string
	commentText := ""

	if req.EventKey == "repo:refs_changed" {
		repo = strings.ToLower(req.Repository.Slug)
		change := req.Changes[0]
		gitRef = change.Ref.DisplayID
		gitFullRef = change.Ref.ID

		projectParam = req.Repository.Project.Key
		commitSHA = change.ToHash
		changeRefType = change.Ref.Type
	} else if strings.HasPrefix(req.EventKey, "pr:") {
		repo = strings.ToLower(req.PullRequest.FromRef.Repository.Slug)
		gitRef = req.PullRequest.FromRef.DisplayID
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
	project = determineProject(bbb.Project, projectParam)
	component = strings.TrimPrefix(repo, project+"-")

	pInfo := PipelineInfo{
		Project:    project,
		Component:  component,
		Repository: repo,
		GitRef:     gitRef,
		GitFullRef: gitFullRef,
		RepoBase:   bbb.RepoBase,
		// Assemble GitURI from scratch instead of using user-supplied URI to
		// protect against attacks from external Bitbucket servers and/or projects.
		GitURI:        fmt.Sprintf("%s/%s/%s.git", bbb.RepoBase, project, repo),
		GitSHA:        commitSHA,
		Namespace:     bbb.Namespace,
		TriggerEvent:  req.EventKey,
		ChangeRefType: changeRefType,
		Comment:       commentText,
	}
	return &pInfo, nil
}

// identifyPipelineConfig finds the first configuration matching the triggering event
func identifyPipelineConfig(pInfo PipelineInfo, odsConfig config.ODS, component string) *PipelineConfig {
	for _, p := range odsConfig.Pipelines {
		if len(p.Triggers) == 0 {
			return &PipelineConfig{
				PipelineInfo: pInfo,
				PVC:          makePVCName(component),
				PipelineSpec: p,
				// no params available
			}
		}
		for _, t := range p.Triggers {
			if triggerMatches(pInfo, t) {
				return &PipelineConfig{
					PipelineInfo: pInfo,
					PVC:          makePVCName(component),
					PipelineSpec: p,
					Params:       t.Params,
				}
			}
		}
	}
	return nil
}

func triggerMatches(pInfo PipelineInfo, trigger config.Trigger) bool {
	return triggerEventsMatch(pInfo, trigger) &&
		((triggerBranchesMatch(pInfo, trigger) && triggerExcludedBranchesDoNotMatch(pInfo, trigger)) ||
			(triggerTagsMatch(pInfo, trigger) && triggerExcludedTagsDoNotMatch(pInfo, trigger))) &&
		triggerPRCommentMatches(pInfo, trigger)
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

// triggerBranchesMatch is true when any pattern matches the Git ref (which refers to a branch).
// If the change type is "TAG" or tag constraints are configured, it will always return "false".
func triggerBranchesMatch(pInfo PipelineInfo, trigger config.Trigger) bool {
	if pInfo.ChangeRefType == changeRefTypeTag || len(trigger.Tags) > 0 {
		return false
	}
	return anyPatternMatches(pInfo.GitRef, trigger.Branches)
}

// triggerExcludedBranchesDoNotMatch is true when no pattern matches the Git ref.
// If the change type is "TAG", it will always return "false".
func triggerExcludedBranchesDoNotMatch(pInfo PipelineInfo, trigger config.Trigger) bool {
	if pInfo.ChangeRefType == changeRefTypeTag {
		return false
	}
	if len(trigger.ExceptBranches) == 0 {
		return true
	}
	return !anyPatternMatches(pInfo.GitRef, trigger.ExceptBranches)
}

// triggerBranchesMatch is true when any pattern matches the Git ref (which refers to a tag).
// If the change type is not "TAG" or branch constraints are configured, it will always return "false".
func triggerTagsMatch(pInfo PipelineInfo, trigger config.Trigger) bool {
	if pInfo.ChangeRefType != changeRefTypeTag || len(trigger.Branches) > 0 {
		return false
	}
	return anyPatternMatches(pInfo.GitRef, trigger.Tags)
}

// triggerExcludedTagsDoNotMatch is true when no pattern matches the Git ref.
// If the change type is not "TAG", it will always return "false".
func triggerExcludedTagsDoNotMatch(pInfo PipelineInfo, trigger config.Trigger) bool {
	if pInfo.ChangeRefType != changeRefTypeTag {
		return false
	}
	if len(trigger.ExceptTags) == 0 {
		return true
	}
	return !anyPatternMatches(pInfo.GitRef, trigger.ExceptTags)
}

// triggerPRCommentMatches is true when the comment is prefixed with trigger.PrComment.
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

// fetchODSConfig returns a *config.ODS for given project/repository at gitFullRef.
// If retrieving fails or not ods.y(a)ml file exists, it errors.
func fetchODSConfig(bitbucketClient bitbucket.RawClientInterface, project, repository, gitFullRef string) (*config.ODS, error) {
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
		return nil, fmt.Errorf("could not find ODS config for repo %s: %w", repository, getErr)
	}

	if body == nil {
		return nil, fmt.Errorf("no ODS config located in repo %s", repository)
	}
	return config.Read(body)
}
