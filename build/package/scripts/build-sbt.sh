#!/bin/bash
set -eu

function timestamped() {
	echo "$(date "+%Y/%m/%d %H:%M:%S") $1"
}

OUTPUT_DIR="docker"
WORKING_DIR="."
ROOT_DIR=$(pwd)
export ARTIFACTS_DIR=$ROOT_DIR/.ods/artifacts

# might be needed for several task executions that would publish to the same artefact path...
ARTIFACT_PREFIX=
DEBUG="${DEBUG:-false}"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${WORKING_DIR}" != "." ]; then
  WORKING_DIR="${ROOT_DIR}"
  ARTIFACT_PREFIX="${WORKING_DIR/\//-}-"
fi

if [ "${DEBUG}" == "true" ]; then
  set -x
fi

echo "Using NEXUS_URL=$NEXUS_URL"
echo "Using ARTIFACTS_DIR=$ARTIFACTS_DIR"

echo
cd "${WORKING_DIR}"
echo "Working on SBT project in '${WORKING_DIR}'..."
echo
export ODS_OUTPUT_DIR=${OUTPUT_DIR}
echo "Exported env var 'ODS_OUTPUT_DIR' with value '${OUTPUT_DIR}'"
echo
echo "Building (Compile and Test) ..."
# shellcheck disable=SC2086

# check format of sbt and source files, activate coverage and test with coverage report
timestamped "run tests and coverage"
UNIT_TEST_ARTIFACTS_DIR="${ARTIFACTS_DIR}/xunit-reports"
CODE_COVERAGE_ARTIFACTS_DIR="${ARTIFACTS_DIR}/code-coverage"
export UNIT_TEST_RESULT_DIR="${UNIT_TEST_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}"
export CODE_COVERAGE_TARGET_FILE="${CODE_COVERAGE_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}scoverage.xml"
sbt -no-colors -v clean scalafmtSbtCheck scalafmtCheckAll coverage test coverageReport copyOdsReports clean stage

timestamped "Verifying unit test report was generated ..."
if ls "${UNIT_TEST_RESULT_DIR}"*.xml >/dev/null 2>&1 ; then
	timestamped "unit test results exist under ${UNIT_TEST_RESULT_DIR}"
else
  timestamped "Build failed: no unit test results found in ${UNIT_TEST_RESULT_DIR}"
  exit 1
fi

timestamped "Verifying unit test coverage report was generated ..."
if [ -f "${CODE_COVERAGE_TARGET_FILE}" ]; then
	timestamped "unit test coverage report was found at ${CODE_COVERAGE_TARGET_FILE}"
else
  timestamped "Build failed: no unit test coverage report was found at ${CODE_COVERAGE_TARGET_FILE}"
  exit 1
fi

BUILD_DIR="target"
STAGING_DIR="${BUILD_DIR}/universal/stage"
timestamped "Copying contents of ${STAGING_DIR} to ${OUTPUT_DIR}/dist ..."
cp -r "${STAGING_DIR}/." "${OUTPUT_DIR}/dist"
