#!/bin/bash
set -eu
# Copy build to cache


# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"

OUTPUT_DIR="docker"
WORKING_DIR="."
CACHE_BUILD_KEY=
CACHE_LOCATION_USED_PATH=
DEBUG="${DEBUG:-false}"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

    --cache-build-key) CACHE_BUILD_KEY="$2"; shift;;
    --cache-build-key=*) CACHE_BUILD_KEY="${1#*=}";;

    --cache-location-used-path) CACHE_LOCATION_USED_PATH="$2"; shift;;
    --cache-location-used-path=*) CACHE_LOCATION_USED_PATH="${1#*=}";;

    --debug) DEBUG="$2"; shift;;
    --debug=*) DEBUG="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ -z "${CACHE_BUILD_KEY}" ]; then
  echo "Param --cache-build-key is required."; exit 1;
elif [ -z "${CACHE_LOCATION_USED_PATH}" ]; then
  echo "Param --cache-location-used-path is required."; exit 1;
fi

if [ "${DEBUG}" == "true" ]; then
  set -x
fi

ROOT_DIR=$(pwd)

if [ "${WORKING_DIR}" == "." ]; then
  git_sha_working_dir=$(git rev-parse "HEAD:")
else
  git_sha_working_dir=$(git rev-parse "HEAD:$WORKING_DIR")
fi
prior_output_dir="$ROOT_DIR/.ods-cache/build-task/$CACHE_BUILD_KEY/$git_sha_working_dir"

if [ "${WORKING_DIR}" != "." ]; then
  cd "${WORKING_DIR}"
fi

rm -rvf "$prior_output_dir"  # should be empty as otherwise cache should be used.
mkdir -p "$prior_output_dir"

# Copying reports
cache_of_reports_dir="$prior_output_dir/reports"
ods_artifacts_dir="${ROOT_DIR}/.ods/artifacts"
echo "Copying build reports to cache: $ods_artifacts_dir -> $cache_of_reports_dir"
mkdir -p "$cache_of_reports_dir"
$CP -r --link "$ods_artifacts_dir/." "$cache_of_reports_dir"

# Copying build output
cache_of_output_dir="$prior_output_dir/output"
echo "Copying build output to cache: $OUTPUT_DIR to $cache_of_output_dir"
mkdir -p "$cache_of_output_dir"
$CP -r --link "$OUTPUT_DIR/." "$cache_of_output_dir"

echo "$prior_output_dir" > "$CACHE_LOCATION_USED_PATH"
touch "$prior_output_dir/.ods-last-used-stamp"
