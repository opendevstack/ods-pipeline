#!/bin/bash
set -eu
# if build is cached skip it and return exit 0
# otherwise exit with an error code.


# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"
LS="${GNU_LS:-ls}"

OUTPUT_DIR="docker"
CACHE_BUILD="true"
CACHE_BUILD_KEY=
CACHE_LOCATION_USED_PATH=
WORKING_DIR="."
DEBUG="${DEBUG:-false}"

while [ "$#" -gt 0 ]; do
    case $1 in

    --cache-build) CACHE_BUILD="$2"; shift;;
    --cache-build=*) CACHE_BUILD="${1#*=}";;

    --cache-build-key) CACHE_BUILD_KEY="$2"; shift;;
    --cache-build-key=*) CACHE_BUILD_KEY="${1#*=}";;

    --cache-location-used-path) CACHE_LOCATION_USED_PATH="$2"; shift;;
    --cache-location-used-path=*) CACHE_LOCATION_USED_PATH="${1#*=}";;

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

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

if [ "$CACHE_BUILD" != "true" ]; then
  echo "Build skipping is not enabled. Continuing with a regular build (CACHE_BUILD==$CACHE_BUILD)"
  exit 1
fi

ROOT_DIR=$(pwd)

if [ "${WORKING_DIR}" == "." ]; then
  git_sha_working_dir=$(git rev-parse "HEAD:")
else
  git_sha_working_dir=$(git rev-parse "HEAD:$WORKING_DIR")
fi
cache_location_dir="$ROOT_DIR/.ods-cache/build-task/$CACHE_BUILD_KEY/$git_sha_working_dir"
if [ ! -d "$cache_location_dir" ]; then
  echo "No prior build output found in cache at $cache_location_dir"
  exit 1  # subsequent build log makes clear what happens next
fi

if [ "${WORKING_DIR}" != "." ]; then
  cd "${WORKING_DIR}"
fi

# Copying ods artifacts which are mostly reports (see artifacts.adoc)
cache_of_artifacts_dir="$cache_location_dir/artifacts"
ods_artifacts_dir="${ROOT_DIR}/.ods/artifacts"
echo "Copying prior build artifacts from cache: $cache_of_artifacts_dir to $ods_artifacts_dir"
mkdir -p "$ods_artifacts_dir"
"$CP" -v -r "$cache_of_artifacts_dir/." "$ods_artifacts_dir"

# Copying build output
cache_of_output_dir="$cache_location_dir/output"
if [ "${DEBUG}" == "true" ]; then
  echo "-- ls OUTPUT IN CACHE -- "
  $LS -Ral "$cache_of_output_dir"
  echo "-- ls OUTPUT_DIR -- "
  $LS -Ral "$OUTPUT_DIR"
fi
echo "Copying prior build output from cache: $cache_of_output_dir to $OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"
start_time=$SECONDS
"$CP" $CP_VERBOSITY_FLAGS -r "$cache_of_output_dir/." "$OUTPUT_DIR"
elapsed=$(( SECONDS - start_time ))
echo "Copying took $elapsed seconds"

echo "$cache_location_dir" > "$CACHE_LOCATION_USED_PATH"
touch "$cache_location_dir/.ods-last-used-stamp"
