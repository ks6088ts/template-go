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

## Quick start

```shell
go run . sandbox hello   # prints "sandbox hello called"
go run . version         # prints "version=dev, revision=unknown"
make ci-test             # the same checks CI runs
make ci-test-docker      # Dockerfile checks, image build, and smoke test
```

Command results go to stdout; errors go to stderr with exit status 1.

## Daily commands

| Command | Purpose |
|---|---|
| `make help` | list all targets |
| `make format` | format code with `gofmt` |
| `make test` | race-enabled tests with coverage |
| `make lint` | golangci-lint and actionlint (pinned, installed into `.tools/bin`) |
| `make build` | build `dist/template-go` with version metadata |
| `make ci-test` | format check, module and Go-version consistency, vet, lint, test, build |
| `make docker-build` / `make docker-run` | build and run the local image (`DOCKER_COMMAND=version`) |
| `make ci-test-docker` | `docker build --check`, build, non-root and version smoke test |
| `make release` | GoReleaser snapshot in `dist/` (never publishes) |
| `make vuln-check` | optional `govulncheck` (needs network) |

Versions are defined once: Go in `go.mod` (the Dockerfile build image must
match, which `make deps-check` enforces), development tools in `Makefile`.
Dependabot proposes grouped weekly updates for Go modules, Actions, Docker
images, and the dev container.

## Adapt this template

1. Replace the module path in `go.mod` and imports (`go mod edit -module ...`).
2. Rename the CLI in `cmd/root.go` and the output name (`OUTPUT`) in `Makefile`.
3. Update the badges and links at the top of this README.
4. Replace the `sandbox` commands, then run `make ci-test ci-test-docker`.

Workflows derive image names from `github.repository`; no workflow edits are
needed.

## CI/CD

| Workflow | Trigger | Does |
|---|---|---|
| `test.yaml` | push / PR to `main` | `make ci-test` on Ubuntu; vet and race tests on macOS and Windows |
| `docker.yaml` | push / PR to `main`, `v*` tags | `test`: `make ci-test-docker`; `publish` (tags only): push to GHCR and Docker Hub |
| `release.yaml` | `v*` tags | GoReleaser archives on the GitHub release |
| `copilot-setup-steps.yml` | manual / own changes | prepares the Copilot coding agent environment |

Settings (all optional; with none of them set, every workflow succeeds):

| Kind | Name | Used for |
|---|---|---|
| Variable | `DOCKERHUB_USERNAME` | Docker Hub namespace; Docker Hub is skipped when unset |
| Secret | `DOCKERHUB_TOKEN` | Docker Hub [access token](https://app.docker.com/settings/personal-access-tokens/create) with write access |

```shell
gh variable set DOCKERHUB_USERNAME --body "<user>"
gh secret set DOCKERHUB_TOKEN    # paste the token when prompted
gh secret delete DOCKERHUB_USERNAME   # if migrating from the old secret
```

GHCR uses `GITHUB_TOKEN` with `packages: write`. Images are
`ghcr.io/<owner>/<repo>` and `docker.io/<DOCKERHUB_USERNAME>/<repo>` (lowercase),
built for `linux/amd64` and `linux/arm64` and pushed once with the same digest,
an SBOM, and (public repositories only) a build provenance attestation.

## Release

```shell
make release-check RELEASE_TAG=v1.2.3   # validate the tag format
git tag v1.2.3 && git push origin v1.2.3
```

A `v*` tag runs `release.yaml` and the `publish` job. Images are tagged `1.2.3`
and `1.2`; `latest` moves only for stable versions (prereleases such as
`v1.2.3-rc.1` get only `1.2.3-rc.1`). Push tags yourself or with a personal
token: tags created with `GITHUB_TOKEN` do not trigger workflows.
