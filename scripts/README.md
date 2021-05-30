# Running Tekton tasks in KinD

```
kind create cluster
./run-tekton-task.sh
```


Run an E2E example:
```
./scripts/upload-dir-to-pvc.sh --pvc-name tektontest --source-directory test/tekton-e2e/input -v
./scripts/run-tekton-task.sh --task-file test/tekton-e2e/sample-task.yaml --task-run-file test/tekton-e2e/sample-taskrun.yaml --task-run-name sample-task-examplerun -v
./scripts/download-dir-from-pvc.sh --pvc-name tektontest --target-directory test/tekton-e2e/output -v
```

This will place a file "out.txt" into test/tekton-e2e/output.


# Running supporting services

Start container on kind network:
```
docker run -d -p "8081:8081" --net kind -e HTTP_PROXY="" -e HTTPS_PROXY="" -e NO_PROXY="" --name nexustest sonatype/nexus3:3.27.0
```

Now networking is possible via:
```
curl nexustest.kind:8081
```

Bitbucket licenses: https://developer.atlassian.com/platform/marketplace/timebomb-licenses-for-testing-server-apps/.
