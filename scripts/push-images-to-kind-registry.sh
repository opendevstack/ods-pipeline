#!/bin/bash

# To carry out normal operations like running ODS Tekton Tasks,
# we need the ODS tasks images available in the KinD cluster.
REGISTRY="localhost:5000"
ODS_IMAGES=(
    "go-toolset:latest"
    "sonar:latest"
    "buildah:latest"
    "helm:latest"
)

for i in "${ODS_IMAGES[@]}"; do
   echo "Pushing image to $REGISTRY/ods/$i"
   docker push "$REGISTRY/ods/$i"
done
