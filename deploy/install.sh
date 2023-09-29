#!/usr/bin/env bash
set -ue

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

VERBOSE="false"
DRY_RUN="false"
DIFF="true"
NAMESPACE=""
RELEASE_NAME="ods-pipeline"
SERVICEACCOUNT="pipeline"
VALUES_FILE="values.yaml"
CHART_DIR="./chart"
# Secrets
AUTH_SEPARATOR=":"
BITBUCKET_AUTH=""
BITBUCKET_WEBHOOK_SECRET=""
NEXUS_AUTH=""
PRIVATE_CERT=""

# Check prerequisites.
KUBECTL_BIN=""
if command -v oc &> /dev/null; then
    KUBECTL_BIN="oc"
elif command -v kubectl &> /dev/null; then
    KUBECTL_BIN="kubectl"
else
    echo "ERROR: Neither oc nor kubectl in \$PATH"; exit 1
fi
HELM_BIN=""
if command -v helm &> /dev/null; then
    HELM_BIN="helm"
else
    echo "ERROR: helm is not in \$PATH"; exit 1
fi

function usage {
    printf "Usage:\n"
    printf "\t-h|--help\t\t\tPrints this usage information.\n"
    printf "\t-v|--verbose\t\t\tTurn on verbose output.\n"
    printf "\t-n|--namespace\t\t\tK8s namespace to target.\n"
    printf "\t-f|--values\t\t\tValues file to supply to Helm (defaults to '%s'). Multiple files can be specified comma-separated.\n" "$VALUES_FILE"
    printf "\t-s|--serviceaccount\t\tServiceaccount to use (defaults to '%s').\n" "$SERVICEACCOUNT"
    printf "\t--no-diff\t\t\tDo not run Helm diff before running Helm upgrade.\n"
    printf "\t--dry-run\t\t\tDo not apply any changes, instead just print what the script would do.\n"
    printf "\t--auth-separator\t\tCharacter to use as a separator for basic auth flags (defaults to '%s')\n" "$AUTH_SEPARATOR"
    printf "\t--bitbucket-auth\t\tAccess token of a Bitbucket user (if not given, script will prompt for this).\n"
    printf "\t--bitbucket-webhook-secret\tSecret to protect webhook endpoint with (if not given, script will generate this).\n"
    printf "\t--nexus-auth\t\t\tUsername and password (separated by '%s') of a Nexus user (if not given, script will prompt for this).\n" "$AUTH_SEPARATOR"
    printf "\t--private-cert\t\t\tHost from which to download private certificate (if not given, script will skip this).\n"
    printf "\nExample:\n\n"
    printf "\t%s \ \
      \n\t\t--namespace foo \ \
      \n\t\t--bitbucket-auth 'personal-access-token' \ \
      \n\t\t--nexus-auth 'user:password' \n\n" "$0"
}

while [ "$#" -gt 0 ]; do
    # shellcheck disable=SC2034
    case $1 in

    -h|--help) shift; usage; exit 0;;

    -v|--verbose) VERBOSE="true";;

    -n|--namespace) NAMESPACE="$2"; shift;;
    -n=*|--namespace=*) NAMESPACE="${1#*=}";;

    -f|--values) VALUES_FILE="$2"; shift;;
    -f=*|--values=*) VALUES_FILE="${1#*=}";;

    -s|--serviceaccount) SERVICEACCOUNT="$2"; shift;;
    -s=*|--serviceaccount=*) SERVICEACCOUNT="${1#*=}";;

    --no-diff) DIFF="false";;

    --dry-run) DRY_RUN="true";;

    --auth-separator) AUTH_SEPARATOR="$2"; shift;;
    --auth-separator=*) AUTH_SEPARATOR="${1#*=}";;

    --bitbucket-auth) BITBUCKET_AUTH="$2"; shift;;
    --bitbucket-auth=*) BITBUCKET_AUTH="${1#*=}";;

    --bitbucket-webhook-secret) BITBUCKET_WEBHOOK_SECRET="$2"; shift;;
    --bitbucket-webhook-secret=*) BITBUCKET_WEBHOOK_SECRET="${1#*=}";;

    --nexus-auth) NEXUS_AUTH="$2"; shift;;
    --nexus-auth=*) NEXUS_AUTH="${1#*=}";;

    --private-cert) PRIVATE_CERT="$2"; shift;;
    --private-cert=*) PRIVATE_CERT="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd "${SCRIPT_DIR}"

VALUES_FILES=$(echo "$VALUES_FILE" | tr "," "\n")
VALUES_ARGS=()
for valueFile in ${VALUES_FILES}; do
    VALUES_ARGS+=(--values="${valueFile}")
done

if [ "${VERBOSE}" == "true" ]; then
    set -x
fi

if [ -z "${NAMESPACE}" ]; then
    echo "--namespace is required"
    exit 1
fi

kubectlApplySecret () {
    local secretName="$1"
    local secretTemplate="$2"
    local username="$3"
    local password="$4"
    # Render variables into template, then apply.
    # To avoid forward slashes messing up sed, escape forward slashes first.
    # See https://tldp.org/LDP/abs/html/string-manipulation.html.
    # shellcheck disable=SC2002
    cat "${secretTemplate}" | sed "s/{{name}}/${secretName}/" | sed "s/{{username}}/${username//\//\\/}/" | sed "s/{{password}}/${password//\//\\/}/" | "${KUBECTL_BIN}" -n "${NAMESPACE}" apply -f -
}

installSecret () {
    local secretName="$1"
    local secretTemplate="$2"
    local flagValue="$3"
    local usernamePrompt="$4"
    local passwordPrompt="$5"

    # Split flag value on first occurence of auth separator.
    local authUser="${flagValue%%"${AUTH_SEPARATOR}"*}"
    local authPassword="${flagValue#*"${AUTH_SEPARATOR}"}"

    # If the secret exists and the flag is present, update the secret.
    if "${KUBECTL_BIN}" -n "${NAMESPACE}" get "secret/${secretName}" &> /dev/null; then
        # In case the secret was previously managed by Helm, we want to instruct Helm
        # to keep the resource during helm upgrade.
        "${KUBECTL_BIN}" -n "${NAMESPACE}" annotate --overwrite secret "${secretName}" "helm.sh/resource-policy=keep"
        if [ -n "${flagValue}" ]; then
            echo "Updating secret ${secretName} ..."
            kubectlApplySecret "${secretName}" "${secretTemplate}" "${authUser}" "${authPassword}"
        else
            echo "Secret ${secretName} already exists and will not be updated."
        fi
    else # If the secret does not exist, create the secret. If the flag is not present, ask for values.
        if [ -z "${flagValue}" ]; then
            if [ -n "${usernamePrompt}" ]; then
                echo "${usernamePrompt}"
                read -r authUser
            fi
            if [ -n "${passwordPrompt}" ]; then
                echo "${passwordPrompt}"
                read -rs authPassword
            else
                authUser=''
                authPassword=$(LC_ALL=C tr -dc 'A-Za-z0-9#%*+\-<=>_{|}' </dev/urandom | head -c 32 ; echo)
            fi
        fi
        echo "Creating secret ${secretName} ..."
        kubectlApplySecret "${secretName}" "${secretTemplate}" "${authUser}" "${authPassword}"
    fi
}

installTLSSecret () {
    local secretName="$1"
    local privateCert="$2"
    local certFile=""
    if [ -z "${privateCert}" ]; then
        echo "No private cert given, skipping ..."
    else
        if [ "${privateCert:0:1}" == '/' ] || [ "${privateCert:0:2}" == './' ]; then
            if [ ! -f "${privateCert}" ]; then
                echo "No cert file exists at ${privateCert}"; exit 1
            fi
            certFile="${privateCert}"
        else
            certFile="private-cert.pem.tmp"
            openssl s_client -showcerts -connect "${privateCert}" </dev/null \
                | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > "${certFile}"
        fi
        if "${KUBECTL_BIN}" -n "${NAMESPACE}" get "secret/${secretName}" &> /dev/null; then
            echo "Re-creating secret ${secretName} ..."
            "${KUBECTL_BIN}" -n "${NAMESPACE}" delete secret "${secretName}"
        else
            echo "Creating secret ${secretName} ..."
        fi
        "${KUBECTL_BIN}" -n "${NAMESPACE}" create secret generic "${secretName}" \
            --from-file=tls.crt="${certFile}"
        rm private-cert.pem.tmp &>/dev/null || true
    fi
}

# Manage serviceaccount ...
if "${KUBECTL_BIN}" -n "${NAMESPACE}" get serviceaccount/"${SERVICEACCOUNT}" &> /dev/null; then
    echo "Serviceaccount exists already ..."
else
    echo "Creating serviceaccount ..."
    if [ "${DRY_RUN}" == "true" ]; then
        echo "(skipping in dry-run)"
    else
        "${KUBECTL_BIN}" -n "${NAMESPACE}" create serviceaccount "${SERVICEACCOUNT}"

        "${KUBECTL_BIN}" -n "${NAMESPACE}" \
            create rolebinding "${SERVICEACCOUNT}-edit" \
            --clusterrole edit \
            --serviceaccount "${NAMESPACE}:${SERVICEACCOUNT}"
    fi
fi

echo "Installing secrets ..."
if [ "${DRY_RUN}" == "true" ]; then
    echo "(skipping in dry-run)"
else
    # Bitbucket username is not required as PAT alone is enough.
    installSecret "ods-bitbucket-auth" \
        "basic-auth-secret.yaml.tmpl" \
        "${BITBUCKET_AUTH}" \
        "Please enter the username of Bitbucket user with write permission." \
        "Please enter a personal access token of this Bitbucket user (input will be hidden):"

    # Webhook secret is a special case, as we do not want the user to set it.
    # No prompts -> password will be auto-generated if not given.
    installSecret "ods-bitbucket-webhook" \
        "opaque-secret.yaml.tmpl" \
        "${BITBUCKET_WEBHOOK_SECRET}" \
        "" ""

    installSecret "ods-nexus-auth" \
        "basic-auth-secret.yaml.tmpl" \
        "${NEXUS_AUTH}" \
        "Please enter the username of a Nexus user with write permission:" \
        "Please enter the password of this Nexus user (input will be hidden):"

    installTLSSecret "ods-private-cert" "${PRIVATE_CERT}"
fi

echo "Installing Helm release ${RELEASE_NAME} ..."
if [ "${DIFF}" == "true" ]; then
    if "${HELM_BIN}" -n "${NAMESPACE}" \
            diff upgrade --install --detailed-exitcode --three-way-merge --normalize-manifests \
            "${VALUES_ARGS[@]}" \
            ${RELEASE_NAME} ${CHART_DIR}; then
        echo "Helm release already up-to-date."
    else
        if [ "${DRY_RUN}" == "true" ]; then
            echo "(skipping in dry-run)"
        else
            "${HELM_BIN}" -n "${NAMESPACE}" \
                upgrade --install \
                "${VALUES_ARGS[@]}" \
                ${RELEASE_NAME} ${CHART_DIR}
        fi
    fi
else
    if [ "${DRY_RUN}" == "true" ]; then
        echo "(skipping in dry-run)"
    else
        "${HELM_BIN}" -n "${NAMESPACE}" \
            upgrade --install \
            "${VALUES_ARGS[@]}" \
            ${RELEASE_NAME} ${CHART_DIR}
    fi
fi

echo "Adding Tekton annotation to ods-bitbucket-auth secret ..."
if [ "${DRY_RUN}" == "true" ]; then
    echo "(skipping in dry-run)"
else
    bitbucketUrl=$("${KUBECTL_BIN}" -n "${NAMESPACE}" get cm/ods-bitbucket -ojsonpath='{.data.url}')
    "${KUBECTL_BIN}" -n "${NAMESPACE}" annotate --overwrite secret ods-bitbucket-auth "tekton.dev/git-0=${bitbucketUrl}"
fi

echo "Adding ods-bitbucket-auth secret to ${SERVICEACCOUNT} serviceaccount ..."
if [ "${DRY_RUN}" == "true" ]; then
    echo "(skipping in dry-run)"
else
    "${KUBECTL_BIN}" -n "${NAMESPACE}" \
        patch sa "${SERVICEACCOUNT}" \
        --type json \
        -p '[{"op": "add", "path": "/secrets", "value":[{"name": "ods-bitbucket-auth"}]}]'
fi
