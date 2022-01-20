# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Note that changes which ONLY affect documentation or the testsuite will not be
listed in the changelog.

## [Unreleased]

### Added

- Apply labels to pipelines allowing easier identification for cleanup ([#358](https://github.com/opendevstack/ods-pipeline/issues/358))
- Configurable workspace PVC size ([#368](https://github.com/opendevstack/ods-pipeline/issues/368))
- Customizable Helm flags ([#388](https://github.com/opendevstack/ods-pipeline/issues/388))
- Run gradle in non daemon mode by default and enabling stacktraces ([#386](https://github.com/opendevstack/ods-pipeline/issues/386)) 
- Enable setting `GRADLE_OPTS` via task parameters ([#387](https://github.com/opendevstack/ods-pipeline/issues/387))
- Export `ARTIFACTS_DIR` environment variable to be visible to the gradle build ([#408](https://github.com/opendevstack/ods-pipeline/issues/408))

### Changed

- Prevent existing image streams from being cleaned up if they are renamed in future versions ([#366](https://github.com/opendevstack/ods-pipeline/issues/366))
- Add `build-dir` and `copy-node-modules` parameters to TypeScript build task to make it more suitable for FE builds ([#356](https://github.com/opendevstack/ods-pipeline/issues/356))
- Update gradle version to 7.3.3 to address log4j vulnerability and improved JDK 17 support. ([#395](https://github.com/opendevstack/ods-pipeline/issues/395))

### Fixed
- Cannot enable debug mode in some tasks ([#377](https://github.com/opendevstack/ods-pipeline/issues/377))
- Gradle task does not expose Nexus env variables ([#373](https://github.com/opendevstack/ods-pipeline/issues/373))
- Gradle build fails when it contains more than one test class ([#414](https://github.com/opendevstack/ods-pipeline/issues/414))

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
