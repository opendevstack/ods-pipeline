#!/bin/bash
set -eu

KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"
TKN_VERSION="v0.23.0"
TASKRUN_NAME="hello-world-example"

if ! which kubectl &> /dev/null; then
    echo "kubectl is required"
fi

if ! which tkn &> /dev/null; then
    echo "tkn is required"
fi

if ! which jq &> /dev/null; then
    echo "jq is required"
fi

# Install Tekton
# https://tekton.dev/docs/getting-started/#installation
if ! $KUBECTL_BIN get namespace tekton-pipelines &> /dev/null; then
    echo "Installing Tekton ..."
    $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/${TKN_VERSION}/release.notags.yaml
else
    echo "Tekton already installed."
fi

echo "Install Tekton Task ..."
$KUBECTL_BIN apply -f task.yaml

echo "Run Task ..."
$KUBECTL_BIN delete taskrun/${TASKRUN_NAME} || true
$KUBECTL_BIN apply -f taskrun.yaml
STATUS=$($KUBECTL_BIN get taskrun ${TASKRUN_NAME} -o json | jq -rc .status.conditions[0].status)
LIMIT=$((SECONDS+180))
while [ "${STATUS}" != "True" ]; do
    if [ $SECONDS -gt $LIMIT ]; then
    echo "Timeout waiting for TaskRun to complete."
    exit 2
    fi
    sleep 10
    echo "Waiting for TaskRun to complete ..."
    STATUS=$($KUBECTL_BIN get taskrun ${TASKRUN_NAME} -o json | jq -rc .status.conditions[0].status)
done

echo "Get TaskRun Logs ..."
tkn taskrun logs ${TASKRUN_NAME} -a -f

echo "Getting result ..."
REASON=$($KUBECTL_BIN get taskrun ${TASKRUN_NAME} -o json | jq -rc .status.conditions[0].reason)
echo "The job ${REASON}"
test "${REASON}" != "Failed"
