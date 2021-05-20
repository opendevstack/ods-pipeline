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
