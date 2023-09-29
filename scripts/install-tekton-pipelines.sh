#!/bin/bash
set -eu

kube_context="--context kind-ods-pipeline"
kubectl_bin="kubectl $kube_context"

# Tekton version is aligned with Red Hat OpenShift Pipelines General Availability 1.10.
# See https://docs.openshift.com/container-platform/latest/cicd/pipelines/op-release-notes.html.
tkn_version="v0.44.4"
tkn_dashboard_version="v0.17.0"

install_tkn_dashboard="false"

if ! command -v kubectl &> /dev/null; then
    echo "kubectl is required"
fi

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    --tekton-dashboard) install_tkn_dashboard="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

# Install Tekton
# https://tekton.dev/docs/getting-started/#installation
if ! $kubectl_bin get namespace tekton-pipelines &> /dev/null; then
    echo "Installing Tekton ..."
    $kubectl_bin apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/${tkn_version}/release.notags.yaml

    if [ "${install_tkn_dashboard}" != "false" ]; then
        echo "Installing Tekton Dashboard..."
        $kubectl_bin apply --filename https://storage.googleapis.com/tekton-releases/dashboard/previous/${tkn_dashboard_version}/tekton-dashboard-release.yaml
    fi
else
    echo "Tekton already installed."
fi
