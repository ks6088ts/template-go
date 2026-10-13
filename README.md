[![test](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/ks6088ts/template-go/actions/workflows/test.yaml?query=branch%3Amain)
[![release](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ks6088ts/template-go)](https://goreportcard.com/report/github.com/ks6088ts/template-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/ks6088ts/template-go.svg)](https://pkg.go.dev/github.com/ks6088ts/template-go)

[![docker](https://github.com/ks6088ts/template-go/actions/workflows/docker.yaml/badge.svg?branch=main)](https://github.com/ks6088ts/template-go/actions/workflows/docker.yaml?query=branch%3Amain)
[![docker-release](https://github.com/ks6088ts/template-go/actions/workflows/docker-release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/docker-release.yaml)
[![ghcr-release](https://github.com/ks6088ts/template-go/actions/workflows/ghcr-release.yaml/badge.svg)](https://github.com/ks6088ts/template-go/actions/workflows/ghcr-release.yaml)

# template-go

A small Go CLI template built with [Cobra](https://github.com/spf13/cobra).
It includes nested example commands, standard-library tests, cross-platform
release archives, and non-root container images. Cobra is the only direct
application dependency; development and release tools are installed separately.

## Prerequisites

- [Go 1.27.2+](https://go.dev/doc/install); keep the supported 1.27 patch up to date
- [GNU Make](https://www.gnu.org/software/make/) for the convenience targets
- [Docker](https://docs.docker.com/get-docker/) only for container commands

Go alone is sufficient to build, run, and test the CLI:

```shell
go run . sandbox
go run . sandbox hello
go run . version
go test -mod=readonly ./...
```

Command results go to stdout; errors go to stderr with exit status 1. The example
commands accept no positional arguments. Help and shell completion are provided
by Cobra. Bare Go builds report `version=dev, revision=unknown`; Make and
GoReleaser inject build metadata without changing the output format.

The unused Viper configuration scaffold, `--config`, and `--toggle` have been
removed. No configuration files or environment variables are read by the CLI.
Add configuration only when a real command needs it.

## Development instructions

### Local development

Use Makefile to run the project locally.

```shell
# help
make

# install pinned lint tools into version-specific paths under .tools/bin
make install-deps-dev

# check formatting, module consistency, and Go vet diagnostics
make format-check deps-check vet

# run static analysis (reuses installed tools of the pinned versions)
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

Tool versions are pinned in `Makefile`. Changing a version selects a new tool
directory; repeated checks reuse the installed binary. `TOOLS_DIR` may point to
an absolute shared cache directory. Application dependencies remain in `go.mod`,
not in the tool cache.

Prefer reviewing one dependency update at a time (`go get module@version`,
`go mod tidy`, then `make ci-test`) over the broad `make update` convenience
target. Dependabot tracks Go modules, Actions, and Dockerfile base images.
Keep `go.mod` and the Docker build image on the same supported Go patch.
The development container uses the 1.27 image line; Go's automatic toolchain
selection honors the minimum patch in `go.mod`.

CI runs the full non-Docker gate on Ubuntu and native Go vet/race tests on
macOS and Windows. Real-executable tests check output streams and exit codes;
release-tag tests also run when Make is available on Unix. Windows development
does not require Make. Docker checks are a separate workflow.

`make build` uses a release version only for a clean, exactly tagged commit.
Untagged or modified source reports `dev` and its commit revision; a source
archive without Git metadata reports `dev` / `unknown`. Explicit
`GIT_TAG` and `GIT_REVISION` values override discovery.

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
separately when appropriate. `make -j ci-test-docker` also waits for the build
before running the smoke check, which verifies the complete version output.
The build stage cross-compiles with Go 1.27.2 on the builder's native CPU;
the runtime image uses a non-root distroless base.

```shell
# build both supported container platforms without publishing
docker buildx build --platform linux/amd64,linux/arm64 \
  --output type=oci,dest=dist/container.tar .
```

Loading or running both platforms locally requires an appropriate image store
and native execution or emulation support. Building them does not require
emulating the Go compiler.

## Deployment instructions

### Adapt this template

Before releasing a new project, replace the module path and imports, the
CLI name in `cmd/root.go`, the output name in `Makefile`, image names in the
workflows and Makefile, and repository links in this README. Build metadata
package paths are derived from the module in both Make and GoReleaser.
Adapt or replace the sandbox commands and run `make ci-test`.

Release archives retain Linux (386, amd64, arm, arm64), macOS (amd64, arm64),
and Windows (386, amd64, arm64) support. Containers target Linux amd64/arm64.
Archive names include the version and platform; Windows archives use zip.
GoReleaser snapshots build all release targets without publishing.

Version tags must be valid semantic versions such as `v1.2.3` or
`v1.2.3-rc.1`. Docker Hub and GHCR publish the matching version tag for both,
but only stable versions update `latest`. The GoReleaser workflow publishes
release artifacts from these tags.

All publication workflows validate the tag and run formatting, module
consistency, vet, and tests against the tagged source before publishing:

```shell
make release-check RELEASE_TAG=v1.2.3-rc.1
```

Registry tags omit the `v` prefix and SemVer build metadata. GHCR publishes to
`ghcr.io/<repository-owner>/template-go` using `GITHUB_TOKEN` with package write
permission. Docker Hub requires the secrets below. If a publication channel
is not needed, remove its workflow (`release.yaml`, `docker-release.yaml`, or
`ghcr-release.yaml`) and corresponding badge; keep the test workflows.

### Docker Hub

To publish the docker image to Docker Hub, you need to [create access token](https://app.docker.com/settings/personal-access-tokens/create) and set the following secrets in the repository settings.

```shell
printf '%s' "$DOCKERHUB_USERNAME" | gh secret set DOCKERHUB_USERNAME
printf '%s' "$DOCKERHUB_TOKEN" | gh secret set DOCKERHUB_TOKEN
```
