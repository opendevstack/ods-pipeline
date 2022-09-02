# Use PipelineRun resources with inlined spec

Date: 2022-07-29

## Status

Accepted

## Context

The current design starts pipeline runs referencing a pipeline. Users of ODS pipeline can control the `Pipeline` resource, but cannot control the `PipelineRun` as it is created automatically. This limits the user to features available at Tekton "authoring time" (see ["authoring time" (Task/Pipeline) vs "runtime" (TaskRun/PipelineRun)](https://github.com/tektoncd/community/blob/main/design-principles.md#reusability)). For example, in contrast to `Pipeline` resources, `PipelineRun` resources allow:

* [specifying pipeline specs and task specs inline](https://github.com/tektoncd/pipeline/blob/release-v0.38.x/examples/v1beta1/pipelineruns/pipelinerun-with-pipelinespec-and-taskspec.yaml), enabling ad-hoc tasks in `ods.yaml` (which soon can make use of [propagated parameters](https://tekton.dev/docs/pipelines/pipelineruns/#propagated-parameters))
* [overriding task steps and sidecars](https://tekton.dev/docs/pipelines/taskruns/#overriding-task-steps-and-sidecars) and [specifying task-level compute resources](https://tekton.dev/docs/pipelines/pipelineruns/#specifying-task-level-computeresources) (soon)
* [specifying pod templates](https://tekton.dev/docs/pipelines/pipelineruns/#specifying-a-pod-template)
* [configuring of timeouts](https://tekton.dev/docs/pipelines/pipelineruns/#configuring-a-failure-timeout)

## Decision

Stop creating and updating a `Pipeline` resource per branch, and instead create a `PipelineRun` resource with an inline spec for each pipeline run.

## Design

Since `Pipeline#spec` and `PipelineRun#pipelineSpec` are the same `PipelineSpec` type, the configuration of `tasks` and `finally` does not need to change. The new `pipeline` definition looks like this:

```
pipeline:
  tasks: []PipelineTask
  finally: []PipelineTask
  timeouts: *TimeoutFields
  podTemplate: *PodTemplate
  taskRunSpecs: []PipelineTaskRunSpec
```

As an alternative, the following `pipeline` definition was considered:

```
pipeline:
  pipelineSpec:
    tasks: []PipelineTask
    finally: []PipelineTask
  timeouts: *TimeoutFields
  podTemplate: *PodTemplate
  taskRunSpecs: []PipelineTaskRunSpec
```

This would be closer to native Tekton but it is more verbose and not backwards-compatible. Therefore, this alternative was rejected.

## Consequences

The OpenShift console UI allows to navigate from `Pipeline` to `PipelineRun`. If we do not have `Pipeline` resources any longer, this grouping mechanism is gone. This drawback is acceptable because we assume most users navigate directly from Bitbucket builds to pipeline runs. Further, the OpenShift UI allows to filter pipeline runs by label, which can be used to see all runs belonging to a branch.

The change should simplify quite a bit of the internal logic, as management of pipelines (creating, updating, pruning) can be removed.

Most importantly, the new design provides access to Tekton features available only via `PipelineRun` like inline task specs, compute resource overrides and control over timeouts and pod templates.
