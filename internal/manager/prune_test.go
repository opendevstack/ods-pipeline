package manager

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPrunerConfig(t *testing.T) {
	tclient := &tektonClient.TestClient{}
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}
	minKeepHours := 1
	maxKeepRuns := 1
	tests := map[string]struct {
		tektonClient tektonClient.ClientInterface
		logger       logging.LeveledLoggerInterface
		minKeepHours int
		maxKeepRuns  int
		wantErr      string
	}{
		"no tekton client": {
			tektonClient: nil,
			logger:       logger,
			minKeepHours: minKeepHours,
			maxKeepRuns:  maxKeepRuns,
			wantErr:      "tekton client is required",
		},
		"no logger": {
			tektonClient: tclient,
			logger:       nil,
			minKeepHours: minKeepHours,
			maxKeepRuns:  maxKeepRuns,
			wantErr:      "logger is required",
		},
		"wrong minKeepHours": {
			tektonClient: tclient,
			logger:       logger,
			minKeepHours: 0,
			maxKeepRuns:  maxKeepRuns,
			wantErr:      "minKeepHours must be at least 1, got 0",
		},
		"wrong maxKeepRuns": {
			tektonClient: tclient,
			logger:       logger,
			minKeepHours: minKeepHours,
			maxKeepRuns:  0,
			wantErr:      "maxKeepRuns must be at least 1, got 0",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewPipelineRunPrunerByStage(
				tc.tektonClient, tc.logger, tc.minKeepHours, tc.maxKeepRuns,
			)
			if err == nil || err.Error() != tc.wantErr {
				t.Fatalf("want err \"%s\", got: %s", tc.wantErr, err)
			}
		})
	}
}

func TestPrune(t *testing.T) {
	tclient := &tektonClient.TestClient{}
	minKeepHours := 2
	maxKeepRuns := 1
	logger := &logging.LeveledLogger{Level: logging.LevelDebug}
	p, err := NewPipelineRunPrunerByStage(tclient, logger, minKeepHours, maxKeepRuns)
	if err != nil {
		t.Fatal(err)
	}
	prs := []tekton.PipelineRun{
		// not pruned
		pipelineRun("pr-a", "p-one", config.DevStage, time.Now().Add(time.Minute*-1)),
		// would be pruned by maxKeepRuns, but is protected by minKeepHours
		pipelineRun("pr-b", "p-one", config.DevStage, time.Now().Add(time.Minute*-3)),
		// pruned
		pipelineRun("pr-c", "p-one", config.DevStage, time.Now().Add(time.Hour*-4)),
		// pruned through pipeline p-two
		pipelineRun("pr-d", "p-two", config.DevStage, time.Now().Add(time.Hour*-5)),
		// pruned through pipeline p-two
		pipelineRun("pr-e", "p-two", config.DevStage, time.Now().Add(time.Hour*-6)),
		// not pruned because different stage (QA)
		pipelineRun("pr-e", "p-three", config.QAStage, time.Now()),
		// not pruned because different stage (PROD)
		pipelineRun("pr-f", "p-four", config.ProdStage, time.Now().Add(time.Hour*-7)),
		// pruned
		pipelineRun("pr-g", "p-four", config.ProdStage, time.Now().Add(time.Hour*-8)),
	}
	err = p.Prune(context.TODO(), prs)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(tclient.DeletedPipelineRuns)
	if diff := cmp.Diff([]string{"pr-c", "pr-g"}, tclient.DeletedPipelineRuns); diff != "" {
		t.Fatalf("pr prune mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"p-two"}, tclient.DeletedPipelines); diff != "" {
		t.Fatalf("p prune mismatch (-want +got):\n%s", diff)
	}
}

func pipelineRun(name, pipeline string, stage config.Stage, creationTime time.Time) tekton.PipelineRun {
	return tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			CreationTimestamp: metav1.Time{Time: creationTime},
			Labels: map[string]string{
				stageLabel:          string(stage),
				tektonPipelineLabel: pipeline,
			},
		},
	}
}
