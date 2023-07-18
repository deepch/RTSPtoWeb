# RTSPtoWeb API

  * [Streams](#streams)
    * [List streams](#list-streams)
    * [Add a stream](#add-a-stream)
    * [Update a stream](#update-a-stream)
    * [Reload a stream](#reload-a-stream)
    * [Get stream info](#get-stream-info)
    * [Delete a stream](#delete-a-stream)
  * [Channels](#channels)
    * [Add a channel to a stream](#add-a-channel-to-a-stream)
    * [Update a stream channel](#update-a-stream-channel)
    * [Reload a stream channel](#reload-a-stream-channel)
    * [Get stream channel info](#get-stream-channel-info)
    * [Get stream channel codec](#get-stream-channel-codec)
    * [Delete a stream channel](#delete-a-stream-channel)
  * [Video endpoints](#video-endpoints)
    * [HLS](#hls)
    * [HLS-LL](#hls-ll)
    * [MSE](#mse)
    * [WebRTC](#webrtc)
    * [RTSP](#rtsp)

## Streams

### List streams

#### Request

`GET /streams`

```bash
curl http://demo:demo@127.0.0.1:8083/streams
```

#### Response

```json
{
    "status": 1,
    "payload": {
        "demo1": {
            "name": "test video",
            "channels": {
                "0": {
                    "name": "ch1",
                    "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                },
                "1": {
                    "name": "ch2",
                    "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
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
                    "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                },
                "1": {
                    "name": "ch2",
                    "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                    "on_demand": true,
                    "debug": false,
                    "status": 0
                }
            }
        }
    }
}
```

### Add a stream

#### Request

`POST /stream/{STREAM_ID}/add`

```bash
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{
              "name": "test video",
              "channels": {
                  "0": {
                      "name": "ch1",
                      "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                      "on_demand": true,
                      "debug": false,
                      "status": 0
                  },
                  "1": {
                      "name": "ch2",
                      "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                      "on_demand": true,
                      "debug": false,
                      "status": 0
                  }
              }
          }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/add
```

#### Response

```json
{
    "status": 1,
    "payload": "success"
}
```

### Update a stream

#### Request

`POST /stream/{STREAM_ID}/edit`

```bash
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{
              "name": "test video",
              "channels": {
                  "0": {
                      "name": "ch1",
                      "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                      "on_demand": true,
                      "debug": false,
                      "status": 0
                  },
                  "1": {
                      "name": "ch2",
                      "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                      "on_demand": true,
                      "debug": false,
                      "status": 0
                  }
              }
          }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/edit
```

#### Response

```json
{
    "status": 1,
    "payload": "success"
}
```

### Reload a stream

#### Request

`GET /stream/{STREAM_ID}/reload`

```bash
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/reload
```

#### Response

```json
{
    "status": 1,
    "payload": "success"
}
```

### Get stream info

#### Request

`GET /stream/{STREAM_ID}/info`

```bash
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/info
```

#### Response

```json
{
    "status": 1,
    "payload": {
        "name": "test video",
        "channels": {
            "0": {
                "name": "ch1",
                "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                "on_demand": true,
                "debug": false,
                "status": 0
            },
            "1": {
                "name": "ch2",
                "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
                "on_demand": true,
                "debug": false,
                "status": 0
            }
        }
    }
}
```

### Delete a stream

#### Request

`GET /stream/{STREAM_ID}/delete`

```bash
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/delete
```

#### Response

```json
{
    "status": 1,
    "payload": "success"
}
```

## Channels

### Add a channel to a stream

#### Request

`POST /stream/{STREAM_ID}/channel/{CHANNEL_ID}/add`

```bash
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{
              "name": "ch4",
              "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
              "on_demand": false,
              "debug": false,
              "status": 0
          }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/add
```

#### Response

```json
{
    "status": 1,
    "payload": "success"
}
```

### Update a stream channel

#### Request

`POST /stream/{STREAM_ID}/channel/{CHANNEL_ID}/edit`

```bash
curl \
  --header "Content-Type: application/json" \
  --request POST \
  --data '{
              "name": "ch4",
              "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
              "on_demand": true,
              "debug": false,
              "status": 0
          }' \
  http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/edit
```

#### Response

```json
{
    "status": 1,
    "payload": "success"
}
```

### Reload a stream channel

#### Request

`GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/reload`

```bash
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/reload
```

#### Response

```json
{
    "status": 1,
    "payload": "success"
}
```

### Get stream channel info

#### Request

`GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/info`

```bash
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/info
```

#### Response

```json
{
    "status": 1,
    "payload": {
        "name": "ch4",
        "url": "rtsp://admin:admin@{YOUR_CAMERA_IP}/uri",
        "on_demand": false,
        "debug": false,
        "status": 1
    }
}
```

### Get stream channel codec

#### Request
`GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/codec`

```bash
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/codec
```

#### Response
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

### Delete a stream channel

#### Request

`GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/delete`

```bash
curl http://demo:demo@127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/delete
```

#### Response
```json
{
    "status": 1,
    "payload": "success"
}
```

## Video endpoints

### HLS

`GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/hls/live/index.m3u8`

```bash
curl http://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/hls/live/index.m3u8
```

```bash
ffplay http://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/hls/live/index.m3u8
```

### HLS-LL

`GET /stream/{STREAM_ID}/channel/{CHANNEL_ID}/hlsll/live/index.m3u8`

```bash
curl http://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/hlsll/live/index.m3u8
```

```bash
ffplay http://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/hlsll/live/index.m3u8
```

### MSE

`/stream/{STREAM_ID}/channel/{CHANNEL_ID}/mse?uuid={STREAM_ID}&channel={CHANNEL_ID}`

```
ws://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/mse?uuid={STREAM_ID}&channel={CHANNEL_ID}
```

NOTE: Use `wss` for a secure connection.

### WebRTC

`/stream/{STREAM_ID}/channel/{CHANNEL_ID}/webrtc`

```
http://127.0.0.1:8083/stream/{STREAM_ID}/channel/{CHANNEL_ID}/webrtc
```

#### Request

The request is an HTTP `POST` with a FormData parameter `data` that is a base64 encoded SDP offer (e.g. `v=0...`) from a WebRTC client.

#### Response

The response is a base64 encoded SDP Answer.

### RTSP

`/{STREAM_ID}/{CHANNEL_ID}`

```
rtsp://127.0.0.1:{RTSP_PORT}/{STREAM_ID}/{CHANNEL_ID}
```

```bash
ffplay -rtsp_transport tcp rtsp://127.0.0.1/{STREAM_ID}/{CHANNEL_ID}
```
