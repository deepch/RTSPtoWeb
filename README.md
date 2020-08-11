# RTSPtoWeb

RTSP Stream to WebBrowser MSE or WebRTC or HLS

full native! not use ffmpeg or gstreamer

## Team

Andrey - https://github.com/deepch video streaming developer

Dmitry - https://github.com/vdalex25 player's and web UI developer

## Installation
1.
```bash
GO111MODULE=on go get github.com/deepch/RTSPtoWeb
```
2.
```bash
cd src/github.com/deepch/RTSPtoWeb
```
3.
```bash
go run *.go
```
4.
```bash
open web browser http://127.0.0.1:8083
```

## Configuration

### Edit file config.json

format:

```bash
{
  "server": {
    "debug": true,
    "http_demo": true,
    "http_debug": true,
    "http_login": "demo",
    "http_password": "demo",
    "http_port": ":8083"
  },
  "streams": {
    "demo1": {
      "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
      "on_demand": false,
      "debug": false,
    },
    "demo2": {
        "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
      "on_demand": false,
      "debug": false,
    }
  }
}
```

## Limitations

Video Codecs Supported: H264 all profiles

Audio Codecs Supported: no

## Test

CPU usage 0.2% one core cpu intel core i7 / stream
