package manager

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// pipelineRunQueue manages multiple queues. These queues
// can be polled in certain intervals.
type pipelineRunQueue struct {
	queues       map[string]bool
	pollInterval time.Duration
	// logger is the logger to send logging messages to.
	logger logging.LeveledLoggerInterface
}

// StartPolling periodically checks status for given identifier.
// The time until the first time is not more than maxInitialWait.
func (q *pipelineRunQueue) StartPolling(pt QueueAdvancer, identifier string, maxInitialWait time.Duration) chan bool {
	quit := make(chan bool)
	if q.queues[identifier] {
		close(quit)
		return quit
	}
	q.queues[identifier] = true

	wait(maxInitialWait)

	ticker := time.NewTicker(q.pollInterval)
	go func() {
		for {
			select {
			case <-quit:
				q.queues[identifier] = false
				ticker.Stop()
				return
			case <-ticker.C:
				q.logger.Debugf("Advancing queue for %s ...", identifier)
				queueLength, err := pt.AdvanceQueue(identifier)
				if err != nil {
					q.logger.Warnf("error during poll tick: %s", err)
				}
				if queueLength == 0 {
					q.logger.Debugf("Stopping to poll for %s ...", identifier)
					close(quit)
				}
			}
		}
	}()

	return quit
}

// QueueAdvancer is the interface passed to
// *pipelineRunQueue#StartPolling.
type QueueAdvancer interface {
	// AdvanceQueue is called for each poll step.
	AdvanceQueue(repository string) (int, error)
}

// Queue represents a pipeline run Queue. Pipelines of one repository must
// not run in parallel.
type Queue struct {
	TektonClient tektonClient.ClientPipelineRunInterface
}

// AdvanceQueue starts the oldest pending pipeline run if there is no
// progressing pipeline run at the moment.
// It returns the queue length.
func (s *Server) AdvanceQueue(repository string) (int, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	pipelineRuns, err := listPipelineRuns(s.TektonClient, ctxt, repository)
	if err != nil {
		return 0, fmt.Errorf("could not retrieve existing pipeline runs: %w", err)
	}
	s.Logger.Debugf("Found %d pipeline runs related to repository %s.", len(pipelineRuns.Items), repository)
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
	s.Logger.Debugf("Found runs for repo %s in state running=%v, pending=%d.", repository, foundRunning, len(pendingPrs))

	if !foundRunning && len(pendingPrs) > 0 {
		// update oldest pending PR
		sortPipelineRunsDescending(pendingPrs)
		oldestPR := pendingPrs[len(pendingPrs)-1]
		pendingPrs = pendingPrs[:len(pendingPrs)-1]
		s.Logger.Infof("Starting pending pipeline run %s ...", oldestPR.Name)
		oldestPR.Spec.Status = "" // remove pending status -> starts pipeline run
		_, err := s.TektonClient.UpdatePipelineRun(ctxt, &oldestPR, metav1.UpdateOptions{})
		if err != nil {
			return len(pendingPrs), fmt.Errorf("could not update pipeline run %s: %w", oldestPR.Name, err)
		}
	}
	return len(pendingPrs), nil
}

// needsQueueing checks if any run has either:
// - pending status set OR
// - is progressing
func needsQueueing(pipelineRuns *tekton.PipelineRunList) bool {
	for _, pr := range pipelineRuns.Items {
		if pr.Spec.Status == tekton.PipelineRunSpecStatusPending || pipelineRunIsProgressing(pr) {
			return true
		}
	}
	return false
}

// pipelineRunIsProgressing returns true if the PR is not done, not pending,
// not cancelled, and not timed out.
func pipelineRunIsProgressing(pr tekton.PipelineRun) bool {
	return !(pr.IsDone() || pr.IsPending() || pr.IsCancelled() || pr.IsTimedOut())
}

// wait waits for up to maxInitialWait. The exact wait time is
// pseudo-randomized if maxInitialWait is longer than one second.
func wait(maxInitialWait time.Duration) {
	initialWait := time.Second
	if maxInitialWait > time.Second {
		initialWait = time.Duration(rand.Intn(int(maxInitialWait.Seconds())-1) + 1)
	}
	timer := time.NewTimer(initialWait)
	<-timer.C
}
