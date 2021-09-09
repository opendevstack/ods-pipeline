#!/bin/bash
set -eu

OUTPUT_DIR="docker"
WORKING_DIR="."
ROOT_DIR=$(pwd)
ARTIFACTS_DIR=$ROOT_DIR/.ods/artifacts
ARTIFACT_PREFIX=
DEBUG="false"
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

    # Gradle project properties ref: https://docs.gradle.org/current/userguide/build_environment.html#sec:gradle_configuration_properties
    # Gradle options ref: https://docs.gradle.org/current/userguide/command_line_interface.html
    --gradle-options) GRADLE_OPTIONS="$2"; shift;;
    --gradle-options=*) GRADLE_OPTIONS="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ -z "${WORKING_DIR}" ]; then
  WORKING_DIR="${ROOT_DIR}"
fi

if [ "${DEBUG}" == "true" ]; then
  set -x
  echo
  echo "Debug mode enabled (DEBUG=true)"
  echo "Generating debug logs ..."
  echo "... current folder:"
  pwd
  echo "... folder content:"
  ls -lart
  echo
fi

cd "${WORKING_DIR}"
echo "Working on Java/Gradle project in ${WORKING_DIR} ..."
echo
echo "... gladlew version: "
./gradlew -version

echo
echo "Note on build environment variables available:"
echo
echo " ODS_OUTPUT_DIR: this environment variable points to the folder "
echo " that this build expects generated application artifacts to be copied to."
echo " The project gradle script should read this env var to copy all the "
echo " generated application artifacts."
echo
echo " NEXUS_* env vars:"
echo " following env vars NEXUS_HOST, NEXUS_USER and NEXUS_PASSWORD"
echo " are available and should be read by your gradle script."
echo
export ODS_OUTPUT_DIR=${OUTPUT_DIR}
echo "Exported env var 'ODS_OUTPUT_DIR' with value '${OUTPUT_DIR}'"
echo
echo "Building (Compile and Test) ..."
./gradlew clean build ${GRADLE_ADDITIONAL_TASKS} ${GRADLE_OPTIONS}
echo

if [ "${DEBUG}" == "true" ]; then
  set -x
  echo
  echo "List content of ${OUTPUT_DIR}"
  ls -lart "${OUTPUT_DIR}"
  echo
fi

echo "Verifying unit test report was generated  ..."
BUILD_DIR="build"
UNIT_TEST_RESULT_DIR="${BUILD_DIR}/test-results/test"

if [ "$(ls -A ${UNIT_TEST_RESULT_DIR})" ]; then
    UNIT_TEST_ARTIFACTS_DIR="${ARTIFACTS_DIR}/xunit-reports"
    mkdir -p "${UNIT_TEST_ARTIFACTS_DIR}"
    echo "... copy unit test report: from ${UNIT_TEST_RESULT_DIR}/*.xml to ${UNIT_TEST_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}report.xml"
    cp "${UNIT_TEST_RESULT_DIR}/"*.xml "${UNIT_TEST_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}report.xml"
    echo "... copied unit test report!"
    echo
else
  echo "Build failed: no unit test results found in ${UNIT_TEST_RESULT_DIR}"
  exit 1
fi

echo "Verifying unit test coverage report was generated  ..."
COVERAGE_RESULT_DIR="${BUILD_DIR}/reports/jacoco/test"
if [ "$(ls -A ${COVERAGE_RESULT_DIR})" ]; then
    CODE_COVERAGE_ARTIFACTS_DIR="${ARTIFACTS_DIR}/code-coverage"
    mkdir -p "${CODE_COVERAGE_ARTIFACTS_DIR}"
    echo "... copy unit test coverage report: from ${COVERAGE_RESULT_DIR}/jacocoTestReport.xml to ${CODE_COVERAGE_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}coverage.xml"
    cp "${COVERAGE_RESULT_DIR}/jacocoTestReport.xml" "${CODE_COVERAGE_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}coverage.xml"
    echo "... copied unit test coverage report!"
    echo
else
  echo "Build failed: no unit test coverage report was found in ${COVERAGE_RESULT_DIR}"
  exit 1
fi

echo "... working on Java/Gradle project is done!"
