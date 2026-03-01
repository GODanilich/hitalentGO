FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

ENV CGO_ENABLED=0
ENV GOFLAGS=-mod=vendor

RUN go build -o /out/app ./cmd/api

RUN GOFLAGS= go build -mod=mod -o /out/goose github.com/pressly/goose/v3/cmd/goose

FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/app /app/app
COPY --from=builder /out/goose /usr/local/bin/goose
COPY migrations /app/migrations

EXPOSE 8080
CMD ["/app/app"]