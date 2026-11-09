# Development

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
3. Update the badges and links in `README.md`.
4. Replace the `sandbox` commands, then run `make ci-test ci-test-docker`.

Workflows derive image names from `github.repository`; no workflow edits are
needed.
