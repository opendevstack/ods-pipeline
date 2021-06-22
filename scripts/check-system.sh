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

if [ "${ok}" == 'true' ]; then
    echo "Please make sure that the Docker host has at least 8GB memory (no automated check exists for this yet)."
    echo ""

    if ! which k9s &> /dev/null; then
        echo "Note: k9s is recommended."
    fi

    echo "OK"
else
    echo "NOT OK"
    exit 1
fi
