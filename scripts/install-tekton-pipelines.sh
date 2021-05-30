#!/bin/bash
set -eu

KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"
TKN_VERSION="v0.22.0"
NAMESPACE="default"

if ! which kubectl &> /dev/null; then
    echo "kubectl is required"
fi

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    # -h|--help) usage; exit 0;;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

# Install Tekton
# https://tekton.dev/docs/getting-started/#installation
if ! $KUBECTL_BIN get namespace tekton-pipelines &> /dev/null; then
    echo "Installing Tekton ..."
    $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/${TKN_VERSION}/release.notags.yaml
else
    echo "Tekton already installed."
fi
