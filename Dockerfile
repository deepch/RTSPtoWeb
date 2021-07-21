FROM golang:alpine as builder

WORKDIR /build

ADD . .

RUN go build


FROM golang:alpine

RUN mkdir -p /opt/RTSPtoWeb

WORKDIR /opt/RTSPtoWeb

COPY --from=builder /build/RTSPtoWeb /opt/RTSPtoWeb

ADD ./web /opt/RTSPtoWeb/web

CMD /opt/RTSPtoWeb/RTSPtoWeb