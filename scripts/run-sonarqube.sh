#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

INSECURE=""
HOST_HTTP_PORT="9000"
HOST_HTTPS_PORT="9443"
IMAGE_NAME="ods-test-sonarqube"
CONTAINER_NAME="ods-test-sonarqube"
SONAR_VERSION="8.9"
SONAR_USERNAME="admin"
SONAR_PASSWORD="admin"
SONAR_EDITION="community"
SONAR_IMAGE_TAG="${SONAR_VERSION}-${SONAR_EDITION}"
kind_values_dir="/tmp/ods-pipeline/kind-values"
mkdir -p "${kind_values_dir}"
reuse="false"

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    -i|--insecure) INSECURE="--insecure";;

    -r|--reuse) reuse="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${reuse}" = "true" ]; then
    if [ "$(docker inspect ${CONTAINER_NAME} -f '{{.State.Running}}')" = "true" ]; then
        echo "Reusing running SonarQube container ${CONTAINER_NAME} ..."
        exit 0
    else
        echo "No running SonarQube container ${CONTAINER_NAME} found ..."
    fi
fi

echo "Run container using image tag ${SONAR_IMAGE_TAG}"
docker rm -f ${CONTAINER_NAME} || true

cd "${SCRIPT_DIR}"/sonarqube

if [ "$(uname -m)" = "arm64" ]; then
    SONAR_ARM_IMAGE_NAME="ods-test-sonarqube-arm"
    if [ "$(docker images -q ${SONAR_ARM_IMAGE_NAME}:${SONAR_IMAGE_TAG} 2> /dev/null)" == "" ]; then
        echo "Building SonarQube arm64 image ..."
        rm -rf docker-sonarqube || true
        git clone https://github.com/SonarSource/docker-sonarqube
        cd docker-sonarqube
        git checkout refs/tags/9.7.0 # Last available Git tag
        cd 9/community
        docker build -t sonarqube-arm:${SONAR_IMAGE_TAG} .
        cd "${SCRIPT_DIR}"/sonarqube
        docker build -t ${SONAR_ARM_IMAGE_NAME}:${SONAR_IMAGE_TAG} --build-arg=from=sonarqube-arm:${SONAR_IMAGE_TAG} .
    else
        echo "Using existing ${SONAR_ARM_IMAGE_NAME}:${SONAR_IMAGE_TAG} image"
    fi
    IMAGE_NAME="${SONAR_ARM_IMAGE_NAME}"
else
    docker build -t ${IMAGE_NAME}:${SONAR_IMAGE_TAG} --build-arg=from=sonarqube:${SONAR_IMAGE_TAG} .
fi
cd - &> /dev/null
docker run -d --net kind --name ${CONTAINER_NAME} -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p "${HOST_HTTP_PORT}:9000" ${IMAGE_NAME}:${SONAR_IMAGE_TAG}

SONARQUBE_URL="http://localhost:${HOST_HTTP_PORT}"
if ! bash "${SCRIPT_DIR}/waitfor-sonarqube.sh" ; then
    docker logs ${CONTAINER_NAME}
    exit 1
fi

echo "Creating token for '${SONAR_USERNAME}' ..."
tokenResponse=$(curl ${INSECURE} -X POST -sSf --user "${SONAR_USERNAME}:${SONAR_PASSWORD}" \
    "${SONARQUBE_URL}/api/user_tokens/generate?login=${SONAR_USERNAME}&name=odspipeline")
# Example response:
# {"login":"cd_user","name":"foo","token":"bar","createdAt":"2020-04-22T13:21:54+0000"}
token=$(echo "${tokenResponse}" | jq -r .token)

echo "Launch TLS proxy"
TLS_CONTAINER_NAME="${CONTAINER_NAME}-tls"
bash "${SCRIPT_DIR}/run-tls-proxy.sh" \
  --container-name "${TLS_CONTAINER_NAME}" \
  --https-port "${HOST_HTTPS_PORT}" \
  --nginx-conf "nginx-sonarqube.conf"

# Write values / secrets so that it can be picked up by install.sh later.
mkdir -p "${kind_values_dir}"
echo -n "https://${TLS_CONTAINER_NAME}.kind:${HOST_HTTPS_PORT}" > "${kind_values_dir}/sonar-https"
echo -n "http://${CONTAINER_NAME}.kind:${HOST_HTTP_PORT}" > "${kind_values_dir}/sonar-http"
echo -n ":${token}" > "${kind_values_dir}/sonar-auth"
