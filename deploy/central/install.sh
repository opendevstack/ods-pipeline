#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

VERBOSE="false"
DRY_RUN="false"
DIFF="true"
NAMESPACE=""
RELEASE_NAME=""
VALUES_FILE=""
CHART_DIR=""
CHART=""

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) VERBOSE="true";;

    -n|--namespace) NAMESPACE="$2"; shift;;
    -n=*|--namespace=*) NAMESPACE="${1#*=}";;

    -f|--values) VALUES_FILE="$2"; shift;;
    -f=*|--values=*) VALUES_FILE="${1#*=}";;

    -c|--chart) CHART="$2"; shift;;
    -c=*|--chart=*) CHART="${1#*=}";;

    --no-diff) DIFF="false";;

    --dry-run) DRY_RUN="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd ${SCRIPT_DIR}

VALUES_FILES=$(echo $VALUES_FILE | tr "," "\n")
VALUES_ARGS=""
for valueFile in ${VALUES_FILES}; do
    VALUES_ARGS="${VALUES_ARGS} --values=${valueFile}"
done

if [ -z "${CHART}" ]; then
    echo "--chart is required"
    exit 1
elif [ "${CHART}" == "tasks" ]; then
    RELEASE_NAME="ods-pipeline-tasks"
    CHART_DIR="./tasks-chart"
elif [ "${CHART}" == "images" ]; then
    RELEASE_NAME="ods-pipeline-images"
    CHART_DIR="./images-chart"
    if [ -z "${NAMESPACE}" ]; then
        echo "--namespace is required"
        exit 1
    fi
else
    echo "--chart is not valid. Use 'tasks' or 'images'."
    exit 1
fi



if [ "${VERBOSE}" == "true" ]; then
    set -x
fi

echo "Installing Helm release ${RELEASE_NAME} ..."
if [ "${DIFF}" == "true" ]; then
    if helm -n ${NAMESPACE} \
            diff upgrade --install --detailed-exitcode \
            ${VALUES_ARGS} \
            ${RELEASE_NAME} ${CHART_DIR}; then
        echo "Helm release already up-to-date."
    else
        if [ "${DRY_RUN}" == "true" ]; then
            echo "(skipping in dry-run)"
        else
            helm -n ${NAMESPACE} \
                upgrade --install \
                ${VALUES_ARGS} \
                ${RELEASE_NAME} ${CHART_DIR}
        fi
    fi
else
    if [ "${DRY_RUN}" == "true" ]; then
        echo "(skipping in dry-run)"
    else
        NAMESPACE_FLAG=""
        if [ -n "${NAMESPACE}" ]; then
            NAMESPACE_FLAG="-n ${NAMESPACE}"
        fi
        helm ${NAMESPACE_FLAG} \
            upgrade --install \
            ${VALUES_ARGS} \
            ${RELEASE_NAME} ${CHART_DIR}
    fi
fi
