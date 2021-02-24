GOGET ?= go get -u -v
GOBUILD ?= go build
GOFMT ?= gofmt -s
GOTEST ?= go test

GOFILES ?= $(shell find . -name "*.go")
PKGS ?= $(shell go list ./...)
VERSION ?= $(shell git rev-parse --short HEAD)

# https://stackoverflow.com/questions/18136918/how-to-get-current-relative-directory-of-your-makefile
MAKEFILE_PATH ?= $(abspath $(lastword $(MAKEFILE_LIST)))
REPO_NAME ?= $(notdir $(patsubst %/,%,$(dir $(MAKEFILE_PATH))))
LDFLAGS ?= -ldflags="-s -w -X '$(shell echo $(MAKEFILE_PATH) | grep -o "github.com\/.*\/$(REPO_NAME)")/cmd/$(REPO_NAME)/cmd.version=$(VERSION)'"
OUT_DIR ?= outputs
BIN_PATH ?= $(OUT_DIR)/$(REPO_NAME)

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

.PHONY: rename
rename:
	find . \( -type d -name .git -prune -or -type d -name docs -prune -or -type f -name Makefile -prune \) -o -type f -print0 | xargs -0 sed -i 's/template-go/$(REPO_NAME)/g'

.PHONY: install-dev
install-dev: ## install dev dependencies
	# $(GOGET) golang.org/x/tools/gopls
	$(GOGET) golang.org/x/lint/golint

.PHONY: format
format: ## format codes
	$(GOFMT) -w $(GOFILES)

.PHONY: lint
lint: ## lint codes
	for PKG in $(PKGS); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: vet
vet: ## vet
	for PKG in $(PKGS); do go vet $$PKG || exit 1; done;

.PHONY: test
test: ## run tests
	$(GOTEST) -cover -v ./...

.PHONY: build
build: ## build an app
	$(GOBUILD) $(LDFLAGS) -o $(BIN_PATH) .

# https://qiita.com/dtan4/items/8c417b629b6b2033a541
.PHONY: crossbuild
crossbuild: ## cross build
	mkdir -p $(OUT_DIR)/release
	for os in darwin linux windows; do \
		for arch in amd64; do \
			make build GOBUILD="GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build" BIN_PATH="$(OUT_DIR)/$$os-$$arch/$(REPO_NAME)"; \
			zip $(OUT_DIR)/release/$(REPO_NAME)_$(VERSION)_$$os-$$arch.zip $(OUT_DIR)/$$os-$$arch/$(REPO_NAME); \
		done; \
	done

.PHONY: tidy
tidy: ## tidy
	go mod tidy

# FIXME: update commands as you like
.PHONY: ci
ci: install-dev install-cobra ## ci
	make cobra-init
	make cobra-add COBRA_CMD=hello
	# https://github.com/spf13/cobra/pull/1099
	# make lint # uncomment to activate lint 
	make vet
	make test
	make crossbuild
	$(OUT_DIR)/linux-amd64/$(REPO_NAME) --help
	$(OUT_DIR)/linux-amd64/$(REPO_NAME) hello --help

# ---
# Cobra: https://github.com/spf13/cobra
# ---

COBRA_CONFIG ?= .cobra.yml
COBRA_PKG_NAME ?= github.com/ks6088ts/$(REPO_NAME)
COBRA_CMD ?= hello
COBRA_PARENT_CMD ?= rootCmd

.PHONY: install-cobra
install-cobra: ## install cobra cli
	$(GOGET) github.com/spf13/cobra/cobra

.PHONY: cobra-init
cobra-init: ## initialize cobra cli
	cobra init \
		--pkg-name $(COBRA_PKG_NAME) \
		--config ./$(COBRA_CONFIG)

.PHONY: cobra-add
cobra-add: ## add cobra command
	cobra add $(COBRA_CMD) \
		--parent $(COBRA_PARENT_CMD) \
		--config ./$(COBRA_CONFIG)
