# Git
# Source archives have no Git metadata; untagged builds are development versions.
GIT_REVISION ?= $(shell git rev-parse --short HEAD 2>/dev/null || printf unknown)
GIT_TAG ?= $(shell if test -z "$$(git status --porcelain 2>/dev/null)"; then git describe --tags --exact-match --match 'v[0-9]*' HEAD 2>/dev/null | sed -e 's/^v//'; fi)
ifeq ($(strip $(GIT_TAG)),)
GIT_TAG := dev
endif

# Go
MODULE_PATH := $(shell go list -m)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD ?= GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -mod=readonly
LDFLAGS ?= '-s -w \
	-X "$(MODULE_PATH)/internal.Revision=$(GIT_REVISION)" \
	-X "$(MODULE_PATH)/internal.Version=$(GIT_TAG)" \
'

# Docker
DOCKER_REPO_NAME ?= ks6088ts
DOCKER_IMAGE_NAME ?= template-go
DOCKER_COMMAND ?=

# Tools
TOOLS_DIR ?= $(CURDIR)/.tools/bin
# https://github.com/golangci/golangci-lint/releases
GOLANGCI_LINT_VERSION ?= 2.14.0
# https://github.com/rhysd/actionlint/releases
ACTIONLINT_VERSION ?= 1.7.12
# https://github.com/goreleaser/goreleaser/releases
GORELEASER_VERSION ?= 2.18.2
# https://github.com/hadolint/hadolint/releases
HADOLINT_VERSION ?= 2.15.1
# https://github.com/aquasecurity/trivy/releases
TRIVY_VERSION ?= 0.74.0
# https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
GOVULNCHECK_VERSION ?= 1.8.0
# Version-specific paths make installs incremental without running a version probe.
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)/golangci-lint
ACTIONLINT := $(TOOLS_DIR)/actionlint-$(ACTIONLINT_VERSION)/actionlint
GORELEASER := $(TOOLS_DIR)/goreleaser-$(GORELEASER_VERSION)/goreleaser
GOVULNCHECK := $(TOOLS_DIR)/govulncheck-$(GOVULNCHECK_VERSION)/govulncheck

# Misc
OUTPUT_DIR ?= dist
OUTPUT ?= $(OUTPUT_DIR)/template-go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:

.PHONY: install-deps-dev
install-deps-dev: $(GOLANGCI_LINT) $(ACTIONLINT) ## install pinned development tools if missing

$(GOLANGCI_LINT):
	@mkdir -p "$(@D)"
	GOBIN="$(@D)" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)

$(ACTIONLINT):
	@mkdir -p "$(@D)"
	GOBIN="$(@D)" go install github.com/rhysd/actionlint/cmd/actionlint@v$(ACTIONLINT_VERSION)

.PHONY: install-release-tool
install-release-tool: $(GORELEASER) ## install the pinned release tool if missing

$(GORELEASER):
	@mkdir -p "$(@D)"
	GOBIN="$(@D)" go install github.com/goreleaser/goreleaser/v2@v$(GORELEASER_VERSION)

$(GOVULNCHECK):
	@mkdir -p "$(@D)"
	GOBIN="$(@D)" go install golang.org/x/vuln/cmd/govulncheck@v$(GOVULNCHECK_VERSION)

.PHONY: format-check
format-check: ## format check
	@unformatted="$$(gofmt -e -l -s .)" && \
		if test -n "$$unformatted"; then printf 'Run make format for:\n%s\n' "$$unformatted"; exit 1; fi

.PHONY: format
format: ## format code
	gofmt -e -l -s -w .

.PHONY: deps-check
deps-check: ## check go.mod and go.sum without changing them
	go mod tidy -diff
	go mod verify

.PHONY: vet
vet: ## run Go's built-in analyzers
	go vet -mod=readonly ./...

.PHONY: lint
lint: install-deps-dev ## lint
	"$(GOLANGCI_LINT)" run
	"$(ACTIONLINT)"

.PHONY: test
test: ## run tests
	go test -mod=readonly -race -cover ./...

.PHONY: build
build: ## build applications
	mkdir -p "$(OUTPUT_DIR)"
	$(GOBUILD) -ldflags=$(LDFLAGS) -trimpath -o "$(OUTPUT)" .

.PHONY: ci-test
ci-test: format-check deps-check vet lint test build ## run CI test

.PHONY: update
update: ## update
	@# https://stackoverflow.com/a/67202539/4457856
	go get -u ./...
	go mod tidy

.PHONY: release
release: install-release-tool ## release applications
	"$(GORELEASER)" release --snapshot --clean

# Validate the environment value, never interpolate a release tag into shell code.
export RELEASE_TAG
.PHONY: release-check
release-check: ## validate RELEASE_TAG as a v-prefixed semantic version
	@case "$$RELEASE_TAG" in *[[:space:]]*) printf 'Invalid release tag: %s\n' "$$RELEASE_TAG" >&2; exit 1;; esac; \
		printf '%s\n' "$$RELEASE_TAG" | grep -Eq \
		'^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-(0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)(\.(0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$$' || \
		{ printf 'Invalid release tag: %s\n' "$$RELEASE_TAG" >&2; exit 1; }

# ---
# Docker
# ---

.PHONY: docker-build
docker-build: ## build Docker image
	docker build \
		-t $(DOCKER_REPO_NAME)/$(DOCKER_IMAGE_NAME):$(GIT_TAG) \
		--build-arg GIT_REVISION=$(GIT_REVISION) \
		--build-arg GIT_TAG=$(GIT_TAG) \
		.

.PHONY: docker-run
docker-run: ## run Docker container
	docker run --rm $(DOCKER_REPO_NAME)/$(DOCKER_IMAGE_NAME):$(GIT_TAG) $(DOCKER_COMMAND)

.PHONY: docker-smoke
docker-smoke: docker-build ## check the container's version output
	@test "$$(docker run --rm $(DOCKER_REPO_NAME)/$(DOCKER_IMAGE_NAME):$(GIT_TAG) version)" = \
		"version=$(GIT_TAG), revision=$(GIT_REVISION)"

.PHONY: docker-lint
docker-lint: ## lint Dockerfile
	docker run --rm -i hadolint/hadolint:v$(HADOLINT_VERSION) < Dockerfile

.PHONY: docker-scan
docker-scan: ## scan Docker image
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:$(TRIVY_VERSION) image $(DOCKER_REPO_NAME)/$(DOCKER_IMAGE_NAME):$(GIT_TAG)

.PHONY: vuln-check
vuln-check: $(GOVULNCHECK) ## check Go packages for known vulnerabilities (optional)
	"$(GOVULNCHECK)" ./...

.PHONY: ci-test-docker
ci-test-docker: docker-lint docker-build docker-smoke ## run CI test for Docker
