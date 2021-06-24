#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

INSECURE=""
HOST_PORT="9000"
CONTAINER_NAME="sonarqubetest"
SONAR_VERSION="8.4"
SONAR_USERNAME="admin"
SONAR_PASSWORD="admin"
SONAR_EDITION="community"
SONAR_IMAGE_TAG="${SONAR_VERSION}-${SONAR_EDITION}"
HELM_VALUES_FILE="${ODS_PIPELINE_DIR}/deploy/cd-namespace/chart/values.generated.yaml"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    -i|--insecure) INSECURE="--insecure";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

# TODO: Set forceAuthentication to mirror actual settings.
# echo "Build image"
# docker build \
#     -t "${CONTAINER_IMAGE}" \
#     --build-arg sonarDistributionUrl="${SONAR_DISTRIBUTION_URL}" \
#     --build-arg sonarVersion="${SONAR_VERSION}" \
#     --build-arg sonarEdition="${SONAR_EDITION}" \
#     --build-arg idpDns="" \
#     ./docker

echo "Run container using image tag ${SONAR_IMAGE_TAG}"
docker rm -f ${CONTAINER_NAME} || true
docker run -d --net kind --name ${CONTAINER_NAME} -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p "${HOST_PORT}:9000" sonarqube:${SONAR_IMAGE_TAG}

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

echo "sonarUrl: 'http://${CONTAINER_NAME}.kind:9000'" >> ${HELM_VALUES_FILE}
echo "sonarUsername: '${SONAR_USERNAME}'" >> ${HELM_VALUES_FILE}
echo "sonarPassword: '${token}'" >> ${HELM_VALUES_FILE}
