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

WORKING_DIR="."
ARTIFACT_PREFIX=""
DEBUG="${DEBUG:-false}"

while [ "$#" -gt 0 ]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --debug) DEBUG="$2"; shift;;
    --debug=*) DEBUG="${1#*=}";;

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

echo "Testing ..."
npm run test

mkdir -p "${tmp_artifacts_dir}/xunit-reports"
cp build/test-results/test/report.xml "${tmp_artifacts_dir}/xunit-reports/${ARTIFACT_PREFIX}report.xml"

mkdir -p "${tmp_artifacts_dir}/code-coverage"
cp build/coverage/clover.xml "${tmp_artifacts_dir}/code-coverage/${ARTIFACT_PREFIX}clover.xml"

cp build/coverage/coverage-final.json "${tmp_artifacts_dir}/code-coverage/${ARTIFACT_PREFIX}coverage-final.json"

cp build/coverage/lcov.info "${tmp_artifacts_dir}/code-coverage/${ARTIFACT_PREFIX}lcov.info"
