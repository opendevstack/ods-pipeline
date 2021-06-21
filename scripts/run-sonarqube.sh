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
K8S_SECRET_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/secret-sonar-auth.yml"
K8S_CONFIGMAP_FILE="${ODS_PIPELINE_DIR}/test/testdata/deploy/cd-kind/configmap-sonar.yml"


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
echo "Waiting up to 4 minutes for SonarQube to start ..."
n=0
health="RED"
set +e
until [ $n -ge 24 ]; do
    health=$(curl -s ${INSECURE} --user "${SONAR_USERNAME}:${SONAR_PASSWORD}" \
        "${SONARQUBE_URL}/api/system/health" | jq -r .health)
    if [ "${health}" == "GREEN" ]; then
        echo " success"
        break
    else
        echo -n "."
        sleep 10s
        n=$((n+1))
    fi
done
set -e
if [ "${health}" != "GREEN" ]; then
    echo "SonarQube did not start, got health=${health}."
    docker logs ${CONTAINER_NAME}
    exit 1
fi

echo "Creating token for '${SONAR_USERNAME}' ..."
tokenResponse=$(curl ${INSECURE} -X POST -sSf --user "${SONAR_USERNAME}:${SONAR_PASSWORD}" \
    "${SONARQUBE_URL}/api/user_tokens/generate?login=${SONAR_USERNAME}&name=odspipeline")
# Example response:
# {"login":"cd_user","name":"foo","token":"bar","createdAt":"2020-04-22T13:21:54+0000"}
token=$(echo "${tokenResponse}" | jq -r .token)

cat <<EOF >${K8S_SECRET_FILE}
apiVersion: v1
stringData:
  password: ${token}
  username: ${SONAR_USERNAME}
kind: Secret
metadata:
  name: ods-sonar-auth
type: kubernetes.io/basic-auth
EOF

cat <<EOF >${K8S_CONFIGMAP_FILE}
kind: ConfigMap
apiVersion: v1
metadata:
  name: ods-sonar
data:
  url: 'http://${CONTAINER_NAME}.kind:9000'
EOF

echo "Created secret with token for '${SONAR_USERNAME}' in ${K8S_SECRET_FILE}."
