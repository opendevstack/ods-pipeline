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

OUTPUT_DIR="docker"
WORKING_DIR="."
ARTIFACT_PREFIX=""
DEBUG="false"

while [[ "$#" -gt 0 ]]; do
    case $1 in

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

ROOT_DIR=$(pwd)
if [ "${WORKING_DIR}" != "." ]; then
  cd "${WORKING_DIR}"
  ARTIFACT_PREFIX="${WORKING_DIR/\//-}-"
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

echo "Building ..."
npm ci
npm run build
mkdir -p "${OUTPUT_DIR}/dist"
cp -r dist "${OUTPUT_DIR}/dist"

echo "Copying node_modules to ${OUTPUT_DIR}/dist/node_modules ..."
cp -r node_modules "${OUTPUT_DIR}/dist/node_modules"

echo "Testing ..."
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
