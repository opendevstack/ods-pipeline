# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Note that changes which ONLY affect documentation or the testsuite will not be
listed in the changelog.

## [Unreleased]

### Changed

- Update Git from 2.31 to 2.39 ([#693](https://github.com/opendevstack/ods-pipeline/pull/693))
- Update Skopeo from 1.9 to 1.11 ([#693](https://github.com/opendevstack/ods-pipeline/pull/693))
- Update Buildah from 1.27 to 1.29 ([#693](https://github.com/opendevstack/ods-pipeline/pull/693))

### Fixed

- Artifact upload fails when no manifest is present ([#694](https://github.com/opendevstack/ods-pipeline/issues/694))
- Unexpected event payload leads to server panic which is hard to debug ([#693](https://github.com/opendevstack/ods-pipeline/pull/693))

## [0.12.0] - 2023-05-12

### Added

- Support webhook events related to Git tags ([#608](https://github.com/opendevstack/ods-pipeline/issues/608))
- Support Git tags for subrepositories (related to [#630](https://github.com/opendevstack/ods-pipeline/issues/630))

### Changed

- IMPORTANT! The trigger mechanism allows for dynamic parameterisation of pipelines now (see [#677](https://github.com/opendevstack/ods-pipeline/issues/677) for the original idea). As a consequence, a few things have changed. The `pipeline` field is now named `pipelines`, and must specify a list of pipelines. Further, the `trigger` field of a pipeline is now named `triggers` and also specifies a list now. Inside each trigger, the `event` field was renamed to `events` for consistency. Further, trigger has learned a new field, `params`, which allows to specify pipeline and task parameters. The name of task parameters are suffixed with the task name, e.g. `<task-name>.some-param. Finally, the `branchToEnvironmentMapping` and `environments` fields have been dropped and equivalent behaviour must now be configured through the use of trigger parameters.
- IMPORTANT! Setting a version in `ods.yaml` is no longer supported. Consequently, subrepositories are always checked out at the specified revision - no release branch "matching" the specified version is preferred. Further the version of the Helm chart isn't modified anymore. ([#630](https://github.com/opendevstack/ods-pipeline/issues/630))
- The created pipeline run artifact records the Git commit SHAs of each checked out subrepository now (related to [#630](https://github.com/opendevstack/ods-pipeline/issues/630)).
- The `artifact-download` binary is expected to be run from the Git commit in the umbrella repository now for which artifacts should be downloaded (related to [#630](https://github.com/opendevstack/ods-pipeline/issues/630)).

### Fixed

- Removed addition of `always-auth=true` to the npm config file for nodeJS builds ([#687](https://github.com/opendevstack/ods-pipeline/issues/687))

## [0.11.1] - 2023-03-31

### Fixed

- Configure Git to use bearer token auth mechanism ([#683](https://github.com/opendevstack/ods-pipeline/issues/683))
- `make create-kind-with-registry` failed locally on Mac ([#679](https://github.com/opendevstack/ods-pipeline/issues/679))

## [0.11.0] - 2023-03-14

### Changed

- Upgrade to SonarQube v8.9 LTS (see [#424](https://github.com/opendevstack/ods-pipeline/issues/424)). Note that this is a breaking change: 0.10.1 and prior will not work with SonarQube >= 8.9, and all future ODS Pipeline versions will not work with SonarQube < 8.9.

## [0.10.1] - 2023-03-03

### Changed

- Updated `_auth` value in npm config to be scoped to the npm registry in Nexus [#668](https://github.com/opendevstack/ods-pipeline/issues/668)

### Fixed

- Setup of secrets during `install.sh` does not work if secret contains `/` ([#670](https://github.com/opendevstack/ods-pipeline/issues/670))
- `ods-start` is unable to cleanup workspace for some storage configurations due to changes in `ods-package-image` ([#672](https://github.com/opendevstack/ods-pipeline/issues/672))
- `runAfter: [start]` not set for all parallel tasks ([#671](https://github.com/opendevstack/ods-pipeline/issues/671))

## [0.10.0] - 2023-02-27

### Added

- Rendered task versions are available under `tasks/` now. These can be referenced directly from pipeline runs through [remote resolution](https://tekton.dev/docs/pipelines/pipelines/#specifying-remote-tasks). In future versions, tasks may be removed from the Helm chart and only be accessible via Git. See [#665](https://github.com/opendevstack/ods-pipeline/issues/665).

### Changed

- Update Tekton to v.041.1 (matching OpenShift Pipelines operator 1.9). Unfortunately the `package-image` task of v0.9.0 breaks in OpenShift Pipelines operator 1.9, so v0.9.0 is not compatible with 1.9, and v0.10.0 will not compatible with 1.6. For details of the change, see [#663](https://github.com/opendevstack/ods-pipeline/pull/663).


## [0.9.0] - 2023-02-22

### Added

- New image for `ods-build-npm` task with Node.js 18 ([#585](https://github.com/opendevstack/ods-pipeline/issues/585))
- Add `extra-tags` parameter to `ods-package-image` ([#631](https://github.com/opendevstack/ods-pipeline/issues/631))

### Changed

- Node.js 18 is now the default for `ods-build-npm` task ([#585](https://github.com/opendevstack/ods-pipeline/issues/585))
- Images used in tasks are now pulled directly from the GitHub registry. "Wrapping" the images in the OpenShift/K8s cluster is not required anymore. If tasks need to trust a private certificate, it needs to be present as a K8s secret, which will then be mounted as a file in the pods. To add the secret to an existing installation, pass `--private-cert <host>` to `./install.sh`. For more details, see [#621](https://github.com/opendevstack/ods-pipeline/issues/621).
- Remove PVC use protection ([#647](https://github.com/opendevstack/ods-pipeline/issues/647))
- Use Go 1.19 for building ([#659](https://github.com/opendevstack/ods-pipeline/issues/659))
- Pipeline manager now returns `application/problem+json` when it encounters an error. Further, it now returns different, better fitting error status codes for some responses. See ([#661](https://github.com/opendevstack/ods-pipeline/issues/661)) for details.

### Fixed

- npm-toolset tests fail with new release of ubi8 Node.js image ([#650](https://github.com/opendevstack/ods-pipeline/issues/650))
- Installation does not ask for Bitbucket username ([#652](https://github.com/opendevstack/ods-pipeline/issues/652))
- e2e environment name not allowed ([#634](https://github.com/opendevstack/ods-pipeline/issues/634))

## [0.8.0] - 2023-01-06

### Added

- Add trivy security scanner CLI for SBOM generation ([#592](https://github.com/opendevstack/ods-pipeline/issues/592))

### Changed

- Normalize K8s manifests to exclude style differences from Helm diff output. The change is applied to both the helm execution in the `ods-deploy-helm` task and in the install script. See [#591](https://github.com/opendevstack/ods-pipeline/issues/591).
- Update skopeo (1.8 to 1.9) ([#616](https://github.com/opendevstack/ods-pipeline/issues/616))
- Update buildah (1.26 to 1.27) ([#626](https://github.com/opendevstack/ods-pipeline/issues/626))
- Stream Helm upgrade log output ([#615](https://github.com/opendevstack/ods-pipeline/issues/615))
- Update Go to 1.18 ([#623](https://github.com/opendevstack/ods-pipeline/issues/623))
- Update go-junit-report to 2.0.0 ([#625](https://github.com/opendevstack/ods-pipeline/issues/625))
- Enable build skipping by default ([#642](https://github.com/opendevstack/ods-pipeline/issues/642))
- Remove secrets from installation Helm chart. Secrets are now managed when running the `install.sh` script. See [#629](https://github.com/opendevstack/ods-pipeline/issues/629).
- Change name of `buildah` task to `package-image` ([#592](https://github.com/opendevstack/ods-pipeline/issues/592))
- Package image task now skips creating an image if the image artifact exists (as opposed to checking for an image in the registry) ([#592](https://github.com/opendevstack/ods-pipeline/issues/592))

### Fixed

- Errors during output collection of binaries such as `buildah`, `aqua-scanner` are not handled ([#611](https://github.com/opendevstack/ods-pipeline/issues/611))
- STDOUT and STDERR is not interleaved as expected ([#613](https://github.com/opendevstack/ods-pipeline/issues/613))


## [0.7.0] - 2022-10-11

### Changed

- Stream Buildah and Aqua log output ([#596](https://github.com/opendevstack/ods-pipeline/issues/596))
- Update skopeo (1.6 to 1.8) and buildah (1.24 to 1.26) ([#598](https://github.com/opendevstack/ods-pipeline/issues/598))
- Support running different pipelines on different webhook events (current `ods.yaml` format still supported for the moment, but will be deprecated and removed in upcoming releases) ([#562](https://github.com/opendevstack/ods-pipeline/issues/562))

### Fixed

- Aqua and helm-diff log output is incomplete ([#593](https://github.com/opendevstack/ods-pipeline/issues/593))
- Image is tagged with `latest` instead of correct tag when pushed to external registry ([#606](https://github.com/opendevstack/ods-pipeline/issues/606))

## [0.6.0] - 2022-07-20

### Changed

- Use `PipelineRun` resources with inlined spec instead of managing and referencing a `Pipeline` resource ([#573](https://github.com/opendevstack/ods-pipeline/issues/573))


### Added

- Enable build caching for gradle builds according to `docs/adr/20220314-caching-build-tasks.md`.
- Add script to install from inside an OpenShift Web Terminal ([#581](https://github.com/opendevstack/ods-pipeline/issues/581))

### Changed

- Add timestamp and tag to log messages in pipeline manager deployment ([#554](https://github.com/opendevstack/ods-pipeline/issues/554))
- Perform 3-way merge in `install.sh` script ([#569](https://github.com/opendevstack/ods-pipeline/issues/569))
- Repurpose `diff-flags` parameter of `deploy-helm` task and use 3-way merge in diff by default ([#574](https://github.com/opendevstack/ods-pipeline/issues/574) and [#569](https://github.com/opendevstack/ods-pipeline/issues/569))

### Fixed
- Trailing slash in service URLs is not handled properly ([#526](https://github.com/opendevstack/ods-pipeline/issues/526))
- Helm diff result log filtering does not work anymore ([#563](https://github.com/opendevstack/ods-pipeline/issues/563))
- Handle Helm diff error and diff detection separately ([#584](https://github.com/opendevstack/ods-pipeline/issues/584))
- Handle Aqua error and compliance problems separately ([#586](https://github.com/opendevstack/ods-pipeline/issues/586))

## [0.5.1] - 2022-06-10

### Fixed
- Pipeline runs are not pruned ([#557](https://github.com/opendevstack/ods-pipeline/issues/557))
- Pipeline pods are listed under deployment pods ([#555](https://github.com/opendevstack/ods-pipeline/issues/555))

## [0.5.0] - 2022-06-09

### Added

- Automatically build images after Helm upgrade ([#525](https://github.com/opendevstack/ods-pipeline/issues/525))
- Allow to use build script located in repository ([#536](https://github.com/opendevstack/ods-pipeline/issues/536))

### Changed

- Avoid need to specify `imageTag` and `taskSuffix` in `values.yaml` ([#551](https://github.com/opendevstack/ods-pipeline/issues/551))

### Fixed
- Wrapper image cannot write aquasec binary ([#539](https://github.com/opendevstack/ods-pipeline/issues/539))
- When a commit is skipped, the log message contains weird output ([#542](https://github.com/opendevstack/ods-pipeline/issues/542))
- `imageTag` not defaulting to `.Chart.AppVersion` in `ods-finish` task ([#547](https://github.com/opendevstack/ods-pipeline/issues/547))
- `taskSuffix` defaults to `-v0-3-0` in release 0.4.0 ([#546](https://github.com/opendevstack/ods-pipeline/issues/546))
- Add (missing) common labels to resources in images and tasks charts ([#543](https://github.com/opendevstack/ods-pipeline/issues/543))


## [0.4.0] - 2022-05-31

### Added

- Support for optional build task caching. The main use case is to avoid lengthy builds in repos with multiple build tasks ([#461](https://github.com/opendevstack/ods-pipeline/issues/461)). See the `docs/adr/20220314-caching-build-tasks.md`for details.
- Provide Apple silicon builds of artifact-download binary ([#510](https://github.com/opendevstack/ods-pipeline/issues/510))

### Changed

- Default imageTag to appVersion + release images without leading v ([#504](https://github.com/opendevstack/ods-pipeline/issues/504))
- Display both version and Git commit SHA in `artifact-download -version` ([#507](https://github.com/opendevstack/ods-pipeline/issues/507))
- Add more context to Bitbucket client errors ([#515](https://github.com/opendevstack/ods-pipeline/issues/515))
- Update skopeo to 1.6, buildah to 1.24 and git to 2.31 ([#519](https://github.com/opendevstack/ods-pipeline/issues/519))
- Update Go to 1.17 ([#528](https://github.com/opendevstack/ods-pipeline/issues/528))
- Rename `ods-build-typescript` task to `ods-build-npm` ([#503](https://github.com/opendevstack/ods-pipeline/issues/503))
- Implement global caching for Gradle build task ([#490](https://github.com/opendevstack/ods-pipeline/issues/490))
- Run `lint` script instead of `eslint` directly ([#532](https://github.com/opendevstack/ods-pipeline/issues/532))

### Fixed
- Pipelines fail in clusters with private / self-signed certificates ([#518](https://github.com/opendevstack/ods-pipeline/issues/518))
- HTTP_PROXY setting is not taken into account when building the gradle-toolset image via a wrapper image in the target cluster. ([#530](https://github.com/opendevstack/ods-pipeline/pull/530))

### Removed

- Remove test skipping from Go build task ([#493](https://github.com/opendevstack/ods-pipeline/issues/493))
- Remove test skipping from TypeScript build task ([#494](https://github.com/opendevstack/ods-pipeline/issues/494))
- Remove `sonarUsername` from `values.yaml` as only the auth token is used ([#514](https://github.com/opendevstack/ods-pipeline/issues/514))

## [0.3.0] - 2022-04-07

### Added

- Apply labels to pipelines allowing easier identification for cleanup ([#358](https://github.com/opendevstack/ods-pipeline/issues/358))
- Configurable workspace PVC size ([#368](https://github.com/opendevstack/ods-pipeline/issues/368))
- Customizable Helm flags ([#388](https://github.com/opendevstack/ods-pipeline/issues/388))
- Run gradle in non daemon mode by default and enabling stacktraces ([#386](https://github.com/opendevstack/ods-pipeline/issues/386))
- Enable setting `GRADLE_OPTS` via task parameters ([#387](https://github.com/opendevstack/ods-pipeline/issues/387))
- Export `ARTIFACTS_DIR` environment variable to be visible to the gradle build ([#408](https://github.com/opendevstack/ods-pipeline/issues/408))
- Add notifications via configurable webhook call from `ods-finish` ([#140](https://github.com/opendevstack/ods-pipeline/issues/140))
- Git LFS support enabled ([#420](https://github.com/opendevstack/ods-pipeline/issues/420))
- Publish images to public registry (ghcr.io) ([#440](https://github.com/opendevstack/ods-pipeline/issues/440))
- Allow to cache dependencies and support for caching go dependencies ([#147](https://github.com/opendevstack/ods-pipeline/issues/147)). See also proposal on caching ([#412](https://github.com/opendevstack/ods-pipeline/pull/412))
- Support node production builds in docker context. It is now required that both `package.json` and `package-lock.json` are available to the build. ([#357](https://github.com/opendevstack/ods-pipeline/issues/357))
- Allow to select which tasks (and related BC/IS resources) to install ([#486](https://github.com/opendevstack/ods-pipeline/issues/486))
- Upload artifacts of unsuccessful pipeline runs as well ([#379](https://github.com/opendevstack/ods-pipeline/issues/379))

### Changed

- Use configmap `ods-sonar` to configure SonarQube edition ([#410](https://github.com/opendevstack/ods-pipeline/issues/410))
- Prevent existing image streams from being cleaned up if they are renamed in future versions ([#366](https://github.com/opendevstack/ods-pipeline/issues/366))
- Add `build-dir` and `copy-node-modules` (defaulting to `false`) parameters to TypeScript build task to make it more suitable for FE builds. A non-obvious but breaking change is that files inside the directory specified `build-dir` are now copied to folder `${output-dir}/dist` whereas previously they were copied to `${output-dir}/dist/dist`([#356](https://github.com/opendevstack/ods-pipeline/issues/356))
- Update gradle version to 7.3.3 to address log4j vulnerability and improved JDK 17 support. ([#395](https://github.com/opendevstack/ods-pipeline/issues/395))
- Create and use one PVC per repository ([#160](https://github.com/opendevstack/ods-pipeline/issues/160))
- Leaner NodeJS 16 Typescript image task, removed cypress and its dependencies ([#426](https://github.com/opendevstack/ods-pipeline/issues/426))
- Update skopeo (from 1.4 to 1.5) and buildah (from 1.22 to 1.23) ([#430](https://github.com/opendevstack/ods-pipeline/issues/430))
- Use `--ignore-scripts` when building TypeScript apps ([#434](https://github.com/opendevstack/ods-pipeline/issues/434))
- Prune pipelines and pipeline runs ([#153](https://github.com/opendevstack/ods-pipeline/issues/153)). Note that any pipeline runs created with 0.2.0 or earlier will not be cleaned up and need to be dealt with manually.
- Log artifact URL after upload ([#384](https://github.com/opendevstack/ods-pipeline/issues/384))
- Remove Tekton Triggers, moving the required functionality it provided into the new ODS pipeline manager ([#438](https://github.com/opendevstack/ods-pipeline/issues/438))
- Use UBI8 provided Python 3.9 toolset image ([#457](https://github.com/opendevstack/ods-pipeline/issues/457))
- Change installation mode from centralized to local/namespaced ([#404](https://github.com/opendevstack/ods-pipeline/pull/404))
- Removed logging of test reports for TypeScript and Python build tasks ([#470](https://github.com/opendevstack/ods-pipeline/issues/470))
- Don't remove tasks on `helm` upgrades, rollbacks, etc. ([#477](https://github.com/opendevstack/ods-pipeline/issues/477))
- Run go fmt over packages, not entire directory ([#484](https://github.com/opendevstack/ods-pipeline/issues/484))
- Update `golangci-lint` from 1.41 to 1.45 ([#497](https://github.com/opendevstack/ods-pipeline/pull/497))
- Improve build time of subsequent local container image builds ([#499](https://github.com/opendevstack/ods-pipeline/pull/499))
- Refactor pipeline manager. This moves the endpoint of the webhook receiver to `/bitbucket`, as a consequence every webhook configuration in Bitbucket needs to be updated ([#491](https://github.com/opendevstack/ods-pipeline/issues/491))

### Fixed
- Cannot enable debug mode in some tasks ([#377](https://github.com/opendevstack/ods-pipeline/issues/377))
- Gradle task does not expose Nexus env variables ([#373](https://github.com/opendevstack/ods-pipeline/issues/373))
- Gradle build fails when it contains more than one test class ([#414](https://github.com/opendevstack/ods-pipeline/issues/414))
- Gradle proxy settings are set during prepare-local-env ([#291](https://github.com/opendevstack/ods-pipeline/issues/291))
- Add `xargs` to helm image as `helm-secrets` depends on it ([#465](https://github.com/opendevstack/ods-pipeline/issues/465))
- Pipeline creation fails when branch names contain slashes ([#466](https://github.com/opendevstack/ods-pipeline/issues/466))
- Race conditions between pipelines of the same repository ([#394](https://github.com/opendevstack/ods-pipeline/issues/394))

## [0.2.0] - 2021-12-22
### Added

- Provide nexus `-build-arg` variables during image building with ods-package-image ([#327](https://github.com/opendevstack/ods-pipeline/issues/327))
- Provide separate binary to download all artifacts related to one version easily ([#167](https://github.com/opendevstack/ods-pipeline/issues/167))
- Allow namespaced installation. This provides a way to give ODS pipeline a try without requiring the buy-in from a cluster admin. The OpenShift Pipelines operator is still required though. See [#263](https://github.com/opendevstack/ods-pipeline/issues/263).
- Automated check if the Docker host has enough memory ([#283](https://github.com/opendevstack/ods-pipeline/issues/283))
- Create SonarQube quality gate artifact ([#273](https://github.com/opendevstack/ods-pipeline/issues/273))
- Make task prefix customizable ([#289](https://github.com/opendevstack/ods-pipeline/issues/289))
- Add overridable test timeout to Makefile ([#284](https://github.com/opendevstack/ods-pipeline/issues/284))
- Skipping Tests in TypeScript build task if test artifacts are present already ([#238](https://github.com/opendevstack/ods-pipeline/issues/238))
- Provide make target for ShellCheck and added ShellCheck to GitHub actions ([#240](https://github.com/opendevstack/ods-pipeline/issues/240))
- Supply default `sonar-project.properties` when none is present, configuring SonarQube out-of-the-box ([#296](https://github.com/opendevstack/ods-pipeline/issues/296))
- Linting step in TypeScript build task ([#325](https://github.com/opendevstack/ods-pipeline/issues/325))
- Set `CI=true` in build tasks ([#336](https://github.com/opendevstack/ods-pipeline/issues/336))
- Generate report for successful linting step in Go build task ([#215](https://github.com/opendevstack/ods-pipeline/issues/215))

### Changed

- Changed encryption tool for helm secrets plugin from `gpg` to `age` ([#293](https://github.com/opendevstack/ods-pipeline/issues/293), [#292](https://github.com/opendevstack/ods-pipeline/pull/292))
- Automatically roll webhook interceptor deployment when related config map or secret changes ([#252](https://github.com/opendevstack/ods-pipeline/issues/252))
- Hide confusing error message in Helm output ([#262](https://github.com/opendevstack/ods-pipeline/issues/262))
- Update gcc (from 8.4 to 8.5), skopeo (from 1.3 to 1.4) and buildah (from 1.21 to 1.22) in container images ([#276](https://github.com/opendevstack/ods-pipeline/pull/276))
- Iterating over Dockerfiles in `build/package` instead of using hardcoded list ([#286](https://github.com/opendevstack/ods-pipeline/issues/286) and [#287](https://github.com/opendevstack/ods-pipeline/pull/287))
- Upgrade Python toolset to v3.9, with migration from Flask to FastAPI sample app ([#312](https://github.com/opendevstack/ods-pipeline/issues/312))
- Upgrade Java toolset to JDK 17 ([#294](https://github.com/opendevstack/ods-pipeline/issues/294))
- Set Helm value `image.tag` instead of `gitCommitSha` ([#342](https://github.com/opendevstack/ods-pipeline/pull/342))
- Provide TypeScript toolset with Node.js 16 ([#337](https://github.com/opendevstack/ods-pipeline/issues/337))
- Use `ubi8/go-toolset` as consistent builder image and as image for the `ods-build-go` task ([#295](https://github.com/opendevstack/ods-pipeline/issues/295))

### Fixed

- Release branches of subrepos are not detected ([#269](https://github.com/opendevstack/ods-pipeline/pull/269))
- `directory` values in the artifact manifest (`.ods/artifacts/manifest.json`) contain an erronous leading slash. This should only be an issue if you relied on this value in a custom task. ([#269](https://github.com/opendevstack/ods-pipeline/pull/269))
- `ods-finish` does not upload artifacts of subrepos ([#257](https://github.com/opendevstack/ods-pipeline/issues/257))
- Waitfor-...sh scripts are not waiting for the expected 5 minutes ([#280](https://github.com/opendevstack/ods-pipeline/issues/280))
- Tagging in `ods-start` causes second pipeline run ([#331](https://github.com/opendevstack/ods-pipeline/issues/331))
- Helm resource names differ between component and umbrella repository ([#340](https://github.com/opendevstack/ods-pipeline/issues/340))
- Commercial SonarQube capabilities are not detected because `SONAR_EDITION` is not set in `ods-sonar` ([#350](https://github.com/opendevstack/ods-pipeline/issues/350))

## [0.1.1] - 2021-10-28
### Fixed

- Incorrect ods-gradle-toolset base image in BuildConfig ([#250](https://github.com/opendevstack/ods-pipeline/issues/250))
- Generating a SonarQube report fails when PR exists for scanned branch ([#227](https://github.com/opendevstack/ods-pipeline/issues/227))
- Generating a SonarQube report fails when background task does not finish immediately ([#227](https://github.com/opendevstack/ods-pipeline/issues/227))
- Quality Gate check fails due to incorrect API authentication (detected while working on #227)
- Uploading of artifacts in `ods-finish` may fail when artifact is already present from previous pipeline run ([#255](https://github.com/opendevstack/ods-pipeline/issues/255))
- Misleading error message when interceptor is forbidden to retrieve ods.yaml ([#254](https://github.com/opendevstack/ods-pipeline/issues/254))

### Changed

- Suffix Helm release name of cluster tasks with version to enable retention of tasks from previous versions ([#234](https://github.com/opendevstack/ods-pipeline/issues/234))

## [0.1.0] - 2021-10-05

Initial version.
