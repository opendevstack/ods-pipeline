package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/nexus"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
)

func main() {
	bitbucketAccessTokenFlag := flag.String("bitbucket-access-token", os.Getenv("BITBUCKET_ACCESS_TOKEN"), "bitbucket-access-token")
	bitbucketURLFlag := flag.String("bitbucket-url", os.Getenv("BITBUCKET_URL"), "bitbucket-url")
	namespaceFlag := flag.String("namespace", "", "namespace")
	projectFlag := flag.String("project", "", "project")
	environmentFlag := flag.String("environment", "", "environment")
	versionFlag := flag.String("version", "", "version")
	gitRefSpecFlag := flag.String("git-ref-spec", "", "(optional) git refspec to fetch before checking out revision")
	prKeyFlag := flag.String("pr-key", "", "pull request key")
	prBaseFlag := flag.String("pr-base", "", "pull request base")
	httpProxyFlag := flag.String("http-proxy", ".", "HTTP_PROXY")
	httpsProxyFlag := flag.String("https-proxy", ".", "HTTPS_PROXY")
	noProxyFlag := flag.String("no-proxy", ".", "NO_PROXY")
	urlFlag := flag.String("url", ".", "URL to clone")
	gitFullRefFlag := flag.String("git-full-ref", "", "Git (full) ref to clone")
	sslVerifyFlag := flag.String("ssl-verify", "true", "defines if http.sslVerify should be set to true or false in the global git config")
	submodulesFlag := flag.String("submodules", "true", "defines if the resource should initialize and fetch the submodules")
	depthFlag := flag.String("depth", "1", "performs a shallow clone where only the most recent commit(s) will be fetched")
	consoleURLFlag := flag.String("console-url", os.Getenv("CONSOLE_URL"), "web console URL")
	pipelineRunNameFlag := flag.String("pipeline-run-name", "", "name of pipeline run")
	nexusURLFlag := flag.String("nexus-url", os.Getenv("NEXUS_URL"), "Nexus URL")
	nexusUsernameFlag := flag.String("nexus-username", os.Getenv("NEXUS_USERNAME"), "Nexus username")
	nexusPasswordFlag := flag.String("nexus-password", os.Getenv("NEXUS_PASSWORD"), "Nexus password")
	nexusTemporaryRepositoryFlag := flag.String("nexus-temporary-repository", os.Getenv("NEXUS_TEMPORARY_REPOSITORY"), "Nexus temporary repository")
	//nexusPermanentRepositoryFlag := flag.String("nexus-permanent-repository", os.Getenv("NEXUS_PERMANENT_REPOSITORY"), "Nexus permanent repository")
	flag.Parse()

	checkoutDir := "."

	fmt.Println("Cleaning checkout directory ...")
	err := deleteDirectoryContents(checkoutDir)
	if err != nil {
		log.Fatal(err)
	}

	// set proxy env vars
	if len(*httpProxyFlag) > 0 {
		err = os.Setenv("HTTP_PROXY", *httpProxyFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(*httpsProxyFlag) > 0 {
		err = os.Setenv("HTTPS_PROXY", *httpsProxyFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(*noProxyFlag) > 0 {
		err = os.Setenv("NO_PROXY", *noProxyFlag)
		if err != nil {
			log.Fatal(err)
		}
	}

	baseCtxt := &pipelinectxt.ODSContext{
		Namespace:       *namespaceFlag,
		Project:         *projectFlag,
		Environment:     *environmentFlag,
		Version:         *versionFlag,
		PullRequestBase: *prBaseFlag,
		PullRequestKey:  *prKeyFlag,
	}
	ctxt, err := checkoutAndAssembleContext(
		checkoutDir,
		*urlFlag,
		*gitFullRefFlag,
		*gitRefSpecFlag,
		*sslVerifyFlag,
		*submodulesFlag,
		*depthFlag,
		baseCtxt,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Assembled pipeline context: %+v\n", ctxt)

	fmt.Println("Setting Bitbucket build status to 'in progress' ...")
	bitbucketClient := bitbucket.NewClient(&bitbucket.ClientConfig{
		APIToken: *bitbucketAccessTokenFlag,
		BaseURL:  *bitbucketURLFlag,
	})
	pipelineRunURL := fmt.Sprintf(
		"%s/k8s/ns/%s/tekton.dev~v1beta1~PipelineRun/%s/",
		*consoleURLFlag,
		ctxt.Namespace,
		*pipelineRunNameFlag,
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
					*urlFlag,
					fmt.Sprintf("/%s.git", ctxt.Repository),
					fmt.Sprintf("/%s.git", subrepo.Name),
					1,
				)
			}
			subrepoGitFullRef := subrepo.Branch
			if len(subrepoGitFullRef) == 0 {
				subrepoGitFullRef = config.DefaultBranch
			}
			subrepoCtxt, err := checkoutAndAssembleContext(
				subrepoCheckoutDir,
				subrepoURL,
				subrepoGitFullRef,
				*gitRefSpecFlag,
				*sslVerifyFlag,
				*submodulesFlag,
				*depthFlag,
				baseCtxt,
			)
			if err != nil {
				log.Fatal(err)
			}
			subrepoContexts = append(subrepoContexts, subrepoCtxt)
		}
	}

	fmt.Println("Downloading any artifacts ...")
	// If there are subrepos, then all of them need to have a successful pipeline run.
	nexusClient, err := nexus.NewClient(&nexus.ClientConfig{
		BaseURL:    *nexusURLFlag,
		Username:   *nexusUsernameFlag,
		Password:   *nexusPasswordFlag,
		Repository: *nexusTemporaryRepositoryFlag,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = downloadArtifacts(nexusClient, ctxt, pipelinectxt.ArtifactsPath)
	if err != nil {
		log.Fatal(err)
	}
	if len(subrepoContexts) > 0 {
		for _, src := range subrepoContexts {
			artifactsDir := filepath.Join(pipelinectxt.SubreposPath, src.Repository, pipelinectxt.ArtifactsPath)
			err = downloadArtifacts(nexusClient, src, artifactsDir)
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

func downloadArtifacts(nexusClient *nexus.Client, ctxt *pipelinectxt.ODSContext, artifactsDir string) error {
	group := fmt.Sprintf("/%s/%s/%s", ctxt.Project, ctxt.Repository, ctxt.GitCommitSHA)
	// We want to target all artifacts underneath the group, hence the '*'.
	urls, err := nexusClient.Search(group + "/*")
	if err != nil {
		return err
	}
	if len(urls) > 0 {
		fmt.Printf("Found artifacts in repository %s inside group %s, downloading ...\n", nexusClient.Repository(), group)
	} else {
		fmt.Printf("No artifacts found in repository %s inside group %s.\n", nexusClient.Repository(), group)
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
