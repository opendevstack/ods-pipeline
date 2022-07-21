#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

PERSONAL_DIR=$1
shift

if [ -z "$PERSONAL_DIR" ]; then
  echo "Missing <personal-directory> as first argument. The personal-directory contains your personal ods-pipeline deployment values (values.yaml, secrets.yaml) and the install.sh scripts."
  exit 1
fi
set -ux




# Delegate to install.sh of the project checked out cd repository typically at <project>-cd/deploy.
cd "$PERSONAL_DIR"
./install.sh -f "$PERSONAL_DIR/values.yaml,$PERSONAL_DIR/secrets.yaml,$ODS_PIPELINE_DIR/deploy/ods-pipeline/values.kind-personal.yaml" "$@"
