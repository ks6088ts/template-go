[![test](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml?query=branch%3Amain)
[![release](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ks6088ts/template-go)](https://goreportcard.com/report/github.com/ks6088ts/template-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/ks6088ts/template-go.svg)](https://pkg.go.dev/github.com/ks6088ts/template-go)

[![docker](https://github.com/ks6088ts/template-go/actions/workflows/docker.yaml/badge.svg?branch=main)](https://github.com/ks6088ts/template-go/actions/workflows/docker.yaml?query=branch%3Amain)
[![docker-release](https://github.com/ks6088ts/template-go/actions/workflows/docker-release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/docker-release.yaml)
[![ghcr-release](https://github.com/ks6088ts/template-go/actions/workflows/ghcr-release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/ghcr-release.yaml)

# template-go

A GitHub template repository for Go

## Prerequisites

- [Go 1.27+](https://go.dev/doc/install) (use the latest 1.27 patch release)
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/get-docker/) for container commands

## Development instructions

### Local development

Use Makefile to run the project locally.

```shell
# help
make

# install pinned lint tools into .tools/bin
make install-deps-dev

# check formatting, module consistency, and Go vet diagnostics
make format-check deps-check vet

# run static analysis (installs the pinned lint tools if needed)
make lint

# run tests with race detection and coverage
make test

# build applications without modifying go.mod or go.sum
make build

# run all required non-Docker checks
make ci-test

# create a local GoReleaser snapshot
make release
```

When changing dependencies, run `go mod tidy` explicitly and commit the resulting
`go.mod` and `go.sum`. The project uses the standard `testing` package for unit
tests. `make vuln-check` runs the pinned Go vulnerability scanner against the
current Go vulnerability database; it is an optional network-dependent check,
not part of the required CI gate.

### Docker development

```shell
# build docker image
make docker-build

# run docker container
make docker-run

# lint the Dockerfile, build the image, and run a smoke check
make ci-test-docker

# optionally scan the built image (requires access to the Docker socket)
make docker-scan
```

The Docker test workflow builds and smoke-tests the image; the scanner is run
separately when appropriate. The build stage uses Go 1.27 and the runtime image
uses a non-root distroless base.

## Deployment instructions

### Adapt this template

Before releasing a new project, replace the module path and imports, the
`-ldflags` package paths in `Makefile` and `.goreleaser.yaml`, the executable
and image names in the workflows, and the example CLI/configuration names.
Run `make ci-test` after making those changes.

Version tags must be valid semantic versions such as `v1.2.3` or
`v1.2.3-rc.1`. Docker Hub and GHCR publish the matching version tag for both,
but only stable versions update `latest`. The GoReleaser workflow publishes
release artifacts from these tags.

### Docker Hub

To publish the docker image to Docker Hub, you need to [create access token](https://app.docker.com/settings/personal-access-tokens/create) and set the following secrets in the repository settings.

```shell
gh secret set DOCKERHUB_USERNAME --body $DOCKERHUB_USERNAME
gh secret set DOCKERHUB_TOKEN --body $DOCKERHUB_TOKEN
```
