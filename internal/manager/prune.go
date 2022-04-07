package manager

import (
	"context"
	"time"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/config"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// tektonPipelineLabel is set by Tekton identifying the pipeline of a run
	tektonPipelineLabel = "tekton.dev/pipeline"
	// pruneDelay defines when pruning kicks in after a pipeline was triggered.
	pruneDelay = time.Minute
	// pruneTimeout defines how long pruning is allowed to take before it gets cancelled.
	pruneTimeout = 5 * time.Minute
)

// Pruner prunes pipeline runs by target stage.
// It's behaviour can be controlled through MinKeepHours and MaxKeepRuns.
// When pruning, it keeps MaxKeepRuns number of pipeline runs per stage,
// however it always keeps all pipelines less than MinKeepHours old.
// If pruning would prune all pipeline runs of one pipeline (identified by
// label "tekton.dev/pipeline"), then the pipeline is pruned instead (removing
// all dependent pipeline runs through propagation).
type Pruner struct {
	// TriggeredRepos receives repo names for each triggered repository.
	TriggeredRepos chan string
	// TektonClient is a client to interact with Tekton.
	TektonClient tektonClient.ClientInterface
	Logger       logging.LeveledLoggerInterface
	// MinKeepHours specifies the minimum hours to keep a pipeline run.
	// This setting has precendence over MaxKeepRuns.
	MinKeepHours int
	// MaxKeepRuns is the maximum number of pipeline runs to keep per stage.
	MaxKeepRuns   int
	upcomingPrune map[string]time.Time
}

// prunableResources holds pipelines and runs that can be pruned.
type prunableResources struct {
	pipelineRuns []string
	pipelines    []string
}

// pruner is an interface to facilitate testing.
type pruner interface {
	prune(ctx context.Context, repository string) error
}

// Run starts the pruning process by calling run. This indirection facilitates
// testing.
func (p *Pruner) Run(ctx context.Context) {
	p.run(ctx, p, pruneDelay)
}

// run actually starts the pruning process.
func (p *Pruner) run(ctx context.Context, pr pruner, delay time.Duration) {
	p.upcomingPrune = make(map[string]time.Time)
	p.Logger.Debugf("Prune settings: MinKeepHours=%d MaxKeepRuns=%d", p.MinKeepHours, p.MaxKeepRuns)
	for {
		select {
		case repo := <-p.TriggeredRepos:
			// Avoid pruning if pruning for this repo is already planned.
			if t, ok := p.upcomingPrune[repo]; ok && t.After(time.Now()) {
				break
			}
			p.upcomingPrune[repo] = time.Now().Add(delay)
			time.AfterFunc(delay, func() {
				err := pr.prune(ctx, repo)
				if err != nil {
					p.Logger.Errorf(err.Error())
				}
			})
		case <-ctx.Done():
			return
		}
	}
}

// prune prunes runs within pipelineRuns which can be cleaned up according to
// the strategy in Pruner.
func (p *Pruner) prune(ctx context.Context, repository string) error {
	ctxt, cancel := context.WithTimeout(ctx, pruneTimeout)
	defer cancel()
	pipelineRuns, err := listPipelineRuns(ctxt, p.TektonClient, repository)
	if err != nil {
		return err
	}
	p.Logger.Debugf("Found %d pipeline runs related to repository %s.", len(pipelineRuns.Items), repository)
	prByStage := p.categorizePipelineRunsByStage(pipelineRuns.Items)
	for stage, prs := range prByStage {
		p.Logger.Debugf("Calculating prunable pipelines / pipeline runs for stage %s ...", stage)
		prunable := p.findPrunableResources(prs)

		p.Logger.Debugf("Pruning %d \"%s\" stage pipelines and their dependent runs ...", len(prunable.pipelines), stage)
		for _, name := range prunable.pipelines {
			err := p.prunePipeline(ctxt, name)
			if err != nil {
				p.Logger.Warnf("Failed to prune pipeline %s: %s", name, err)
			}
		}

		p.Logger.Debugf("Pruning %d \"%s\" stage pipeline runs ...", len(prunable.pipelineRuns), stage)
		for _, name := range prunable.pipelineRuns {
			err := p.pruneRun(ctxt, name)
			if err != nil {
				p.Logger.Warnf("Failed to prune pipeline run %s: %s", name, err)
			}
		}
	}
	return nil
}

// categorizePipelineRunsByStage assigns the given pipelineRuns into buckets
// by target stages (DEV, QA, PROD).
func (p *Pruner) categorizePipelineRunsByStage(pipelineRuns []tekton.PipelineRun) map[string][]tekton.PipelineRun {
	pipelineRunsByStage := map[string][]tekton.PipelineRun{
		string(config.DevStage):  {},
		string(config.QAStage):   {},
		string(config.ProdStage): {},
	}
	for _, pr := range pipelineRuns {
		stage := pr.Labels[stageLabel]
		if _, ok := pipelineRunsByStage[stage]; !ok {
			p.Logger.Warnf("Unknown stage '%s' for pipeline run %s", stage, pr.Name)
		}
		pipelineRunsByStage[stage] = append(pipelineRunsByStage[stage], pr)
	}
	return pipelineRunsByStage
}

// findPrunableResources finds resources that can be pruned within the given
// pipeline runs. Returned resources are either pipelines or pipeline runs.
// If all pipeline runs of one pipeline can be pruned, the pipeline is
// returned instead of the individual pipeline runs.
func (s *Pruner) findPrunableResources(pipelineRuns []tekton.PipelineRun) *prunableResources {
	sortPipelineRunsDescending(pipelineRuns)

	// Apply cleanup to each bucket.
	prunablePipelines := []string{}
	prunablePipelineRuns := []string{}

	cutoff := time.Now().Add(time.Duration(s.MinKeepHours*-1) * time.Hour)
	protectedRuns := []tekton.PipelineRun{}
	prunableRuns := []tekton.PipelineRun{}
	// Categorize runs as either "protected" or "prunable".
	// A run is protected if it is newer than the cutoff time, or if MaxKeepRuns
	// is not reached yet.
	for _, p := range pipelineRuns {
		if p.CreationTimestamp.Time.After(cutoff) || len(protectedRuns) < s.MaxKeepRuns {
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
func (p *Pruner) pruneRun(ctxt context.Context, name string) error {
	p.Logger.Debugf("Pruning pipeline run %s ...", name)
	ppPolicy := metav1.DeletePropagationForeground
	return p.TektonClient.DeletePipelineRun(
		ctxt,
		name,
		metav1.DeleteOptions{PropagationPolicy: &ppPolicy},
	)
}

// prunePipeline removes the pipeline identified by name. The deletion is
// propagated to dependents.
func (p *Pruner) prunePipeline(ctxt context.Context, name string) error {
	p.Logger.Debugf("Pruning pipeline %s and its dependent runs ...", name)
	ppPolicy := metav1.DeletePropagationForeground
	return p.TektonClient.DeletePipeline(
		ctxt,
		name,
		metav1.DeleteOptions{PropagationPolicy: &ppPolicy},
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
