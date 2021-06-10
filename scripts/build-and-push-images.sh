#!/bin/bash

# To carry out normal operations like running ODS Tekton Tasks,
# we need the ODS tasks images available in the KinD cluster.
REGISTRY="localhost:5000"
NAMESPACE="ods"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

SKIP_BUILD="false"
IMAGES="go-toolset sonar buildah helm start"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    --skip-build) SKIP_BUILD="true";;

    -i|--image) IMAGES="$2"; shift;;
    -i=*|--image=*) IMAGES="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd $ODS_PIPELINE_DIR

for image in $IMAGES; do
    if [ "${SKIP_BUILD}" != "true" ]; then
        echo "Building image $REGISTRY/$NAMESPACE/$image ..."
        docker build -f build/package/Dockerfile.$image -t $REGISTRY/$NAMESPACE/$image .
    fi
    echo "Pushing image to $REGISTRY/$NAMESPACE/$image ..."
    docker push "$REGISTRY/$NAMESPACE/$image"
done
