#!/bin/bash
set -eu

# the copy commands are based on GNU cp tools
# On a mac `brew install coreutils` gives `g` prefixed cmd line tools such as gcp
# to use these define env variable GNU_CP=gcp before invoking this script.
CP="${GNU_CP:-cp}"

output_dir="docker"
working_dir="."
artifact_prefix=""
debug="${DEBUG:-false}"
gradle_build_dir="build"
gradle_additional_tasks=
gradle_options=

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --working-dir) working_dir="$2"; shift;;
    --working-dir=*) working_dir="${1#*=}";;

    --output-dir) output_dir="$2"; shift;;
    --output-dir=*) output_dir="${1#*=}";;

    --gradle-build-dir) gradle_build_dir="$2"; shift;;
    --gradle-build-dir=*) gradle_build_dir="${1#*=}";;

    --gradle-additional-tasks) gradle_additional_tasks="$2"; shift;;
    --gradle-additional-tasks=*) gradle_additional_tasks="${1#*=}";;

    # Gradle project properties ref: https://docs.gradle.org/7.4.2/userguide/build_environment.html#sec:gradle_configuration_properties
    # Gradle options ref: https://docs.gradle.org/7.4.2/userguide/command_line_interface.html
    --gradle-options) gradle_options="$2"; shift;;
    --gradle-options=*) gradle_options="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

root_dir=$(pwd)
tmp_artifacts_dir="${root_dir}/.ods/tmp-artifacts"
# tmp_artifacts_dir enables keeping artifacts created by this build
# separate from other builds in the same repo to facilitate caching.
rm -rf "${tmp_artifacts_dir}"
if [ "${working_dir}" != "." ]; then
  cd "${working_dir}"
  artifact_prefix="${working_dir/\//-}-"
fi

if [ "${debug}" == "true" ]; then
  set -x
fi

echo "Using NEXUS_URL=$NEXUS_URL"
echo "Using GRADLE_OPTS=$GRADLE_OPTS"
echo "Using GRADLE_USER_HOME=$GRADLE_USER_HOME"
mkdir -p "${GRADLE_USER_HOME}"

configure-gradle

echo
echo "Working on Gradle project in '${working_dir}'..."
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
export ODS_OUTPUT_DIR=${output_dir}
echo "Exported env var 'ODS_OUTPUT_DIR' with value '${output_dir}'"
echo
echo "Building (Compile and Test) ..."
# shellcheck disable=SC2086
./gradlew clean build ${gradle_additional_tasks} ${gradle_options}
echo

echo "Verifying unit test report was generated  ..."
unit_test_result_dir="${gradle_build_dir}/test-results/test"
if [ -d "${unit_test_result_dir}" ]; then
    unit_test_artifacts_dir="${tmp_artifacts_dir}/xunit-reports"
    mkdir -p "${unit_test_artifacts_dir}"
    # Each test class produces its own report file, but they contain a fully qualified class
    # name in their file name. Due to that, we do not need to add an artifact prefix to
    # distinguish them with reports from other artifacts of the same repo/pipeline build.
    "$CP" "${unit_test_result_dir}/"*.xml "${unit_test_artifacts_dir}"
else
  echo "Build failed: no unit test results found in ${unit_test_result_dir}"
  exit 1
fi

echo "Verifying unit test coverage report was generated  ..."
coverage_result_dir="${gradle_build_dir}/reports/jacoco/test"
if [ -d "${coverage_result_dir}" ]; then
    code_coverage_artifacts_dir="${tmp_artifacts_dir}/code-coverage"
    mkdir -p "${code_coverage_artifacts_dir}"
    "$CP" "${coverage_result_dir}/jacocoTestReport.xml" "${code_coverage_artifacts_dir}/${artifact_prefix}coverage.xml"
else
  echo "Build failed: no unit test coverage report was found in ${coverage_result_dir}"
  exit 1
fi
