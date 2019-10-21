.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

SRC_PACKAGES=$(shell go list -mod=vendor ./... | grep -v "vendor")
BUILD?=$(shell git describe --always --dirty 2> /dev/null)
GOLINT:=$(shell command -v golint 2> /dev/null)
RICHGO=$(shell command -v richgo 2> /dev/null)
GOMETA_LINT=$(shell command -v golangci-lint 2> /dev/null)
GOLANGCI_LINT_VERSION=v1.20.0
GO111MODULE=on
SHELL=bash -o pipefail

ifeq ($(GOMETA_LINT),)
	GOMETA_LINT=$(shell command -v $(PWD)/bin/golangci-lint 2> /dev/null)
endif

ifeq ($(RICHGO),)
	GO_BINARY=go
else
	GO_BINARY=richgo
endif

ifeq ($(BUILD),)
	BUILD=dev
endif

ifdef CI_COMMIT_SHORT_SHA
	BUILD=$(CI_COMMIT_SHORT_SHA)
endif

all: setup build

ci: setup-common build-common

ensure-build-dir:
	mkdir -p out

build-deps: ## Install dependencies
	go get
	go mod tidy
	go mod vendor

update-deps: ## Update dependencies
	go get -u

compile: compile-app ## Compile library

compile-app: ensure-build-dir
	$(GO_BINARY) build -mod=vendor ...

build: fmt build-common ## Build the library

build-common: vet lint-all test compile

fmt:
	GOFLAGS="-mod=vendor" $(GO_BINARY) fmt $(SRC_PACKAGES)

vet:
	$(GO_BINARY) vet -mod=vendor $(SRC_PACKAGES)

setup-common:
ifeq ($(GOMETA_LINT),)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION)
endif
ifeq ($(GOLINT),)
	GO111MODULE=off $(GO_BINARY) get -u golang.org/x/lint/golint
endif

setup-richgo:
ifeq ($(RICHGO),)
	GO111MODULE=off $(GO_BINARY) get -u github.com/kyoh86/richgo
endif

setup: setup-richgo setup-common ensure-build-dir ## Setup environment

lint-all: lint setup-common
	$(GOMETA_LINT) run

target:
	 for number in 1 2 3 4 ; do \
            echo $$number ; \
     done

lint:
	golint $(SRC_PACKAGES)

test: ensure-build-dir ## Run tests
	ENVIRONMENT=test $(GO_BINARY) test -mod=vendor $(SRC_PACKAGES) -race -coverprofile ./out/coverage -short -v | grep -viE "start|no test files"

test-cover-html: ## Run tests with coverage
	mkdir -p ./out
	@echo "mode: count" > coverage-all.out
	$(foreach pkg, $(SRC_PACKAGES),\
	ENVIRONMENT=test $(GO_BINARY) test -mod=vendor -coverprofile=coverage.out -covermode=count $(pkg);\
	tail -n +2 coverage.out >> coverage-all.out;)
	$(GO_BINARY) tool cover -html=coverage-all.out -o out/coverage.html

generate-test-summary:
	ENVIRONMENT=test $(GO_BINARY) test -mod=vendor $(SRC_PACKAGES) -race -coverprofile ./out/coverage -short -v -json | grep -viE "start|no test files" | tee test-summary.json; \
    sed -i '' -E "s/^(.+\| {)/{/" test-summary.json; \
	passed=`cat test-summary.json | jq | rg '"Action": "pass"' | wc -l`; \
	skipped=`cat test-summary.json | jq | rg '"Action": "skip"' | wc -l`; \
	failed=`cat test-summary.json | jq | rg '"Action": "fail"' | wc -l`; \
	echo "Passed: $$passed | Failed: $$failed | Skipped: $$skipped"