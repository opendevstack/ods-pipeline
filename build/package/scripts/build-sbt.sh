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
# old command: "sbt clean scalafmtSbtCheck scalafmtCheckAll coverage test coverageReport copyDockerFiles"
# plugins needed by the build: scalafmt scoverage native-packager
# TODO: with scoverage the build needs to run two times, one time with coverage and one time without, otherwise the production code
# will end up instrumented in production... This is due to the fact that scoverage has no way for on-the-fly instrumentation which only happens in memory

# check format of sbt and source files, activate coverage and test with coverage report
timestamped "run tests and coverage"
sbt -no-colors -v clean scalafmtSbtCheck scalafmtCheckAll coverage test coverageReport

# copy reports
timestamped "Verifying unit test report was generated ..."
BUILD_DIR="target"
UNIT_TEST_RESULT_DIR="${BUILD_DIR}/test-reports"
if [ -d "${UNIT_TEST_RESULT_DIR}" ]; then
    UNIT_TEST_ARTIFACTS_DIR="${ARTIFACTS_DIR}/xunit-reports"
    mkdir -p "${UNIT_TEST_ARTIFACTS_DIR}"
    cp "${UNIT_TEST_RESULT_DIR}/"*.xml "${UNIT_TEST_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}"
else
  echo "Build failed: no unit test results found in ${UNIT_TEST_RESULT_DIR}"
  exit 1
fi

timestamped "Verifying unit test coverage report was generated  ..."
COVERAGE_RESULT_DIR="${BUILD_DIR}/scala-2.13"
if [ -d "${COVERAGE_RESULT_DIR}" ]; then
    CODE_COVERAGE_ARTIFACTS_DIR="${ARTIFACTS_DIR}/code-coverage"
    mkdir -p "${CODE_COVERAGE_ARTIFACTS_DIR}"
    cp "${COVERAGE_RESULT_DIR}/scoverage-report/scoverage.xml" "${CODE_COVERAGE_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}scoverage.xml"
else
  echo "Build failed: no unit test coverage report was found in ${COVERAGE_RESULT_DIR}"
  exit 1
fi

# create a clean binary as the previous compiled sources where instrumented for the coverage report
timestamped "creating build artefacts"
sbt -no-colors -v clean stage

STAGING_DIR="${BUILD_DIR}/universal/stage"
timestamped "Copying contents of ${STAGING_DIR} to ${OUTPUT_DIR}/dist ..."
cp -r "${STAGING_DIR}/." "${OUTPUT_DIR}/dist"

# TODO oder alles in einem command:
# sbt clean coverage test coverageReport coverageOff clean compile / publishMyReports / copy stuff to docker...
# dann muss man nur dafür sorgen, dass die reports vor dem 2. clean gesaved werden...
