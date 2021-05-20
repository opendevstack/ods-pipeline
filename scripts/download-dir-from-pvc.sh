#!/bin/bash
set -eu

KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"

if ! which kubectl &> /dev/null; then
    echo "kubectl is required"
fi

NAMESPACE="default"
PVC_NAME=""
TARGET_DIRECTORY=""
POD_NAME="my-pod"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    # -h|--help) usage; exit 0;;

    -p|--pvc-name) PVC_NAME="$2"; shift;;
    -p=*|--pvc-name=*) PVC_NAME="${1#*=}";;

    -t|--target-directory) TARGET_DIRECTORY="$2"; shift;;
    -t=*|--target-directory=*) TARGET_DIRECTORY="${1#*=}";;    

    -n|--namespace) NAMESPACE="$2"; shift;;
    -n=*|--namespace=*) NAMESPACE="${1#*=}";;

    *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

if [ -z "${PVC_NAME}" ]; then
  echo "--pvc-name is required"
  exit 1
fi

if [ -z "${TARGET_DIRECTORY}" ]; then
  echo "--target-directory is required"
  exit 1
fi

KUBECTL_BIN_WITH_NS="$KUBECTL_BIN -n $NAMESPACE"

echo "Create Pod with PVC ..."
echo "apiVersion: v1
kind: Pod
metadata:
  labels:
    run: pod
  name: $POD_NAME
spec:
  containers:
    - image: alpine
      name: pod
      command: [\"/bin/sh\", \"-c\", \"--\"]
      args: [\"while true; do sleep 30; done;\"]
      resources: {}
      volumeMounts:
        - name: test-volume
          mountPath: \"/tmp/mydir\"
  volumes:
    - name: test-volume
      persistentVolumeClaim:
        claimName: $PVC_NAME
  dnsPolicy: ClusterFirst
  restartPolicy: Always
" > pod.yml
$KUBECTL_BIN_WITH_NS apply -f pod.yml

echo "Download from PVC into target directory ..."
$KUBECTL_BIN_WITH_NS cp $POD_NAME:/tmp/mydir $TARGET_DIRECTORY

echo "Delete pod ..."
$KUBECTL_BIN_WITH_NS delete pod/${POD_NAME}
