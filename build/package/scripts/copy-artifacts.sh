#!/bin/bash
set -eu

# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"

DEBUG="${DEBUG:-false}"

while [[ "$#" -gt 0 ]]; do
    case $1 in
    --debug) DEBUG="$2"; shift;;
    --debug=*) DEBUG="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${DEBUG}" == "true" ]; then
  set -x
fi

ROOT_DIR=$(pwd)
ODS_ARTIFACTS_DIR="${ROOT_DIR}/.ods/artifacts"
TMP_ARTIFACTS_DIR="${ROOT_DIR}/.ods/tmp-artifacts"

# Copying ods artifacts which are mostly reports (see artifacts.adoc)
echo "Copying build artifacts from $TMP_ARTIFACTS_DIR to $ODS_ARTIFACTS_DIR"
mkdir -p "$ODS_ARTIFACTS_DIR"
# "$CP" -v -r --link "$TMP_ARTIFACTS_DIR/." "$ODS_ARTIFACTS_DIR"
"$CP" -v -r "$TMP_ARTIFACTS_DIR/." "$ODS_ARTIFACTS_DIR"
