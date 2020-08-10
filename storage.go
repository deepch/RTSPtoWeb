package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
)

var Storage = NewStreamCore()

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
	HTTPDemo     bool   `json:"http_demo"`
	HTTPDebug    bool   `json:"http_debug"`
	HTTPLogin    string `json:"http_login"`
	HTTPPassword string `json:"http_password"`
	HTTPPort     string `json:"http_port"`
}

type StreamST struct {
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

func NewStreamCore() *StorageST {
	var tmp StorageST
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	for i, i2 := range tmp.Streams {
		i2.clients = make(map[string]ClientST)
		i2.hlsSegmentBuffer = make(map[int]Segment)
		tmp.Streams[i] = i2
	}
	return &tmp
}

/*
 Server Sections
*/

func (obj *StorageST) ServerHTTPDebug() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDebug
}

func (obj *StorageST) ServerHTTPDemo() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDemo
}

func (obj *StorageST) ServerHTTPLogin() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPLogin
}

func (obj *StorageST) ServerHTTPPassword() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPassword
}

func (obj *StorageST) ServerHTTPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPort
}

/*
 Stream Sections
*/

func (obj *StorageST) StreamExist(key string) bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	_, ok := obj.Streams[key]
	return ok
}

func (obj *StorageST) StreamRunAll() {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	for k, v := range obj.Streams {
		if !v.OnDemand {
			v.runLock = true
			go StreamServerRunStreamDo(k)
			obj.Streams[k] = v
		}
	}
}

func (obj *StorageST) StreamRun(key string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok && !tmp.runLock {
		tmp.runLock = true
		obj.Streams[key] = tmp
		go StreamServerRunStreamDo(key)
	}
}

func (obj *StorageST) StreamUnlock(key string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.runLock = false
		obj.Streams[key] = tmp
	}
}

func (obj *StorageST) StreamControl(key string) (*StreamST, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		return &tmp, nil
	}
	return nil, ErrorNotFound
}

func (obj *StorageST) List() map[string]StreamST {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	tmp := make(map[string]StreamST)
	for i, i2 := range obj.Streams {
		tmp[i] = i2
	}
	return tmp
}

func (obj *StorageST) StreamAdd(uuid string, val StreamST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[uuid]; ok {
		return ErrorFound
	}
	obj.Streams[uuid] = val
	err := obj.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (obj *StorageST) StreamEdit(uuid string, val StreamST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if tmp.runLock {
			tmp.signals <- SignalStreamStop
		}
		obj.Streams[uuid] = val
		err := obj.SaveConfig()
		if err != nil {
			return err
		}
		return nil
	}
	return ErrorNotFound
}

func (obj *StorageST) StreamReload(uuid string) error {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if tmp.runLock {
			tmp.signals <- SignalStreamRestart
		}
		return nil
	}
	return ErrorNotFound
}

func (obj *StorageST) StreamDelete(uuid string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if tmp.runLock {
			tmp.signals <- SignalStreamStop
		}
		delete(obj.Streams, uuid)
		err := obj.SaveConfig()
		if err != nil {
			return err
		}
		return nil
	}
	return ErrorNotFound
}

func (obj *StorageST) StreamInfo(uuid string) (*StreamST, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		return &tmp, nil
	}
	return nil, ErrorNotFound
}

func (obj *StorageST) StreamCodecsUpdate(key string, val []av.CodecData) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.codecs = val
		obj.Streams[key] = tmp
	}
}

func (obj *StorageST) StreamCodecs(key string) ([]av.CodecData, error) {
	for i := 0; i < 100; i++ {
		obj.mutex.RLock()
		tmp, ok := obj.Streams[key]
		obj.mutex.RUnlock()
		if !ok {
			return nil, ErrorNotFound
		}
		if tmp.codecs != nil {
			return tmp.codecs, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil, ErrorCodecNotFound
}
func (obj *StorageST) StreamHLSAdd(uuid string, val []*av.Packet, dur time.Duration) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		tmp.hlsSegmentNumber++
		tmp.hlsSegmentBuffer[tmp.hlsSegmentNumber] = Segment{data: val, dur: dur}
		if len(tmp.hlsSegmentBuffer) >= 6 {
			delete(tmp.hlsSegmentBuffer, tmp.hlsSegmentNumber-6-1)
		}
		obj.Streams[uuid] = tmp
	}
}

func (obj *StorageST) StreamHLSm3u8(uuid string) (string, int, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		var out string
		out += "#EXTM3U\r\n#EXT-X-TARGETDURATION:4\r\n#EXT-X-VERSION:4\r\n#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(tmp.hlsSegmentNumber) + "\r\n"
		var keys []int
		for k := range tmp.hlsSegmentBuffer {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		var count int
		for _, i := range keys {
			count++
			out += "#EXTINF:" + strconv.FormatFloat(tmp.hlsSegmentBuffer[i].dur.Seconds(), 'f', 1, 64) + ",\r\nsegment/" + strconv.Itoa(i) + "/file.ts\r\n"

		}
		return out, count, nil
	}
	return "", 0, ErrorNotFound
}

//ready
//StreamHLSTS send hls segment buffer to clients
func (obj *StorageST) StreamHLSTS(key string, seq int) ([]*av.Packet, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[key].hlsSegmentBuffer[seq]; ok {
		return tmp.data, nil
	}
	return nil, ErrorNotFound
}

//ready
//Cast broadcast stream
func (obj *StorageST) Cast(key string, val *av.Packet) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[key]; ok {
		for _, i2 := range tmp.clients {
			if len(i2.outgoingPacket) < 1000 {
				i2.outgoingPacket <- val
			} else {
				//send stop signals to client
				i2.signals <- SignalStreamStop
				i2.socket.Close()
			}
		}
	}
}

//ready
//StreamStatus change stream status
func (obj *StorageST) StreamStatus(key string, val int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.Status = val
		obj.Streams[key] = tmp
	}
}

/*
 Client Sections
*/

//ready
//ClientAdd Add New Client to Translations
func (obj *StorageST) ClientAdd(uuid string) (string, chan *av.Packet, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[uuid]; !ok {
		return "", nil, ErrorNotFound
	}
	//Generate UUID client
	cid := pseudoUUID()
	ch := make(chan *av.Packet, 2000)
	obj.Streams[uuid].clients[cid] = ClientST{outgoingPacket: ch}
	return cid, ch, nil

}

//ready
//ClientDelete Delete Client
func (obj *StorageST) ClientDelete(uuid string, cid string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[uuid]; ok {
		delete(obj.Streams[uuid].clients, cid)
	}
}

//ready
//ClientDelete Delete Client
func (obj *StorageST) SaveConfig() error {
	res, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.json", res, 0644)
	if err != nil {
		return err
	}
	return nil
}
