#!/usr/bin/env bash
set -ue

# Starts a Bitbucket instance with a timebomb license (3 hours).
# The instance is setup with an admin account (pw: admin) and an "ODSPIPELINETEST" project.

INSECURE=""
BITBUCKET_SERVER_HOST_PORT="7990"

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

BITBUCKET_URL="http://localhost:${BITBUCKET_SERVER_HOST_PORT}"

echo "Waiting up to 5 minutes for Bitbucket to start ..."
# https://confluence.atlassian.com/bitbucketserverkb/how-to-monitor-if-bitbucket-server-is-up-and-running-975014635.html
n=0
status="STARTING"
set +e
until [ $n -ge 30 ]; do
    status=$(curl -s "${INSECURE}" "${BITBUCKET_URL}/status" | jq -r .state)
    if [ "${status}" == "RUNNING" ]; then
        echo " success"
        break
    else
        echo -n "."
        sleep 5
        n=$((n+1))
    fi
done
set -e
if [ "${status}" != "RUNNING" ]; then
    echo "Bitbucket did not start, got status=${status}."
    exit 1
fi
exit 0
