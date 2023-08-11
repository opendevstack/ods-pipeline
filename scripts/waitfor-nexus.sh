#!/usr/bin/env bash
set -ue

INSECURE=""
HOST_PORT="8081"
NEXUS_URL=

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    -i|--insecure) INSECURE="--insecure";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

NEXUS_URL="http://localhost:${HOST_PORT}"
function waitForReady {
    echo "Waiting up to 5 minutes for Nexus to start ..."
    local n=0
    local http_code=
    set +e
    until [ $n -ge 30 ]; do
        http_code=$(curl ${INSECURE} -s -o /dev/null -w "%{http_code}" "${NEXUS_URL}/service/rest/v1/status/writable")
        if [ "${http_code}" == "200" ]; then
            echo " success"
            break
        else
            echo -n "."
            sleep 10
            n=$((n+1))
        fi
    done
    set -e

    if [ "${http_code}" != "200" ]; then
        echo "Nexus did not start, got http_code=${http_code}."
        exit 1
    fi
}

waitForReady
