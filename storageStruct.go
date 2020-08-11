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

var (
	ErrorNotFound      = errors.New("stream not found")
	ErrorFound         = errors.New("stream already exists")
	ErrorCodecNotFound = errors.New("stream codec data not found")
)

type StorageST struct {
	mutex   sync.RWMutex
	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

type ServerST struct {
	Debug        bool   `json:"debug"`
	HTTPDemo     bool   `json:"http_demo"`
	HTTPDebug    bool   `json:"http_debug"`
	HTTPLogin    string `json:"http_login"`
	HTTPPassword string `json:"http_password"`
	HTTPPort     string `json:"http_port"`
}

type StreamST struct {
	Name             string `json:"name"`
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

type ClientST struct {
	mode           int
	signals        chan int
	outgoingPacket chan *av.Packet
	socket         net.Conn
}

type Segment struct {
	dur  time.Duration
	data []*av.Packet
}
