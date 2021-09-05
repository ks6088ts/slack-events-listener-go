GOOS ?= $(shell go env GOOS)
GOPATH ?= $(shell go env GOPATH)
GOFILES ?= $(shell find . -name "*.go")

GOLANGCILINT_VERSION ?= 1.42.0
GORELEASER_VERSION ?= 0.174.2
COBRA_VERSION ?= 1.2.1
LDFLAGS ?= '-s -w \
	-X "github.com/ks6088ts/template-go/internal.Revision=$(shell git rev-parse --short HEAD)" \
	-X "github.com/ks6088ts/template-go/internal.Version=$(shell git describe --tags $$(git rev-list --tags --max-count=1))" \
'
BIN_NAME ?= template-go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

# Prerequisites
ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: make install-deps-lint")
endif
ifeq (, $(shell which cobra))
$(warning "could not find cobra in $(PATH), run: make install-deps-dev")
endif
ifeq (, $(shell which goreleaser))
$(warning "could not find goreleaser in $(PATH), run: make install-deps-release")
endif

# GNU Make version >= 3.81
.PHONY: install-deps-lint
install-deps-lint: ## install dependencies for lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v$(GOLANGCILINT_VERSION)

.PHONY: install-deps-dev
install-deps-dev: ## install dependencies for development
	go install github.com/spf13/cobra/cobra@v$(COBRA_VERSION)

.PHONY: install-deps-release
install-deps-release: ## install dependencies for release
	go install github.com/goreleaser/goreleaser@v$(GORELEASER_VERSION)

.PHONY: lint
lint: ## lint
	golangci-lint run -v

.PHONY: format
format: ## format codes
	gofmt -s -w $(GOFILES)

.PHONY: test
test: ## run tests
	go test -cover -v ./...

.PHONY: build
build: ## build applications
	go build -ldflags=$(LDFLAGS) -o dist/$(BIN_NAME) ./cli

.PHONY: ci-test
ci-test: lint test build ## run ci tests
