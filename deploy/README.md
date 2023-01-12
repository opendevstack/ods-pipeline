# Deployment

This directory contains container orchestration deployment configurations and templates.

Manifests in `ods-pipeline` are applied once per project by a project administrator.

## Subcharts

The `tasks` and `setup` subcharts are maintained in https://github.com/opendevstack/ods-pipeline, and may be used by project admins to control the deployment of ODS pipeline resources in the respective project namespace in OpenShift.

### Subcharts Contents

The resources are defined using Helm:
* `Task` resources (in `tasks` subchart)
* `ConfigMap` and `Secret` resources used by ODS tasks (in `setup` subchart)
* ODS pipeline manager (`Service`/`Deployment`) (in `setup` subchart)

### Versioning

In a KinD cluster there are no versions. Images use the implicit `latest` tag. That makes testing and local development easy.

In OpenShift, however, images and tasks are versioned. That provides the greatest stability.

Remember to adjust the `values.yaml` files every time there is a new version.
