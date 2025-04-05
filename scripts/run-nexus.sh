#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

INSECURE=""
HOST_HTTP_PORT="8081"
HOST_HTTPS_PORT="8443"
ADMIN_USER="admin"
ADMIN_PASSWORD=""
DEVELOPER_USERNAME="developer"
DEVELOPER_PASSWORD="s3cr3t"
NEXUS_URL=
IMAGE_NAME="ods-test-nexus"
CONTAINER_NAME="ods-test-nexus"
NEXUS_IMAGE_TAG="3.30.1"
kind_values_dir="/tmp/ods-pipeline/kind-values"
mkdir -p "${kind_values_dir}"
DOCKER_CONTEXT_DIR="${ODS_PIPELINE_DIR}/test/testdata/private-cert"
reuse="false"

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    -i|--insecure) INSECURE="--insecure";;

    -r|--reuse) reuse="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${reuse}" = "true" ]; then
    if ! docker container inspect ${CONTAINER_NAME} &> /dev/null; then
        echo "No existing Nexus container ${CONTAINER_NAME} found ..."
    else
        echo "Reusing existing Nexus container ${CONTAINER_NAME} ..."
        exit 0
    fi
fi

echo "Run container using image tag ${NEXUS_IMAGE_TAG}"
docker rm -f ${CONTAINER_NAME} || true
cd "${SCRIPT_DIR}"/nexus
docker build -t ${IMAGE_NAME} -f "Dockerfile.$(uname -m)" "${DOCKER_CONTEXT_DIR}"
cd - &> /dev/null
docker run -d -p "${HOST_HTTP_PORT}:8081" --net kind --name ${CONTAINER_NAME} ${IMAGE_NAME}

if ! bash "${SCRIPT_DIR}"/waitfor-nexus.sh ; then
    docker logs ${CONTAINER_NAME}
    exit 1
fi

NEXUS_URL="http://localhost:${HOST_HTTP_PORT}"

function runJsonScript {
    local jsonScriptName=$1
    shift 1
    # shellcheck disable=SC2124
    local runParams="$@"
    echo "uploading ${jsonScriptName}.json"
    curl ${INSECURE} -X POST -sSf \
        --user "${ADMIN_USER}:${ADMIN_PASSWORD}" \
        --header 'Content-Type: application/json' \
        "${NEXUS_URL}/service/rest/v1/script" -d @"${SCRIPT_DIR}"/nexus/"${jsonScriptName}".json
    echo "running ${jsonScriptName}"
    # shellcheck disable=SC2086
    curl ${INSECURE} -X POST -sSf \
        --user "${ADMIN_USER}:${ADMIN_PASSWORD}" \
        --header 'Content-Type: text/plain' \
        "${NEXUS_URL}/service/rest/v1/script/${jsonScriptName}/run" ${runParams} > /dev/null
    echo "deleting ${jsonScriptName}"
    curl ${INSECURE} -X DELETE -sSf \
        --user "${ADMIN_USER}:${ADMIN_PASSWORD}" \
        "${NEXUS_URL}/service/rest/v1/script/${jsonScriptName}"
}

echo "Retrieving admin password"
DEFAULT_ADMIN_PASSWORD_FILE="/nexus-data/admin.password"
ADMIN_PASSWORD=$(docker exec -t "${CONTAINER_NAME}" sh -c "cat ${DEFAULT_ADMIN_PASSWORD_FILE} 2> /dev/null || true")

echo "Install Blob Stores"
runJsonScript "createBlobStores"

echo "Install Repositories"
runJsonScript "createRepos"

echo "Deactivate anonymous access"
runJsonScript "deactivateAnonymous"

echo "Setup developer role"
runJsonScript "createRole" "-d @${SCRIPT_DIR}/nexus/developer-role.json"

echo "Setup developer user"
sed "s|@developer_password@|${DEVELOPER_PASSWORD}|g" "${SCRIPT_DIR}"/nexus/developer-user.json > /tmp/nexus-developer-user-with-password.json
runJsonScript "createUser" "-d @/tmp/nexus-developer-user-with-password.json"
rm /tmp/nexus-developer-user-with-password.json

echo "Launch TLS proxy"
TLS_CONTAINER_NAME="${CONTAINER_NAME}-tls"
bash "${SCRIPT_DIR}/run-tls-proxy.sh" \
  --container-name "${TLS_CONTAINER_NAME}" \
  --https-port "${HOST_HTTPS_PORT}" \
  --nginx-conf "nginx-nexus.conf"

# Write values / secrets so that it can be picked up by install.sh later.
mkdir -p "${kind_values_dir}"
echo -n "https://${TLS_CONTAINER_NAME}.kind:${HOST_HTTPS_PORT}" > "${kind_values_dir}/nexus-https"
echo -n "http://${CONTAINER_NAME}.kind:${HOST_HTTP_PORT}" > "${kind_values_dir}/nexus-http"
echo -n "${DEVELOPER_USERNAME}:${DEVELOPER_PASSWORD}" > "${kind_values_dir}/nexus-auth"
