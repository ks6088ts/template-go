# CI/CD and release

## Workflows

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
