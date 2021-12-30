#!/usr/bin/env bash
set -ue

kind_registry='kind-registry'
kind_control_plane='kind-control-plane'
BITBUCKET_POSTGRES_CONTAINER_NAME="ods-test-bitbucket-postgres"
BITBUCKET_SERVER_CONTAINER_NAME="ods-test-bitbucket-server"
NEXUS_CONTAINER_NAME="ods-test-nexus"
SQ_CONTAINER_NAME="ods-test-sonarqube"

container_names_in_stop_order=( "$SQ_CONTAINER_NAME" "$NEXUS_CONTAINER_NAME" "$BITBUCKET_SERVER_CONTAINER_NAME" "$BITBUCKET_POSTGRES_CONTAINER_NAME" "$kind_control_plane" "$kind_registry" ) 

for cn in "${container_names_in_stop_order[@]}"; do
    echo docker stop "$cn"
    docker stop "$cn"  || true
done
