#!/usr/bin/env bash
set -ue

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ods_pipeline_dir=${script_dir%/*}

https_port="8443"
container_name=""
nginx_conf=""

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    --container-name) container_name="$2"; shift;;
    --container-name=*) container_name="${1#*=}";;

    --nginx-conf) nginx_conf="$2"; shift;;
    --nginx-conf=*) nginx_conf="${1#*=}";;

    --https-port) https_port="$2"; shift;;
    --https-port=*) https_port="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

nginx_image="nginx:alpine-slim"
private_cert_dir="${ods_pipeline_dir}/test/testdata/private-cert"
if [ "$(uname -m)" == "arm64" ]; then
    nginx_image="arm64v8/nginx:alpine-slim"
fi
docker rm -f "${container_name}" &> /dev/null || true
docker run --name "${container_name}" \
  -v "${script_dir}/nginx/${nginx_conf}:/etc/nginx/nginx.conf:ro" \
  -v "${private_cert_dir}/tls.crt:/etc/nginx/tls.crt:ro" \
  -v "${private_cert_dir}/tls.key:/etc/nginx/tls.key:ro" \
  -d --net kind -p "${https_port}:${https_port}" "${nginx_image}"
