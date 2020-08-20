# RTSPtoWeb share you ip camera to world!

RTSP Stream to WebBrowser MSE or WebRTC or HLS, full native! not use ffmpeg or gstreamer

## Table of Contents

- [Installation from binary](#Installation from binary)
- [Installation from source](#Installation from source)
- [Configuration](#Configuration)
- [Command-Line Arguments](#Command-Line Arguments)
- [API Documentation](#API documentation)
- [Limitations](#Limitations)
- [Performance](#Performance)
- [Authors](#Authors)
- [License](#license) 

## Installation from binary

Select the latest release
```shell
go to https://github.com/deepch/RTSPtoWeb/releases
```

Download the latest version [$version].tar.gz
```shell
wget https://github.com/deepch/RTSPtoWeb/archive/v0.0.1.tar.gz
```

Extract the archive
```shell
tar -xvzf v0.0.1.tar.gz
```

Change permission
```shell
chmod 777 RTSPtoWeb
```

Run the application
 ```shell
./RTSPtoWeb
 ```

## Installation from source

Enable the go module and get the source code
```shell
GO111MODULE=on go get github.com/deepch/RTSPtoWeb
```
Go to working directory

```shell
cd src/github.com/deepch/RTSPtoWeb
```

Run the source
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
    "demo": {
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

## Command-Line Arguments

######Use help show arg

```bash
./RTSPtoWeb --help
```

######Response 

```bash
Usage of ./RTSPtoWeb:
  -config string
        config patch (/etc/server/config.json or config.json) (default "config.json")
  -debug
        set debug mode (default true)
```

## API documentation
   
#### Streams List
###### Query
```bash
GET /streams

curl http://demo:demo@127.0.0.1:8083/streams
```

###### Response
```json
{
    "status": 1,
    "payload": {
        "demo": {
            "name": "test name 1",
            "url": "rtsp://admin:123456@127.0.0.1:550/mpeg4",
            "on_demand": true,
            "debug": false
        },
        "3demo": {
            "name": "test name 2",
             "url": "rtsp://admin:123456@127.0.0.1:551/mpeg4",
            "on_demand": false,
            "debug": false
        },
        "demo2": {
            "name": "test name 3",
             "url": "rtsp://admin:123456@127.0.0.1:552/mpeg4",
            "on_demand": true,
            "debug": false
        }
    }
}
```
#### Stream Add
###### Query
```bash
POST /stream/:uuid/add
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name": "test name 1","url": "rtsp://admin:123456@127.0.0.1:550/mpeg4", "on_demand": false,"debug": false}' \
  http://demo:demo@127.0.0.1:8083/stream/demo/add
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```
#### Stream Edit
###### Query
```bash
POST /stream/:uuid/edit
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name": "test name 1","url": "rtsp://admin:123456@127.0.0.1:550/mpeg4", "on_demand": false,"debug": false}' \
  http://demo:demo@127.0.0.1:8083/stream/demo/edit
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```
#### Stream Reload
###### Query
```bash
GET /stream/:uuid/reload
curl http://demo:demo@127.0.0.1:8083/stream/demo/reload
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```

#### Stream Info
###### Query
```bash
GET /stream/:uuid/info
curl http://demo:demo@127.0.0.1:8083/stream/demo/info
```

###### Response
```json
{
    "status": 1,
    "payload": {
        "name": "test name 1",
        "url": "rtsp://admin:123456@10.128.18.211/mpeg4",
        "on_demand": false,
        "debug": false,
        "status": 1
    }
}
```

#### Stream Codec
###### Query
```bash
GET /stream/:uuid/codec
curl http://demo:demo@127.0.0.1:8083/stream/demo/codec
```

###### Response
```json
{
    "status": 1,
    "payload": [
        {
            "Record": "AUKAKP/hACRnQoAo2gHgCJeWVIAAADwAAA4QMCAAHoSAAAiVRXvfC8IhGoABAARozjyA",
            "RecordInfo": {
                "AVCProfileIndication": 66,
                "ProfileCompatibility": 128,
                "AVCLevelIndication": 40,
                "LengthSizeMinusOne": 3,
                "SPS": [
                    "Z0KAKNoB4AiXllSAAAA8AAAOEDAgAB6EgAAIlUV73wvCIRqA"
                ],
                "PPS": [
                    "aM48gA=="
                ]
            },
            "SPSInfo": {
                "ProfileIdc": 66,
                "LevelIdc": 40,
                "MbWidth": 120,
                "MbHeight": 68,
                "CropLeft": 0,
                "CropRight": 0,
                "CropTop": 0,
                "CropBottom": 4,
                "Width": 1920,
                "Height": 1080
            }
        }
    ]
}
```

#### Stream Delete
###### Query
```bash
GET /stream/:uuid/delete
curl http://demo:demo@127.0.0.1:8083/stream/demo/delete
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```

#### Stream hls play
###### Query
```bash
GET /stream/:uuid/hls/live/index.m3u8
curl http://127.0.0.1:8083/stream/demo/hls/live/index.m3u8
```

###### Response
```bash
index.m3u8
```
```bash
ffplay http://127.0.0.1:8083/stream/demo/hls/live/index.m3u8
```

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