#!/bin/bash
set -eu

copyLintReport() {
  cat golangci-lint-report.txt
  mkdir -p "${ROOT_DIR}/.ods/artifacts/lint-reports"
  cp golangci-lint-report.txt "${ROOT_DIR}/.ods/artifacts/lint-reports/${ARTIFACT_PREFIX}report.txt"
}

ENABLE_CGO="false"
GO_OS=""
GO_ARCH=""
OUTPUT_DIR="docker"
WORKING_DIR="."
ARTIFACT_PREFIX=""
PRE_TEST_SCRIPT=""
DEBUG="${DEBUG:-false}"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --working-dir) WORKING_DIR="$2"; shift;;
    --working-dir=*) WORKING_DIR="${1#*=}";;

    --enable-cgo) ENABLE_CGO="$2"; shift;;
    --enable-cgo=*) ENABLE_CGO="${1#*=}";;

    --go-os) GO_OS="$2"; shift;;
    --go-os=*) GO_OS="${1#*=}";;

    --go-arch) GO_ARCH="$2"; shift;;
    --go-arch=*) GO_ARCH="${1#*=}";;

    --output-dir) OUTPUT_DIR="$2"; shift;;
    --output-dir=*) OUTPUT_DIR="${1#*=}";;

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
if [ "${WORKING_DIR}" != "." ]; then
  cd "${WORKING_DIR}"
  ARTIFACT_PREFIX="${WORKING_DIR/\//-}-"
fi

echo "Working on Go module in $(pwd) ..."

go version
if [ "${ENABLE_CGO}" = "false" ]; then
  export CGO_ENABLED=0
fi
if [ -n "${GO_OS}" ]; then
  export GOOS="${GO_OS}"
fi
if [ -n "${GO_ARCH}" ]; then
  export GOARCH="${GO_ARCH}"
fi
export GOMODCACHE="$ROOT_DIR/.ods-cache/deps/gomod"
echo INFO: Using gomodule cache on repo pvc
echo GOMODCACHE="$GOMODCACHE"
df -h "$ROOT_DIR"

echo "Checking format ..."
# shellcheck disable=SC2046
unformatted=$(go fmt $(go list ./...))
if [ -n "${unformatted}" ]; then
  echo "Unformatted files:"
  echo "${unformatted}"
  echo "All files need to be gofmt'd. Please run: gofmt -w ."
  exit 1
fi

echo "Linting ..."
golangci-lint version
set +e
rm golangci-lint-report.txt &>/dev/null
golangci-lint run > golangci-lint-report.txt
exitcode=$?
set -e
if [ $exitcode == 0 ]; then
  echo "OK" > golangci-lint-report.txt
  copyLintReport
else
  copyLintReport
  exit $exitcode
fi

echo "Testing ..."
if [ -f "${ROOT_DIR}/.ods/artifacts/xunit-reports/${ARTIFACT_PREFIX}report.xml" ]; then
  echo "Test artifacts already present, skipping tests ..."
  # Copy artifacts to working directory so that the SonarQube scanner can pick them up later.
  cp "${ROOT_DIR}/.ods/artifacts/xunit-reports/${ARTIFACT_PREFIX}report.xml" report.xml
  cp "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}coverage.out" coverage.out
else
  if [ -n "${PRE_TEST_SCRIPT}" ]; then
    echo "Executing pre-test script ..."
    ./"${PRE_TEST_SCRIPT}"
  fi
  GOPKGS=$(go list ./... | grep -v /vendor)
  set +e
  rm coverage.out test-results.txt report.xml &>/dev/null
  go test -v -coverprofile=coverage.out "$GOPKGS" > test-results.txt 2>&1
  exitcode=$?
  set -e
  df -h "$ROOT_DIR"
  if [ -f test-results.txt ]; then
      cat test-results.txt
      go-junit-report < test-results.txt > report.xml
      mkdir -p "${ROOT_DIR}/.ods/artifacts/xunit-reports"
      cp report.xml "${ROOT_DIR}/.ods/artifacts/xunit-reports/${ARTIFACT_PREFIX}report.xml"
  else
    echo "No test results found"
    exit 1
  fi
  if [ -f coverage.out ]; then
      mkdir -p "${ROOT_DIR}/.ods/artifacts/code-coverage"
      cp coverage.out "${ROOT_DIR}/.ods/artifacts/code-coverage/${ARTIFACT_PREFIX}coverage.out"
  else
    echo "No code coverage found"
    exit 1
  fi
  if [ $exitcode != 0 ]; then
    exit $exitcode
  fi
fi

echo "Building ..."
go build -gcflags "all=-trimpath=$(pwd)" -o "${OUTPUT_DIR}/app"

supply-sonar-project-properties-default
