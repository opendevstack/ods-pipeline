#!/bin/bash
set -eu

md5bin="${MD5_BIN:-"md5sum --tag"}"
debug="${DEBUG:-false}"
aquaScannerUrl=""
binDir=".ods-cache/bin"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --bin-dir) binDir="$2"; shift;;
    --bin-dir=*) binDir="${1#*=}";;

    --aqua-scanner-url) aquaScannerUrl="$2"; shift;;
    --aqua-scanner-url=*) aquaScannerUrl="${1#*=}";;

    --debug) debug="$2"; shift;;
    --debug=*) debug="${1#*=}";;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ "${debug}" == "true" ]; then
  set -x
fi

aquaScannerPath="${binDir}/aquasec"
md5AquaScannerUrlPath="${binDir}/.md5-aquasec"

# Optionally install Aqua scanner.
# If the binary already exists and was downloaded from the
# URL given by aquaScannerUrl, skip download.
if [ -n "${aquaScannerUrl}" ] && [ "${aquaScannerUrl}" != "none" ]; then
    md5AquaScannerUrl=$(${md5bin} -s "${aquaScannerUrl}")
    if [ ! -f "${md5AquaScannerUrlPath}" ] || [ "${md5AquaScannerUrl}" != "$(cat "${md5AquaScannerUrlPath}")" ]; then
        echo 'Installing Aqua scanner...'
        curl -v -sSf -L "${aquaScannerUrl}" -o aquasec
        mv aquasec "${aquaScannerPath}"
        chmod +x "${aquaScannerPath}"
        echo "${md5AquaScannerUrl}" > "${md5AquaScannerUrlPath}"
        echo 'Installed Aqua scanner version:'
        "${aquaScannerPath}" version
    fi
fi
