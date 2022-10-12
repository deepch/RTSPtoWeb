# supervisor

https://codesahara.com/blog/how-to-deploy-golang-with-supervisor/

```
apt-get install supervisor
sudo service supervisor reload
supervisorctl status
```

```
[program:RTSPtoWeb]
directory=/home/rocoders/RTSPtoWeb
command=/home/rocoders/RTSPtoWeb/bin/RTSPtoWeb
autostart=true
autorestart=true
stderr_logfile=/var/log/RTSPtoWeb.err
stdout_logfile=/var/log/RTSPtoWeb.log

```

## Installation

### Installation from source

1. Download source
   ```bash
   $ git clone https://github.com/deepch/RTSPtoWeb
   ```
1. CD to Directory
   ```bash
    $ cd RTSPtoWeb/
   ```
1. Test Run
   ```bash
    $ GO111MODULE=on go run *.go
   ```
1. Open Browser
   ```bash
   open web browser http://127.0.0.1:8083 work chrome, safari, firefox
   ```

## Configuration

### Server settings

Install go lang > v1.3

```
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt-get update
sudo apt-get install golang-go
```

proxy pass web UI with authentication

```
location /streams/ {
    auth_basic           "Administrator’s Area";
    auth_basic_user_file /etc/apache2/.htpasswd;
    proxy_pass http://172.16.20.242:8083/;
}

location /static/ {
    proxy_pass http://172.16.20.242:8083/static/;
}

```

### App settings

```text
debug           - enable debug output
log_level       - log level (trace, debug, info, warning, error, fatal, or panic)

http_demo       - serve static files
http_debug      - debug http api server
http_login      - http auth login
http_password   - http auth password
http_port       - http server port
http_dir        - path to serve static files from
ice_servers     - array of servers to use for STUN/TURN
ice_username    - username to use for STUN/TURN
ice_credential  - credential to use for STUN/TURN
webrtc_port_min - minimum WebRTC port to use (UDP)
webrtc_port_max - maximum WebRTC port to use (UDP)

https
https_auto_tls
https_auto_tls_name
https_cert
https_key
https_port

rtsp_port       - rtsp server port
```

### Stream settings

```text
name            - stream name
```

### Channel settings

```text
name            - channel name
url             - channel rtsp url
on_demand       - stream mode static (run any time) or ondemand (run only has viewers)
debug           - enable debug output (RTSP client)
audio           - enable audio
status          - default stream status
```

#### Authorization play video

1 - enable config

```text
"token": {
"enable": true,
"backend": "http://127.0.0.1/file.php"
}
```

2 - try

```text
rtsp://127.0.0.1:5541/demo/0?token=you_key
```

file.php need response json

```text
   status: "1" or "0"
```

#### RTSP pull modes

- **on demand** (on_demand=true) - only pull video from the source when there's a viewer
- **static** (on_demand=false) - pull video from the source constantly

### Example config.json

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
    "ice_servers": ["stun:stun.l.google.com:19302"],
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
          "audio": true,
          "status": 0
        },
        "1": {
          "name": "ch2",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "audio": true,
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
  },
  "channel_defaults": {
    "on_demand": true
  }
}
```

## Command-line

### Use help to show available args

```bash
./RTSPtoWeb --help
```

#### Response

```bash
Usage of ./RTSPtoWeb:
  -config string
        config patch (/etc/server/config.json or config.json) (default "config.json")
  -debug
        set debug mode (default true)
```

## API documentation

See the [API docs](/docs/api.md)

## Limitations

Video Codecs Supported: H264 all profiles

Audio Codecs Supported: no

## Performance

```bash
CPU usage ≈0.2%-1% one (thread) core cpu intel core i7 per stream
```

## Authors

- **Andrey Semochkin** - _Initial work video_ - [deepch](https://github.com/deepch)
- **Dmitriy Vladykin** - _Initial work web UI_ - [vdalex25](https://github.com/vdalex25)

See also the list of [contributors](https://github.com/deepch/RTSPtoWeb/contributors) who participated in this project.

## License

This project licensed. License - see the [LICENSE.md](LICENSE.md) file for details

[webrtc](https://github.com/pion/webrtc) follows license MIT [license](https://raw.githubusercontent.com/pion/webrtc/master/LICENSE).

[joy4](https://github.com/nareix/joy4) follows license MIT [license](https://raw.githubusercontent.com/nareix/joy4/master/LICENSE).

## Other Example

Examples of working with video on golang

- [RTSPtoWeb](https://github.com/deepch/RTSPtoWeb)
- [RTSPtoWebRTC](https://github.com/deepch/RTSPtoWebRTC)
- [RTSPtoWSMP4f](https://github.com/deepch/RTSPtoWSMP4f)
- [RTSPtoImage](https://github.com/deepch/RTSPtoImage)
- [RTSPtoHLS](https://github.com/deepch/RTSPtoHLS)
- [RTSPtoHLSLL](https://github.com/deepch/RTSPtoHLSLL)

[![paypal.me/AndreySemochkin](https://ionicabizau.github.io/badges/paypal.svg)](https://www.paypal.me/AndreySemochkin) - You can make one-time donations via PayPal. I'll probably buy a ~~coffee~~ tea. :tea:
