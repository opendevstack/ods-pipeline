SHELL = /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# Namespace to use (applied for some targets)
namespace =

## Check if system meets prerequisites.
check-system:
	cd scripts && ./check-system.sh
.PHONY: check-system

## Create KinD cluster with local registry.
create-kind-with-registry:
	cd scripts && ./kind-with-registry.sh
.PHONY: create-kind-with-registry

## Build and push images to local registry.
build-and-push-images:
		cd scripts && ./build-and-push-images.sh
.PHONY: build-and-push-images

## Install Tekton pipelines in kind cluster.
install-tekton-pipelines:
	cd scripts && ./install-tekton-pipelines.sh
.PHONY: install-tekton-pipelines

## Install resources in CD namespace via Helm.
install-cd-namespace:
ifeq ($(strip $(namespace)),)
	@echo "Argument 'namespace' is required, e.g. make install-cd-namespace namespace=foo-cd"
	@exit 1
endif
	cd scripts && ./install-cd-namespace-resources.sh -n $(namespace)
.PHONY: install-cd-namespace

## Apply ODS BuildConfig, ImageStream and ClusterTask manifests (OpenShift only)
install-ods-central:
ifeq ($(strip $(namespace)),)
	@echo "Argument 'namespace' is required, e.g. make install-ods-central namespace=ods"
	@exit 1
endif
	cd deploy/central && kubectl -n $(namespace) apply -k openshift
.PHONY: install-ods-central

## Start builds for each ODS BuildConfig (OpenShift only)
start-ods-central-builds: deploy-ods-central
	oc start-build ods-build-go
	oc start-build ods-buildah
	oc start-build ods-helm
	oc start-build ods-sonar
	oc start-build ods-start
	oc start-build ods-finish
.PHONY: start-ods-central-builds

## Apply ODS ClusterTask manifests in KinD
deploy-ods-tasks-kind:
	cd deploy/central && kubectl apply -k kind
.PHONY: deploy-ods-tasks-kind

## Run Bitbucket server (using timebomb license, in "kind" network).
run-bitbucket:
	cd scripts && ./run-bitbucket.sh
.PHONY: run-bitbucket

## Restart Bitbucket server (re-activating timebomb license).
restart-bitbucket:
	cd scripts && ./restart-bitbucket.sh
.PHONY: restart-bitbucket

## Run Nexus server (in "kind" network).
run-nexus:
	cd scripts && ./run-nexus.sh
.PHONY: run-nexus

## Run SonarQube server (in "kind" network).
run-sonarqube:
	cd scripts && ./run-sonarqube.sh
.PHONY: run-sonarqube

## Prepare local environment from scratch.
prepare-local-env: create-kind-with-registry build-and-push-images install-tekton-pipelines run-bitbucket run-nexus run-sonarqube deploy-ods-tasks-kind
.PHONY: prepare-local-env

## Stop local environment.
stop-local-env:
	cd scripts && ./stop-local-env.sh
.PHONY: stop-local-env

## Restart stopped local environment.
start-local-env:
	cd scripts && ./start-local-env.sh
.PHONY: start-local-env

## Render documentation for tasks.
docs:
	go run cmd/docs/main.go
.PHONY: docs

## Run testsuite.
test: test-internal test-pkg test-tasks
.PHONY: test

## Run testsuite of internal packages.
test-internal:
	go test -v ./internal/...
.PHONY: test-internal

## Run testsuite of public packages.
test-pkg:
	go test -v ./pkg/...
.PHONY: test-pkg

## Run testsuite of Tekton tasks.
test-tasks:
	go test -v -count=1 ./test/tasks/...
.PHONY: test-tasks

## Clear temporary workspaces created in testruns.
clear-tmp-workspaces:
	rm -rf test/testdata/workspaces/workspace-*
.PHONY: clear-tmp-workspaces

## Apply Tasks ODS manifests
deploy-ods-tasks:
	kubectl create -f deploy/tasks
.PHONY: deploy-ods-tasks

### HELP
### Based on https://gist.github.com/prwhite/8168133#gistcomment-2278355.
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:|^# .*/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-35s %s\n", helpCommand, helpMessage; \
		} else { \
			printf "\n"; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
.PHONY: help
