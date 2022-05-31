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

As part of the preconditions to run the tests, we shall:

- Create a KinD cluster with a registry. (needed to push images for buildah)
- Pull the images from DockerHub into the KinD registry (Tasks by default will use localhost:5001, make it easy later to override).
- Install Tekton Core components
- Apply the custom Tekton tasks under /deploy.
- Verify we can run Tekton tests.
- [LATER] Install SonarQube (plugins, auth token, etc.) and Nexus.
- [LATER] Apply any ConfigMap, Secrets, etc for SonarQube and Nexus.

* Decide whether testing should be made public or internal.

* User will provide an array of parameters and a task name and a local folder optionally in case they need it to be mounted as a workspace.

In Go, create a temporary folder within /test. Defer cleaning up. https://yourbasic.org/golang/temporary-file-directory/

* User will write a Go test that will create a TaskRun on the fly based on the input parameters and run it with the tekton pkg.

* We get back a TaskRun that can be queried to check the status of the run (success or fail).

## Further resources

* https://martinheinz.dev/blog/47
* https://github.com/MartinHeinz/tekton-kickstarter
* https://github.com/tektoncd/catalog/blob/main/CONTRIBUTING.md#end-to-end-testing
* https://github.com/tektoncd/catalog/tree/main/test
* Example test: https://github.com/tektoncd/catalog/tree/main/task/buildah/0.2/tests

Compared to the catalog tests we likely do not want to hack a sidecar into the task for external dependencies. I'd rather use the standard ConfigMap/Secret/Deployment setup we have in the actual OpenShift namespaces.
