SHELL := /bin/bash

# constant variables
PROJECT_NAME 	= ns-cleaner
BINARY_NAME 	= ns-cleaner
GIT_COMMIT 		= $(shell git rev-parse HEAD)
BINARY_TAR_DIR 	= $(BINARY_NAME)-$(GIT_COMMIT)
BINARY_TAR_FILE	= $(BINARY_TAR_DIR).tar.gz
BUILD_VERSION 	= $(shell cat VERSION.txt)
BUILD_DATE 		= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint fmt

fmt:
	@gofmt -l -w $(SRC)

lint:
	@echo 'running linter...'
	@golangci-lint run ./...
#
# Build
#
build:
	@echo 'compiling binary...'
	@cd cmd/ && GOARCH=amd64 GOOS=linux go build -ldflags "-X main.buildTimestamp=$(BUILD_DATE) -X main.gitHash=$(GIT_COMMIT) -X main.buildVersion=$(BUILD_VERSION)" -o ../../$(BINARY_NAME)

