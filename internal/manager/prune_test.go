package manager

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	tektonClient "github.com/opendevstack/ods-pipeline/internal/tekton"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type fakePruner struct {
	called chan string
}

func (p *fakePruner) prune(ctx context.Context, repository string) error {
	p.called <- repository
	return nil
}

func TestRun(t *testing.T) {
	repoCh := make(chan string)
	p := &Pruner{
		TriggeredRepos: repoCh,
		Logger:         &logging.LeveledLogger{Level: logging.LevelNull},
		MinKeepHours:   2,
		MaxKeepRuns:    1,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pruneCh := make(chan string)
	fp := &fakePruner{called: pruneCh}
	delay := 50 * time.Millisecond
	go p.run(ctx, fp, delay)
	// Trigger the same repository twice.
	repoCh <- "a"
	repoCh <- "a"
	// Expect that pruning runs after delay.
	select {
	case <-pruneCh:
		t.Log("pruning triggered")
	case <-time.After(delay + time.Second):
		t.Fatal("pruning should have been triggered")
	}
	// Expect that pruning didn't run twice.
	select {
	case <-pruneCh:
		t.Fatal("pruning should not have been triggered twice")
	case <-time.After(time.Millisecond):
		t.Log("pruning was not triggered twice")
	}
	// // Trigger the same repository again.
	repoCh <- "a"
	// Expect that pruning runs after delay.
	select {
	case <-pruneCh:
		t.Log("pruning triggered again")
	case <-time.After(delay + time.Second):
		t.Fatal("pruning should have been triggered again")
	}
}

func TestPrune(t *testing.T) {
	tclient := &tektonClient.TestClient{
		PipelineRuns: []*tekton.PipelineRun{
			// not pruned
			pipelineRun("pr-a", time.Now().Add(time.Minute*-1)),
			// would be pruned by maxKeepRuns, but is protected by minKeepHours
			pipelineRun("pr-b", time.Now().Add(time.Minute*-3)),
			// pruned
			pipelineRun("pr-c", time.Now().Add(time.Hour*-4)),
			// pruned
			pipelineRun("pr-d", time.Now().Add(time.Hour*-5)),
		},
	}
	minKeepHours := 2
	maxKeepRuns := 1
	logger := &logging.LeveledLogger{Level: logging.LevelNull}
	ch := make(chan string)
	p := &Pruner{
		TriggeredRepos: ch,
		TektonClient:   tclient,
		Logger:         logger,
		MinKeepHours:   minKeepHours,
		MaxKeepRuns:    maxKeepRuns,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := p.prune(ctx, "repo")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(tclient.DeletedPipelineRuns)
	if diff := cmp.Diff([]string{"pr-c", "pr-d"}, tclient.DeletedPipelineRuns); diff != "" {
		t.Fatalf("pipeline run prune mismatch (-want +got):\n%s", diff)
	}
}

func pipelineRun(name string, creationTime time.Time) *tekton.PipelineRun {
	return &tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			CreationTimestamp: metav1.Time{Time: creationTime},
		},
	}
}
