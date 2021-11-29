#!/usr/bin/env bash
set -ue

# Checks whether the system has all prerequisites to run the tests.

ok="true"

if ! which go &> /dev/null; then
    ok="false"
    echo "go is required. For documentation on how to install please visit https://go.dev/doc/install."
fi

if ! which docker &> /dev/null; then
    ok="false"
    echo "docker is required. For documentation on how to install please visit https://docs.docker.com/engine/install/."
fi

if ! which jq &> /dev/null; then
    ok="false"
    echo "jq is required. For documentation on how to install please visit https://stedolan.github.io/jq/download/."
fi

if ! which kind &> /dev/null; then
    ok="false"
    echo "kind is required. For documentation on how to install please visit https://kind.sigs.k8s.io/#installation-and-usage."
fi

if ! which kubectl &> /dev/null; then
    ok="false"
    echo "kubectl is required. For documentation on how to install please visit https://kubernetes.io/docs/tasks/tools/#kubectl."
fi

if ! which helm &> /dev/null; then
    ok="false"
    echo "helm is required. For documentation on how to install please visit https://helm.sh/docs/intro/install/."
fi

if which docker &> /dev/null; then
    docker_host_memory=$(docker info --format "{{.MemTotal}}")
    if [ ${docker_host_memory} -lt $((8 * 10 ** 9)) ]; then
        ok="false"
        echo "The Docker host must have at least 8GB of memory."
    fi
fi


if [ "${ok}" == 'true' ]; then

    if ! which k9s &> /dev/null; then
        echo "Note: k9s is recommended."
    fi

    echo "OK"
else
    echo "NOT OK"
    exit 1
fi
