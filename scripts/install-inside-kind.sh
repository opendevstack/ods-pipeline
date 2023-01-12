#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}
kind_values_dir="${ODS_PIPELINE_DIR}/deploy/.kind-values"
HELM_GENERATED_VALUES_FILE="${ODS_PIPELINE_DIR}/deploy/ods-pipeline/values.generated.yaml"

URL_SUFFIX="http"
BITBUCKET_AUTH="unavailable"
NEXUS_AUTH="unavailable:unavailable"
SONAR_AUTH="unavailable"

if [ "$#" -gt 0 ]; then
    case $1 in
    --private-cert=*) URL_SUFFIX="https";
esac; fi

if [ -f "${kind_values_dir}/bitbucket-auth" ]; then
    BITBUCKET_AUTH=$(cat "${kind_values_dir}/bitbucket-auth")
fi
if [ -f "${kind_values_dir}/nexus-auth" ]; then
    NEXUS_AUTH=$(cat "${kind_values_dir}/nexus-auth")
fi
if [ -f "${kind_values_dir}/sonar-auth" ]; then
    SONAR_AUTH=$(cat "${kind_values_dir}/sonar-auth")
fi

if [ ! -e "${HELM_GENERATED_VALUES_FILE}" ]; then
    echo "setup:" > "${HELM_GENERATED_VALUES_FILE}"
fi
if [ -f "${kind_values_dir}/bitbucket-${URL_SUFFIX}" ]; then
    BITBUCKET_URL=$(cat "${kind_values_dir}/bitbucket-${URL_SUFFIX}")
    echo "  bitbucketUrl: '${BITBUCKET_URL}'" >> "${HELM_GENERATED_VALUES_FILE}"
fi
if [ -f "${kind_values_dir}/nexus-${URL_SUFFIX}" ]; then
    NEXUS_URL=$(cat "${kind_values_dir}/nexus-${URL_SUFFIX}")
    echo "  nexusUrl: '${NEXUS_URL}'" >> "${HELM_GENERATED_VALUES_FILE}"
fi
if [ -f "${kind_values_dir}/sonar-${URL_SUFFIX}" ]; then
    SONAR_URL=$(cat "${kind_values_dir}/sonar-${URL_SUFFIX}")
    echo "  sonarUrl: '${SONAR_URL}'" >> "${HELM_GENERATED_VALUES_FILE}"
fi

"${ODS_PIPELINE_DIR}"/deploy/install.sh \
    --aqua-auth "unavailable:unavailable" \
    --aqua-scanner-url "none" \
    --bitbucket-auth "${BITBUCKET_AUTH}" \
    --nexus-auth "${NEXUS_AUTH}" \
    --sonar-auth "${SONAR_AUTH}" \
    -f "./ods-pipeline/values.kind.yaml,${HELM_GENERATED_VALUES_FILE}" "$@"
