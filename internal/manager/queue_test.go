package manager

import (
	"testing"
	"time"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// fakeAdvancerDone is always done advancing the queue.
type fakeAdvancerDone struct {
	Logger logging.LeveledLoggerInterface
}

func (f *fakeAdvancerDone) AdvanceQueue(repository string) (int, error) {
	return 0, nil
}

func TestPollIdentifier(t *testing.T) {
	p := &pipelineRunQueue{
		queues: map[string]bool{
			"a": true,
			"b": false,
		},
		pollInterval: time.Second,
		logger:       &logging.LeveledLogger{Level: logging.LevelNull},
	}
	f := &fakeAdvancerDone{
		Logger: &logging.LeveledLogger{Level: logging.LevelNull},
	}
	p.StartPolling(f, "a", time.Second)
	p.StartPolling(f, "b", time.Second)
	if !p.queues["a"] {
		t.Fatal("polling state for 'a' should be true")
	}
	if !p.queues["b"] {
		t.Fatal("polling state for 'b' should be true")
	}
}

// fakeAdvancerSteps can be called a few times before it is done advancing the queue.
type fakeAdvancerSteps struct {
	count  int
	Logger logging.LeveledLoggerInterface
}

func (f *fakeAdvancerSteps) AdvanceQueue(repository string) (int, error) {
	if f.count < 2 {
		f.count++
		return 1, nil
	}
	return 0, nil
}

func TestAdvanceQueueAndQuit(t *testing.T) {
	p := &pipelineRunQueue{
		queues:       map[string]bool{},
		pollInterval: time.Millisecond,
		logger:       &logging.LeveledLogger{Level: logging.LevelNull},
	}
	f := &fakeAdvancerSteps{
		Logger: &logging.LeveledLogger{Level: logging.LevelNull},
	}
	done := p.StartPolling(f, "a", time.Second)
	select {
	case <-done:
		t.Log("quit occured")
	case <-time.After(5 * time.Second):
		t.Fatal("quit should have occured")
	}
}

func TestAdvanceQueue(t *testing.T) {
	tests := map[string]struct {
		runs         []*tekton.PipelineRun
		wantStart    string
		wantPollDone bool
	}{
		"none": {
			runs:         []*tekton.PipelineRun{},
			wantStart:    "",
			wantPollDone: true,
		},
		"one cancelled, none pending": {
			runs: []*tekton.PipelineRun{
				cancelledPipelineRun(t, "one", time.Now()),
			},
			wantStart:    "",
			wantPollDone: true,
		},
		"one cancelled, one pending": {
			runs: []*tekton.PipelineRun{
				cancelledPipelineRun(t, "one", time.Now()),
				pendingPipelineRun(t, "two", time.Now()),
			},
			wantStart:    "two",
			wantPollDone: true,
		},
		"one cancelled, two pending": {
			runs: []*tekton.PipelineRun{
				cancelledPipelineRun(t, "one", time.Now()),
				pendingPipelineRun(t, "two", time.Now().Add(time.Minute*-1)),
				pendingPipelineRun(t, "three", time.Now().Add(time.Minute*-2)),
			},
			wantStart:    "three",
			wantPollDone: false,
		},
		"two pending": {
			runs: []*tekton.PipelineRun{
				pendingPipelineRun(t, "one", time.Now().Add(time.Minute*-2)),
				pendingPipelineRun(t, "two", time.Now().Add(time.Minute*-1)),
			},
			wantStart:    "one",
			wantPollDone: false,
		},
		"one timed out, one pending": {
			runs: []*tekton.PipelineRun{
				timedOutPipelineRun(t, "one", time.Now().Add(time.Minute*-2)),
				pendingPipelineRun(t, "two", time.Now().Add(time.Minute*-1)),
			},
			wantStart:    "two",
			wantPollDone: true,
		},
		"one running, one pending": {
			runs: []*tekton.PipelineRun{
				runningPipelineRun(t, "one", time.Now().Add(time.Minute*-2)),
				pendingPipelineRun(t, "two", time.Now().Add(time.Minute*-1)),
			},
			wantStart:    "",
			wantPollDone: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tclient := &tektonClient.TestClient{PipelineRuns: tc.runs}
			s := &Server{TektonClient: tclient, Logger: &logging.LeveledLogger{Level: logging.LevelNull}}
			queueLength, err := s.AdvanceQueue("a")
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantStart != "" {
				if len(tclient.UpdatedPipelineRuns) != 1 {
					t.Fatal("should have updated one run")
				}
				if tclient.UpdatedPipelineRuns[0] != tc.wantStart {
					t.Fatalf("should have updated run '%s'", tc.wantStart)
				}
			} else {
				if len(tclient.UpdatedPipelineRuns) > 0 {
					t.Fatal("should not have updated any run")
				}
			}
			if (queueLength == 0) != tc.wantPollDone {
				t.Fatalf("want polling to be done: %v, but queue length is: %d", tc.wantPollDone, queueLength)
			}
		})
	}
}

func pendingPipelineRun(t *testing.T, name string, creationTime time.Time) *tekton.PipelineRun {
	pr := &tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			CreationTimestamp: metav1.Time{Time: creationTime},
		},
		Spec: tekton.PipelineRunSpec{
			Status: tekton.PipelineRunSpecStatusPending,
		},
	}
	if !pr.IsPending() {
		t.Fatal("pr should be pending")
	}
	return pr
}

func cancelledPipelineRun(t *testing.T, name string, creationTime time.Time) *tekton.PipelineRun {
	pr := &tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			CreationTimestamp: metav1.Time{Time: creationTime},
		},
		Spec: tekton.PipelineRunSpec{
			Status: tekton.PipelineRunSpecStatusCancelled,
		},
	}
	if !pr.IsCancelled() || pr.IsPending() || pr.IsDone() || pr.IsTimedOut() {
		t.Fatal("pr should be cancelled")
	}
	return pr
}

func timedOutPipelineRun(t *testing.T, name string, creationTime time.Time) *tekton.PipelineRun {
	// pipelineTimeout := pr.Spec.Timeout
	// startTime := pr.Status.StartTime
	pr := &tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			CreationTimestamp: metav1.Time{Time: creationTime},
		},
		Spec: tekton.PipelineRunSpec{
			Timeout: &metav1.Duration{Duration: time.Second},
		},
		Status: tekton.PipelineRunStatus{
			PipelineRunStatusFields: tekton.PipelineRunStatusFields{
				StartTime: &metav1.Time{Time: time.Now().Add(-2 * time.Second)},
			},
		},
	}
	if !pr.IsTimedOut() || pr.IsPending() || pr.IsDone() || pr.IsCancelled() {
		t.Fatal("pr should be timed out")
	}
	return pr
}

func runningPipelineRun(t *testing.T, name string, creationTime time.Time) *tekton.PipelineRun {
	pr := &tekton.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			CreationTimestamp: metav1.Time{Time: creationTime},
		},
	}
	if pr.IsDone() || pr.IsPending() || pr.IsTimedOut() || pr.IsCancelled() {
		t.Fatal("pr should be running")
	}
	return pr
}
