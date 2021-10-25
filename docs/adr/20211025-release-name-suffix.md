# Add version suffix to tasks Helm release

Date: 2021-10-25

## Status

Accepted

## Context

Central tasks are installed as part of the `ods-pipeline-tasks` Helm release. If a user upgrades the release, then the current tasks (such as `ods-build-go-v0-1-0`) would be be deleted as the new tasks have a different name (such as `ods-build-go-v0-2-0`). This is not desired: we want the old tasks to stick around until decommissioned by the cluster admin.

## Decision

The current version will be suffixed to the Helm release name (e.g. `ods-pipeline-tasks-v0-1-0` instead of `ods-pipeline-tasks`).

This version suffix is not needed for the images Helm release, nor the "local" installation as they simply reflect the latest state.

## Consequences

This will be a little odd as there will only ever be one installation per release, but it does allow for multiple task versions to co-exist in the cluster, as well as simple decommission of old tasks (by simple uninstalling the corresponding Helm release).
