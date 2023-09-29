package odstasktest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	ttr "github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
)

// InstallOptions configure the installation of ODS Pipeline.
type InstallOptions struct {
	// PrivateCert specifies if services should be accessed through TLS
	// with a private certificate.
	PrivateCert bool
}

// InstallODSPipeline installs the ODS Pipeline Helm chart in the namespace
// given in NamespaceConfig.
func InstallODSPipeline(opts *InstallOptions) ttr.NamespaceOpt {
	if opts == nil {
		opts = &InstallOptions{PrivateCert: false}
	}
	return func(cc *ttr.ClusterConfig, nc *ttr.NamespaceConfig) error {
		return installCDNamespaceResources(nc.Name, "pipeline", opts.PrivateCert)
	}
}

func installCDNamespaceResources(ns, serviceaccount string, privateCert bool) error {
	scriptArgs := []string{"-n", ns, "-s", serviceaccount, "--no-diff"}
	// if testing.Verbose() {
	// 	scriptArgs = append(scriptArgs, "-v")
	// }
	if privateCert {
		// Insert as first flag because install-inside-kind.sh won't recognize it otherwise.
		scriptArgs = append(
			[]string{fmt.Sprintf("--private-cert=%s", filepath.Join(projectpath.Root, "test/testdata/private-cert/tls.crt"))},
			scriptArgs...,
		)
	}

	return command.Run(
		"bash",
		append([]string{filepath.Join(projectpath.Root, "scripts/install-inside-kind.sh")}, scriptArgs...),
		[]string{},
		os.Stdout,
		os.Stderr,
	)
}
