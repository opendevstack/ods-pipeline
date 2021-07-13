# deploy/central

This folder contains resources deployed centrally in the cluster. The resources
are defined using kustomize.

Central resources are:
* `ClusterTask` resources (in folder `tasks`)
* `BuildConfig` and `ImageStream` resources  (in folder `images`)

The resources under `images` are only applicable for OpenShift clusters.

## Versioning

By default (which is used in a KinD cluster) there are no versions. Images use
the implicit `latest` tag. That makes testing and local development easy.

In OpenShift, however, images and tasks are versioned. That provides the
greatest stability.

Remember that the patches in the `openshift` need to be adjusted every time
there is a new version.
