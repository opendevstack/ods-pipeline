# ODS Pipeline

[![Tekton Tasks Tests](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yaml/badge.svg)](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yaml)

**IMPORTANT: This approach is currently being evaluated. This may or may not become part of ODS one day.**

## Introduction

ODS Pipeline provides a  CI/CD pipeline based on OpenShift Pipelines as an alternative to Jenkins. This repository contains everything that relates to it, such as Tekton tasks, container images, Go packages, services, documentation, ...

The [ODS Pipeline Introduction](/docs/introduction.adoc) describes what ODS pipeline is and how it works. It is important to understand this before looking at further documentation or any other repository content.

ODS Pipeline is well suited for regulated development (e.g. medical device software development) and many trade-offs and design decisions have been made to support this.

## Documentation

### Getting Started
* [Installation & Updating](/docs/installation.adoc)
* [Add ODS Pipeline to a repository](/docs/add-to-repository.adoc)
* [Convert an ODS quickstarter based component](/docs/convert-quickstarter-component.adoc)
* [Example Project](/docs/example-project.adoc)

### Technical Reference
* [Repository configuration (ods.yaml)](/docs/ods-configuration.adoc)
* Plumbing tasks: [ods-start](/docs/tasks/ods-start.adoc), [ods-finish](/docs/tasks/ods-finish.adoc)
* Build tasks: [ods-build-go](/docs/tasks/ods-build-go.adoc), [ods-build-gradle](/docs/tasks/ods-build-gradle.adoc), [ods-build-python](/docs/tasks/ods-build-python.adoc), [ods-build-typescript](/docs/tasks/ods-build-typescript.adoc)
* Package tasks: [ods-package-image](/docs/tasks/ods-package-image.adoc)
* Deploy tasks: [ods-deploy-helm](/docs/tasks/ods-deploy-helm.adoc)

### How-To Guides
* [Working with secrets in Helm](/docs/helm-secrets.adoc)
* [Accessing artifacts](/docs/accessing-artifacts.adoc)
* [Debugging](/docs/debugging.adoc)
* [Authoring your own tasks](/docs/authoring-tasks.adoc)
* [FAQ](https://github.com/opendevstack/ods-pipeline/wiki/FAQ)

### Background Information
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
| [0.3](https://github.com/opendevstack/ods-pipeline/releases/tag/v0.3.0) | 1.5 | 4.0.0 |
| [0.2](https://github.com/opendevstack/ods-pipeline/releases/tag/v0.2.0) | 1.5 | 4.0.0 |
| [0.1](https://github.com/opendevstack/ods-pipeline/releases/tag/v0.1.1) | 1.5 | 4.0.0 |

## Contributing

* [Repository Layout](/docs/repository-layout.adoc)
* [Development & Running Tests](/docs/development.adoc)
* [Artifacts](/docs/artifacts.adoc)
* [Creating an ODS task](/docs/creating-an-ods-task.adoc)
* [Releasing a new version](/docs/releasing.adoc)
