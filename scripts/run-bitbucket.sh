#!/usr/bin/env bash
set -ue

# Starts a Bitbucket instance with a timebomb license (3 hours).
# The instance is setup with an admin account (pw: admin) and an "ODSPIPELINETEST" project.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

BITBUCKET_SERVER_HOST_PORT="7990"
BITBUCKET_SERVER_CONTAINER_NAME="ods-test-bitbucket-server"
BITBUCKET_SERVER_IMAGE_NAME="atlassian/bitbucket"
BITBUCKET_SERVER_IMAGE_TAG="7.6.5"
BITBUCKET_POSTGRES_HOST_PORT="5432"
BITBUCKET_POSTGRES_CONTAINER_NAME="ods-test-bitbucket-postgres"
BITBUCKET_POSTGRES_IMAGE_TAG="12"
HELM_VALUES_FILE="${ODS_PIPELINE_DIR}/deploy/ods-pipeline/values.generated.yaml"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

echo "Run Postgres container"
docker rm -f ${BITBUCKET_POSTGRES_CONTAINER_NAME} || true
docker run  --name ${BITBUCKET_POSTGRES_CONTAINER_NAME} \
  -v "${ODS_PIPELINE_DIR}"/test/testdata/bitbucket-sql:/docker-entrypoint-initdb.d \
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
  "${BITBUCKET_SERVER_IMAGE_NAME}:${BITBUCKET_SERVER_IMAGE_TAG}"

if ! "${SCRIPT_DIR}/waitfor-bitbucket.sh" ; then
    docker logs ${BITBUCKET_SERVER_CONTAINER_NAME}
    exit 1
fi 
BITBUCKET_URL_FULL="http://${BITBUCKET_SERVER_CONTAINER_NAME}.kind:7990"

if [ ! -e "${HELM_VALUES_FILE}" ]; then
    echo "setup:" > "${HELM_VALUES_FILE}"
fi

{
  echo "  bitbucketUrl: '${BITBUCKET_URL_FULL}'"
  echo "  bitbucketUsername: 'admin'"
  echo "  bitbucketAccessToken: 'NzU0OTk1MjU0NjEzOpzj5hmFNAaawvupxPKpcJlsfNgP'"
  echo "  bitbucketWebhookSecret: 's3cr3t'"
} >> "${HELM_VALUES_FILE}"
