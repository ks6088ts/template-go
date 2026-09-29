# Git
GIT_REVISION ?= $(shell git rev-parse --short HEAD)
GIT_TAG ?= $(shell git describe --tags --abbrev=0 --always | sed -e 's/^v//')

# Go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD ?= GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -mod=readonly
LDFLAGS ?= '-s -w \
	-X "github.com/ks6088ts/template-go/internal.Revision=$(GIT_REVISION)" \
	-X "github.com/ks6088ts/template-go/internal.Version=$(GIT_TAG)" \
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

# Misc
OUTPUT_DIR ?= dist
OUTPUT ?= $(OUTPUT_DIR)/template-go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

.PHONY: install-deps-dev
install-deps-dev: ## install dependencies for development
	@mkdir -p $(TOOLS_DIR)
	GOBIN=$(TOOLS_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)
	GOBIN=$(TOOLS_DIR) go install github.com/rhysd/actionlint/cmd/actionlint@v$(ACTIONLINT_VERSION)

.PHONY: install-release-tool
install-release-tool: ## install the release tool
	@mkdir -p $(TOOLS_DIR)
	GOBIN=$(TOOLS_DIR) go install github.com/goreleaser/goreleaser/v2@v$(GORELEASER_VERSION)

.PHONY: format-check
format-check: ## format check
	@# https://stackoverflow.com/a/67962664
	test -z "$$(gofmt -e -l -s .)"

.PHONY: format
format: ## format code
	gofmt -e -l -s -w .

.PHONY: deps-check
deps-check: ## check go.mod and go.sum without changing them
	go mod tidy -diff
	go mod verify

.PHONY: vet
vet: ## run Go's built-in analyzers
	go vet ./...

.PHONY: lint
lint: install-deps-dev ## lint
	$(TOOLS_DIR)/golangci-lint run
	$(TOOLS_DIR)/actionlint

.PHONY: test
test: ## run tests
	go test -race -cover ./...

.PHONY: build
build: ## build applications
	mkdir -p $(OUTPUT_DIR)
	$(GOBUILD) -ldflags=$(LDFLAGS) -trimpath -o $(OUTPUT) .

.PHONY: ci-test
ci-test: format-check deps-check vet lint test build ## run CI test

.PHONY: update
update: ## update
	@# https://stackoverflow.com/a/67202539/4457856
	go get -u ./...
	go mod tidy

.PHONY: release
release: install-release-tool ## release applications
	$(TOOLS_DIR)/goreleaser release --snapshot --clean

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
docker-smoke: ## check the container's version command
	docker run --rm $(DOCKER_REPO_NAME)/$(DOCKER_IMAGE_NAME):$(GIT_TAG) version

.PHONY: docker-lint
docker-lint: ## lint Dockerfile
	docker run --rm -i hadolint/hadolint:v$(HADOLINT_VERSION) < Dockerfile

.PHONY: docker-scan
docker-scan: ## scan Docker image
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:$(TRIVY_VERSION) image $(DOCKER_REPO_NAME)/$(DOCKER_IMAGE_NAME):$(GIT_TAG)

.PHONY: vuln-check
vuln-check: ## check Go packages for known vulnerabilities (optional)
	@mkdir -p $(TOOLS_DIR)
	GOBIN=$(TOOLS_DIR) go install golang.org/x/vuln/cmd/govulncheck@v$(GOVULNCHECK_VERSION)
	$(TOOLS_DIR)/govulncheck ./...

.PHONY: ci-test-docker
ci-test-docker: docker-lint docker-build docker-smoke ## run CI test for Docker
