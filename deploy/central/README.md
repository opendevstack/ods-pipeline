# Deployment of central "ODS namespace"

This directory is maintained in https://github.com/opendevstack/ods-pipeline, and may be used by ODS pipeline admins to control the deployment of ODS pipeline resources in the centrals "ODS namespace" in OpenShift. For this purpose, this directory may be added to another Git repository via `git subtree` as explained in the [Admin Installation Guide](/docs/admin-installation.adoc).

## Directory Contents

The resources are defined using Helm:
* `BuildConfig` and `ImageStream` resources (in folder `images-chart`)
* `Task` resources (in folder `tasks-chart`)

The resources under `images-chart` are only applicable for OpenShift clusters.

## Versioning

In a KinD cluster there are no versions. Images use the implicit `latest` tag. That makes testing and local development easy.

In OpenShift, however, images and tasks are versioned. That provides the greatest stability.

Remember to adjust the `values.yaml` files every time there is a new version.
