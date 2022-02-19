# syntax=docker/dockerfile:1

FROM golang:1.18rc1-alpine3.15 AS builder

RUN apk add git

WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0
RUN go get \
    && go mod download \
    && go build -a -o rtsp-to-web

FROM alpine:3.15

WORKDIR /app

COPY --from=builder /go/src/app/rtsp-to-web /app/
COPY --from=builder /go/src/app/web /app/web
COPY --from=builder /go/src/app/config.json /app/

ENV GO111MODULE="on"
ENV GIN_MODE="release"

CMD ["./rtsp-to-web"]
