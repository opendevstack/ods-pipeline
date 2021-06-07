# Release Management

We like to have a way to manage releases of the application via Git. It should be possible to see in Git which target environments there are, and which versions are on which environment, as well as full traceability how the pipeline looked like that deployed a certain version of the application to a target environment.

The approach should somehow be compatible with "ODS classic", as in there should be a way to integrate this with Jira (triggering a release from there), as well as work without the Jira integration.

While the approach may sound initially like GitOps, it is not pure GitOps for several reasons. For some background of the current problems with GitOps, read https://codefresh.io/about-gitops/pains-gitops-1-0/. Apart from the issues mentioned, we are also in a situation where we have targets which are not Kubernetes, such as AWS, Salesforce, SAP. They are easier integrated into a pipeline-style approach.

## Approach

The above talks about an "application" but does not define what it is. In ODS, one repository typically corresponds to one application component only. In "ODS classic" the release manager repository ties multiple components together into one "application". At the same time, this repository also contains information about environments.

This proposal wants to split the repo into two:

* an "environment" / "release" repo, containing information about the environments ("what version is deployed where?")

* an "app" / "deployable" repository, containing the pipeline definition, the sub-repos and e.g. a Helm umbrella chart ("what was deployed and how?").

The "app" / "deployable" repo (and any referenced sub-repos / components) has a list of Git tags that can be deployed or have been already deployed in the the past. By default the "app" / "deployable" repo only needs one branch (typically `master`), however people can make use of additional (e.g. `release/`) branches if wanted/needed.

The "environment" / "release" repo has one branch per environment, for example an `environment/qa` branch on which it describes what is deployed in the QA environment. Promotion between environments simply means opening a pull request between the environment branches, such as from `environment/qa` to `environment/prod`. The content of the repository is basically just one file, which contains:

* description of the environments (e.g. API and credentials to use)
* a version (such as `1.0.0`)
* the name of the "app" / "deployable" repo

When a push is made to a branch in the "environment" / "release" repo, a pipeline starts which deploys the "app" / "deployable" repo in the desired version in the corresponding environment, using the pipeline definition in the "app" / "deployable" repo.

We still need to define if we auto-create the tags or not, and which form they take. My current thinking is:

* development environment uses `version: WIP`, creates no tags, just deploys the configured branch (typically `master`)
* QA environment uses e.g. `version: 1.0.0` and auto-creates tags of form `vVERSION-rc.NUMBER` from either the configured branch (typically `master`) or a more specific release branch (`release/VERSION`) if it exists.
* PROD env uses e.g. `version: 1.0.0` and selects the matching tag with the highest `rc` number, auto-creates a corresponding `vVERSION` tag. The pipeline aborts if there is no matching rc tag.

## Questions

* If the Git reference is a tag, should we auto-create them? Maybe we should allow different kinds of "versions": either a Git ref (which then is looked up), or a conceptual version (such as 1.1.3) which creates corresponding tags.
* How to pass the values to the Helm chart properly, and how to include the charts from the subrepos? Maybe we need to have an "environment type" such as dev, qa, staging, prod so that multiple prod environments can make use of the same settings.

## Nice properties

* The state of each environment can be seen by looking at the branch's `version` information.
* Promoting from QA to PROD is by pull request from `qa` branch to `prod` branch.
* No unwanted deploys if environments are edited in the wrong branch. Only the environment matching the branch is deployed.
* In theory the app repo can be merged into one monorepo with the subrepos.
* One can create a separate namespace (e.g. `foo-deploy`) and a second set of "release" and "app" repos. Their Bitbucket webhook points to an event listener in the separate namespace. This 2nd repo/namespace can be used for restricted access so that developers can build the application, but do not see or have access to e.g. the production environment
* VERSION is opaque - it can be any string that is valid in a Git tag.
