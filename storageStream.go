package main

import (
	"log"
	"time"

	"github.com/deepch/vdk/av"
)

//StreamMake check stream exist
func (obj *StorageST) StreamMake(val ChannelST) ChannelST {
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
func (obj *StorageST) StreamChannelExist(streamID string, channelID int) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := streamTmp.Channels[channelID]; ok {
			channelTmp.ack = time.Now()
			streamTmp.Channels[channelID] = channelTmp
			obj.Streams[streamID] = streamTmp
			return ok
		}
	}
	return false
}

//StreamRunAll run all stream go
func (obj *StorageST) StreamRunAll() {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	for k, v := range obj.Streams {
		for ks, vs := range v.Channels {
			if !vs.OnDemand {
				vs.runLock = true
				go StreamServerRunStreamDo(k, ks)
				v.Channels[ks] = vs
				obj.Streams[k] = v
			}
		}
	}
}

//StreamRun one stream and lock
func (obj *StorageST) StreamRun(streamID string, channelID int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := streamTmp.Channels[channelID]; ok {
			if !channelTmp.runLock {
				channelTmp.runLock = true
				streamTmp.Channels[channelID] = channelTmp
				obj.Streams[streamID] = streamTmp
				go StreamServerRunStreamDo(streamID, channelID)
			}
		}
	}
}

//StreamUnlock unlock status to no lock
func (obj *StorageST) StreamUnlock(streamID string, channelID int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := streamTmp.Channels[channelID]; ok {
			channelTmp.runLock = false
			streamTmp.Channels[channelID] = channelTmp
			obj.Streams[streamID] = streamTmp
		}
	}
}

//StreamControl get stream
func (obj *StorageST) StreamControl(key string, channelID int) (*ChannelST, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[key]; ok {
		if channelTmp, ok := streamTmp.Channels[channelID]; ok {
			return &channelTmp, nil
		}
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
	for i, i2 := range val.Channels {
		//val = obj.StreamMake(i, i2)
		log.Println("fix it make")
		if !i2.OnDemand {
			i2.runLock = true
			val.Channels[i] = i2
			go StreamServerRunStreamDo(uuid, i)
		}
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
		//val = obj.StreamMake(val)
		for i, i2 := range tmp.Channels {

			log.Println("fix it make")
			if channelTmp, ok := tmp.Channels[i]; ok {
				//copy global stream status need safe it
				i2.runLock = channelTmp.runLock
				//if stream running send restart stream
				if i2.runLock {
					channelTmp.signals <- SignalStreamRestart
				}
			}
			//if stream no running and no OnDemand
			if !i2.runLock && !i2.OnDemand {
				i2.runLock = true
				go StreamServerRunStreamDo(uuid, i)
			}
			//replace map
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
		for _, i2 := range tmp.Channels {
			if i2.runLock {
				i2.signals <- SignalStreamRestart
			}
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
		for _, i2 := range tmp.Channels {
			if i2.runLock {
				i2.signals <- SignalStreamStop
			}
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
func (obj *StorageST) StreamCodecsUpdate(streamID string, channelID int, val []av.CodecData) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.codecs = val
			tmp.Channels[channelID] = channelTmp
			obj.Streams[streamID] = tmp
		}
	}
}

//StreamCodecs get stream codec storage or wait
func (obj *StorageST) StreamCodecs(streamID string, channelID int) ([]av.CodecData, error) {
	for i := 0; i < 100; i++ {
		obj.mutex.RLock()
		tmp, ok := obj.Streams[streamID]
		obj.mutex.RUnlock()
		if !ok {
			return nil, ErrorStreamNotFound
		}
		channelTmp, ok := tmp.Channels[channelID]
		if !ok {
			return nil, ErrorChannelNotFound
		}

		if channelTmp.codecs != nil {
			return channelTmp.codecs, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil, ErrorStreamNotFound
}

//Cast broadcast stream
func (obj *StorageST) Cast(key string, channelID int, val *av.Packet) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			if len(channelTmp.clients) > 0 {
				for ic, i2 := range channelTmp.clients {
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
				channelTmp.ack = time.Now()
				tmp.Channels[channelID] = channelTmp
				obj.Streams[key] = tmp
			}
		}
	}
}

//StreamStatus change stream status
func (obj *StorageST) StreamStatus(key string, channelID int, val int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.Status = val
			tmp.Channels[channelID] = channelTmp
			obj.Streams[key] = tmp
		}
	}
}
