#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

HOST_PORT="7990"
CONTAINER_NAME="bitbuckettest"
BITBUCKET_IMAGE_TAG="latest"
K8S_SECRET_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/secret-bitbucket-auth.yml"
K8S_CONFIGMAP_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/configmap-bitbucket.yml"

echo "Run container using image tag ${BITBUCKET_IMAGE_TAG}"
docker rm -f ${CONTAINER_NAME} || true
docker run -d --name ${CONTAINER_NAME} -e ELASTICSEARCH_ENABLED=false --net kind -p "${HOST_PORT}:7990" -p 7999:7999 atlassian/bitbucket-server:${BITBUCKET_IMAGE_TAG}

BITBUCKET_URL="http://localhost:${HOST_PORT}"
cat <<EOF >${K8S_SECRET_FILE}
apiVersion: v1
data:
  password: YWRtaW4=
  username: YWRtaW4=
kind: Secret
metadata:
  name: bitbucket-auth
type: kubernetes.io/basic-auth
EOF

cat <<EOF >${K8S_CONFIGMAP_FILE}
kind: ConfigMap
apiVersion: v1
metadata:
  name: bitbucket
data:
  url: 'http://${CONTAINER_NAME}.kind:7999'
EOF

echo "Visit ${BITBUCKET_URL} to setup Bitbucket."
echo "Use timebomb license from https://developer.atlassian.com/platform/marketplace/timebomb-licenses-for-testing-server-apps/ (3 hour expiration for all Atlassian host products)"
