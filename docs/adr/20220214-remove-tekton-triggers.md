# Remove Tekton Triggers

Date: 2022-02-14

## Status

Accepted

## Context

Tekton Triggers is used to trigger pipelines in response to Bitbucket webhooks. ODS Pipeline adds an interceptor in the Tekton Triggers chain to pull off "CI as code". This was done in an attempt to stay close to "official" usage of Tekton.

It turns out we gain very little from this and needlessly complicate the architecture. Further, pipeline queueing requires that we control the creation of PipelineRun resources, which prevents us from using Tekton Triggers.

## Decision

Remove Tekton Triggers and promote the webhook interceptor to a "pipeline manager". In addition to its current functionality, it needs to learn two things:

* validation of webhook
* creation of pipeline run

## Consequences

This will simplify the architecture quite a bit, removing several Tekton resources. Further, it also makes it easier to debug issues when pipelines are not triggered as only one deployment/pod needs to be checked.
