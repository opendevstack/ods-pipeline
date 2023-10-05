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

check-system: ## Check if system meets prerequisites for development.
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
	shellcheck scripts/*.sh build/images/scripts/* deploy/*.sh
.PHONY: lint-shell

docs: ## Render documentation for tasks.
	renderedStartTask=$(shell mktemp); \
	helm template ods-pipeline deploy/chart --show-only=templates/task-start.yaml > $$renderedStartTask; \
	go run github.com/opendevstack/ods-pipeline/cmd/taskdoc \
		-task $$renderedStartTask \
		-description build/docs/task-start.adoc \
		-destination docs/task-start.adoc; \
	rm $$renderedStartTask

	renderedFinishTask=$(shell mktemp); \
	helm template ods-pipeline deploy/chart --show-only=templates/task-finish.yaml > $$renderedFinishTask; \
	go run github.com/opendevstack/ods-pipeline/cmd/taskdoc \
		-task $$renderedFinishTask \
		-description build/docs/task-finish.adoc \
		-destination docs/task-finish.adoc; \
	rm $$renderedFinishTask
.PHONY: docs

##@ Building

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

test: test-cmd test-internal test-pkg test-e2e ## Run complete testsuite.
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

test-e2e: test-e2e-tasks test-e2e-pipelineruns ## Run testsuite of tasks and full pipeline run.
.PHONY: test-e2e

test-e2e-tasks: ## Run testsuite of tasks.
	go test -v -count=1 -timeout 20m -skip ^TestPipelineRun  ./test/e2e/...
.PHONY: test-e2e-tasks

test-e2e-pipelineruns: ## Run testsuite of full pipeline run.
	go test -v -count=1 -timeout 10m -run ^TestPipelineRun ./test/e2e/...
.PHONY: test-e2e-pipelineruns
