# Deployment of "cd-namespace"

This directory is maintained in https://github.com/opendevstack/ods-pipeline, and may be used by ODS pipeline users to control the deployment of ODS pipeline resources in their "cd namespace" in OpenShift. For this purpose, this directory may be added to another Git repository via `git subtree` as explained in the [User Installation Guide](/docs/user-installation.adoc).

## Directory Contents

The resources are defined using Helm:
* `ConfigMap` and `Secret` resources used by ODS tasks
* Tekton Triggers related resources (`EventListener` etc.), including the custom ODS pipeline webhook interceptor
