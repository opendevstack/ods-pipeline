#!/bin/bash
set -eu

KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"

# Tekton versions are aligned with Red Hat OpenShift Pipelines General Availability 1.4.
# See https://docs.openshift.com/container-platform/4.7/cicd/pipelines/op-release-notes.html.
TKN_VERSION="v0.22.0"
TKN_DASHBOARD_VERSION="v0.17.0"
TKN_TRIGGERS="v0.12.0"

INSTALL_TKN_DASHBOARD="false"

if ! which kubectl &> /dev/null; then
    echo "kubectl is required"
fi

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    --tekton-dashboard) INSTALL_TKN_DASHBOARD="true";;

    # -h|--help) usage; exit 0;;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

# Install Tekton
# https://tekton.dev/docs/getting-started/#installation
if ! $KUBECTL_BIN get namespace tekton-pipelines &> /dev/null; then
    echo "Installing Tekton ..."
    $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/${TKN_VERSION}/release.notags.yaml
    $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/triggers/previous/${TKN_TRIGGERS}/release.yaml
    if [ "${INSTALL_TKN_DASHBOARD}" != "false" ]; then
        echo "Installing Tekton Dashboard..."
        $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/dashboard/previous/${TKN_DASHBOARD_VERSION}/tekton-dashboard-release.yaml
    fi
else
    echo "Tekton already installed."
fi
