#!/bin/bash
set -eu

# To carry out normal operations like running ODS Tekton Tasks,
# we need the ODS tasks images available in the KinD cluster.
REGISTRY="localhost:5000"
NAMESPACE="ods"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ODS_PIPELINE_DIR=${SCRIPT_DIR%/*}

SKIP_BUILD="false"
IMAGES=""
http_proxy="${http_proxy:-}"
https_proxy="${https_proxy:-}"

PLATFORM=""
# eg. --platform linux/amd64

while [ "$#" -gt 0 ]; do
    case $1 in

    -v|--verbose) set -x;;

    --skip-build) SKIP_BUILD="true";;

    -i|--image) IMAGES="$2"; shift;;
    -i=*|--image=*) IMAGES="${1#*=}";;

    -p|--platform) PLATFORM="$2"; shift;;
    -p=*|--platform=*) PLATFORM="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

cd "$ODS_PIPELINE_DIR"

build_and_push_image() {
    odsImage="ods-$image"
    if [ "${SKIP_BUILD}" != "true" ]; then
        echo "Building image $REGISTRY/$NAMESPACE/$odsImage..."
        # shellcheck disable=SC2086
        docker build $platform_arg \
            --build-arg http_proxy="$http_proxy" \
            --build-arg https_proxy="$https_proxy" \
            --build-arg HTTP_PROXY="$http_proxy" \
            --build-arg HTTPS_PROXY="$https_proxy" \
            -f build/package/Dockerfile."$image" -t $REGISTRY/$NAMESPACE/"$odsImage" .
    fi
    echo "Pushing image to $REGISTRY/$NAMESPACE/$odsImage ..."
    docker push "$REGISTRY/$NAMESPACE/$odsImage"
}

platform_arg=
if [ -n "$PLATFORM" ]; then
    platform_arg="--platform=${PLATFORM}"
fi

if [ -z "$IMAGES" ]; then
    for file in build/package/Dockerfile.*; do
        image=${file##*Dockerfile.}
        build_and_push_image
    done
else
    for image in $IMAGES; do
        build_and_push_image
    done
fi
