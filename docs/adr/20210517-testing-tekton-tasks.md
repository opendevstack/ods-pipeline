# Testing Tekton Tasks

We will make use of [Kind](https://kind.sigs.k8s.io) to run a Kubernetes cluster. That cluster should have a container registry and Tekton installed.

With this setup, we will be able to run tests locally and on GitHub Actions.

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


## Further resources

* https://martinheinz.dev/blog/47
* https://github.com/MartinHeinz/tekton-kickstarter
* https://github.com/tektoncd/catalog/blob/main/CONTRIBUTING.md#end-to-end-testing
* https://github.com/tektoncd/catalog/tree/main/test
* Example test: https://github.com/tektoncd/catalog/tree/main/task/buildah/0.2/tests

Compared to the catalog tests we likely do not want to hack a sidecar into the task for external dependencies. I'd rather use the standard ConfigMap/Secret/Deployment setup we have in the actual OpenShift namespaces.
