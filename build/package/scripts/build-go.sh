#!/bin/bash
set -eu

ODS_PROJECT=""
ODS_REPOSITORY=""
ODS_GIT_COMMIT=""

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --project) ODS_PROJECT="$2"; shift;;
    --project=*) ODS_PROJECT="${1#*=}";;

    --repository) ODS_REPOSITORY="$2"; shift;;
    --repository=*) ODS_REPOSITORY="${1#*=}";;

    --git-commit) ODS_GIT_COMMIT="$2"; shift;;
    --git-commit=*) ODS_GIT_COMMIT="${1#*=}";;

  *) echo_error "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

echo "Run Go build"
# TODO: CGO might be required at times - we need to make this configurable.
export CGO_ENABLED=0
go version

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
  go build -o docker/app-linux-amd64
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
      # nexus-upload \
      #   -nexus-url=${NEXUS_URL} \
      #   -nexus-user=${NEXUS_USERNAME} \
      #   -nexus-password=${NEXUS_PASSWORD} \
      #   -repository=${ODS_PROJECT} \
      #   -group=/${ODS_REPOSITORY}/${ODS_GIT_COMMIT}/test-reports \
      #   -file=build/test-results/test/report.xml || true
  else
    echo "no test results found"
  fi
  exit $exitcode
else
  make ci-test
fi
