#!/usr/bin/env bash
set -ue

# Checks whether the system has all prerequisites to run the tests.

ok="true"

for exe in go docker jq kind kubectl helm; do
    if ! which $exe &> /dev/null; then
        ok="false"
        echo "$exe is required"
    fi
done

if which docker &> /dev/null; then
    docker_host_memory=$(docker info --format "{{.MemTotal}}")
    if [ ! ${docker_host_memory} -gt $((8 * 10 ** 9)) ]; then
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
