#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}
ODS_KIND_CREDENTIALS_DIR="${ODS_PIPELINE_DIR}/deploy/.kind-credentials"

BITBUCKET_AUTH="unavailable"
NEXUS_AUTH="unavailable:unavailable"
SONAR_AUTH="unavailable"

if [ -f "${ODS_KIND_CREDENTIALS_DIR}/bitbucket-auth" ]; then
    BITBUCKET_AUTH=$(cat "${ODS_KIND_CREDENTIALS_DIR}/bitbucket-auth")
fi
if [ -f "${ODS_KIND_CREDENTIALS_DIR}/nexus-auth" ]; then
    NEXUS_AUTH=$(cat "${ODS_KIND_CREDENTIALS_DIR}/nexus-auth")
fi
if [ -f "${ODS_KIND_CREDENTIALS_DIR}/sonar-auth" ]; then
    SONAR_AUTH=$(cat "${ODS_KIND_CREDENTIALS_DIR}/sonar-auth")
fi

"${ODS_PIPELINE_DIR}"/deploy/install.sh \
    --aqua-auth "unavailable:unavailable" \
    --bitbucket-auth "${BITBUCKET_AUTH}" \
    --nexus-auth "${NEXUS_AUTH}" \
    --sonar-auth "${SONAR_AUTH}" \
    -f ./ods-pipeline/values.kind.yaml,./ods-pipeline/values.generated.yaml "$@"
