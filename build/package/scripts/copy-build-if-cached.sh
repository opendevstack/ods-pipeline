#!/usr/bin/env bash  
# avoids ancient bash on macos
set -eu
# if build is cached skip it and return exit 0
# otherwise exit with an error code.


# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"
LS="${GNU_LS:-ls}"

join() {
  local IFS="$1"
  shift
  echo "$*"
}

splitAtColon() { 
  # colon is 
  echo $1 | tr ":" "\n"
}
# https://stackoverflow.com/a/918931

outputs_str=
extra_inputs_str=
declare -a outputs
outputs=()
declare -a inputs
inputs=()
working_dir="."
cache_build="true"
cache_build_key=
cache_location_used_path=
debug="${DEBUG:-false}"
dry_run=false

while [ "$#" -gt 0 ]; do
    case $1 in

    --working-dir) working_dir="$2"; shift;;
    --working-dir=*) working_dir="${1#*=}";;

    --cached-outputs) outputs_str="$2"; shift;;
    --cached-outputs=*) outputs_str="${1#*=}";;

    --cache-extra-inputs) extra_inputs_str="$2"; shift;;
    --cache-extra-inputs=*) extra_inputs_str="${1#*=}";;

    --cache-build) cache_build="$2"; shift;;
    --cache-build=*) cache_build="${1#*=}";;

    --cache-build-key) cache_build_key="$2"; shift;;
    --cache-build-key=*) cache_build_key="${1#*=}";;

    --cache-location-used-path) cache_location_used_path="$2"; shift;;
    --cache-location-used-path=*) cache_location_used_path="${1#*=}";;

    --debug) debug="$2"; shift;;
    --debug=*) debug="${1#*=}";;

    --dry-run) dry_run=true;;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ -z "${cache_build_key}" ]; then
  echo "Param --cache-build-key is required."; exit 1;
elif [ -z "${cache_location_used_path}" ]; then
  echo "Param --cache-location-used-path is required."; exit 1;
fi

cp_verbosity_flags=
if [ "${debug}" == "true" ]; then
  set -x
  cp_verbosity_flags="-v"
fi

IFS=":" read -r -a outputs <<< "$outputs_str"

# note leads to undefined variable if extra_inputs_str is empty on ancient bash
IFS=":" read -r -a extra_inputs <<< "$extra_inputs_str"
inputs=("$working_dir")
for f in "${extra_inputs[@]}"; do
  inputs+=( "$f" )
done

if [ "$cache_build" != "true" ]; then
  echo "Build skipping is not enabled. Continuing with a regular build (cache_build==$cache_build)"
  exit 1
fi

root_dir=$(pwd)

declare -a git_shas  #  relative to root
for f in "${inputs[@]}"; do
  if [ "${f}" == "." ]; then
    git_shas+=( "$(git rev-parse --short "HEAD:")" )
  else
    git_shas+=( "$(git rev-parse --short "HEAD:$f")")
  fi
done
# shellcheck disable=SC2048,SC2086
git_sha_combined=$(join "-" ${git_shas[*]})
cache_location_dir="$root_dir/.ods-cache/build-task/$cache_build_key/$git_sha_combined"

if [ ! -d "$cache_location_dir" ]; then
  echo "No prior build output found in cache at $cache_location_dir"
  exit 1  # subsequent build log makes clear what happens next
fi

if [ "${working_dir}" != "." ]; then
  cd "${working_dir}"
fi

# Copying ods artifacts which are mostly reports (see artifacts.adoc)
cache_of_artifacts_dir="$cache_location_dir/artifacts"
ods_artifacts_dir="${root_dir}/.ods/artifacts"
echo "Copying prior ods build artifacts from cache: $cache_of_artifacts_dir to $ods_artifacts_dir"
if [ "${dry_run}" == "true" ]; then
  echo "(skipping copying ods build artifacts)"
else
  mkdir -p "$ods_artifacts_dir"
  "$CP" -v -r "$cache_of_artifacts_dir/." "$ods_artifacts_dir"
fi

# Copying build output
for i in "${!outputs[@]}"; do
  cache_of_output_dir="$cache_location_dir/output/$i"
  output_dir="${outputs[$i]}"
  echo "Copying prior build output from cache: $cache_of_output_dir to $output_dir"
  if [ "${debug}" == "true" ]; then
    echo "-- ls OUTPUT IN CACHE -- "
    $LS -Ral "$cache_of_output_dir"
    echo "-- ls OUTPUT_DIR -- "
    $LS -Ral "$output_dir"
  fi
  if [ "${dry_run}" == "true" ]; then
    echo "(skipping copying build outputs)"
  else 
    mkdir -p "$output_dir"
    start_time=$SECONDS
    "$CP" $cp_verbosity_flags -r "$cache_of_output_dir/." "$output_dir"
    elapsed=$(( SECONDS - start_time ))
    echo "Copying took $elapsed seconds"
  fi
done

if [ "${dry_run}" == "true" ]; then
  echo "(skipping saving $cache_location_dir in $cache_location_used_path)"
  echo "(skipping touch of $cache_location_dir/.ods-last-used-stamp"
else
  echo "$cache_location_dir" > "$cache_location_used_path"
  touch "$cache_location_dir/.ods-last-used-stamp"
fi
