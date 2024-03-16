GIT_REVISION ?= $(shell git rev-parse --short HEAD)
GIT_TAG ?= $(shell git describe --tags --abbrev=0 | sed -e s/v//g)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD ?= GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build
GOLANGCI_LINT_VERSION ?= 1.56.2
LDFLAGS ?= '-s -w \
	-X "github.com/ks6088ts/template-go/internal.Revision=$(GIT_REVISION)" \
	-X "github.com/ks6088ts/template-go/internal.Version=$(GIT_TAG)" \
'

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

.PHONY: install-deps-dev
install-deps-dev: ## install dependencies for development
	@# https://golangci-lint.run/welcome/install/#local-installation
	@which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION)

.PHONY: format-check
format-check: ## format check
	@# https://stackoverflow.com/a/67962664
	test -z "$$(gofmt -e -l -s .)"

.PHONY: format
format: ## format code
	gofmt -e -l -s -w .

.PHONY: lint
lint: ## lint
	golangci-lint run -v

.PHONY: test
test: ## run tests
	go test -cover -v ./...

.PHONY: build
build: ## build applications
	mkdir -p dist
	$(GOBUILD) -ldflags=$(LDFLAGS) -o dist/template-go_$(GIT_TAG)_$(GIT_REVISION)_$(GOOS)_$(GOARCH) .

.PHONY: ci-test
ci-test: install-deps-dev format-check lint test build ## run CI test
