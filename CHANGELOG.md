# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Note that changes which ONLY affect documentation or the testsuite will not be
listed in the changelog.

## [Unreleased]
### Added

- Provide separate binary to download all artifacts related to one version easily ([#167](https://github.com/opendevstack/ods-pipeline/issues/167))
- Allow namespaced installation. This provides a way to give ODS pipeline a try without requiring the buy-in from a cluster admin. The OpenShift Pipelines operator is still required though. See [#263](https://github.com/opendevstack/ods-pipeline/issues/263).
- Automated check if the Docker host has enough memory ([#283](https://github.com/opendevstack/ods-pipeline/issues/283))
- Create SonarQube quality gate artifact ([#273](https://github.com/opendevstack/ods-pipeline/issues/273))
- Make task prefix customizable ([#289](https://github.com/opendevstack/ods-pipeline/issues/289))
- Add overridable test timeout to Makefile ([#284](https://github.com/opendevstack/ods-pipeline/issues/284))
- Skipping Tests in TypeScript build task if test artifacts are present already ([#238](https://github.com/opendevstack/ods-pipeline/issues/238))

### Changed

- Changed encryption tool for helm secrets plugin from `gpg` to `age` ([#292](https://github.com/opendevstack/ods-pipeline/pull/292))
- Automatically roll webhook interceptor deployment when related config map or secret changes ([#252](https://github.com/opendevstack/ods-pipeline/issues/252))
- Hide confusing error message in Helm output ([#262](https://github.com/opendevstack/ods-pipeline/issues/262))
- Update gcc (from 8.4 to 8.5), skopeo (from 1.3 to 1.4) and buildah (from 1.21 to 1.22) in container images ([#276](https://github.com/opendevstack/ods-pipeline/pull/276))
- Iterating over Dockerfiles in build/package instead of using hardcoded list ([#286](https://github.com/opendevstack/ods-pipeline/issues/286))
- Upgrade Python toolset to v3.9, with migration from Flask to FastAPI sample app ([#312](https://github.com/opendevstack/ods-pipeline/issues/312))

### Fixed

- Release branches of subrepos are not detected ([#269](https://github.com/opendevstack/ods-pipeline/pull/269))
- `Directory` values in the artifact manifest (`.ods/artifacts/manifest.json`) contain an erronous leading slash. This should only be an issue if you relied on this value in a custom task. ([#269](https://github.com/opendevstack/ods-pipeline/pull/269))
- `ods-finish` does not upload artifacts of subrepos ([#257](https://github.com/opendevstack/ods-pipeline/issues/257))
- Waitfor-...sh scripts are not waiting for the expected 5 minutes ([#280](https://github.com/opendevstack/ods-pipeline/issues/280))
- Specifying images to be build and pushed is not working anymore due to changes made in [#287](https://github.com/opendevstack/ods-pipeline/pull/287) ([#299](https://github.com/opendevstack/ods-pipeline/issues/299))

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
