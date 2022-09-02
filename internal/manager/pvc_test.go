package manager

import (
	"context"
	"testing"

	kubernetesClient "github.com/opendevstack/pipeline/internal/kubernetes"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreatePVCIfRequired(t *testing.T) {
	tests := map[string]struct {
		existingPVC string
		wantNewPVC  bool
	}{
		"PVC exists already": {
			existingPVC: "pvc",
			wantNewPVC:  false,
		},
		"different PVC exists already": {
			existingPVC: "pvc-other",
			wantNewPVC:  true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := &kubernetesClient.TestClient{
				PVCs: []*corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: tc.existingPVC,
						},
					},
				},
			}
			s := &Scheduler{
				KubernetesClient: c,
				StorageConfig: StorageConfig{
					Provisioner: "prov",
					ClassName:   "class",
					Size:        "1Gi",
				},
				Logger: &logging.LeveledLogger{Level: logging.LevelNull},
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			pData := PipelineConfig{
				PipelineInfo: PipelineInfo{
					Repository: "repo",
					GitRef:     "branch",
					Stage:      config.DevStage,
				},
				PVC: "pvc",
			}
			err := s.createPVCIfRequired(ctx, pData)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantNewPVC && len(c.CreatedPVCs) == 0 {
				t.Fatal("should have created a PVC")
			}
			if !tc.wantNewPVC && len(c.CreatedPVCs) > 0 {
				t.Fatal("should not have created a PVC")
			}
		})
	}
}
