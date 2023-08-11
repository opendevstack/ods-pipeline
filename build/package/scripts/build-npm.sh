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
  mkdir -p "${tmp_artifacts_dir}/lint-reports"
  cp eslint-report.txt "${tmp_artifacts_dir}/lint-reports/${ARTIFACT_PREFIX}report.txt"
}

# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"
BUILD_DIR="dist"
OUTPUT_DIR="docker"
WORKING_DIR="."
ARTIFACT_PREFIX=""
DEBUG="${DEBUG:-false}"
COPY_NODE_MODULES="false"

while [ "$#" -gt 0 ]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

    --debug) DEBUG="$2"; shift;;
    --debug=*) DEBUG="${1#*=}";;

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
tmp_artifacts_dir="${ROOT_DIR}/.ods/tmp-artifacts"
# tmp_artifacts_dir enables keeping artifacts created by this build 
# separate from other builds in the same repo to facilitate caching.
rm -rf "${tmp_artifacts_dir}"
if [ "${WORKING_DIR}" != "." ]; then
  cd "${WORKING_DIR}"
  ARTIFACT_PREFIX="${WORKING_DIR/\//-}-"
fi

echo "Configuring npm to use Nexus (${NEXUS_URL}) ..."
# Remove the protocol segment from NEXUS_URL
NEXUS_HOST=$(echo "${NEXUS_URL}" | sed -E 's/^\s*.*:\/\///g')
if [ -n "${NEXUS_URL}" ] && [ -n "${NEXUS_USERNAME}" ] && [ -n "${NEXUS_PASSWORD}" ]; then
    NEXUS_AUTH="$(urlencode "${NEXUS_USERNAME}"):$(urlencode "${NEXUS_PASSWORD}")"
    npm config set registry="$NEXUS_URL"/repository/npmjs/
    npm config set "//${NEXUS_HOST}/repository/npmjs/:_auth"="$(echo -n "$NEXUS_AUTH" | base64)"
    npm config set email=no-reply@opendevstack.org
    if [ -f /etc/ssl/certs/private-cert.pem ]; then
      echo "Configuring private cert ..."
      npm config set cafile=/etc/ssl/certs/private-cert.pem
    fi
fi;

echo "package-*.json checks ..."
if [ ! -f package.json ]; then
  echo "File package.json not found"
  exit 1
fi
if [ ! -f package-lock.json ]; then
  echo "File package-lock.json not found"
  exit 1
fi

echo "Installing dependencies ..."
npm ci --ignore-scripts

echo "Linting ..."
set +e
npm run lint > eslint-report.txt
exitcode=$?
set -e

if [ $exitcode == 0 ]; then
  echo "OK" > eslint-report.txt
  copyLintReport
else
  copyLintReport
  exit $exitcode
fi

echo "Building ..."
npm run build
# Copying most build output before testing so
# that additional modules which may be installed by testing
# is not included.
# However copying package.json too early can confuse the tests.
mkdir -p "${OUTPUT_DIR}"
echo "Copying contents of ${BUILD_DIR} into ${OUTPUT_DIR}/dist ..."
# see https://unix.stackexchange.com/questions/228597/how-to-copy-a-folder-recursively-in-an-idempotent-way-using-cp
"$CP" -r "${BUILD_DIR}/." "${OUTPUT_DIR}/dist"

if [ "${COPY_NODE_MODULES}" = true ]; then
  echo "Copying node_modules to ${OUTPUT_DIR}/dist/node_modules ..."
  # note "${OUTPUT_DIR}/dist" exists now and node_modules name will be maintained.
  "$CP"  -r node_modules "${OUTPUT_DIR}/dist"
fi

echo "Testing ..."
npm run test

mkdir -p "${tmp_artifacts_dir}/xunit-reports"
cp build/test-results/test/report.xml "${tmp_artifacts_dir}/xunit-reports/${ARTIFACT_PREFIX}report.xml"

mkdir -p "${tmp_artifacts_dir}/code-coverage"
cp build/coverage/clover.xml "${tmp_artifacts_dir}/code-coverage/${ARTIFACT_PREFIX}clover.xml"

cp build/coverage/coverage-final.json "${tmp_artifacts_dir}/code-coverage/${ARTIFACT_PREFIX}coverage-final.json"

cp build/coverage/lcov.info "${tmp_artifacts_dir}/code-coverage/${ARTIFACT_PREFIX}lcov.info"

# Doing this earlier can confuse jest.
# test build_javascript_app_with_custom_build_directory fails with
#  jest-haste-map: Haste module naming collision: src
#    The following files share their name; please adjust your hasteImpl:
#      * <rootDir>/package.json
#      * <rootDir>/docker/package.json
#  No tests found, exiting with code 1
# While one could demand to change the config of the test, there is no need
# to copy this earlier
echo "Copying package.json and package-lock.json to ${OUTPUT_DIR}/dist ..."
cp package.json package-lock.json "${OUTPUT_DIR}/dist"
