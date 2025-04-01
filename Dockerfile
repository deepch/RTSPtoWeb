# syntax=docker/dockerfile:1

FROM --platform=${BUILDPLATFORM} golang:1.23-alpine3.20 AS builder

RUN apk add git

WORKDIR /go/src/app
COPY . .

ARG TARGETOS TARGETARCH TARGETVARIANT

# Initialize go module if it doesn't exist
RUN go mod init rtsp-to-web || true

ENV CGO_ENABLED=0
RUN go mod tidy \
    && go mod download \
    && GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT#"v"} go build -a -o rtsp-to-web

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache bash

COPY --from=builder /go/src/app/rtsp-to-web /app/
COPY --from=builder /go/src/app/web /app/web
COPY docker-entrypoint.sh /app/

RUN mkdir -p /config && \
    chmod +x /app/docker-entrypoint.sh

ENV GO111MODULE="on"
ENV GIN_MODE="release"

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["./rtsp-to-web", "--config=/config/config.json"]