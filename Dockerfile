# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/filegen

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-s -w -extldflags=-static -X 'main.Version=${APP_VERSION}' -X 'main.BuildDate=${BUILD_TIME}'" -o /go/bin/filegen /go/src/github.com/artarts36/filegen/cmd/main.go

######################################################

FROM alpine

RUN apk add tini git

COPY --from=builder /go/bin/filegen /go/bin/filegen

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="filegen"
LABEL org.opencontainers.image.description="console application for generate files of your templates"
LABEL org.opencontainers.image.url="https://github.com/artarts36/filegen"
LABEL org.opencontainers.image.source="https://github.com/artarts36/filegen"
LABEL org.opencontainers.image.vendor="ArtARTs36"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="MIT"

ENTRYPOINT ["/go/bin/filegen"]
