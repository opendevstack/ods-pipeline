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

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is not installed. Please see https://kind.sigs.k8s.io/"
  exit 1
fi

# desired cluster name; default is "kind"
kind_cluster_name="ods-pipeline"
recreate_kind_cluster="false"
registry_port="5000"
kind_mount_path="/tmp/ods-pipeline/kind-mount"

# K8S version is aligned with OpenShift GA 4.11.
# See https://docs.openshift.com/container-platform/4.11/release_notes/ocp-4-11-release-notes.html
k8s_version="v1.24.7"

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    --name) kind_cluster_name="$2"; shift;;
    --name=*) kind_cluster_name="${1#*=}";;

    --recreate) recreate_kind_cluster="true";;

    --registry-port) registry_port="$2"; shift;;
    --registry-port=*) registry_port="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

registry_name="${kind_cluster_name}-registry"
reg_ip_selector='{{.NetworkSettings.Networks.kind.IPAddress}}'
reg_network='kind'

# create registry container unless it already exists
running="$(docker inspect -f '{{.State.Running}}' "${registry_name}" 2>/dev/null || true)"

# If the registry already exists, but is in the wrong network, we have to
# re-create it.
if [ "${running}" = 'true' ]; then
  reg_ip="$(docker inspect -f ${reg_ip_selector} "${registry_name}")"
  if [ "${reg_ip}" = '' ]; then
    docker kill "${registry_name}"
    docker rm "${registry_name}"
    running="false"
  fi
fi

if [ "${running}" != 'true' ]; then
  net_driver=$(docker network inspect "${reg_network}" -f '{{.Driver}}' || true)
  if [ "${net_driver}" != "bridge" ]; then
    docker network create "${reg_network}"
  fi
  if docker inspect "${registry_name}" >/dev/null 2>&1; then
    docker rm "${registry_name}"
  fi
  docker run \
    -d --restart=always -p "${registry_port}:5000" --name "${registry_name}" --net "${reg_network}" \
    registry:2
fi

reg_ip="$(docker inspect -f ${reg_ip_selector} "${registry_name}")"
if [ "${reg_ip}" = "" ]; then
    echo "Error creating registry: no IPAddress found at: ${reg_ip_selector}"
    exit 1
fi

if [ "${recreate_kind_cluster}" == "false" ]; then
  if kind get clusters | grep "${kind_cluster_name}" >/dev/null 2>&1; then
    echo "Reusing existing cluster ..."
    exit 0
  fi 
fi

if [ "${recreate_kind_cluster}" == "true" ]; then
  kind delete cluster --name "${kind_cluster_name}"
fi

# create a cluster with the local registry enabled in containerd
mkdir -p ${kind_mount_path}
chmod -R 0755 ${kind_mount_path}
cat <<EOF | kind create cluster --name "${kind_cluster_name}" --image "kindest/node:${k8s_version}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraMounts:
  - hostPath: ${kind_mount_path}
    containerPath: /files
  kubeadmConfigPatches:
  - |
    kind: KubeletConfiguration
    cgroupDriver: systemd
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${registry_port}"]
    endpoint = ["http://${reg_ip}:${registry_port}"]
EOF

for node in $(kind get nodes --name "${kind_cluster_name}"); do
  kubectl annotate node "${node}" tilt.dev/registry=localhost:"${registry_port}";

  # Add the registry config to the nodes.
  # This is necessary because localhost resolves to loopback addresses that are
  # network-namespace local.
  # In other words: localhost in the container is not localhost on the host.
  # We want a consistent name that works from both ends, so we tell containerd to
  # alias localhost:${reg_port} to the registry container when pulling images.
  registry_dir="/etc/containerd/certs.d/localhost:${registry_port}"
  docker exec "${node}" mkdir -p "${registry_dir}"
  cat <<EOF | docker exec -i "${node}" cp /dev/stdin "${registry_dir}/hosts.toml"
[host."http://${registry_name}:5000"]
EOF
done
