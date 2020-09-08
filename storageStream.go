package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

//StreamMake check stream exist
func (obj *StorageST) StreamMake(val StreamST) StreamST {
	//make client's
	val.clients = make(map[string]ClientST)
	//make last ack
	val.ack = time.Now().Add(-255 * time.Hour)
	//make hls buffer
	val.hlsSegmentBuffer = make(map[int]Segment)
	//make signals buffer chain
	val.signals = make(chan int, 100)

	return val
}

//StreamExist check stream exist
func (obj *StorageST) StreamExist(key string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.ack = time.Now()
		obj.Streams[key] = tmp
		return ok
	}
	return false
}

//StreamRunAll run all stream go
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

//StreamRun one stream and lock
func (obj *StorageST) StreamRun(key string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if !tmp.runLock {
			tmp.runLock = true
			go StreamServerRunStreamDo(key)
		}
		obj.Streams[key] = tmp
	}
}

//StreamUnlock unlock status to no lock
func (obj *StorageST) StreamUnlock(key string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.runLock = false
		obj.Streams[key] = tmp
	}
}

//StreamControl get stream
func (obj *StorageST) StreamControl(key string) (*StreamST, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		return &tmp, nil
	}
	return nil, ErrorStreamNotFound
}

//List list all stream
func (obj *StorageST) List() map[string]StreamST {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	tmp := make(map[string]StreamST)
	for i, i2 := range obj.Streams {
		tmp[i] = i2
	}
	return tmp
}

//curl --header "Content-Type: application/json"   --request POST   --data '{"name": "test name 1","url": "rtsp://admin:123456@127.0.0.1:550/mpeg4", "on_demand": false,"debug": false}'   http://demo:demo@127.0.0.1:8083/stream/demo5/add
//StreamAdd add stream
func (obj *StorageST) StreamAdd(uuid string, val StreamST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[uuid]; ok {
		return ErrorStreamAlreadyExists
	}
	val = obj.StreamMake(val)
	if !val.OnDemand {
		val.runLock = true
		go StreamServerRunStreamDo(uuid)
	}
	obj.Streams[uuid] = val
	err := obj.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

//StreamEdit edit stream
func (obj *StorageST) StreamEdit(uuid string, val StreamST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		val = obj.StreamMake(val)
		if tmp.runLock {
			tmp.signals <- SignalStreamRestart
		} else if !val.OnDemand {
			val.runLock = true
			go StreamServerRunStreamDo(uuid)
		}
		obj.Streams[uuid] = val
		err := obj.SaveConfig()
		if err != nil {
			return err
		}
		return nil
	}
	return ErrorStreamNotFound
}

//StreamReload reload stream
func (obj *StorageST) StreamReload(uuid string) error {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if tmp.runLock {
			tmp.signals <- SignalStreamRestart
		}
		return nil
	}
	return ErrorStreamNotFound
}

//StreamDelete stream
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
	return ErrorStreamNotFound
}

//StreamInfo return stream info
func (obj *StorageST) StreamInfo(uuid string) (*StreamST, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		return &tmp, nil
	}
	return nil, ErrorStreamNotFound
}

//StreamCodecsUpdate update stream codec storage
func (obj *StorageST) StreamCodecsUpdate(key string, val []av.CodecData) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.codecs = val
		obj.Streams[key] = tmp
	}
}

//StreamCodecs get stream codec storage or wait
func (obj *StorageST) StreamCodecs(key string) ([]av.CodecData, error) {
	for i := 0; i < 100; i++ {
		obj.mutex.RLock()
		tmp, ok := obj.Streams[key]
		obj.mutex.RUnlock()
		if !ok {
			return nil, ErrorStreamNotFound
		}
		if tmp.codecs != nil {
			return tmp.codecs, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil, ErrorStreamNotFound
}

//Cast broadcast stream
func (obj *StorageST) Cast(key string, val *av.Packet) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if len(tmp.clients) > 0 {
			for ic, i2 := range tmp.clients {
				if len(i2.outgoingPacket) < 1000 {
					i2.outgoingPacket <- val
				} else if len(i2.signals) < 10 {
					//send stop signals to client
					i2.signals <- SignalStreamStop
					err := i2.socket.Close()
					if err != nil {
						loggingPrintln(ic, "close client error", err)
					}
				}
			}
			tmp.ack = time.Now()
			obj.Streams[key] = tmp
		}
	}
}

//StreamStatus change stream status
func (obj *StorageST) StreamStatus(key string, val int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		tmp.Status = val
		obj.Streams[key] = tmp
	}
}
