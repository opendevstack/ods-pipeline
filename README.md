# ODS Pipeline

[![Tekton Tasks Tests](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yml/badge.svg)](https://github.com/opendevstack/ods-pipeline/actions/workflows/main.yml)

**IMPORTANT: This is EXPERIMENTAL ONLY. This may or may not become part of ODS one day.**

## Introduction

ODS provides CI/CD pipeline support based on OpenShift Pipelines. This repository contains everything that relates to it, such as Tekton tasks, container images, Go packages, services, documentation, ...

The [ODS Pipeline Introduction](/docs/introduction.adoc) describes what ODS pipeline is and how it works. It is important to understand this before looking at further documentation or any other repository content.

## Documentation

The documentation provided by ODS pipeline has three audiences:

* **Admins** install and maintain a central ODS pipeline installation in an OpenShift cluster that can be used by many users.

* **Users** consume an existing central ODS pipeline installation. Users install ODS pipeline resources into a namespace they own to run CI/CD pipelines for their repositories.

* **Developers** contribute to the ODS pipeline project itself, for example by improving existing tasks, adding new ones, updating documentation, etc.

**Admin Guide**
* [Installation](/docs/admin-installation.adoc)

**User Guide**
* [Installation](/docs/user-installation.adoc)
* [Task Reference](/docs/tasks)
* [Authoring Tasks](/docs/authoring-tasks.adoc)

**Developer Guide**
* [Repository Layout](/docs/repository-layout.adoc)
* [Development & Running Tests](/docs/development.adoc)
* [Artifacts](/docs/artifacts.adoc)
* [Creating an ODS task](/docs/creating-an-ods-task.adoc)

This repository also hosts the design documents that describe ODS pipeline more formally. Those design documents provide more detail and background on goals, requirements and architecture decisions and may be useful for all audiences.

* [Goals and Non-Goals](/docs/goals-and-nongoals.adoc)
* [Architecture Decision Records](/docs/adr)
* [Task Requirements](/docs/task-requirements.adoc)

## Compatibility

 For OpenShift Pipelines releases and its relationship to Tekton and OpenShift versions, see https://docs.openshift.com/container-platform/4.8/cicd/pipelines/op-release-notes.html

 | ods-pipeline | OpenShift Pipelines | ODS Core/Quickstarters |
 |---|---|---|
 | 0.1 (to be released end of August 2021) | 1.5 | 4.0.0 |
 | 0.2 (to be released TBD) | 1.5 | 4.0.0 |
