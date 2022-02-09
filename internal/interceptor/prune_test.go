package interceptor

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
		{ // not pruned
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-a",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(time.Minute * -1)},
				Labels: map[string]string{
					stageLabel:          config.DevStage,
					tektonPipelineLabel: "p-one",
				},
			},
		},
		{ // would be pruned by maxKeepRuns, but is protected by minKeepHours
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-b",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(time.Minute * -3)},
				Labels: map[string]string{
					stageLabel:          config.DevStage,
					tektonPipelineLabel: "p-one",
				},
			},
		},
		{ // pruned
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-c",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(time.Hour * -4)},
				Labels: map[string]string{
					stageLabel:          config.DevStage,
					tektonPipelineLabel: "p-one",
				},
			},
		},
		{ // pruned through pipeline p-two
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-d",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(time.Hour * -5)},
				Labels: map[string]string{
					stageLabel:          config.DevStage,
					tektonPipelineLabel: "p-two",
				},
			},
		},
		{ // pruned through pipeline p-two
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-e",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(time.Hour * -6)},
				Labels: map[string]string{
					stageLabel:          config.DevStage,
					tektonPipelineLabel: "p-two",
				},
			},
		},
		{ // not pruned because different stage (QA)
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-e",
				CreationTimestamp: metav1.Time{Time: time.Now()},
				Labels: map[string]string{
					stageLabel:          config.QAStage,
					tektonPipelineLabel: "p-three",
				},
			},
		},
		{ // not pruned because different stage (PROD)
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-f",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(time.Hour * -7)},
				Labels: map[string]string{
					stageLabel:          config.ProdStage,
					tektonPipelineLabel: "p-four",
				},
			},
		},
		{ // pruned
			ObjectMeta: metav1.ObjectMeta{
				Name:              "pr-g",
				CreationTimestamp: metav1.Time{Time: time.Now().Add(time.Hour * -8)},
				Labels: map[string]string{
					stageLabel:          config.ProdStage,
					tektonPipelineLabel: "p-four",
				},
			},
		},
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
