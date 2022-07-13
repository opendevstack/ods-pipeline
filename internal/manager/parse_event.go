package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// allowedChangeRefType is the Bitbucket change ref handled by this service.
	allowedChangeRefType = "BRANCH"
)

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

func (s *BitbucketWebhookReceiver) HandleParseBitbucketWebhookEvent(w http.ResponseWriter, r *http.Request) {
	parsedGitInfo := s.handleParseBitbucketWebhookEvent(w, r)
	if parsedGitInfo == nil {
		return
	}
	s.Handle(w, r, parsedGitInfo)
}

// Parses Bitbucket webhook requests and handle them if invalid by writing a response and logging issues.
func (s *BitbucketWebhookReceiver) handleParseBitbucketWebhookEvent(w http.ResponseWriter, r *http.Request) *ParseGitEvent {
	body := s.handleValidateBodyHasProperSignature(w, r)
	if body == nil {
		return nil
	}

	req := &requestBitbucket{}
	if err := json.Unmarshal(body, &req); err != nil {
		msg := fmt.Sprintf("cannot parse JSON: %s", err)
		s.Logger.Errorf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return nil
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
			return nil
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
		msg := fmt.Sprintf("Unsupported event key: %s", req.EventKey)
		s.Logger.Warnf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return nil
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
	return &ParseGitEvent{
		pInfo, commitSHA,
	}
}

type RequestPersonalBuild struct {
	RepositorySlug string `json:"repo"`
	FullRefName    string `json:"fullRefName"`
	BranchName     string `json:"branchName"`
	ToHash         string `json:"toHash"`
}

func (s *BitbucketWebhookReceiver) HandleParsePersonalBuildEvent(w http.ResponseWriter, r *http.Request) {
	parsedGitInfo := s.handleParsePersonalBuildEvent(w, r)
	if parsedGitInfo == nil {
		return
	}
	s.Handle(w, r, parsedGitInfo)
}

// Parses Personal build requests and handle them if invalid by writing a response and logging issues.
func (s *BitbucketWebhookReceiver) handleParsePersonalBuildEvent(w http.ResponseWriter, r *http.Request) *ParseGitEvent {
	body := s.handleValidateBodyHasProperSignature(w, r)
	if body == nil {
		return nil
	}

	req := &RequestPersonalBuild{}
	if err := json.Unmarshal(body, &req); err != nil {
		msg := fmt.Sprintf("cannot parse JSON: %s", err)
		s.Logger.Errorf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return nil
	}

	repo := strings.ToLower(req.RepositorySlug)
	gitRef := req.BranchName
	gitFullRef := req.FullRefName
	project := strings.ToLower(s.Project)
	component := strings.TrimPrefix(repo, project+"-")
	var commitSHA string
	commentText := ""
	commitSHA = req.ToHash

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
		TriggerEvent: "personalBuild",
		Comment:      commentText,
	}
	return &ParseGitEvent{
		pInfo, commitSHA,
	}
}

func (s *BitbucketWebhookReceiver) handleValidateBodyHasProperSignature(w http.ResponseWriter, r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "could not read body"
		s.Logger.Errorf("%s: %s", msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return nil
	}

	if err := validatePayload(r.Header, body, []byte(s.WebhookSecret)); err != nil {
		msg := "failed to validate incoming request"
		s.Logger.Errorf("%s: %s", msg, err)
		http.Error(w, msg, http.StatusBadRequest)
		return nil
	}
	return body
}

// determineProject returns the project from given serverProject/projectParam.
func determineProject(serverProject, projectParam string) string {
	projectParam = strings.ToLower(projectParam)
	if len(projectParam) > 0 {
		return projectParam
	}
	return strings.ToLower(serverProject)
}
