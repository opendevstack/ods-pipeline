#!/usr/bin/env bash
set -ue

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ods_pipeline_dir=${script_dir%/*}
kind_deploy_path="/tmp/ods-pipeline/kind-deploy"
kind_values_dir="/tmp/ods-pipeline/kind-values"
helm_generated_values_file="${kind_deploy_path}/ods-pipeline/values.generated.yaml"

url_suffix="http"
bitbucket_auth="unavailable"
nexus_auth="unavailable:unavailable"
sonar_auth="unavailable"

if [ "$#" -gt 0 ]; then
    case $1 in
    --private-cert=*) url_suffix="https";
esac; fi

# Copy deploy path to tmp dir as the deploy path may be used through the Go package.
# The source directories of Go packages are placed into a non-writable location.
rm -rf "${kind_deploy_path}"
cp -r "${ods_pipeline_dir}/deploy" "${kind_deploy_path}"
chmod -R u+w "${kind_deploy_path}"

if [ -f "${kind_values_dir}/bitbucket-auth" ]; then
    bitbucket_auth=$(cat "${kind_values_dir}/bitbucket-auth")
fi
if [ -f "${kind_values_dir}/nexus-auth" ]; then
    nexus_auth=$(cat "${kind_values_dir}/nexus-auth")
fi
if [ -f "${kind_values_dir}/sonar-auth" ]; then
    sonar_auth=$(cat "${kind_values_dir}/sonar-auth")
fi

if [ ! -e "${helm_generated_values_file}" ]; then
    echo "setup:" > "${helm_generated_values_file}"
fi
if [ -f "${kind_values_dir}/bitbucket-${url_suffix}" ]; then
    bitbucket_url=$(cat "${kind_values_dir}/bitbucket-${url_suffix}")
    echo "  bitbucketUrl: '${bitbucket_url}'" >> "${helm_generated_values_file}"
fi
if [ -f "${kind_values_dir}/nexus-${url_suffix}" ]; then
    nexus_url=$(cat "${kind_values_dir}/nexus-${url_suffix}")
    echo "  nexusUrl: '${nexus_url}'" >> "${helm_generated_values_file}"
fi
if [ -f "${kind_values_dir}/sonar-${url_suffix}" ]; then
    sonar_url=$(cat "${kind_values_dir}/sonar-${url_suffix}")
    echo "  sonarUrl: '${sonar_url}'" >> "${helm_generated_values_file}"
fi

values_arg="${kind_deploy_path}/ods-pipeline/values.kind.yaml"
if [ "$(cat "${helm_generated_values_file}")" != "setup:" ]; then
    values_arg="${values_arg},${helm_generated_values_file}"
fi

cd "${kind_deploy_path}"
bash ./install.sh \
    --aqua-auth "unavailable:unavailable" \
    --aqua-scanner-url "none" \
    --bitbucket-auth "${bitbucket_auth}" \
    --nexus-auth "${nexus_auth}" \
    --sonar-auth "${sonar_auth}" \
    -f "${values_arg}" "$@"
