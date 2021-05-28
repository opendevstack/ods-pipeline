#!/bin/bash

# To carry out normal operations like running ODS Tekton Tasks,
# we need the ODS tasks images available in the KinD cluster.
# Those images will be hosted in DockerHub so we need to pull them
# into our localhost to later on push them to the internal KinD registry.
REGISTRY=""
DOCKERHUB_REGISTRY="index.docker.io"
KIND_INTERNAL_REGISTRY="localhost:5000"
ODS_IMAGES=(
    "ods-build-go:latest" 
    "ods-sonar:latest"
)
LOAD_IMAGES_FROM_LOCALHOST=true

for i in "${ODS_IMAGES[@]}"
do
   if [ "$LOAD_IMAGES_FROM_LOCALHOST" = true ]; then
     REGISTRY=$KIND_INTERNAL_REGISTRY
   else
     REGISTRY=$DOCKERHUB_REGISTRY
   fi
   
   echo "Pushing image to $REGISTRY/ods/$i"
   docker push "$REGISTRY/ods/$i"
done