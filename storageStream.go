package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

/*
 Stream Sections
*/

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
	if tmp, ok := obj.Streams[key]; ok {
		if !tmp.runLock {
			tmp.runLock = true
			go StreamServerRunStreamDo(key)
		}
		obj.Streams[key] = tmp
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

//ready
//Cast broadcast stream
func (obj *StorageST) Cast(key string, val *av.Packet) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if len(tmp.clients) > 0 {
			for _, i2 := range tmp.clients {
				if len(i2.outgoingPacket) < 1000 {
					i2.outgoingPacket <- val
				} else {
					//send stop signals to client
					i2.signals <- SignalStreamStop
					i2.socket.Close()
				}
			}
			tmp.ack = time.Now()
			obj.Streams[key] = tmp
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
