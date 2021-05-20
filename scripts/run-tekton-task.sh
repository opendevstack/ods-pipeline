#!/bin/bash
set -eu

KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"
TKN_VERSION="v0.22.0"
TASKRUN_NAME="hello-world-example"
NAMESPACE="default"
TASK_FILE=""
TASKRUN_FILE=""

if ! which kubectl &> /dev/null; then
    echo "kubectl is required"
fi

if ! which tkn &> /dev/null; then
    echo "tkn is required"
fi

if ! which jq &> /dev/null; then
    echo "jq is required"
fi

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    # -h|--help) usage; exit 0;;

    -t|--task-file) TASK_FILE="$2"; shift;;
    -t=*|--task-file=*) TASK_FILE="${1#*=}";;

    -r|--task-run-file) TASKRUN_FILE="$2"; shift;;
    -r=*|--task-run-file=*) TASKRUN_FILE="${1#*=}";;

    --task-run-name) TASKRUN_NAME="$2"; shift;;
    --task-run-name=*) TASKRUN_NAME="${1#*=}";;

    -n|--namespace) NAMESPACE="$2"; shift;;
    -n=*|--namespace=*) NAMESPACE="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ -z "${TASK_FILE}" ]; then
  echo "--task-file is required"
  exit 1
fi

if [ -z "${TASKRUN_FILE}" ]; then
  echo "--task-run-file is required"
  exit 1
fi

if [ -z "${TASKRUN_NAME}" ]; then
  echo "--task-run-name is required"
  exit 1
fi

KUBECTL_BIN_WITH_NS="$KUBECTL_BIN -n $NAMESPACE"

# Install Tekton
# https://tekton.dev/docs/getting-started/#installation
if ! $KUBECTL_BIN get namespace tekton-pipelines &> /dev/null; then
    echo "Installing Tekton ..."
    $KUBECTL_BIN apply --filename https://storage.googleapis.com/tekton-releases/pipeline/previous/${TKN_VERSION}/release.notags.yaml
else
    echo "Tekton already installed."
fi

echo "Install Tekton Task ..."
$KUBECTL_BIN_WITH_NS apply -f ${TASK_FILE}

echo "Run Task ..."
$KUBECTL_BIN_WITH_NS delete taskrun/${TASKRUN_NAME} || true
$KUBECTL_BIN_WITH_NS apply -f ${TASKRUN_FILE}
STATUS=$($KUBECTL_BIN_WITH_NS get taskrun ${TASKRUN_NAME} -o json | jq -rc .status.conditions[0].status)
LIMIT=$((SECONDS+180))
while [ "${STATUS}" != "True" ]; do
    if [ $SECONDS -gt $LIMIT ]; then
    echo "Timeout waiting for TaskRun to complete."
    exit 2
    fi
    sleep 10
    echo "Waiting for TaskRun to complete ..."
    STATUS=$($KUBECTL_BIN_WITH_NS get taskrun ${TASKRUN_NAME} -o json | jq -rc .status.conditions[0].status)
done

echo "Get TaskRun Logs ..."
tkn taskrun logs ${TASKRUN_NAME} -a -f

echo "Getting result ..."
REASON=$($KUBECTL_BIN_WITH_NS get taskrun ${TASKRUN_NAME} -o json | jq -rc .status.conditions[0].reason)
echo "The job ${REASON}"
test "${REASON}" != "Failed"
