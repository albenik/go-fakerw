SHELL := /bin/bash

GO = vgo
MKFILE_DIR := $(abspath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SRC := $(shell find $(MKFILE_DIR) -type f -name '*.go' | grep -v /vendor/)

.PHONY: all
all: list test build

.PHONY: list
list:
	$(GO) list -m

.PHONY: test
test:
	@test -z $(shell gofmt -s -l $(SRC) | tee /dev/stderr) || echo "[WARN] Project has a formatting issues!"
	@$(GO) vet -v .
	@$(GO) test -v -race ./...

.PHONY: build
build:
	@$(GO) build -v ./...
