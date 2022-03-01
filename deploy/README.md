# Deployment

This directory contains container orchestration deployment configurations and templates.

Manifests in `ods-pipeline` are applied once per project by a project administrator.

## Tasks & Images Subcharts

The task subchart is maintained in https://github.com/opendevstack/ods-pipeline, and may be used by ODS pipeline admins to control the deployment of ODS pipeline resources in the respecitve project namespace in OpenShift.

### Subcharts Contents

The resources are defined using Helm:
* `BuildConfig` and `ImageStream` resources (in the `images` subchart)
* `Task` resources (in `tasks` subchart)

The resources of the `images` subchart are only applicable for OpenShift clusters. The subcharts may individually be enabled or disabled via the umbrella chart's `values.yaml`.

### Versioning

In a KinD cluster there are no versions. Images use the implicit `latest` tag. That makes testing and local development easy.

In OpenShift, however, images and tasks are versioned. That provides the greatest stability.

Remember to adjust the `values.yaml` files every time there is a new version.

## Setup Subchart

The setup subchart maintained in https://github.com/opendevstack/ods-pipeline, and may be used by ODS pipeline users to control the deployment of ODS pipeline resources in their project(s) in OpenShift.

## Subchart Contents

The resources are defined using Helm:
* `ConfigMap` and `Secret` resources used by ODS tasks
* ODS pipeline manager (`Service`/`Deployment`)
