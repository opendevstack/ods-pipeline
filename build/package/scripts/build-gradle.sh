#!/bin/bash
set -eu

OUTPUT_DIR="docker"
WORKING_DIR="."
ROOT_DIR=$(pwd)
export ARTIFACTS_DIR=$ROOT_DIR/.ods/artifacts
ARTIFACT_PREFIX=
DEBUG="${DEBUG:-false}"
GRADLE_ADDITIONAL_TASKS=
GRADLE_OPTIONS=

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

    --gradle-additional-tasks) GRADLE_ADDITIONAL_TASKS="$2"; shift;;
    --gradle-additional-tasks=*) GRADLE_ADDITIONAL_TASKS="${1#*=}";;

    # Gradle project properties ref: https://docs.gradle.org/7.4.2/userguide/build_environment.html#sec:gradle_configuration_properties
    # Gradle options ref: https://docs.gradle.org/7.4.2/userguide/command_line_interface.html
    --gradle-options) GRADLE_OPTIONS="$2"; shift;;
    --gradle-options=*) GRADLE_OPTIONS="${1#*=}";;

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
echo "Using GRADLE_OPTS=$GRADLE_OPTS"
echo "Using GRADLE_USER_HOME=$GRADLE_USER_HOME"
echo "Using ARTIFACTS_DIR=$ARTIFACTS_DIR"
mkdir -p "${GRADLE_USER_HOME}"

configure-gradle

echo
cd "${WORKING_DIR}"
echo "Working on Gradle project in '${WORKING_DIR}'..."
echo
echo "Gradlew version: "
./gradlew -version
echo
echo "Note on build environment variables available:"
echo
echo " ODS_OUTPUT_DIR: this environment variable points to the folder "
echo " that this build expects generated application artifacts to be copied to."
echo " The project gradle script should read this env var to copy all the "
echo " generated application artifacts."
echo
export ODS_OUTPUT_DIR=${OUTPUT_DIR}
echo "Exported env var 'ODS_OUTPUT_DIR' with value '${OUTPUT_DIR}'"
echo
echo "Building (Compile and Test) ..."
# shellcheck disable=SC2086
./gradlew clean build ${GRADLE_ADDITIONAL_TASKS} ${GRADLE_OPTIONS}
echo

echo "Verifying unit test report was generated  ..."
BUILD_DIR="build"
UNIT_TEST_RESULT_DIR="${BUILD_DIR}/test-results/test"

if [ -d "${UNIT_TEST_RESULT_DIR}" ]; then
    UNIT_TEST_ARTIFACTS_DIR="${ARTIFACTS_DIR}/xunit-reports"
    mkdir -p "${UNIT_TEST_ARTIFACTS_DIR}"
    # Each test class produces its own report file, but they contain a fully qualified class
    # name in their file name. Due to that, we do not need to add an artifact prefix to
    # distinguish them with reports from other artifacts of the same repo/pipeline build.
    cp "${UNIT_TEST_RESULT_DIR}/"*.xml "${UNIT_TEST_ARTIFACTS_DIR}"
else
  echo "Build failed: no unit test results found in ${UNIT_TEST_RESULT_DIR}"
  exit 1
fi

echo "Verifying unit test coverage report was generated  ..."
COVERAGE_RESULT_DIR="${BUILD_DIR}/reports/jacoco/test"
if [ -d "${COVERAGE_RESULT_DIR}" ]; then
    CODE_COVERAGE_ARTIFACTS_DIR="${ARTIFACTS_DIR}/code-coverage"
    mkdir -p "${CODE_COVERAGE_ARTIFACTS_DIR}"
    cp "${COVERAGE_RESULT_DIR}/jacocoTestReport.xml" "${CODE_COVERAGE_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}coverage.xml"
else
  echo "Build failed: no unit test coverage report was found in ${COVERAGE_RESULT_DIR}"
  exit 1
fi
