# ODS Pipeline

[![Tekton Tasks Tests](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yml/badge.svg)](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yml)

**IMPORTANT: This is EXPERIMENTAL ONLY. This may or may not become part of ODS one day.**

## Introduction

ODS provides CI/CD pipeline support based on OpenShift Pipelines. This repository contains everything that relates to it, such as Tekton tasks, container images, Go packages, services, documentation, ...

The "user perspective" of what the ODS pipeline is and how it works is described in the [ODS Pipeline Introduction](/docs/introduction.adoc). It is important to understand this before looking at this repository, which is the actual "plumbing".

## How is this repository organized?

The repo follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

The most important pieces are:

* **build/package**: `Dockerfile`s for the various container images in use. These images back Tekton tasks or the webhook interceptor.
* **cmd**: Main executables. These are installed (in different combinations) into the contaier images.
* **deploy**: OpenShift/K8S resource definitions, such as `BuildConfig`/`ImageStream` or `ClusterTask` resources. The tasks typically make use of the images built via `build/package` and their `script` calls one or more executables built from the `cmd` folder.
* **docs**: Design and user documents
* **internal/interceptor**: Implementation of Tekton trigger interceptor - it creates and modifies the actual Tekton pipelines on the fly based on the config found in the repository triggering the webhook request.
* **pkg**: Packages shared by the various main executables and the interceptor. These packages are the public interface and may be used outside this repo (e.g. by custom tasks). Example of packages are `bitbucket` (a Bitbucket Server API v1.0 client), `sonar` (a SonarQube client exposing API endpoints, scanner CLI and report CLI in one unified interface), `nexus` (a Nexus client for uploading, downloading and searching for assets) and `config` (the ODS configuration specification).
* **test**: Test scripts and test data

## Details / Documentation

* [Introduction](/docs/introduction.adoc)
* [Goals and Non-Goals](/docs/goals-and-nongoals.adoc)
* [Architecture Decision Records](/docs/adr)
* [Task Library](/docs/tasks)
* [Authoring Tasks](/docs/authoring-tasks.adoc)

## Development & Running tests

First, check if your system meets the prerequisites:
```
make check-system
```

Then, launch a KinD cluster, install Tekton, build&push images and run services:
```
make prepare-local-env
```

Finally, run the tests:
```
make test
```

More fine-grained make targets are available, see:
```
make help
```

## Compatibility

 For OpenShift Pipelines releases and its relationship to Tekton and OpenShift versions, see https://docs.openshift.com/container-platform/4.8/cicd/pipelines/op-release-notes.html

 | ods-pipeline | OpenShift Pipelines | ODS Core/Quickstarters |
 |---|---|---|
 | 0.1 (to be released end of August 2021) | 1.5 | 4.0.0 |
 | 0.2 (to be released TBD) | 1.5 | 4.0.0 |
