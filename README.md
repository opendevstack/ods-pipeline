# ODS Pipeline

[![Tekton Tasks Tests](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yaml/badge.svg)](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yaml)

**IMPORTANT: This approach is currently being evaluated. This may or may not become part of ODS one day.**

## Introduction

This repository provides an alternative OpenDevStack CI/CD pipeline based on OpenShift Pipelines. This repository contains everything that relates to it, such as Tekton tasks, container images, Go packages, services, documentation, ...

The [ODS Pipeline Introduction](/docs/introduction.adoc) describes what ODS pipeline is and how it works. It is important to understand this before looking at further documentation or any other repository content.

ODS Pipeline is well suited for regulated development (e.g. medical device software development) and many trade-offs and design decisions have been made to support this.

## Documentation

**User Guide**
* [Installation & Updating](/docs/installation.adoc)
* [Getting Started](/docs/getting-started.adoc)
* [ODS.YAML Reference](/docs/ods-configuration.adoc)
* [Task Reference](/docs/tasks)
* [Working with secrets in Helm](/docs/helm-secrets.adoc)
* [Accessing artifacts](/docs/accessing-artifacts.adoc)
* [Debugging](/docs/debugging.adoc)
* [Authoring Tasks](/docs/authoring-tasks.adoc)
* [Example Project](/docs/example-project.adoc)
* [FAQ](https://github.com/opendevstack/ods-pipeline/wiki/FAQ)

**Contributor Guide**
* [Repository Layout](/docs/repository-layout.adoc)
* [Development & Running Tests](/docs/development.adoc)
* [Artifacts](/docs/artifacts.adoc)
* [Creating an ODS task](/docs/creating-an-ods-task.adoc)
* [Releasing a new version](/docs/releasing.adoc)

This repository also hosts the design documents that describe ODS pipeline more formally. Those design documents provide more detail and background on goals, requirements and architecture decisions and may be useful for all audiences.

* [Stakeholder Requirements](/docs/design/stakeholder-requirements.adoc)
* [Software Architecture](/docs/design/software-architecture.adoc)
* [Software Requirements Specification](/docs/design/software-requirements-specification.adoc)
* [Software Design Specification](/docs/design/software-design-specification.adoc)
* [Architecture Decision Records](/docs/adr)
* [Goals and Non-Goals](/docs/design/goals-and-nongoals.adoc)
* [Relationship to Jenkins Shared Library](/docs/design/relationship-shared-library.adoc)

## Compatibility

For OpenShift Pipelines releases and its relationship to Tekton and OpenShift versions, see https://docs.openshift.com/container-platform/4.8/cicd/pipelines/op-release-notes.html

| ods-pipeline | OpenShift Pipelines | ODS Core/Quickstarters |
|---|---|---|
| [0.2](https://github.com/opendevstack/ods-pipeline/milestone/2) | 1.5 | 4.0.0 |
| [0.1](https://github.com/opendevstack/ods-pipeline/milestone/1) | 1.5 | 4.0.0 |
