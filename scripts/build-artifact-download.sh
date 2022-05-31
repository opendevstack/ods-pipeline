#!/usr/bin/env bash
set -ue

GO_OS=""
GO_ARCH=""

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --go-os) GO_OS="$2"; shift;;
    --go-os=*) GO_OS="${1#*=}";;

    --go-arch) GO_ARCH="$2"; shift;;
    --go-arch=*) GO_ARCH="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd ../cmd/artifact-download

GIT_COMMIT=$(git rev-parse --short HEAD)
VERSION=$(git tag --contains "${GIT_COMMIT}" | grep ^v | sort -V | tail -n 1)
if [ -z "${VERSION}" ]; then
    VERSION="dev"
fi

OUTFILE="artifact-download-${GO_OS}-${GO_ARCH}"
if [ "${GO_OS}" == "windows" ]; then
    OUTFILE="${OUTFILE}.exe"
fi

GOOS=${GO_OS} GOARCH=${GO_ARCH} CGO_ENABLED=0 go build \
	-gcflags "all=-trimpath=$(pwd);$(go env GOPATH)" \
	-ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" \
	-o "${OUTFILE}"
