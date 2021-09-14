#!/bin/bash

set -o errexit

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

# desired cluster name; default is "kind"
KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-kind}"
RECREATE_KIND_CLUSTER="false"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    --recreate) RECREATE_KIND_CLUSTER="true";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

kind_version=$(kind version)
reg_name='kind-registry'
reg_port='5000'
reg_ip_selector='{{.NetworkSettings.Networks.kind.IPAddress}}'
reg_network='kind'

echo "Removing kind cluster ${KIND_CLUSTER_NAME} ..."
kind delete cluster --name "${KIND_CLUSTER_NAME}"
echo "... done!"
echo

echo "Stopping all docker containers in docker network $(reg_network) ..."
docker ps -qf "network=kind"
docker stop $(docker ps -qf "network=kind")
echo "... done!"
echo

echo "Removing all docker images in docker network $(reg_network) ..."
docker ps -af "network=kind"
docker rm $(docker ps -aqf "network=kind")
echo "... done!"
echo
