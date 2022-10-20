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

### Web server settings (nginx)

The main server will be for PTZ commands, and handle all requests.

```
server {
    listen 80;

    server_name streaming.acsdns.co.uk;

    root /var/www/ptz/public_html;

    index index.html index.htm index.php;

    error_log  /var/log/nginx/streaming.error.log;
    access_log /var/log/nginx/streaming.access.log;


    location / {
        try_files $uri $uri/ /index.php$is_args$args;
    }

    location ~ \.php$ {
        include snippets/fastcgi-php.conf;
        fastcgi_pass unix:/var/run/php/php7.4-fpm.sock;
    }

    location /webui/ {
        auth_basic           "Administrator’s Area";
        auth_basic_user_file /etc/apache2/.htpasswd;
        # proxy_pass http://172.16.20.242:8083/;
        proxy_pass http://10.0.100.3:8083/;
    }
}
```

# MISC

If ufw is enabled, need to allow ports from config

```
ufw allow 61000:61000/upd
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

# supervisor

https://codesahara.com/blog/how-to-deploy-golang-with-supervisor/

```
apt-get install supervisor
sudo service supervisor reload
supervisorctl status
```

```
 /etc/supervisor/conf.d/rtsptoweb.conf

[program:RTSPtoWeb]
directory=/home/cywareadmin/RTSPtoWeb
command=/home/cywareadmin/RTSPtoWeb/bin/RTSPtoWeb
autostart=true
autorestart=true
stderr_logfile=/var/log/RTSPtoWeb.err
stdout_logfile=/var/log/RTSPtoWeb.log

```

## Performance

```bash
CPU usage ≈0.2%-1% one (thread) core cpu intel core i7 per stream
```
