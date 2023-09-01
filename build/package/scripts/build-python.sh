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

MAX_LINE_LENGTH="120"
WORKING_DIR="."
ARTIFACT_PREFIX=""
PRE_TEST_SCRIPT=""
DEBUG="${DEBUG:-false}"

while [ "$#" -gt 0 ]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --max-line-length) MAX_LINE_LENGTH="$2"; shift;;
    --max-line-length=*) MAX_LINE_LENGTH="${1#*=}";;

    --pre-test-script) PRE_TEST_SCRIPT="$2"; shift;;
    --pre-test-script=*) PRE_TEST_SCRIPT="${1#*=}";;

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

echo "Configuring pip to use Nexus (${NEXUS_URL}) ..."
# Remove the protocol segment from NEXUS_URL
NEXUS_HOST=$(echo "${NEXUS_URL}" | sed -E 's/^\s*.*:\/\///g')
if [ -n "${NEXUS_HOST}" ] && [ -n "${NEXUS_USERNAME}" ] && [ -n "${NEXUS_PASSWORD}" ]; then
    NEXUS_AUTH="$(urlencode "${NEXUS_USERNAME}"):$(urlencode "${NEXUS_PASSWORD}")"
    NEXUS_URL_WITH_AUTH="$(echo "${NEXUS_URL}" | sed -E 's/:\/\//:\/\/'"${NEXUS_AUTH}"@'/g')"
    pip3 config set global.index-url "${NEXUS_URL_WITH_AUTH}"/repository/pypi-all/simple
    pip3 config set global.trusted-host "${NEXUS_HOST}"
    pip3 config set global.extra-index-url https://pypi.org/simple
fi;

echo "Installing test requirements ..."
# shellcheck source=/dev/null
pip install --upgrade pip
pip install -r tests_requirements.txt
pip check

echo "Linting ..."
mypy src
flake8 --max-line-length="${MAX_LINE_LENGTH}" src

if [ -n "${PRE_TEST_SCRIPT}" ]; then
  echo "Executing pre-test script ..."
  ./"${PRE_TEST_SCRIPT}"
fi

echo "Testing ..."
rm report.xml coverage.xml &>/dev/null || true
PYTHONPATH=src python -m pytest --junitxml=report.xml -o junit_family=xunit2 --cov-report term-missing --cov-report xml:coverage.xml --cov=src -o testpaths=tests

mkdir -p "${tmp_artifacts_dir}/xunit-reports"
cp report.xml "${tmp_artifacts_dir}/xunit-reports/${ARTIFACT_PREFIX}report.xml"
mkdir -p "${tmp_artifacts_dir}/code-coverage"
cp coverage.xml "${tmp_artifacts_dir}/code-coverage/${ARTIFACT_PREFIX}coverage.xml"
