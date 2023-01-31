#!/bin/bash
set -eu

md5_bin="${MD5_BIN:-"md5sum --tag"}"
aqua_scanner_url=""
bin_dir=".ods-cache/bin"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    --bin-dir) bin_dir="$2"; shift;;
    --bin-dir=*) bin_dir="${1#*=}";;

    --aqua-scanner-url) aqua_scanner_url="$2"; shift;;
    --aqua-scanner-url=*) aqua_scanner_url="${1#*=}";;

    --debug) set -x; shift;;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

aqua_scanner_path="${bin_dir}/aquasec"
md5_aqua_scanner_url_path="${bin_dir}/.md5-aquasec"

# Optionally install Aqua scanner.
# If the binary already exists and was downloaded from the
# URL given by aqua_scanner_url, skip download.
if [ -n "${aqua_scanner_url}" ] && [ "${aqua_scanner_url}" != "none" ]; then
    md5_aqua_scanner_url=$(${md5_bin} -s "${aqua_scanner_url}")
    if [ ! -f "${md5_aqua_scanner_url_path}" ] || [ "${md5_aqua_scanner_url}" != "$(cat "${md5_aqua_scanner_url_path}")" ]; then
        echo 'Installing Aqua scanner...'
        curl -v -sSf -L "${aqua_scanner_url}" -o aquasec
        mv aquasec "${aqua_scanner_path}"
        chmod +x "${aqua_scanner_path}"
        echo "${md5_aqua_scanner_url}" > "${md5_aqua_scanner_url_path}"
        echo 'Installed Aqua scanner version:'
        "${aqua_scanner_path}" version
    fi
fi
