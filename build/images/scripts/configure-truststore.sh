#!/bin/bash
set -u

md5_bin="${MD5_BIN:-"md5sum"}"
private_cert="/etc/ssl/certs/private-cert.pem"
src_truststore="${JAVA_HOME}/lib/security/cacerts"
src_pass="changeit"
dest_pass="changeit"

while [ "$#" -gt 0 ]; do
    case $1 in

    --src-store) src_truststore="$2"; shift;;
    --src-store=*) src_truststore="${1#*=}";;

    --src-storepass) src_pass="$2"; shift;;
    --src-storepass=*) src_pass="${1#*=}";;

    --dest-store) dest_truststore="$2"; shift;;
    --dest-store=*) dest_truststore="${1#*=}";;

    --dest-storepass) dest_pass="$2"; shift;;
    --dest-storepass=*) dest_pass="${1#*=}";;

    --debug) set -x;;

  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

dest_truststore_dir="${dest_truststore%/*}"
mkdir -p "${dest_truststore_dir}"
md5_private_cert_path="${dest_truststore_dir}/.md5-private-cert"
md5_private_cert=$(${md5_bin} < "${private_cert}" | cut -d- -f1)

if [ ! -f "${dest_truststore}" ] || [ "${md5_private_cert}" != "$(cat "${md5_private_cert_path}")" ]; then
    echo "Creating truststore with private cert ..."
    # Copy global keystone to location where we can write to (hide output containing warnings).
    if [ -f "${dest_truststore}" ]; then
        rm "${dest_truststore}"
    fi
    keytool -importkeystore \
        -srckeystore "${src_truststore}" -destkeystore "${dest_truststore}" \
        -deststorepass "${dest_pass}" -srcstorepass "${src_pass}" &> keytool-output.txt
    # shellcheck disable=SC2181
    if [ $? -ne 0 ]; then
        echo "error importing keystore:"
        cat keytool-output.txt; exit 1
    fi
    # Trust private cert (hide output containing warnings).
    keytool -importcert -noprompt -trustcacerts \
        -alias private-cert -file "${private_cert}" \
        -keystore "${dest_truststore}" -storepass "${dest_pass}" &> keytool-output.txt
    # shellcheck disable=SC2181
    if [ $? -ne 0 ]; then
        echo "error importing cert:"
        cat keytool-output.txt; exit 1
    fi
    echo "${md5_private_cert}" > "${md5_private_cert_path}"
else
    echo "Trustore with private cert exists already and is up-to-date."
fi
