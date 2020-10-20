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

##Installation from binary

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

##Installation from source

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
## Stream mode

on_demand true  - receive video from source only has viewer

on_demand false - receive video from source any time

you can set mode use config "on_demand": true or "on_demand": false

## Configuration

###Options

####Server section's
```text
debug         - enable debug output
log_level     - log level
http_debug    - debug http api server
http_login    - http auth login
http_password - http auth password
http_port     - http server port
rtsp_port     - rtsp server port
```
####Stream section's
```text
name          - stream name
```
####Stream section's
```text
name          - channel name
url           - channel rtsp url
on_demand     - stream mode static (run any time) or ondaemand (run only has viewers)
debug         - enable debug output (RTSP client)
status        - default stream status

```

### example

```json
{
  "server": {
    "debug": true,
    "log_level": "info",
    "http_demo": true,
    "http_debug": false,
    "http_login": "demo",
    "http_password": "demo",
    "http_port": ":8083",
    "rtsp_port": ":5541"
  },
  "streams": {
    "demo1": {
      "name": "test video stream 1",
      "channels": {
        "0": {
          "name": "ch1",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        },
        "1": {
          "name": "ch2",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        }
      }
    },
    "demo2": {
      "name": "test video stream 2",
      "channels": {
        "0": {
          "name": "ch1",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        },
        "1": {
          "name": "ch2",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        }
      }
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

##API documentation
   
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
        "demo1": {
            "name": "test video",
            "channels": {
                "0": {
                    "name": "ch1",
                    "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                },
                "1": {
                    "name": "ch2",
                    "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                }
            }
        },
        "demo2": {
            "name": "test video",
            "channels": {
                "0": {
                    "name": "ch1",
                    "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                },
                "1": {
                    "name": "ch2",
                    "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                }
            }
        }
    }
}
```
### Stream Control
#### Stream Add
###### Query
```bash
POST /stream/{STREAM_ID}/add
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
              "name": "test video",
              "channels": {
                  "0": {
                      "name": "ch1",
                      "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                      "on_demand": true,
                      "debug": false,
                      "status": 0
                  },
                  "1": {
                      "name": "ch2",
                      "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                      "on_demand": true,
                      "debug": false,
                      "status": 0
                  }
              }
          }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/add
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
POST /stream/{STREAM_ID}/edit
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
            "name": "test video",
            "channels": {
                "0": {
                    "name": "ch1",
                    "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                },
                "1": {
                    "name": "ch2",
                    "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                }
            }
        }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/edit
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
GET /stream/{STREAM_ID}/reload
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/reload
```
###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```
#### Stream Channel Reload
###### Query
```bash
GET /stream/{STREAM_ID}/reload
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/reload
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
GET /stream/{STREAM_ID}/info
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/info
```

###### Response
```json
{
    "status": 1,
    "payload": {
        "name": "test video",
        "channels": {
            "0": {
                "name": "ch1",
                "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                "on_demand": true,
                "debug": false,
                "status": 0
            },
            "1": {
                "name": "ch2",
                "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                "on_demand": true,
                "debug": false,
                "status": 0
            }
        }
    }
}
```

#### Stream Delete
###### Query
```bash
GET /stream/{STREAM_ID}/delete
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/delete
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```

### Channel Control
#### Channel Add
###### Query
```bash
POST /stream/{STREAM_ID}/channel/{CHANNEL_ID}/add
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
                      "name": "ch4",
                      "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                      "on_demand": false,
                      "debug": false,
                      "status": 0
            }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/add
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```
#### Channel Edit
###### Query
```bash
POST /stream/{STREAM_ID}/channel/{CHANNEL_ID}/edit
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
                      "name": "ch4",
                      "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
                      "on_demand": true,
                      "debug": false,
                      "status": 0
            }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/edit
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```
#### Channel Reload
###### Query
```bash
GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/reload
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/reload
```
###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```

#### Channel Info
###### Query
```bash
GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/info
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/info
```

###### Response
```json
{
    "status": 1,
    "payload": {
        "name": "ch4",
        "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
        "on_demand": false,
        "debug": false,
        "status": 1
    }
}
```

#### Stream Codec
###### Query
```bash
GET /stream/{STREAM_ID}/{CHANNEL_ID}/codec
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/{CHANNEL_ID}/codec
```

###### Response
```json
{
    "status": 1,
    "payload": [
        {
            "Record": "AU0AFP/hABRnTQAUlahQfoQAAAMABAAAAwCiEAEABGjuPIA=",
            "RecordInfo": {
                "AVCProfileIndication": 77,
                "ProfileCompatibility": 0,
                "AVCLevelIndication": 20,
                "LengthSizeMinusOne": 3,
                "SPS": [
                    "Z00AFJWoUH6EAAADAAQAAAMAohA="
                ],
                "PPS": [
                    "aO48gA=="
                ]
            },
            "SPSInfo": {
                "ProfileIdc": 77,
                "LevelIdc": 20,
                "MbWidth": 20,
                "MbHeight": 15,
                "CropLeft": 0,
                "CropRight": 0,
                "CropTop": 0,
                "CropBottom": 0,
                "Width": 320,
                "Height": 240
            }
        }
    ]
}
```

#### Channel Delete
###### Query
```bash
GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/delete
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/delete
```

###### Response
```json
{
    "status": 1,
    "payload": "success"
}
```

#### Channel hls play
###### Query
```bash
GET /stream/{STREAM_ID}/hls/live/index.m3u8
curl http://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/hls/live/index.m3u8
```

###### Response
```bash
index.m3u8
```
```bash
ffplay http://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/hls/live/index.m3u8
```

#### Stream rtsp play
###### Query

```bash
ffplay -rtsp_transport tcp  rtsp://127.0.0.1/{STREAM_ID}/{CHANNEL_ID}
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