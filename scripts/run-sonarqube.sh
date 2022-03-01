#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

INSECURE=""
HOST_PORT="9000"
IMAGE_NAME="ods-test-sonarqube"
CONTAINER_NAME="ods-test-sonarqube"
SONAR_VERSION="8.4"
SONAR_USERNAME="admin"
SONAR_PASSWORD="admin"
SONAR_EDITION="community"
SONAR_IMAGE_TAG="${SONAR_VERSION}-${SONAR_EDITION}"
HELM_VALUES_FILE="${ODS_PIPELINE_DIR}/deploy/ods-pipeline/values.generated.yaml"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    -i|--insecure) INSECURE="--insecure";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

echo "Run container using image tag ${SONAR_IMAGE_TAG}"
docker rm -f ${CONTAINER_NAME} || true
cd "${SCRIPT_DIR}"/sonarqube
docker build -t ${IMAGE_NAME} .
cd - &> /dev/null
docker run -d --net kind --name ${CONTAINER_NAME} -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p "${HOST_PORT}:9000" ${IMAGE_NAME}

SONARQUBE_URL="http://localhost:${HOST_PORT}"
if ! "${SCRIPT_DIR}/waitfor-sonarqube.sh" ; then
    docker logs ${CONTAINER_NAME}
    exit 1
fi 

echo "Creating token for '${SONAR_USERNAME}' ..."
tokenResponse=$(curl ${INSECURE} -X POST -sSf --user "${SONAR_USERNAME}:${SONAR_PASSWORD}" \
    "${SONARQUBE_URL}/api/user_tokens/generate?login=${SONAR_USERNAME}&name=odspipeline")
# Example response:
# {"login":"cd_user","name":"foo","token":"bar","createdAt":"2020-04-22T13:21:54+0000"}
token=$(echo "${tokenResponse}" | jq -r .token)

if [ ! -e "${HELM_VALUES_FILE}" ]; then
    echo "ods-pipeline-setup:" > "${HELM_VALUES_FILE}"
fi

{
    echo "  sonarUrl: 'http://${CONTAINER_NAME}.kind:9000'"
    echo "  sonarUsername: '${SONAR_USERNAME}'"
    echo "  sonarAuthToken: '${token}'"
} >> "${HELM_VALUES_FILE}"
