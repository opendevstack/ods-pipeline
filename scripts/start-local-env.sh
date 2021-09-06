#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

kind_registry='kind-registry'
kind_control_plane='kind-control-plane'
BITBUCKET_POSTGRES_CONTAINER_NAME="ods-test-bitbucket-postgres"
BITBUCKET_SERVER_CONTAINER_NAME="ods-test-bitbucket-server"
NEXUS_CONTAINER_NAME="ods-test-nexus"
SQ_CONTAINER_NAME="ods-test-sonarqube"

container_names_in_start_order=( "$kind_registry" "$kind_control_plane" "$BITBUCKET_POSTGRES_CONTAINER_NAME" "$BITBUCKET_SERVER_CONTAINER_NAME" "$NEXUS_CONTAINER_NAME" 
    "$SQ_CONTAINER_NAME" ) 

for cn in "${container_names_in_start_order[@]}"; do
    echo docker start "$cn"
    docker start "$cn"
done

echo "Waiting for tools to start..."
echo "If this times out you can run this script again."

"$SCRIPT_DIR/waitfor-bitbucket.sh"
"$SCRIPT_DIR/waitfor-nexus.sh"
"$SCRIPT_DIR/waitfor-sonarqube.sh"

echo "Please start k9s and see pods are all ready before using cluster."
