#!/bin/bash
set -eu
# Copy build to cache


# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"
LS="${GNU_LS:-ls}"

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

CP_VERBOSITY_FLAGS=
if [ "${DEBUG}" == "true" ]; then
  set -x
  CP_VERBOSITY_FLAGS="-v"
fi

ROOT_DIR=$(pwd)

git_sha_working_dir=""
if [ "${WORKING_DIR}" == "." ]; then
  git_sha_working_dir=$(git rev-parse "HEAD:")
else
  git_sha_working_dir=$(git rev-parse "HEAD:$WORKING_DIR")
fi
cache_location_dir="$ROOT_DIR/.ods-cache/build-task/$CACHE_BUILD_KEY/$git_sha_working_dir"

if [ "${WORKING_DIR}" != "." ]; then
  cd "${WORKING_DIR}"
fi

rm -rvf "$cache_location_dir"  # should be empty as otherwise cache should be used.
mkdir -p "$cache_location_dir"

# Copying ods artifacts which are mostly reports (see artifacts.adoc)
# TODO: consistent casing and naming across scripts regarding dir variables
cache_of_artifacts_dir="$cache_location_dir/artifacts"
tmp_artifacts_dir="${ROOT_DIR}/.ods/tmp-artifacts"
echo "Copying build artifacts to cache: $tmp_artifacts_dir -> $cache_of_artifacts_dir"
mkdir -p "$cache_of_artifacts_dir"
"$CP" -v -r "$tmp_artifacts_dir/." "$cache_of_artifacts_dir"

# Copying build output
cache_of_output_dir="$cache_location_dir/output"
echo "Copying build output to cache: $OUTPUT_DIR to $cache_of_output_dir"
mkdir -p "$cache_of_output_dir"
start_time=$SECONDS
"$CP" $CP_VERBOSITY_FLAGS -r "$OUTPUT_DIR/." "$cache_of_output_dir"
elapsed=$(( SECONDS - start_time ))
echo "Copying took $elapsed seconds"
if [ "${DEBUG}" == "true" ]; then
  echo "-- ls OUTPUT IN CACHE -- "
  $LS -Ral "$cache_of_output_dir"
fi

echo "$cache_location_dir" > "$CACHE_LOCATION_USED_PATH"
touch "$cache_location_dir/.ods-last-used-stamp"
