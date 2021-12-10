//go:build external
// +build external

package tasks

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/pipeline/internal/command"
	"github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/internal/random"
	"github.com/opendevstack/pipeline/pkg/artifact"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/pipeline/pkg/tasktesting"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// To test deployment to external cluster, you must provide the token for a
// serviceaccount in an externa cluster, and a matching configuration like this:
//
// environments:
// - name: dev
//   stage: dev
//   namespace: foobar
//   apiServer: https://api.example.openshift.com:443
//   registryHost: default-route-openshift-image-registry.apps.example.openshiftapps.com
//
// You do not need to specify "apiCredentialsSecret", it is set automatically to
// the secret created from the token given via -external-cluster-token.
//
// The test will not create or delete any namespaces. It will install a Helm
// release into the specified namespace, and delete the release again after the
// test. The Helm release and related resources are prefixed with the temporary
// workspace directory (e.g. "workspace-476709422") so any clashes even in none-
// empty namespace are very unlikely. Nonetheless, it is always recommended to
// use an empty namespace setup solely for the purpose of testing.
var (
	externalClusterTokenFlag  = flag.String("external-cluster-token", "", "Token of serviceaccount in external cluster")
	externalClusterConfigFlag = flag.String("external-cluster-config", "", "ods.yaml describing external cluster")
)

func TestTaskODSDeployHelmExternal(t *testing.T) {
	var externalEnv *config.Environment
	var imageStream string
	runTaskTestCases(t,
		"ods-deploy-helm",
		[]tasktesting.Service{},
		map[string]tasktesting.TestCase{
			"external deployment": {
				Timeout:             10 * time.Minute,
				WorkspaceDirMapping: map[string]string{"source": "helm-sample-app"},
				PreRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					wsDir := ctxt.Workspaces["source"]

					if *externalClusterConfigFlag == "" || *externalClusterTokenFlag == "" {
						t.Fatal(
							"-external-cluster-token and -external-cluster-config are required to run this test. " +
								"Use -short to skip this test.",
						)
					}

					t.Log("Create token secret for external cluster")
					secret, err := kubernetes.CreateSecret(ctxt.Clients.KubernetesClientSet, ctxt.Namespace, &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{Name: "ext"},
						Data: map[string][]byte{
							"token": []byte(*externalClusterTokenFlag),
						},
					})
					if err != nil {
						t.Fatal(err)
					}

					t.Log("Create private key secret for sample app")
					createSampleAppPrivateKeySecret(t, ctxt.Clients.KubernetesClientSet, ctxt.Namespace)

					t.Log("Read ods.yaml from flag and write into working dir")
					externalClusterConfig := *externalClusterConfigFlag
					if !filepath.IsAbs(externalClusterConfig) {
						externalClusterConfig = filepath.Join(projectpath.Root, externalClusterConfig)
					}
					o, err := config.ReadFromFile(externalClusterConfig)
					if err != nil {
						t.Fatal(err)
					}
					externalEnv := o.Environments[0]
					externalEnv.APICredentialsSecret = secret.Name
					externalEnv.APIToken = *externalClusterTokenFlag
					o.Environments[0] = externalEnv
					err = createODSYML(wsDir, o)
					if err != nil {
						t.Fatal(err)
					}

					imageStream = random.PseudoString()
					tag := "latest"
					fullTag := fmt.Sprintf("localhost:5000/%s/%s:%s", ctxt.Namespace, imageStream, tag)
					buildAndPushImageWithLabel(t, ctxt, fullTag, wsDir)
					ia := artifact.Image{
						Image:      fmt.Sprintf("kind-registry.kind:5000/%s/%s:%s", ctxt.Namespace, imageStream, tag),
						Registry:   "kind-registry.kind:5000",
						Repository: ctxt.Namespace,
						Name:       imageStream,
						Tag:        tag,
						Digest:     "abc",
					}
					imageArtifactFilename := fmt.Sprintf("%s.json", imageStream)
					err = pipelinectxt.WriteJsonArtifact(ia, filepath.Join(wsDir, pipelinectxt.ImageDigestsPath), imageArtifactFilename)
					if err != nil {
						t.Fatal(err)
					}

					ctxt.ODS = tasktesting.SetupGitRepo(t, ctxt.Namespace, wsDir)
				},
				WantRunSuccess: true,
				PostRunFunc: func(t *testing.T, ctxt *tasktesting.TaskRunContext) {
					t.Log("Check image")
					_, _, err := command.Run("skopeo", []string{
						"inspect",
						fmt.Sprintf("--registry-token=%s", externalEnv.APIToken),
						fmt.Sprintf("docker://%s/%s/%s:%s", externalEnv.RegistryHost, ctxt.Namespace, imageStream, "latest"),
					})
					if err != nil {
						t.Fatal(err)
					}
					t.Log("Remove Helm release again")
					command.Run("helm", []string{
						fmt.Sprintf("--kube-apiserver=%s", externalEnv.APIServer),
						fmt.Sprintf("--kube-token=%s", externalEnv.APIToken),
						fmt.Sprintf("--namespace=%s", externalEnv.Namespace),
						"uninstall",
						ctxt.ODS.Component,
					})
				},
			},
		},
	)
}
