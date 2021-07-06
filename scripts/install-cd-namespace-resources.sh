#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

VERBOSE="false"
DIFF="true"
NAMESPACE=""
RELEASE_NAME="ods-pipeline"
SERVICEACCOUNT="pipeline"
VALUES_FILE="values.custom.yaml"
CHART_DIR="./chart"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) VERBOSE="true";;

    -n|--namespace) NAMESPACE="$2"; shift;;
    -n=*|--namespace=*) NAMESPACE="${1#*=}";;

    -f|--values) VALUES_FILE="$2"; shift;;
    -f=*|--values=*) VALUES_FILE="${1#*=}";;

    -s|--serviceaccount) SERVICEACCOUNT="$2"; shift;;
    -s=*|--serviceaccount=*) SERVICEACCOUNT="${1#*=}";;

    --no-diff) DIFF="false"; shift;;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd ${ODS_PIPELINE_DIR}/deploy/cd-namespace

VALUES_FILES=$(echo $VALUES_FILE | tr "," "\n")
VALUES_ARGS=""
for valueFile in ${VALUES_FILES}; do
    VALUES_ARGS="${VALUES_ARGS} --values=${CHART_DIR}/${valueFile}"
done

if [ "${VERBOSE}" == "true" ]; then
    set -x
fi

if kubectl -n ${NAMESPACE} get serviceaccount/${SERVICEACCOUNT} &> /dev/null; then
    echo "Serviceaccount exists already ..."
else
    echo "Creating serviceaccount ..."
    kubectl -n ${NAMESPACE} create serviceaccount ${SERVICEACCOUNT}

    kubectl -n ${NAMESPACE} \
        create rolebinding "${SERVICEACCOUNT}-edit" \
        --clusterrole edit \
        --serviceaccount "${NAMESPACE}:${SERVICEACCOUNT}"
fi

echo "Installing Helm release ..."
if [ "${DIFF}" == "true" ]; then
    if helm -n ${NAMESPACE} \
            diff upgrade --install --detailed-exitcode \
            ${VALUES_ARGS} \
            ${RELEASE_NAME} ${CHART_DIR}; then
        echo "Helm release already up-to-date."
    else
        helm -n ${NAMESPACE} \
            upgrade --install \
            ${VALUES_ARGS} \
            ${RELEASE_NAME} ${CHART_DIR}
    fi
else
    helm -n ${NAMESPACE} \
            upgrade --install \
            ${VALUES_ARGS} \
            ${RELEASE_NAME} ${CHART_DIR}
fi

echo "Adding ods-bitbucket-auth secret to serviceaccount ..."
kubectl -n ${NAMESPACE} \
    patch sa ${SERVICEACCOUNT} \
    --type json \
    -p '[{"op": "add", "path": "/secrets", "value":[{"name": "ods-bitbucket-auth"}]}]'

# echo "Exposing event listener ..."
# oc -n ${NAMESPACE} expose svc el-ods-pipeline
