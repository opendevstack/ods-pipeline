package manager

import (
	"context"
	"time"

	tektonClient "github.com/opendevstack/pipeline/internal/tekton"
	"github.com/opendevstack/pipeline/pkg/logging"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// pruneDelay defines when pruning kicks in after a pipeline was triggered.
	pruneDelay = time.Minute
	// pruneTimeout defines how long pruning is allowed to take before it gets cancelled.
	pruneTimeout = 5 * time.Minute
)

// Pruner prunes pipeline runs.
// It's behaviour can be controlled through MinKeepHours and MaxKeepRuns.
// When pruning, it keeps MaxKeepRuns number of pipeline runs,
// however it always keeps all pipeline runs less than MinKeepHours old.
type Pruner struct {
	// TriggeredRepos receives repo names for each triggered repository.
	TriggeredRepos chan string
	// TektonClient is a client to interact with Tekton.
	TektonClient tektonClient.ClientInterface
	Logger       logging.LeveledLoggerInterface
	// MinKeepHours specifies the minimum hours to keep a pipeline run.
	// This setting has precendence over MaxKeepRuns.
	MinKeepHours int
	// MaxKeepRuns is the maximum number of pipeline runs to keep.
	MaxKeepRuns   int
	upcomingPrune map[string]time.Time
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
	prunable := p.findPrunableResources(pipelineRuns.Items)

	p.Logger.Debugf("Pruning %d pipeline runs ...", len(prunable))
	for _, name := range prunable {
		err := p.pruneRun(ctxt, name)
		if err != nil {
			p.Logger.Warnf("Failed to prune pipeline run %s: %s", name, err)
		}
	}
	return nil
}

// findPrunableResources finds resources that can be pruned.
func (s *Pruner) findPrunableResources(pipelineRuns []tekton.PipelineRun) []string {
	sortPipelineRunsDescending(pipelineRuns)

	cutoff := time.Now().Add(time.Duration(s.MinKeepHours*-1) * time.Hour)
	protectedRuns := []string{}
	prunableRuns := []string{}
	// Categorize runs as either "protected" or "prunable".
	// A run is protected if it is newer than the cutoff time, or if MaxKeepRuns
	// is not reached yet.
	for _, p := range pipelineRuns {
		if p.CreationTimestamp.Time.After(cutoff) || len(protectedRuns) < s.MaxKeepRuns {
			protectedRuns = append(protectedRuns, p.Name)
		} else {
			prunableRuns = append(prunableRuns, p.Name)
		}
	}

	return unique(prunableRuns)
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
