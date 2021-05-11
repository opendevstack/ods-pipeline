# Put version into task name

## Rationale

Tekton has no task versioning concept built in. `ClusterTask` resources provided by the Tekton catalog have the version in the name, e.g. `buildah-v0-16-3` and `buildah-v0-19-0`. For lack of a better alternative, we use the same approach. Every ODS release will install new tasks with new names.

## Context

The only helpful related thing we see is Tekton Bundles, see e.g. https://github.com/tektoncd/pipeline/issues/1839. However:

* right now it's only behind a feature flag
* it stores tasks in a container images which sound odd
* it prevents users from inspecting the tasks in the OpenShift UI

