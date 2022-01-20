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

BUILD_DIR="dist"
OUTPUT_DIR="docker"
WORKING_DIR="."
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


if [ -d "$PIPELINE_CACHE_DIR" ]; then
  # NOTE: rsync is currently not available in build script
  echo "Updating cache dir with sources ..."
  # rsync -auvhp --itemize-changes --delete --exclude=node_modules "${WORKING_DIR}" "$PIPELINE_CACHE_DIR"
  # rsync -auvhp --progress --delete --exclude=node_modules "${WORKING_DIR}" "$PIPELINE_CACHE_DIR"
  rsync -auhp --info=progress2 --delete --exclude=node_modules "${WORKING_DIR}" "$PIPELINE_CACHE_DIR"  # the info progress2 does not work with verbose, needs rather new rsync version.
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

build_dir="${WORKING_DIR}"
output_dir_abs="${ROOT_DIR}/${WORKING_DIR}/${OUTPUT_DIR}"
if [ -d  "$PIPELINE_CACHE_DIR" ]; then
  build_dir="$PIPELINE_CACHE_DIR/$WORKING_DIR"
fi

if [ "${WORKING_DIR}" != "." ]; then
  ARTIFACT_PREFIX="${WORKING_DIR/\//-}-"
fi
if [ "${build_dir}" != "." ]; then
  cd "${build_dir}" 
fi

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
mkdir -p "${output_dir_abs}"
cp -rv "${BUILD_DIR}" "${output_dir_abs}/dist"

if [ "${COPY_NODE_MODULES}" = true ]; then
  echo "Copying node_modules to ${output_dir_abs}/dist/node_modules ..."
  start_time=$SECONDS
  cp -r node_modules "${output_dir_abs}/dist/node_modules"
  elapsed=$(( SECONDS - start_time ))
  echo "copying node_modules took $elapsed seconds"
fi

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
fi

# Provide disk usage so that one can use this to estimate needs


supply-sonar-project-properties-default
