# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25.2
FROM golang:${GO_VERSION}-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /out/broker ./cmd/broker

FROM alpine:3.19
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 kafscale
USER 10001
WORKDIR /app

COPY --from=builder /out/broker /usr/local/bin/kafscale-broker

EXPOSE 19092 19093 19094
ENTRYPOINT ["/usr/local/bin/kafscale-broker"]
