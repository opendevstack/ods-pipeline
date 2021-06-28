#!/bin/bash
set -eu

ENABLE_CGO="false"
GO_OS=""
GO_ARCH=""

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --enable-cgo) ENABLE_CGO="$2"; shift;;
    --enable-cgo=*) ENABLE_CGO="${1#*=}";;

    --go-os) GO_OS="$2"; shift;;
    --go-os=*) GO_OS="${1#*=}";;

    --go-arch) GO_ARCH="$2"; shift;;
    --go-arch=*) GO_ARCH="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

echo "Run Go build"
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

echo "Check format"
make -q ci-checkfmt &> /dev/null || makeErrorCode=$?
if [ "${makeErrorCode}" -eq 2 ]; then
  unformatted=$(gofmt -l .)
  if [ -n "${unformatted}" ]; then
    echo "Unformatted files:"
    echo "${unformatted}"
    echo "All files need to be gofmt'd. Please run: gofmt -w ."
    exit 1
  fi
else
  make ci-checkfmt
fi

echo "Lint"
make -q ci-lint &> /dev/null || makeErrorCode=$?
if [ "${makeErrorCode}" -eq 2 ]; then
  golangci-lint version
  set +e
  golangci-lint run > go-lint-report.txt
  exitcode=$?
  set -e
  if [ -s go-lint-report.txt ]; then
    cat go-lint-report.txt
    mkdir -p .ods/artifacts/lint-report
    cp go-lint-report.txt .ods/artifacts/lint-report/report.txt
    exit $exitcode
  fi
else
  make ci-lint
fi

echo "Build"
make -q ci-build &> /dev/null || makeErrorCode=$?
if [ "${makeErrorCode}" -eq 2 ]; then
  go build -o docker/app
else
  make ci-build
fi

echo "Test"
make -q ci-test &> /dev/null || makeErrorCode=$?
if [ "${makeErrorCode}" -eq 2 ]; then
  mkdir -p build/test-results/test
  if [ -f .ods/artifacts/xunit-reports/report.xml ]; then
    echo "restoring archived test artifacts ..."
    cp .ods/artifacts/xunit-reports/report.xml build/test-results/test/report.xml
    cp .ods/artifacts/code-coverage/coverage.out coverage.out
  else
    GOPKGS=$(go list ./... | grep -v /vendor)
    set +e
    echo "running tests ..."
    go test -v -coverprofile=coverage.out $GOPKGS 2>&1 > test-results.txt
    exitcode=$?
    set -e
    if [ -f test-results.txt ]; then
        go-junit-report < test-results.txt > build/test-results/test/report.xml
        mkdir -p .ods/artifacts/xunit-reports
        cp build/test-results/test/report.xml .ods/artifacts/xunit-reports/report.xml
    else
      echo "no test results found"
    fi
    if [ -f coverage.out ]; then
        mkdir -p .ods/artifacts/code-coverage
        cp coverage.out .ods/artifacts/code-coverage/coverage.out
    else
      echo "no code coverage found"
    fi
    exit $exitcode
  fi
else
  make ci-test
fi
