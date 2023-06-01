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
KIND_CLUSTER_NAME="kind"
RECREATE_KIND_CLUSTER="false"
REGISTRY_PORT="5000"

# K8S version is aligned with OpenShift GA 4.11.
# See https://docs.openshift.com/container-platform/4.11/release_notes/ocp-4-11-release-notes.html
K8S_VERSION="v1.24.7"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    --name) KIND_CLUSTER_NAME="$2"; shift;;
    --name=*) KIND_CLUSTER_NAME="${1#*=}";;

    --recreate) RECREATE_KIND_CLUSTER="true";;

    --registry-port) REGISTRY_PORT="$2"; shift;;
    --registry-port=*) REGISTRY_PORT="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

kind_version=$(kind version)
REGISTRY_NAME="${KIND_CLUSTER_NAME}-registry"
reg_ip_selector='{{.NetworkSettings.Networks.kind.IPAddress}}'
reg_network='kind'
case "${kind_version}" in
  "kind v0.7."* | "kind v0.6."* | "kind v0.5."*)
    reg_ip_selector='{{.NetworkSettings.IPAddress}}'
    reg_network='bridge'
    ;;
esac

# create registry container unless it already exists
running="$(docker inspect -f '{{.State.Running}}' "${REGISTRY_NAME}" 2>/dev/null || true)"

# If the registry already exists, but is in the wrong network, we have to
# re-create it.
if [ "${running}" = 'true' ]; then
  reg_ip="$(docker inspect -f ${reg_ip_selector} "${REGISTRY_NAME}")"
  if [ "${reg_ip}" = '' ]; then
    docker kill "${REGISTRY_NAME}"
    docker rm "${REGISTRY_NAME}"
    running="false"
  fi
fi

if [ "${running}" != 'true' ]; then
  if [ "${reg_network}" != "bridge" ]; then
    docker network create "${reg_network}" || true
  fi

  docker run \
    -d --restart=always -p "${REGISTRY_PORT}:5000" --name "${REGISTRY_NAME}" --net "${reg_network}" \
    registry:2
fi

reg_ip="$(docker inspect -f ${reg_ip_selector} "${REGISTRY_NAME}")"
if [ "${reg_ip}" = "" ]; then
    echo "Error creating registry: no IPAddress found at: ${reg_ip_selector}"
    exit 1
fi
echo "Registry IP: ${reg_ip}"

if [ "${RECREATE_KIND_CLUSTER}" == "true" ]; then
  kind delete cluster --name "${KIND_CLUSTER_NAME}"
fi

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --name "${KIND_CLUSTER_NAME}" --image "kindest/node:${K8S_VERSION}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  # add a mount from /path/to/my/files on the host to /files on the node
  extraMounts:
  - hostPath: ${ODS_PIPELINE_DIR}/test
    containerPath: /files
  kubeadmConfigPatches:
  - |
    kind: KubeletConfiguration
    cgroupDriver: systemd
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${REGISTRY_PORT}"]
    endpoint = ["http://${reg_ip}:${REGISTRY_PORT}"]
EOF

for node in $(kind get nodes --name "${KIND_CLUSTER_NAME}"); do
  kubectl annotate node "${node}" tilt.dev/registry=localhost:"${REGISTRY_PORT}";

  # Add the registry config to the nodes.
  # This is necessary because localhost resolves to loopback addresses that are
  # network-namespace local.
  # In other words: localhost in the container is not localhost on the host.
  # We want a consistent name that works from both ends, so we tell containerd to
  # alias localhost:${reg_port} to the registry container when pulling images.
  REGISTRY_DIR="/etc/containerd/certs.d/localhost:${REGISTRY_PORT}"
  docker exec "${node}" mkdir -p "${REGISTRY_DIR}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${REGISTRY_DIR}/hosts.toml"
[host."http://${REGISTRY_NAME}:5000"]
EOF
done
