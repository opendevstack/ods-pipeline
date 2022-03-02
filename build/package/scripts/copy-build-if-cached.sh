#!/bin/bash
set -eu
# if build is cached skip it and return exit 0
# otherwise exit with an error code.


# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"

OUTPUT_DIR="docker"
CACHE_OUTPUT_DIR="true"
WORKING_DIR="."
DEBUG="${DEBUG:-false}"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --cache-output-dir) CACHE_OUTPUT_DIR="$2"; shift;;
    --cache-output-dir=*) CACHE_OUTPUT_DIR="${1#*=}";;

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

    --debug) DEBUG="$2"; shift;;
    --debug=*) DEBUG="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${DEBUG}" == "true" ]; then
  set -x
fi

if [ "$CACHE_OUTPUT_DIR" != "true" ]; then
  echo "Build skipping is not enabled. Continuing with a regular build (CACHE_OUTPUT_DIR==$CACHE_OUTPUT_DIR)"
  exit 1
fi

ROOT_DIR=$(pwd)

git_sha_working_dir=$(git rev-parse "HEAD:$WORKING_DIR")
prior_output_dir="$ROOT_DIR/.ods-cache/build-task/$git_sha_working_dir"
if [ ! -d "$prior_output_dir" ]; then
  echo "No prior build output found in cache at $prior_output_dir"
  exit 1  # no message really needed here
fi

if [ "${WORKING_DIR}" != "." ]; then
  cd "${WORKING_DIR}"
fi

# TODO: consider ensure cache problems self repair?

# Copying reports
cache_of_reports_dir="$prior_output_dir/reports"
ods_artifacts_dir="${ROOT_DIR}/.ods/artifacts"
echo "Copying prior build reports from cache: $cache_of_reports_dir to $ods_artifacts_dir"
mkdir -p "$ods_artifacts_dir"
$CP -r --link "$cache_of_reports_dir/." "$ods_artifacts_dir"

# Copying build output
cache_of_output_dir="$prior_output_dir/output"
echo "Copying prior build output from cache: $cache_of_output_dir to $OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"
$CP -r --link "$cache_of_output_dir/." "$OUTPUT_DIR"
