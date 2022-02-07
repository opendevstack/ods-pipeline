#!/usr/bin/env bash

set -o -pipefail -c

HELM_VERSION=3.5.2
SOPS_VERSION=3.7.1
AGE_VERSION=1.0.0
GOBIN=/tmp/bin
SAVED_PATH=$PATH
TARGETARCH=

PATH=$GOBIN:$SAVED_PATH

mkdir -p $GOBIN

mkdir -p /tmp/helm-unpacked \
    && cd /tmp \
    && curl -LO https://get.helm.sh/helm-v${HELM_VERSION}-linux-${TARGETARCH}.tar.gz \
    && tar -zxvf helm-v${HELM_VERSION}-linux-${TARGETARCH}.tar.gz -C /tmp/helm-unpacked \
    && mv /tmp/helm-unpacked/linux-${TARGETARCH}/helm /tmp/helm \
    && chmod a+x /tmp/helm \
    && ./tmp/helm version \
    && ./tmp/helm env

cd ../../../cmd/deploy-with-helm && CGO_ENABLED=0 go build -o /tmp/deploy-with-helm

# install sops
go install go.mozilla.org/sops/v3/cmd/sops@v${SOPS_VERSION} \
    && sops --version

# install age
go install filippo.io/age/cmd/...@v${AGE_VERSION} \
    && age --version

PATH=$SAVED_PATH

docker build -f ../Docker.helm-singlestage ../
