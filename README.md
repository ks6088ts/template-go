[![test](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml?query=branch%3Amain)
[![release](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ks6088ts/template-go)](https://goreportcard.com/report/github.com/ks6088ts/template-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/ks6088ts/template-go.svg)](https://pkg.go.dev/github.com/ks6088ts/template-go)

# template-go

A GitHub template repository for Go

## Prerequisites

- [Go 1.21+](https://go.dev/doc/install)
- [GNU Make](https://www.gnu.org/software/make/)

## Development instructions

### Local development

Use Makefile to run the project locally.

```shell
# help
make

# install dependencies for development
make install-deps-dev

# run tests
make test

# build applications
make build

# run CI tests
make ci-test

# release applications
make release
```
