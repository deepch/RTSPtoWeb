package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

//ClientAdd Add New Client to Translations
func (obj *StorageST) ClientAdd(streamID string, channelID int) (string, chan *av.Packet, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return "", nil, ErrorStreamNotFound
	}
	//Generate UUID client
	cid, err := generateUUID()
	if err != nil {
		return "", nil, err
	}
	ch := make(chan *av.Packet, 2000)

	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return "", nil, ErrorStreamNotFound
	}

	channelTmp.clients[cid] = ClientST{outgoingPacket: ch, signals: make(chan int, 100)}
	channelTmp.ack = time.Now()
	streamTmp.Channels[channelID] = channelTmp
	obj.Streams[streamID] = streamTmp
	return cid, ch, nil

}

//ClientDelete Delete Client
func (obj *StorageST) ClientDelete(streamID string, cid string, channelID int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[streamID]; ok {
		delete(obj.Streams[streamID].Channels[channelID].clients, cid)
	}
}

//ClientHas check is client ext
func (obj *StorageST) ClientHas(streamID string, channelID int) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return false
	}
	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return false
	}
	if time.Now().Sub(channelTmp.ack).Seconds() > 30 {
		return false
	}
	return true
}
