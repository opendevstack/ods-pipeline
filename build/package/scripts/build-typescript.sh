#!/bin/bash
set -eu

urlencode() {
    local LC_COLLATE=C
    local length="${#1}"
    for (( i = 0; i < length; i++ )); do
        local c="${1:$i:1}"
        case $c in
            [a-zA-Z0-9.~_-]) printf '%s' "$c" ;;
            *) printf '%%%02X' "'$c" ;;
        esac
    done
}

copyLintReport() {
  cat eslint-report.txt
  mkdir -p "${ROOT_DIR}/.ods/artifacts/lint-reports"
  cp eslint-report.txt "${ROOT_DIR}/.ods/artifacts/lint-reports/${ARTIFACT_PREFIX}report.txt"
}

copyTestReports() {
  mkdir -p "${ROOT_DIR}/.ods/artifacts/xunit-reports"
  cat build/test-results/test/report.xml
  cp build/test-results/test/report.xml "${ROOT_DIR}/.ods/artifacts/xunit-reports/${ARTIFACT_PREFIX}report.xml"

  mkdir -p "${ROOT_DIR}/.ods/artifacts/code-coverage"
  cat build/coverage/clover.xml
  cp build/coverage/clover.xml "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}clover.xml"

  cat build/coverage/coverage-final.json
  cp build/coverage/coverage-final.json "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}coverage-final.json"

  cat build/coverage/lcov.info
  cp build/coverage/lcov.info "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}lcov.info"
}

copyBuildResults() {
  mkdir -p "${output_dir_abs}"
  cp -rv "${BUILD_DIR}" "${output_dir_abs}/dist"

  if [ "${COPY_NODE_MODULES}" = true ]; then
    echo "Copying node_modules to ${output_dir_abs}/dist/node_modules ..."
    start_time=$SECONDS
    cp -r node_modules "${output_dir_abs}/dist/node_modules"
    elapsed=$(( SECONDS - start_time ))
    echo "copying node_modules took $elapsed seconds"
  fi
}

touch_dir_commit_hash() {
  # this would enables cleanup based on these timestamps
  # however this function should ideally be done outside of the individual 
  # build scripts.
  mkdir -p "$build_root_dir/.ods/"
  echo -n "$WORKING_DIR_COMMIT_SHA" > "$build_root_dir/.ods/git-dir-commit-sha"
}

report_disk_usage() {
  echo "Disk usage to estimate caching needs."
  du -hs -- *
}


BUILD_DIR="dist"
OUTPUT_DIR="docker"
WORKING_DIR="."
WORKING_DIR_COMMIT_SHA=""
PIPELINE_CACHE_DIR=""
ARTIFACT_PREFIX=""
DEBUG="${DEBUG:-false}"
MAX_LINT_WARNINGS="0"
LINT_FILE_EXT=".js,.ts,.jsx,.tsx,.svelte"
COPY_NODE_MODULES="false"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --working-dir-commit-sha) WORKING_DIR_COMMIT_SHA="$2"; shift;;
    --working-dir-commit-sha=*) WORKING_DIR_COMMIT_SHA="${1#*=}";;

    # build pipeline should support skipping for example via commit tag
    # in this case --pipeline-cache-dir would be omitted or set to an empty value.  
    --pipeline-cache-dir) PIPELINE_CACHE_DIR="$2"; shift;;
    --pipeline-cache-dir=*) PIPELINE_CACHE_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

    --debug) DEBUG="$2"; shift;;
    --debug=*) DEBUG="${1#*=}";;

    --max-lint-warnings) MAX_LINT_WARNINGS="$2"; shift;;
    --max-lint-warnings=*) MAX_LINT_WARNINGS="${1#*=}";;

    --lint-file-ext) LINT_FILE_EXT="$2"; shift;;
    --lint-file-ext=*) LINT_FILE_EXT="${1#*=}";;

    --build-dir) BUILD_DIR="$2"; shift;;
    --build-dir=*) BUILD_DIR="${1#*=}";;

    --copy-node-modules) COPY_NODE_MODULES="$2"; shift;;
    --copy-node-modules=*) COPY_NODE_MODULES="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ -z "$WORKING_DIR_COMMIT_SHA" ]; then
  echo "--working-dir-commit-sha parameter is required."; exit 1
fi 

if [ "${DEBUG}" == "true" ]; then
  set -x
fi

ROOT_DIR=$(pwd)
# NOTE: git is currently not available in build script
git_ref=$(cat .ods/git-ref)
case $git_ref in
  master|develop)
    echo "INFO: not using cache in master or develop branch"
    PIPELINE_CACHE_DIR=""
    ;;
  *);;
esac

same_build_in_cache=false
build_root_dir="${WORKING_DIR}"
output_dir_abs="${ROOT_DIR}/${WORKING_DIR}/${OUTPUT_DIR}"
if [ -n "$PIPELINE_CACHE_DIR" ]; then
  build_root_dir="$PIPELINE_CACHE_DIR/$WORKING_DIR"
  if [ ! -d "$build_root_dir" ]; then
    mkdir -p "$build_root_dir"
  fi
  if [ -f "$build_root_dir/.ods/git-dir-commit-sha" ]; then
    previous_dir_commit_sha=$(cat "$build_root_dir/.ods/git-dir-commit-sha")
    if [ "$previous_dir_commit_sha" = "$WORKING_DIR_COMMIT_SHA" ]; then
      same_build_in_cache=true
      echo "INFO: build with same commit hash of dir $build_root_dir already in cache ($previous_dir_commit_sha)."
    fi
  fi
  if [ "$same_build_in_cache" = "false" ]; then
    # NOTE: rsync is currently not available in build script
    echo "Updating cache dir with sources ..."
    rsync -ah -v --delete --exclude=node_modules "${WORKING_DIR}" "$PIPELINE_CACHE_DIR"
    # -a = -rlptgoD
    #   -r = --recursive
    #   -l = --links recreate the symlink on the destination
    #   -p = --perms (details see man page): aims to set dest permissions same as source
    #   -t = --times tells rsync to transfer modification times 
    #   -goD = --group --owner --devices --special
    # -h = --hard-links reestablish hard links if they are both in the copied set  
    # -v = --verbose (multiple would be more verbose)
    # copied content is shown by sending incremental file list. Which maybe empty. 
    # seems like this is all that is needed 
  fi
fi

if [ "${WORKING_DIR}" != "." ]; then
  ARTIFACT_PREFIX="${WORKING_DIR/\//-}-"
fi
if [ "${build_root_dir}" != "." ]; then
  cd "${build_root_dir}" 
fi

if [ "$same_build_in_cache" = "true" ]; then
  echo "Using prior build for same commit hash to copy build results and reports"
  copyBuildResults
  copyLintReport
  copyTestReports
  report_disk_usage
  exit 0
fi
echo "Configuring npm to use Nexus ..."
# Remove the protocol segment from NEXUS_URL
NEXUS_HOST=$(echo "${NEXUS_URL}" | sed -E 's/^\s*.*:\/\///g')
if [ -n "${NEXUS_HOST}" ] && [ -n "${NEXUS_USERNAME}" ] && [ -n "${NEXUS_PASSWORD}" ]; then
    NEXUS_AUTH="$(urlencode "${NEXUS_USERNAME}"):$(urlencode "${NEXUS_PASSWORD}")"
    npm config set registry="$NEXUS_URL"/repository/npmjs/
    npm config set always-auth=true
    npm config set _auth="$(echo -n "$NEXUS_AUTH" | base64)"
    npm config set email=no-reply@opendevstack.org
    npm config set ca=null
    npm config set strict-ssl=false
fi;

echo "Installing dependencies ..."
<<<<<<< HEAD
npm ci --ignore-scripts
=======
start_time=$SECONDS
if [ -d "$PIPELINE_CACHE_DIR" ]; then
  npm i
else
  npm ci
fi
elapsed=$(( SECONDS - start_time ))
echo "Installing dependencies took $elapsed seconds"
>>>>>>> 4dc6da42 (build-typescript draft changes)

echo "Linting ..."
start_time=$SECONDS
set +e
npx eslint src --ext "${LINT_FILE_EXT}" --format compact --max-warnings "${MAX_LINT_WARNINGS}" > eslint-report.txt
exitcode=$?
set -e
elapsed=$(( SECONDS - start_time ))
echo "linting took $elapsed seconds"

if [ $exitcode == 0 ]; then
  echo "OK" > eslint-report.txt
  copyLintReport
else
  copyLintReport
  exit $exitcode
fi

echo "Building ..."
start_time=$SECONDS
npm run build
elapsed=$(( SECONDS - start_time ))
echo "build took $elapsed seconds"
copyBuildResults

echo "Testing ..."
# Implement skipping tests for TypeScript #238
# TODO: is this still needed?
# https://github.com/opendevstack/ods-pipeline/issues/238
if [ -f "${ROOT_DIR}/.ods/artifacts/xunit-reports/${ARTIFACT_PREFIX}report.xml" ]; then
  echo "Test artifacts already present, skipping tests ..."
  # Copy artifacts to working directory so that the SonarQube scanner can pick them up later.
  cp "${ROOT_DIR}/.ods/artifacts/xunit-reports/${ARTIFACT_PREFIX}report.xml" report.xml
  cp "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}clover.xml" clover.xml
  cp "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}coverage-final.json" coverage-final.json
  cp "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}lcov.info" lcov.info
else
  npm run test
  copyTestReports
  touch_dir_commit_hash
fi

report_disk_usage

supply-sonar-project-properties-default
