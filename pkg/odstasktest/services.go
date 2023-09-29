package odstasktest

import (
	"flag"
	"os"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
)

var restartNexusFlag = flag.Bool("ods-restart-nexus", false, "Whether to force a restart of Nexus")
var restartSonarQubeFlag = flag.Bool("ods-restart-sonarqube", false, "Whether to force a restart of SonarQube")
var restartBitbucketFlag = flag.Bool("ods-restart-bitbucket", false, "Whether to force a restart of Bitbucket")

// StartNexus starts a Nexus instance in a Docker container (named
// ods-test-nexus). If a container of the same name already exists, it will be
// reused unless -ods-restart-nexus is passed.
func StartNexus() ttr.NamespaceOpt {
	flag.Parse()
	return runService("run-nexus.sh", *restartNexusFlag)
}

// StartSonarQube starts a SonarQube instance in a Docker container (named
// ods-test-sonarqube). If a container of the same name already exists, it will
// be reused unless -ods-restart-sonarqube is passed.
func StartSonarQube() ttr.NamespaceOpt {
	flag.Parse()
	return runService("run-sonarqube.sh", *restartSonarQubeFlag)
}

// StartBitbucket starts a Bitbucket instance in a Docker container (named
// ods-test-bitbucket-server). If a container of the same name already exists,
// it will be reused unless -ods-restart-bitbucket is passed.
func StartBitbucket() ttr.NamespaceOpt {
	flag.Parse()
	return runService("run-bitbucket.sh", *restartBitbucketFlag)
}

func runService(script string, restart bool) ttr.NamespaceOpt {
	return func(cc *ttr.ClusterConfig, nc *ttr.NamespaceConfig) error {
		args := []string{projectpath.RootedPath("scripts/" + script)}
		if !restart {
			args = append(args, "--reuse")
		}
		return command.Run("bash", args, []string{}, os.Stdout, os.Stderr)
	}
}
