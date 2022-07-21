#!/bin/bash
set -eu

KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"

# Tekton version is aligned with Red Hat OpenShift Pipelines General Availability 1.6.
# See https://docs.openshift.com/container-platform/4.9/cicd/pipelines/op-release-notes.html.
TKN_VERSION="v0.28.3"
TKN_DASHBOARD_VERSION="v0.17.0"

INSTALL_TKN_DASHBOARD="false"

if ! which kubectl &> /dev/null; then
    echo "kubectl is required"
fi

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    --tekton-dashboard) INSTALL_TKN_DASHBOARD="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

# Install Tekton
# https://tekton.dev/docs/getting-started/#installation
if ! $KUBECTL_BIN get namespace tekton-pipelines &> /dev/null; then
    echo "Installing Tekton ..."
    $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/${TKN_VERSION}/release.notags.yaml

    if [ "${INSTALL_TKN_DASHBOARD}" != "false" ]; then
        echo "Installing Tekton Dashboard..."
        $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/dashboard/previous/${TKN_DASHBOARD_VERSION}/tekton-dashboard-release.yaml
    fi
else
    echo "Tekton already installed."
fi
