package manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Label set by Tekton identifying the pipeline of a run
	tektonPipelineLabel = "tekton.dev/pipeline"
)

// PipelineRunPruner is the interface for a pruner implementation.
type PipelineRunPruner interface {
	// Prune removes pipeline runs (and potentially pipelines) from the
	// given list of pipeline runs based on a strategy as decided by the
	// implemnter.
	Prune(ctxt context.Context, pipelineRuns []tekton.PipelineRun) error
}

// PipelineRunPrunerByStage prunes pipeline runs by target stage.
// It's behaviour can be controlled through minKeepHours and maxKeepRuns.
// When pruning, it keeps maxKeepRuns number of pipeline runs per stage,
// however it always keeps all pipelines less than minKeepHours old.
// If pruning would prune all pipeline runs of one pipeline (identified by
// label "tekton.dev/pipeline"), then the pipeline is pruned instead (removing
// all dependent pipeline runs through propagation).
type pipelineRunPrunerByStage struct {
	client tektonClient.ClientInterface
	logger logging.LeveledLoggerInterface
	// minKeepHours specifies the minimum hours to keep a pipeline run.
	// This setting has precendence over maxKeepRuns.
	minKeepHours int
	// maxKeepRuns is the maximum number of pipeline runs to keep per stage.
	maxKeepRuns int
}

type prunableResources struct {
	pipelineRuns []string
	pipelines    []string
}

// NewPipelineRunPrunerByStage returns an instance of pipelineRunPrunerByStage.
func NewPipelineRunPrunerByStage(
	client tektonClient.ClientInterface,
	logger logging.LeveledLoggerInterface,
	minKeepHours, maxKeepRuns int) (*pipelineRunPrunerByStage, error) {
	if client == nil {
		return nil, errors.New("tekton client is required")
	}
	if logger == nil {
		return nil, errors.New("logger is required")
	}
	if minKeepHours < 1 {
		return nil, fmt.Errorf("minKeepHours must be at least 1, got %d", minKeepHours)
	}
	if maxKeepRuns < 1 {
		return nil, fmt.Errorf("maxKeepRuns must be at least 1, got %d", maxKeepRuns)
	}
	return &pipelineRunPrunerByStage{
		client:       client,
		logger:       logger,
		minKeepHours: minKeepHours,
		maxKeepRuns:  maxKeepRuns,
	}, nil
}

// Prune prunes runs within pipelineRuns which can be cleaned up according to
// the strategy in pipelineRunPrunerByStage.
func (p *pipelineRunPrunerByStage) Prune(ctxt context.Context, pipelineRuns []tekton.PipelineRun) error {
	p.logger.Debugf("Prune settings: minKeepHours=%d maxKeepRuns=%d", p.minKeepHours, p.maxKeepRuns)
	prByStage := p.categorizePipelineRunsByStage(pipelineRuns)
	for stage, prs := range prByStage {
		p.logger.Debugf("Calculating prunable pipelines / pipeline runs for stage %s ...", stage)
		prunable := p.findPrunableResources(prs)

		p.logger.Debugf("Pruning %d \"%s\" stage pipelines and their dependent runs ...", len(prunable.pipelines), stage)
		for _, name := range prunable.pipelines {
			err := p.prunePipeline(ctxt, name)
			if err != nil {
				p.logger.Warnf("Failed to prune pipeline %s: %s", name, err)
			}
		}

		p.logger.Debugf("Pruning %d \"%s\" stage pipeline runs ...", len(prunable.pipelineRuns), stage)
		for _, name := range prunable.pipelineRuns {
			err := p.pruneRun(ctxt, name)
			if err != nil {
				p.logger.Warnf("Failed to prune pipeline run %s: %s", name, err)
			}
		}
	}
	return nil
}

// categorizePipelineRunsByStage assigns the given pipelineRuns into buckets
// by target stages (DEV, QA, PROD).
func (p *pipelineRunPrunerByStage) categorizePipelineRunsByStage(pipelineRuns []tekton.PipelineRun) map[string][]tekton.PipelineRun {
	pipelineRunsByStage := map[string][]tekton.PipelineRun{
		string(config.DevStage):  {},
		string(config.QAStage):   {},
		string(config.ProdStage): {},
	}
	for _, pr := range pipelineRuns {
		stage := pr.Labels[stageLabel]
		if _, ok := pipelineRunsByStage[stage]; !ok {
			p.logger.Warnf("Unknown stage '%s' for pipeline run %s", stage, pr.Name)
		}
		pipelineRunsByStage[stage] = append(pipelineRunsByStage[stage], pr)
	}
	return pipelineRunsByStage
}

// findPrunableResources finds resources that can be pruned within the given
// pipeline runs. Returned resources are either pipelines or pipeline runs.
// If all pipeline runs of one pipeline can be pruned, the pipeline is
// returned instead of the individual pipeline runs.
func (s *pipelineRunPrunerByStage) findPrunableResources(pipelineRuns []tekton.PipelineRun) *prunableResources {
	sortPipelineRunsDescending(pipelineRuns)

	// Apply cleanup to each bucket.
	prunablePipelines := []string{}
	prunablePipelineRuns := []string{}

	cutoff := time.Now().Add(time.Duration(s.minKeepHours*-1) * time.Hour)
	protectedRuns := []tekton.PipelineRun{}
	prunableRuns := []tekton.PipelineRun{}
	// Categorize runs as either "protected" or "prunable".
	// A run is protected if it is newer than the cutoff time, or if maxKeepRuns
	// is not reached yet.
	for _, p := range pipelineRuns {
		if p.CreationTimestamp.Time.After(cutoff) || len(protectedRuns) < s.maxKeepRuns {
			protectedRuns = append(protectedRuns, p)
		} else {
			prunableRuns = append(prunableRuns, p)
		}
	}
	// Check for each prunable run, if there is another run for the same pipeline
	// which is protected. If no such run exists, we want to prune the pipeline
	// as a whole instead of individual runs.
	for _, pruneableRun := range prunableRuns {
		if pipelineIsProtected(pruneableRun.Labels[tektonPipelineLabel], protectedRuns) {
			prunablePipelineRuns = append(prunablePipelineRuns, pruneableRun.Name)
		} else {
			prunablePipelines = append(prunablePipelines, pruneableRun.Labels[tektonPipelineLabel])
		}
	}

	return &prunableResources{
		pipelineRuns: unique(prunablePipelineRuns),
		pipelines:    unique(prunablePipelines),
	}
}

// pipelineIsProtected checks if the pipelineName exists in the given pipeline
// runs.
func pipelineIsProtected(pipelineName string, protected []tekton.PipelineRun) bool {
	for _, protect := range protected {
		if protect.Labels[tektonPipelineLabel] == pipelineName {
			return true
		}
	}
	return false
}

// pruneRun removes the pipeline run identified by name. The deletion is
// propagated to dependents.
func (p *pipelineRunPrunerByStage) pruneRun(ctxt context.Context, name string) error {
	p.logger.Debugf("Pruning pipeline run %s ...", name)
	ppPolicy := v1.DeletePropagationForeground
	return p.client.DeletePipelineRun(
		ctxt,
		name,
		v1.DeleteOptions{PropagationPolicy: &ppPolicy},
	)
}

// prunePipeline removes the pipeline identified by name. The deletion is
// propagated to dependents.
func (p *pipelineRunPrunerByStage) prunePipeline(ctxt context.Context, name string) error {
	p.logger.Debugf("Pruning pipeline %s and its dependent runs ...", name)
	ppPolicy := v1.DeletePropagationForeground
	return p.client.DeletePipeline(
		ctxt,
		name,
		v1.DeleteOptions{PropagationPolicy: &ppPolicy},
	)
}

// unique returns a slice of strings where all items appear only once.
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
