# Trigger-based parameters

Date: 2023-04-25

## Status

Accepted

## Context

The fields `branchToEnvironmentMapping` and `environments` allow for a little bit of dynamic configuration of pipelines. In particular, they allow to set different environments to deploy to based on the checked out branch, and they allow to control behaviour around artifacts storage and promotion protection via the stage concept.

However, there is no way to generally parameterise pipelines based on checked out Git refs. Further, the `environments` concept is targeted towards K8s/Helm and may not be appropriate for other environments. Finally, the stge concept is also quite limited and does introduce some magic that is not immediately obvious.

## Decision

Extend the trigger mechansim to allow dynamic parameterisation of pipelines.

## Consequences

* Triggers will allow to specify pipeline parameters and task parameters
* Multiple triggers per pipeline should be possible to have different parameters per trigger, while using the same pipeline spec
* No more `branchToEnvironmentMapping` and `environments fields. The functionality they provided must now be availble through other means, e.g. through task parameters
* The `ods-deploy-helm` task needs to learn many new parameters that allow to configure what was previously possible through `environments`.
* The stage concept needs to be replaced. The Git tag support will just be dropped. Its main purpose was to ensure promotion to prod is impossible without promotion to QA first. A much better approach will be to allow users to specify an "artifact source" and and "artifact target". Artifacts for a pipeline will be pulled from the "artifact source" (a Nexus repo) and pushed to the "artifact target" (also a Nexus repo). Using this mechanism, promotion protection can be implemented using different Nexus repos for DEV/QA/PROD. The mechanism has more benefits, such as being able to not retrain any artifacts for some branches.

## Rejected Alternatives

I went through a few iterations how exactly to configure triggers and parameters, but in the end the chosen solution should be the most flexible one.
