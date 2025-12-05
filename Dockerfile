# syntax=docker/dockerfile:1

FROM --platform=${BUILDPLATFORM} golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /go/src/app

# Copy go mod/sum first to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

ARG TARGETOS TARGETARCH TARGETVARIANT

ENV CGO_ENABLED=0
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT#"v"} go build -a -o rtsp-to-web

FROM alpine:3.22

WORKDIR /app

# Copy binary
COPY --from=builder /go/src/app/rtsp-to-web /app/
# Copy web static assets
COPY --from=builder /go/src/app/web /app/web

# Create config dir
RUN mkdir -p /config
# Copy default config (will be overridden if mounted or updated via env vars logic)
COPY --from=builder /go/src/app/config.json /config

ENV GO111MODULE="on"
ENV GIN_MODE="release"

# Expose ports
EXPOSE 8083
EXPOSE 554

CMD ["./rtsp-to-web", "--config=/config/config.json"]
