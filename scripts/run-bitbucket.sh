#!/usr/bin/env bash
set -ue

# Starts a Bitbucket instance with a timebomb license (3 hours).
# The instance is setup with an admin account (pw: admin) and an "ODSPIPELINETEST" project.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

HOST_HTTP_PORT="7990"
HOST_HTTPS_PORT="7993"
BITBUCKET_SERVER_CONTAINER_NAME="ods-test-bitbucket-server"
BITBUCKET_SERVER_IMAGE_NAME="atlassian/bitbucket"
BITBUCKET_SERVER_IMAGE_TAG="7.6.5"
BITBUCKET_POSTGRES_CONTAINER_NAME="ods-test-bitbucket-postgres"
BITBUCKET_POSTGRES_IMAGE_TAG="12"
kind_values_dir="/tmp/ods-pipeline/kind-values"
mkdir -p "${kind_values_dir}"
reuse="false"

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    -r|--reuse) reuse="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${reuse}" = "true" ]; then
    if ! docker container inspect ${BITBUCKET_SERVER_CONTAINER_NAME} &> /dev/null; then
        echo "No existing Bitbucket container ${BITBUCKET_SERVER_CONTAINER_NAME} found ..."
    else
        echo "Reusing existing Bitbucket container ${BITBUCKET_SERVER_CONTAINER_NAME} ..."
        exit 0
    fi
fi

echo "Run Postgres container"
docker rm -f ${BITBUCKET_POSTGRES_CONTAINER_NAME} || true
docker run  --name ${BITBUCKET_POSTGRES_CONTAINER_NAME} \
  -v "${ODS_PIPELINE_DIR}"/test/testdata/bitbucket-sql:/docker-entrypoint-initdb.d \
  -e POSTGRES_PASSWORD=jellyfish -e POSTGRES_USER=bitbucketuser -e POSTGRES_DB=bitbucket \
  -d --net kind \
  postgres:${BITBUCKET_POSTGRES_IMAGE_TAG}

echo "Run Bitbucket Server pointing to Postgres"
docker rm -f ${BITBUCKET_SERVER_CONTAINER_NAME} || true

cd "${SCRIPT_DIR}"/bitbucket

if [ "$(uname -m)" = "arm64" ]; then
  BITBUCKET_SERVER_ARM_IMAGE_NAME="ods-test-bitbucket-arm"
  if [ "$(docker images -q ${BITBUCKET_SERVER_ARM_IMAGE_NAME}:${BITBUCKET_SERVER_IMAGE_TAG} 2> /dev/null)" == "" ]; then
    echo "Building Bitbucket Server arm64 image ..."
    rm -rf docker-atlassian-bitbucket-server || true
    git clone --recurse-submodules https://bitbucket.org/atlassian-docker/docker-atlassian-bitbucket-server.git
    cd docker-atlassian-bitbucket-server
    git checkout 0e57b62 # Last known working Git commit (no branches / tags available)
    docker build -t ${BITBUCKET_SERVER_ARM_IMAGE_NAME}:${BITBUCKET_SERVER_IMAGE_TAG} --build-arg BITBUCKET_VERSION=${BITBUCKET_SERVER_IMAGE_TAG} .
  else
    echo "Using existing ${BITBUCKET_SERVER_ARM_IMAGE_NAME}:${BITBUCKET_SERVER_IMAGE_TAG} image"
  fi
  BITBUCKET_SERVER_IMAGE_NAME="${BITBUCKET_SERVER_ARM_IMAGE_NAME}"
fi

docker run --name ${BITBUCKET_SERVER_CONTAINER_NAME} \
  -e JDBC_DRIVER=org.postgresql.Driver \
  -e JDBC_USER=bitbucketuser \
  -e JDBC_PASSWORD=jellyfish \
  -e JDBC_URL=jdbc:postgresql://${BITBUCKET_POSTGRES_CONTAINER_NAME}.kind:5432/bitbucket \
  -e ELASTICSEARCH_ENABLED=false \
  -d --net kind -p "${HOST_HTTP_PORT}:7990" -p 7999:7999 \
  "${BITBUCKET_SERVER_IMAGE_NAME}:${BITBUCKET_SERVER_IMAGE_TAG}"

if ! bash "${SCRIPT_DIR}/waitfor-bitbucket.sh" ; then
    docker logs ${BITBUCKET_SERVER_CONTAINER_NAME}
    exit 1
fi 

echo "Launch TLS proxy"
TLS_CONTAINER_NAME="${BITBUCKET_SERVER_CONTAINER_NAME}-tls"
bash "${SCRIPT_DIR}/run-tls-proxy.sh" \
  --container-name "${TLS_CONTAINER_NAME}" \
  --https-port "${HOST_HTTPS_PORT}" \
  --nginx-conf "nginx-bitbucket.conf"

# Write values / secrets so that it can be picked up by install.sh later.
mkdir -p "${kind_values_dir}"
echo -n "https://${TLS_CONTAINER_NAME}.kind:${HOST_HTTPS_PORT}" > "${kind_values_dir}/bitbucket-https"
echo -n "http://${BITBUCKET_SERVER_CONTAINER_NAME}.kind:${HOST_HTTP_PORT}" > "${kind_values_dir}/bitbucket-http"
echo -n "admin:NzU0OTk1MjU0NjEzOpzj5hmFNAaawvupxPKpcJlsfNgP" > "${kind_values_dir}/bitbucket-auth"
echo -n "webhook:s3cr3t" > "${kind_values_dir}/bitbucket-webhook-secret"
