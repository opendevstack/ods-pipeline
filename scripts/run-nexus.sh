#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

INSECURE=""
HOST_PORT="8081"
ADMIN_USER="admin"
ADMIN_PASSWORD=""
DEVELOPER_USERNAME="developer"
DEVELOPER_PASSWORD="s3cr3t"
NEXUS_URL=
CONTAINER_NAME="nexustest"
NEXUS_IMAGE_TAG="3.30.1"
K8S_SECRET_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/secret-nexus-auth.yml"
K8S_CONFIGMAP_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/configmap-nexus.yml"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    -i|--insecure) INSECURE="--insecure";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

echo "Run container using image tag ${NEXUS_IMAGE_TAG}"
docker rm -f ${CONTAINER_NAME} || true
cd ${SCRIPT_DIR}/nexus
docker build -t nexustest .
cd -
docker volume create --name nexus-vol --driver local --opt type=tmpfs --opt device=tmpfs --opt o=size=2G
docker run -d -p "${HOST_PORT}:8081" -v nexus-vol --net kind --name ${CONTAINER_NAME} nexustest

NEXUS_URL="http://localhost:${HOST_PORT}"
function waitForReady {
    echo "Waiting up to 4 minutes for Nexus to start ..."
    local n=0
    local http_code=
    set +e
    until [ $n -ge 24 ]; do
        http_code=$(curl ${INSECURE} -s -o /dev/null -w "%{http_code}" "${NEXUS_URL}/service/rest/v1/status/writable")
        if [ "${http_code}" == "200" ]; then
            echo " success"
            break
        else
            echo -n "."
            sleep 10s
            n=$((n+1))
        fi
    done
    set -e

    if [ "${http_code}" != "200" ]; then
        echo "Nexus did not start, got http_code=${http_code}."
        docker logs ${CONTAINER_NAME}
        exit 1
    fi
}

function runJsonScript {
    local jsonScriptName=$1
    shift 1
    # shellcheck disable=SC2124
    local runParams="$@"
    echo "uploading ${jsonScriptName}.json"
    curl ${INSECURE} -X POST -sSf \
        --user "${ADMIN_USER}:${ADMIN_PASSWORD}" \
        --header 'Content-Type: application/json' \
        "${NEXUS_URL}/service/rest/v1/script" -d @${SCRIPT_DIR}/nexus/"${jsonScriptName}".json
    echo "running ${jsonScriptName}"
    curl ${INSECURE} -X POST -sSf \
        --user "${ADMIN_USER}:${ADMIN_PASSWORD}" \
        --header 'Content-Type: text/plain' \
        "${NEXUS_URL}/service/rest/v1/script/${jsonScriptName}/run" ${runParams} > /dev/null
    echo "deleting ${jsonScriptName}"
    curl ${INSECURE} -X DELETE -sSf \
        --user "${ADMIN_USER}:${ADMIN_PASSWORD}" \
        "${NEXUS_URL}/service/rest/v1/script/${jsonScriptName}"
}

waitForReady

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
sed "s|@developer_password@|${DEVELOPER_PASSWORD}|g" ${SCRIPT_DIR}/nexus/developer-user.json > ${SCRIPT_DIR}/nexus/developer-user-with-password.json
runJsonScript "createUser" "-d @${SCRIPT_DIR}/nexus/developer-user-with-password.json"
rm ${SCRIPT_DIR}/nexus/developer-user-with-password.json

cat <<EOF >${K8S_SECRET_FILE}
apiVersion: v1
stringData:
  password: ${DEVELOPER_PASSWORD}
  username: ${DEVELOPER_USERNAME}
kind: Secret
metadata:
  name: nexus-auth
type: kubernetes.io/basic-auth
EOF

cat <<EOF >${K8S_CONFIGMAP_FILE}
kind: ConfigMap
apiVersion: v1
metadata:
  name: nexus
data:
  url: 'http://${CONTAINER_NAME}.kind:8081'
EOF

echo "Created secret with password for '${DEVELOPER_USERNAME}' in ${K8S_SECRET_FILE}."
