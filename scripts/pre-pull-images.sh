#!/bin/bash
# ODS images can be pre-pulled to reduce the start-up time of containers when a TaskRun is executed.
# For this, we use kube-fledged to pre-pull ODS images as part of the local setup, so they're already available in the KinD cluster.
# The advantage of this approach is to reduce the time the Golang tests take when testing ODS ClusterTasks and Tasks that make use of ODS images.
# See: https://github.com/senthilrch/kube-fledged
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}
KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"
KUBEFLEDGED_VERSION="v0.7.1"
KUBE_FLEDGED_FOLDER="${ODS_PIPELINE_DIR}/scripts/kube-fledged"

# Install kube-fledged
if ! $KUBECTL_BIN get namespace kube-fledged &> /dev/null; then

    echo "Installing kube-fledged ${KUBEFLEDGED_VERSION} ..."
    if [ ! -d "$KUBE_FLEDGED_FOLDER" ] ; then
        git clone --depth 1 https://github.com/senthilrch/kube-fledged.git --branch ${KUBEFLEDGED_VERSION}
    fi

    echo "kube-fledge repo already exists. Using local directory."
    cd kube-fledged/ && make deploy-using-yaml

    # Wait for kube-fledge deployments to be ready
    $KUBECTL_BIN rollout status deploy/kubefledged-controller -n kube-fledged
    $KUBECTL_BIN rollout status deploy/kubefledged-webhook-server -n kube-fledged
fi

echo "kube-fledge already installed."

# Apply ImageCache Resource with the list of ODS images to pre-pull.
# Edit it as per your needs before creating image cache.
# If images are in private repositories requiring credentials to pull, add "imagePullSecrets" to the end.
echo "Applying ImageCache resourece for ODS images..."
$KUBECTL_BIN apply -f "${ODS_PIPELINE_DIR}/scripts/deploy/kubefledged-imagecache-ods-images.yaml"