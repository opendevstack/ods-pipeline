# ODS Pipeline

[![Tekton Tasks Tests](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yaml/badge.svg)](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yaml)

**IMPORTANT: This approach is currently being evaluated. This may or may not become part of ODS one day.**

## Introduction

ODS Pipeline provides a  CI/CD pipeline based on OpenShift Pipelines as an alternative to Jenkins. This repository contains everything that relates to it, such as Tekton tasks, container images, Go packages, services, documentation, ...

The [ODS Pipeline Introduction](/docs/introduction.adoc) describes what ODS pipeline is and how it works. It is important to understand this before looking at further documentation or any other repository content.

Note that this ODS Pipeline repository does not provide tasks for building, packaging and deploying your application. You can use any Tekton task to fullfil those needs, however, there are a few "companion" tasks specifically designed for ODS Pipeline to cover common use cases. See the "Technical Reference" section below for information on those tasks.

## Documentation

### Getting Started
* [Installation & Updating](/docs/installation.adoc)
* [Add ODS Pipeline to a repository](/docs/add-to-repository.adoc)
* [Convert an ODS quickstarter based component](/docs/convert-quickstarter-component.adoc)

### Technical Reference
* [Repository configuration (ods.yaml)](/docs/ods-configuration.adoc)
* [Start task](/docs/task-start.adoc)
* [Finish task](/docs/task-finish.adoc)
* Companion tasks: [Go build task](https://github.com/opendevstack/ods-pipeline-go), [Gradle build task](https://github.com/opendevstack/ods-pipeline-gradle), [NPM build task](https://github.com/opendevstack/ods-pipeline-npm), [SonarQube scan task](https://github.com/opendevstack/ods-pipeline-sonar), [Package image task](https://github.com/opendevstack/ods-pipeline-image), [Helm deploy task](https://github.com/opendevstack/ods-pipeline-helm)

### How-To Guides
* [Accessing artifacts](/docs/accessing-artifacts.adoc)
* [Debugging](/docs/debugging.adoc)
* [Authoring your own tasks](/docs/authoring-tasks.adoc)
* [FAQ](https://github.com/opendevstack/ods-pipeline/wiki/FAQ)

### Examples
* [Example Project](/docs/example-project.adoc)
* [BIX-Digital Repository Templates](https://github.com/BIX-Digital/ods-pipeline-examples)

### Background Information
* [Stakeholder Requirements](/docs/design/stakeholder-requirements.adoc)
* [Software Architecture](/docs/design/software-architecture.adoc)
* [Software Requirements Specification](/docs/design/software-requirements-specification.adoc)
* [Software Design Specification](/docs/design/software-design-specification.adoc)
* [Architecture Decision Records](/docs/adr)
* [Goals and Non-Goals](/docs/design/goals-and-nongoals.adoc)
* [Relationship to Jenkins Shared Library](/docs/design/relationship-shared-library.adoc)

## Compatibility

For OpenShift Pipelines releases and its relationship to Tekton and OpenShift versions, see https://docs.openshift.com/container-platform/latest/cicd/pipelines/op-release-notes.html

| ods-pipeline | OpenShift Pipelines | ODS Core/Quickstarters |
|---|---|---|
| [0.13](https://github.com/opendevstack/ods-pipeline/releases/tag/v0.13.2) | 1.9 | 4.x |
| [0.12](https://github.com/opendevstack/ods-pipeline/releases/tag/v0.12.0) | 1.9 | 4.x |
| [0.11](https://github.com/opendevstack/ods-pipeline/releases/tag/v0.11.1) | 1.9 | 4.x |

## Contributing

* [Repository Layout](/docs/repository-layout.adoc)
* [Development & Running Tests](/docs/development.adoc)
* [Artifacts](/docs/artifacts.adoc)
* [Releasing a new version](/docs/releasing.adoc)
