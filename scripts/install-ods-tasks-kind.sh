#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

# Delegate to install.sh within deploy/central.
# The script here exists only for consistency (all scripts are located under /scripts)
# and the install.sh script is under deploy/central so that the whole
# deployment is self-contained within that folder, making it easy for consumers
# to pull in the deployment logic into their repositories via "git subtree".
"${ODS_PIPELINE_DIR}"/deploy/central/install.sh \
    --namespace=default \
    --chart=tasks \
    --no-diff \
    --values=./tasks-chart/values.kind.yaml "$@"
