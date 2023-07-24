//go:build external
// +build external

package tasks

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/opendevstack/ods-pipeline/internal/command"
	"github.com/opendevstack/ods-pipeline/internal/kubernetes"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/internal/random"
	"github.com/opendevstack/ods-pipeline/pkg/artifact"
	"github.com/opendevstack/ods-pipeline/pkg/config"
	"github.com/opendevstack/ods-pipeline/pkg/pipelinectxt"
	"github.com/opendevstack/ods-pipeline/pkg/tasktesting"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// To test deployment to external cluster, you must provide the token for a
// serviceaccount in an externa cluster, and a matching configuration like this:
//
// TODO: make this part of triggers, and supply
// tasks:
//   - name: deploy
//     taskRef:
//     kind: Task
//     name: ods-deploy-helm
//     params:
//   - name: namespace
//     value: foobar
//     apiServer: https://api.example.openshift.com:443
//     registryHost: default-route-openshift-image-registry.apps.example.openshiftapps.com
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
						Ref:        fmt.Sprintf("kind-registry.kind:5000/%s/%s:%s", ctxt.Namespace, imageStream, tag),
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

// buildAndPushImageWithLabel builds an image and pushes it to the registry.
// The used image tag equals the Git SHA that is being built, so the task
// will pick up the existing image.
// The image is labelled with "tasktestrun=true" so that it is possible to
// verify that the image has not been rebuild in the task.
func buildAndPushImageWithLabel(t *testing.T, ctxt *tasktesting.TaskRunContext, tag string, wsDir string) {
	t.Logf("Build image %s ahead of taskrun", tag)
	_, stderr, err := command.RunBuffered("docker", []string{
		"build", "--label", "tasktestrun=true", "-t", tag, filepath.Join(wsDir, "docker"),
	})
	if err != nil {
		t.Fatalf("could not build image: %s, stderr: %s", err, string(stderr))
	}
	_, stderr, err = command.RunBuffered("docker", []string{
		"push", tag,
	})
	if err != nil {
		t.Fatalf("could not push image: %s, stderr: %s", err, string(stderr))
	}
}
