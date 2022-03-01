#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

VERBOSE="false"
DRY_RUN="false"
DIFF="true"
NAMESPACE=""
RELEASE_NAME="ods-pipeline"
SERVICEACCOUNT="pipeline"
VALUES_FILE="values.custom.yaml"
CHART_DIR="./ods-pipeline"

while [[ "$#" -gt 0 ]]; do
    # shellcheck disable=SC2034
    case $1 in

    -v|--verbose) VERBOSE="true";;

    -n|--namespace) NAMESPACE="$2"; shift;;
    -n=*|--namespace=*) NAMESPACE="${1#*=}";;

    -f|--values) VALUES_FILE="$2"; shift;;
    -f=*|--values=*) VALUES_FILE="${1#*=}";;

    -s|--serviceaccount) SERVICEACCOUNT="$2"; shift;;
    -s=*|--serviceaccount=*) SERVICEACCOUNT="${1#*=}";;

    --no-diff) DIFF="false";;

    --dry-run) DRY_RUN="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd "${SCRIPT_DIR}"

VALUES_FILES=$(echo "$VALUES_FILE" | tr "," "\n")
VALUES_ARGS=()
for valueFile in ${VALUES_FILES}; do
    VALUES_ARGS+=(--values="${valueFile}")
done

if [ "${VERBOSE}" == "true" ]; then
    set -x
fi

if [ -z "${NAMESPACE}" ]; then
    echo "--namespace is required"
    exit 1
fi

if kubectl -n "${NAMESPACE}" get serviceaccount/"${SERVICEACCOUNT}" &> /dev/null; then
    echo "Serviceaccount exists already ..."
else
    echo "Creating serviceaccount ..."
    if [ "${DRY_RUN}" == "true" ]; then
        echo "(skipping in dry-run)"
    else
        kubectl -n "${NAMESPACE}" create serviceaccount "${SERVICEACCOUNT}"

        kubectl -n "${NAMESPACE}" \
            create rolebinding "${SERVICEACCOUNT}-edit" \
            --clusterrole edit \
            --serviceaccount "${NAMESPACE}:${SERVICEACCOUNT}"
    fi
fi

DIFF_UPGRADE_ARGS=(diff upgrade)
UPGRADE_ARGS=(upgrade)
if helm plugin list | grep secrets &> /dev/null; then
    DIFF_UPGRADE_ARGS=(secrets diff upgrade)
    UPGRADE_ARGS=(secrets upgrade)
fi

echo "Installing Helm release ${RELEASE_NAME} ..."
if [ "${DIFF}" == "true" ]; then
    if helm -n "${NAMESPACE}" \
            "${DIFF_UPGRADE_ARGS[@]}" --install --detailed-exitcode \
            "${VALUES_ARGS[@]}" \
            ${RELEASE_NAME} ${CHART_DIR}; then
        echo "Helm release already up-to-date."
    else
        if [ "${DRY_RUN}" == "true" ]; then
            echo "(skipping in dry-run)"
        else
            helm -n "${NAMESPACE}" \
                "${UPGRADE_ARGS[@]}" --install \
                "${VALUES_ARGS[@]}" \
                ${RELEASE_NAME} ${CHART_DIR}
        fi
    fi
else
    if [ "${DRY_RUN}" == "true" ]; then
        echo "(skipping in dry-run)"
    else
        helm -n "${NAMESPACE}" \
            "${UPGRADE_ARGS[@]}" --install \
            "${VALUES_ARGS[@]}" \
            ${RELEASE_NAME} ${CHART_DIR}
    fi
fi

echo "Adding ods-bitbucket-auth secret to serviceaccount ..."
if [ "${DRY_RUN}" == "true" ]; then
    echo "(skipping in dry-run)"
else
    kubectl -n "${NAMESPACE}" \
        patch sa "${SERVICEACCOUNT}" \
        --type json \
        -p '[{"op": "add", "path": "/secrets", "value":[{"name": "ods-bitbucket-auth"}]}]'
fi
