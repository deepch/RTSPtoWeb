# RTSPtoWeb share you ip camera to world!

RTSP Stream to WebBrowser MSE or WebRTC or HLS, full native! not use ffmpeg or gstreamer

## Table of Contents

- [Installation from binary](#Installation from binary)
- [Installation from source](#Installation from source)
- [Configuration](#Configuration)
- [API Documentation](#API documentation)
- [Limitations](#Limitations)
- [Performance](#Performance)
- [Authors](#Authors)
- [License](#license) 

## Installation from binary

To achieve it, after Darknet compilation (via make) execute following command:
```shell
GO111MODULE=on go get github.com/deepch/RTSPtoWeb
```
To achieve it, after Darknet compilation (via make) execute following command:
```shell
cd src/github.com/deepch/RTSPtoWeb
```
To achieve it, after Darknet compilation (via make) execute following command:
```shell
go run *.go
```
To access the web interface, you open a browser.
 ```shell
http://127.0.0.1:8083
 ```

## Installation from source

To achieve it, after Darknet compilation (via make) execute following command:
```shell
GO111MODULE=on go get github.com/deepch/RTSPtoWeb
```
To achieve it, after Darknet compilation (via make) execute following command:
```shell
cd src/github.com/deepch/RTSPtoWeb
```
To achieve it, after Darknet compilation (via make) execute following command:
```shell
go run *.go
```
To access the web interface, you open a browser.
 ```shell
http://127.0.0.1:8083
 ```

## Configuration

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
## API documentation

nope

## Limitations

Video Codecs Supported: H264 all profiles

Audio Codecs Supported: no

## Performance

```bash
CPU usage ≈0.2%-1% one (thread) core cpu intel core i7 per stream
```

## Authors

* **Andrey Semochkin** - *Initial work video* - [deepch](https://github.com/deepch)
* **Dmitry Vladikin** - *Initial work web UI* - [vdalex25](https://github.com/vdalex25)

See also the list of [contributors](https://github.com/deepch/RTSPtoWeb/contributors) who participated in this project.

## License

This project licensed. License - see the [LICENSE.md](LICENSE.md) file for details

[webrtc](https://github.com/pion/webrtc) follows license MIT [license](https://raw.githubusercontent.com/pion/webrtc/master/LICENSE).

[joy4](https://github.com/nareix/joy4) follows license MIT [license](https://raw.githubusercontent.com/nareix/joy4/master/LICENSE).

See also included packages.