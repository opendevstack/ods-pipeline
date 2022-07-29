package manager

import (
	"context"
	"testing"
	"time"

	kubernetesClient "github.com/opendevstack/pipeline/internal/kubernetes"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSchedule(t *testing.T) {

	cfg := PipelineConfig{
		PipelineInfo: PipelineInfo{
			Repository: "repo",
		},
		PVC: "pvc",
	}

	tests := map[string]struct {
		requestBodyFixture string
		kubernetesClient   *kubernetesClient.TestClient
		tektonClient       *tektonClient.TestClient
		wantQueuedRun      bool
		check              func(t *testing.T, kc *kubernetesClient.TestClient, tc *tektonClient.TestClient)
	}{
		"creates pipeline run and starts it": {
			kubernetesClient: &kubernetesClient.TestClient{},
			tektonClient:     &tektonClient.TestClient{},
			wantQueuedRun:    false,
		},
		"creates pipeline run and queues if necessary": {
			kubernetesClient: &kubernetesClient.TestClient{
				PVCs: []*corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{Name: "pvc"}},
				},
			},
			tektonClient: &tektonClient.TestClient{
				PipelineRuns: []*tekton.PipelineRun{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "pipeline-abcdef",
							CreationTimestamp: metav1.Time{Time: time.Now().Add(-2 * time.Second)},
						},
					},
				},
			},
			wantQueuedRun: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := &Scheduler{
				TektonClient:     tc.tektonClient,
				KubernetesClient: tc.kubernetesClient,
				StorageConfig: StorageConfig{
					Provisioner: "prov",
					ClassName:   "class",
					Size:        "1Gi",
				},
				Logger:         &logging.LeveledLogger{Level: logging.LevelNull},
				TriggeredRepos: make(chan string, 100),
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			gotQueued := s.schedule(ctx, cfg)

			if len(tc.tektonClient.CreatedPipelineRuns) != 1 {
				t.Fatal("one pipeline run should have been created")
			}
			if gotTriggered := <-s.TriggeredRepos; gotTriggered != "repo" {
				t.Fatal("channel should have received the repository name")
			}
			if tc.wantQueuedRun != gotQueued {
				t.Fatalf("want queued: %v, got: %v", tc.wantQueuedRun, gotQueued)
			}
		})
	}
}
