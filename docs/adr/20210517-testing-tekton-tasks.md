# Testing Tekton Tasks

Tekton Tasks should be easily testeable. Every Tekton Task shall have a `TaskRun` associated that will serve to test the functionality of the `Task`.


Either extend the existing `./scripts/run-tekton-task.sh` script or move it to Golang.

It shall do:
1. Create a PVC.
2. Start a Pod that mounts the PVC.
3. Copy local files into the Pod with `kubectl cp`.
4. Run the `TaskRun`
5. In the TaskRun mount the PVC
6. Start a Pod that mounts the PVC.
7. Copy files from Pod into local with `kubectl cp`.
8. Depending on the test, the expectations should be different (e.g. a file must be present in a specific dir, etc.)

## How test cases should be defined?

Assuming a `TaskRun` has been defined, you also need to reference a folder upload and write the actual test.

Use [this](https://github.com/opendevstack/tailor/blob/master/pkg/openshift/filter_test.go#L105) as a reference.

Simplest-case scenario:

Run a task that prints a hello world message and define a test in Go using Use [this](https://github.com/opendevstack/tailor/blob/master/pkg/openshift/filter_test.go#L105) approach as a reference.

* Decide whether testing should be made public or internal.