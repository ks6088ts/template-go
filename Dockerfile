# Compile on the builder's native platform instead of emulating the target CPU.
FROM --platform=$BUILDPLATFORM golang:1.27.2-bookworm AS build

ARG TARGETOS
ARG TARGETARCH
ARG GIT_REVISION="unknown"
ARG GIT_TAG="dev"

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY Makefile ./

RUN make build GOOS=$TARGETOS GOARCH=$TARGETARCH OUTPUT=/go/bin/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /go/bin/app /
# Pass docker run arguments to the application, such as the version subcommand.
ENTRYPOINT ["/app"]
