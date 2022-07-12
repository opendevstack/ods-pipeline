# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Note that changes which ONLY affect documentation or the testsuite will not be
listed in the changelog.

## [Unreleased]

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
