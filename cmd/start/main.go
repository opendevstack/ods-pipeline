package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
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
	flag.StringVar(&opts.nexusTemporaryRepository, nexus.TemporaryRepositoryDefault, os.Getenv("NEXUS_TEMPORARY_REPOSITORY"), "Nexus temporary repository")
	flag.StringVar(&opts.nexusPermanentRepository, nexus.PermanentRepositoryDefault, os.Getenv("NEXUS_PERMANENT_REPOSITORY"), "Nexus permanent repository")
	flag.BoolVar(&opts.debug, "debug", (os.Getenv("DEBUG") == "true"), "debug mode")
	flag.Parse()

	checkoutDir := "."

	var logger logging.LeveledLoggerInterface
	if opts.debug {
		logger = &logging.LeveledLogger{Level: logging.LevelDebug}
	}

	fmt.Println("Cleaning checkout directory ...")
	err := deleteDirectoryContents(checkoutDir)
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
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Assembled pipeline context: %+v\n", ctxt)

	fmt.Println("Setting Bitbucket build status to 'in progress' ...")
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

	fmt.Println("Reading configuration from ods.y(a)ml ...")
	odsConfig, err := config.ReadFromDir(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}
	subrepoContexts := []*pipelinectxt.ODSContext{}
	if len(odsConfig.Repositories) > 0 {
		fmt.Println("Detected subrepos, checking out subrepos ...")
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
			// Checkout subrepo at the following refs, in order of specifity:
			// - release branch (if there is a version and the branch exists)
			// - configured branch (if configured in ODS config file)
			// - default branch
			subrepoGitFullRef := config.DefaultBranch
			if len(subrepo.Branch) > 0 {
				subrepoGitFullRef = subrepo.Branch
			}
			if ctxt.Version == pipelinectxt.WIP {
				releaseBranch, err := findReleaseBranch(bitbucketClient, ctxt.Project, subrepo.Name, ctxt.Version)
				if err != nil {
					log.Fatal(err)
				}
				if len(releaseBranch) > 0 {
					subrepoGitFullRef = releaseBranch
				}
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
			)
			if err != nil {
				log.Fatal(err)
			}
			subrepoContexts = append(subrepoContexts, subrepoCtxt)
		}
	}

	if len(ctxt.Environment) > 0 {
		env, err := odsConfig.Environment(ctxt.Environment)
		if err != nil {
			log.Fatal(fmt.Sprintf("err during namespace extraction: %s", err))
		}
		err = applyVersionTags(os.Stdout, bitbucketClient, ctxt, subrepoContexts, env)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Downloading any artifacts ...")
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
	err = downloadArtifacts(nexusClient, ctxt, opts, pipelinectxt.ArtifactsPath)
	if err != nil {
		log.Fatal(err)
	}
	if len(subrepoContexts) > 0 {
		for _, src := range subrepoContexts {
			artifactsDir := filepath.Join(pipelinectxt.SubreposPath, src.Repository, pipelinectxt.ArtifactsPath)
			err = downloadArtifacts(nexusClient, src, opts, artifactsDir)
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

func applyVersionTags(out io.Writer, bitbucketClient *bitbucket.Client, ctxt *pipelinectxt.ODSContext, subrepoContexts []*pipelinectxt.ODSContext, env *config.Environment) error {
	var tags []bitbucket.Tag
	tagVersion := ctxt.Version
	if env.Stage != config.DevStage {
		fmt.Fprintln(out, "Applying version tags ...")
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
		if tagListContainsFinalVersion(tags, tagVersion) {
			fmt.Fprintln(out, "Final version tag exists already.")
		} else {
			_, num := latestReleaseCandidate(tags, tagVersion)
			rcNum := num + 1
			tagName := fmt.Sprintf("v%s-rc.%d", tagVersion, rcNum)
			_, err := createTag(bitbucketClient, ctxt, tagName)
			if err != nil {
				return fmt.Errorf("could not create tag %s in %s/%s: %w", tagName, ctxt.Project, ctxt.Repository, err)
			}
			// subrepos
			for _, sctxt := range subrepoContexts {
				_, err := createTag(bitbucketClient, sctxt, tagName)
				if err != nil {
					return fmt.Errorf("could not create tag %s in %s/%s: %w", tagName, sctxt.Project, sctxt.Repository, err)
				}
			}
		}
	} else if env.Stage == config.ProdStage {
		if tagListContainsFinalVersion(tags, tagVersion) {
			fmt.Fprintln(out, "Final version tag exists already.")
		} else {
			err := checkProdTagRequirements(tags, ctxt, tagVersion)
			if err != nil {
				return fmt.Errorf("cannot proceed to prod stage: %w", err)
			}
			tagName := fmt.Sprintf("v%s", tagVersion)
			_, err = createTag(bitbucketClient, ctxt, tagName)
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
				_, err = createTag(bitbucketClient, sctxt, tagName)
				if err != nil {
					return fmt.Errorf("could not create tag %s in %s/%s: %w", tagName, sctxt.Project, sctxt.Repository, err)
				}
			}
		}
	}
	return nil
}

// findReleaseBranch returns the full Git ref of the release branch corresponding
// to given version. If none is found, it returns an empty string.
func findReleaseBranch(bitbucketClient *bitbucket.Client, projectKey, repositorySlug, version string) (string, error) {
	releaseBranch := fmt.Sprintf("release/%s", version)
	branchPage, err := bitbucketClient.BranchList(projectKey, repositorySlug, bitbucket.BranchListParams{
		FilterText:   fmt.Sprintf("release/%s", version),
		BoostMatches: true,
	})
	if err != nil {
		return "", err
	}
	for _, b := range branchPage.Values {
		if b.DisplayId == releaseBranch {
			return b.ID, nil
		}
	}
	return "", nil
}

func checkProdTagRequirements(tags []bitbucket.Tag, ctxt *pipelinectxt.ODSContext, version string) error {
	tag, _ := latestReleaseCandidate(tags, version)
	if tag == nil {
		return fmt.Errorf("no release candidate tag found for %s. Deploy to QA before deploying to Prod", version)
	}
	if tag.LatestCommit != ctxt.GitCommitSHA {
		return fmt.Errorf("latest release candidate tag for %s does not point to checked out commit, cowardly refusing to deploy", version)
	}
	return nil
}

func tagListContainsFinalVersion(tags []bitbucket.Tag, version string) bool {
	searchID := fmt.Sprintf("refs/tags/v%s", version)
	for _, t := range tags {
		if t.ID == searchID {
			return true
		}
	}
	return false
}

// latestReleaseCandidate returns the highest number of all tags of
// format "v<VERSION>-rc.<NUMBER>".
func latestReleaseCandidate(tags []bitbucket.Tag, version string) (*bitbucket.Tag, int) {
	var highestNumber int
	var latestTag *bitbucket.Tag
	prefix := fmt.Sprintf("refs/tags/v%s-rc.", version)
	for _, t := range tags {
		if strings.HasPrefix(t.ID, prefix) {
			i, err := strconv.Atoi(strings.TrimPrefix(t.ID, prefix))
			if err == nil && i > highestNumber {
				highestNumber = i
				latestTag = &t
			}
		}
	}
	return latestTag, highestNumber
}

func createTag(bitbucketClient *bitbucket.Client, ctxt *pipelinectxt.ODSContext, name string) (*bitbucket.Tag, error) {
	return bitbucketClient.TagCreate(
		ctxt.Project,
		ctxt.Repository,
		bitbucket.TagCreatePayload{
			Name:       name,
			StartPoint: ctxt.GitCommitSHA,
		},
	)
}

func getNexusURLs(nexusClient *nexus.Client, repository, group string) ([]string, error) {
	urls, err := nexusClient.Search(repository, group)
	if err != nil {
		return nil, err
	}
	if len(urls) > 0 {
		fmt.Printf("Found artifacts in repository %s inside group %s, downloading ...\n", repository, group)
	} else {
		fmt.Printf("No artifacts found in repository %s inside group %s.\n", repository, group)
	}
	return urls, nil
}

func downloadArtifacts(nexusClient *nexus.Client, ctxt *pipelinectxt.ODSContext, opts options, artifactsDir string) error {
	group := fmt.Sprintf("/%s/%s/%s", ctxt.Project, ctxt.Repository, ctxt.GitCommitSHA)
	// We want to target all artifacts underneath the group, hence the trailing '*'.
	nexusSearchGroup := fmt.Sprintf("%s/*", group)
	urls, err := getNexusURLs(nexusClient, opts.nexusPermanentRepository, nexusSearchGroup)
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		u, err := getNexusURLs(nexusClient, opts.nexusTemporaryRepository, nexusSearchGroup)
		if err != nil {
			return err
		}
		urls = u
	}
	for _, s := range urls {
		u, err := url.Parse(s)
		if err != nil {
			return err
		}
		urlPathParts := strings.Split(u.Path, group)
		fileWithSubPath := urlPathParts[1] // e.g. "/pipeline-runs/foo-zh9gt0.json"
		artifactsSubPath := filepath.Join(artifactsDir, path.Dir(fileWithSubPath))
		if _, err := os.Stat(artifactsSubPath); os.IsNotExist(err) {
			if err := os.MkdirAll(artifactsSubPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %s, error: %w", artifactsSubPath, err)
			}
		}
		outfile := filepath.Join(artifactsDir, fileWithSubPath)
		_, err = nexusClient.Download(s, outfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkoutAndAssembleContext(checkoutDir, url, gitFullRef, gitRefSpec, sslVerify, submodules, depth string, baseCtxt *pipelinectxt.ODSContext) (*pipelinectxt.ODSContext, error) {
	absCheckoutDir, err := filepath.Abs(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Checking out %s@%s into %s ...\n", url, gitFullRef, absCheckoutDir)
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
		log.Println(string(stderr))
		log.Fatal(err)
	}
	fmt.Println(string(stdout))

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

func deleteDirectoryContents(directory string) error {
	// Open the directory and read all its files.
	dirRead, err := os.Open(directory)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", directory, err)
	}
	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		return fmt.Errorf("could not read files in %s: %w", directory, err)
	}

	// Loop over the directory's files and remove them.
	for _, f := range dirFiles {
		filename := filepath.Join(directory, f.Name())
		err := os.RemoveAll(filename)
		if err != nil {
			return fmt.Errorf("could not remove file %s: %w", filename, err)
		}
	}
	return nil
}

func getCommitSHA(dir string) (string, error) {
	content, err := ioutil.ReadFile(filepath.Join(dir, ".git/HEAD"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}
