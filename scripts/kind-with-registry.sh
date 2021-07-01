#!/bin/bash
#	
# Adapted from:	
# https://github.com/kubernetes-sigs/kind/commits/master/site/static/examples/kind-with-registry.sh	
#	
# Copyright 2020 The Kubernetes Project	
#	
# Licensed under the Apache License, Version 2.0 (the "License");	
# you may not use this file except in compliance with the License.	
# You may obtain a copy of the License at	
#	
#     http://www.apache.org/licenses/LICENSE-2.0	
#	
# Unless required by applicable law or agreed to in writing, software	
# distributed under the License is distributed on an "AS IS" BASIS,	
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.	
# See the License for the specific language governing permissions and	
# limitations under the License.

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
case "${kind_version}" in
  "kind v0.7."* | "kind v0.6."* | "kind v0.5."*)
    reg_ip_selector='{{.NetworkSettings.IPAddress}}'
    reg_network='bridge'
    ;;
esac

# create registry container unless it already exists
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"

# If the registry already exists, but is in the wrong network, we have to
# re-create it.
if [ "${running}" = 'true' ]; then
  reg_ip="$(docker inspect -f ${reg_ip_selector} "${reg_name}")"
  if [ "${reg_ip}" = '' ]; then
    docker kill ${reg_name}
    docker rm ${reg_name}
    running="false"
  fi
fi

if [ "${running}" != 'true' ]; then
  if [ "${reg_network}" != "bridge" ]; then
    docker network create "${reg_network}" || true
  fi
  
  docker run \
    -d --restart=always -p "${reg_port}:5000" --name "${reg_name}" --net "${reg_network}" \
    registry:2
fi

reg_ip="$(docker inspect -f ${reg_ip_selector} "${reg_name}")"
if [ "${reg_ip}" = "" ]; then
    echo "Error creating registry: no IPAddress found at: ${reg_ip_selector}"
    exit 1
fi
echo "Registry IP: ${reg_ip}"

if [ "${RECREATE_KIND_CLUSTER}" == "true" ]; then
  kind delete cluster --name "${KIND_CLUSTER_NAME}"
fi

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name "${KIND_CLUSTER_NAME}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30950
    hostPort: 30950
  # add a mount from /path/to/my/files on the host to /files on the node
  extraMounts:
  - hostPath: ${ODS_PIPELINE_DIR}/test
    containerPath: /files
containerdConfigPatches: 
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_ip}:${reg_port}"]
EOF

for node in $(kind get nodes --name "${KIND_CLUSTER_NAME}"); do
  kubectl annotate node "${node}" tilt.dev/registry=localhost:${reg_port};
done
