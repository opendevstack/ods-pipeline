#!/usr/bin/env bash
set -ue

kind_registry='kind-registry'
kind_control_plane='kind-control-plane'
BITBUCKET_POSTGRES_CONTAINER_NAME="bitbucket-postgres-test"
BITBUCKET_SERVER_CONTAINER_NAME="bitbucket-server-test"
NEXUS_CONTAINER_NAME="nexustest"
SQ_CONTAINER_NAME="sonarqubetest"

container_names_in_stop_order=( "$SQ_CONTAINER_NAME" "$NEXUS_CONTAINER_NAME" "$BITBUCKET_SERVER_CONTAINER_NAME" "$BITBUCKET_POSTGRES_CONTAINER_NAME" "$kind_control_plane" "$kind_registry" ) 

for cn in "${container_names_in_stop_order[@]}"; do
    echo docker stop "$cn"
    docker stop "$cn"
done
