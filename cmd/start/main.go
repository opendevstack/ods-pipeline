package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/ods-pipeline/internal/tekton"
	"github.com/opendevstack/ods-pipeline/pkg/bitbucket"
	"github.com/opendevstack/ods-pipeline/pkg/config"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
	"github.com/opendevstack/ods-pipeline/pkg/nexus"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
)

type options struct {
	bitbucketAccessToken   string
	bitbucketURL           string
	consoleURL             string
	pipelineRunName        string
	nexusURL               string
	nexusUsername          string
	nexusPassword          string
	artifactSource         string
	project                string
	prKey                  string
	prBase                 string
	httpProxy              string
	httpsProxy             string
	noProxy                string
	url                    string
	gitFullRef             string
	submodules             string
	cloneDepth             string
	cacheBuildTasksForDays int
	debug                  bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.bitbucketAccessToken, "bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	flag.StringVar(&opts.bitbucketURL, "bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	flag.StringVar(&opts.project, "project", "", "project")
	flag.StringVar(&opts.prKey, "pr-key", "", "pull request key")
	flag.StringVar(&opts.prBase, "pr-base", "", "pull request base")
	flag.StringVar(&opts.httpProxy, "http-proxy", ".", "HTTP_PROXY")
	flag.StringVar(&opts.httpsProxy, "https-proxy", ".", "HTTPS_PROXY")
	flag.StringVar(&opts.noProxy, "no-proxy", ".", "NO_PROXY")
	flag.StringVar(&opts.url, "url", ".", "URL to clone")
	flag.StringVar(&opts.gitFullRef, "git-full-ref", "", "Git (full) ref to clone")
	flag.StringVar(&opts.submodules, "submodules", "true", "defines if the resource should initialize and fetch the submodules")
	flag.StringVar(&opts.cloneDepth, "clone-depth", "", "perform a shallow clone where only the most recent commit(s) will be fetched")
	flag.IntVar(&opts.cacheBuildTasksForDays, "cache-build-tasks-for-days", 7, "the number of days build outputs are cached. A negative number can be used to clear the cache.")
	flag.StringVar(&opts.consoleURL, "console-url", os.Getenv("CONSOLE_URL"), "web console URL")
	flag.StringVar(&opts.pipelineRunName, "pipeline-run-name", "", "name of pipeline run")
	flag.StringVar(&opts.nexusURL, "nexus-url", os.Getenv("NEXUS_URL"), "Nexus URL")
	flag.StringVar(&opts.nexusUsername, "nexus-username", os.Getenv("NEXUS_USERNAME"), "Nexus username")
	flag.StringVar(&opts.nexusPassword, "nexus-password", os.Getenv("NEXUS_PASSWORD"), "Nexus password")
	flag.StringVar(&opts.artifactSource, "artifact-source", "", "Artifacts source repository")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "debug mode")
	flag.Parse()

	checkoutDir := "."

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	} else {
		logger = &logging.LeveledLogger{Level: logging.LevelInfo}
	}

	logger.Infof("Cleaning checkout directory ...")
	checkoutDirFSB := FileSystemBase{os.DirFS(checkoutDir), checkoutDir}
	err := deleteDirectoryContentsSpareCache(checkoutDirFSB, removeFileOrDir)
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof("Cleaning cache at %s ...", odsCacheDirName)
	err = cleanCache(checkoutDirFSB, removeFileOrDir, opts.cacheBuildTasksForDays)
	if err != nil {
		log.Fatal(err)
	}

	// set proxy env vars
	if len(opts.httpProxy) > 0 {
		err = os.Setenv("HTTP_PROXY", opts.httpProxy)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(opts.httpsProxy) > 0 {
		err = os.Setenv("HTTPS_PROXY", opts.httpsProxy)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(opts.noProxy) > 0 {
		err = os.Setenv("NO_PROXY", opts.noProxy)
		if err != nil {
			log.Fatal(err)
		}
	}

	baseCtxt := &pipelinectxt.ODSContext{
		Project:         opts.project,
		PullRequestBase: opts.prBase,
		PullRequestKey:  opts.prKey,
	}
	ctxt, err := checkoutAndAssembleContext(
		checkoutDir,
		opts.url,
		opts.gitFullRef,
		opts,
		baseCtxt,
		logger,
	)
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof("Assembled pipeline context: %+v", ctxt)

	logger.Infof("Setting Bitbucket build status to 'in progress' ...")
	bitbucketClient, err := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: opts.bitbucketAccessToken,
		BaseURL:  opts.bitbucketURL,
		Logger:   logger,
	})
	if err != nil {
		log.Fatal("bitbucket client:", err)
	}

	prURL, err := tekton.PipelineRunURL(opts.consoleURL, ctxt.Namespace, opts.pipelineRunName)
	if err != nil {
		log.Fatal("pipeline run URL:", err)
	}

	err = bitbucketClient.BuildStatusCreate(ctxt.GitCommitSHA, bitbucket.BuildStatusCreatePayload{
		State:       bitbucket.BuildStatusInProgress,
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         prURL,
		Description: "ODS Pipeline Build",
	})
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof("Reading configuration from ods.y(a)ml ...")
	odsConfig, err := config.ReadFromDir(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}
	subrepoContexts := []*pipelinectxt.ODSContext{}
	if len(odsConfig.Repositories) > 0 {
		logger.Infof("Detected subrepos, checking out subrepos ...")
		for _, subrepo := range odsConfig.Repositories {
			subrepoCheckoutDir := filepath.Join(pipelinectxt.SubreposPath, subrepo.Name)
			err = os.MkdirAll(subrepoCheckoutDir, 0755)
			if err != nil {
				log.Fatalf("could not create checkout dir for subrepo %s: %s", subrepo.Name, err)
			}
			subrepoURL := subrepo.URL
			if len(subrepoURL) == 0 {
				subrepoURL = strings.Replace(
					opts.url,
					fmt.Sprintf("/%s.git", ctxt.Repository),
					fmt.Sprintf("/%s.git", subrepo.Name),
					1,
				)
			}
			subrepoCtxt, err := checkoutAndAssembleContext(
				subrepoCheckoutDir,
				subrepoURL,
				findBestMatchingRef(subrepo),
				opts,
				baseCtxt,
				logger,
			)
			if err != nil {
				log.Fatal(err)
			}
			logger.Infof("Assembled context for sub-repo %q: %+v", subrepo.Name, subrepoCtxt)
			subrepoContexts = append(subrepoContexts, subrepoCtxt)
		}
	}

	if err := os.MkdirAll(pipelinectxt.ArtifactsPath, 0755); err != nil {
		log.Fatalf("could not create %s: %s", pipelinectxt.ArtifactsPath, err)
	}
	if opts.artifactSource != "" {
		logger.Infof("Downloading any artifacts from %s ...", opts.nexusURL)
		nexusClient, err := nexus.NewClient(&nexus.ClientConfig{
			BaseURL:  opts.nexusURL,
			Username: opts.nexusUsername,
			Password: opts.nexusPassword,
			Logger:   logger,
		})
		if err != nil {
			log.Fatal(err)
		}
		err = downloadArtifacts(logger, nexusClient, ctxt, opts.artifactSource, pipelinectxt.ArtifactsPath)
		if err != nil {
			log.Fatal(err)
		}
		// If there are subrepos, then all of them need to have a successful pipeline run
		// for the currently checkout out commit.
		for _, src := range subrepoContexts {
			artifactsDir := filepath.Join(pipelinectxt.SubreposPath, src.Repository, pipelinectxt.ArtifactsPath)
			err = downloadArtifacts(logger, nexusClient, src, opts.artifactSource, artifactsDir)
			if err != nil {
				log.Fatal(err)
			}
			// check that a pipeline run exists
			// TODO: actually check for success.
			pipelineRunDir := filepath.Join(artifactsDir, pipelinectxt.PipelineRunsDir)
			if _, err := os.Stat(pipelineRunDir); os.IsNotExist(err) {
				log.Fatalf(
					"Pipeline runs with subrepos require a successful pipeline run artifact "+
						"for all checked out subrepo commits, however no such artifact was found for %s. "+
						"Re-run this pipeline once there is a successful pipeline run.", src.Repository,
				)
			}
		}
	} else {
		err := writeEmptyArtifactManifests(subrepoContexts)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func writeEmptyArtifactManifests(subrepoContexts []*pipelinectxt.ODSContext) error {
	emptyManifest := pipelinectxt.NewArtifactsManifest("")
	err := pipelinectxt.WriteJsonArtifact(emptyManifest, pipelinectxt.ArtifactsPath, pipelinectxt.ArtifactsManifestFilename)
	if err != nil {
		return fmt.Errorf("write repo empty manifest: %w", err)
	}
	for _, src := range subrepoContexts {
		artifactsDir := filepath.Join(pipelinectxt.SubreposPath, src.Repository, pipelinectxt.ArtifactsPath)
		err := pipelinectxt.WriteJsonArtifact(emptyManifest, artifactsDir, pipelinectxt.ArtifactsManifestFilename)
		if err != nil {
			return fmt.Errorf("write subrepo %s empty manifest: %w", src.Repository, err)
		}
	}
	return nil
}

// findBestMatchingRef returns a full Git ref, from either tag, branch or default.
func findBestMatchingRef(subrepo config.Repository) string {
	if subrepo.Tag != "" {
		if !strings.HasPrefix(subrepo.Tag, "refs/tags/") {
			return fmt.Sprintf("refs/tags/%s", subrepo.Tag)
		}
		return subrepo.Tag
	}
	if subrepo.Branch != "" {
		if !strings.HasPrefix(subrepo.Branch, "refs/heads/") {
			return fmt.Sprintf("refs/heads/%s", subrepo.Branch)
		}
		return subrepo.Branch
	}
	return config.DefaultBranch
}

func downloadArtifacts(
	logger logging.LeveledLoggerInterface,
	nexusClient *nexus.Client,
	ctxt *pipelinectxt.ODSContext,
	artifactSource, artifactsDir string) error {
	group := pipelinectxt.ArtifactGroupBase(ctxt)
	am, err := pipelinectxt.DownloadGroup(
		nexusClient,
		artifactSource,
		group,
		artifactsDir,
		logger,
	)
	if err != nil {
		return fmt.Errorf("download group %s: %w", group, err)
	}
	return pipelinectxt.WriteJsonArtifact(am, artifactsDir, pipelinectxt.ArtifactsManifestFilename)

}

func checkoutAndAssembleContext(
	checkoutDir, url, gitFullRef string, opts options,
	baseCtxt *pipelinectxt.ODSContext,
	logger logging.LeveledLoggerInterface) (ctxt *pipelinectxt.ODSContext, err error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return
	}
	// Change back to working dir after checkout.
	defer func(wd string) {
		if err != nil { // if there are previous errors, give them predence.
			return
		}
		err = os.Chdir(wd)
	}(workingDir)

	absCheckoutDir, err := filepath.Abs(checkoutDir)
	if err != nil {
		return nil, fmt.Errorf("absolute path: %w", err)
	}
	logger.Infof("Checking out %s@%s into %s ...", url, gitFullRef, absCheckoutDir)
	if err := os.Chdir(absCheckoutDir); err != nil {
		return nil, fmt.Errorf("change dir: %w", err)
	}
	if err := gitCheckout(gitCheckoutParams{
		repoURL:              url,
		bitbucketAccessToken: opts.bitbucketAccessToken,
		recurseSubmodules:    opts.submodules,
		depth:                opts.cloneDepth,
		fullRef:              gitFullRef,
	}); err != nil {
		return nil, fmt.Errorf("git checkout: %w", err)
	}

	odsPipelineIgnoreFile := filepath.Join(absCheckoutDir, ".git", "info", "exclude")
	if err := pipelinectxt.WriteGitIgnore(odsPipelineIgnoreFile); err != nil {
		return nil, fmt.Errorf("write git ignore: %w", err)
	}
	logger.Infof("Wrote gitignore exclude at %s", odsPipelineIgnoreFile)

	// check git LFS state and maybe pull
	lfs, err := gitLfsInUse(logger, absCheckoutDir)
	if err != nil {
		return nil, fmt.Errorf("check if git LFS is in use: %w", err)
	}
	if lfs {
		logger.Infof("Git LFS detected, enabling and pulling files...")
		err := gitLfsEnableAndPullFiles(logger, absCheckoutDir)
		if err != nil {
			return nil, fmt.Errorf("git LFS enable and pull: %w", err)
		}
	}

	// write ODS cache
	sha, err := getCommitSHA(absCheckoutDir)
	if err != nil {
		return nil, fmt.Errorf("commit SHA: %w", err)
	}
	ctxt = baseCtxt.Copy()
	ctxt.GitFullRef = gitFullRef
	ctxt.GitCommitSHA = sha
	err = ctxt.Assemble(absCheckoutDir)
	if err != nil {
		return nil, fmt.Errorf("assemble ODS context: %w", err)
	}
	err = ctxt.WriteCache(absCheckoutDir)
	if err != nil {
		return nil, fmt.Errorf("write ODS context cache: %w", err)
	}
	return
}
