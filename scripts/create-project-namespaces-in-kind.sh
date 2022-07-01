#!/usr/bin/env bash
# compare to ods-core/create-projects/create-projects.sh

#!/usr/bin/env bash
set -e

PROJECT_ID=""

function usage {
  printf "usage: %s [options]\n" "$0"
  printf "\t-h|--help\tPrints the usage\n"
  printf "\t-v|--verbose\tVerbose output\n"
  printf "\t-p|--project\tProject ID\n"
}

while [[ "$#" -gt 0 ]]; do case $1 in
  -v|--verbose) set -x;;

  -h|--help) usage; exit 0;;

  -p=*|--project=*) PROJECT_ID="${1#*=}";;
  -p|--project)     PROJECT_ID="$2"; shift;;

   *) echo "Unknown parameter passed: $1"; usage; exit 1;;
esac; shift; done

# check required parameters
if [ -z "${PROJECT_ID}" ]; then
  echo "PROJECT_ID is unset"; usage
  exit 1
else
	echo "PROJECT_ID=${PROJECT_ID}"
fi

echo "Create namespaces ${PROJECT_ID}-cd, ${PROJECT_ID}-dev and ${PROJECT_ID}-test"
kubectl create namespace "${PROJECT_ID}-cd"
kubectl create namespace "${PROJECT_ID}-dev"
kubectl create namespace "${PROJECT_ID}-test"
