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
			Name: "pipeline",
		},
		PVC: "pvc",
	}

	tests := map[string]struct {
		requestBodyFixture  string
		kubernetesClient    *kubernetesClient.TestClient
		tektonClient        *tektonClient.TestClient
		wantCreatedPipeline bool
		wantUpdatedPipeline bool
		wantQueuedRun       bool
		check               func(t *testing.T, kc *kubernetesClient.TestClient, tc *tektonClient.TestClient)
	}{
		"creates a new pipeline and starts run": {
			kubernetesClient:    &kubernetesClient.TestClient{},
			tektonClient:        &tektonClient.TestClient{},
			wantCreatedPipeline: true,
			wantUpdatedPipeline: false,
			wantQueuedRun:       false,
		},
		"updates an existing pipeline and starts run": {
			kubernetesClient: &kubernetesClient.TestClient{
				PVCs: []*corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{Name: "pvc"}},
				},
			},
			tektonClient: &tektonClient.TestClient{
				Pipelines: []*tekton.Pipeline{
					{ObjectMeta: metav1.ObjectMeta{Name: "pipeline"}},
				},
			},
			wantCreatedPipeline: false,
			wantUpdatedPipeline: true,
			wantQueuedRun:       false,
		},
		"updates an existing pipeline and queues run if required": {
			kubernetesClient: &kubernetesClient.TestClient{
				PVCs: []*corev1.PersistentVolumeClaim{
					{ObjectMeta: metav1.ObjectMeta{Name: "pvc"}},
				},
			},
			tektonClient: &tektonClient.TestClient{
				Pipelines: []*tekton.Pipeline{
					{ObjectMeta: metav1.ObjectMeta{Name: "pipeline"}},
				},
				PipelineRuns: []*tekton.PipelineRun{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "pipeline-abcdef",
							CreationTimestamp: metav1.Time{Time: time.Now().Add(-2 * time.Second)},
						},
					},
				},
			},
			wantCreatedPipeline: false,
			wantUpdatedPipeline: true,
			wantQueuedRun:       true,
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
				Logger: &logging.LeveledLogger{Level: logging.LevelNull},
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			gotQueued := s.schedule(ctx, cfg)

			if (tc.wantCreatedPipeline && len(tc.tektonClient.CreatedPipelines) != 1) ||
				(!tc.wantCreatedPipeline && len(tc.tektonClient.CreatedPipelines) != 0) {
				t.Fatal("one pipeline should have been created")
			}
			if (tc.wantUpdatedPipeline && len(tc.tektonClient.UpdatedPipelines) != 1) ||
				(!tc.wantUpdatedPipeline && len(tc.tektonClient.UpdatedPipelines) != 0) {
				t.Fatal("one pipeline should have been updated")
			}
			if len(tc.tektonClient.CreatedPipelineRuns) != 1 {
				t.Fatal("one pipeline run should have been created")
			}
			if tc.wantQueuedRun != gotQueued {
				t.Fatalf("want queued: %v, got: %v", tc.wantQueuedRun, gotQueued)
			}
		})
	}
}
