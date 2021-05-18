#!/bin/bash
set -eu

KUBE_CONTEXT="--context kind-kind"
KUBECTL_BIN="kubectl $KUBE_CONTEXT"

if ! which kubectl &> /dev/null; then
    echo "kubectl is required"
fi

NAMESPACE="default"
PVC_NAME=""
SOURCE_DIRECTORY=""
POD_NAME="my-pod"

while [[ "$#" -gt 0 ]]; do
    case $1 in

    -v|--verbose) set -x;;

    # -h|--help) usage; exit 0;;

    # -i|--insecure) INSECURE="--insecure";;

    -p|--pvc-name) PVC_NAME="$2"; shift;;
    -p=*|--pvc-name=*) PVC_NAME="${1#*=}";;

    -s|--source-directory) SOURCE_DIRECTORY="$2"; shift;;
    -s=*|--source-directory=*) SOURCE_DIRECTORY="${1#*=}";;    

    -n|--namespace) NAMESPACE="$2"; shift;;
    -n=*|--namespace=*) NAMESPACE="${1#*=}";;

    *) echo_error "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

# Create PVC
echo "apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: $PVC_NAME
spec:
  # storageClassName: manual
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
" > pvc.yml

$KUBECTL_BIN -n $NAMESPACE apply -f pvc.yml

# Create Pod with PVC

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

$KUBECTL_BIN -n $NAMESPACE apply -f pod.yml

$KUBECTL_BIN -n $NAMESPACE cp -r $SOURCE_DIRECTORY $POD_NAME:/tmp/mydir