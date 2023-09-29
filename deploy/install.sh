#!/usr/bin/env bash
set -ue

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

verbose="false"
dry_run="false"
diff="true"
namespace=""
release_name="ods-pipeline"
serviceaccount="pipeline"
values_file="values.yaml"
chart_dir="./chart"
# Secrets
auth_separator=":"
bitbucket_auth=""
bitbucket_webhook_secret=""
nexus_auth=""
private_cert=""
# Templates
basicAuthSecretTemplate="apiVersion: v1
kind: Secret
metadata:
  name: '{{name}}'
  labels:
    app.kubernetes.io/name: ods-pipeline
stringData:
  password: '{{password}}'
  username: '{{username}}'
type: kubernetes.io/basic-auth"
opaqueSecretTemplate="apiVersion: v1
kind: Secret
metadata:
  name: '{{name}}'
  labels:
    app.kubernetes.io/name: ods-pipeline
stringData:
  secret: '{{password}}'
type: Opaque"

# Check prerequisites.
kubectl_bin=""
if command -v oc &> /dev/null; then
    kubectl_bin="oc"
elif command -v kubectl &> /dev/null; then
    kubectl_bin="kubectl"
else
    echo "ERROR: Neither oc nor kubectl in \$PATH"; exit 1
fi
helm_bin=""
if command -v helm &> /dev/null; then
    helm_bin="helm"
else
    echo "ERROR: helm is not in \$PATH"; exit 1
fi

function usage {
    printf "Usage:\n"
    printf "\t-h|--help\t\t\tPrints this usage information.\n"
    printf "\t-v|--verbose\t\t\tTurn on verbose output.\n"
    printf "\t-n|--namespace\t\t\tK8s namespace to target.\n"
    printf "\t-f|--values\t\t\tValues file to supply to Helm (defaults to '%s'). Multiple files can be specified comma-separated.\n" "$values_file"
    printf "\t-s|--serviceaccount\t\tServiceaccount to use (defaults to '%s').\n" "$serviceaccount"
    printf "\t--no-diff\t\t\tDo not run Helm diff before running Helm upgrade.\n"
    printf "\t--dry-run\t\t\tDo not apply any changes, instead just print what the script would do.\n"
    printf "\t--auth-separator\t\tCharacter to use as a separator for basic auth flags (defaults to '%s')\n" "$auth_separator"
    printf "\t--bitbucket-auth\t\tAccess token of a Bitbucket user (if not given, script will prompt for this).\n"
    printf "\t--bitbucket-webhook-secret\tSecret to protect webhook endpoint with (if not given, script will generate this).\n"
    printf "\t--nexus-auth\t\t\tUsername and password (separated by '%s') of a Nexus user (if not given, script will prompt for this).\n" "$auth_separator"
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

    -v|--verbose) verbose="true";;

    -n|--namespace) namespace="$2"; shift;;
    -n=*|--namespace=*) namespace="${1#*=}";;

    -f|--values) values_file="$2"; shift;;
    -f=*|--values=*) values_file="${1#*=}";;

    -s|--serviceaccount) serviceaccount="$2"; shift;;
    -s=*|--serviceaccount=*) serviceaccount="${1#*=}";;

    --no-diff) diff="false";;

    --dry-run) dry_run="true";;

    --auth-separator) auth_separator="$2"; shift;;
    --auth-separator=*) auth_separator="${1#*=}";;

    --bitbucket-auth) bitbucket_auth="$2"; shift;;
    --bitbucket-auth=*) bitbucket_auth="${1#*=}";;

    --bitbucket-webhook-secret) bitbucket_webhook_secret="$2"; shift;;
    --bitbucket-webhook-secret=*) bitbucket_webhook_secret="${1#*=}";;

    --nexus-auth) nexus_auth="$2"; shift;;
    --nexus-auth=*) nexus_auth="${1#*=}";;

    --private-cert) private_cert="$2"; shift;;
    --private-cert=*) private_cert="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd "${script_dir}"

values_fileS=$(echo "$values_file" | tr "," "\n")
values_args=()
for valueFile in ${values_fileS}; do
    values_args+=(--values="${valueFile}")
done

if [ "${verbose}" == "true" ]; then
    set -x
fi

if [ -z "${namespace}" ]; then
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
    echo "${secretTemplate}" | sed "s/{{name}}/${secretName}/" | sed "s/{{username}}/${username//\//\\/}/" | sed "s/{{password}}/${password//\//\\/}/" | "${kubectl_bin}" -n "${namespace}" apply -f -
}

installSecret () {
    local secretName="$1"
    local secretTemplate="$2"
    local flagValue="$3"
    local usernamePrompt="$4"
    local passwordPrompt="$5"

    # Split flag value on first occurence of auth separator.
    local authUser="${flagValue%%"${auth_separator}"*}"
    local authPassword="${flagValue#*"${auth_separator}"}"

    # If the secret exists and the flag is present, update the secret.
    if "${kubectl_bin}" -n "${namespace}" get "secret/${secretName}" &> /dev/null; then
        # In case the secret was previously managed by Helm, we want to instruct Helm
        # to keep the resource during helm upgrade.
        "${kubectl_bin}" -n "${namespace}" annotate --overwrite secret "${secretName}" "helm.sh/resource-policy=keep"
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
        if "${kubectl_bin}" -n "${namespace}" get "secret/${secretName}" &> /dev/null; then
            echo "Re-creating secret ${secretName} ..."
            "${kubectl_bin}" -n "${namespace}" delete secret "${secretName}"
        else
            echo "Creating secret ${secretName} ..."
        fi
        "${kubectl_bin}" -n "${namespace}" create secret generic "${secretName}" \
            --from-file=tls.crt="${certFile}"
        rm private-cert.pem.tmp &>/dev/null || true
    fi
}

# Manage serviceaccount ...
if "${kubectl_bin}" -n "${namespace}" get serviceaccount/"${serviceaccount}" &> /dev/null; then
    echo "Serviceaccount exists already ..."
else
    echo "Creating serviceaccount ..."
    if [ "${dry_run}" == "true" ]; then
        echo "(skipping in dry-run)"
    else
        "${kubectl_bin}" -n "${namespace}" create serviceaccount "${serviceaccount}"

        "${kubectl_bin}" -n "${namespace}" \
            create rolebinding "${serviceaccount}-edit" \
            --clusterrole edit \
            --serviceaccount "${namespace}:${serviceaccount}"
    fi
fi

echo "Installing secrets ..."
if [ "${dry_run}" == "true" ]; then
    echo "(skipping in dry-run)"
else
    # Bitbucket username is not required as PAT alone is enough.
    installSecret "ods-bitbucket-auth" \
        "${basicAuthSecretTemplate}" \
        "${bitbucket_auth}" \
        "Please enter the username of Bitbucket user with write permission." \
        "Please enter a personal access token of this Bitbucket user (input will be hidden):"

    # Webhook secret is a special case, as we do not want the user to set it.
    # No prompts -> password will be auto-generated if not given.
    installSecret "ods-bitbucket-webhook" \
        "${opaqueSecretTemplate}" \
        "${bitbucket_webhook_secret}" \
        "" ""

    installSecret "ods-nexus-auth" \
        "${basicAuthSecretTemplate}" \
        "${nexus_auth}" \
        "Please enter the username of a Nexus user with write permission:" \
        "Please enter the password of this Nexus user (input will be hidden):"

    installTLSSecret "ods-private-cert" "${private_cert}"
fi

echo "Discovering Helm repository ..."
helm_repo_alias="ods-pipeline"
chart_name="ods-pipeline"
"${helm_bin}" repo add "${helm_repo_alias}" https://opendevstack.github.io/ods-pipeline
"${helm_bin}" repo update "${helm_repo_alias}"

echo "Installing Helm release ${release_name} ..."
if [ "${diff}" == "true" ]; then
    if "${helm_bin}" -n "${namespace}" \
            diff upgrade --install --detailed-exitcode --three-way-merge --normalize-manifests \
            "${values_args[@]}" \
            "${release_name}" "${helm_repo_alias}/${chart_name}"; then
        echo "Helm release already up-to-date."
    else
        if [ "${dry_run}" == "true" ]; then
            echo "(skipping in dry-run)"
        else
            "${helm_bin}" -n "${namespace}" \
                upgrade --install \
                "${values_args[@]}" \
                ${release_name} ${chart_dir}
        fi
    fi
else
    if [ "${dry_run}" == "true" ]; then
        echo "(skipping in dry-run)"
    else
        "${helm_bin}" -n "${namespace}" \
            upgrade --install \
            "${values_args[@]}" \
            "${release_name}" "${helm_repo_alias}/${chart_name}"
    fi
fi

echo "Adding Tekton annotation to ods-bitbucket-auth secret ..."
if [ "${dry_run}" == "true" ]; then
    echo "(skipping in dry-run)"
else
    bitbucketUrl=$("${kubectl_bin}" -n "${namespace}" get cm/ods-bitbucket -ojsonpath='{.data.url}')
    "${kubectl_bin}" -n "${namespace}" annotate --overwrite secret ods-bitbucket-auth "tekton.dev/git-0=${bitbucketUrl}"
fi

echo "Adding ods-bitbucket-auth secret to ${serviceaccount} serviceaccount ..."
if [ "${dry_run}" == "true" ]; then
    echo "(skipping in dry-run)"
else
    "${kubectl_bin}" -n "${namespace}" \
        patch sa "${serviceaccount}" \
        --type json \
        -p '[{"op": "add", "path": "/secrets", "value":[{"name": "ods-bitbucket-auth"}]}]'
fi
