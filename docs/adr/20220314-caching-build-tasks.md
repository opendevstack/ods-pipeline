# Build Task Caching

Date: 2022-03-14

Last Modified: 2023-08-30 - build task caching to be based on multiple input locations  

## Status

Accepted

## Context

For a repo with a single build, the pipeline build will only be triggered if any changes trigger a new pipeline run.

However if one has multiple builds in one repo, the pipeline build runs all build tasks even for those areas where nothing has changed.

A simple scenario is a repo with a backend and frontend component whose build outputs are combined by a subsequent docker build in a subsequent deploy task.

There is no good reason when only the frontend code changes that the backend build is redone as well and vice versa of course.

## Decision

Build task caching persists a build's produced reports and built files (if any) so that future builds can be skipped if the specified builds inputs did not change.

The built output files (if any) and reports (ods artifacts) are cached in a dedicated cache area named `build-task` on the caching PVC described in the [caching ADR](20220225-caching.md).

In most cases as ods-pipeline builds should enable build task caching by specifying an appropriate value in their `cache-sources` parameter (colon separated paths). This parameter must list any directories which impact the build. When there are no files impacting the build outside of the `working-dir` one can set the a tasks `cache-sources` to the same value. By default this value is empty which disables build caching.

Build tasks with caching enabled shall copy their build outputs and ods artifacts to `.ods-cache/build-task/<technology-build-key>-<cache-build-key>/$(git-sha-combined)`, where

- `<technology-build-key>` is a value to allow humans to easily distinguish for what the artifacts are, while also allowing to have different keys for different versions of a build task. This might be useful for example to separate jar files with different class file versions for example. Such differentiation could be used in later versions. Initially they will likely be simple names such as `go`, `python` etc.

- `<cache-build-key>` allows keeping multiple build tasks associated with the same working directory separate. Such build variants may allow to create builds for different platforms or architectures while keeping the cached files in separate locations. In most cases such a key can be derived from the existing build task's parameters.

- `git-sha-combined` is the git sha of the internal tree object of the specified `cache-sources` directories (specified colon separated)

<aside class="notice">
****
`git-sha-working-dir=$(git rev-parse "HEAD:$working_dir")`

This is the git cryptographic hash **(sha) of the internal Git tree object** which stores the working-dir (and all the git managed files and directories beneath it).
In particular this is **not a git commit hash** which for example last modified that directory.
For more information see

- https://git-scm.com/book/en/v2/Git-Internals-Git-Objects
- https://stackoverflow.com/a/55459589
- https://stackoverflow.com/a/58600111
****
</aside>

For traceability the tekton build tasks supporting build skipping introduce a result named `build-reused-from-location`. The file contains a path to the cache folder used or is empty if no build was reused.

The build tasks support proper cleanup by touching file named `.ods-last-used-stamp` when a cache location is used or created.

### Cleanup

`ods-start` is changed to delete prior build output from the build cache if it is older than the number of days configured in task parameter `cache-build-tasks-for-days`:

- `cache-build-tasks-for-days` an integer number provided as string, which defines the number of days build outputs are cached.  Defaulting to `'7'`. A negative number can be used to clear the build task cache.

## Consequences

* With build skipping build tasks for which their working directory does not change can skip rebuilding and instead reuse their prior build outputs (if any) and other produced artifacts (see artifacts.adoc).

* For build skipping to work as expected a build task must be have specify all directories which impact the builds in parameter `cache-sources`. In particular a build must not use any files of the repo outside these paths. You can disable build caching by setting `cache-sources` to an empty value which is the default.

* The build skipping cache directories are cleaned up by `ods-start` if they have not been used for 7 days or what you configure with the `ods-start`'s  parameter `cache-build-tasks-for-days`.

* It is not readily apparent why the build step will be skipped while the sonarqube step is redone. The latter ensures that the quality gate attribution occurs for example when a feature branch is merged.

## Alternatives

As an alternative it was considered to use Nexus to store build outputs. The reasons this was not further pursued was:

* There would be a lot of additional pressure on Nexus as the build outputs can be rather large.
* There were doubts that the performance would be satisfactory.
* The complexity of the proposal appeared concerning.
