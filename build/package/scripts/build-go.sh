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
  golangci-lint run
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
  GOPKGS=$(go list ./... | grep -v /vendor)
  set +e
  go test -v -coverprofile=coverage.out $GOPKGS 2>&1 > test-results.txt
  exitcode=$?
  if [ -f test-results.txt ]; then
      set -e
      go-junit-report < test-results.txt > build/test-results/test/report.xml
      mkdir -p ./ods/artifacts/xunit-reports
      cp build/test-results/test/report.xml ./ods/artifacts/xunit-reports/report.xml
  else
    echo "no test results found"
  fi
  exit $exitcode
else
  make ci-test
fi
