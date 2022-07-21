package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	intrepo "github.com/opendevstack/pipeline/internal/repository"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
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

type ParseGitEvent struct {
	PInfo     PipelineInfo
	CommitSHA string
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
func (s *BitbucketWebhookReceiver) Handle(w http.ResponseWriter, r *http.Request, parsedGitInfo *ParseGitEvent) {
	pInfo := parsedGitInfo.PInfo
	commitSHA := parsedGitInfo.CommitSHA
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
		PVC:          makePVCName(pInfo.Component),
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
