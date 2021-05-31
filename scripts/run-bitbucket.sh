#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

INSECURE=""
HOST_PORT="7990"
CONTAINER_NAME="bitbuckettest"
BITBUCKET_IMAGE_TAG="latest"
K8S_SECRET_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/secret-bitbucket-auth.yml"
K8S_CONFIGMAP_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/configmap-bitbucket.yml"

echo "Run container using image tag ${BITBUCKET_IMAGE_TAG}"
docker rm -f ${CONTAINER_NAME} || true
docker run -d --name ${CONTAINER_NAME} -e ELASTICSEARCH_ENABLED=false --net kind -p "${HOST_PORT}:7990" -p 7999:7999 atlassian/bitbucket-server:${BITBUCKET_IMAGE_TAG}

BITBUCKET_URL="http://localhost:${HOST_PORT}"
echo "Waiting up to 3 minutes for Bitbucket to start ..."
# https://confluence.atlassian.com/bitbucketserverkb/how-to-monitor-if-bitbucket-server-is-up-and-running-975014635.html
n=0
set +e
until [ $n -ge 18 ]; do
    health=$(curl -s ${INSECURE} "${BITBUCKET_URL}/status" | jq -r .state)
    if [ "${health}" == "RUNNING" ]; then
        echo " success"
        break
    else
        if [ "${health}" == "FIRST_RUN" ]; then
            echo " up but still needs setup"
            break
        else
            echo -n "."
            sleep 10s
            n=$((n+1))
        fi
    fi
done
set -e

cat <<EOF >${K8S_SECRET_FILE}
apiVersion: v1
stringData:
  password: create_access_token_and_enter_here
  username: admin
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

echo "If setup is needed, visit ${BITBUCKET_URL} and create an access token and a project 'FOO'."
echo "Use timebomb license from https://developer.atlassian.com/platform/marketplace/timebomb-licenses-for-testing-server-apps/ (3 hour expiration for all Atlassian host products)"
