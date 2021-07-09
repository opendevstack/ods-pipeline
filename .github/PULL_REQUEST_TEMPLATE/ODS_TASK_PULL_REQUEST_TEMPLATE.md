To create a technology-specific Task (e.g. python), the following files should be **added**:

- [ ] Dockerfile.python-toolset
- [ ] build-python.sh - Bash script to carry out the build, linting, testing operations. 
- [ ] bc-ods-build-python.yml - BuildConfig to generate the ods-build-python image.
- [ ] is-ods-build-python.yml - Create ImageStream resource in OpenShift.
- [ ] task-ods-build-python.yml - The Tekton Task.
- [ ] ods-build-python_test.go - A test file to test the behavior of the Tekton Task.
- [ ] python-flask-sample-app/** - Sample app.

The following files should be **updated**:

- [ ] .github/workflows/main.yml - Build ods python image and push it to the internal registry.
- [ ] `build-and-push-images.sh` - Extend it to include the new ods python image.
- [ ] Makefile - Extend it to start BuildConfig.