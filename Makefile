SHELL = /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# Namespace to use (applied for some targets)
namespace =

##@ General

# help target is based on https://github.com/operator-framework/operator-sdk/blob/master/release/Makefile.
.DEFAULT_GOAL := help
help: ## Show this help screen.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

check-system: ## Check if system meets prerequisites.
	cd scripts && ./check-system.sh
.PHONY: check-system

##@ Development

fmt: ## Run gofmt.
	gofmt -w .
.PHONY: fmt

lint: lint-go lint-shell ## Run all linters.
.PHONY: lint

lint-go: ## Run golangci-lint.
	golangci-lint run
.PHONY: lint-go

lint-shell: ## Run shellcheck.
	shellcheck scripts/*.sh build/package/scripts/* deploy/*/*.sh
.PHONY: lint-shell

##@ Building

sidecar-tasks: ## Render sidecar task variants.
	go run cmd/sidecar-tasks/main.go
.PHONY: sidecar-tasks

docs: sidecar-tasks ## Render documentation for tasks.
	go run cmd/docs/main.go
.PHONY: docs

build-artifact-download: build-artifact-download-linux build-artifact-download-darwin build-artifact-download-windows ## Build artifact-download binary for each supported OS/arch.
.PHONY: build-artifact-download

build-artifact-download-linux: ## Build artifact-download Linux binary.
	cd cmd/artifact-download && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-trimpath=$(CURDIR);$(shell go env GOPATH)" -o artifact-download-linux-amd64
.PHONY: build-artifact-download-linux

build-artifact-download-darwin: ## Build artifact-download macOS binary.
	cd cmd/artifact-download && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-trimpath=$(CURDIR);$(shell go env GOPATH)" -o artifact-download-darwin-amd64
.PHONY: build-artifact-download-darwin

build-artifact-download-windows: ## Build artifact-download Windows binary.
	cd cmd/artifact-download && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -gcflags "all=-trimpath=$(CURDIR);$(shell go env GOPATH)" -o artifact-download-windows-amd64.exe
.PHONY: build-artifact-download-windows

##@ Testing

test: test-cmd test-internal test-pkg test-tasks test-e2e ## Run complete testsuite.
.PHONY: test

test-cmd: ## Run testsuite of cmd packages.
	go test -cover ./cmd/...
.PHONY: test-cmd

test-internal: ## Run testsuite of internal packages.
	go test -cover ./internal/...
.PHONY: test-internal

test-pkg: ## Run testsuite of public packages.
	go test -cover ./pkg/...
.PHONY: test-pkg

test-tasks: ## Run testsuite of Tekton tasks.
	go test -v -count=1 -timeout $${ODS_TESTTIMEOUT:-30m} ./test/tasks/...
.PHONY: test-tasks

test-e2e: ## Run testsuite of end-to-end pipeline run.
	go test -v -count=1 -timeout $${ODS_TESTTIMEOUT:-10m} ./test/e2e/...
.PHONY: test-e2e

clear-tmp-workspaces: ## Clear temporary workspaces created in testruns.
	rm -rf test/testdata/workspaces/workspace-*
.PHONY: clear-tmp-workspaces

##@ KinD (local development environment)

prepare-local-env: create-kind-with-registry build-and-push-images install-tekton-pipelines run-bitbucket run-nexus run-sonarqube install-ods-tasks-kind ## Prepare local environment from scratch.
.PHONY: prepare-local-env

create-kind-with-registry: ## Create KinD cluster with local registry.
	cd scripts && ./kind-with-registry.sh
.PHONY: create-kind-with-registry

install-tekton-pipelines: ## Install Tekton pipelines in KinD cluster.
	cd scripts && ./install-tekton-pipelines.sh
.PHONY: install-tekton-pipelines

build-and-push-images: ## Build and push images to local registry.
		cd scripts && ./build-and-push-images.sh --image start
		cd scripts && ./build-and-push-images.sh --image finish
		cd scripts && ./build-and-push-images.sh --image buildah
		cd scripts && ./build-and-push-images.sh --image sonar
		cd scripts && ./build-and-push-images.sh --image webhook-interceptor
		cd scripts && ./build-and-push-images.sh --image go-toolset
		cd scripts && ./build-and-push-images.sh --image gradle-toolset
		cd scripts && ./build-and-push-images.sh --image python-toolset
		cd scripts && ./build-and-push-images.sh --image helm --platform linux/amd64
		cd scripts && ./build-and-push-images.sh --image node16-typescript-toolset --platform linux/amd64
.PHONY: build-and-push-images

install-ods-tasks-kind: ## KinD only! Apply ODS ClusterTask manifests in KinD
	cd scripts && ./install-ods-tasks-kind.sh
.PHONY: install-ods-tasks-kind

run-bitbucket: ## Run Bitbucket server (using timebomb license, in "kind" network).
	cd scripts && ./run-bitbucket.sh
.PHONY: run-bitbucket

restart-bitbucket: ## Restart Bitbucket server (re-activating timebomb license).
	cd scripts && ./restart-bitbucket.sh
.PHONY: restart-bitbucket

run-nexus: ## Run Nexus server (in "kind" network).
	cd scripts && ./run-nexus.sh
.PHONY: run-nexus

run-sonarqube: ## Run SonarQube server (in "kind" network).
	cd scripts && ./run-sonarqube.sh
.PHONY: run-sonarqube

recreate-kind-cluster: ## Recreate KinD cluster including Tekton tasks.
	cd scripts && ./kind-with-registry.sh --recreate
	cd scripts && ./install-tekton-pipelines.sh
	cd scripts && ./install-ods-tasks-kind.sh
.PHONY: recreate-kind-cluster

stop-local-env: ## Stop local environment.
	cd scripts && ./stop-local-env.sh
.PHONY: stop-local-env

start-local-env: ## Restart stopped local environment.
	cd scripts && ./start-local-env.sh
.PHONY: start-local-env

##@ OpenShift

install-ods-central: ## OpenShift only! Apply ODS BuildConfig, ImageStream and ClusterTask manifests
ifeq ($(strip $(namespace)),)
	@echo "Argument 'namespace' is required, e.g. make install-ods-central namespace=ods"
	@exit 1
endif
	cd scripts && ./install-ods-central-resources.sh -n $(namespace)
.PHONY: install-ods-central

start-ods-central-builds: ## OpenShift only! Start builds for each ODS BuildConfig
	oc start-build ods-buildah
	oc start-build ods-finish
	oc start-build ods-go-toolset
	oc start-build ods-gradle-toolset
	oc start-build ods-helm
	oc start-build ods-python-toolset
	oc start-build ods-sonar
	oc start-build ods-start
	oc start-build ods-node16-typescript-toolset
	oc start-build ods-webhook-interceptor
.PHONY: start-ods-central-builds

##@ User Installation

install-cd-namespace: ## Install resources in CD namespace via Helm.
ifeq ($(strip $(namespace)),)
	@echo "Argument 'namespace' is required, e.g. make install-cd-namespace namespace=foo-cd"
	@exit 1
endif
	cd scripts && ./install-cd-namespace-resources.sh -n $(namespace)
.PHONY: install-cd-namespace
