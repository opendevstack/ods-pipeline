#!/usr/bin/env bash
set -ue

# Starts a Bitbucket instance with a timebomb license (3 hours).
# The instance is setup with an admin account (pw: admin) and an "ODSPIPELINETEST" project.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

INSECURE=""
BITBUCKET_SERVER_HOST_PORT="7990"
BITBUCKET_SERVER_CONTAINER_NAME="bitbucket-server"
BITBUCKET_SERVER_IMAGE_TAG="latest"
BITBUCKET_POSTGRES_HOST_PORT="7990"
BITBUCKET_POSTGRES_CONTAINER_NAME="bitbucket-postgres"
BITBUCKET_POSTGRES_IMAGE_TAG="12"
K8S_SECRET_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/secret-bitbucket-auth.yml"
K8S_CONFIGMAP_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/configmap-bitbucket.yml"


echo "Run Postgres container"
docker rm -f ${BITBUCKET_POSTGRES_CONTAINER_NAME} || true
docker run  --name ${BITBUCKET_POSTGRES_CONTAINER_NAME} \
  -v ${ODS_PIPELINE_DIR}/testdata/bitbucket.sql:/docker-entrypoint-initdb.d/init.sql \
  -e POSTGRES_PASSWORD=jellyfish -e POSTGRES_USER=bitbucketuser -e POSTGRES_DB=bitbucket \
  -d --net kind -p "${BITBUCKET_POSTGRES_HOST_PORT}:5432" \
  postgres:${BITBUCKET_POSTGRES_IMAGE_TAG}

echo "Run Bitbucket Server pointing to Postgres"
docker rm -f ${BITBUCKET_SERVER_CONTAINER_NAME} || true
docker run --name ${BITBUCKET_SERVER_CONTAINER_NAME} \
  -e JDBC_DRIVER=org.postgresql.Driver \
  -e JDBC_USER=bitbucketuser \
  -e JDBC_PASSWORD=jellyfish \
  -e JDBC_URL=jdbc:postgresql://${BITBUCKET_POSTGRES_CONTAINER_NAME}.kind:5432/bitbucket \
  -e ELASTICSEARCH_ENABLED=false \
  -d --net kind -p "${BITBUCKET_SERVER_HOST_PORT}:7990" -p 7999:7999 \
  atlassian/bitbucket-server:${BITBUCKET_SERVER_IMAGE_TAG}

BITBUCKET_URL="http://localhost:${BITBUCKET_SERVER_HOST_PORT}"
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

# Personal access token (PAT) is baked into the SQL dump.
cat <<EOF >${K8S_SECRET_FILE}
apiVersion: v1
stringData:
  password: NzU0OTk1MjU0NjEzOpzj5hmFNAaawvupxPKpcJlsfNgP
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
  url: 'http://${BITBUCKET_SERVER_CONTAINER_NAME}.kind:7999'
EOF
