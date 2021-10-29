# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Note that changes which ONLY affect documentation or the testsuite will not be
listed in the changelog.

## [Unreleased]
### Changed

- Automatically roll webhook interceptor deployment when related config map or secret changes ([#252](https://github.com/opendevstack/ods-pipeline/issues/252))

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
