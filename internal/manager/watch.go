package manager

import (
	"context"
	"fmt"
	"time"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// watchInterval defines in which interval queues are inspected.
	watchInterval = 30 * time.Second
	// advanceTimeout defines how long queue advancing is allowed to take before it gets cancelled.
	advanceTimeout = 5 * time.Minute
)

// Watcher watches pending pipeline run queues. Runs are queued per repository.
type Watcher struct {
	// PendingRunRepos receives repositories for which there is a new
	// pending pipeline run.
	PendingRunRepos chan string
	// Queues holds repositories for which there are pending runs.
	Queues       map[string]bool
	Logger       logging.LeveledLoggerInterface
	TektonClient tektonClient.ClientPipelineRunInterface
}

// Run starts monitoring and advancing the queues.
func (w *Watcher) Run(ctx context.Context) {
	t := time.NewTimer(watchInterval)
	for {
		select {
		case r := <-w.PendingRunRepos:
			w.Queues[r] = true
		case <-t.C:
			for _, q := range w.activeQueues() {
				w.Logger.Debugf("Advancing pipeline run queue for queue '%s' ...", q)
				remaining, err := w.advanceQueue(ctx, q)
				if err != nil {
					w.Logger.Errorf("could not advance queue '%s': %s", q, err)
				}
				if remaining == 0 {
					w.Queues[q] = false
				}
			}
		case <-ctx.Done():
			t.Stop()
			return
		}
		if w.hasActiveQueue() {
			t = time.NewTimer(watchInterval) // use Reset instead?
		} else {
			t.Stop()
		}
	}
}

// hasActiveQueue indicates whether any queue is active.
func (w *Watcher) hasActiveQueue() bool {
	for _, a := range w.Queues {
		if a {
			return true
		}
	}
	return false
}

// activeQueues returns all queues currently active.
func (w *Watcher) activeQueues() []string {
	var qs []string
	for x, a := range w.Queues {
		if a {
			qs = append(qs, x)
		}
	}
	return qs
}

// advanceQueue starts the oldest pending pipeline run if there is no
// progressing pipeline run at the moment.
// It returns the queue length.
func (w *Watcher) advanceQueue(ctx context.Context, queue string) (int, error) {
	ctxt, cancel := context.WithTimeout(ctx, advanceTimeout)
	defer cancel()
	pipelineRuns, err := listPipelineRuns(ctxt, w.TektonClient, queue)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve existing pipeline runs: %w", err)
	}
	w.Logger.Debugf("Found %d pipeline runs related to repository %s.", len(pipelineRuns.Items), queue)
	if len(pipelineRuns.Items) == 0 {
		return 0, nil
	}

	var foundRunning bool
	pendingPrs := []tekton.PipelineRun{}
	for _, pr := range pipelineRuns.Items {
		if pr.IsPending() {
			pendingPrs = append(pendingPrs, pr)
			continue
		}
		if pipelineRunIsProgressing(pr) {
			foundRunning = true
			continue
		}
	}
	w.Logger.Debugf("Found runs for repo %s in state running=%v, pending=%d.", queue, foundRunning, len(pendingPrs))

	if !foundRunning && len(pendingPrs) > 0 {
		// update oldest pending PR
		sortPipelineRunsDescending(pendingPrs)
		oldestPR := pendingPrs[len(pendingPrs)-1]
		pendingPrs = pendingPrs[:len(pendingPrs)-1]
		w.Logger.Infof("Starting pending pipeline run %s ...", oldestPR.Name)
		oldestPR.Spec.Status = "" // remove pending status -> starts pipeline run
		_, err := w.TektonClient.UpdatePipelineRun(ctxt, &oldestPR, metav1.UpdateOptions{})
		if err != nil {
			return len(pendingPrs), fmt.Errorf("could not update pipeline run %s: %w", oldestPR.Name, err)
		}
	}
	return len(pendingPrs), nil
}
