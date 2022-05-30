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
	shellcheck scripts/*.sh build/package/scripts/* deploy/*.sh
.PHONY: lint-shell

##@ Building

docs: ## Render documentation for tasks.
	go run cmd/docs/main.go
.PHONY: docs

build-artifact-download: build-artifact-download-linux build-artifact-download-darwin-amd64 build-artifact-download-darwin-arm64 build-artifact-download-windows ## Build artifact-download binary for each supported OS/arch.
.PHONY: build-artifact-download

build-artifact-download-linux: ## Build artifact-download Linux binary.
	cd scripts && ./build-artifact-download.sh --go-os=linux --go-arch=amd64
.PHONY: build-artifact-download-linux

build-artifact-download-darwin-amd64: ## Build artifact-download macOS binary (Intel).
	cd scripts && ./build-artifact-download.sh --go-os=darwin --go-arch=amd64
.PHONY: build-artifact-download-darwin-amd64

build-artifact-download-darwin-arm64: ## Build artifact-download macOS binary (Apple Silicon).
	cd scripts && ./build-artifact-download.sh --go-os=darwin --go-arch=arm64
.PHONY: build-artifact-download-darwin-arm64

build-artifact-download-windows: ## Build artifact-download Windows binary.
	cd scripts && ./build-artifact-download.sh --go-os=windows --go-arch=amd64
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
	chmod -R u+w test/testdata/workspaces/workspace-*
	rm -rf test/testdata/workspaces/workspace-*
.PHONY: clear-tmp-workspaces

##@ KinD (local development environment)

prepare-local-env: create-kind-with-registry build-and-push-images install-tekton-pipelines run-bitbucket run-nexus run-sonarqube ## Prepare local environment from scratch.
.PHONY: prepare-local-env

create-kind-with-registry: ## Create KinD cluster with local registry.
	cd scripts && ./kind-with-registry.sh
.PHONY: create-kind-with-registry

install-tekton-pipelines: ## Install Tekton pipelines in KinD cluster.
	cd scripts && ./install-tekton-pipelines.sh
.PHONY: install-tekton-pipelines

build-and-push-images: ## Build and push images to local registry.
		cd scripts && ./build-and-push-images.sh
.PHONY: build-and-push-images

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
.PHONY: recreate-kind-cluster

stop-local-env: ## Stop local environment.
	cd scripts && ./stop-local-env.sh
.PHONY: stop-local-env

start-local-env: ## Restart stopped local environment.
	cd scripts && ./start-local-env.sh
.PHONY: start-local-env

deploy: ## Install ODS pipeline resources in namespace.
ifeq ($(strip $(namespace)),)
	@echo "Argument 'namespace' is required, e.g. make deploy namespace=foo-cd"
	@exit 1
endif
	cd scripts && ./install-inside-kind.sh -n $(namespace)
.PHONY: deploy

##@ OpenShift

start-ods-builds: ## Start builds for each ODS BuildConfig
	oc start-build ods-buildah
	oc start-build ods-finish
	oc start-build ods-go-toolset
	oc start-build ods-gradle-toolset
	oc start-build ods-helm
	oc start-build ods-node16-npm-toolset
	oc start-build ods-pipeline-manager
	oc start-build ods-python-toolset
	oc start-build ods-sonar
	oc start-build ods-start
.PHONY: start-ods-builds
