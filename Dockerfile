FROM golang:1.27-bookworm AS build

ARG GIT_REVISION="0000000"
ARG GIT_TAG="x.x.x"

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN make build OUTPUT=/go/bin/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /go/bin/app /
CMD ["/app"]
