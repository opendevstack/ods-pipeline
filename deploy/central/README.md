# deploy/central

This folder contains resources deployed centrally in the cluster. The resources
are defined using Helm.

Central resources are:
* `BuildConfig` and `ImageStream` resources  (in folder `images-chart`)
* `ClusterTask` resources (in folder `tasks-chart`)

The resources under `images-chart` are only applicable for OpenShift clusters.

## Versioning

In a KinD cluster there are no versions. Images use the implicit `latest` tag. That makes testing and local development easy.

In OpenShift, however, images and tasks are versioned. That provides the greatest stability.

Remember to adjust the `values.yaml` files every time there is a new version.
