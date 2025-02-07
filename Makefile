SHELL=/bin/bash -o pipefail
$(shell git config core.hooksPath ops/git-hooks)
PROJECT_NAME := $(shell basename $(CURDIR))
GO_MODULE := $(shell grep "^module " go.mod | awk '{print $$2}')
GO_PATH := $(shell go env GOPATH 2> /dev/null)
PATH := $(GO_PATH)/bin:$(PATH)

build:
	go generate
	go build -ldflags='-s -w' -o tmp/build .
	go install .

## Run the docs server for the project
docs-go:
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo "listening on http://127.0.0.1:6060/pkg/${GO_MODULE}"
	@godoc -http=127.0.0.1:6060
