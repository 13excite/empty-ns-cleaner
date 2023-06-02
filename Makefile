SHELL := /bin/bash

# constant variables
PROJECT_NAME 	= ns-cleaner
BINARY_NAME 	= ns-cleaner
GIT_COMMIT 		= $(shell git rev-parse HEAD)
BINARY_TAR_DIR 	= $(BINARY_NAME)-$(GIT_COMMIT)
BINARY_TAR_FILE	= $(BINARY_TAR_DIR).tar.gz
BUILD_VERSION 	= $(shell cat VERSION.txt)
BUILD_DATE 		= $(shell date -u '+%Y-%m-%d_%H:%M:%S')


# golangci-lint config
golangci_lint_version=latest
vols=-v `pwd`:/app -w /app
run_lint=docker run --rm $(vols) golangci/golangci-lint:$(golangci_lint_version)

# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint fmt

fmt:
	@gofmt -l -w $(SRC)

lint:
	@printf "$(OK_COLOR)==> Running golang-ci-linter via Docker$(NO_COLOR)\n"
	@$(run_lint) golangci-lint run --timeout=5m --verbose

build:
	@echo 'compiling binary...'
	# @cd cmd/ && GOARCH=amd64 GOOS=linux go build -ldflags "-X main.buildTimestamp=$(BUILD_DATE) -X main.gitHash=$(GIT_COMMIT) -X main.buildVersion=$(BUILD_VERSION)" -o ../$(BINARY_NAME)
	@cd cmd/ && GOARCH=arm64 GOOS=linux go build -ldflags "-X main.buildTimestamp=$(BUILD_DATE) -X main.gitHash=$(GIT_COMMIT) -X main.buildVersion=$(BUILD_VERSION)" -o ../$(BINARY_NAME)

