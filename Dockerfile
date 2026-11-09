# check=error=true
# Compile on the builder's native platform instead of emulating the target CPU.
# Keep the Go patch equal to go.mod; make deps-check enforces this.
FROM --platform=$BUILDPLATFORM golang:1.27.3-trixie@sha256:af1bfc2614ebb1dbd31af514b78740fd678850e76498762b425ecee4c5c6ea80 AS build

ARG TARGETOS
ARG TARGETARCH
ARG GIT_REVISION="unknown"
ARG GIT_TAG="dev"

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY main.go Makefile ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build GOOS="$TARGETOS" GOARCH="$TARGETARCH" OUTPUT=/go/bin/app

FROM gcr.io/distroless/static-debian13:nonroot@sha256:25ac8f0c3b2334a9eb05452266bd4697864a8c8ef0ef71e5e6480d7f92868329

COPY --from=build /go/bin/app /
USER nonroot:nonroot
# Pass docker run arguments to the application, such as the version subcommand.
ENTRYPOINT ["/app"]
