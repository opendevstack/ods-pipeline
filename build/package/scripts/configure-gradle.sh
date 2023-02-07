#!/bin/bash

# This script checks for env variable HTTP_PROXY and adds them to gradle.properties.
CONTENT=""

if [ -f /etc/ssl/certs/private-cert.pem ]; then
	truststore_location="$(pwd)/.ods-cache/truststore/cacerts"
	truststore_pass="changeit"
	echo "Configuring Gradle to trust private cert ..."
	configure-truststore --dest-store="${truststore_location}" --dest-storepass="${truststore_pass}"
	# shellcheck disable=SC2181
	if [ $? -ne 0 ]; then
		exit 1
	fi
	# Configure Gradle to use the modified trust store.
	CONTENT+="systemProp.javax.net.ssl.trustStore=${truststore_location}\n"
	CONTENT+="systemProp.javax.net.ssl.trustStorePassword=${truststore_pass}\n"
fi

if [ "${HTTP_PROXY}" != "" ]; then
	echo "Configuring Gradle to honor HTTP_PROXY ..."
	proxy=$(echo "$HTTP_PROXY" | sed -e "s|https://||g" | sed -e "s|http://||g")
	proxy_hostp=$(echo "$proxy" | cut -d "@" -f2)

	CONTENT+="systemProp.proxySet=\"true\"\n"

	proxy_host=$(echo "$proxy_hostp" | cut -d ":" -f1)
	CONTENT+="systemProp.http.proxyHost=${proxy_host}\n"
	CONTENT+="systemProp.https.proxyHost=${proxy_host}\n"

	proxy_port=$(echo "$proxy_hostp" | cut -d ":" -f2)
	CONTENT+="systemProp.http.proxyPort=${proxy_port}\n"
	CONTENT+="systemProp.https.proxyPort=${proxy_port}\n"

	proxy_userp=$(echo "$proxy" | cut -d "@" -f1)
	if [[ $proxy_userp != "$proxy_hostp" ]];
	then
		proxy_user=$(echo "$proxy_userp" | cut -d ":" -f1)
		CONTENT+="systemProp.http.proxyUser=${proxy_user}\n"
		CONTENT+="systemProp.https.proxyUser=${proxy_user}\n"

		# shellcheck disable=SC2001
		proxy_pw=$(echo "$proxy_userp" | sed -e "s|$proxy_user:||g")
		CONTENT+="systemProp.http.proxyPassword=${proxy_pw}\n"
		CONTENT+="systemProp.https.proxyPassword=${proxy_pw}\n"
 	fi
fi

if [ "${NO_PROXY}" != "" ]; then
	echo "Configuring Gradle to honor NO_PROXY ..."
	# shellcheck disable=SC2001
	noproxy_host=$(echo "$NO_PROXY" | sed -e 's|\,\.|\,\*\.|g')
	# shellcheck disable=SC2001
	noproxy_host=$(echo "$noproxy_host" | sed -e "s/,/|/g")
	CONTENT+="systemProp.http.nonProxyHosts=$noproxy_host\n"
	CONTENT+="systemProp.https.nonProxyHosts=$noproxy_host\n"
fi

if [ "${CONTENT}" != "" ]; then
  echo -e "$CONTENT" > "${GRADLE_USER_HOME}/gradle.properties"
fi
