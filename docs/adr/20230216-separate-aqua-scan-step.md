# Extrace Aqua scan to separate container

Date: 2023-02-16

## Status

Accepted

## Context

Aqua is baked into the `ods-package-image` task. This has the advantage of fewer tasks (= faster) and that Aqua is more likely to be executed. However, due to how the Aqua scanner is distributed (via basic auth protected endpoint), having Aqua in the `ods-package-image` is not trivial and leads to permission issues as the `ods-package-image` forces gid=1001, which is not allowed to write to the workspace.

## Decision

Create a separate container for Aqua scanning.

## Consequences

Pipelines using Aqua will be a bit slower.

## Rejected Alternatives

Extracing Aqua scan into a separate task. This would provide more flexibility, but it would also be even slower. Further, while exploring this option, a few issues were identified:

* the logic which image should be scanned is not trivial. One could opt to scan all images existing as image digests, but this would break if the task is used only in a promotion pipeline, as artifacts belonging to the subrepo are not uploaded
* the Bitbucket logic does not make sense in a promotion pipeline
