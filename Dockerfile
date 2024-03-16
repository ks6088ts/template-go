FROM golang:1.21 as build

WORKDIR /go/src/app
COPY . .

RUN make build OUTPUT=/go/bin/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /go/bin/app /
CMD ["/app"]
