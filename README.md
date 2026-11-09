[![test](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml?query=branch%3Amain)
[![docker](https://github.com/ks6088ts/template-go/actions/workflows/docker.yaml/badge.svg?branch=main)](https://github.com/ks6088ts/template-go/actions/workflows/docker.yaml?query=branch%3Amain)
[![release](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ks6088ts/template-go.svg)](https://pkg.go.dev/github.com/ks6088ts/template-go)

# template-go

A small Go CLI template built with [Cobra](https://github.com/spf13/cobra), the
only direct application dependency. It ships standard-library tests,
cross-platform release archives, and a non-root distroless container image.

## Prerequisites

- [Go](https://go.dev/doc/install) matching the `go` line in `go.mod` (older Go
  downloads it automatically)
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/get-docker/) with BuildKit, only for container commands

## Documentation

- [Development](docs/development.md): quick start, Make targets, version management, adapting this template
- [CI/CD and release](docs/ci-cd.md): workflows, Docker Hub / GHCR settings, release steps
