package manager

import (
	"context"
	"time"

	kubernetesClient "github.com/opendevstack/ods-pipeline/internal/kubernetes"
	tektonClient "github.com/opendevstack/ods-pipeline/internal/tekton"
	"github.com/opendevstack/ods-pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

type StorageConfig struct {
	Provisioner string
	ClassName   string
	Size        string
}

// Scheduler creates or updates pipelines based on PipelineConfig received from
// the TriggeredPipelines channel. It then schedules a pipeline run
// connected to the pipeline. If the run cannot start immediately because
// of another run, the new pipeline run is created in pending status.
type Scheduler struct {
	// Channel to read newly received runs from
	TriggeredPipelines chan PipelineConfig
	// Channel to send triggered repos on (signalling to start pruning)
	TriggeredRepos chan string
	// Channel to send pending runs on (singalling to start watching)
	PendingRunRepos  chan string
	TektonClient     tektonClient.ClientInterface
	KubernetesClient kubernetesClient.ClientInterface
	Logger           logging.LeveledLoggerInterface
	// TaskKind is the Tekton resource kind for tasks.
	// Either "ClusterTask" or "Task".
	TaskKind tekton.TaskKind
	// TaskSuffic is the suffix applied to tasks (version information).
	TaskSuffix string

	StorageConfig StorageConfig
}

// Run starts the scheduling process.
func (s *Scheduler) Run(ctx context.Context) {
	for {
		select {
		case pData := <-s.TriggeredPipelines:
			needQueueing := s.schedule(ctx, pData)
			if needQueueing {
				s.PendingRunRepos <- pData.Repository
			}
		case <-ctx.Done():
			return
		}
	}
}

// schedule turns a PipelineConfig into a pipeline (run).
func (s *Scheduler) schedule(ctx context.Context, pData PipelineConfig) bool {
	ctxt, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Create PVC if it does not exist yet
	err := s.createPVCIfRequired(ctxt, pData)
	if err != nil {
		s.Logger.Errorf(err.Error())
		return false
	}

	pipelineRuns, err := listPipelineRuns(ctxt, s.TektonClient, pData.Repository)
	if err != nil {
		s.Logger.Errorf(err.Error())
		return false
	}
	s.Logger.Debugf("Found %d pipeline runs related to repository %s.", len(pipelineRuns.Items), pData.Repository)
	needQueueing := needsQueueing(pipelineRuns)
	s.Logger.Debugf("Creating run for pipeline %s (queued=%v) ...", pData.Component, needQueueing)
	_, err = createPipelineRun(s.TektonClient, ctxt, pData, s.TaskKind, s.TaskSuffix, needQueueing)
	if err != nil {
		s.Logger.Errorf(err.Error())
		return false
	}
	// Trigger pruning
	s.TriggeredRepos <- pData.Repository
	return needQueueing
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
