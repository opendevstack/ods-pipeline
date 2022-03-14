# Build Task Caching

Date: 2022-03-14

## Status

Accepted

## Context

For a repo with a single build, the pipeline build will only be triggered if any changes trigger a new pipeline run.

However if one has multiple builds in one repo, the pipeline build runs all build tasks even for those areas where nothing has changed.

A simple scenario is a repo with a backend and frontend component whose build outputs are combined by a subsequent docker build in a subsequent deploy task.

There is no good reason when only the frontend code changes that the backend build is redone as well and vice versa of course.

## Decision

Build task caching persists a build's produced files so that future builds can be skipped if the working directory specified in the build task did not change.

The build files (output and reports) are cached in a dedicated cash area named `build-task` on the caching PVC described in the [caching ADR](20220225-caching.md).

By default build task caching is not enabled for the following reasons:

- Build task skipping has limited benefits for single build repos which is the main design focus of `ods-pipeline`. The build task cache could only be reused when the working dir did not change while some other change in the repo triggers a build. This might be the case when the working-dir is not at the root or when there is a merge commit without file changes.
- Build caching requires space on the caching PVC and a (typically) small runtime overhead on initial cache population.

In repos with multiple build tasks you would typically want to enable build task skipping. However when enabled builds are skipped if the contents of their `working-dir` does not change. 

**NOTE: Do not enable if your build consumes files in your repo located outside of the `working-dir` and changes of them are not reflected as version controlled file changes in or below the working dir!**
For example this could happen if you consume code maintained in your multi-repo in several build tasks. To enable build skipping in this case you will have to provide a checksum file of those other areas and check them into the working directory. This should be achievable via a git pre-commit hook.  

Tasks supporting builds skipping will support a parameter `cache-build` which by default is `'false'`. If `cache-build` is `'true'` build tasks shall copy their build outputs and reports to `.ods-cache/build-task/<technology-build-key>-<cache-build-key>/$(git-sha-working-dir)`, where

- `<technology-build-key>` is a value to allow humans to easily distinguish for what the artifacts are, while also allowing to have different keys for different versions of a build task. This might be useful for example to separate jar files with different class file versions for example. Such differentiation could be used in later versions. Initially they will likely be simple names such as `go`, `python` etc.

- `<cache-build-key>` allows keeping multiple build tasks associated with the same working directory separate. Such build variants may allow to create builds for different platforms or architectures while keeping the cached files in separate locations. In most cases such a key can be derived from the existing build task's parameters. 

- `git-sha-working-dir` is the git sha of the internal tree object of the working directory

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

* You can enable build skipping so that build tasks for which their working directory does not change can skip rebuilding and instead reuse their prior build outputs and reports.

* For build skipping to work as expected a build task must be self contained by their working directory. In particular a build must not use any files of the repo outside the `working-dir`. See also NOTE of section Decision above for more examples and potential workarounds.

* It is recommended to enable this in multi build repos unless the build is not self contained.

* The build skipping cache directories are cleaned up by `ods-start` if they have not been used for 7 days or what you configure with the `ods-start`'s  parameter `cache-build-tasks-for-days`.

* It is not readily apparent why the build step will be skipped while the sonarqube step is redone. The latter ensures that the quality gate attribution occurs for example when a feature branch is merged.

## Alternatives

As an alternative it was considered to use Nexus to store build outputs. The reasons this was not further pursued was:

* There would be a lot of additional pressure on Nexus as the build outputs can be rather large.
* There were doubts that the performance would be satisfactory.
* The complexity of the proposal appeared concerning.

