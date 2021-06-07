#!/bin/bash
# ODS images can be pre-pulled to reduce the start-up time of containers when a TaskRun is executed.
# For this, we use kube-fledged to pre-pull ODS images as part of the local setup, so they're already available in the KinD cluster.
# The advantage of this approach is to reduce the time the Golang tests take when testing ODS ClusterTasks and Tasks that make use of ODS images.
# See: https://github.com/senthilrch/kube-fledged
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}
KUBEFLEDGED_VERSION="v0.7.1"
KUBE_FLEDGED_FOLDER="${ODS_PIPELINE_DIR}/scripts/kube-fledged"

echo $KUBE_FLEDGED_FOLDER

# Install kube-fledged
if [ ! -d "$KUBE_FLEDGED_FOLDER" ] ; then
    git clone --depth 1 https://github.com/senthilrch/kube-fledged.git --branch ${KUBEFLEDGED_VERSION}
    cd kube-fledged/ && make deploy-using-yaml
fi

# Apply ImageCache Resource with the list of ODS images to pre-pull.
# Edit it as per your needs before creating image cache.
# If images are in private repositories requiring credentials to pull, add "imagePullSecrets" to the end.
kubectl apply -f "${ODS_PIPELINE_DIR}/scripts/deploy/kubefledged-imagecache-ods-images.yaml"