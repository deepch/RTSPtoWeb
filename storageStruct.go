package main

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
)

var Storage = NewStreamCore()

//Default stream status type
const (
	OFFLINE = iota
	ONLINE
)

//Default stream errors
var (
	Success                   = "success"
	ErrorStreamNotFound       = errors.New("stream not found")
	ErrorStreamAlreadyExists  = errors.New("stream already exists")
	ErrorStreamNotHLSSegments = errors.New("stream hls not ts seq found")
	ErrorStreamNoVideo        = errors.New("stream no video")
	ErrorStreamNoClients      = errors.New("stream no clients")
	ErrorStreamRestart        = errors.New("stream restart")
	ErrorStreamStopCoreSignal = errors.New("stream stop core signal")
	ErrorStreamStopRTSPSignal = errors.New("stream stop rtsp signal")
	ErrorChannelNotFound      = errors.New("channel not found")
)

//StorageST main storage struct
type StorageST struct {
	mutex   sync.RWMutex
	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

//ServerST server storage section
type ServerST struct {
	Debug        bool   `json:"debug"`
	HTTPDemo     bool   `json:"http_demo"`
	HTTPDebug    bool   `json:"http_debug"`
	HTTPLogin    string `json:"http_login"`
	HTTPPassword string `json:"http_password"`
	HTTPPort     string `json:"http_port"`
}

//ServerST stream storage section
type StreamST struct {
	Name     string            `json:"name"`
	Channels map[int]ChannelST `json:"channels"`
}
type ChannelST struct {
	URL              string `json:"url"`
	OnDemand         bool   `json:"on_demand"`
	Debug            bool   `json:"debug"`
	runLock          bool
	Status           int `json:"status"`
	codecs           []av.CodecData
	signals          chan int
	hlsSegmentBuffer map[int]Segment
	hlsSegmentNumber int
	clients          map[string]ClientST
	ack              time.Time
}

//ClientST client storage section
type ClientST struct {
	mode           int
	signals        chan int
	outgoingPacket chan *av.Packet
	socket         net.Conn
}

//Segment HLS cache section
type Segment struct {
	dur  time.Duration
	data []*av.Packet
}
