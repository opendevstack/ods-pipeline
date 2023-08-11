#!/bin/bash
set -eu

HELM_PLUGIN_DIFF_VERSION=3.3.2
REPOSITORY=""
NAMESPACE=$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace)

while [ "$#" -gt 0 ]; do
    # shellcheck disable=SC2034
    case $1 in

    -v|--verbose) VERBOSE="true";;

    -r|--repository) REPOSITORY="$2"; shift;;
    -r=*|--repository=*) REPOSITORY="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

echo "Extending PATH to user-writable location ..."
mkdir -p bin
export PATH=/home/user/bin:$PATH
 
echo "Installing Helm plugins ..."
NO_DIFF_FLAG=""
if [ "$(helm plugin list | grep ^diff)" != "" ]; then
  echo "Plugin helm-diff is already installed."
else
  if command -v tar; then
    helm plugin install https://github.com/databus23/helm-diff --version "v${HELM_PLUGIN_DIFF_VERSION}"
  else
    NO_DIFF_FLAG="--no-diff"
  fi
fi

echo "Cloning Git repository ..."
if oc -n "${NAMESPACE}" get secrets/ods-bitbucket-auth &> /dev/null; then
  repoBase=$(oc -n "${NAMESPACE}" get configmaps/ods-bitbucket -o jsonpath='{.data.repoBase}')
  authToken=$(oc -n "${NAMESPACE}" get secrets/ods-bitbucket-auth -o jsonpath='{.data.password}' | base64 --decode)
  if [ -z "${REPOSITORY}" ]; then
    REPOSITORY="${repoBase}/${NAMESPACE%-cd}/${NAMESPACE}.git"
  fi
  repoName="${REPOSITORY##*/}"
  rm -rf "${repoName%.git}" || true
  git clone -c http.extraHeader="Authorization: Bearer ${authToken}" "${REPOSITORY}"
else
  echo 'No secret ods-bitbucket-auth found.'
  echo 'Most likely, there is no ODS Pipeline installation yet.'
  echo 'Clone the Git repository and run install.sh manually.'
  exit 1
fi

echo "Installing ..."
repoName="${REPOSITORY##*/}"
cd "${repoName%.git}/deploy"
./install.sh -n "${NAMESPACE}" ${NO_DIFF_FLAG}
