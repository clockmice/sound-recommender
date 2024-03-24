FROM golang:1.22.1 as golang

WORKDIR /src

COPY . /src

RUN CGO_ENABLED=0 go test -mod=vendor ./...

RUN CGO_ENABLED=0 go build \
    -mod vendor \
    -installsuffix 'static' \
    -o /bin/app cmd/server/main.go

FROM gcr.io/distroless/base

COPY --from=golang /bin/app /application
CMD ["/application"]
