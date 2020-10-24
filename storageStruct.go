package main

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/sirupsen/logrus"
)

var Storage = NewStreamCore()

//Default stream  type
const (
	MSE = iota
	WEBRTC
	RTSP
)

//Default stream status type
const (
	OFFLINE = iota
	ONLINE
)

//Default stream errors
var (
	Success                         = "success"
	ErrorStreamNotFound             = errors.New("stream not found")
	ErrorStreamAlreadyExists        = errors.New("stream already exists")
	ErrorStreamChannelAlreadyExists = errors.New("stream channel already exists")
	ErrorStreamNotHLSSegments       = errors.New("stream hls not ts seq found")
	ErrorStreamNoVideo              = errors.New("stream no video")
	ErrorStreamNoClients            = errors.New("stream no clients")
	ErrorStreamRestart              = errors.New("stream restart")
	ErrorStreamStopCoreSignal       = errors.New("stream stop core signal")
	ErrorStreamStopRTSPSignal       = errors.New("stream stop rtsp signal")
	ErrorStreamChannelNotFound      = errors.New("stream channel not found")
	ErrorStreamChannelCodecNotFound = errors.New("stream channel codec not ready, possible stream offline")
)

//StorageST main storage struct
type StorageST struct {
	mutex   sync.RWMutex
	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

//ServerST server storage section
type ServerST struct {
	Debug        bool         `json:"debug"`
	LogLevel     logrus.Level `json:"log_level"`
	HTTPDemo     bool         `json:"http_demo"`
	HTTPDebug    bool         `json:"http_debug"`
	HTTPLogin    string       `json:"http_login"`
	HTTPPassword string       `json:"http_password"`
	HTTPPort     string       `json:"http_port"`
	RTSPPort     string       `json:"rtsp_port"`
}

//ServerST stream storage section
type StreamST struct {
	Name     string            `json:"name"`
	Channels map[int]ChannelST `json:"channels"`
}
type ChannelST struct {
	Name             string `json:"name"`
	URL              string `json:"url"`
	OnDemand         bool   `json:"on_demand"`
	Debug            bool   `json:"debug"`
	runLock          bool
	Status           int `json:"status"`
	codecs           []av.CodecData
	sdp              []byte
	signals          chan int
	hlsSegmentBuffer map[int]Segment
	hlsSegmentNumber int
	clients          map[string]ClientST
	ack              time.Time
}

//ClientST client storage section
type ClientST struct {
	mode              int
	signals           chan int
	outgoingAVPacket  chan *av.Packet
	outgoingRTPPacket chan *[]byte
	socket            net.Conn
}

//Segment HLS cache section
type Segment struct {
	dur  time.Duration
	data []*av.Packet
}
