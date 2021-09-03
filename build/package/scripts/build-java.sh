#!/bin/bash
set -eu

OUTPUT_DIR="docker"
WORKING_DIR="."
ROOT_DIR=$(pwd)
ARTIFACTS_DIR=$ROOT_DIR/.ods/artifacts
DEBUG="false"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

    --gradle-additional-tasks) GRADLE_ADDITIONAL_TASKS="$2"; shift;;
    --gradle-additional-tasks=*) GRADLE_ADDITIONAL_TASKS="${1#*=}";;

    # Gradle options ref: https://docs.gradle.org/current/userguide/command_line_interface.html
    --gradle-options) GRADLE_OPTIONS="$2"; shift;;
    --gradle-options=*) GRADLE_OPTIONS="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${DEBUG}" == "true" ]; then
  set -x
fi

if [ -z "${WORKING_DIR}" ]; then
  WORKING_DIR="${ROOT_DIR}"
fi
cd "${WORKING_DIR}"

echo "Working on Java/Gradle module in ${WORKING_DIR} ..."

gradlew -version

# NOTE regarding nexus: no need to pass NEXUS env vars to gradle as gradle options. The should be configured to read those envs vars
echo "Note regarding nexus: your gradle script can configure nexus by reading the env vars: NEXUS_HOST, NEXUS_USER and , NEXUS_PASSWORD"

echo "Building (Compile and Test) ..."
graldlew clean build "${GRADLE_ADDITIONAL_TASKS}" "${GRADLE_OPTIONS}" "-PodsPipelineOutputDir=${OUTPUT_DIR}"

echo "Verifying unit test report was generated  ..."
BUILD_DIR="${WORKING_DIR}/build"

UNIT_TEST_RESULT_DIR="${BUILD_DIR}/tests-results"
if [ "$(ls -A ${UNIT_TEST_RESULT_DIR})" ]; then
    UNIT_TEST_ARTIFACTS_DIR="${ARTIFACTS_DIR}/xunit-reports"
    mkdir -p "${UNIT_TEST_ARTIFACTS_DIR}"
    cp "${UNIT_TEST_RESULT_DIR}/${ARTIFACT_PREFIX}*.xml" "${UNIT_TEST_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}*.xml"
else
  echo "Build failed: no unit test results found in ${UNIT_TEST_RESULT_DIR}"
  exit 1
fi

echo "Verifying unit test coverage report was generated  ..."
COVERAGE_RESULT_DIR="${BUILD_DIR}/jacoco/reports/test"
if [ "$(ls -A ${COVERAGE_RESULT_DIR})" ]; then
    CODE_COVERAGE_ARTIFACTS_DIR="${ARTIFACTS_DIR}/code-coverage"
    mkdir -p "${CODE_COVERAGE_ARTIFACTS_DIR}"
    cp "${COVERAGE_RESULT_DIR}/${ARTIFACT_PREFIX}*.xml" "${CODE_COVERAGE_ARTIFACTS_DIR}/${ARTIFACT_PREFIX}*.xml"
else
  echo "Build failed: no unit test coverage report was found in ${COVERAGE_RESULT_DIR}"
  exit 1
fi


