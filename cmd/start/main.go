package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/repository"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

type options struct {
	bitbucketAccessToken     string
	bitbucketURL             string
	consoleURL               string
	pipelineRunName          string
	nexusURL                 string
	nexusUsername            string
	nexusPassword            string
	nexusTemporaryRepository string
	nexusPermanentRepository string
	project                  string
	environment              string
	version                  string
	prKey                    string
	prBase                   string
	gitRefSpec               string
	httpProxy                string
	httpsProxy               string
	noProxy                  string
	url                      string
	gitFullRef               string
	sslVerify                string
	submodules               string
	depth                    string
	debug                    bool
}

func main() {
	opts := options{}
	flag.StringVar(&opts.bitbucketAccessToken, "bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	flag.StringVar(&opts.bitbucketURL, "bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	flag.StringVar(&opts.project, "project", "", "project")
	flag.StringVar(&opts.environment, "environment", "", "environment")
	flag.StringVar(&opts.version, "version", "", "version")
	flag.StringVar(&opts.gitRefSpec, "git-ref-spec", "", "(optional) git refspec to fetch before checking out revision")
	flag.StringVar(&opts.prKey, "pr-key", "", "pull request key")
	flag.StringVar(&opts.prBase, "pr-base", "", "pull request base")
	flag.StringVar(&opts.httpProxy, "http-proxy", ".", "HTTP_PROXY")
	flag.StringVar(&opts.httpsProxy, "https-proxy", ".", "HTTPS_PROXY")
	flag.StringVar(&opts.noProxy, "no-proxy", ".", "NO_PROXY")
	flag.StringVar(&opts.url, "url", ".", "URL to clone")
	flag.StringVar(&opts.gitFullRef, "git-full-ref", "", "Git (full) ref to clone")
	flag.StringVar(&opts.sslVerify, "ssl-verify", "true", "defines if http.sslVerify should be set to true or false in the global git config")
	flag.StringVar(&opts.submodules, "submodules", "true", "defines if the resource should initialize and fetch the submodules")
	flag.StringVar(&opts.depth, "depth", "1", "performs a shallow clone where only the most recent commit(s) will be fetched")
	flag.StringVar(&opts.consoleURL, "console-url", os.Getenv("CONSOLE_URL"), "web console URL")
	flag.StringVar(&opts.pipelineRunName, "pipeline-run-name", "", "name of pipeline run")
	flag.StringVar(&opts.nexusURL, "nexus-url", os.Getenv("NEXUS_URL"), "Nexus URL")
	flag.StringVar(&opts.nexusUsername, "nexus-username", os.Getenv("NEXUS_USERNAME"), "Nexus username")
	flag.StringVar(&opts.nexusPassword, "nexus-password", os.Getenv("NEXUS_PASSWORD"), "Nexus password")
	flag.StringVar(&opts.nexusTemporaryRepository, "nexus-temporary-repository", os.Getenv("NEXUS_TEMPORARY_REPOSITORY"), "Nexus temporary repository")
	flag.StringVar(&opts.nexusPermanentRepository, "nexus-permanent-repository", os.Getenv("NEXUS_PERMANENT_REPOSITORY"), "Nexus permanent repository")
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
	err = cleanCache(checkoutDirFSB, removeFileOrDir)
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
		Environment:     opts.environment,
		Version:         opts.version,
		PullRequestBase: opts.prBase,
		PullRequestKey:  opts.prKey,
	}
	ctxt, err := checkoutAndAssembleContext(
		checkoutDir,
		opts.url,
		opts.gitFullRef,
		opts.gitRefSpec,
		opts.sslVerify,
		opts.submodules,
		opts.depth,
		baseCtxt,
		logger,
	)
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof("Assembled pipeline context: %+v", ctxt)

	logger.Infof("Setting Bitbucket build status to 'in progress' ...")
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: opts.bitbucketAccessToken,
		BaseURL:  opts.bitbucketURL,
		Logger:   logger,
	})
	pipelineRunURL := fmt.Sprintf(
		"%s/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s/",
		opts.consoleURL,
		ctxt.Namespace,
		opts.pipelineRunName,
	)
	err = bitbucketClient.BuildStatusCreate(ctxt.GitCommitSHA, bitbucket.BuildStatusCreatePayload{
		State:       bitbucket.BuildStatusInProgress,
		Key:         ctxt.GitCommitSHA,
		Name:        ctxt.GitCommitSHA,
		URL:         pipelineRunURL,
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
			subrepoGitFullRef, err := repository.BestMatchingBranch(bitbucketClient, ctxt.Project, subrepo, ctxt.Version)
			if err != nil {
				log.Fatal(err)
			}
			subrepoCtxt, err := checkoutAndAssembleContext(
				subrepoCheckoutDir,
				subrepoURL,
				subrepoGitFullRef,
				opts.gitRefSpec,
				opts.sslVerify,
				opts.submodules,
				opts.depth,
				baseCtxt,
				logger,
			)
			if err != nil {
				log.Fatal(err)
			}
			subrepoContexts = append(subrepoContexts, subrepoCtxt)
		}
	}

	if ctxt.Environment != "" {
		env, err := odsConfig.Environment(ctxt.Environment)
		if err != nil {
			log.Fatal(fmt.Sprintf("err during namespace extraction: %s", err))
		}
		err = applyVersionTags(logger, bitbucketClient, ctxt, subrepoContexts, env)
		if err != nil {
			log.Fatal(err)
		}
	}

	logger.Infof("Downloading any artifacts ...")
	// If there are subrepos, then all of them need to have a successful pipeline run.
	nexusClient, err := nexus.NewClient(&nexus.ClientConfig{
		BaseURL:  opts.nexusURL,
		Username: opts.nexusUsername,
		Password: opts.nexusPassword,
		Logger:   logger,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = downloadArtifacts(logger, nexusClient, ctxt, opts, pipelinectxt.ArtifactsPath)
	if err != nil {
		log.Fatal(err)
	}
	if len(subrepoContexts) > 0 {
		for _, src := range subrepoContexts {
			artifactsDir := filepath.Join(pipelinectxt.SubreposPath, src.Repository, pipelinectxt.ArtifactsPath)
			err = downloadArtifacts(logger, nexusClient, src, opts, artifactsDir)
			if err != nil {
				log.Fatal(err)
			}
			// check that a pipeline run exists
			// TODO: actually check for success.
			pipelineRunDir := filepath.Join(artifactsDir, pipelinectxt.PipelineRunsDir)
			if _, err := os.Stat(pipelineRunDir); os.IsNotExist(err) {
				log.Fatalf(
					"Pipeline runs with subrepos require a successful pipeline run "+
						"for all checked out subrepo commits, "+
						"however no such run was found for %s. "+
						"Re-run this pipeline once there is a successful pipeline run.", src.Repository,
				)
			}
		}
	}
}

func applyVersionTags(logger logging.LeveledLoggerInterface, bitbucketClient *bitbucket.Client, ctxt *pipelinectxt.ODSContext, subrepoContexts []*pipelinectxt.ODSContext, env *config.Environment) error {
	var tags []bitbucket.Tag
	tagVersion := ctxt.Version
	if env.Stage != config.DevStage {
		logger.Infof("Applying version tags ...")
		if tagVersion == pipelinectxt.WIP {
			return errors.New("when stage != dev, you must provide a version")
		}
		t, err := bitbucketClient.TagList(
			ctxt.Project,
			ctxt.Repository,
			bitbucket.TagListParams{
				FilterText: fmt.Sprintf("v%s", tagVersion),
			},
		)
		if err != nil {
			return fmt.Errorf("could not list tags in %s/%s: %w", ctxt.Project, ctxt.Repository, err)
		}
		tags = t.Values
	}
	if env.Stage == config.QAStage {
		if repository.TagListContainsFinalVersion(tags, tagVersion) {
			logger.Infof("Final version tag exists already.")
		} else {
			_, num := repository.LatestReleaseCandidate(tags, tagVersion)
			rcNum := num + 1
			tagName := fmt.Sprintf("v%s-rc.%d", tagVersion, rcNum)
			_, err := repository.CreateTag(bitbucketClient, ctxt, tagName)
			if err != nil {
				return fmt.Errorf("could not create tag %s in %s/%s: %w", tagName, ctxt.Project, ctxt.Repository, err)
			}
			// subrepos
			for _, sctxt := range subrepoContexts {
				_, err := repository.CreateTag(bitbucketClient, sctxt, tagName)
				if err != nil {
					return fmt.Errorf("could not create tag %s in %s/%s: %w", tagName, sctxt.Project, sctxt.Repository, err)
				}
			}
		}
	} else if env.Stage == config.ProdStage {
		if repository.TagListContainsFinalVersion(tags, tagVersion) {
			logger.Infof("Final version tag exists already.")
		} else {
			err := checkProdTagRequirements(tags, ctxt, tagVersion)
			if err != nil {
				return fmt.Errorf("cannot proceed to prod stage: %w", err)
			}
			tagName := fmt.Sprintf("v%s", tagVersion)
			_, err = repository.CreateTag(bitbucketClient, ctxt, tagName)
			if err != nil {
				return fmt.Errorf("could not create tag %s in %s/%s: %w", tagName, ctxt.Project, ctxt.Repository, err)
			}
			// subrepos
			for _, sctxt := range subrepoContexts {
				var subtags []bitbucket.Tag
				t, err := bitbucketClient.TagList(
					sctxt.Project,
					sctxt.Repository,
					bitbucket.TagListParams{
						FilterText: tagName,
					},
				)
				if err != nil {
					return fmt.Errorf("could not list tags in %s/%s: %w", sctxt.Project, sctxt.Repository, err)
				}
				subtags = t.Values
				err = checkProdTagRequirements(subtags, sctxt, tagVersion)
				if err != nil {
					return fmt.Errorf("cannot proceed to prod stage: %w", err)
				}
				_, err = repository.CreateTag(bitbucketClient, sctxt, tagName)
				if err != nil {
					return fmt.Errorf("could not create tag %s in %s/%s: %w", tagName, sctxt.Project, sctxt.Repository, err)
				}
			}
		}
	}
	return nil
}

func checkProdTagRequirements(tags []bitbucket.Tag, ctxt *pipelinectxt.ODSContext, version string) error {
	tag, _ := repository.LatestReleaseCandidate(tags, version)
	if tag == nil {
		return fmt.Errorf("no release candidate tag found for %s. Deploy to QA before deploying to Prod", version)
	}
	if tag.LatestCommit != ctxt.GitCommitSHA {
		return fmt.Errorf("latest release candidate tag for %s does not point to checked out commit, cowardly refusing to deploy", version)
	}
	return nil
}

func downloadArtifacts(
	logger logging.LeveledLoggerInterface,
	nexusClient *nexus.Client,
	ctxt *pipelinectxt.ODSContext,
	opts options,
	artifactsDir string) error {
	group := pipelinectxt.ArtifactGroupBase(ctxt)
	am, err := pipelinectxt.DownloadGroup(
		nexusClient,
		[]string{opts.nexusPermanentRepository, opts.nexusTemporaryRepository},
		group,
		artifactsDir,
		logger,
	)
	if err != nil {
		return err
	}
	return pipelinectxt.WriteJsonArtifact(am, artifactsDir, pipelinectxt.ArtifactsManifestFilename)

}

func checkoutAndAssembleContext(
	checkoutDir, url, gitFullRef, gitRefSpec, sslVerify, submodules, depth string,
	baseCtxt *pipelinectxt.ODSContext,
	logger logging.LeveledLoggerInterface) (*pipelinectxt.ODSContext, error) {
	absCheckoutDir, err := filepath.Abs(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof("Checking out %s@%s into %s ...", url, gitFullRef, absCheckoutDir)
	stdout, stderr, err := command.Run("/ko-app/git-init", []string{
		"-url", url,
		"-revision", gitFullRef,
		"-refspec", gitRefSpec,
		"-path", absCheckoutDir,
		"-sslVerify", sslVerify,
		"-submodules", submodules,
		"-depth", depth,
	})
	if err != nil {
		logger.Errorf(string(stderr))
		log.Fatal(err)
	}
	logger.Infof(string(stdout))

	// check git LFS state and maybe pull
	lfs, err := gitLfsInUse(logger, absCheckoutDir)
	if err != nil {
		log.Fatal(err)
	}
	if lfs {
		logger.Infof("Git LFS detected, enabling and pulling files...")
		err := gitLfsEnableAndPullFiles(logger, absCheckoutDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// write ODS cache
	sha, err := getCommitSHA(absCheckoutDir)
	if err != nil {
		log.Fatal(err)
	}
	ctxt := baseCtxt.Copy()
	ctxt.GitFullRef = gitFullRef
	ctxt.GitCommitSHA = sha
	err = ctxt.Assemble(absCheckoutDir)
	if err != nil {
		log.Fatal(err)
	}
	err = ctxt.WriteCache(absCheckoutDir)
	if err != nil {
		log.Fatal(err)
	}
	return ctxt, nil
}

func getCommitSHA(dir string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(dir, ".git/HEAD"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func gitLfsInUse(logger logging.LeveledLoggerInterface, dir string) (lfs bool, err error) {
	stdout, stderr, err := command.RunInDir("git", []string{"lfs", "ls-files", "--all"}, dir)
	if err != nil {
		return false, fmt.Errorf("cannot list git lfs files: %s (%w)", stderr, err)
	}
	return strings.TrimSpace(string(stdout)) != "", err
}

func gitLfsEnableAndPullFiles(logger logging.LeveledLoggerInterface, dir string) (err error) {
	stdout, stderr, err := command.RunInDir("git", []string{"lfs", "install"}, dir)
	if err != nil {
		return fmt.Errorf("cannot enable git lfs: %s (%w)", stderr, err)
	}
	logger.Infof(string(stdout))
	stdout, stderr, err = command.RunInDir("git", []string{"lfs", "pull"}, dir)
	if err != nil {
		return fmt.Errorf("cannot git pull lfs files: %s (%w)", stderr, err)
	}
	logger.Infof(string(stdout))
	return err
}
