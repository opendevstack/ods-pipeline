#!/usr/bin/env bash
set -ue

# Restarts a Bitbucket instance with a timebomb license (3 hours).

INSECURE=""
BITBUCKET_SERVER_HOST_PORT="7990"
BITBUCKET_SERVER_CONTAINER_NAME="ods-test-bitbucket-server"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done


echo "Restart Bitbucket Server"
docker stop ${BITBUCKET_SERVER_CONTAINER_NAME}
docker start ${BITBUCKET_SERVER_CONTAINER_NAME}

BITBUCKET_URL="http://localhost:${BITBUCKET_SERVER_HOST_PORT}"
echo "Waiting up to 3 minutes for Bitbucket to start ..."
# https://confluence.atlassian.com/bitbucketserverkb/how-to-monitor-if-bitbucket-server-is-up-and-running-975014635.html
n=0
status="STARTING"
set +e
until [ $n -ge 18 ]; do
    status=$(curl -s ${INSECURE} "${BITBUCKET_URL}/status" | jq -r .state)
    if [ "${status}" == "RUNNING" ]; then
        echo " success"
        break
    else
        echo -n "."
        sleep 10
        n=$((n+1))
    fi
done
set -e
if [ "${status}" != "RUNNING" ]; then
    echo "Bitbucket did not start, got status=${status}."
    exit 1
fi
